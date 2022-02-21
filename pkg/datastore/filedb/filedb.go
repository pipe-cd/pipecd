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
	"fmt"

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
	return func(f *FileDB) {
		f.logger = logger
	}
}

func NewFileDB(fs filestore.Store, opts ...Option) (*FileDB, error) {
	fd := &FileDB{
		backend: fs,
		logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(fd)
	}

	return fd, nil
}

func (f *FileDB) fetch(ctx context.Context, path string) ([]byte, error) {
	raw, err := f.backend.Get(ctx, path)
	if err == filestore.ErrNotFound {
		return nil, datastore.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (f *FileDB) Find(ctx context.Context, col datastore.Collection, opts datastore.ListOptions) (datastore.Iterator, error) {
	_, ok := col.(datastore.ShardStorable)
	if !ok {
		return nil, datastore.ErrUnsupported
	}
	return nil, datastore.ErrUnimplemented
}

func (f *FileDB) Get(ctx context.Context, col datastore.Collection, id string, v interface{}) error {
	fcol, ok := col.(datastore.ShardStorable)
	if !ok {
		return datastore.ErrUnsupported
	}

	kind := col.Kind()
	shards := fcol.ListInUsedShards()
	paths := make([]string, 0, len(shards))
	for _, s := range shards {
		paths = append(paths, makeHotStorageFilePath(kind, id, s))
	}

	parts := make([][]byte, 0, len(paths))
	for _, path := range paths {
		part, err := f.fetch(ctx, path)
		if err != nil {
			f.logger.Error("failed to fetch entity",
				zap.String("id", id),
				zap.String("kind", kind),
				zap.Error(err),
			)
			return err
		}
		parts = append(parts, part)
	}

	if len(parts) == 1 {
		return json.Unmarshal(parts[0], v)
	}

	// TODO: Add merge based on UpdatedAt field in case there are multiple parts of object are fetched.

	return datastore.ErrUnsupported
}

func (f *FileDB) Create(ctx context.Context, col datastore.Collection, id string, entity interface{}) error {
	_, ok := col.(datastore.ShardStorable)
	if !ok {
		return datastore.ErrUnsupported
	}
	return datastore.ErrUnimplemented
}

func (f *FileDB) Update(ctx context.Context, col datastore.Collection, id string, updater datastore.Updater) error {
	_, ok := col.(datastore.ShardStorable)
	if !ok {
		return datastore.ErrUnsupported
	}
	return datastore.ErrUnimplemented
}

func (f *FileDB) Close() error {
	return f.backend.Close()
}

func makeHotStorageFilePath(kind, id string, shard datastore.Shard) string {
	// TODO: Find a way to separate files by project to avoid fetch resources cross project.
	return fmt.Sprintf("%s/%s/%s.json", kind, shard, id)
}
