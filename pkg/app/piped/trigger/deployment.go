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

func (t *Trigger) triggerDeployment(ctx context.Context, app *model.Application, branch string, commit git.Commit, commander string) (runErr error) {
	deployment, err := buildDeployment(app, branch, commit, commander, time.Now())
	if err != nil {
		return err
	}

	defer func() {
		if runErr != nil {
			return
		}
		var envName string
		if env, ok := t.environmentLister.Get(deployment.EnvId); ok {
			envName = env.Name
		}
		t.notifier.Notify(model.Event{
			Type: model.EventType_EVENT_DEPLOYMENT_TRIGGERED,
			Metadata: &model.EventDeploymentTriggered{
				Deployment: deployment,
				EnvName:    envName,
			},
		})
	}()

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

	if err := t.reportMostRecentlyTriggeredDeployment(ctx, deployment); err != nil {
		t.logger.Error("failed to report most recently triggered deployment", zap.Error(err))
	}

	return nil
}

func (t *Trigger) loadDeploymentConfiguration(repoPath string, app *model.Application) (*config.Config, error) {
	path := filepath.Join(repoPath, app.GitPath.GetDeploymentConfigFilePath())
	cfg, err := config.LoadFromYAML(path)
	if err != nil {
		return nil, err
	}
	if appKind, ok := config.ToApplicationKind(cfg.Kind); !ok || appKind != app.Kind {
		return nil, fmt.Errorf("application in deployment configuration file is not match, got: %s, expected: %s", appKind, app.Kind)
	}
	return cfg, nil
}

func (t *Trigger) reportMostRecentlyTriggeredDeployment(ctx context.Context, d *model.Deployment) error {
	var (
		err error
		req = &pipedservice.ReportApplicationMostRecentDeploymentRequest{
			ApplicationId: d.ApplicationId,
			Status:        model.DeploymentStatus_DEPLOYMENT_PENDING,
			Deployment: &model.ApplicationDeploymentReference{
				DeploymentId: d.Id,
				Trigger:      d.Trigger,
				Summary:      d.Summary,
				Version:      d.Version,
				StartedAt:    d.CreatedAt,
				CompletedAt:  d.CompletedAt,
			},
		}
		retry = pipedservice.NewRetry(10)
	)

	for retry.WaitNext(ctx) {
		if _, err = t.apiClient.ReportApplicationMostRecentDeployment(ctx, req); err == nil {
			return nil
		}
		err = fmt.Errorf("failed to report most recent successful deployment: %w", err)
	}
	return err
}

func buildDeployment(app *model.Application, branch string, commit git.Commit, commander string, now time.Time) (*model.Deployment, error) {
	commitURL := ""
	if r := app.GitPath.Repo; r != nil {
		var err error
		commitURL, err = git.MakeCommitURL(r.Remote, commit.Hash)
		if err != nil {
			return nil, err
		}
	}

	deployment := &model.Deployment{
		Id:              uuid.New().String(),
		ApplicationId:   app.Id,
		ApplicationName: app.Name,
		EnvId:           app.EnvId,
		PipedId:         app.PipedId,
		ProjectId:       app.ProjectId,
		Kind:            app.Kind,
		Trigger: &model.DeploymentTrigger{
			Commit: &model.Commit{
				Hash:      commit.Hash,
				Message:   commit.Message,
				Author:    commit.Author,
				Branch:    branch,
				Url:       commitURL,
				CreatedAt: int64(commit.CreatedAt),
			},
			Commander: commander,
			Timestamp: now.Unix(),
		},
		GitPath:       app.GitPath,
		CloudProvider: app.CloudProvider,
		Status:        model.DeploymentStatus_DEPLOYMENT_PENDING,
		StatusReason:  "The deployment is waiting to be planned",
		CreatedAt:     now.Unix(),
		UpdatedAt:     now.Unix(),
	}

	return deployment, nil
}
