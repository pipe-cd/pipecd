// Copyright 2022 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ecs

import (
	"context"
	"fmt"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/git"
)

type Loader interface {
	LoadDefinitionFiles(ctx context.Context) (types.Service, types.TaskDefinition, error)
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type loader struct {
	appName        string
	appDir         string
	repoDir        string
	configFileName string
	input          config.ECSDeploymentInput
	gc             gitClient
	logger         *zap.Logger
}

func NewLoader(
	appName, appDir, repoDir, configFileName string,
	input config.ECSDeploymentInput,
	gc gitClient,
	logger *zap.Logger,
) Loader {
	return &loader{
		appName:        appName,
		appDir:         appDir,
		repoDir:        repoDir,
		configFileName: configFileName,
		input:          input,
		gc:             gc,
		logger:         logger.Named("ecs-loader"),
	}
}

func (l *loader) LoadDefinitionFiles(ctx context.Context) (types.Service, types.TaskDefinition, error) {
	var ecsInput config.ECSDeploymentInput = l.input
	serviceDefinitionFilePath := filepath.Join(l.appDir, ecsInput.ServiceDefinitionFile)
	taskDefinitionFilePath := filepath.Join(l.appDir, ecsInput.TaskDefinitionFile)
	serviceDefinition, err := loadServiceDefinition(serviceDefinitionFilePath)
	if err != nil {
		return types.Service{}, types.TaskDefinition{}, fmt.Errorf("failed to load service definition %s (%w)", serviceDefinitionFilePath, err)
	}
	taskDefinition, err := loadTaskDefinition(taskDefinitionFilePath)
	if err != nil {
		return types.Service{}, types.TaskDefinition{}, fmt.Errorf("failed to load task definition %s (%w)", serviceDefinitionFilePath, err)
	}
	return serviceDefinition, taskDefinition, nil
}
