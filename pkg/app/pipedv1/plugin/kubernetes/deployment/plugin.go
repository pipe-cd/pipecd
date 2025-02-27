package deployment

import (
	"context"
	"errors"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

// Plugin implements the sdk.DeploymentPlugin interface.
type Plugin struct{}

var _ sdk.DeploymentPlugin[sdk.ConfigNone, kubeconfig.KubernetesDeployTargetConfig] = (*Plugin)(nil)

func (p *Plugin) Name() string {
	return "kubernetes"
}

func (p *Plugin) Version() string {
	return "0.0.1" // TODO
}

func (p *Plugin) FetchDefinedStages() []string {
	stages := make([]string, 0, len(AllStages))
	for _, s := range AllStages {
		stages = append(stages, string(s))
	}

	return stages
}

// FIXME
func (p *Plugin) BuildPipelineSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	return nil, nil
}

// FIXME
func (p *Plugin) ExecuteStage(ctx context.Context, _ *sdk.ConfigNone, dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], input *sdk.ExecuteStageInput) (*sdk.ExecuteStageResponse, error) {
	switch input.Request.StageName {
	case StageK8sSync.String():
		return &sdk.ExecuteStageResponse{
			Status: p.executeK8sSyncStage(ctx, input),
		}, nil
	case StageK8sRollback.String():
		return &sdk.ExecuteStageResponse{
			Status: p.executeK8sRollbackStage(ctx, input),
		}, nil
	default:
		return nil, errors.New("unimplemented or unsupported stage")
	}
}

// FIXME
func (p *Plugin) executeK8sSyncStage(ctx context.Context, input *sdk.ExecuteStageInput) sdk.StageStatus {
	return sdk.StageStatusFailure
}

// FIXME
func (p *Plugin) executeK8sRollbackStage(ctx context.Context, input *sdk.ExecuteStageInput) sdk.StageStatus {
	return sdk.StageStatusFailure
}

// FIXME
func (p *Plugin) DetermineVersions(context.Context, *sdk.ConfigNone, *sdk.Client, sdk.TODO) (sdk.TODO, error) {
	return sdk.TODO{}, nil
}

// FIXME
func (p *Plugin) DetermineStrategy(context.Context, *sdk.ConfigNone, *sdk.Client, sdk.TODO) (sdk.TODO, error) {
	return sdk.TODO{}, nil
}

func (p *Plugin) BuildQuickSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	return &sdk.BuildQuickSyncStagesResponse{
		Stages: BuildQuickSyncPipeline(input.Request.Rollback),
	}, nil
}
