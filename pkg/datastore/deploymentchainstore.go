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
	"time"

	"github.com/pipe-cd/pipe/pkg/model"
)

const DeploymentChainModelKind = "DeploymentChain"

type DeploymentChainStore interface {
	AddDeploymentChain(ctx context.Context, d *model.DeploymentChain) error
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
