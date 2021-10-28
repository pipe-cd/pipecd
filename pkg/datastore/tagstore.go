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

const TagModelKind = "Tag"

var tagFactory = func() interface{} {
	return &model.Tag{}
}

type TagStore interface {
	AddTag(ctx context.Context, app *model.Tag) error
	GetTag(ctx context.Context, id string) (*model.Tag, error)
}

type tagStore struct {
	backend
	nowFunc func() time.Time
}

func NewTagStore(ds DataStore) TagStore {
	return &tagStore{
		backend: backend{
			ds: ds,
		},
		nowFunc: time.Now,
	}
}

func (s *tagStore) AddTag(ctx context.Context, tag *model.Tag) error {
	now := s.nowFunc().Unix()
	if tag.CreatedAt == 0 {
		tag.CreatedAt = now
	}
	if tag.UpdatedAt == 0 {
		tag.UpdatedAt = now
	}
	if err := tag.Validate(); err != nil {
		return err
	}
	return s.ds.Create(ctx, TagModelKind, tag.Id, tag)
}

func (s *tagStore) GetTag(ctx context.Context, id string) (*model.Tag, error) {
	var entity model.Tag
	if err := s.ds.Get(ctx, TagModelKind, id, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}
