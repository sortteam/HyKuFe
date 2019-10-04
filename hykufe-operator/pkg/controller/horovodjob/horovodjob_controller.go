package horovodjob

import (
	"context"
	"fmt"
	volcanov1alpha1 "github.com/volcano-sh/volcano/pkg/apis/batch/v1alpha1"
	hykufev1alpha1 "hykufe-operator/pkg/apis/hykufe/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_horovodjob")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new HorovodJob Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileHorovodJob{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("horovodjob-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource HorovodJob
	err = c.Watch(&source.Kind{Type: &hykufev1alpha1.HorovodJob{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner HorovodJob
	// HorovodJob에 속해 있는 Volcano Job을 Watch
	err = c.Watch(&source.Kind{Type: &volcanov1alpha1.Job{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &hykufev1alpha1.HorovodJob{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileHorovodJob implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileHorovodJob{}

// ReconcileHorovodJob reconciles a HorovodJob object
type ReconcileHorovodJob struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a HorovodJob object and makes changes based on the state read
// and what is in the HorovodJob.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileHorovodJob) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling HorovodJob")

	// Fetch the HorovodJob instance
	instance := &hykufev1alpha1.HorovodJob{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	jsonByte, err := json.Marshal(instance)
	if err != nil {

	}
	reqLogger.Info("CR definition", "horovodjob", string(jsonByte))
	// Define a new Pod object
	volcanoJob, err := newVolcanoJobForCR(instance)
	if err != nil {
		updateErr := r.client.Status().Update(context.TODO(), instance)
		if updateErr != nil {
			reqLogger.Error(updateErr, "fail to update horovodjob instance")
		}

		return reconcile.Result{}, err
	}


	// Set HorovodJob instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, volcanoJob, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this HorovodJob already exists
	found := &volcanov1alpha1.Job{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: volcanoJob.Name, Namespace: volcanoJob.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new VolcanoJob", "VolcanoJob.Namespace", volcanoJob.Namespace, "VolcanoJob.Name", volcanoJob.Name)
		err = r.client.Create(context.TODO(), volcanoJob)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: VolcanoJob already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
//func newPodForCR(cr *hykufev1alpha1.HorovodJob) *corev1.Pod {
//	labels := map[string]string{
//		"app": cr.Name,
//	}
//	return &corev1.Pod{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      cr.Name + "-pod",
//			Namespace: cr.Namespace,
//			Labels:    labels,
//		},
//		Spec: corev1.PodSpec{
//			Containers: []corev1.Container{
//				{
//					Name:    "busybox",
//					Image:   "busybox",
//					Command: []string{"sleep", "3600"},
//				},
//			},
//		},
//	}
//}

func newVolcanoJobForCR(cr *hykufev1alpha1.HorovodJob) (*volcanov1alpha1.Job, error) {
	labels := map[string]string {
		"app": cr.Name,
	}
	volcanojob := &volcanov1alpha1.Job {
		ObjectMeta: metav1.ObjectMeta{
			Name:		cr.Name + "-volcanojob",
			Namespace:	cr.Namespace,
			Labels:		labels,
		},
		Spec: volcanov1alpha1.JobSpec{
			// SchedulerName:           "",
			MinAvailable:            cr.Spec.Worker.Replicas + 1,
			Tasks: []volcanov1alpha1.TaskSpec{
				volcanov1alpha1.TaskSpec{
					Name:     cr.Spec.Master.Name,
					Replicas: 1,
					Template: cr.Spec.Master.Template,
					Policies: nil,
				},
				volcanov1alpha1.TaskSpec{
					Name:     cr.Spec.Worker.Name,
					Replicas: cr.Spec.Worker.Replicas,
					Template: cr.Spec.Worker.Template,
					Policies: nil,
				},
			},
			//Volumes:                 nil,
			//Policies:                {

			Plugins:                 map[string][]string{
				"ssh": []string{},
				"svc": []string{},
			},
			//Queue:                   "",
			MaxRetry:                cr.Spec.MaxRetry,
			TTLSecondsAfterFinished: cr.Spec.TTLSecondsAfterFinished,
			PriorityClassName:       cr.Spec.PriorityClassName,
		},
	}


	// add Sidecar Container
	//volcanojob.Spec.Tasks[0].Template.Spec.Container
	masterJobSpec := &volcanojob.Spec.Tasks[0].Template.Spec
	workerJobSpec := &volcanojob.Spec.Tasks[1].Template.Spec


	// Sync Process namespace with all containers
	t := true
	masterJobSpec.ShareProcessNamespace = &t

	//masterJobSpec.Containers[0].LivenessProbe.Exec.Command = []string{
	//	"/bin/sh",
	//	"-c",
	//	"horovod_pid=$(ps -A | grep mpiexec | awk '/!(grep)/ { print $1 }')",
	//	"if [ \"$horovod_pid\" != \"\" ]\"",
	//	"then",
	//	"exit 0",
	//	"else",
	//	"exit 1",
	//	"fi",
	//}
	//masterJobSpec.Containers[0].LivenessProbe.InitialDelaySeconds = 20;

	// Add EmptyDir Volume for saving model, log, etc...
	masterJobSpec.Volumes = append(masterJobSpec.Volumes, v1.Volume{
		Name:         "result-data-volume",
		VolumeSource: v1.VolumeSource{
			EmptyDir: &v1.EmptyDirVolumeSource{
			},
		},
	})

	// Mount Volume to main container
	masterJobSpec.Containers[0].VolumeMounts = append(masterJobSpec.Containers[0].VolumeMounts, v1.VolumeMount{
		Name:      "result-data-volume",
		ReadOnly:  false,
		MountPath: "/result",
		//MountPropagation: nil,
	},
	)

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

	// Add NFS Volume For data save
	if cr.Spec.DataShareMode.NFSMode != nil {
		masterJobSpec.Volumes = append(masterJobSpec.Volumes, v1.Volume{
			Name:         "data-volume",
			VolumeSource: v1.VolumeSource{
				NFS: &v1.NFSVolumeSource{
					// FIXME : 임시로 지정
					Server:   cr.Spec.DataShareMode.NFSMode.IPAddress,
					Path:     cr.Spec.DataShareMode.NFSMode.Path,
					ReadOnly: false,
				},
			},
		})

		// Mount data volume to master
		masterJobSpec.Containers[0].VolumeMounts = append(masterJobSpec.Containers[0].VolumeMounts, v1.VolumeMount{
			Name:             "data-volume",
			ReadOnly:         true,
			MountPath:        "/data",
		})

		// Add NFS Volume For data save
		workerJobSpec.Volumes = append(workerJobSpec.Volumes, v1.Volume{
			Name:         "data-volume",
			VolumeSource: v1.VolumeSource{

				NFS: &v1.NFSVolumeSource{
					Server:   cr.Spec.DataShareMode.NFSMode.IPAddress,
					Path:     cr.Spec.DataShareMode.NFSMode.Path,
					ReadOnly: true,
				},
			},
		})

		//Mount data volume to worker
		workerJobSpec.Containers[0].VolumeMounts = append(workerJobSpec.Containers[0].VolumeMounts, v1.VolumeMount{
			Name:             "data-volume",
			ReadOnly:         true,
			MountPath:        "/data",
		})

	}


	if len(masterJobSpec.Containers) == 0 {
		return nil, fmt.Errorf("must exist master pods spec")
	}

	// Add Sidecar Container and attach volumes
	masterJobSpec.Containers = append(masterJobSpec.Containers, v1.Container{
		Name:                     "sidecar-container",
		Image:                    "banst/awscli",
		Command:                  []string{ "/bin/sh", "/exec/sidecar.sh" },
		Args:                     nil,
		WorkingDir:               "/",
		Ports:                    nil,
		EnvFrom:                  nil,
		Env:                      nil,
		Resources:                v1.ResourceRequirements{},
		VolumeMounts:             []v1.VolumeMount{
			{
				Name:             "result-data-volume",
				ReadOnly:         false,
				MountPath:        "/result",
			},
			{
				Name:             "horovod-cm",
				ReadOnly:         false,
				MountPath:        "/exec",
			},
		},
		ImagePullPolicy:          "",
		SecurityContext:          nil,
	})

	// add initContainer for data sync from data source
	volcanojob.Spec.Tasks[0].Template.Spec.InitContainers = []v1.Container{}

	for i, dataSource := range cr.Spec.DataSources {
		// S3 데이터 처리용 initContainer 추가
		if dataSource.S3Source != nil {
			volcanojob.Spec.Tasks[0].Template.Spec.InitContainers = append(volcanojob.Spec.Tasks[0].Template.Spec.InitContainers, v1.Container{
				Name:                     fmt.Sprintf("initcontainer-%d", i),
				Image:                    "banst/awscli",
				Command:                  []string{
					"/bin/sh",
				},
				Args:                     []string{
					"-c",
					"aws s3 cp --recursive s3://${AWS_S3_BUCKET}/${AWS_S3_DIRECTORY} /data/${DATA_SOURCE_NAME}",
				},
				WorkingDir:               "/data",
				Ports:                    nil,
				Env:                      []v1.EnvVar{
					{
						Name:	"AWS_ACCESS_KEY_ID",
						Value:	dataSource.S3Source.AccessKeyId,
					},
					{
						Name:	"AWS_SECRET_ACCESS_KEY",
						Value: 	dataSource.S3Source.SecretAccessKey,
					},
					{
						Name:	"AWS_DEFAULT_REGION",
						Value:	dataSource.S3Source.Region,
					},
					{
						Name:	"AWS_S3_BUCKET",
						Value: dataSource.S3Source.Bucket,
					},
					{
						Name: "AWS_S3_DIRECTORY",
						Value: dataSource.S3Source.Directory,
					},
					{
						Name: "DATA_SOURCE_NAME",
						Value: dataSource.Name,
					},
				},
				VolumeMounts:             []v1.VolumeMount{
					{
						Name:             "data-volume",
						ReadOnly:         false,
						MountPath:        "/data",
					},
				},
			})
		}
	}

	jsonByte, err := json.Marshal(volcanojob)
	if err != nil {

	}
	log.Info(string(jsonByte))

	return volcanojob, nil
}

func validateHorovodJobCR(cr *hykufev1alpha1.HorovodJob) error {

	// Validate DataSource
	for _, dataSource := range cr.Spec.DataSources {
		if dataSource.S3Source.AccessKeyId == "" {
			return fmt.Errorf("Access Key ID must entered")
		}
		if dataSource.S3Source.SecretAccessKey == "" {
			return fmt.Errorf("Secret Access Key must entered")
		}
		if dataSource.S3Source.Region == "" {
			return fmt.Errorf("Region must entered")
		}
		if dataSource.S3Source.Bucket == "" {
			return fmt.Errorf("Bucket must entered")
		}
		if dataSource.S3Source.Directory == "" {
			return fmt.Errorf("DirectoryName must entered")
		}
	}

	return nil
}

