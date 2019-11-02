package horovodjob

import (
	"fmt"
	volcanov1alpha1 "github.com/volcano-sh/volcano/pkg/apis/batch/v1alpha1"
	hykufev1alpha1 "hykufe-operator/pkg/apis/hykufe/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileHorovodJob) newVolcanoJobForCR(cr *hykufev1alpha1.HorovodJob) (*volcanov1alpha1.Job, error) {
	labels := map[string]string {
		"app": cr.Name,
	}

	copiedHorovodJob := cr.DeepCopy()

	volcanojob := &volcanov1alpha1.Job {
		ObjectMeta: metav1.ObjectMeta{
			Name:		copiedHorovodJob.Name + "-volcanojob",
			Namespace:	copiedHorovodJob.Namespace,
			Labels:		labels,
		},
		Spec: volcanov1alpha1.JobSpec{
			// SchedulerName:           "",
			MinAvailable:            copiedHorovodJob.Spec.Worker.Replicas + 1,
			Tasks: []volcanov1alpha1.TaskSpec{
				{
					Name:     copiedHorovodJob.Spec.Master.Name,
					Replicas: 1,
					Template: copiedHorovodJob.Spec.Master.Template,
					Policies: []volcanov1alpha1.LifecyclePolicy{
						{
							Action:   "CompleteJob",
							Event:    "TaskCompleted",
						},
					},
				},
				{
					Name:     copiedHorovodJob.Spec.Worker.Name,
					Replicas: copiedHorovodJob.Spec.Worker.Replicas,
					Template: copiedHorovodJob.Spec.Worker.Template,
					Policies: nil,
				},
			},
			Volumes: copiedHorovodJob.Spec.Volumes,
			//Volumes:                 nil,
			//Policies:                {

			Plugins:                 map[string][]string{
				"ssh": []string{},
				"svc": []string{},
			},
			//Queue:                   "",
			MaxRetry:                copiedHorovodJob.Spec.MaxRetry,
			TTLSecondsAfterFinished: copiedHorovodJob.Spec.TTLSecondsAfterFinished,
			PriorityClassName:       copiedHorovodJob.Spec.PriorityClassName,
		},
	}


	// add Sidecar Container
	//volcanojob.Spec.Tasks[0].Template.Spec.Container
	masterJobSpec := &volcanojob.Spec.Tasks[0].Template.Spec
	//workerJobSpec := &volcanojob.Spec.Tasks[1].Template.Spec

	// Sync Process namespace with all containers
	t := true
	masterJobSpec.ShareProcessNamespace = &t

	// Add Configmap Volume for sidecar container
	mode := int32(365)
	masterJobSpec.Volumes = append(masterJobSpec.Volumes, v1.Volume{
		Name:         "horovod-cm",
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name:"horovod-cm",
				},
				Items: []v1.KeyToPath{
					{
						Key:  "sidecar.run",
						Path: "sidecar.sh",
						Mode: &mode,
					},
				},
			},
		},
	})

	if len(masterJobSpec.Containers) == 0 {
		return nil, fmt.Errorf("must exist master pods spec")
	}

	//  만약 S3 Source가 있다면 다운로드 진행
	//if len(cr.Spec.DataSources) != 0 {
	//	masterJobSpec.InitContainers = []v1.Container{}
	//	for _, dataSource := range cr.Spec.DataSources {
	//		if dataSource.S3Source != nil  {
	//
	//			masterJobSpec.InitContainers = append(masterJobSpec.InitContainers, v1.Container{
	//				Name:                     "init-s3-" + dataSource.Name,
	//				Image:                    "banst/awscli",
	//				Command:                  []string {
	//					"/bin/sh",
	//				},
	//				Args:                     []string {
	//					"-c",
	//					"aws s3 cp --recursive s3://${AWS_S3_BUCKET}/${AWS_S3_DIRECTORY} /data/${DATA_SOURCE_NAME}",
	//				},
	//				WorkingDir:               "",
	//				EnvFrom:                  []v1.EnvFromSource{
	//					{
	//						Prefix:       "",
	//						ConfigMapRef: nil,
	//						SecretRef:    &v1.SecretEnvSource{
	//							LocalObjectReference: v1.LocalObjectReference{
	//								Name: dataSource.Name,
	//							},
	//							Optional:             nil,
	//						},
	//					},
	//				},
	//				Env:                      []v1.EnvVar{
	//					{
	//						Name:      "DIRECTORY",
	//						Value:     dataSource.S3Source.Directory,
	//					},
	//				},
	//				VolumeMounts:             getVolumeMount(cr),
	//
	//				ImagePullPolicy:          "",
	//
	//			})
	//		}
	//	}
	//}


	//jsonByte, err := json.Marshal(volcanojob)
	//if err != nil {
	//
	//}
	//log.Info(string(jsonByte))

	// Set HorovodJob instance as the owner of the VolcanoJob

	if err := controllerutil.SetControllerReference(cr, volcanojob, r.scheme); err != nil {
		log.Error(err, "Volcanojob의 owner를 Horovodjob으로 설정할 수 없습니다.")
	}
	return volcanojob, nil
}

func getVolumeMount(cr *hykufev1alpha1.HorovodJob)  []v1.VolumeMount {
	vm := []v1.VolumeMount{}

	for _, volume := range cr.Spec.Volumes {
		vm = append(vm, v1.VolumeMount{
			Name:             volume.VolumeClaimName,
			ReadOnly:         false,
			MountPath:        volume.MountPath,
		})
	}

	return vm
}