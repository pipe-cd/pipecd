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

package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pipe-cd/pipecd/pkg/model"
)

type deploymentChainCollection struct {
	requestedBy Commander
}

func (d *deploymentChainCollection) Kind() string {
	return "DeploymentChain"
}

func (d *deploymentChainCollection) Factory() Factory {
	return func() interface{} {
		return &model.DeploymentChain{}
	}
}

func (d *deploymentChainCollection) ListInUsedShards() []Shard {
	return []Shard{
		OpsShard,
	}
}

func (d *deploymentChainCollection) GetUpdatableShard() (Shard, error) {
	switch d.requestedBy {
	case OpsCommander:
		return OpsShard, nil
	default:
		return "", ErrUnsupported
	}
}

func (d *deploymentChainCollection) Encode(e interface{}) (map[Shard][]byte, error) {
	const errFmt = "failed while encode DeploymentChain object: %s"

	me, ok := e.(*model.DeploymentChain)
	if !ok {
		return nil, fmt.Errorf(errFmt, "type not matched")
	}

	data, err := json.Marshal(me)
	if err != nil {
		return nil, fmt.Errorf(errFmt, "unable to marshal entity data")
	}
	return map[Shard][]byte{
		OpsShard: data,
	}, nil
}

var (
	addDeploymentToBlockUpdateFunc = func(deployment *model.Deployment) func(*model.DeploymentChain) error {
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

	nodeDeploymentStatusUpdateFunc = func(blockIndex uint32, deploymentID string, status model.DeploymentStatus, reason string) func(*model.DeploymentChain) error {
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
	Add(ctx context.Context, d *model.DeploymentChain) error
	Get(ctx context.Context, id string) (*model.DeploymentChain, error)
	List(ctx context.Context, opts ListOptions) ([]*model.DeploymentChain, string, error)
	AddNodeDeployment(ctx context.Context, chainID string, deployment *model.Deployment) error
	UpdateNodeDeploymentStatus(ctx context.Context, chainID string, blockIndex uint32, deploymentID string, status model.DeploymentStatus, reason string) error
}

type deploymentChainStore struct {
	backend
	commander Commander
	nowFunc   func() time.Time
}

func NewDeploymentChainStore(ds DataStore, c Commander) DeploymentChainStore {
	return &deploymentChainStore{
		backend: backend{
			ds:  ds,
			col: &deploymentChainCollection{requestedBy: c},
		},
		commander: c,
		nowFunc:   time.Now,
	}
}

func (s *deploymentChainStore) Add(ctx context.Context, dc *model.DeploymentChain) error {
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
	return s.ds.Create(ctx, s.col, dc.Id, dc)
}

func (s *deploymentChainStore) Get(ctx context.Context, id string) (*model.DeploymentChain, error) {
	var entity model.DeploymentChain
	if err := s.ds.Get(ctx, s.col, id, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (s *deploymentChainStore) List(ctx context.Context, opts ListOptions) ([]*model.DeploymentChain, string, error) {
	it, err := s.ds.Find(ctx, s.col, opts)
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

func (s *deploymentChainStore) update(ctx context.Context, id string, updater func(*model.DeploymentChain) error) error {
	now := s.nowFunc().Unix()
	return s.ds.Update(ctx, s.col, id, func(e interface{}) error {
		dc := e.(*model.DeploymentChain)
		if err := updater(dc); err != nil {
			return err
		}
		dc.UpdatedAt = now
		return dc.Validate()
	})
}

func (s *deploymentChainStore) AddNodeDeployment(ctx context.Context, chainID string, deployment *model.Deployment) error {
	updater := addDeploymentToBlockUpdateFunc(deployment)
	return s.update(ctx, chainID, updater)
}

func (s *deploymentChainStore) UpdateNodeDeploymentStatus(ctx context.Context, chainID string, blockIndex uint32, deploymentID string, status model.DeploymentStatus, reason string) error {
	updater := nodeDeploymentStatusUpdateFunc(blockIndex, deploymentID, status, reason)
	return s.update(ctx, chainID, updater)
}
