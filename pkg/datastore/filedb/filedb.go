// Copyright 2022 The PipeCD Authors.
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

package filedb

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/filestore"
)

type FileDB struct {
	backend filestore.Store
	logger  *zap.Logger
}

type Option func(*FileDB)

func WithLogger(logger *zap.Logger) Option {
	return func(m *FileDB) {
		m.logger = logger
	}
}

func NewFileDB(fs filestore.Store, opts ...Option) (*FileDB, error) {
	fd := &FileDB{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(fd)
	}

	fd.backend = fs

	return fd, nil
}

func (f *FileDB) Find(ctx context.Context, kind string, opts datastore.ListOptions) (datastore.Iterator, error) {
	return nil, datastore.ErrUnimplemented
}

func (f *FileDB) Get(ctx context.Context, kind, id string, v interface{}) error {
	// TODO: Find a better way to avoid knowledge leak about Application kind.
	var path string
	switch kind {
	case datastore.ApplicationModelKind:
		return datastore.ErrUnimplemented
	default:
		path = buildHotPath(kind, id)
	}

	raw, err := f.backend.Get(ctx, path)
	if err == filestore.ErrNotFound {
		return datastore.ErrNotFound
	}
	if err != nil {
		f.logger.Error("failed to get entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	if err = json.Unmarshal(raw, v); err != nil {
		f.logger.Error("failed to get entity: failed to decode entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (f *FileDB) Create(ctx context.Context, kind, id string, entity interface{}) error {
	// TODO: Find a better way to avoid knowledge leak about Application kind.
	var path string
	switch kind {
	case datastore.ApplicationModelKind:
		return datastore.ErrUnimplemented
	default:
		path = buildHotPath(kind, id)
	}

	// Note: To enable the current check existence logic works, the filestore
	// bucket need to be created with `Retention policy` configuration.
	// ref:
	//  - gcs: https://cloud.google.com/storage/docs/bucket-lock
	//  - s3: https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-lock.html
	//  - minio: https://docs.min.io/docs/minio-bucket-object-lock-guide.html
	// Depends on the network, retention time set to 1s (one second) is enough
	// to avoid create objects with duplicated key.
	_, err := f.backend.Get(ctx, path)
	if err == nil {
		return datastore.ErrAlreadyExists
	}
	if err != filestore.ErrNotFound {
		f.logger.Error("failed to create entity: failed to check entity existence",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	val, err := json.Marshal(entity)
	if err != nil {
		f.logger.Error("failed to create entity: failed to encode entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	if err = f.backend.Put(ctx, path, val); err != nil {
		f.logger.Error("failed to create entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (f *FileDB) Put(ctx context.Context, kind, id string, entity interface{}) error {
	// TODO: Find a better way to avoid knowledge leak about Application kind.
	var path string
	switch kind {
	case datastore.ApplicationModelKind:
		return datastore.ErrUnimplemented
	default:
		path = buildHotPath(kind, id)
	}

	val, err := json.Marshal(entity)
	if err != nil {
		f.logger.Error("failed to put entity: failed to encode entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	if err = f.backend.Put(ctx, path, val); err != nil {
		f.logger.Error("failed to put entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (f *FileDB) Update(ctx context.Context, kind, id string, factory datastore.Factory, updater datastore.Updater) error {
	// Note: PipeCD follows `single writer pattern`, which means
	// there will be no two or more processes which try to update
	// a specified object at once, so we don't need to open transaction here.

	// TODO: Find a better way to avoid knowledge leak about Application kind.
	var path string
	switch kind {
	case datastore.ApplicationModelKind:
		return datastore.ErrUnimplemented
	default:
		path = buildHotPath(kind, id)
	}

	raw, err := f.backend.Get(ctx, path)
	if err == filestore.ErrNotFound {
		return datastore.ErrNotFound
	}
	if err != nil {
		f.logger.Error("failed to update entity: failed to get entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	entity := factory()
	if err = json.Unmarshal(raw, entity); err != nil {
		f.logger.Error("failed to update entity: failed to decode entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	if err = updater(entity); err != nil {
		f.logger.Error("failed to update entity: failed to apply updater",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	val, err := json.Marshal(entity)
	if err != nil {
		f.logger.Error("failed to update entity: failed to encode entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	if err = f.backend.Put(ctx, path, val); err != nil {
		f.logger.Error("failed to update entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (f *FileDB) Close() error {
	return f.backend.Close()
}
