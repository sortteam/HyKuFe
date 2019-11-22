package v1alpha1

type DataSourceSpec struct {
	Name string `json:"name"`
	S3Source *S3Spec `json:"s3Source,omitempty"`
}

type S3Spec struct {
	// S3 Secret Name
	S3SecretName string `json:"s3SecretName"`

	// Directory Name
	Directory string `json:"directory"`
}



