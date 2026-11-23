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

type fakeAPIClient struct {
	resp *pipedservice.GetApplicationMostRecentDeploymentResponse
	err  error
}

func (c *fakeAPIClient) GetApplicationMostRecentDeployment(_ context.Context, _ *pipedservice.GetApplicationMostRecentDeploymentRequest, _ ...grpc.CallOption) (*pipedservice.GetApplicationMostRecentDeploymentResponse, error) {
	return c.resp, c.err
}

type fakePluginRegistry struct {
	plugins []pluginapi.PluginClient
	err     error
}

func (r *fakePluginRegistry) GetPluginClientByStageName(_ string) (pluginapi.PluginClient, error) {
	return nil, nil
}

func (r *fakePluginRegistry) GetPluginClientsByAppConfig(_ *config.GenericApplicationSpec) ([]pluginapi.PluginClient, error) {
	return r.plugins, r.err
}

type fakePluginClient struct {
	pluginapi.PluginClient
	name             string
	strategyResp     *deployment.DetermineStrategyResponse
	strategyErr      error
	planPreviewResp  *planpreviewapi.GetPlanPreviewResponse
	planPreviewErr   error
	gotRunningSource bool
}

func (p *fakePluginClient) Name() string {
	return p.name
}

func (p *fakePluginClient) DetermineStrategy(_ context.Context, _ *deployment.DetermineStrategyRequest, _ ...grpc.CallOption) (*deployment.DetermineStrategyResponse, error) {
	return p.strategyResp, p.strategyErr
}

func (p *fakePluginClient) GetPlanPreview(_ context.Context, in *planpreviewapi.GetPlanPreviewRequest, _ ...grpc.CallOption) (*planpreviewapi.GetPlanPreviewResponse, error) {
	p.gotRunningSource = in.GetRunningDeploymentSource() != nil
	return p.planPreviewResp, p.planPreviewErr
}

type fakeGitRepo struct {
	git.Repo
	dir         string
	appPath     string
	configByRev map[string]string
	copyErr     error
	checkoutErr map[string]error
}

func (r *fakeGitRepo) Copy(dest string) (git.Worktree, error) {
	if r.copyErr != nil {
		return nil, r.copyErr
	}
	return &fakeGitRepo{
		dir:         dest,
		appPath:     r.appPath,
		configByRev: r.configByRev,
		checkoutErr: r.checkoutErr,
	}, nil
}

func (r *fakeGitRepo) Checkout(_ context.Context, commitish string) error {
	if err, ok := r.checkoutErr[commitish]; ok && err != nil {
		return err
	}
	cfgContent, ok := r.configByRev[commitish]
	if !ok {
		return errors.New("revision not found")
	}
	appDir := filepath.Join(r.dir, r.appPath)
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(appDir, "app.pipecd.yaml"), []byte(cfgContent), 0o644)
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

