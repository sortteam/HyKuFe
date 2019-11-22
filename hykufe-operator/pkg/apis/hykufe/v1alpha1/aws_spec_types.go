package v1alpha1

type AwsSpec struct {
	InstanceType string `json:"instanceType"`
	Replicas int64 `json:"replicas"`
}