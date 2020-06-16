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

	"github.com/pipe-cd/pipe/pkg/model"
)

const pipedModelKind = "Piped"

var pipedFactory = func() interface{} {
	return &model.Piped{}
}

type PipedStore interface {
	AddPiped(ctx context.Context, piped *model.Piped) error
	GetPiped(ctx context.Context, id string) (*model.Piped, error)
	ListPipeds(ctx context.Context, opts ListOptions) ([]model.Piped, error)
}

type pipedStore struct {
	backend
	nowFunc func() time.Time
}

func NewPipedStore(ds DataStore) PipedStore {
	return &pipedStore{
		backend: backend{
			ds: ds,
		},
		nowFunc: time.Now,
	}
}

func (s *pipedStore) AddPiped(ctx context.Context, piped *model.Piped) error {
	now := s.nowFunc().Unix()
	if piped.CreatedAt == 0 {
		piped.CreatedAt = now
	}
	if piped.UpdatedAt == 0 {
		piped.UpdatedAt = now
	}
	if err := piped.Validate(); err != nil {
		return err
	}
	return s.ds.Create(ctx, pipedModelKind, piped.Id, piped)
}

func (s *pipedStore) GetPiped(ctx context.Context, id string) (*model.Piped, error) {
	var entity model.Piped
	if err := s.ds.Get(ctx, pipedModelKind, id, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (s *pipedStore) ListPipeds(ctx context.Context, opts ListOptions) ([]model.Piped, error) {
	it, err := s.ds.Find(ctx, pipedModelKind, opts)
	if err != nil {
		return nil, err
	}
	ps := make([]model.Piped, 0)
	for {
		var p model.Piped
		err := it.Next(&p)
		if err == ErrIteratorDone {
			break
		}
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}
