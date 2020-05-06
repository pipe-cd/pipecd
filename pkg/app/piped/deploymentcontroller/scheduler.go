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

	"github.com/kapetaniosci/pipe/pkg/app/piped/executor"
	"github.com/kapetaniosci/pipe/pkg/app/piped/logpersister"
	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type repoStore interface {
	CloneReadOnlyRepo(repo, branch, revision string) (string, error)
}

type executorRegistry interface {
	Executor(stage model.Stage, in executor.Input) (executor.Executor, error)
}

// scheduler is a dedicated object for a specific deployment of a single application.
type scheduler struct {
	deployment        *model.Deployment
	pipedConfig       *config.PipedSpec
	workingDir        string
	executorRegistry  executorRegistry
	logPersister      logpersister.Persister
	metadataPersister metadataPersister
	logger            *zap.Logger

	// Deployment configuration for this application.
	appConfig *config.Config
}

func newScheduler(d *model.Deployment, cfg *config.PipedSpec, workingDir string, lp logpersister.Persister, mdp metadataPersister, logger *zap.Logger) *scheduler {
	logger = logger.Named("scheduler").With(
		zap.String("deployment-id", d.Id),
		zap.String("application-id", d.ApplicationId),
		zap.String("env-id", d.EnvId),
		zap.String("project-id", d.ProjectId),
		zap.String("application-kind", d.Kind.String()),
		zap.String("working-dir", workingDir),
	)
	return &scheduler{
		deployment:        d,
		pipedConfig:       cfg,
		workingDir:        workingDir,
		executorRegistry:  executor.DefaultRegistry(),
		logPersister:      lp,
		metadataPersister: mdp,
		logger:            logger,
	}
}

func (s *scheduler) Id() string {
	return s.deployment.Id
}

func (s *scheduler) IsCompleted() bool {
	return false
}

func (s *scheduler) IsDone() bool {
	return false
}

func (s *scheduler) Run(ctx context.Context) error {
	// Prepare a working space for this deployment.
	// Load deployment configuration data.
	// Restore previous executed state.
	// Start executing the next stages.
	s.logger.Info("start running scheduler")

	return nil
}

// prepare does all needed things before start executing the deployment.
// Includes:
// - Clone a readonly repository at the required revision
// - Restore previous executed state from deployment data.
func (s *scheduler) prepare(ctx context.Context) error {
	return nil
}

func (s *scheduler) run(ctx context.Context) error {
	// Loop until one of the following conditions occurs:
	// - context has done
	// - no stage to execute
	// - executing stage has completed with an error
	// Determine the next stage that should be executed.
	var (
		stageName = model.Stage("")
		input     = executor.Input{
			Deployment:        s.deployment,
			AppConfig:         s.appConfig,
			WorkingDir:        s.workingDir,
			LogPersister:      s.logPersister.StageLogPersister("", ""),
			MetadataPersister: s.metadataPersister.StageMetadataPersister("", ""),
			Logger:            s.logger,
		}
	)
	ex, err := s.executorRegistry.Executor(stageName, input)
	if err != nil {
		return nil
	}
	_, err = ex.Execute(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *scheduler) determineNextStages() []string {
	return nil
}
