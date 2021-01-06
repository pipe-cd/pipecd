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

const imageMetadataModelKind = "ImageMetadata"

var (
	imageMetadataFactory = func() interface{} {
		return &model.ImageMetadata{}
	}
)

type ImageMetadataStore interface {
	AddImageMetadata(ctx context.Context, im model.ImageMetadata) error
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

func (s *imageMetadataStore) AddImageMetadata(ctx context.Context, im model.ImageMetadata) error {
	now := s.nowFunc().Unix()
	if im.CreatedAt == 0 {
		im.CreatedAt = now
	}
	if im.UpdatedAt == 0 {
		im.UpdatedAt = now
	}
	if err := im.Validate(); err != nil {
		return err
	}
	return s.ds.Create(ctx, imageMetadataModelKind, im.Id, &im)
}
