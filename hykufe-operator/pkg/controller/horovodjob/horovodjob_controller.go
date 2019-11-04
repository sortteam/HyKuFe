package horovodjob

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	volcanov1alpha1 "github.com/volcano-sh/volcano/pkg/apis/batch/v1alpha1"
	hykufev1alpha1 "hykufe-operator/pkg/apis/hykufe/v1alpha1"
	v1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"
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
	return &ReconcileHorovodJob{client: mgr.GetClient(), scheme: mgr.GetScheme(), awsController: NewAWSController()}
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

	err = c.Watch(&source.Kind{Type: &v1.Job{}}, &handler.EnqueueRequestForOwner{
		OwnerType:    &hykufev1alpha1.HorovodJob{},
		IsController: true,
	})

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
	awsController *AWSController
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

	err2 := r.controlProvisioning(instance, reqLogger)
	if err2 != nil {
		return reconcile.Result{}, err2
	}

	time, err := r.controlPreprocessingJob(instance, reqLogger)
	if err != nil || time != 0 {
		if time != 0 {
			return reconcile.Result{
				Requeue:      true,
				RequeueAfter: time,
			}, err
		} else {
			return reconcile.Result{}, err
		}
	}


	// Define a new Object
	volcanoJob, err := r.newVolcanoJobForCR(instance)
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
	foundVolcanoJob := &volcanov1alpha1.Job{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: volcanoJob.Name, Namespace: volcanoJob.Namespace}, foundVolcanoJob)
	if err != nil && errors.IsNotFound(err) && instance.Status.State.Phase == hykufev1alpha1.Preprocessed{
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

	// Sync Status


	oldHorovodJobState := instance.Status.State
	//oldVolcanoJobState := foundVolcanoJob.Status.State

	if updateState(instance, foundVolcanoJob) == nil {
		if !reflect.DeepEqual(oldHorovodJobState, instance.Status.State) {
			log.Info("update horovod job status!!!")
			// Update HorovodJob Status
			jsonByte, err := json.Marshal(instance)
			if err != nil {

			}
			reqLogger.Info("CR definition", "horovodjob", string(jsonByte))
			log.Error(err, "Faid to update horovodjob resource status")
			if err := r.client.Status().Update(context.TODO(), instance); err != nil {
				jsonByte, err := json.Marshal(instance)
				if err != nil {

				}
				reqLogger.Info("CR definition", "horovodjob", string(jsonByte))
				log.Error(err, "Faid to update horovodjob resource status")
				return reconcile.Result{}, err
			}
		}

		//if !reflect.DeepEqual(oldVolcanoJobState, foundVolcanoJob.Status.State) {
		//	// Update VolcanoJob Status
		//	if err := r.client.Status().Update(context.TODO(), foundVolcanoJob); err != nil {
		//		log.Error(err, "Fail to update volcanojob resource status")
		//		return reconcile.Result{}, err
		//	}
		//}
	}



	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: VolcanoJob already exists", "Pod.Namespace", foundVolcanoJob.Namespace, "Pod.Name", foundVolcanoJob.Name)
	return reconcile.Result{}, nil
}




func updateState(horovodJob *hykufev1alpha1.HorovodJob, volcanoJob *volcanov1alpha1.Job) error {
	nowHorovodJobState := &horovodJob.Status.State
	nowVolcanoJobState := &volcanoJob.Status.State



	if nowVolcanoJobState.Phase == volcanov1alpha1.Pending {
		if nowHorovodJobState.Phase != hykufev1alpha1.Pending {
			nowHorovodJobState.Phase = hykufev1alpha1.Pending
			nowHorovodJobState.LastTransitionTime = metav1.Now()
		}
	} else if nowVolcanoJobState.Phase == volcanov1alpha1.Failed {
		if nowHorovodJobState.Phase != hykufev1alpha1.Failed {
			nowHorovodJobState.Phase = hykufev1alpha1.Failed
			nowHorovodJobState.LastTransitionTime = metav1.Now()
		}
	} else if nowVolcanoJobState.Phase == volcanov1alpha1.Running {
		if nowHorovodJobState.Phase != hykufev1alpha1.Running {
			nowHorovodJobState.Phase = hykufev1alpha1.Running
			nowHorovodJobState.LastTransitionTime = metav1.Now()
		}
	} else if nowVolcanoJobState.Phase == volcanov1alpha1.Aborted {
		if nowHorovodJobState.Phase != hykufev1alpha1.Aborted {
			nowHorovodJobState.Phase = hykufev1alpha1.Aborted
			nowHorovodJobState.LastTransitionTime = metav1.Now()
		}
	} else if nowVolcanoJobState.Phase == volcanov1alpha1.Completed {
		if nowHorovodJobState.Phase != hykufev1alpha1.PostProcessing {
			nowHorovodJobState.Phase = hykufev1alpha1.PostProcessing
			nowHorovodJobState.LastTransitionTime = metav1.Now()
		}
	}

	if nowHorovodJobState.Phase == hykufev1alpha1.Completed {
		if nowVolcanoJobState.Phase != volcanov1alpha1.Terminating {
			nowVolcanoJobState.Phase = volcanov1alpha1.Terminating
			nowVolcanoJobState.LastTransitionTime = metav1.Now()

			nowHorovodJobState.Phase = hykufev1alpha1.Terminating
			nowHorovodJobState.LastTransitionTime = metav1.Now()
		}
	}
	// last Transition 수정하기

	return nil
}

