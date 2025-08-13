package deployment

import (
	"context"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/cloudrun/config"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

type Plugin struct{}

const (
	// StageCloudRunSync does quick sync by rolling out the new version
	// and switching all traffic to it.
	StageCloudRunSync = "CLOUDRUN_SYNC"
	// StageCloudRunPromote promotes the new version to receive amount of traffic.
	StageCloudRunPromote = "CLOUDRUN_PROMOTE"
	// StageRollback the legacy generic rollback stage name
	StageRollback = "ROLLBACK"
)

func (p Plugin) FetchDefinedStages() []string {
	return []string{StageCloudRunSync, StageCloudRunPromote, StageRollback}
}

func (p Plugin) BuildPipelineSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p Plugin) ExecuteStage(ctx context.Context, _ *sdk.ConfigNone, dts []*sdk.DeployTarget[config.CloudRunDeployTargetConfig], input *sdk.ExecuteStageInput[config.CloudRunApplicationSpec]) (*sdk.ExecuteStageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p Plugin) DetermineVersions(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineVersionsInput[config.CloudRunApplicationSpec]) (*sdk.DetermineVersionsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p Plugin) DetermineStrategy(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineStrategyInput[config.CloudRunApplicationSpec]) (*sdk.DetermineStrategyResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p Plugin) BuildQuickSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	//TODO implement me
	panic("implement me")
}
