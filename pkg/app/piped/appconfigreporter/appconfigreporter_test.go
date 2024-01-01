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

package appconfigreporter

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type fakeApplicationLister struct {
	apps []*model.Application
}

func (f *fakeApplicationLister) List() []*model.Application {
	return f.apps
}

func TestReporter_findRegisteredApps(t *testing.T) {
	t.Parallel()

	type args struct {
		repoPath string
		repoID   string
	}
	testcases := []struct {
		name     string
		reporter *Reporter
		args     args
		want     []*model.ApplicationInfo
		wantErr  bool
	}{
		{
			name: "no app registered",
			reporter: &Reporter{
				applicationLister: &fakeApplicationLister{},
				logger:            zap.NewNop(),
			},
			args: args{
				repoPath: "path/to/repo-1",
				repoID:   "repo-1",
			},
			want:    []*model.ApplicationInfo{},
			wantErr: false,
		},
		{
			name: "no app registered in the repo",
			reporter: &Reporter{
				applicationLister: &fakeApplicationLister{apps: []*model.Application{
					{Id: "id-1", Name: "app-1", Labels: map[string]string{"key-1": "value-1"}, GitPath: &model.ApplicationGitPath{Repo: &model.ApplicationGitRepository{Id: "different-repo"}, Path: "app-1", ConfigFilename: "app.pipecd.yaml"}},
				}},
				logger: zap.NewNop(),
			},
			args: args{
				repoPath: "path/to/repo-1",
				repoID:   "repo-1",
			},
			want:    []*model.ApplicationInfo{},
			wantErr: false,
		},
		{
			name: "invalid app config is contained",
			reporter: &Reporter{
				applicationLister: &fakeApplicationLister{apps: []*model.Application{
					{Id: "id-1", Name: "app-1", Labels: map[string]string{"key-1": "value-1"}, GitPath: &model.ApplicationGitPath{Repo: &model.ApplicationGitRepository{Id: "repo-1"}, Path: "app-1", ConfigFilename: "app.pipecd.yaml"}},
				}},
				fileSystem: fstest.MapFS{
					"path/to/repo-1/app-1/app.pipecd.yaml": &fstest.MapFile{Data: []byte("invalid-text")},
				},
				logger: zap.NewNop(),
			},
			args: args{
				repoPath: "path/to/repo-1",
				repoID:   "repo-1",
			},
			want:    []*model.ApplicationInfo{},
			wantErr: false,
		},
		{
			name: "app not changed",
			reporter: &Reporter{
				config: &config.PipedSpec{PipedID: "piped-1"},
				applicationLister: &fakeApplicationLister{apps: []*model.Application{
					{Id: "id-1", Name: "app-1", Labels: map[string]string{"key-1": "value-1"}, GitPath: &model.ApplicationGitPath{Repo: &model.ApplicationGitRepository{Id: "repo-1"}, Path: "app-1", ConfigFilename: "app.pipecd.yaml"}},
				}},
				fileSystem: fstest.MapFS{
					"path/to/repo-1/app-1/app.pipecd.yaml": &fstest.MapFile{Data: []byte(`
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: app-1
  labels:
    key-1: value-1`)},
				},
				lastScannedCommits: make(map[string]string),
				logger:             zap.NewNop(),
			},
			args: args{
				repoPath: "path/to/repo-1",
				repoID:   "repo-1",
			},
			want:    []*model.ApplicationInfo{},
			wantErr: false,
		},
		{
			name: "app changed",
			reporter: &Reporter{
				config: &config.PipedSpec{PipedID: "piped-1"},
				applicationLister: &fakeApplicationLister{apps: []*model.Application{
					{Id: "id-1", Name: "app-1", Labels: map[string]string{"key-1": "value-1"}, GitPath: &model.ApplicationGitPath{Repo: &model.ApplicationGitRepository{Id: "repo-1"}, Path: "app-1", ConfigFilename: "app.pipecd.yaml"}},
				}},
				fileSystem: fstest.MapFS{
					"path/to/repo-1/app-1/app.pipecd.yaml": &fstest.MapFile{Data: []byte(`
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: new-app-1
  labels:
    key-1: value-1`)},
				},
				logger: zap.NewNop(),
			},
			args: args{
				repoPath: "path/to/repo-1",
				repoID:   "repo-1",
			},
			want: []*model.ApplicationInfo{
				{
					Id:             "id-1",
					Name:           "new-app-1",
					Labels:         map[string]string{"key-1": "value-1"},
					RepoId:         "repo-1",
					Path:           "app-1",
					ConfigFilename: "app.pipecd.yaml",
					PipedId:        "piped-1",
				},
			},
			wantErr: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.reporter.findOutOfSyncRegisteredApps(tc.args.repoPath, tc.args.repoID, "not-existed-head-commit")
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestReporter_findUnregisteredApps(t *testing.T) {
	t.Parallel()

	type args struct {
		registeredAppPaths map[string]string
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
			reporter: &Reporter{
				applicationLister: &fakeApplicationLister{},
				fileSystem: fstest.MapFS{
					"path/to/repo-1/app-1/app.pipecd.yaml": &fstest.MapFile{Data: []byte("")},
				},
				logger: zap.NewNop(),
				config: &config.PipedSpec{
					PipedID: "piped-1",
				},
			},
			args: args{
				repoPath:           "invalid",
				repoID:             "repo-1",
				registeredAppPaths: map[string]string{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "all are registered",
			reporter: &Reporter{
				applicationLister: &fakeApplicationLister{},
				fileSystem: fstest.MapFS{
					"path/to/repo-1/app-1/app.pipecd.yaml": &fstest.MapFile{Data: []byte("")},
				},
				logger: zap.NewNop(),
				config: &config.PipedSpec{
					PipedID: "piped-1",
				},
			},
			args: args{
				repoPath: "path/to/repo-1",
				repoID:   "repo-1",
				registeredAppPaths: map[string]string{
					"repo-1:app-1/app.pipecd.yaml": "id-1",
				},
			},
			want:    []*model.ApplicationInfo{},
			wantErr: false,
		},
		{
			name: "invalid app config is contained",
			reporter: &Reporter{
				applicationLister: &fakeApplicationLister{},
				fileSystem: fstest.MapFS{
					"path/to/repo-1/app-1/app.pipecd.yaml": &fstest.MapFile{Data: []byte("invalid-text")},
				},
				config: &config.PipedSpec{
					PipedID: "piped-1",
				},
				logger: zap.NewNop(),
			},
			args: args{
				repoPath:           "path/to/repo-1",
				repoID:             "repo-1",
				registeredAppPaths: map[string]string{},
			},
			want:    []*model.ApplicationInfo{},
			wantErr: false,
		},
		{
			name: "valid app config that is unregistered",
			reporter: &Reporter{
				config: &config.PipedSpec{
					PipedID: "piped-1",
				},
				applicationLister: &fakeApplicationLister{},
				fileSystem: fstest.MapFS{
					"path/to/repo-1/app-1/app.pipecd.yaml": &fstest.MapFile{Data: []byte(`
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: app-1
  labels:
    key-1: value-1`)},
				},
				logger: zap.NewNop(),
			},
			args: args{
				repoPath:           "path/to/repo-1",
				repoID:             "repo-1",
				registeredAppPaths: map[string]string{},
			},
			want: []*model.ApplicationInfo{
				{
					Name:           "app-1",
					Labels:         map[string]string{"key-1": "value-1"},
					RepoId:         "repo-1",
					Path:           "app-1",
					ConfigFilename: "app.pipecd.yaml",
					PipedId:        "piped-1",
				},
			},
			wantErr: false,
		},
		{
			name: "valid app config that name isn't default",
			reporter: &Reporter{
				config: &config.PipedSpec{
					PipedID: "piped-1",
				},
				applicationLister: &fakeApplicationLister{},
				fileSystem: fstest.MapFS{
					"path/to/repo-1/app-1/dev.pipecd.yaml": &fstest.MapFile{Data: []byte(`
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: app-1
  labels:
    key-1: value-1`)},
				},
				logger: zap.NewNop(),
			},
			args: args{
				repoPath:           "path/to/repo-1",
				repoID:             "repo-1",
				registeredAppPaths: map[string]string{},
			},
			want: []*model.ApplicationInfo{
				{
					Name:           "app-1",
					Labels:         map[string]string{"key-1": "value-1"},
					RepoId:         "repo-1",
					Path:           "app-1",
					ConfigFilename: "dev.pipecd.yaml",
					PipedId:        "piped-1",
				},
			},
			wantErr: false,
		},
		{
			name: "filtered by appSelector",
			reporter: &Reporter{
				config: &config.PipedSpec{
					PipedID: "piped-1",
					AppSelector: map[string]string{
						"env": "test",
					},
				},
				applicationLister: &fakeApplicationLister{},
				fileSystem: fstest.MapFS{
					"path/to/repo-1/app-1/app.pipecd.yaml": &fstest.MapFile{Data: []byte(`
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: app-1
  labels:
    key-1: value-1
    env: dev
`)},
				},
				logger: zap.NewNop(),
			},
			args: args{
				repoPath:           "path/to/repo-1",
				repoID:             "repo-1",
				registeredAppPaths: map[string]string{},
			},
			want:    []*model.ApplicationInfo{},
			wantErr: false,
		},
		{
			name: "match labels with appSelector",
			reporter: &Reporter{
				config: &config.PipedSpec{
					PipedID: "piped-1",
					AppSelector: map[string]string{
						"env": "test",
					},
				},
				applicationLister: &fakeApplicationLister{},
				fileSystem: fstest.MapFS{
					"path/to/repo-1/app-1/dev.pipecd.yaml": &fstest.MapFile{Data: []byte(`
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: app-1
  labels:
    key-1: value-1
    env: test
`)},
				},
				logger: zap.NewNop(),
			},
			args: args{
				repoPath:           "path/to/repo-1",
				repoID:             "repo-1",
				registeredAppPaths: map[string]string{},
			},
			want: []*model.ApplicationInfo{
				{
					Name: "app-1",
					Labels: map[string]string{
						"key-1": "value-1",
						"env":   "test",
					},
					RepoId:         "repo-1",
					Path:           "app-1",
					ConfigFilename: "dev.pipecd.yaml",
					PipedId:        "piped-1",
				},
			},
			wantErr: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.reporter.findUnregisteredApps(tc.args.repoPath, tc.args.repoID)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}
