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

package main

import (
	"context"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

type plugin struct{}

type config struct{}

type deployTargetConfig struct{}

var _ sdk.StagePlugin[config, deployTargetConfig] = (*plugin)(nil)

// Name implements sdk.Plugin.
func (p *plugin) Name() string {
	return "kubernetes_multicluster"
}

// Version implements sdk.Plugin.
func (p *plugin) Version() string {
	return "0.0.1"
}

// BuildPipelineSyncStages implements sdk.StagePlugin.
func (p *plugin) BuildPipelineSyncStages(ctx context.Context, _ *config, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	stages := deployment.BuildPipelineStages(input)

	return &sdk.BuildPipelineSyncStagesResponse{
		Stages: stages,
	}, nil
}

// ExecuteStage implements sdk.StagePlugin.
func (p *plugin) ExecuteStage(context.Context, *config, []*sdk.DeployTarget[deployTargetConfig], *sdk.ExecuteStageInput) (*sdk.ExecuteStageResponse, error) {
	return &sdk.ExecuteStageResponse{}, nil
}

// FetchDefinedStages implements sdk.StagePlugin.
func (p *plugin) FetchDefinedStages() []string {
	return []string{"K8S_SYNC"}
}
