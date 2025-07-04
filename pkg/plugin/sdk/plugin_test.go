// Copyright 2025 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sdk

import (
	"context"
	"log"
)

var (
	_ StagePlugin[ExampleConfig, ExampleDeployTargetConfig, ExampleApplicationConfigSpec]       = ExampleStagePlugin{}
	_ DeploymentPlugin[ExampleConfig, ExampleDeployTargetConfig, ExampleApplicationConfigSpec]  = ExampleDeploymentPlugin{}
	_ LivestatePlugin[ExampleConfig, ExampleDeployTargetConfig, ExampleApplicationConfigSpec]   = ExampleLivestatePlugin{}
	_ PlanPreviewPlugin[ExampleConfig, ExampleDeployTargetConfig, ExampleApplicationConfigSpec] = ExamplePlanPreviewPlugin{}
)

type (
	ExampleStagePlugin           struct{}
	ExampleDeploymentPlugin      struct{}
	ExampleLivestatePlugin       struct{}
	ExamplePlanPreviewPlugin     struct{}
	ExampleConfig                struct{}
	ExampleDeployTargetConfig    struct{}
	ExampleApplicationConfigSpec struct{}
)

// BuildPipelineSyncStages implements StagePlugin.
func (e ExampleStagePlugin) BuildPipelineSyncStages(context.Context, *ExampleConfig, *BuildPipelineSyncStagesInput) (*BuildPipelineSyncStagesResponse, error) {
	panic("unimplemented")
}

// ExecuteStage implements StagePlugin.
func (e ExampleStagePlugin) ExecuteStage(context.Context, *ExampleConfig, []*DeployTarget[ExampleDeployTargetConfig], *ExecuteStageInput[ExampleApplicationConfigSpec]) (*ExecuteStageResponse, error) {
	panic("unimplemented")
}

// FetchDefinedStages implements StagePlugin.
func (e ExampleStagePlugin) FetchDefinedStages() []string {
	panic("unimplemented")
}

// BuildPipelineSyncStages implements DeploymentPlugin.
func (e ExampleDeploymentPlugin) BuildPipelineSyncStages(context.Context, *ExampleConfig, *BuildPipelineSyncStagesInput) (*BuildPipelineSyncStagesResponse, error) {
	panic("unimplemented")
}

// BuildQuickSyncStages implements DeploymentPlugin.
func (e ExampleDeploymentPlugin) BuildQuickSyncStages(context.Context, *ExampleConfig, *BuildQuickSyncStagesInput) (*BuildQuickSyncStagesResponse, error) {
	panic("unimplemented")
}

// DetermineStrategy implements DeploymentPlugin.
func (e ExampleDeploymentPlugin) DetermineStrategy(context.Context, *ExampleConfig, *DetermineStrategyInput[ExampleApplicationConfigSpec]) (*DetermineStrategyResponse, error) {
	panic("unimplemented")
}

// DetermineVersions implements DeploymentPlugin.
func (e ExampleDeploymentPlugin) DetermineVersions(context.Context, *ExampleConfig, *DetermineVersionsInput[ExampleApplicationConfigSpec]) (*DetermineVersionsResponse, error) {
	panic("unimplemented")
}

// ExecuteStage implements DeploymentPlugin.
func (e ExampleDeploymentPlugin) ExecuteStage(context.Context, *ExampleConfig, []*DeployTarget[ExampleDeployTargetConfig], *ExecuteStageInput[ExampleApplicationConfigSpec]) (*ExecuteStageResponse, error) {
	panic("unimplemented")
}

// FetchDefinedStages implements DeploymentPlugin.
func (e ExampleDeploymentPlugin) FetchDefinedStages() []string {
	panic("unimplemented")
}

// GetLivestate implements LivestatePlugin.
func (e ExampleLivestatePlugin) GetLivestate(context.Context, *ExampleConfig, []*DeployTarget[ExampleDeployTargetConfig], *GetLivestateInput[ExampleApplicationConfigSpec]) (*GetLivestateResponse, error) {
	panic("unimplemented")
}

// GetPlanPreview implements PlanPreviewPlugin.
func (e ExamplePlanPreviewPlugin) GetPlanPreview(context.Context, *ExampleConfig, []*DeployTarget[ExampleDeployTargetConfig], *GetPlanPreviewInput[ExampleApplicationConfigSpec]) (*GetPlanPreviewResponse, error) {
	panic("unimplemented")
}

func ExampleNewPlugin() {
	plugin, err := NewPlugin("1.0.0",
		WithDeploymentPlugin(ExampleDeploymentPlugin{}),
		WithLivestatePlugin(ExampleLivestatePlugin{}),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Run runs the plugin and blocks until the signal is received.
	// So you can use it like this:
	/*
		if err := plugin.Run(); err != nil {
			log.Fatal(err)
		}
	*/

	_ = plugin
}

func ExampleWithStagePlugin() {
	plugin, err := NewPlugin("1.0.0",
		WithStagePlugin(ExampleStagePlugin{}),
	)
	if err != nil {
		log.Fatal(err)
	}

	// plugin.Run()
	_ = plugin
}

func ExampleWithDeploymentPlugin() {
	plugin, err := NewPlugin("1.0.0",
		WithDeploymentPlugin(ExampleDeploymentPlugin{}),
	)
	if err != nil {
		log.Fatal(err)
	}

	// plugin.Run()
	_ = plugin
}

func ExampleWithLivestatePlugin() {
	plugin, err := NewPlugin("1.0.0",
		WithLivestatePlugin(ExampleLivestatePlugin{}),
	)
	if err != nil {
		log.Fatal(err)
	}

	// plugin.Run()
	_ = plugin
}

func ExampleWithPlanPreviewPlugin() {
	plugin, err := NewPlugin("1.0.0",
		WithPlanPreviewPlugin(ExamplePlanPreviewPlugin{}),
	)
	if err != nil {
		log.Fatal(err)
	}

	// plugin.Run()
	_ = plugin
}
