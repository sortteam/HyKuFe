package apis

import (
	"hykufe-operator/pkg/apis/hykufe/v1alpha1"
	//"k8s.io/api/batch/v1"
	volcanov1alpha1 "github.com/volcano-sh/volcano/pkg/apis/batch/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v1alpha1.SchemeBuilder.AddToScheme)
	//AddToSchemes = append(AddToSchemes, v1.SchemeBuilder.AddToScheme)
	AddToSchemes = append(AddToSchemes, volcanov1alpha1.SchemeBuilder.AddToScheme)
}
