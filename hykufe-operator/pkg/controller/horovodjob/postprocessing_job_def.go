package horovodjob

import (
	v12 "k8s.io/api/batch/v1"
	hykufev1alpha1 "hykufe-operator/pkg/apis/hykufe/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	v1 "k8s.io/api/core/v1"
)

func (r *ReconcileHorovodJob) newPostProcessingPod(cr *hykufev1alpha1.HorovodJob) (*v12.Job, error) {
	labels := map[string]string {
		"app": cr.Name,
	}
	parallelism := int32(1)
	completions := int32(1)
	backOffLimit := int32(6)
	ttlSeconds := int32(60)

	jobDefinition := &v12.Job{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
			Namespace: cr.Namespace,
			Name: cr.Name + "-postprocessing-job",
		},
		Spec:       v12.JobSpec{
			Parallelism:             &parallelism,
			Completions:             &completions,
			BackoffLimit:            &backOffLimit,
			Template:                v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:                      	cr.Name + "-postprocessing-pod",
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

	if cr.Spec.AwsSpec != nil {
		jobDefinition.Spec.Template.Spec.NodeSelector = map[string]string {
			"cloud-type": "aws",
		}
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
					"aws s3 cp --recursive /data s3://${AWS_BUCKET}/${AWS_S3_DIRECTORY}/model/",
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