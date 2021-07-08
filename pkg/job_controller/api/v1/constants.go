package v1

const (
	// ReplicaIndexLabel represents the label key for the replica-index, e.g. the value is 0, 1, 2.. etc
	ReplicaIndexLabel = "replica-index"

	// ReplicaTypeLabel represents the label key for the replica-type, e.g. the value is ps , worker etc.
	ReplicaTypeLabel = "replica-type"

	// ReplicaNameLabel represents the label key for the replica-name, the value is replica name.
	ReplicaNameLabel = "replica-name"

	// GroupNameLabel represents the label key for group name, e.g. the value is kubeflow.org
	GroupNameLabel = "group-name"

	// JobNameLabel represents the label key for the job name, the value is job name
	JobNameLabel = "job-name"

	// JobRoleLabel represents the label key for the job role, e.g. the value is master
	JobRoleLabel = "job-role"
)

// Constant label/annotation keys for job configuration.
const (
	KubeDLPrefix = "kubedl.io"

	// AnnotationGitSyncConfig annotate git sync configurations.
	AnnotationGitSyncConfig = KubeDLPrefix + "/git-sync-config"
	// AnnotationTenancyInfo annotate tenancy information.
	AnnotationTenancyInfo = KubeDLPrefix + "/tenancy"
	// AnnotationNetworkMode annotate job network mode.
	AnnotationNetworkMode = KubeDLPrefix + "/network-mode"

	// AnnotationTensorBoardConfig annotate tensorboard configurations.
	AnnotationTensorBoardConfig = KubeDLPrefix + "/tensorboard-config"
	// ReplicaTypeTensorBoard is the type for TensorBoard.
	ReplicaTypeTensorBoard ReplicaType = "TensorBoard"
)

const (
	// LabelInferenceName represents the inference service name.
	LabelInferenceName = KubeDLPrefix + "/inference-name"
	// LabelPredictorName represents the predictor name of served model.
	LabelPredictorName = KubeDLPrefix + "/predictor-name"
	// LabelModelVersion represents the model version value for inference role.
	LabelModelVersion = KubeDLPrefix + "/model-version"
	// LabelCronName indicates the name of cron who created this job.
	LabelCronName = KubeDLPrefix + "/cron-name"
)

// NetworkMode defines network mode for intra job communicating.
type NetworkMode string

const (
	// HostNetworkMode indicates that replicas use host-network to communicate with each other.
	HostNetworkMode NetworkMode = "host"
)
