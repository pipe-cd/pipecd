// Copyright 2023 The PipeCD Authors.
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

	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/filestore/gcs"
)

func newEmulator(bucket string, objects map[string]string, now time.Time) (*fakestorage.Server, error) {
	initialObjects := make([]fakestorage.Object, 0, len(objects))
	for k, v := range objects {
		initialObjects = append(initialObjects, fakestorage.Object{
			BucketName: bucket,
			Name:       k,
			Content:    []byte(v),
			Created:    now,
			Updated:    now,
		})
	}
	return fakestorage.NewServerWithOptions(fakestorage.Options{
		InitialObjects: initialObjects,
		Host:           "127.0.0.1",
		Port:           8081,
	})
}

func TestGet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	bucket := "test"
	server, err := newEmulator(
		bucket,
		map[string]string{
			"path/to/file.txt": "foo",
		},
		time.Now(),
	)
	assert.Nil(t, err)
	defer server.Stop()

	store, err := gcs.NewStore(ctx, bucket, gcs.WithHTTPClient(server.HTTPClient()))
	assert.Nil(t, err)

	tests := []struct {
		name    string
		path    string
		want    []byte
		wantErr error
	}{
		{
			name:    "found content",
			path:    "path/to/file.txt",
			want:    []byte("foo"),
			wantErr: nil,
		},
		{
			name:    "not found",
			path:    "path/to/wrong.txt",
			want:    nil,
			wantErr: filestore.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.Get(ctx, tt.path)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPut(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	bucket := "test"
	server, err := newEmulator(
		bucket,
		map[string]string{
			"path/to/fileA.txt": "foo",
		},
		time.Now(),
	)
	assert.Nil(t, err)
	defer server.Stop()

	store, err := gcs.NewStore(ctx, bucket, gcs.WithHTTPClient(server.HTTPClient()))
	assert.Nil(t, err)

	tests := []struct {
		name    string
		path    string
		content string
		wantErr bool
	}{
		{
			name:    "write new content",
			path:    "path/to/fileB.txt",
			content: "foo",
			wantErr: false,
		},
		{
			name:    "overwrite content",
			path:    "path/to/fileA.txt",
			content: "bar",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Put(ctx, tt.path, []byte(tt.content))
			assert.Equal(t, tt.wantErr, err != nil)

			got, err := store.Get(ctx, tt.path)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, []byte(tt.content), got)
		})
	}
}

func TestList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	bucket := "test"
	now := time.Now()
	server, err := newEmulator(
		bucket, map[string]string{
			"path/to/fileA.txt": "foo",
			"path/to/fileB.txt": "hello",
		},
		now,
	)
	assert.Nil(t, err)
	defer server.Stop()

	store, err := gcs.NewStore(ctx, bucket, gcs.WithHTTPClient(server.HTTPClient()))
	assert.Nil(t, err)

	tests := []struct {
		name    string
		prefix  string
		want    []filestore.ObjectAttrs
		wantErr bool
	}{
		{
			name:   "found contents",
			prefix: "path/to",
			want: []filestore.ObjectAttrs{
				{
					Path:      "path/to/fileA.txt",
					Size:      3,
					UpdatedAt: now.Unix(),
				},
				{
					Path:      "path/to/fileB.txt",
					Size:      5,
					UpdatedAt: now.Unix(),
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
			got, err := store.List(ctx, tt.prefix)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
