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
	return nil, nil
}

func (f *FileDB) Get(ctx context.Context, kind, id string, v interface{}) error {
	return nil
}

func (f *FileDB) Create(ctx context.Context, kind, id string, entity interface{}) error {
	return nil
}

func (f *FileDB) Put(ctx context.Context, kind, id string, entity interface{}) error {
	return nil
}

func (f *FileDB) Update(ctx context.Context, kind, id string, factory datastore.Factory, updater datastore.Updater) error {
	return nil
}

func (f *FileDB) Close() error {
	return f.backend.Close()
}
