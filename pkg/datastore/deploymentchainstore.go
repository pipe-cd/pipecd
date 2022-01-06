// Copyright 2021 The PipeCD Authors.
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

	"github.com/pipe-cd/pipe/pkg/model"
)

const DeploymentChainModelKind = "DeploymentChain"

var deploymentChainFactory = func() interface{} {
	return &model.DeploymentChain{}
}

var (
	DeploymentChainAddDeploymentToBlock = func(deployment *model.Deployment) func(*model.DeploymentChain) error {
		return func(dc *model.DeploymentChain) error {
			if deployment.DeploymentChainBlockIndex >= uint32(len(dc.Blocks)) {
				return fmt.Errorf("invalid block index (%d) provided", deployment.DeploymentChainBlockIndex)
			}

			block := dc.Blocks[deployment.DeploymentChainBlockIndex]
			var updated bool
			for _, node := range block.Nodes {
				if node.ApplicationRef.ApplicationId != deployment.ApplicationId {
					continue
				}
				node.DeploymentRef = &model.ChainDeploymentRef{
					DeploymentId: deployment.Id,
					Status:       deployment.Status,
					StatusReason: deployment.StatusReason,
				}
				updated = true
				break
			}
			if !updated {
				return fmt.Errorf("unable to find the right node in chain to assign deployment to")
			}
			return nil
		}
	}

	DeploymentChainNodeDeploymentStatusUpdater = func(blockIndex uint32, deploymentID string, status model.DeploymentStatus, reason string) func(*model.DeploymentChain) error {
		return func(dc *model.DeploymentChain) error {
			if blockIndex >= uint32(len(dc.Blocks)) {
				return fmt.Errorf("invalid block index %d provided", blockIndex)
			}

			block := dc.Blocks[blockIndex]
			node, err := block.GetNodeByDeploymentID(deploymentID)
			if err != nil {
				return err
			}
			node.DeploymentRef.Status = status
			node.DeploymentRef.StatusReason = reason

			// If the block is already finished, keep its finished status.
			if block.IsCompleted() {
				return nil
			}

			// Update block status based on its state after update latest submitted deployment status.
			block.Status = block.DesiredStatus()
			if block.IsCompleted() {
				block.CompletedAt = time.Now().Unix()
			}

			// Update chain status based on its updated blocks state.
			dc.Status = dc.DesiredStatus()
			if dc.IsCompleted() {
				dc.CompletedAt = time.Now().Unix()
			}

			return nil
		}
	}
)

type DeploymentChainStore interface {
	AddDeploymentChain(ctx context.Context, d *model.DeploymentChain) error
	UpdateDeploymentChain(ctx context.Context, id string, updater func(*model.DeploymentChain) error) error
	GetDeploymentChain(ctx context.Context, id string) (*model.DeploymentChain, error)
	ListDeploymentChains(ctx context.Context, opts ListOptions) ([]*model.DeploymentChain, string, error)
}

type deploymentChainStore struct {
	backend
	nowFunc func() time.Time
}

func NewDeploymentChainStore(ds DataStore) DeploymentChainStore {
	return &deploymentChainStore{
		backend: backend{
			ds: ds,
		},
		nowFunc: time.Now,
	}
}

func (s *deploymentChainStore) AddDeploymentChain(ctx context.Context, dc *model.DeploymentChain) error {
	now := s.nowFunc().Unix()
	if dc.CreatedAt == 0 {
		dc.CreatedAt = now
	}
	if dc.UpdatedAt == 0 {
		dc.UpdatedAt = now
	}
	if err := dc.Validate(); err != nil {
		return err
	}
	return s.ds.Create(ctx, DeploymentChainModelKind, dc.Id, dc)
}

func (s *deploymentChainStore) UpdateDeploymentChain(ctx context.Context, id string, updater func(*model.DeploymentChain) error) error {
	now := s.nowFunc().Unix()
	return s.ds.Update(ctx, DeploymentChainModelKind, id, deploymentChainFactory, func(e interface{}) error {
		dc := e.(*model.DeploymentChain)
		if err := updater(dc); err != nil {
			return err
		}
		dc.UpdatedAt = now
		return dc.Validate()
	})
}

func (s *deploymentChainStore) GetDeploymentChain(ctx context.Context, id string) (*model.DeploymentChain, error) {
	var entity model.DeploymentChain
	if err := s.ds.Get(ctx, DeploymentChainModelKind, id, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (s *deploymentChainStore) ListDeploymentChains(ctx context.Context, opts ListOptions) ([]*model.DeploymentChain, string, error) {
	it, err := s.ds.Find(ctx, DeploymentChainModelKind, opts)
	if err != nil {
		return nil, "", err
	}
	dcs := make([]*model.DeploymentChain, 0)
	for {
		var dc model.DeploymentChain
		err := it.Next(&dc)
		if err == ErrIteratorDone {
			break
		}
		if err != nil {
			return nil, "", err
		}
		dcs = append(dcs, &dc)
	}

	// In case there is no more elements found, cursor should be set to empty too.
	if len(dcs) == 0 {
		return dcs, "", nil
	}
	cursor, err := it.Cursor()
	if err != nil {
		return nil, "", err
	}
	return dcs, cursor, nil
}