func TestBuilder_buildApp(t *testing.T) {
	t.Parallel()

	const validAppConfig = `apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: test-app
`

	newTestApp := func() *model.Application {
		return &model.Application{
			Id:   "app-1",
			Name: "test-app",
			GitPath: &model.ApplicationGitPath{
				Repo: &model.ApplicationGitRepository{
					Id:     "repo-1",
					Remote: "git@github.com:org/repo-1.git",
					Branch: "main",
				},
				Path:           "app",
				ConfigFilename: "app.pipecd.yaml",
			},
		}
	}

	t.Run("failure when finding last successful deployment returns non-NotFound error", func(t *testing.T) {
		t.Parallel()

		b := &builder{
			apiClient: &fakeAPIClient{
				err: status.Error(codes.PermissionDenied, "permission denied"),
			},
			logger: zap.NewNop(),
		}

		res := b.buildApp(context.Background(), 1, "cmd-1", newTestApp(), &fakeGitRepo{}, "target-commit")
		assert.Contains(t, res.Error, "failed while finding the last successful deployment")
		assert.Equal(t, []string{"<unknown>"}, res.PluginNames)
	})

	t.Run("failure when getting target deploy source fails", func(t *testing.T) {
		t.Parallel()

		b := &builder{
			apiClient:  &fakeAPIClient{err: status.Error(codes.NotFound, "not found")},
			workingDir: t.TempDir(),
			logger:     zap.NewNop(),
		}
		repo := &fakeGitRepo{
			copyErr: errors.New("copy failed"),
		}

		res := b.buildApp(context.Background(), 1, "cmd-1", newTestApp(), repo, "target-commit")
		assert.Contains(t, res.Error, "failed to get the target deploy source")
		assert.Equal(t, []string{"<unknown>"}, res.PluginNames)
	})

	t.Run("failure when target application config is invalid", func(t *testing.T) {
		t.Parallel()

		b := &builder{
			apiClient:  &fakeAPIClient{err: status.Error(codes.NotFound, "not found")},
			workingDir: t.TempDir(),
			logger:     zap.NewNop(),
		}
		repo := &fakeGitRepo{
			appPath:     "app",
			configByRev: map[string]string{"target-commit": ":::invalid-yaml"},
		}

		res := b.buildApp(context.Background(), 1, "cmd-1", newTestApp(), repo, "target-commit")
		assert.Contains(t, res.Error, "failed to get the target deploy source")
		assert.Equal(t, []string{"<unknown>"}, res.PluginNames)
	})

	t.Run("failure when getting plugin clients fails", func(t *testing.T) {
		t.Parallel()

		b := &builder{
			apiClient:      &fakeAPIClient{err: status.Error(codes.NotFound, "not found")},
			pluginRegistry: &fakePluginRegistry{err: errors.New("registry error")},
			workingDir:     t.TempDir(),
			logger:         zap.NewNop(),
		}
		repo := &fakeGitRepo{
			appPath:     "app",
			configByRev: map[string]string{"target-commit": validAppConfig},
		}

		res := b.buildApp(context.Background(), 1, "cmd-1", newTestApp(), repo, "target-commit")
		assert.Contains(t, res.Error, "failed to get plugin clients")
		assert.Equal(t, []string{"<unknown>"}, res.PluginNames)
	})

	t.Run("success on first deployment when last deployment is NotFound", func(t *testing.T) {
		t.Parallel()

		plugin1 := &fakePluginClient{
			name: "plugin-1",
			planPreviewResp: &planpreviewapi.GetPlanPreviewResponse{
				Results: []*planpreviewapi.PlanPreviewResult{
					{
						DeployTarget: "dt-1",
						Summary:      "1 to add",
						Details:      []byte("+ resource"),
						DiffLanguage: "hcl",
						NoChange:     false,
					},
				},
			},
		}

		b := &builder{
			apiClient:      &fakeAPIClient{err: status.Error(codes.NotFound, "not found")},
			pluginRegistry: &fakePluginRegistry{plugins: []pluginapi.PluginClient{plugin1}},
			pipedCfg:       &config.PipedSpec{PipedID: "piped-1"},
			repoCfg:        config.PipedRepository{Branch: "main"},
			workingDir:     t.TempDir(),
			logger:         zap.NewNop(),
		}
		repo := &fakeGitRepo{
			appPath:     "app",
			configByRev: map[string]string{"target-commit": validAppConfig},
		}

		res := b.buildApp(context.Background(), 1, "cmd-1", newTestApp(), repo, "target-commit")
		require.Empty(t, res.Error)
		assert.Equal(t, model.SyncStrategy_QUICK_SYNC, res.SyncStrategy)
		assert.Equal(t, []string{"plugin-1"}, res.PluginNames)
		assert.False(t, plugin1.gotRunningSource)
		require.Len(t, res.PluginPlanResults, 1)
		assert.Equal(t, &model.PluginPlanPreviewResult{
			PluginName:   "plugin-1",
			DeployTarget: "dt-1",
			PlanSummary:  []byte("1 to add"),
			PlanDetails:  []byte("+ resource"),
			DiffLanguage: "hcl",
		}, res.PluginPlanResults[0])
	})

	t.Run("success with previous deployment and skips unimplemented plugin", func(t *testing.T) {
		t.Parallel()

		unimplementedPlugin := &fakePluginClient{
			name: "plugin-unimplemented",
			strategyResp: &deployment.DetermineStrategyResponse{
				SyncStrategy: model.SyncStrategy_QUICK_SYNC,
			},
			planPreviewErr: status.Error(codes.Unimplemented, "not implemented"),
		}
		activePlugin := &fakePluginClient{
			name: "plugin-active",
			strategyResp: &deployment.DetermineStrategyResponse{
				SyncStrategy: model.SyncStrategy_QUICK_SYNC,
			},
			planPreviewResp: &planpreviewapi.GetPlanPreviewResponse{
				Results: []*planpreviewapi.PlanPreviewResult{
					{
						DeployTarget: "dt-1",
						Summary:      "no changes",
						Details:      []byte(""),
						NoChange:     true,
					},
				},
			},
		}

		b := &builder{
			apiClient: &fakeAPIClient{
				resp: &pipedservice.GetApplicationMostRecentDeploymentResponse{
					Deployment: &model.ApplicationDeploymentReference{
						Trigger: &model.DeploymentTrigger{
							Commit: &model.Commit{Hash: "running-commit"},
						},
					},
				},
			},
			pluginRegistry: &fakePluginRegistry{
				plugins: []pluginapi.PluginClient{unimplementedPlugin, activePlugin},
			},
			pipedCfg:   &config.PipedSpec{PipedID: "piped-1"},
			repoCfg:    config.PipedRepository{Branch: "main"},
			workingDir: t.TempDir(),
			logger:     zap.NewNop(),
		}
		repo := &fakeGitRepo{
			appPath: "app",
			configByRev: map[string]string{
				"running-commit": validAppConfig,
				"target-commit":  validAppConfig,
			},
		}

		res := b.buildApp(context.Background(), 1, "cmd-1", newTestApp(), repo, "target-commit")
		require.Empty(t, res.Error)
		assert.Equal(t, []string{"plugin-unimplemented", "plugin-active"}, res.PluginNames)
		assert.True(t, activePlugin.gotRunningSource)
		require.Len(t, res.PluginPlanResults, 1)
		assert.Equal(t, "plugin-active", res.PluginPlanResults[0].PluginName)
	})

	t.Run("failure when plugin GetPlanPreview returns an error", func(t *testing.T) {
		t.Parallel()

		failingPlugin := &fakePluginClient{
			name:           "plugin-fail",
			planPreviewErr: status.Error(codes.Internal, "internal failure"),
		}

		b := &builder{
			apiClient:      &fakeAPIClient{err: status.Error(codes.NotFound, "not found")},
			pluginRegistry: &fakePluginRegistry{plugins: []pluginapi.PluginClient{failingPlugin}},
			pipedCfg:       &config.PipedSpec{PipedID: "piped-1"},
			repoCfg:        config.PipedRepository{Branch: "main"},
			workingDir:     t.TempDir(),
			logger:         zap.NewNop(),
		}
		repo := &fakeGitRepo{
			appPath:     "app",
			configByRev: map[string]string{"target-commit": validAppConfig},
		}

		res := b.buildApp(context.Background(), 1, "cmd-1", newTestApp(), repo, "target-commit")
		assert.Contains(t, res.Error, "failed to get plan preview")
		assert.Contains(t, res.Error, "[plugin-fail]")
		assert.Equal(t, []string{"plugin-fail"}, res.PluginNames)
	})
}
