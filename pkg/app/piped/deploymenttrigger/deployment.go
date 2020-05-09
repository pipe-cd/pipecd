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

package deploymenttrigger

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/kapetaniosci/pipe/pkg/app/api/service/pipedservice"
	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/git"
	"github.com/kapetaniosci/pipe/pkg/model"
)

func (t *DeploymentTrigger) triggerDeployment(ctx context.Context, app *model.Application, repo git.Repo, branch string, commit git.Commit) error {
	// Load deployment configuration at the commit.
	cfg, err := t.loadDeploymentConfiguration(ctx, repo.GetPath(), app)
	if err != nil {
		t.logger.Error("failed to load application configuration",
			zap.String("commit-hash", commit.Hash),
			zap.Error(err),
		)
		return err
	}

	// TODO: Detect the update type (just scale or need rollout with pipeline) by checking the change
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

func (t *DeploymentTrigger) loadDeploymentConfiguration(ctx context.Context, repoPath string, app *model.Application) (*config.Config, error) {
	path := filepath.Join(repoPath, app.GetDeploymentConfigFilePath(config.DeploymentConfigurationFileName))
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
	var (
		stages []*model.PipelineStage
		err    error
	)

	switch cfg.Kind {
	case config.KindKubernetesApp:
		stages, err = buildKubernetesPipelineStages(cfg, now)
	case config.KindTerraformApp:
		stages, err = buildTerraformPipelineStages(cfg, now)
	default:
		err = fmt.Errorf("unsupported application kind: %s", cfg.Kind)
	}
	if err != nil {
		return nil, err
	}

	deployment := &model.Deployment{
		Id:            uuid.New().String(),
		ApplicationId: app.Id,
		EnvId:         app.EnvId,
		PipedId:       app.PipedId,
		ProjectId:     app.ProjectId,
		Kind:          app.Kind,
		Trigger: &model.DeploymentTrigger{
			Commit: &model.Commit{
				Revision:  commit.Hash,
				Message:   commit.Message,
				Author:    commit.Author,
				Branch:    branch,
				CreatedAt: int64(commit.CreatedAt),
			},
			User:      commit.Author,
			Timestamp: now.Unix(),
		},
		GitPath:   app.GitPath,
		Status:    model.DeploymentStatus_DEPLOYMENT_NOT_STARTED_YET,
		Stages:    stages,
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	}
	return deployment, nil
}

var kubernetesDefaultPipeline = &config.AppPipeline{
	Stages: []config.PipelineStage{
		{
			Name: model.StageK8sPrimaryUpdate,
			Desc: "Update primary to new version",
		},
	},
}

func buildKubernetesPipelineStages(cfg *config.Config, now time.Time) ([]*model.PipelineStage, error) {
	p := cfg.KubernetesAppSpec.Pipeline
	if p == nil {
		p = kubernetesDefaultPipeline
	}
	stages := make([]*model.PipelineStage, 0, len(p.Stages))
	for i, s := range p.Stages {
		id := s.Id
		if id == "" {
			id = fmt.Sprintf("stage-%d", i)
		}
		stage := &model.PipelineStage{
			Id:        id,
			Name:      s.Name.String(),
			Desc:      s.Desc,
			Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
			CreatedAt: now.Unix(),
			UpdatedAt: now.Unix(),
		}
		stages = append(stages, stage)
	}
	return stages, nil
}

var terraformDefaultPipeline = &config.AppPipeline{
	Stages: []config.PipelineStage{
		{
			Name: model.StageTerraformPlan,
			Desc: "Terraform Plan",
		},
		{
			Name: model.StageWaitApproval,
			Desc: "Wait for an approval",
		},
		{
			Name: model.StageTerraformApply,
			Desc: "Terraform Apply",
		},
	},
}

func buildTerraformPipelineStages(cfg *config.Config, now time.Time) ([]*model.PipelineStage, error) {
	p := cfg.TerraformAppSpec.Pipeline
	if p == nil {
		p = terraformDefaultPipeline
	}
	stages := make([]*model.PipelineStage, 0, len(p.Stages))
	for i, s := range p.Stages {
		id := s.Id
		if id == "" {
			id = fmt.Sprintf("stage-%d", i)
		}
		stage := &model.PipelineStage{
			Id:        id,
			Name:      s.Name.String(),
			Desc:      s.Desc,
			Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
			CreatedAt: now.Unix(),
			UpdatedAt: now.Unix(),
		}
		stages = append(stages, stage)
	}
	return stages, nil
}
