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
	"context"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/model"
)

type fakeEnvGetter struct {
	env *model.Environment
	err error
}

func (f fakeEnvGetter) Get(_ context.Context, _ string) (*model.Environment, error) {
	return f.env, f.err
}

func (f fakeEnvGetter) GetByName(_ context.Context, _ string) (*model.Environment, error) {
	return f.env, f.err
}

func TestReporter_findUnregisteredApps(t *testing.T) {
	type args struct {
		registeredAppPaths map[string]struct{}
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
				envGetter: &fakeEnvGetter{env: &model.Environment{Id: "env-1"}},
				fsys: fstest.MapFS{
					"path/to/repo-1/app-1/app.pipecd.yaml": &fstest.MapFile{Data: []byte("")},
				},
				logger: zap.NewNop(),
			},
			args: args{
				repoPath:           "invalid",
				repoID:             "repo-1",
				registeredAppPaths: map[string]struct{}{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "all are registered",
			reporter: &Reporter{
				envGetter: &fakeEnvGetter{env: &model.Environment{Id: "env-1"}},
				fsys: fstest.MapFS{
					"path/to/repo-1/app-1/app.pipecd.yaml": &fstest.MapFile{Data: []byte("")},
				},
				logger: zap.NewNop(),
			},
			args: args{
				repoPath: "path/to/repo-1",
				repoID:   "repo-1",
				registeredAppPaths: map[string]struct{}{
					"repo-1:path/to/repo-1/app-1/app.pipecd.yaml": {},
				},
			},
			want:    []*model.ApplicationInfo{},
			wantErr: false,
		},
		{
			name: "invalid app config is contained",
			reporter: &Reporter{
				envGetter: &fakeEnvGetter{env: &model.Environment{Id: "env-1"}},
				fsys: fstest.MapFS{
					"path/to/repo-1/app-1/app.pipecd.yaml": &fstest.MapFile{Data: []byte("invalid-text")},
				},
				logger: zap.NewNop(),
			},
			args: args{
				repoPath:           "path/to/repo-1",
				repoID:             "repo-1",
				registeredAppPaths: map[string]struct{}{},
			},
			want:    []*model.ApplicationInfo{},
			wantErr: false,
		},
		{
			name: "valid app config that is unregistered",
			reporter: &Reporter{
				envGetter: &fakeEnvGetter{env: &model.Environment{Id: "env-1"}},
				fsys: fstest.MapFS{
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
				registeredAppPaths: map[string]struct{}{},
			},
			want: []*model.ApplicationInfo{
				{
					Name:           "app-1",
					EnvId:          "env-1",
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
			got, err := tc.reporter.findUnregisteredApps(context.Background(), tc.args.repoPath, tc.args.repoID, tc.args.registeredAppPaths)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}

type fakeGitRepo struct {
	path         string
	changedFiles []string
	err          error
}

func (f *fakeGitRepo) GetPath() string {
	return f.path
}

func (f *fakeGitRepo) ChangedFiles(_ context.Context, _, _ string) ([]string, error) {
	return f.changedFiles, f.err
}

func TestReporter_findRegisteredApps(t *testing.T) {
	type args struct {
		repoID             string
		repo               gitRepo
		headCommitHash     string
		registeredAppPaths map[string]struct{}
	}
	testcases := []struct {
		name     string
		reporter *Reporter
		args     args
		want     []*model.ApplicationInfo
		wantErr  bool
	}{
		{
			name: "already scanned commit",
			reporter: &Reporter{
				lastScannedCommits: map[string]string{"repo-1": "commit-hash"},
				logger:             zap.NewNop(),
			},
			args: args{
				repoID:         "repo-1",
				repo:           &fakeGitRepo{path: "path/to/repo-1", changedFiles: nil},
				headCommitHash: "commit-hash",
			},
			want:    []*model.ApplicationInfo{},
			wantErr: false,
		},
		{
			name: "no changed file",
			reporter: &Reporter{
				lastScannedCommits: make(map[string]string),
				envGetter:          &fakeEnvGetter{env: &model.Environment{Id: "env-1"}},
				fsys: fstest.MapFS{
					"path/to/repo-1/app-1/app.pipecd.yaml": &fstest.MapFile{Data: []byte("invalid-text")},
				},
				logger: zap.NewNop(),
			},
			args: args{
				repoID: "repo-1",
				repo:   &fakeGitRepo{path: "path/to/repo-1", changedFiles: nil},
			},
			want:    []*model.ApplicationInfo{},
			wantErr: false,
		},
		{
			name: "all are unregistered",
			reporter: &Reporter{
				lastScannedCommits: make(map[string]string),
				envGetter:          &fakeEnvGetter{env: &model.Environment{Id: "env-1"}},
				fsys: fstest.MapFS{
					"path/to/repo-1/app-1/app.pipecd.yaml": &fstest.MapFile{Data: []byte("")},
				},
				logger: zap.NewNop(),
			},
			args: args{
				repoID: "repo-1",
				repo:   &fakeGitRepo{path: "path/to/repo-1", changedFiles: []string{"app-1/app.pipecd.yaml"}},
			},
			want:    []*model.ApplicationInfo{},
			wantErr: false,
		},
		{
			name: "invalid app config is contained",
			reporter: &Reporter{
				lastScannedCommits: make(map[string]string),
				envGetter:          &fakeEnvGetter{env: &model.Environment{Id: "env-1"}},
				fsys: fstest.MapFS{
					"path/to/repo-1/app-1/app.pipecd.yaml": &fstest.MapFile{Data: []byte("invalid-text")},
				},
				logger: zap.NewNop(),
			},
			args: args{
				repoID: "repo-1",
				repo:   &fakeGitRepo{path: "path/to/repo-1", changedFiles: []string{"app-1/app.pipecd.yaml"}},
				registeredAppPaths: map[string]struct{}{
					"repo-1:path/to/repo-1/app-1/app.pipecd.yaml": {},
				},
			},
			want:    []*model.ApplicationInfo{},
			wantErr: false,
		},
		{
			name: "valid app config that is registered",
			reporter: &Reporter{
				lastScannedCommits: make(map[string]string),
				envGetter:          &fakeEnvGetter{env: &model.Environment{Id: "env-1"}},
				fsys: fstest.MapFS{
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
				repoID: "repo-1",
				repo:   &fakeGitRepo{path: "path/to/repo-1", changedFiles: []string{"app-1/app.pipecd.yaml"}},
				registeredAppPaths: map[string]struct{}{
					"repo-1:path/to/repo-1/app-1/app.pipecd.yaml": {},
				},
			},
			want: []*model.ApplicationInfo{
				{
					Name:           "app-1",
					EnvId:          "env-1",
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
			got, err := tc.reporter.findRegisteredApps(context.Background(), tc.args.repoID, tc.args.repo, "", tc.args.registeredAppPaths)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}
