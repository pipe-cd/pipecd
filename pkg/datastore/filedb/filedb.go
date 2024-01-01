// Copyright 2024 The PipeCD Authors.
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
	"path/filepath"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/datastore/filedb/objectcache"
	"github.com/pipe-cd/pipecd/pkg/filestore"
)

const (
	datastoreRootPath = "datastore"
)

type FileDB struct {
	backend     filestore.Store
	objectCache objectcache.Cache
	logger      *zap.Logger
}

type Option func(*FileDB)

func WithLogger(logger *zap.Logger) Option {
	return func(f *FileDB) {
		f.logger = logger
	}
}

func NewFileDB(fs filestore.Store, c cache.Cache, opts ...Option) (*FileDB, error) {
	fd := &FileDB{
		backend:     fs,
		objectCache: objectcache.NewCache(c),
		logger:      zap.NewNop(),
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
		objects map[string]map[datastore.Shard][]byte
	)

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
			objects = make(map[string]map[datastore.Shard][]byte, len(parts))
		}
		for _, obj := range parts {
			id := filepath.Base(obj.Path)

			if objects[id] == nil {
				objects[id] = make(map[datastore.Shard][]byte, len(scol.ListInUsedShards()))
			}

			// Try to get object content from objectCache.
			cdata, err := f.objectCache.Get(shard, id, obj.Etag)
			if err == nil {
				objects[id][shard] = cdata
				continue
			}

			// If there is no object content found from objectCache, try fetching
			// content under the object path.
			data, err := f.fetch(ctx, obj.Path)
			if err != nil {
				f.logger.Error("failed to fetch entity part",
					zap.String("kind", kind),
					zap.String("id", id),
					zap.Error(err),
				)
				return nil, err
			}

			// Store fetched data to cache.
			if err = f.objectCache.Put(shard, id, obj.Etag, data); err != nil {
				f.logger.Error("failed to store entity part to cache",
					zap.String("kind", kind),
					zap.String("id", id),
					zap.String("etag", obj.Etag),
					zap.Error(err),
				)
			}

			objects[id][shard] = data
		}
	}

	entities := make([]interface{}, 0, len(objects))
	for id, obj := range objects {
		e := col.Factory()()
		if err := decode(col, e, obj); err != nil {
			f.logger.Error("failed to build entity",
				zap.String("kind", kind),
				zap.String("id", id),
				zap.Error(err),
			)
			return nil, err
		}

		// TODO: Remove this unnecessary log print.
		f.logger.Info("filtering...",
			zap.Any("entity", e),
			zap.Any("filter", opts.Filters),
		)

		pass, err := filter(col, e, opts.Filters)
		if err != nil {
			f.logger.Error("failed to filter entity",
				zap.String("kind", kind),
				zap.String("id", id),
				zap.Error(err),
			)
			return nil, err
		}

		// TODO: Remove this unnecessary log print.
		f.logger.Info("check filter result",
			zap.Any("entity", e),
			zap.Any("filter", opts.Filters),
			zap.Bool("result", pass),
		)

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
	paths := make(map[datastore.Shard]string, len(shards))
	for _, s := range shards {
		paths[s] = makeHotStorageFilePath(kind, id, s)
	}

	parts := make(map[datastore.Shard][]byte, len(paths))
	for shard, path := range paths {
		part, err := f.fetch(ctx, path)
		if err != nil {
			f.logger.Error("failed to fetch entity",
				zap.String("id", id),
				zap.String("kind", kind),
				zap.Error(err),
			)
			return err
		}
		parts[shard] = part
	}

	return decode(col, v, parts)
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
	scol, ok := col.(datastore.ShardStorable)
	if !ok {
		return datastore.ErrUnsupported
	}

	kind := col.Kind()
	shard, err := scol.GetUpdatableShard()
	if err != nil {
		f.logger.Error("failed to prepare updatable shard",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	path := makeHotStorageFilePath(kind, id, shard)
	raw, err := f.fetch(ctx, path)
	if err != nil {
		f.logger.Error("failed to fetch updatable shard",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	entity := col.Factory()()
	if err = json.Unmarshal(raw, entity); err != nil {
		f.logger.Error("failed to decode entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	if err = updater(entity); err != nil {
		f.logger.Error("failed to apply updater",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	val, err := json.Marshal(entity)
	if err != nil {
		f.logger.Error("failed to encode entity",
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

func makeHotStorageFilePath(kind, id string, shard datastore.Shard) string {
	// TODO: Find a way to separate files by project to avoid fetch resources cross project.
	return fmt.Sprintf("%s/%s/%s/%s.json", datastoreRootPath, kind, shard, id)
}

func makeHotStorageDirPath(kind string, shard datastore.Shard) string {
	return fmt.Sprintf("%s/%s/%s/", datastoreRootPath, kind, shard)
}
