// Copyright 2025 The PipeCD Authors.
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

package planpreview

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
	pluginapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	planpreviewapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/planpreview"
)

type fakeApplicationLister struct {
	apps []*model.Application
}

func (l *fakeApplicationLister) List() []*model.Application {
	return l.apps
}

func TestListApplications(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		repo     config.PipedRepository
		apps     []*model.Application
		expected []*model.Application
	}{
		{
			name: "no applications",
			repo: config.PipedRepository{
				RepoID: "repo-1",
				Remote: "git@github.com:org/repo-1.git",
				Branch: "main",
			},
			apps:     []*model.Application{},
			expected: []*model.Application{},
		},
		{
			name: "filter by repository configuration",
			repo: config.PipedRepository{
				RepoID: "repo-1",
				Remote: "git@github.com:org/repo-1.git",
				Branch: "main",
			},
			apps: []*model.Application{
				{
					Id: "app-1",
					GitPath: &model.ApplicationGitPath{
						Repo: &model.ApplicationGitRepository{
							Id:     "repo-1",
							Remote: "git@github.com:org/repo-1.git",
							Branch: "main",
						},
					},
				},
				{
					Id: "app-2",
					GitPath: &model.ApplicationGitPath{
						Repo: &model.ApplicationGitRepository{
							Id:     "repo-2", // Different repo ID
							Remote: "git@github.com:org/repo-1.git",
							Branch: "main",
						},
					},
				},
				{
					Id: "app-3",
					GitPath: &model.ApplicationGitPath{
						Repo: &model.ApplicationGitRepository{
							Id:     "repo-1",
							Remote: "git@github.com:org/repo-2.git", // Different remote
							Branch: "main",
						},
					},
				},
				{
					Id: "app-4",
					GitPath: &model.ApplicationGitPath{
						Repo: &model.ApplicationGitRepository{
							Id:     "repo-1",
							Remote: "git@github.com:org/repo-1.git",
							Branch: "develop", // Different branch
						},
					},
				},
			},
			expected: []*model.Application{
				{
					Id: "app-1",
					GitPath: &model.ApplicationGitPath{
						Repo: &model.ApplicationGitRepository{
							Id:     "repo-1",
							Remote: "git@github.com:org/repo-1.git",
							Branch: "main",
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			b := &builder{
				applicationLister: &fakeApplicationLister{
					apps: tc.apps,
				},
			}

			got := b.listApplications(tc.repo)
			assert.Equal(t, tc.expected, got)
		})
	}
}

type fakeAPIClient struct {
	resp *pipedservice.GetApplicationMostRecentDeploymentResponse
	err  error
}

func (c *fakeAPIClient) GetApplicationMostRecentDeployment(_ context.Context, _ *pipedservice.GetApplicationMostRecentDeploymentRequest, _ ...grpc.CallOption) (*pipedservice.GetApplicationMostRecentDeploymentResponse, error) {
	return c.resp, c.err
}

type fakeGitWorktree struct {
	git.Worktree
}

func (w *fakeGitWorktree) Checkout(_ context.Context, _ string) error {
	return nil
}

func (w *fakeGitWorktree) Clean() error {
	return nil
}

type fakeGitRepo struct {
	git.Repo
	appConfig string
}

func (r *fakeGitRepo) Copy(dest string) (git.Worktree, error) {
	appDir := filepath.Join(dest, "app")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(appDir, "app.pipecd.yaml"), []byte(r.appConfig), 0o644); err != nil {
		return nil, err
	}
	return &fakeGitWorktree{}, nil
}

type fakePluginRegistry struct {
	plugin.PluginRegistry
	clients []pluginapi.PluginClient
	err     error
}

func (r *fakePluginRegistry) GetPluginClientsByAppConfig(_ *config.GenericApplicationSpec) ([]pluginapi.PluginClient, error) {
	return r.clients, r.err
}

type fakePluginClient struct {
	pluginapi.PluginClient
	name         string
	noChange     bool
	emptyResults bool
	err          error
}

func (p *fakePluginClient) Name() string {
	return p.name
}

func (p *fakePluginClient) DetermineStrategy(_ context.Context, _ *deployment.DetermineStrategyRequest, _ ...grpc.CallOption) (*deployment.DetermineStrategyResponse, error) {
	return &deployment.DetermineStrategyResponse{
		SyncStrategy: model.SyncStrategy_QUICK_SYNC,
	}, nil
}

func (p *fakePluginClient) GetPlanPreview(_ context.Context, _ *planpreviewapi.GetPlanPreviewRequest, _ ...grpc.CallOption) (*planpreviewapi.GetPlanPreviewResponse, error) {
	if p.err != nil {
		return nil, p.err
	}
	if p.emptyResults {
		return &planpreviewapi.GetPlanPreviewResponse{}, nil
	}
	return &planpreviewapi.GetPlanPreviewResponse{
		Results: []*planpreviewapi.PlanPreviewResult{
			{
				DeployTarget: "target",
				Summary:      "summary",
				NoChange:     p.noChange,
			},
		},
	}, nil
}

func TestBuilder_buildApp_NoChange(t *testing.T) {
	t.Parallel()

	const appConfig = "apiVersion: pipecd.dev/v1beta1\nkind: Application\nspec:\n  name: test-app\n"

	newTestApp := func() *model.Application {
		return &model.Application{
			Id:   "app-1",
			Name: "test-app",
			GitPath: &model.ApplicationGitPath{
				Path:           "app",
				ConfigFilename: "app.pipecd.yaml",
			},
		}
	}

	testcases := []struct {
		name         string
		plugins      []pluginapi.PluginClient
		wantNoChange bool
		wantErr      bool
	}{
		{
			name: "all plugins report no change",
			plugins: []pluginapi.PluginClient{
				&fakePluginClient{name: "plugin-1", noChange: true},
				&fakePluginClient{name: "plugin-2", noChange: true},
			},
			wantNoChange: true,
		},
		{
			name: "mixed plugin no change results",
			plugins: []pluginapi.PluginClient{
				&fakePluginClient{name: "plugin-1", noChange: true},
				&fakePluginClient{name: "plugin-2", noChange: false},
			},
			wantNoChange: false,
		},
		{
			name: "plugin returns error",
			plugins: []pluginapi.PluginClient{
				&fakePluginClient{name: "plugin-1", noChange: true},
				&fakePluginClient{name: "plugin-2", err: errors.New("preview failed")},
			},
			wantNoChange: false,
			wantErr:      true,
		},
		{
			name: "all plugins return Unimplemented",
			plugins: []pluginapi.PluginClient{
				&fakePluginClient{name: "plugin-1", err: status.Error(codes.Unimplemented, "unimplemented")},
			},
			wantNoChange: false,
		},
		{
			name: "plugin returns empty results",
			plugins: []pluginapi.PluginClient{
				&fakePluginClient{name: "plugin-1", emptyResults: true},
			},
			wantNoChange: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			b := &builder{
				apiClient:      &fakeAPIClient{err: status.Error(codes.NotFound, "not found")},
				pluginRegistry: &fakePluginRegistry{clients: tc.plugins},
				pipedCfg:       &config.PipedSpec{PipedID: "piped-1"},
				workingDir:     t.TempDir(),
				logger:         zap.NewNop(),
			}

			res := b.buildApp(context.Background(), 1, "cmd-1", newTestApp(), &fakeGitRepo{appConfig: appConfig}, "target-commit")
			if tc.wantErr {
				require.NotEmpty(t, res.Error)
			} else {
				require.Empty(t, res.Error)
			}
			assert.Equal(t, tc.wantNoChange, res.NoChange)
		})
	}
}
