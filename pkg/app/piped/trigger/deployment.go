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

package trigger

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func (t *Trigger) triggerDeployment(
	ctx context.Context,
	deployment *model.Deployment,
) error {
	if _, err := t.apiClient.CreateDeployment(ctx, &pipedservice.CreateDeploymentRequest{
		Deployment: deployment,
	}); err != nil {
		return fmt.Errorf("cound not register a new deployment to control-plane: %w", err)
	}
	return nil
}

func buildDeployment(
	app *model.Application,
	branch string,
	commit git.Commit,
	commander string,
	syncStrategy model.SyncStrategy,
	strategySummary string,
	now time.Time,
	noti *config.DeploymentNotification,
	deploymentChainID string,
	deploymentChainBlockIndex uint32,
) (*model.Deployment, error) {

	var commitURL string
	if r := app.GitPath.Repo; r != nil {
		url, err := git.MakeCommitURL(r.Remote, commit.Hash)
		if err != nil {
			return nil, err
		}
		commitURL = url
	}

	metadata := make(map[string]string)
	if noti != nil {
		value, err := json.Marshal(noti)
		if err != nil {
			return nil, fmt.Errorf("failed to save notification config to deployment metadata: %w", err)
		}
		metadata[model.MetadataKeyDeploymentNotification] = string(value)
	}

	deployment := &model.Deployment{
		Id:              uuid.New().String(),
		ApplicationId:   app.Id,
		ApplicationName: app.Name,
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
			Commander:       commander,
			Timestamp:       now.Unix(),
			SyncStrategy:    syncStrategy,
			StrategySummary: strategySummary,
		},
		GitPath:                   app.GitPath,
		CloudProvider:             app.CloudProvider,
		PlatformProvider:          app.PlatformProvider,
		DeployTargetsByPlugin:     app.GetDeployTargetsByPlugin(),
		Labels:                    app.Labels,
		Status:                    model.DeploymentStatus_DEPLOYMENT_PENDING,
		StatusReason:              "The deployment is waiting to be planned",
		Metadata:                  metadata,
		CreatedAt:                 now.Unix(),
		UpdatedAt:                 now.Unix(),
		DeploymentChainId:         deploymentChainID,
		DeploymentChainBlockIndex: deploymentChainBlockIndex,
	}

	return deployment, nil
}

func reportMostRecentlyTriggeredDeployment(ctx context.Context, client apiClient, d *model.Deployment) error {
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
				Versions:     d.Versions,
				StartedAt:    d.CreatedAt,
				CompletedAt:  d.CompletedAt,
			},
		}
		retry = pipedservice.NewRetry(10)
	)

	for retry.WaitNext(ctx) {
		if _, err = client.ReportApplicationMostRecentDeployment(ctx, req); err == nil {
			return nil
		}
		err = fmt.Errorf("failed to report most recent successful deployment: %w", err)
	}
	return err
}
