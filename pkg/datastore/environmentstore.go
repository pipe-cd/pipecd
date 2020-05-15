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

const environmentModelKind = "Environment"

var environmentFactory = func() interface{} {
	return &model.Environment{}
}

type EnvironmentStore interface {
	AddEnvironment(ctx context.Context, env *model.Environment) error
	ListEnvironments(ctx context.Context, opts ListOptions) ([]model.Environment, error)
}

type environmentStore struct {
	backend
	nowFunc func() time.Time
}

func NewEnvironmentStore(ds DataStore) EnvironmentStore {
	return &environmentStore{
		backend: backend{
			ds: ds,
		},
		nowFunc: time.Now,
	}
}

func (s *environmentStore) AddEnvironment(ctx context.Context, env *model.Environment) error {
	now := s.nowFunc().Unix()
	if env.CreatedAt == 0 {
		env.CreatedAt = now
	}
	if env.UpdatedAt == 0 {
		env.UpdatedAt = now
	}
	if err := env.Validate(); err != nil {
		return err
	}
	return s.ds.Create(ctx, environmentModelKind, env.Id, env)
}

func (s *environmentStore) ListEnvironments(ctx context.Context, opts ListOptions) ([]model.Environment, error) {
	it, err := s.ds.Find(ctx, environmentModelKind, opts)
	if err != nil {
		return nil, err
	}
	envs := make([]model.Environment, 0)
	for {
		var env model.Environment
		err := it.Next(&env)
		if err == ErrIteratorDone {
			break
		}
		if err != nil {
			return nil, err
		}
		envs = append(envs, env)
	}
	return envs, nil
}
