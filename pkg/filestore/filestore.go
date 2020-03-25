// Copyright 2020 The Pipe Authors.
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

type Object struct {
	Path    string
	Size    int64
	Content []byte
}

type Getter interface {
	// GetObject retrieves the content of a specific object in file storage bucket at path.
	GetObject(ctx context.Context, path string) (Object, error)
}

type Putter interface {
	// PutObject upload an object content to file storage at path.
	PutObject(ctx context.Context, path string, content []byte) error
}

type Lister interface {
	// ListObjects list all objects in file storage bucket at prefix.
	// The returned objects only contain the path to the object without object content.
	ListObjects(ctx context.Context, prefix string) ([]Object, error)
}

type Closer interface {
	Close() error
}

type Store interface {
	Getter
	Putter
	Lister
	Closer
	NewReader(ctx context.Context, path string) (io.ReadCloser, error)
	NewWriter(ctx context.Context, path string) io.WriteCloser
}
