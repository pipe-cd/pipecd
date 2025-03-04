// Copyright 2025 The PipeCD Authors.
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

type deploymentTraceCollection struct {
	requestedBy Commander
}

func (d *deploymentTraceCollection) Kind() string {
	return "DeploymentTrace"
}

func (d *deploymentTraceCollection) Factory() Factory {
	return func() interface{} {
		return &model.DeploymentTrace{}
	}
}

func (d *deploymentTraceCollection) ListInUsedShards() []Shard {
	return []Shard{
		AgentShard,
	}
}

func (d *deploymentTraceCollection) GetUpdatableShard() (Shard, error) {
	switch d.requestedBy {
	case PipedCommander:
		return AgentShard, nil
	default:
		return "", ErrUnsupported
	}
}

func (d *deploymentTraceCollection) Encode(entity interface{}) (map[Shard][]byte, error) {
	const errFmt = "failed while encode Deployment Trace object: %s"

	me, ok := entity.(*model.DeploymentTrace)
	if !ok {
		return nil, fmt.Errorf(errFmt, "type not matched")
	}

	data, err := json.Marshal(me)
	if err != nil {
		return nil, fmt.Errorf(errFmt, "unable to marshal entity data")
	}
	return map[Shard][]byte{
		AgentShard: data,
	}, nil
}

type DeploymentTraceStore interface {
	Add(ctx context.Context, d model.DeploymentTrace) error
}

type deploymentTraceStore struct {
	backend
	commander Commander
	nowFunc   func() time.Time
}

func NewDeploymentTraceStore(ds DataStore, c Commander) DeploymentTraceStore {
	return &deploymentTraceStore{
		backend: backend{
			ds:  ds,
			col: &deploymentTraceCollection{requestedBy: c},
		},
		commander: c,
		nowFunc:   time.Now,
	}
}

func (s *deploymentTraceStore) Add(ctx context.Context, d model.DeploymentTrace) error {
	now := s.nowFunc().Unix()
	if d.CreatedAt == 0 {
		d.CreatedAt = now
	}
	if d.UpdatedAt == 0 {
		d.UpdatedAt = now
	}
	if err := d.Validate(); err != nil {
		return fmt.Errorf("failed to validate deployment trace: %w: %w", ErrInvalidArgument, err)
	}
	return s.ds.Create(ctx, s.col, d.Id, &d)
}
