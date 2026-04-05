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
	"fmt"
	"os"
	"os/exec"
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
	planpreviewapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/planpreview"
)

type fakeApplicationLister struct {
	apps []*model.Application
}

func (l *fakeApplicationLister) List() []*model.Application {
	return l.apps
}

type fakeAPIClient struct {
	deployRef *model.ApplicationDeploymentReference
	err       error
}

func (c *fakeAPIClient) GetApplicationMostRecentDeployment(_ context.Context, _ *pipedservice.GetApplicationMostRecentDeploymentRequest, _ ...grpc.CallOption) (*pipedservice.GetApplicationMostRecentDeploymentResponse, error) {
	if c.err != nil {
		return nil, c.err
	}
	return &pipedservice.GetApplicationMostRecentDeploymentResponse{
		Deployment: c.deployRef,
	}, nil
}

type fakeSecretDecrypter struct{}

func (d *fakeSecretDecrypter) Decrypt(text string) (string, error) {
	return text, nil
}

type fakeWorktree struct {
	git.Worktree
	path string
}

func (w *fakeWorktree) GetPath() string                            { return w.path }
func (w *fakeWorktree) Checkout(_ context.Context, _ string) error { return nil }

type fakeRepo struct {
	git.Repo
	path string
}

func (r *fakeRepo) GetPath() string { return r.path }
func (r *fakeRepo) Copy(dest string) (git.Worktree, error) {
	cmd := exec.Command("cp", "-rf", r.path, dest)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("failed to copy: %s, %w", string(out), err)
	}
	return &fakeWorktree{path: dest}, nil
}

type fakePluginRegistry struct {
	plugin.PluginRegistry
	clients []pluginapi.PluginClient
	err     error
}

func (r *fakePluginRegistry) GetPluginClientsByAppConfig(_ *config.GenericApplicationSpec) ([]pluginapi.PluginClient, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.clients, nil
}

type fakePluginClient struct {
	pluginapi.PluginClient
	name           string
	planPreviewRes *planpreviewapi.GetPlanPreviewResponse
	planPreviewErr error
}

func (c *fakePluginClient) Name() string { return c.name }
func (c *fakePluginClient) Close() error { return nil }
func (c *fakePluginClient) GetPlanPreview(_ context.Context, _ *planpreviewapi.GetPlanPreviewRequest, _ ...grpc.CallOption) (*planpreviewapi.GetPlanPreviewResponse, error) {
	return c.planPreviewRes, c.planPreviewErr
}

func setupRepoDir(t *testing.T, appPath, configFilename, configContent string) string {
	t.Helper()

	repoDir := t.TempDir()
	appDir := filepath.Join(repoDir, appPath)
	require.NoError(t, os.MkdirAll(appDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(appDir, configFilename), []byte(configContent), 0644))

	return repoDir
}

