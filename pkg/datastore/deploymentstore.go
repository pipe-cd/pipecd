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

	"github.com/kapetaniosci/pipe/pkg/model"
)

const deploymentModelKind = "Deployment"

var deploymentFactory = func() interface{} {
	return &model.Deployment{}
}

type DeploymentStore interface {
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
