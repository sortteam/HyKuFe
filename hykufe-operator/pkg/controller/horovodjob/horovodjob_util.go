package horovodjob

import (
	"context"
	"encoding/json"
	"hykufe-operator/pkg/apis/hykufe/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DefinitionToJson(instance interface{}) []byte{
	jsonByte, err := json.Marshal(instance)
	if err != nil {

	}
	return jsonByte
}

func (r *ReconcileHorovodJob) UpdateState(instance *v1alpha1.HorovodJob, phase v1alpha1.JobPhase) error {
	instance.Status.State.Phase = phase
	instance.Status.State.LastTransitionTime = v1.Now()

	if err := r.client.Status().Update(context.TODO(), instance); err != nil {
		return err
	}

	return nil
}
