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

package filestore

import (
	"context"
	"errors"
	"io"
)

var (
	ErrNotFound = errors.New("not found")
)

type ObjectAttrs struct {
	Path      string
	Size      int64
	Etag      string
	UpdatedAt int64
}

type Getter interface {
	// Get returns bytes content of file object at the given path.
	Get(ctx context.Context, path string) ([]byte, error)

	// GetReader returns a new Reader to read the contents of the object.
	// The caller must call Close on the returned Reader when done reading.
	GetReader(ctx context.Context, path string) (io.ReadCloser, error)
}

type Putter interface {
	// Put uploads a file object to store at the given path.
	Put(ctx context.Context, path string, content []byte) error
}

type Lister interface {
	// List finds all objects in storage where its path starts with the given prefix
	// and returns objects' attributes (without content).
	List(ctx context.Context, prefix string) ([]ObjectAttrs, error)
}

type Deleter interface {
	// Delete deletes a file object at the given path.
	Delete(ctx context.Context, path string) error
}

type Closer interface {
	Close() error
}

type Store interface {
	Getter
	Putter
	Deleter
	Lister
	Closer
}
