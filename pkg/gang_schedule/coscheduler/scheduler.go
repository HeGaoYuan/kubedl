/*
Copyright 2021 The Alibaba Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package coscheduler

import (
	"context"
	"github.com/alibaba/kubedl/pkg/util/k8sutil"
	"sigs.k8s.io/scheduler-plugins/pkg/apis/scheduling/v1alpha1"

	"github.com/alibaba/kubedl/apis"
	"github.com/alibaba/kubedl/pkg/gang_schedule"
	apiv1 "github.com/alibaba/kubedl/pkg/job_controller/api/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func init() {
	// Add to runtime scheme so that reflector of go-client will identify this CRD
	// controlled by scheduler.
	apis.AddToSchemes = append(apis.AddToSchemes, v1alpha1.AddToScheme)
}

func NewKubeCoscheduler(mgr controllerruntime.Manager) gang_schedule.GangScheduler {
	return &kubeCoscheduler{client: mgr.GetClient()}
}

var _ gang_schedule.GangScheduler = &kubeCoscheduler{}

type kubeCoscheduler struct {
	client client.Client
}

func (kbs *kubeCoscheduler) Name() string {
	return "kube-coscheduler"
}

func (kbs *kubeCoscheduler) CreateGang(job metav1.Object, replicas map[apiv1.ReplicaType]*apiv1.ReplicaSpec) (runtime.Object, error) {
	// Initialize pod group.
	podGroup := &v1alpha1.PodGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      job.GetName(),
			Namespace: job.GetNamespace(),
		},
		Spec: v1alpha1.PodGroupSpec{
			MinMember: k8sutil.GetTotalReplicas(replicas),
		},
	}
	err := kbs.client.Create(context.Background(), podGroup)
	return podGroup, err
}

func (kbs *kubeCoscheduler) BindPodToGang(obj metav1.Object, entity runtime.Object) error {
	podSpec := obj.(*v1.PodTemplateSpec)
	podGroup := entity.(*v1alpha1.PodGroup)
	// The newly-created pods should be submitted to target gang scheduler.
	if podSpec.Spec.SchedulerName == "" || podSpec.Spec.SchedulerName != kbs.Name() {
		podSpec.Spec.SchedulerName = kbs.Name()
	}
	if podSpec.Labels == nil {
		podSpec.Labels = map[string]string{}
	}
	// Coscheduling based on PodGroup CRD
	podSpec.Labels["pod-group.scheduling.sigs.k8s.io"] = podGroup.Name
	// Lightweight coscheduling based on back-to-back queue sorting
	podSpec.Labels["pod-group.scheduling.sigs.k8s.io/name"] = podGroup.Name
	podSpec.Labels["pod-group.scheduling.sigs.k8s.io/min-available"] = string(podGroup.Spec.MinMember)

	return nil
}

func (kbs *kubeCoscheduler) GetGang(name types.NamespacedName) (runtime.Object, error) {
	podGroup := &v1alpha1.PodGroup{}
	if err := kbs.client.Get(context.Background(), name, podGroup); err != nil {
		return nil, err
	}
	return podGroup, nil
}

func (kbs *kubeCoscheduler) DeleteGang(name types.NamespacedName) error {
	podGroup, err := kbs.GetGang(name)
	if err != nil {
		return err
	}
	err = kbs.client.Delete(context.Background(), podGroup)
	// Discard deleted pod group object.
	if err != nil && errors.IsNotFound(err) {
		return nil
	}
	return err
}
