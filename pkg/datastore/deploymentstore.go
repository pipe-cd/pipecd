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

package datastore

import (
	"context"
	"fmt"
	"time"

	"github.com/pipe-cd/pipecd/pkg/model"
)

type deploymentCollection struct {
}

func (d *deploymentCollection) Kind() string {
	return "Deployment"
}

func (d *deploymentCollection) Factory() Factory {
	return func() interface{} {
		return &model.Deployment{}
	}
}

var (
	DeploymentToPlannedUpdater = func(summary, statusReason, runningCommitHash, runningConfigFilename, version string, stages []*model.PipelineStage) func(*model.Deployment) error {
		return func(d *model.Deployment) error {
			d.Status = model.DeploymentStatus_DEPLOYMENT_PLANNED
			d.Summary = summary
			d.StatusReason = statusReason
			d.RunningCommitHash = runningCommitHash
			d.RunningConfigFilename = runningConfigFilename
			d.Version = version
			d.Stages = stages
			return nil
		}
	}

	DeploymentStatusUpdater = func(status model.DeploymentStatus, statusReason string) func(*model.Deployment) error {
		return func(d *model.Deployment) error {
			d.Status = status
			d.StatusReason = statusReason
			return nil
		}
	}

	DeploymentToCompletedUpdater = func(status model.DeploymentStatus, statuses map[string]model.StageStatus, statusReason string, completedAt int64) func(*model.Deployment) error {
		return func(d *model.Deployment) error {
			if !model.IsCompletedDeployment(status) {
				return fmt.Errorf("deployment status %s is not completed value: %w", status, ErrInvalidArgument)
			}

			d.Status = status
			d.StatusReason = statusReason
			d.CompletedAt = completedAt
			for i := range d.Stages {
				stageID := d.Stages[i].Id
				if status, ok := statuses[stageID]; ok {
					d.Stages[i].Status = status
				}
			}
			return nil
		}
	}

	StageStatusChangedUpdater = func(stageID string, status model.StageStatus, statusReason string, requires []string, visible bool, retriedCount int32, completedAt int64) func(*model.Deployment) error {
		return func(d *model.Deployment) error {
			for _, s := range d.Stages {
				if s.Id == stageID {
					s.Status = status
					s.StatusReason = statusReason
					if len(requires) > 0 {
						s.Requires = requires
					}
					s.Visible = visible
					s.RetriedCount = retriedCount
					s.CompletedAt = completedAt
					return nil
				}
			}
			return fmt.Errorf("stage id %s not found: %w", stageID, ErrInvalidArgument)
		}
	}
)

type DeploymentStore interface {
	AddDeployment(ctx context.Context, d *model.Deployment) error
	GetDeployment(ctx context.Context, id string) (*model.Deployment, error)
	ListDeployments(ctx context.Context, opts ListOptions) ([]*model.Deployment, string, error)
	UpdateDeployment(ctx context.Context, id string, updater func(*model.Deployment) error) error
	UpdateDeploymentMetadata(ctx context.Context, id string, metadata map[string]string) error
	UpdateDeploymentStageMetadata(ctx context.Context, deploymentID, stageID string, metadata map[string]string) error
}

type deploymentStore struct {
	backend
	nowFunc func() time.Time
}

func NewDeploymentStore(ds DataStore) DeploymentStore {
	return &deploymentStore{
		backend: backend{
			ds:  ds,
			col: &deploymentCollection{},
		},
		nowFunc: time.Now,
	}
}

func (s *deploymentStore) AddDeployment(ctx context.Context, d *model.Deployment) error {
	now := s.nowFunc().Unix()
	if d.CreatedAt == 0 {
		d.CreatedAt = now
	}
	if d.UpdatedAt == 0 {
		d.UpdatedAt = now
	}
	if err := d.Validate(); err != nil {
		return err
	}
	return s.ds.Create(ctx, s.col, d.Id, d)
}

func (s *deploymentStore) UpdateDeployment(ctx context.Context, id string, updater func(*model.Deployment) error) error {
	now := s.nowFunc().Unix()
	return s.ds.Update(ctx, s.col, id, func(e interface{}) error {
		d := e.(*model.Deployment)
		if err := updater(d); err != nil {
			return err
		}
		d.UpdatedAt = now
		return d.Validate()
	})
}

func (s *deploymentStore) UpdateDeploymentMetadata(ctx context.Context, id string, metadata map[string]string) error {
	now := s.nowFunc().Unix()
	return s.ds.Update(ctx, s.col, id, func(e interface{}) error {
		d := e.(*model.Deployment)
		d.Metadata = mergeMetadata(d.Metadata, metadata)
		d.UpdatedAt = now
		return nil
	})
}

func (s *deploymentStore) UpdateDeploymentStageMetadata(ctx context.Context, deploymentID, stageID string, metadata map[string]string) error {
	now := s.nowFunc().Unix()
	return s.ds.Update(ctx, s.col, deploymentID, func(e interface{}) error {
		d := e.(*model.Deployment)
		for _, stage := range d.Stages {
			if stage.Id == stageID {
				stage.Metadata = mergeMetadata(stage.Metadata, metadata)
				d.UpdatedAt = now
				return nil
			}
		}
		return fmt.Errorf("stage %s is not found: %w", stageID, ErrInvalidArgument)
	})
}

func mergeMetadata(ori map[string]string, new map[string]string) map[string]string {
	out := make(map[string]string, len(ori)+len(new))
	for k, v := range ori {
		out[k] = v
	}
	for k, v := range new {
		out[k] = v
	}
	return out
}

func (s *deploymentStore) ListDeployments(ctx context.Context, opts ListOptions) ([]*model.Deployment, string, error) {
	it, err := s.ds.Find(ctx, s.col, opts)
	if err != nil {
		return nil, "", err
	}
	ds := make([]*model.Deployment, 0)
	for {
		var d model.Deployment
		err := it.Next(&d)
		if err == ErrIteratorDone {
			break
		}
		if err != nil {
			return nil, "", err
		}
		ds = append(ds, &d)
	}

	// In case there is no more elements found, cursor should be set to empty too.
	if len(ds) == 0 {
		return ds, "", nil
	}
	cursor, err := it.Cursor()
	if err != nil {
		return nil, "", err
	}
	return ds, cursor, nil
}

func (s *deploymentStore) GetDeployment(ctx context.Context, id string) (*model.Deployment, error) {
	var entity model.Deployment
	if err := s.ds.Get(ctx, s.col, id, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}
