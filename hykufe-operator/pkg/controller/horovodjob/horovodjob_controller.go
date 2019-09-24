package horovodjob

import (
	"context"
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

	// Define a new Pod object
	volcanoJob := newVolcanoJobForCR(instance)

	// Set HorovodJob instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, volcanoJob, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this HorovodJob already exists
	found := &hykufev1alpha1.HorovodJob{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: volcanoJob.Name, Namespace: volcanoJob.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new HorovodJob", "HorovodJob.Namespace", volcanoJob.Namespace, "HorovodJob.Name", volcanoJob.Name)
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

func newVolcanoJobForCR(cr *hykufev1alpha1.HorovodJob) *volcanov1alpha1.Job {
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

	// Sync Process namespace with all containers
	t := true
	masterJobSpec.ShareProcessNamespace = &t

	// Add EmptyDir Volume for saving model, log, etc...
	masterJobSpec.Volumes = append(masterJobSpec.Volumes, v1.Volume{
		Name:         "result-data-volume",
		VolumeSource: v1.VolumeSource{
			EmptyDir: &v1.EmptyDirVolumeSource{
			},
		},
	})

	// Add Volume to main container
	masterJobSpec.Containers[0].VolumeMounts = append(masterJobSpec.Containers[0].VolumeMounts, v1.VolumeMount{
			Name:      "result-data-volume",
			ReadOnly:  false,
			MountPath: "/result",
			//MountPropagation: nil,
		},
	)

	// Add Sidecar Container
	masterJobSpec.Containers = append(masterJobSpec.Containers, v1.Container{
		Name:                     "sidecar-container",
		Image:                    "alpine",
		Command:                  []string{ "/bin/sh" },
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
		},
		ImagePullPolicy:          "",
		SecurityContext:          nil,
	})

	jsonByte, err := json.Marshal(volcanojob)
	if err != nil {

	}
	log.Info(string(jsonByte))

	return volcanojob
}