func (r *ReconcileHorovodJob) controlPreprocessingJob(instance *hykufev1alpha1.HorovodJob, reqLogger logr.Logger) (time.Duration, error) {
	if !(instance.Status.State.Phase == hykufev1alpha1.Provisioned || instance.Status.State.Phase == hykufev1alpha1.Preprocessing) {
		return 0, nil
	}
	preprocessingJob, err := r.newPreProcessingJob(instance)
	if err != nil {
		reqLogger.Error(err, "Fail to make preprocessing job")
		return 0, err
	}
	pvc, err := r.newPVCDefinition(instance)
	if err != nil {
		reqLogger.Error(err, "Fail to make PVC")
		return 0, err
	}

	foundPreprocessingJob := &v1.Job{}

	err = r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: preprocessingJob.Namespace,
		Name:      preprocessingJob.Name,
	}, foundPreprocessingJob)

	if err != nil && errors.IsNotFound(err) && (instance.Status.State.Phase == "" || instance.Status.State.Phase == hykufev1alpha1.Provisioned) {
		reqLogger.Info("Creating a new Preprocessing Job")
		err = r.client.Create(context.TODO(), preprocessingJob)
		if err != nil {
			reqLogger.Error(err, "Fail to create preprocessing job")
			return 0, err
		}

		reqLogger.Info("Creating a new PVC")
		if err := r.client.Create(context.TODO(), pvc); err != nil {
			reqLogger.Error(err, "Fail to create PVC")
			return 0, err
		}

		// horovodjob 상태를 업데이트 한다.
		instance.Status.State.Phase = hykufev1alpha1.Preprocessing
		instance.Status.State.LastTransitionTime = metav1.Now()

		// 상태 업데이트
		updateErr := r.client.Status().Update(context.TODO(), instance)
		if updateErr != nil {
			reqLogger.Error(updateErr, "fail to update horovodjob instance")
		}

	}

	if instance.Status.State.Phase == hykufev1alpha1.Preprocessing {
		reqLogger.Info("wait to finalize preprocessing job")
		condition := foundPreprocessingJob.Status.Conditions
		reqLogger.Info(fmt.Sprintf("%s", DefinitionToJson(foundPreprocessingJob)))
		// preprocessing job이 성공했을 때
		if len(condition) != 0 &&
			(condition[0].Type == v1.JobComplete ||
			condition[0].Type == v1.JobFailed) {
			// TODO: Refactor code duplication



			// preprocessing job이 완료되면 horovod job의 상태를 preprocessed로 변경
			if condition[len(condition) - 1].Type == v1.JobComplete{
				instance.Status.State.Phase = hykufev1alpha1.Preprocessed
			} else if condition[len(condition) - 1].Type == v1.JobFailed {
				instance.Status.State.Phase = hykufev1alpha1.Failed
			}
			instance.Status.State.LastTransitionTime = metav1.Now()

			// 임시 PVC 삭제
			//if err := r.client.Delete(context.TODO(), pvc); err != nil {
			//	reqLogger.Error(err, "fail to delete temp pvc")
			//	return 0, err
			//}

			// 상태 업데이트
			updateErr := r.client.Status().Update(context.TODO(), instance)
			if updateErr != nil {
				reqLogger.Error(updateErr, "fail to update horovodjob instance")
				return 0, err
			}
		}


		return 0, nil
	}

	if instance.Status.State.Phase == hykufev1alpha1.Preprocessed {
		reqLogger.Info("completed preprocessing job")
	}

	return 0, nil
}

func (r *ReconcileHorovodJob) controlProvisioning(instance *hykufev1alpha1.HorovodJob, reqLogger logr.Logger) error {

	nowState := instance.Status.State.Phase
	if instance.Spec.AwsSpec == nil {
		return nil
	}

	if nowState == hykufev1alpha1.Provisioned {
		return nil
	}

	if nowState == "" || nowState == hykufev1alpha1.Pending || nowState == hykufev1alpha1.Provisioning {
		// 상태를 Provisining으로 변경
		if err := r.UpdateState(instance, hykufev1alpha1.Provisioning); err != nil {
			reqLogger.Error(err, "Fail to Update Horovod State")
			return err
		}

		// AWS 스펙에 맞게 인스턴스 생성
		reqLogger.Info("Create EC2 Instance...")
		ec2Instances, err := r.awsController.CreateEC2Instance(instance.Spec.AwsSpec.InstanceType, instance.Spec.AwsSpec.Replicas)
		if err != nil {
			reqLogger.Error(err, "Fail to create EC2 Instances")
			return err
		}
		instance.Status.InstanceID = []string{}
		for _, ec2Info := range ec2Instances {
			instance.Status.InstanceID = append(instance.Status.InstanceID, *ec2Info.InstanceId)
		}

		if err := r.UpdateState(instance, hykufev1alpha1.Provisioned); err != nil {
			reqLogger.Error(err, "Fail to Update Horovod State Provisioned")
			return err
		}
		reqLogger.Info("Created EC2 Instance!!!")
	}

	return nil
}