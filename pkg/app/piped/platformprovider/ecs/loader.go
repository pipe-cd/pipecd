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

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/git"
)

type Loader interface {
	LoadManifests(ctx context.Context) ([]Manifest, error)
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

func (l *loader) LoadManifests(ctx context.Context) ([]Manifest, error) {
	var ecsInput config.ECSDeploymentInput = l.input
	manifestNames := []string{ecsInput.ServiceDefinitionFile, ecsInput.TaskDefinitionFile}
	manifests, err := LoadPlainYAMLManifests(l.appDir, manifestNames, l.configFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to load manifest(%w)", err)
	}
	return manifests, nil
}
