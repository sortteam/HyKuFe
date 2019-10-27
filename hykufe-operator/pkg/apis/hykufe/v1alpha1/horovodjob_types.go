package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HorovodJobSpec defines the desired state of HorovodJob
// +k8s:openapi-gen=true
type HorovodJobSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html


	// VolumeSpec
	Volumes []VolumeSpec `json:"volumes,omitempty"`

	DataSources []DataSourceSpec `json:"dataSources,omitempty"`

	// Data share mode
	DataShareMode DataShareSpec `json:"dataShareMode,omitempty"`

	// Tasks specifies the task specification of Job
	Master TaskSpec `json:"master,omitempty"`

	Worker TaskSpec `json:"worker,omitempty"`

	// Specifies the maximum number of retries before marking this Job failed.
	// Defaults to 3.
	// +optional
	MaxRetry int32 `json:"maxRetry,omitempty"`

	// ttlSecondsAfterFinished limits the lifetime of a Job that has finished
	// execution (either Completed or Failed). If this field is set,
	// ttlSecondsAfterFinished after the Job finishes, it is eligible to be
	// automatically deleted. If this field is unset,
	// the Job won't be automatically deleted. If this field is set to zero,
	// the Job becomes eligible to be deleted immediately after it finishes.
	// +optional
	TTLSecondsAfterFinished *int32 `json:"ttlSecondsAfterFinished,omitempty" protobuf:"varint,9,opt,name=ttlSecondsAfterFinished"`

	// If specified, indicates the job's priority.
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty" protobuf:"bytes,10,opt,name=priorityClassName"`
}


// TaskSpec specifies the task specification of Job
type TaskSpec struct {
	// Name specifies the name of tasks
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`

	// Replicas specifies the replicas of this TaskSpec in Job
	Replicas int32 `json:"replicas,omitempty" protobuf:"bytes,2,opt,name=replicas"`

	// Specifies the pod that will be created for this TaskSpec
	// when executing a Job
	Template v1.PodTemplateSpec `json:"template,omitempty" protobuf:"bytes,3,opt,name=template"`

}

// VolumeSpec defines the specification of Volume, e.g. PVC
type VolumeSpec struct {
	// Path within the container at which the volume should be mounted.  Must
	// not contain ':'.
	MountPath string `json:"mountPath" protobuf:"bytes,1,opt,name=mountPath"`

	// defined the PVC name
	VolumeClaimName string `json:"volumeClaimName,omitempty" protobuf:"bytes,2,opt,name=volumeClaimName"`

	// VolumeClaim defines the PVC used by the VolumeMount.
	VolumeClaim *v1.PersistentVolumeClaimSpec `json:"volumeClaim,omitempty" protobuf:"bytes,3,opt,name=volumeClaim"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HorovodJob is the Schema for the horovodjobs API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type HorovodJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HorovodJobSpec   `json:"spec,omitempty"`
	Status HorovodJobStatus `json:"status,omitempty"`
}


// JobPhase defines the phase of the job
type JobPhase string

const (
	// Pending is the phase that job is pending in the queue, waiting for scheduling decision
	Pending JobPhase = "Pending"
	// Aborting is the phase that job is aborted, waiting for releasing pods
	Aborting JobPhase = "Aborting"
	// Aborted is the phase that job is aborted by user or error handling
	Aborted JobPhase = "Aborted"
	// Running is the phase that minimal available tasks of Job are running
	Running JobPhase = "Running"
	// Restarting is the phase that the Job is restarted, waiting for pod releasing and recreating
	Restarting JobPhase = "Restarting"
	// Completing is the phase that required tasks of job are completed, job starts to clean up
	Completing JobPhase = "Completing"
	// Completed is the phase that all tasks of Job are completed
	Completed JobPhase = "Completed"
	// Terminating is the phase that the Job is terminated, waiting for releasing pods
	Terminating JobPhase = "Terminating"
	// Terminated is the phase that the job is finished unexpected, e.g. events
	Terminated JobPhase = "Terminated"

	// job이 완료된 후 처리되는 상태
	PostProcessing JobPhase = "PostProcessing"
	// Failed is the phase that the job is restarted failed reached the maximum number of retries.
	Failed JobPhase = "Failed"
)

// JobState contains details for the current state of the job.
type HorovodJobState struct {
	// The phase of Job.
	// +optional
	Phase JobPhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase"`

	// Unique, one-word, CamelCase reason for the phase's last transition.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,2,opt,name=reason"`

	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,3,opt,name=message"`

	// Last time the condition transit from one phase to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`
}
// HorovodJobStatus defines the observed state of HorovodJob
// +k8s:openapi-gen=true
type HorovodJobStatus struct {
	// Current state of Job.
	State HorovodJobState `json:"state,omitempty" protobuf:"bytes,1,opt,name=state"`

	// The minimal available pods to run for this Job
	// +optional
	MinAvailable int32 `json:"minAvailable,omitempty" protobuf:"bytes,2,opt,name=minAvailable"`

	// The number of pending pods.
	// +optional
	Pending int32 `json:"pending,omitempty" protobuf:"bytes,3,opt,name=pending"`

	// The number of running pods.
	// +optional
	Running int32 `json:"running,omitempty" protobuf:"bytes,4,opt,name=running"`

	// The number of pods which reached phase Succeeded.
	// +optional
	Succeeded int32 `json:"succeeded,omitempty" protobuf:"bytes,5,opt,name=succeeded"`

	// The number of pods which reached phase Failed.
	// +optional
	Failed int32 `json:"failed,omitempty" protobuf:"bytes,6,opt,name=failed"`

	// The number of pods which reached phase Terminating.
	// +optional
	Terminating int32 `json:"terminating,omitempty" protobuf:"bytes,7,opt,name=terminating"`

	// The number of pods which reached phase Unknown.
	// +optional
	Unknown int32 `json:"unknown,omitempty" protobuf:"bytes,8,opt,name=unknown"`

	//Current version of job
	Version int32 `json:"version,omitempty" protobuf:"bytes,9,opt,name=version"`

	// The number of Job retries.
	// +optional
	RetryCount int32 `json:"retryCount,omitempty" protobuf:"bytes,10,opt,name=retryCount"`

	// The resources that controlled by this job, e.g. Service, ConfigMap
	ControlledResources map[string]string `json:"controlledResources,omitempty" protobuf:"bytes,11,opt,name=controlledResources"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HorovodJobList contains a list of HorovodJob
type HorovodJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HorovodJob `json:"items"` 
}

func init() {
	SchemeBuilder.Register(&HorovodJob{}, &HorovodJobList{})
}
