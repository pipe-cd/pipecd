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

package deploymentcontroller

import (
	"context"

	"go.uber.org/zap"

	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type repoStore interface {
	CloneReadOnlyRepo(repo, branch, revision string) (string, error)
}

// scheduler is a dedicated object for a specific deployment of a single application.
type scheduler struct {
	deployment *model.Deployment
	// Deployment configuration for this application.
	appConfig  *config.Config
	workingDir string
	logger     *zap.Logger
}

func newExecutor(d *model.Deployment, logger *zap.Logger) *scheduler {
	logger = logger.Named("scheduler").With(
		zap.String("deployment-id", d.Id),
		zap.String("application-id", d.ApplicationId),
		zap.String("env-id", d.EnvId),
		zap.String("project-id", d.ProjectId),
		zap.String("application-kind", d.Kind.String()),
	)
	return &scheduler{
		deployment: d,
		logger:     logger,
	}
}

func (e *scheduler) Id() string {
	return e.deployment.Id
}

func (e *scheduler) IsCompleted() bool {
	return false
}

func (e *scheduler) IsDone() bool {
	return false
}

func (e *scheduler) Run(ctx context.Context) error {
	// Prepare a working space for this deployment.
	// Load deployment configuration data.
	// Restore previous executed state.
	// Start executing the next stages.
	return nil
}

// prepare does all needed things before start executing the deployment.
// Includes:
// - Clone a readonly repository at the required revision
// - Restore previous executed state from deployment data.
func (e *scheduler) prepare(ctx context.Context) error {
	return nil
}

func (e *scheduler) run(ctx context.Context) error {
	return nil
}
