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
	"encoding/json"
	"fmt"
	"time"

	"github.com/pipe-cd/pipecd/pkg/model"
)

type pipedCollection struct {
	requestedBy Commander
}

func (p *pipedCollection) Kind() string {
	return "Piped"
}

func (p *pipedCollection) Factory() Factory {
	return func() interface{} {
		return &model.Piped{}
	}
}

func (p *pipedCollection) ListInUsedShards() []Shard {
	return []Shard{
		ClientShard,
	}
}

func (p *pipedCollection) GetUpdatableShard() (Shard, error) {
	switch p.requestedBy {
	case WebCommander:
		return ClientShard, nil
	default:
		return "", ErrUnsupported
	}
}

func (p *pipedCollection) Encode(e interface{}) (map[Shard][]byte, error) {
	const errFmt = "failed while encode Piped object: %s"

	me, ok := e.(*model.Piped)
	if !ok {
		return nil, fmt.Errorf(errFmt, "type not matched")
	}

	data, err := json.Marshal(me)
	if err != nil {
		return nil, fmt.Errorf(errFmt, "unable to marshal entity data")
	}
	return map[Shard][]byte{
		ClientShard: data,
	}, nil
}

type PipedStore interface {
	Add(ctx context.Context, piped *model.Piped) error
	Get(ctx context.Context, id string) (*model.Piped, error)
	List(ctx context.Context, opts ListOptions) ([]*model.Piped, error)
	UpdateInfo(ctx context.Context, id, name, desc string, envIds []string) error
	EnablePiped(ctx context.Context, id string) error
	DisablePiped(ctx context.Context, id string) error
	UpdateDesiredVersion(ctx context.Context, id, version string) error
	UpdateMetadata(ctx context.Context, id, version string, cps []*model.Piped_CloudProvider, repos []*model.ApplicationGitRepository, se *model.Piped_SecretEncryption, startedAt int64) error
	AddKey(ctx context.Context, id, keyHash, creator string, createdAt time.Time) error
	DeleteOldKeys(ctx context.Context, id string) error
}

type pipedStore struct {
	backend
	commander Commander
	nowFunc   func() time.Time
}

func NewPipedStore(ds DataStore, c Commander) PipedStore {
	return &pipedStore{
		backend: backend{
			ds:  ds,
			col: &pipedCollection{requestedBy: c},
		},
		commander: c,
		nowFunc:   time.Now,
	}
}

func (s *pipedStore) Add(ctx context.Context, piped *model.Piped) error {
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

func (s *pipedStore) Get(ctx context.Context, id string) (*model.Piped, error) {
	var entity model.Piped
	if err := s.ds.Get(ctx, s.col, id, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (s *pipedStore) List(ctx context.Context, opts ListOptions) ([]*model.Piped, error) {
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

func (s *pipedStore) update(ctx context.Context, id string, updater func(piped *model.Piped) error) error {
	now := s.nowFunc().Unix()
	return s.ds.Update(ctx, s.col, id, func(e interface{}) error {
		p := e.(*model.Piped)
		if err := updater(p); err != nil {
			return err
		}
		p.UpdatedAt = now
		return p.Validate()
	})
}

func (s *pipedStore) UpdateInfo(ctx context.Context, id, name, desc string, envIds []string) error {
	return s.update(ctx, id, func(piped *model.Piped) error {
		piped.Name = name
		piped.Desc = desc
		piped.EnvIds = envIds
		return nil
	})
}

func (s *pipedStore) EnablePiped(ctx context.Context, id string) error {
	return s.update(ctx, id, func(piped *model.Piped) error {
		piped.Disabled = false
		piped.UpdatedAt = time.Now().Unix()
		return nil
	})
}

func (s *pipedStore) DisablePiped(ctx context.Context, id string) error {
	return s.update(ctx, id, func(piped *model.Piped) error {
		piped.Disabled = true
		piped.UpdatedAt = time.Now().Unix()
		return nil
	})
}

func (s *pipedStore) UpdateDesiredVersion(ctx context.Context, id, version string) error {
	return s.update(ctx, id, func(piped *model.Piped) error {
		piped.DesiredVersion = version
		return nil
	})
}

func (s *pipedStore) UpdateMetadata(ctx context.Context, id, version string, cps []*model.Piped_CloudProvider, repos []*model.ApplicationGitRepository, se *model.Piped_SecretEncryption, startedAt int64) error {
	return s.update(ctx, id, func(piped *model.Piped) error {
		piped.CloudProviders = cps
		piped.Repositories = repos

		piped.SecretEncryption = se
		// Remove the legacy data.
		piped.SealedSecretEncryption = nil

		piped.Version = version
		piped.StartedAt = startedAt
		return nil
	})
}

func (s *pipedStore) AddKey(ctx context.Context, id, keyHash, creator string, createdAt time.Time) error {
	return s.update(ctx, id, func(piped *model.Piped) error {
		piped.UpdatedAt = time.Now().Unix()
		return piped.AddKey(keyHash, creator, createdAt)
	})
}

func (s *pipedStore) DeleteOldKeys(ctx context.Context, id string) error {
	return s.update(ctx, id, func(piped *model.Piped) error {
		piped.DeleteOldPipedKeys()
		piped.UpdatedAt = time.Now().Unix()
		return nil
	})
}
