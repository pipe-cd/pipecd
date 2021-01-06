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
	"errors"
	"time"

	"github.com/pipe-cd/pipe/pkg/model"
)

const imageMetadataModelKind = "ImageMetadata"

var (
	imageMetadataFactory = func() interface{} {
		return &model.ImageMetadata{}
	}
)

type ImageMetadataStore interface {
	PutImageMetadata(ctx context.Context, im model.ImageMetadata) error
	GetImageMetadata(ctx context.Context, id string) (*model.ImageMetadata, error)
}

type imageMetadataStore struct {
	backend
	nowFunc func() time.Time
}

func NewImageMetadataStore(ds DataStore) ImageMetadataStore {
	return &imageMetadataStore{
		backend: backend{
			ds: ds,
		},
		nowFunc: time.Now,
	}
}

func (s *imageMetadataStore) PutImageMetadata(ctx context.Context, im model.ImageMetadata) error {
	now := s.nowFunc().Unix()
	// Load the previously saved one to get the CreatedAt value.
	// In the future maybe we should change the interface of Put function to save this DB call.
	cur, err := s.GetImageMetadata(ctx, im.Id)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return err
		}
		im.CreatedAt = now
	} else {
		im.CreatedAt = cur.CreatedAt
	}

	if im.UpdatedAt == 0 {
		im.UpdatedAt = now
	}
	if err := im.Validate(); err != nil {
		return err
	}
	return s.ds.Put(ctx, imageMetadataModelKind, im.Id, &im)
}

func (s *imageMetadataStore) GetImageMetadata(ctx context.Context, id string) (*model.ImageMetadata, error) {
	var entity model.ImageMetadata
	if err := s.ds.Get(ctx, imageMetadataModelKind, id, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}
