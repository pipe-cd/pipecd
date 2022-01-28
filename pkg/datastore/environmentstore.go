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
	"errors"
	"time"

	"github.com/pipe-cd/pipecd/pkg/model"
)

type environmentCollection struct {
}

func (e *environmentCollection) Kind() string {
	return "Environment"
}

func (e *environmentCollection) Factory() Factory {
	return func() interface{} {
		return &model.Environment{}
	}
}

type EnvironmentStore interface {
	AddEnvironment(ctx context.Context, env *model.Environment) error
	GetEnvironment(ctx context.Context, id string) (*model.Environment, error)
	ListEnvironments(ctx context.Context, opts ListOptions) ([]*model.Environment, error)
	EnableEnvironment(ctx context.Context, id string) error
	DisableEnvironment(ctx context.Context, id string) error
	DeleteEnvironment(ctx context.Context, id string) error
}

type environmentStore struct {
	backend
	nowFunc func() time.Time
}

func NewEnvironmentStore(ds DataStore) EnvironmentStore {
	return &environmentStore{
		backend: backend{
			ds:  ds,
			col: &environmentCollection{},
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
	return s.ds.Create(ctx, s.col, env.Id, env)
}

func (s *environmentStore) EnableEnvironment(ctx context.Context, id string) error {
	return s.ds.Update(ctx, s.col, id, func(e interface{}) error {
		env := e.(*model.Environment)
		if env.Deleted {
			return errors.New("unable to enable a deleted environment")
		}
		env.Disabled = false
		env.UpdatedAt = s.nowFunc().Unix()
		return nil
	})
}

func (s *environmentStore) DisableEnvironment(ctx context.Context, id string) error {
	return s.ds.Update(ctx, s.col, id, func(e interface{}) error {
		env := e.(*model.Environment)
		if env.Deleted {
			return errors.New("unable to disable a deleted environment")
		}
		env.Disabled = true
		env.UpdatedAt = s.nowFunc().Unix()
		return nil
	})
}

func (s *environmentStore) DeleteEnvironment(ctx context.Context, id string) error {
	return s.ds.Update(ctx, s.col, id, func(e interface{}) error {
		now := s.nowFunc().Unix()
		env := e.(*model.Environment)
		env.Deleted = true
		env.Disabled = true
		env.DeletedAt = now
		env.UpdatedAt = now
		return nil
	})
}

func (s *environmentStore) GetEnvironment(ctx context.Context, id string) (*model.Environment, error) {
	var entity model.Environment
	if err := s.ds.Get(ctx, s.col, id, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (s *environmentStore) ListEnvironments(ctx context.Context, opts ListOptions) ([]*model.Environment, error) {
	it, err := s.ds.Find(ctx, s.col, opts)
	if err != nil {
		return nil, err
	}
	envs := make([]*model.Environment, 0)
	for {
		var env model.Environment
		err := it.Next(&env)
		if err == ErrIteratorDone {
			break
		}
		if err != nil {
			return nil, err
		}
		envs = append(envs, &env)
	}
	return envs, nil
}
