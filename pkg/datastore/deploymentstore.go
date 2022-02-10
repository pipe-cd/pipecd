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
	toPlannedUpdateFunc = func(summary, statusReason, runningCommitHash, runningConfigFilename, version string, stages []*model.PipelineStage) func(*model.Deployment) error {
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

	toCompletedUpdateFunc = func(status model.DeploymentStatus, stageStatuses map[string]model.StageStatus, statusReason string, completedAt int64) func(*model.Deployment) error {
		return func(d *model.Deployment) error {
			if !model.IsCompletedDeployment(status) {
				return fmt.Errorf("deployment status %s is not completed value: %w", status, ErrInvalidArgument)
			}

			d.Status = status
			d.StatusReason = statusReason
			d.CompletedAt = completedAt
			for i := range d.Stages {
				stageID := d.Stages[i].Id
				if status, ok := stageStatuses[stageID]; ok {
					d.Stages[i].Status = status
				}
			}
			return nil
		}
	}

	statusUpdateFunc = func(status model.DeploymentStatus, statusReason string) func(*model.Deployment) error {
		return func(d *model.Deployment) error {
			d.Status = status
			d.StatusReason = statusReason
			return nil
		}
	}

	stageStatusUpdateFunc = func(stageID string, status model.StageStatus, statusReason string, requires []string, visible bool, retriedCount int32, completedAt int64) func(*model.Deployment) error {
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
	Add(ctx context.Context, d *model.Deployment) error
	Get(ctx context.Context, id string) (*model.Deployment, error)
	List(ctx context.Context, opts ListOptions) ([]*model.Deployment, string, error)
	UpdateToPlanned(ctx context.Context, id, summary, reason, runningCommitHash, runningConfigFilename, version string, stages []*model.PipelineStage) error
	UpdateToCompleted(ctx context.Context, id string, status model.DeploymentStatus, stageStatuses map[string]model.StageStatus, reason string, completedAt int64) error
	UpdateStatus(ctx context.Context, id string, status model.DeploymentStatus, reason string) error
	UpdateStageStatus(ctx context.Context, id, stageID string, status model.StageStatus, reason string, requires []string, visible bool, retriedCount int32, completedAt int64) error
	UpdateMetadata(ctx context.Context, id string, metadata map[string]string) error
	UpdateStageMetadata(ctx context.Context, deploymentID, stageID string, metadata map[string]string) error
}

type deploymentStore struct {
	backend
	writer  Writer
	nowFunc func() time.Time
}

func NewDeploymentStore(ds DataStore, w Writer) DeploymentStore {
	return &deploymentStore{
		backend: backend{
			ds:  ds,
			col: &deploymentCollection{},
		},
		writer:  w,
		nowFunc: time.Now,
	}
}

func (s *deploymentStore) Add(ctx context.Context, d *model.Deployment) error {
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

func (s *deploymentStore) Get(ctx context.Context, id string) (*model.Deployment, error) {
	var entity model.Deployment
	if err := s.ds.Get(ctx, s.col, id, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (s *deploymentStore) List(ctx context.Context, opts ListOptions) ([]*model.Deployment, string, error) {
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

func (s *deploymentStore) update(ctx context.Context, id string, updater func(*model.Deployment) error) error {
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

func (s *deploymentStore) UpdateToPlanned(ctx context.Context, id, summary, reason, runningCommitHash, runningConfigFilename, version string, stages []*model.PipelineStage) error {
	updater := toPlannedUpdateFunc(summary, reason, runningCommitHash, runningConfigFilename, version, stages)
	return s.update(ctx, id, updater)
}

func (s *deploymentStore) UpdateToCompleted(ctx context.Context, id string, status model.DeploymentStatus, stageStatuses map[string]model.StageStatus, reason string, completedAt int64) error {
	updater := toCompletedUpdateFunc(status, stageStatuses, reason, completedAt)
	return s.update(ctx, id, updater)
}

func (s *deploymentStore) UpdateStatus(ctx context.Context, id string, status model.DeploymentStatus, reason string) error {
	updater := statusUpdateFunc(status, reason)
	return s.update(ctx, id, updater)
}

func (s *deploymentStore) UpdateStageStatus(ctx context.Context, id, stageID string, status model.StageStatus, reason string, requires []string, visible bool, retriedCount int32, completedAt int64) error {
	updater := stageStatusUpdateFunc(stageID, status, reason, requires, visible, retriedCount, completedAt)
	return s.update(ctx, id, updater)
}

func (s *deploymentStore) UpdateMetadata(ctx context.Context, id string, metadata map[string]string) error {
	return s.update(ctx, id, func(d *model.Deployment) error {
		d.Metadata = mergeMetadata(d.Metadata, metadata)
		return nil
	})
}

func (s *deploymentStore) UpdateStageMetadata(ctx context.Context, deploymentID, stageID string, metadata map[string]string) error {
	return s.update(ctx, deploymentID, func(d *model.Deployment) error {
		for _, stage := range d.Stages {
			if stage.Id == stageID {
				stage.Metadata = mergeMetadata(stage.Metadata, metadata)
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
