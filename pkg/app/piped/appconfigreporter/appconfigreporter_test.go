// Copyright 2021 The PipeCD Authors.
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

package appconfigreporter

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/model"
)

func Test_findUnregisteredApps(t *testing.T) {
	type args struct {
		registeredAppPaths map[string]struct{}
		fileSystem         fs.FS
		repoPath, repoID   string
	}
	testcases := []struct {
		name     string
		reporter *Reporter
		args     args
		want     []*model.ApplicationInfo
		wantErr  bool
	}{
		{
			name: "file not found",
			args: args{
				fileSystem: fstest.MapFS{
					"path/to/repo-1/app-1/app.pipecd.yaml": &fstest.MapFile{Data: []byte("")},
				},
				repoPath:           "invalid",
				repoID:             "repo-1",
				registeredAppPaths: map[string]struct{}{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "all are registered",
			args: args{
				fileSystem: fstest.MapFS{
					"path/to/repo-1/app-1/app.pipecd.yaml": &fstest.MapFile{Data: []byte("")},
				},
				repoPath: "path/to/repo-1",
				repoID:   "repo-1",
				registeredAppPaths: map[string]struct{}{
					"repo-1:path/to/repo-1/app.pipecd.yaml": {},
				},
			},
			want:    []*model.ApplicationInfo{},
			wantErr: false,
		},
		{
			name: "invalid app config is contained",
			args: args{
				fileSystem: fstest.MapFS{
					"path/to/repo-1/app-1/app.pipecd.yaml": &fstest.MapFile{Data: []byte("invalid-text")},
				},
				repoPath:           "path/to/repo-1",
				repoID:             "repo-1",
				registeredAppPaths: map[string]struct{}{},
			},
			want:    []*model.ApplicationInfo{},
			wantErr: false,
		},
		{
			name: "valid app config that is unregistered",
			args: args{
				fileSystem: fstest.MapFS{
					"path/to/repo-1/app-1/app.pipecd.yaml": &fstest.MapFile{Data: []byte(`
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: app-1
  labels:
    key-1: value-1`)},
				},
				repoPath:           "path/to/repo-1",
				repoID:             "repo-1",
				registeredAppPaths: map[string]struct{}{},
			},
			want: []*model.ApplicationInfo{
				{
					Name:           "app-1",
					Labels:         map[string]string{"key-1": "value-1"},
					Path:           "path/to/repo-1",
					ConfigFilename: "app.pipecd.yaml",
				},
			},
			wantErr: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := findUnregisteredApps(tc.args.fileSystem, tc.args.repoPath, tc.args.repoID, tc.args.registeredAppPaths, zap.NewNop())
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}
