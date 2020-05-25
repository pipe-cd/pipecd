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
	"time"

	"github.com/pkg/errors"

	"github.com/kapetaniosci/pipe/pkg/model"
)

const deploymentModelKind = "Deployment"

var deploymentFactory = func() interface{} {
	return &model.Deployment{}
}

var (
	DeploymentToPlannedUpdater = func(desc, statusDesc string, stages []*model.PipelineStage) func(*model.Deployment) error {
		return func(d *model.Deployment) error {
			d.Status = model.DeploymentStatus_DEPLOYMENT_PLANNED
			d.Description = desc
			d.StatusDescription = statusDesc
			d.Stages = stages
			return nil
		}
	}

	DeploymentToRunningUpdater = func(statusDesc string) func(*model.Deployment) error {
		return func(d *model.Deployment) error {
			d.Status = model.DeploymentStatus_DEPLOYMENT_RUNNING
			d.StatusDescription = statusDesc
			return nil
		}
	}

	DeploymentToCompletedUpdater = func(status model.DeploymentStatus, statuses map[string]model.StageStatus, statusDesc string, completedAt int64) func(*model.Deployment) error {
		return func(d *model.Deployment) error {
			if !model.IsCompletedDeployment(status) {
				return errors.Wrapf(ErrInvalidArgument, "invalid completed deployment status: %s", status)
			}

			d.Status = status
			d.StatusDescription = statusDesc
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

	StageStatusChangedUpdater = func(stageID string, status model.StageStatus, statusDescription string,
		retriedCount int32, completedAt int64) func(*model.Deployment) error {

		return func(d *model.Deployment) error {
			d.StatusDescription = statusDescription
			for _, stage := range d.Stages {
				if stage.Id == stageID {
					stage.Status = status
					stage.RetriedCount = retriedCount
					stage.CompletedAt = completedAt
					return nil
				}
			}
			return errors.Wrapf(ErrInvalidArgument, "stage is not found: %s", stageID)
		}
	}
)

type DeploymentStore interface {
	AddDeployment(ctx context.Context, d *model.Deployment) error
	PutDeployment(ctx context.Context, id string, updater func(*model.Deployment) error) error
	PutDeploymentMetadata(ctx context.Context, id string, metadata map[string]string) error
	PutDeploymentStageMetadata(ctx context.Context, deploymentID, stageID, jsonMetadata string) error
	ListDeployments(ctx context.Context, opts ListOptions) ([]model.Deployment, error)
}

type deploymentStore struct {
	backend
	nowFunc func() time.Time
}

func NewDeploymentStore(ds DataStore) DeploymentStore {
	return &deploymentStore{
		backend: backend{
			ds: ds,
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
	return s.ds.Create(ctx, deploymentModelKind, d.Id, d)
}

func (s *deploymentStore) PutDeployment(ctx context.Context, id string, updater func(*model.Deployment) error) error {
	now := s.nowFunc().Unix()
	return s.ds.Update(ctx, deploymentModelKind, id, deploymentFactory, func(e interface{}) error {
		d := e.(*model.Deployment)
		if err := updater(d); err != nil {
			return err
		}
		d.UpdatedAt = now
		return d.Validate()
	})
}

func (s *deploymentStore) PutDeploymentMetadata(ctx context.Context, id string, metadata map[string]string) error {
	now := s.nowFunc().Unix()
	return s.ds.Update(ctx, deploymentModelKind, id, deploymentFactory, func(e interface{}) error {
		d := e.(*model.Deployment)
		d.Metadata = metadata
		d.UpdatedAt = now
		return nil
	})
}

func (s *deploymentStore) PutDeploymentStageMetadata(ctx context.Context, deploymentID, stageID, jsonMetadata string) error {
	now := s.nowFunc().Unix()
	return s.ds.Update(ctx, deploymentModelKind, deploymentID, deploymentFactory, func(e interface{}) error {
		d := e.(*model.Deployment)
		for _, stage := range d.Stages {
			if stage.Id == stageID {
				stage.JsonMetadata = jsonMetadata
				d.UpdatedAt = now
				return nil
			}
		}
		return errors.Wrapf(ErrInvalidArgument, "stage is not found: %s", stageID)
	})
}

func (s *deploymentStore) ListDeployments(ctx context.Context, opts ListOptions) ([]model.Deployment, error) {
	it, err := s.ds.Find(ctx, deploymentModelKind, opts)
	if err != nil {
		return nil, err
	}
	ds := make([]model.Deployment, 0)
	for {
		var d model.Deployment
		err := it.Next(&d)
		if err == ErrIteratorDone {
			break
		}
		if err != nil {
			return nil, err
		}
		ds = append(ds, d)
	}
	return ds, nil
}
