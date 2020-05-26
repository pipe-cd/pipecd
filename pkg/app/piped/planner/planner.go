// Copyright 2020 The PipeCD Authors.
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

// Package planner provides a piped component
// that decides the deployment pipeline of a deployment.
// The planner bases on the changes from git commits
// then builds the deployment manifests to know the behavior of the deployment.
// From that behavior the planner can decides which pipeline should be applied.
package planner

import (
	"context"

	"go.uber.org/zap"

	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type Planner interface {
	Plan(ctx context.Context, in Input) (Output, error)
}

type Input struct {
	// Readonly deployment model.
	Deployment               *model.Deployment
	LastSuccessfulCommitHash string
	DeploymentConfig         *config.Config
	RepoDir                  string
	AppDir                   string
	Logger                   *zap.Logger
}

type Output struct {
	Stages      []*model.PipelineStage
	Description string
}
