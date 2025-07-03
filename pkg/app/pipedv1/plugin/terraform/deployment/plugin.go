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

package deployment

import (
	"context"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/config"
)

// Plugin implements sdk.DeploymentPlugin for Terraform.
type Plugin struct{}

var _ sdk.DeploymentPlugin[config.Config, config.DeployTargetConfig, config.ApplicationConfigSpec] = (*Plugin)(nil)

// BuildPipelineSyncStages implements sdk.DeploymentPlugin.
func (p *Plugin) BuildPipelineSyncStages(ctx context.Context, _ *config.Config, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	panic("unimplemented")
}

// BuildQuickSyncStages implements sdk.DeploymentPlugin.
func (p *Plugin) BuildQuickSyncStages(ctx context.Context, _ *config.Config, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	panic("unimplemented")
}

// DetermineStrategy implements sdk.DeploymentPlugin.
func (p *Plugin) DetermineStrategy(ctx context.Context, _ *config.Config, input *sdk.DetermineStrategyInput[config.ApplicationConfigSpec]) (*sdk.DetermineStrategyResponse, error) {
	panic("unimplemented")
}

// DetermineVersions implements sdk.DeploymentPlugin.
func (p *Plugin) DetermineVersions(ctx context.Context, _ *config.Config, input *sdk.DetermineVersionsInput[config.ApplicationConfigSpec]) (*sdk.DetermineVersionsResponse, error) {
	panic("unimplemented")
}

// ExecuteStage implements sdk.DeploymentPlugin.
func (p *Plugin) ExecuteStage(ctx context.Context, _ *config.Config, dts []*sdk.DeployTarget[config.DeployTargetConfig], input *sdk.ExecuteStageInput[config.ApplicationConfigSpec]) (*sdk.ExecuteStageResponse, error) {
	panic("unimplemented")
}

// FetchDefinedStages implements sdk.DeploymentPlugin.
func (p *Plugin) FetchDefinedStages() []string {
	panic("unimplemented")
}
