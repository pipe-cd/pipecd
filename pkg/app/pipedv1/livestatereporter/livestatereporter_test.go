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

// Package livestatereporter provides a piped component
// that reports the changes as well as full snapshot about live state of registered applications.
package livestatereporter

import (
	"context"
	"testing"
	"time"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
	pluginapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/livestate"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
)

type fakeAPIClient struct {
	apiClient
}

func (f *fakeAPIClient) ReportApplicationLiveState(ctx context.Context, req *pipedservice.ReportApplicationLiveStateRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationLiveStateResponse, error) {
	return &pipedservice.ReportApplicationLiveStateResponse{}, nil
}

func (f *fakeAPIClient) ReportApplicationSyncState(ctx context.Context, req *pipedservice.ReportApplicationSyncStateRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationSyncStateResponse, error) {
	return &pipedservice.ReportApplicationSyncStateResponse{}, nil
}

type fakePluginClient struct {
	pluginapi.PluginClient
}

func (f *fakePluginClient) GetLivestate(ctx context.Context, in *livestate.GetLivestateRequest, opts ...grpc.CallOption) (*livestate.GetLivestateResponse, error) {
	return &livestate.GetLivestateResponse{
		ApplicationLiveState: &model.ApplicationLiveState{},
		SyncState:            &model.ApplicationSyncState{},
	}, nil
}

type fakeApiLister struct {
	applicationLister
}

func (f *fakeApiLister) ListByPluginName(name string) []*model.Application {
	apps := make([]*model.Application, 0, 100)

	for i := 0; i < 100; i++ {
		apps = append(apps, &model.Application{
			Id:   "app-id",
			Name: "app-name",
			GitPath: &model.ApplicationGitPath{
				Repo: &model.ApplicationGitRepository{
					Id:     "repo-id",
					Remote: "https://github.com/pipe-cd/examples.git",
					Branch: "master",
				},
				Path:           "kubernetes/simple",
				ConfigFilename: "app.pipecd.yaml",
			},
		})
	}

	return apps
}

func Test_pluginReporter_flushSnapshots(t *testing.T) {
	gitClient, err := git.NewClient()
	require.NoError(t, err)

	pr := &pluginReporter{
		pluginName:            "test",
		snapshotFlushInterval: 1 * time.Minute,
		appLister:             &fakeApiLister{},
		apiClient:             &fakeAPIClient{},
		pluginClient:          &fakePluginClient{},
		gitClient:             gitClient,
		pipedConfig: config.PipedSpec{
			Repositories: []config.PipedRepository{
				{
					RepoID: "repo-id",
					Remote: "https://github.com/pipe-cd/examples.git",
					Branch: "master",
				},
			},
		},
		secretDecrypter: nil,
		workingDir:      t.TempDir(),
		logger:          zaptest.NewLogger(t),
	}

	pr.flushSnapshots(context.Background())
}

func Benchmark_pluginReporter_flushSnapshots(b *testing.B) {
	gitClient, _ := git.NewClient()

	pr := &pluginReporter{
		pluginName:            "test",
		snapshotFlushInterval: 1 * time.Minute,
		appLister:             &fakeApiLister{},
		apiClient:             &fakeAPIClient{},
		pluginClient:          &fakePluginClient{},
		gitClient:             gitClient,
		pipedConfig: config.PipedSpec{
			Repositories: []config.PipedRepository{
				{
					RepoID: "repo-id",
					Remote: "https://github.com/pipe-cd/examples.git",
					Branch: "master",
				},
			},
		},
		secretDecrypter: nil,
		workingDir:      b.TempDir(),
		logger:          zaptest.NewLogger(b),
	}

	pr.flushSnapshots(context.Background())
}
