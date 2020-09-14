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
)

func TestGetObject(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	bucket := "test"
	server, err := fakestorage.NewServerWithOptions(fakestorage.Options{
		InitialObjects: []fakestorage.Object{
			{
				BucketName: bucket,
				Name:       "path/to/file.txt",
				Content:    []byte("foo"),
			},
		},
		Host: "127.0.0.1",
		Port: 8081,
	})
	assert.Nil(t, err)
	defer server.Stop()

	store, err := NewStore(ctx, bucket, WithHTTPClient(server.HTTPClient()))
	assert.Nil(t, err)

	tests := []struct {
		name    string
		path    string
		want    filestore.Object
		wantErr bool
	}{
		{
			name: "found content",
			path: "path/to/file.txt",
			want: filestore.Object{
				Path:    "path/to/file.txt",
				Content: []byte("foo"),
				Size:    3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.GetObject(ctx, tt.path)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, got, tt.want)
		})
	}
}
