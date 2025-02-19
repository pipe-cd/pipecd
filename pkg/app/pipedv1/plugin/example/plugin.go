package main

import (
	"context"

	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

type plugin struct{}

type config struct{}

// Name implements sdk.Plugin.
func (p *plugin) Name() string {
	return "example"
}

// Version implements sdk.Plugin.
func (p *plugin) Version() string {
	return "0.0.1"
}

// BuildPipelineSyncStages implements sdk.PipelineSyncPlugin.
func (p *plugin) BuildPipelineSyncStages(context.Context, *config, *sdk.Client, sdk.TODO) (sdk.TODO, error) {
	return sdk.TODO{}, nil
}

// ExecuteStage implements sdk.PipelineSyncPlugin.
func (p *plugin) ExecuteStage(context.Context, *config, sdk.DeployTargetsNone, *sdk.Client, logpersister.StageLogPersister, sdk.TODO) (sdk.TODO, error) {
	return sdk.TODO{}, nil
}

// FetchDefinedStages implements sdk.PipelineSyncPlugin.
func (p *plugin) FetchDefinedStages() []string {
	return []string{"EXAMPLE_PLAN", "EXAMPLE_APPLY"}
}
