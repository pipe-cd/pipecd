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

type apiKeyCollection struct {
}

func (a *apiKeyCollection) Kind() string {
	return "APIKey"
}

func (a *apiKeyCollection) Factory() Factory {
	return func() interface{} {
		return &model.APIKey{}
	}
}

type APIKeyStore interface {
	Add(ctx context.Context, k *model.APIKey) error
	Get(ctx context.Context, id string) (*model.APIKey, error)
	List(ctx context.Context, opts ListOptions) ([]*model.APIKey, error)
	Disable(ctx context.Context, id, projectID string) error
}

type apiKeyStore struct {
	backend
	nowFunc func() time.Time
}

func NewAPIKeyStore(ds DataStore) APIKeyStore {
	return &apiKeyStore{
		backend: backend{
			ds:  ds,
			col: &apiKeyCollection{},
		},
		nowFunc: time.Now,
	}
}

func (s *apiKeyStore) Add(ctx context.Context, k *model.APIKey) error {
	now := s.nowFunc().Unix()
	if k.CreatedAt == 0 {
		k.CreatedAt = now
	}
	if k.UpdatedAt == 0 {
		k.UpdatedAt = now
	}
	if err := k.Validate(); err != nil {
		return err
	}
	return s.ds.Create(ctx, s.col, k.Id, k)
}

func (s *apiKeyStore) Get(ctx context.Context, id string) (*model.APIKey, error) {
	var entity model.APIKey
	if err := s.ds.Get(ctx, s.col, id, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (s *apiKeyStore) List(ctx context.Context, opts ListOptions) ([]*model.APIKey, error) {
	it, err := s.ds.Find(ctx, s.col, opts)
	if err != nil {
		return nil, err
	}
	ks := make([]*model.APIKey, 0)
	for {
		var k model.APIKey
		err := it.Next(&k)
		if err == ErrIteratorDone {
			break
		}
		if err != nil {
			return nil, err
		}
		ks = append(ks, &k)
	}
	return ks, nil
}

func (s *apiKeyStore) Disable(ctx context.Context, id, projectID string) error {
	now := s.nowFunc().Unix()
	return s.ds.Update(ctx, s.col, id, func(e interface{}) error {
		k := e.(*model.APIKey)
		if k.ProjectId != projectID {
			return fmt.Errorf("invalid project id, expected %s, got %s", k.ProjectId, projectID)
		}

		k.Disabled = true
		k.UpdatedAt = now
		return k.Validate()
	})
}
