package mpi

import (
	"context"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/plugins"
	flyteerr "github.com/flyteorg/flyteplugins/go/tasks/errors"
	"github.com/flyteorg/flyteplugins/go/tasks/plugins/k8s/kfoperators/common"
	"github.com/flyteorg/flyteplugins/go/tasks/pluginmachinery"
	commonOp "github.com/kubeflow/common/pkg/apis/common/v1"
	pluginsCore "github.com/flyteorg/flyteplugins/go/tasks/pluginmachinery/core"
	"github.com/flyteorg/flyteplugins/go/tasks/pluginmachinery/flytek8s"
	"github.com/flyteorg/flyteplugins/go/tasks/pluginmachinery/k8s"
	"github.com/flyteorg/flyteplugins/go/tasks/pluginmachinery/utils"
	mpiOp "github.com/kubeflow/mpi-operator/pkg/apis/kubeflow/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type mpiOperatorResourceHandler struct {
}

// Sanity test that the plugin implements method of k8s.Plugin
var _ k8s.Plugin = mpiOperatorResourceHandler{}

func (mpiOperatorResourceHandler) GetProperties() k8s.PluginProperties {
	return k8s.PluginProperties{}
}

// Defines a func to create a query object (typically just object and type meta portions) that's used to query k8s
// resources.
func (mpiOperatorResourceHandler) BuildIdentityResource(ctx context.Context, taskCtx pluginsCore.TaskExecutionMetadata) (client.Object, error) {
	return &mpiOp.MPIJob{
		TypeMeta: metav1.TypeMeta{
			Kind:       mpiOp.Kind,
			APIVersion: mpiOp.SchemeGroupVersion.String(),
		},
	}, nil
}

// Defines a func to create the full resource object that will be posted to k8s.
func (mpiOperatorResourceHandler) BuildResource(ctx context.Context, taskCtx pluginsCore.TaskExecutionContext) (client.Object, error) {
	taskTemplate, err := taskCtx.TaskReader().Read(ctx)

	if err != nil {
		return nil, flyteerr.Errorf(flyteerr.BadTaskSpecification, "unable to fetch task specification [%v]", err.Error())
	} else if taskTemplate == nil {
		return nil, flyteerr.Errorf(flyteerr.BadTaskSpecification, "nil task specification")
	}

	mpiTaskExtraArgs := plugins.DistributedTensorflowTrainingTask{}
	err = utils.UnmarshalStruct(taskTemplate.GetCustom(), &mpiTaskExtraArgs)
	if err != nil {
		return nil, flyteerr.Errorf(flyteerr.BadTaskSpecification, "invalid TaskSpecification [%v], Err: [%v]", taskTemplate.GetCustom(), err.Error())
	}

	podSpec, err := flytek8s.ToK8sPodSpec(ctx, taskCtx)
	if err != nil {
		return nil, flyteerr.Errorf(flyteerr.BadTaskSpecification, "Unable to create pod spec: [%v]", err.Error())
	}

	workers := mpiTaskExtraArgs.GetWorkers()
	launcherReplicas := mpiTaskExtraArgs.GetPsReplicas()

	jobSpec := mpiOp.MPIJobSpec{
		SlotsPerWorker: &workers,
		//CleanPodPolicy: commonOp.CleanPodPolicyRunning,
		MPIReplicaSpecs: map[mpiOp.MPIReplicaType]*commonOp.ReplicaSpec{
			mpiOp.MPIReplicaTypeLauncher: &commonOp.ReplicaSpec{
				Replicas: &launcherReplicas,
				Template: v1.PodTemplateSpec{
					Spec: *podSpec,
				},
				RestartPolicy: commonOp.RestartPolicyNever,
			},
			mpiOp.MPIReplicaTypeWorker: &commonOp.ReplicaSpec{
				Replicas: &workers,
				Template: v1.PodTemplateSpec{
					Spec: *podSpec,
				},
				RestartPolicy: commonOp.RestartPolicyNever,
			},
		},
	}

	job := &mpiOp.MPIJob{
		TypeMeta: metav1.TypeMeta{
			Kind:       mpiOp.Kind,
			APIVersion: mpiOp.SchemeGroupVersion.String(),
		},
		Spec: jobSpec,
	}

	return job, nil
}

// Analyses the k8s resource and reports the status as TaskPhase. This call is expected to be relatively fast,
// any operations that might take a long time (limits are configured system-wide) should be offloaded to the
// background.
func (mpiOperatorResourceHandler) GetTaskPhase(_ context.Context, pluginContext k8s.PluginContext, resource client.Object) (pluginsCore.PhaseInfo, error) {
	app := resource.(*mpiOp.MPIJob)

	//workersCount := app.Spec.SlotsPerWorker
	//launcherReplicasCount := app.Spec.MPIReplicaSpecs

	//taskLogs, err := common.GetLogs(common.MPITaskType, app.Name, app.Namespace,
	//	*workersCount, *launcherReplicasCount, *chiefCount)
	//if err != nil {
	//	return pluginsCore.PhaseInfoUndefined, err
	//}
	//
	//currentCondition, err := common.ExtractCurrentCondition(app.Status.Conditions)
	//if err != nil {
	//	return pluginsCore.PhaseInfoUndefined, err
	//}

	//occurredAt := time.Now()
	//statusDetails, _ := utils.MarshalObjToStruct(app.Status)
	//taskPhaseInfo := pluginsCore.TaskInfo{
	//	//Logs:       taskLogs,
	//	OccurredAt: &occurredAt,
	//	CustomInfo: statusDetails,
	//}

	//return common.GetPhaseInfo(currentCondition, occurredAt, taskPhaseInfo)

}

func init() {
	if err := mpiOp.AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}

	pluginmachinery.PluginRegistry().RegisterK8sPlugin(
		k8s.PluginEntry{
			ID:                  common.MPITaskType,
			RegisteredTaskTypes: []pluginsCore.TaskType{common.MPITaskType},
			ResourceToWatch:     &mpiOp.MPIJob{},
			Plugin:              mpiOperatorResourceHandler{},
			IsDefault:           false,
			DefaultForTaskTypes: []pluginsCore.TaskType{common.MPITaskType},
		})
}