func TestBuildApp(t *testing.T) {
	t.Parallel()

	const validAppCfg = `apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: test-app
`

	makeApp := func() *model.Application {
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

	makeBuilder := func(ac apiClient, pr plugin.PluginRegistry, workingDir string) *builder {
		return &builder{
			apiClient:       ac,
			secretDecrypter: &fakeSecretDecrypter{},
			pluginRegistry:  pr,
			pipedCfg: &config.PipedSpec{
				PipedID: "piped-1",
			},
			repoCfg: config.PipedRepository{
				RepoID: "repo-1",
				Remote: "git@github.com:org/repo-1.git",
				Branch: "main",
			},
			workingDir: workingDir,
			logger:     zap.NewNop(),
		}
	}

	testcases := []struct {
		name              string
		apiClient         *fakeAPIClient
		pluginRegistry    *fakePluginRegistry
		setupRepo         bool
		wantError         string
		wantPluginNames   []string
		wantSyncStrategy  model.SyncStrategy
		wantPluginResults []*model.PluginPlanPreviewResult
	}{
		{
			name: "api error on recent deployment",
			apiClient: &fakeAPIClient{
				err: status.Error(codes.Internal, "internal server error"),
			},
			pluginRegistry:  &fakePluginRegistry{},
			wantError:       "failed while finding the last successful deployment",
			wantPluginNames: []string{"<unknown>"},
		},
		{
			name: "missing deploy source",
			apiClient: &fakeAPIClient{
				err: status.Error(codes.NotFound, "not found"),
			},
			pluginRegistry:  &fakePluginRegistry{},
			wantError:       "failed to get the target deploy source",
			wantPluginNames: []string{"<unknown>"},
		},
		{
			name: "plugin registry error",
			apiClient: &fakeAPIClient{
				err: status.Error(codes.NotFound, "not found"),
			},
			pluginRegistry: &fakePluginRegistry{
				err: fmt.Errorf("no plugins available"),
			},
			setupRepo:       true,
			wantError:       "failed to get plugin clients",
			wantPluginNames: []string{"<unknown>"},
		},
		{
			name: "quick sync with no plugins",
			apiClient: &fakeAPIClient{
				err: status.Error(codes.NotFound, "not found"),
			},
			pluginRegistry:   &fakePluginRegistry{},
			setupRepo:        true,
			wantSyncStrategy: model.SyncStrategy_QUICK_SYNC,
			wantPluginNames:  []string{"<unknown>"},
		},
		{
			name: "plugin returns plan preview",
			apiClient: &fakeAPIClient{
				err: status.Error(codes.NotFound, "not found"),
			},
			pluginRegistry: &fakePluginRegistry{
				clients: []pluginapi.PluginClient{
					&fakePluginClient{
						name: "kubernetes",
						planPreviewRes: &planpreviewapi.GetPlanPreviewResponse{
							Results: []*planpreviewapi.PlanPreviewResult{
								{
									DeployTarget: "default",
									Summary:      "Updated deployment image",
									Details:      []byte("--- a\n+++ b\n"),
									DiffLanguage: "diff",
								},
							},
						},
					},
				},
			},
			setupRepo:        true,
			wantSyncStrategy: model.SyncStrategy_QUICK_SYNC,
			wantPluginNames:  []string{"kubernetes"},
			wantPluginResults: []*model.PluginPlanPreviewResult{
				{
					PluginName:   "kubernetes",
					DeployTarget: "default",
					PlanSummary:  []byte("Updated deployment image"),
					PlanDetails:  []byte("--- a\n+++ b\n"),
					DiffLanguage: "diff",
				},
			},
		},
		{
			name: "plugin returns unimplemented",
			apiClient: &fakeAPIClient{
				err: status.Error(codes.NotFound, "not found"),
			},
			pluginRegistry: &fakePluginRegistry{
				clients: []pluginapi.PluginClient{
					&fakePluginClient{
						name:           "cloudrun",
						planPreviewErr: status.Error(codes.Unimplemented, "not implemented"),
					},
				},
			},
			setupRepo:        true,
			wantSyncStrategy: model.SyncStrategy_QUICK_SYNC,
			wantPluginNames:  []string{"cloudrun"},
		},
		{
			name: "plugin returns error",
			apiClient: &fakeAPIClient{
				err: status.Error(codes.NotFound, "not found"),
			},
			pluginRegistry: &fakePluginRegistry{
				clients: []pluginapi.PluginClient{
					&fakePluginClient{
						name:           "kubernetes",
						planPreviewErr: status.Error(codes.Internal, "plugin internal error"),
					},
				},
			},
			setupRepo:       true,
			wantError:       "failed to get plan preview",
			wantPluginNames: []string{"kubernetes"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := makeApp()
			workingDir := t.TempDir()

			var repoDir string
			if tc.setupRepo {
				repoDir = setupRepoDir(t, "app", "app.pipecd.yaml", validAppCfg)
			} else {
				repoDir = t.TempDir()
			}

			b := makeBuilder(tc.apiClient, tc.pluginRegistry, workingDir)
			repo := &fakeRepo{path: repoDir}
			result := b.buildApp(context.Background(), 0, "cmd-1", app, repo, "abc123")

			assert.NotNil(t, result)
			assert.Equal(t, "app-1", result.ApplicationId)
			assert.Equal(t, "test-app", result.ApplicationName)

			if tc.wantError != "" {
				assert.Contains(t, result.Error, tc.wantError)
			} else {
				assert.Empty(t, result.Error)
				assert.Equal(t, tc.wantSyncStrategy, result.SyncStrategy)
			}

			assert.Equal(t, tc.wantPluginNames, result.PluginNames)

			if len(tc.wantPluginResults) > 0 {
				require.Len(t, result.PluginPlanResults, len(tc.wantPluginResults))
				for i, want := range tc.wantPluginResults {
					got := result.PluginPlanResults[i]
					assert.Equal(t, want.PluginName, got.PluginName)
					assert.Equal(t, want.DeployTarget, got.DeployTarget)
					assert.Equal(t, want.PlanSummary, got.PlanSummary)
					assert.Equal(t, want.PlanDetails, got.PlanDetails)
					assert.Equal(t, want.DiffLanguage, got.DiffLanguage)
				}
			}
		})
	}
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
