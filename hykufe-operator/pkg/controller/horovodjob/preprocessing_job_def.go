package horovodjob

import (
	hykufev1alpha1 "hykufe-operator/pkg/apis/hykufe/v1alpha1"
	v12 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileHorovodJob) newPreProcessingJob(cr *hykufev1alpha1.HorovodJob) (*v12.Job, error) {
	labels := map[string]string {
		"app": cr.Name,
	}
	parallelism := int32(1)
	completions := int32(1)
	backOffLimit := int32(6)
	ttlSeconds := int32(60)


	jobDefinition := &v12.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:                      cr.Name + "-preprocessing-job",
			Namespace:                  cr.Namespace,
			Labels:                     labels,
		},
		Spec:       v12.JobSpec{
			Parallelism:             &parallelism,
			Completions:             &completions,
			BackoffLimit:            &backOffLimit,
			Template:                v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:                      	cr.Name + "-preprocessing-pod",
					Namespace:                  cr.Namespace,
					Labels:                     labels,
				},
				Spec:       v1.PodSpec{
					Volumes:                       []v1.Volume{
					},
					Containers:                    []v1.Container{
						{
							Name:                     "empty",
							Image:                    "alpine",
							Command:                  nil,
						},
					},
					RestartPolicy:                 "OnFailure",
				},
			},
			TTLSecondsAfterFinished: &ttlSeconds,
		},
	}

	if len(cr.Spec.DataSources) == 0 {
		return nil, nil
	}

	for _, datasource := range cr.Spec.DataSources {
		if datasource.S3Source != nil {
			containerRef := &jobDefinition.Spec.Template.Spec.Containers
			*containerRef = append(*containerRef, v1.Container{
				Name:                     "ds-s3-" + datasource.Name,
				Image:                    "banst/awscli",
				Command:                  []string {
					"/bin/sh",
				},
				Args:                     []string {
					"-c",
					//"sleep 200;",
					"aws s3 cp --recursive s3://${AWS_BUCKET}/${AWS_S3_DIRECTORY} /data/",
				},
				WorkingDir:               "",
				EnvFrom:                  []v1.EnvFromSource{
					{
						Prefix:       "",
						ConfigMapRef: nil,
						SecretRef:    &v1.SecretEnvSource{
							LocalObjectReference: v1.LocalObjectReference{
								Name: datasource.Name,
							},
							Optional:             nil,
						},
					},
				},
				Env:                      []v1.EnvVar{
					{
						Name:      "AWS_S3_DIRECTORY",
						Value:     datasource.S3Source.Directory,
					},
				},
				VolumeMounts:             getVolumeMount(cr),
			})
		}
	}

	for _, volume := range cr.Spec.Volumes {
		jobDefinition.Spec.Template.Spec.Volumes = append(jobDefinition.Spec.Template.Spec.Volumes, v1.Volume{
			Name:         volume.VolumeClaimName,
			VolumeSource: v1.VolumeSource{
				PersistentVolumeClaim: & v1.PersistentVolumeClaimVolumeSource{
					ClaimName: volume.VolumeClaimName,
					ReadOnly:  false,
				},
			},
		})
	}

	// Set PostProcessing pod instance as the owner of the Horovodjob
	if err := controllerutil.SetControllerReference(cr, jobDefinition, r.scheme); err != nil {
		log.Error(err, "Post processing pod의 owner를 Horovodjob으로 설정할 수 없습니다.")
		return nil, err
	}

	return jobDefinition, nil
}

func (r *ReconcileHorovodJob) newPVCDefinition(instance *hykufev1alpha1.HorovodJob) (*v1.PersistentVolumeClaim, error) {
	volume := instance.Spec.Volumes[0]
	pvcDefinition := &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:                       volume.VolumeClaimName,
			Namespace:                  instance.Namespace,
		},
		Spec:       *volume.VolumeClaim,
	}


	// Set PostProcessing pod instance as the owner of the Horovodjob
	if err := controllerutil.SetControllerReference(instance, pvcDefinition, r.scheme); err != nil {
		log.Error(err, "PVC의 owner를 Horovodjob으로 설정할 수 없습니다.")
		return nil, err
	}

	return pvcDefinition, nil
}