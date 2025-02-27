package deployment

import (
	"context"

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
func (p *Plugin) ExecuteStage(context.Context, *sdk.ConfigNone, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], *sdk.ExecuteStageInput) (*sdk.ExecuteStageResponse, error) {
	return nil, nil
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
