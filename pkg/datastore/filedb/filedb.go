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
	"fmt"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/filestore"
)

type FileDB struct {
	backend filestore.Store
	cache   cache.Cache
	logger  *zap.Logger
}

type Option func(*FileDB)

func WithLogger(logger *zap.Logger) Option {
	return func(f *FileDB) {
		f.logger = logger
	}
}

func NewFileDB(fs filestore.Store, c cache.Cache, opts ...Option) (*FileDB, error) {
	fd := &FileDB{
		backend: fs,
		cache:   c,
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
	scol, ok := col.(datastore.ShardStorable)
	if !ok {
		return nil, datastore.ErrUnsupported
	}

	var (
		kind   = col.Kind()
		shards = scol.ListInUsedShards()
		// Map of objects values with the first key is the object id.
		objects map[string][][]byte
	)

	// Fetch the first part of all objects under a specified "directory".
	for _, shard := range shards {
		dpath := makeHotStorageDirPath(kind, shard)
		parts, err := f.backend.List(ctx, dpath)
		if err != nil {
			f.logger.Error("failed to find entities",
				zap.String("kind", kind),
				zap.Error(err),
			)
			return nil, err
		}

		if objects == nil {
			objects = make(map[string][][]byte, len(parts))
		}
		for _, obj := range parts {
			id := filepath.Base(obj.Path)

			data, err := f.fetch(ctx, obj.Path)
			if err != nil {
				f.logger.Error("failed to fetch entity part",
					zap.String("kind", kind),
					zap.String("id", id),
					zap.Error(err),
				)
				return nil, err
			}

			objects[id] = append(objects[id], data)
		}
	}

	entities := make([]interface{}, 0, len(objects))
	for id, obj := range objects {
		e := col.Factory()()
		if err := decode(col, e, obj...); err != nil {
			f.logger.Error("failed to build entity",
				zap.String("kind", kind),
				zap.String("id", id),
				zap.Error(err),
			)
			return nil, err
		}

		pass, err := filter(col, e, opts.Filters)
		if err != nil {
			f.logger.Error("failed to filter entity",
				zap.String("kind", kind),
				zap.String("id", id),
				zap.Error(err),
			)
			return nil, err
		}

		if pass {
			entities = append(entities, e)
		}
	}

	return &Iterator{
		data: entities,
	}, nil
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

	return decode(col, v, parts...)
}

func (f *FileDB) Create(ctx context.Context, col datastore.Collection, id string, entity interface{}) error {
	kind := col.Kind()
	sdata, err := encode(col, entity)
	if err != nil {
		f.logger.Error("failed to encode entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	for shard, data := range sdata {
		path := makeHotStorageFilePath(kind, id, shard)
		if err = f.backend.Put(ctx, path, data); err != nil {
			f.logger.Error("failed to store entity",
				zap.String("id", id),
				zap.String("kind", kind),
				zap.Error(err),
			)
			return err
		}
	}

	return nil
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

func makeHotStorageDirPath(kind string, shard datastore.Shard) string {
	return fmt.Sprintf("%s/%s/", kind, shard)
}
