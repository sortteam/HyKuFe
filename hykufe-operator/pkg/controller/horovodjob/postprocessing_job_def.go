package horovodjob

import (
	hykufev1alpha1 "hykufe-operator/pkg/apis/hykufe/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	v1 "k8s.io/api/core/v1"
)

func (r *ReconcileHorovodJob) newPostProcessingPod(cr *hykufev1alpha1.HorovodJob) (*v1.Pod, error) {
	labels := map[string]string {
		"app": cr.Name,
	}

	podDefinition := &v1.Pod{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       cr.Name + "-postprocessing",
			Namespace:                  cr.Namespace,
			Labels:                     labels,
		},
		Spec:       v1.PodSpec{
			Volumes:                       nil,
			InitContainers:                nil,
			Containers:                    nil,
			RestartPolicy:                 "OnFailure",
			DNSPolicy:                     "",
			NodeSelector:                  nil,
			ServiceAccountName:            "",
			DeprecatedServiceAccount:      "",
			AutomountServiceAccountToken:  nil,
			ImagePullSecrets:              nil,
			Hostname:                      "",
			Subdomain:                     "",
			Affinity:                      nil,
			SchedulerName:                 "",
			Tolerations:                   nil,
			HostAliases:                   nil,
			PriorityClassName:             "",
			Priority:                      nil,
			DNSConfig:                     nil,
			ReadinessGates:                nil,
			RuntimeClassName:              nil,
			EnableServiceLinks:            nil,
		},
	}

	// Set PostProcessing pod instance as the owner of the Horovodjob
	if err := controllerutil.SetControllerReference(cr, podDefinition, r.scheme); err != nil {
		log.Error(err, "Post processing pod의 owner를 Horovodjob으로 설정할 수 없습니다.")
		return nil, err
	}

	return podDefinition, nil
}