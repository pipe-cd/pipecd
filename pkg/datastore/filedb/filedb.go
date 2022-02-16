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

type FileStorable interface {
	// GetStoredFileNames returns a list of files' name which object of kind currently storing in.
	GetStoredFileNames(id string) []string
	// GetUpdatableFileName returns the name of the file which should be referred to on Updating object of kind.
	// datastore.ErrUnsupported will be returned if there is no such file name exist.
	GetUpdatableFileName(id string) (string, error)
}

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
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(fd)
	}

	fd.backend = fs

	return fd, nil
}

func (f *FileDB) Find(ctx context.Context, col datastore.Collection, opts datastore.ListOptions) (datastore.Iterator, error) {
	_, ok := col.(FileStorable)
	if !ok {
		return nil, datastore.ErrUnsupported
	}
	return nil, datastore.ErrUnimplemented
}

func (f *FileDB) Get(ctx context.Context, col datastore.Collection, id string, v interface{}) error {
	_, ok := col.(FileStorable)
	if !ok {
		return datastore.ErrUnsupported
	}
	return datastore.ErrUnimplemented
}

func (f *FileDB) Create(ctx context.Context, col datastore.Collection, id string, entity interface{}) error {
	_, ok := col.(FileStorable)
	if !ok {
		return datastore.ErrUnsupported
	}
	return datastore.ErrUnimplemented
}

func (f *FileDB) Update(ctx context.Context, col datastore.Collection, id string, updater datastore.Updater) error {
	_, ok := col.(FileStorable)
	if !ok {
		return datastore.ErrUnsupported
	}
	return datastore.ErrUnimplemented
}

func (f *FileDB) Close() error {
	return f.backend.Close()
}
