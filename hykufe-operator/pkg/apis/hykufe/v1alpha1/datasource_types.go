package v1alpha1

import v1 "k8s.io/api/core/v1"

type DataSourceSpec struct {
	Name string `json:"name"`
	S3Source *S3Spec `json:"s3Source"`
	v1.CephFSVolumeSource
}

type S3Spec struct {
	// Access Key
	AccessKeyId string `json:"accessKeyId"`

	// Secret Access key
	SecretAccessKey string `json:"secretAccessKey"`

	// 사용할 S3의 Region을 입력합니다.
	Region string `json:"region"`

	Bucket string `json:"bucket"`


}



