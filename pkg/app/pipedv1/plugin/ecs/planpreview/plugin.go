// Copyright 2026 The PipeCD Authors.
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

package planpreview

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"sigs.k8s.io/yaml"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/config"
)

var (
	_ sdk.PlanPreviewPlugin[config.ECSPluginConfig, config.ECSDeployTargetConfig, config.ECSApplicationSpec] = (*Plugin)(nil)
)

// Plugin implements the PlanPreview feature for the ECS plugin
type Plugin struct{}

func (p *Plugin) GetPlanPreview(
	ctx context.Context,
	_ *config.ECSPluginConfig,
	dts []*sdk.DeployTarget[config.ECSDeployTargetConfig],
	input *sdk.GetPlanPreviewInput[config.ECSApplicationSpec],
) (*sdk.GetPlanPreviewResponse, error) {
	targetDS := input.Request.TargetDeploymentSource
	targetAppCfg, err := targetDS.AppConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load target app config: %w", err)
	}
	targetInput := targetAppCfg.Spec.Input

	targetTaskDef, err := loadTaskDef(targetDS.ApplicationDirectory, targetInput.TaskDefinitionFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load target task definition: %w", err)
	}

	runningDS := input.Request.RunningDeploymentSource
	var runningTaskDef *types.TaskDefinition

	if runningDS.CommitHash != "" {
		runningAppCfg, err := runningDS.AppConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to load running app config: %w", err)
		}
		td, err := loadTaskDef(runningDS.ApplicationDirectory, runningAppCfg.Spec.Input.TaskDefinitionFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load running task definition: %w", err)
		}
		runningTaskDef = &td
	}

	taskDefDiff, err := diffDefinitions(runningTaskDef, &targetTaskDef, "taskdef")
	if err != nil {
		return nil, fmt.Errorf("failed to diff task definitions: %w", err)
	}

	var serviceDiff string
	if targetInput.ServiceDefinitionFile != "" {
		targetServiceDef, err := loadServiceDef(targetDS.ApplicationDirectory, targetInput.ServiceDefinitionFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load target service definition: %w", err)
		}

		var runningServiceDef *types.Service
		if runningDS.CommitHash != "" {
			runningAppCfg, err := runningDS.AppConfig()
			if err != nil {
				return nil, fmt.Errorf("failed to load running app config: %w", err)
			}
			if runningAppCfg.Spec.Input.ServiceDefinitionFile != "" {
				sd, err := loadServiceDef(runningDS.ApplicationDirectory, runningAppCfg.Spec.Input.ServiceDefinitionFile)
				if err != nil {
					return nil, fmt.Errorf("failed to load running service definition: %w", err)
				}
				runningServiceDef = &sd
			}
		}

		serviceDiff, err = diffDefinitions(runningServiceDef, &targetServiceDef, "servicedef")
		if err != nil {
			return nil, fmt.Errorf("failed to diff service definitions: %w", err)
		}
	}

	return toResponse(dts[0].Name, taskDefDiff, serviceDiff), nil
}

func loadTaskDef(appDir, filename string) (types.TaskDefinition, error) {
	data, err := os.ReadFile(filepath.Join(appDir, filename))
	if err != nil {
		return types.TaskDefinition{}, err
	}
	var obj types.TaskDefinition
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return types.TaskDefinition{}, err
	}
	return obj, nil
}

func loadServiceDef(appDir, filename string) (types.Service, error) {
	data, err := os.ReadFile(filepath.Join(appDir, filename))
	if err != nil {
		return types.Service{}, err
	}
	var obj types.Service
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return types.Service{}, err
	}
	return obj, nil
}

