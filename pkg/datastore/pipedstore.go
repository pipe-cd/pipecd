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

	"github.com/pipe-cd/pipecd/pkg/model"
)

type pipedCollection struct {
}

func (p *pipedCollection) Kind() string {
	return "Piped"
}

func (p *pipedCollection) Factory() Factory {
	return func() interface{} {
		return &model.Piped{}
	}
}

var (
	PipedMetadataUpdater = func(
		cloudProviders []*model.Piped_CloudProvider,
		repos []*model.ApplicationGitRepository,
		se *model.Piped_SecretEncryption,
		version string,
		startedAt int64,
	) func(piped *model.Piped) error {

		return func(piped *model.Piped) error {
			piped.CloudProviders = cloudProviders
			piped.Repositories = repos

			piped.SecretEncryption = se
			// Remove the legacy data.
			piped.SealedSecretEncryption = nil

			piped.Version = version
			piped.StartedAt = startedAt
			return nil
		}
	}
)

type PipedStore interface {
	AddPiped(ctx context.Context, piped *model.Piped) error
	GetPiped(ctx context.Context, id string) (*model.Piped, error)
	ListPipeds(ctx context.Context, opts ListOptions) ([]*model.Piped, error)
	UpdatePiped(ctx context.Context, w Writer, id string, updater func(piped *model.Piped) error) error
	EnablePiped(ctx context.Context, w Writer, id string) error
	DisablePiped(ctx context.Context, w Writer, id string) error
	AddKey(ctx context.Context, w Writer, id, keyHash, creator string, createdAt time.Time) error
	DeleteOldKeys(ctx context.Context, w Writer, id string) error
}

type pipedStore struct {
	backend
	nowFunc func() time.Time
}

func NewPipedStore(ds DataStore) PipedStore {
	return &pipedStore{
		backend: backend{
			ds:  ds,
			col: &pipedCollection{},
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
	return s.ds.Create(ctx, s.col, piped.Id, piped)
}

func (s *pipedStore) GetPiped(ctx context.Context, id string) (*model.Piped, error) {
	var entity model.Piped
	if err := s.ds.Get(ctx, s.col, id, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (s *pipedStore) ListPipeds(ctx context.Context, opts ListOptions) ([]*model.Piped, error) {
	it, err := s.ds.Find(ctx, s.col, opts)
	if err != nil {
		return nil, err
	}
	ps := make([]*model.Piped, 0)
	for {
		var p model.Piped
		err := it.Next(&p)
		if err == ErrIteratorDone {
			break
		}
		if err != nil {
			return nil, err
		}
		ps = append(ps, &p)
	}
	return ps, nil
}

func (s *pipedStore) UpdatePiped(ctx context.Context, w Writer, id string, updater func(piped *model.Piped) error) error {
	now := s.nowFunc().Unix()
	return s.ds.Update(ctx, s.col, id, NewUpdater(w, func(e interface{}) error {
		p := e.(*model.Piped)
		if err := updater(p); err != nil {
			return err
		}
		p.UpdatedAt = now
		return p.Validate()
	}))
}

func (s *pipedStore) EnablePiped(ctx context.Context, w Writer, id string) error {
	return s.UpdatePiped(ctx, w, id, func(piped *model.Piped) error {
		piped.Disabled = false
		piped.UpdatedAt = time.Now().Unix()
		return nil
	})
}

func (s *pipedStore) DisablePiped(ctx context.Context, w Writer, id string) error {
	return s.UpdatePiped(ctx, w, id, func(piped *model.Piped) error {
		piped.Disabled = true
		piped.UpdatedAt = time.Now().Unix()
		return nil
	})
}

func (s *pipedStore) AddKey(ctx context.Context, w Writer, id, keyHash, creator string, createdAt time.Time) error {
	return s.UpdatePiped(ctx, w, id, func(piped *model.Piped) error {
		piped.UpdatedAt = time.Now().Unix()
		return piped.AddKey(keyHash, creator, createdAt)
	})
}

func (s *pipedStore) DeleteOldKeys(ctx context.Context, w Writer, id string) error {
	return s.UpdatePiped(ctx, w, id, func(piped *model.Piped) error {
		piped.DeleteOldPipedKeys()
		piped.UpdatedAt = time.Now().Unix()
		return nil
	})
}
