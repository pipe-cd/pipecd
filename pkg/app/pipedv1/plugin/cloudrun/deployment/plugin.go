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
	"slices"
	"strings"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/cloudrunservice/config"
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

	StageCloudRunSyncDescription = "Deploy the new version and configure all traffic to it"
	StageRollbackDescription     = "Rollback the deployment"
)

func (p *Plugin) FetchDefinedStages() []string {
	return []string{StageCloudRunSync, StageCloudRunPromote, StageRollback}
}

func (p *Plugin) BuildPipelineSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	return &sdk.BuildPipelineSyncStagesResponse{
		Stages: buildPipelineStages(input.Request.Stages, input.Request.Rollback),
	}, nil
}

func (p *Plugin) ExecuteStage(ctx context.Context, _ *sdk.ConfigNone, dts []*sdk.DeployTarget[config.CloudRunDeployTargetConfig], input *sdk.ExecuteStageInput[config.CloudRunApplicationSpec]) (*sdk.ExecuteStageResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (p *Plugin) DetermineVersions(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineVersionsInput[config.CloudRunApplicationSpec]) (*sdk.DetermineVersionsResponse, error) {
	//Getting the image from the spec
	image := input.Request.DeploymentSource.ApplicationConfig.Spec.Input.Image

	// custom ImageVersionExtraction logic to parse the image string 
	version := ImageVersionExtraction(image)

	// Return using the verified SDK fields
	return &sdk.DetermineVersionsResponse{
		Versions: []sdk.ArtifactVersion{
			{
				Version: version,
				Name:    image,
			},
		},
	}, nil
}

func ImageVersionExtraction(image string) string {
	if image == "" {
		return "unknown"
	}

	// Handle digest format: gcr.io/app@sha256:xxx
	if idx := strings.LastIndex(image, "@"); idx != -1 {
		return image[idx+1:]
	}

	// Handle tag format: gcr.io/app:tag
	if idx := strings.LastIndex(image, ":"); idx != -1 {
		// Avoid catching port numbers in registry URL
		slashIdx := strings.LastIndex(image, "/")
		if idx > slashIdx {
			return image[idx+1:]
		}
	}

	return "latest"
}

func (p *Plugin) DetermineStrategy(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineStrategyInput[config.CloudRunApplicationSpec]) (*sdk.DetermineStrategyResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (p *Plugin) BuildQuickSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	return &sdk.BuildQuickSyncStagesResponse{
		Stages: buildQuickSyncPipeline(input.Request.Rollback),
	}, nil
}

func buildQuickSyncPipeline(autoRollback bool) []sdk.QuickSyncStage {
	out := make([]sdk.QuickSyncStage, 0, 2)
	out = append(out, sdk.QuickSyncStage{
		Name:               StageCloudRunSync,
		Description:        StageCloudRunSyncDescription,
		Rollback:           false,
		Metadata:           map[string]string{},
		AvailableOperation: sdk.ManualOperationNone,
	})
	if autoRollback {
		out = append(out, sdk.QuickSyncStage{
			Name:               StageRollback,
			Description:        StageRollbackDescription,
			Rollback:           true,
			Metadata:           map[string]string{},
			AvailableOperation: sdk.ManualOperationNone,
		})
	}
	return out
}

func buildPipelineStages(stages []sdk.StageConfig, autoRollback bool) []sdk.PipelineStage {
	out := make([]sdk.PipelineStage, 0, len(stages)+1)
	for _, stage := range stages {
		out = append(out, sdk.PipelineStage{
			Name:               stage.Name,
			Index:              stage.Index,
			Rollback:           false,
			Metadata:           map[string]string{},
			AvailableOperation: sdk.ManualOperationNone,
		})
	}
	if autoRollback {
		out = append(out, sdk.PipelineStage{
			Name: StageRollback,
			Index: slices.MinFunc(stages, func(a, b sdk.StageConfig) int {
				return a.Index - b.Index
			}).Index,
			Rollback:           true,
			Metadata:           map[string]string{},
			AvailableOperation: sdk.ManualOperationNone,
		})
	}
	return out
}
