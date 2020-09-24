// Copyright 2020 The PipeCD Authors.
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

package gcs

import (
	"context"
	"testing"
	"time"

	"github.com/fsouza/fake-gcs-server/fakestorage"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/filestore/gcs"
)

func newEmulator(bucket string, objects map[string]string) (*fakestorage.Server, error) {
	initialObjects := make([]fakestorage.Object, 0, len(objects))
	for k, v := range objects {
		initialObjects = append(initialObjects, fakestorage.Object{
			BucketName: bucket,
			Name:       k,
			Content:    []byte(v),
		})
	}
	return fakestorage.NewServerWithOptions(fakestorage.Options{
		InitialObjects: initialObjects,
		Host:           "127.0.0.1",
		Port:           8081,
	})
}

func TestGetObject(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	bucket := "test"
	server, err := newEmulator(bucket, map[string]string{
		"path/to/file.txt": "foo",
	})
	assert.Nil(t, err)
	defer server.Stop()

	store, err := gcs.NewStore(ctx, bucket, gcs.WithHTTPClient(server.HTTPClient()))
	assert.Nil(t, err)

	tests := []struct {
		name    string
		path    string
		want    filestore.Object
		wantErr error
	}{
		{
			name: "found content",
			path: "path/to/file.txt",
			want: filestore.Object{
				Path:    "path/to/file.txt",
				Content: []byte("foo"),
				Size:    3,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			path: "path/to/wrong.txt",
			want: filestore.Object{
				Path: "path/to/wrong.txt",
			},
			wantErr: filestore.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.GetObject(ctx, tt.path)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPutObject(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	bucket := "test"
	server, err := newEmulator(bucket, map[string]string{
		"path/to/fileA.txt": "foo",
	})
	assert.Nil(t, err)
	defer server.Stop()

	store, err := gcs.NewStore(ctx, bucket, gcs.WithHTTPClient(server.HTTPClient()))
	assert.Nil(t, err)

	tests := []struct {
		name    string
		path    string
		content string
		want    filestore.Object
		wantErr bool
	}{
		{
			name:    "write new content",
			path:    "path/to/fileB.txt",
			content: "foo",
			want: filestore.Object{
				Path:    "path/to/fileB.txt",
				Content: []byte("foo"),
				Size:    3,
			},
			wantErr: false,
		},
		{
			name:    "overwrite content",
			path:    "path/to/fileA.txt",
			content: "bar",
			want: filestore.Object{
				Path:    "path/to/fileA.txt",
				Content: []byte("bar"),
				Size:    3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.PutObject(ctx, tt.path, []byte(tt.content))
			assert.Equal(t, tt.wantErr, err != nil)

			got, err := store.GetObject(ctx, tt.path)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestListObjects(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	bucket := "test"
	server, err := newEmulator(bucket, map[string]string{
		"path/to/fileA.txt": "foo",
		"path/to/fileB.txt": "bar",
	})
	assert.Nil(t, err)
	defer server.Stop()

	store, err := gcs.NewStore(ctx, bucket, gcs.WithHTTPClient(server.HTTPClient()))
	assert.Nil(t, err)

	tests := []struct {
		name    string
		prefix  string
		want    []filestore.Object
		wantErr bool
	}{
		{
			name:   "found contents",
			prefix: "path/to",
			want: []filestore.Object{
				{
					Path:    "path/to/fileA.txt",
					Content: []byte{},
					Size:    3,
				},
				{
					Path:    "path/to/fileB.txt",
					Content: []byte{},
					Size:    3,
				},
			},
			wantErr: false,
		},
		{
			name:    "not found",
			prefix:  "wrong",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.ListObjects(ctx, tt.prefix)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

// Check if it panics when connect with closed client.
func TestClose(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	bucket := "test"
	server, err := newEmulator(bucket, map[string]string{})
	assert.Nil(t, err)
	defer server.Stop()

	store, err := gcs.NewStore(ctx, bucket, gcs.WithHTTPClient(server.HTTPClient()))
	assert.Nil(t, err)

	err = store.Close()
	assert.Equal(t, false, err != nil)
	assert.Panics(t, func() {
		_, _ = store.GetObject(ctx, "path/to/fileA.txt")
	})
}
