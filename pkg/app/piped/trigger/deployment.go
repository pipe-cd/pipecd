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

package trigger

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/git"
	"github.com/pipe-cd/pipe/pkg/model"
)

func (t *Trigger) triggerDeployment(ctx context.Context, app *model.Application, repo git.Repo, branch string, commit git.Commit) error {
	// Load deployment configuration at the commit.
	cfg, err := t.loadDeploymentConfiguration(repo.GetPath(), app)
	if err != nil {
		t.logger.Error("failed to load application configuration",
			zap.String("commit-hash", commit.Hash),
			zap.Error(err),
		)
		return err
	}

	deployment, err := buildDeploment(app, cfg, branch, commit, time.Now())
	if err != nil {
		return err
	}

	t.logger.Info(fmt.Sprintf("application %s will be triggered to sync", app.Id),
		zap.String("commit-hash", commit.Hash),
	)
	req := &pipedservice.CreateDeploymentRequest{
		Deployment: deployment,
	}
	if _, err = t.apiClient.CreateDeployment(ctx, req); err != nil {
		t.logger.Error("failed to create deployment", zap.Error(err))
		return err
	}

	return nil
}

func (t *Trigger) loadDeploymentConfiguration(repoPath string, app *model.Application) (*config.Config, error) {
	path := filepath.Join(repoPath, app.GitPath.GetDeploymentConfigFilePath(config.DeploymentConfigurationFileName))
	cfg, err := config.LoadFromYAML(path)
	if err != nil {
		return nil, err
	}
	if appKind, ok := config.ToApplicationKind(cfg.Kind); !ok || appKind != app.Kind {
		return nil, fmt.Errorf("application in deployment configuration file is not match, got: %s, expected: %s", appKind, app.Kind)
	}
	return cfg, nil
}

func buildDeploment(app *model.Application, cfg *config.Config, branch string, commit git.Commit, now time.Time) (*model.Deployment, error) {
	deployment := &model.Deployment{
		Id:            uuid.New().String(),
		ApplicationId: app.Id,
		EnvId:         app.EnvId,
		PipedId:       app.PipedId,
		ProjectId:     app.ProjectId,
		Kind:          app.Kind,
		Trigger: &model.DeploymentTrigger{
			Commit: &model.Commit{
				Hash:      commit.Hash,
				Message:   commit.Message,
				Author:    commit.Author,
				Branch:    branch,
				CreatedAt: int64(commit.CreatedAt),
			},
			User:      commit.Author,
			Timestamp: now.Unix(),
		},
		GitPath:       app.GitPath,
		CloudProvider: app.CloudProvider,
		Status:        model.DeploymentStatus_DEPLOYMENT_PENDING,
		CreatedAt:     now.Unix(),
		UpdatedAt:     now.Unix(),
	}

	return deployment, nil
}
