// Copyright 2024 The PipeCD Authors.
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

package initialize

import (
	"github.com/pipe-cd/pipecd/pkg/app/pipectl/prompt"
	"github.com/pipe-cd/pipecd/pkg/config"
)

// Use genericConfigs in order to simplify using the GenericApplicationSpec and keep the order as we want.
type genericECSApplicationSpec struct {
	Name        string                    `json:"name"`
	Input       config.ECSDeploymentInput `json:"input"`
	Description string                    `json:"description,omitempty"`
}

func generateECSConfig(p prompt.Prompt) (*genericConfig, error) {
	// inputs
	var (
		appName        string
		serviceDefFile string
		taskDefFile    string
		targetGroupArn string
		containerName  string
		containerPort  int
	)
	inputs := []prompt.Input{
		{
			Message:       "Name of the application",
			TargetPointer: &appName,
			Required:      true,
		},
		{
			Message:       "Name of the service definition file (e.g. serviceDef.yaml)",
			TargetPointer: &serviceDefFile,
			Required:      true,
		},
		{
			Message:       "Name of the task definition file (e.g. taskDef.yaml)",
			TargetPointer: &taskDefFile,
			Required:      true,
		},
		// target group inputs
		{
			Message:       "ARN of the target group to the service",
			TargetPointer: &targetGroupArn,
			Required:      false,
		},
		{
			Message:       "Name of the container of the target group",
			TargetPointer: &containerName,
			Required:      false,
		},
		{
			Message:       "Port number of the container of the target group",
			TargetPointer: &containerPort,
			Required:      false,
		},
	}

	err := p.Run(inputs)
	if err != nil {
		return nil, err
	}

	spec := &genericECSApplicationSpec{
		Name: appName,
		Input: config.ECSDeploymentInput{
			ServiceDefinitionFile: serviceDefFile,
			TaskDefinitionFile:    taskDefFile,
			TargetGroups: config.ECSTargetGroups{
				Primary: &config.ECSTargetGroup{
					TargetGroupArn: targetGroupArn,
					ContainerName:  containerName,
					ContainerPort:  containerPort,
				},
			},
		},
		Description: "Generated by `pipectl init`. See https://pipecd.dev/docs/user-guide/configuration-reference/ for more.",
	}

	return &genericConfig{
		Kind:            config.KindECSApp,
		APIVersion:      config.VersionV1Beta1,
		ApplicationSpec: spec,
	}, nil
}
