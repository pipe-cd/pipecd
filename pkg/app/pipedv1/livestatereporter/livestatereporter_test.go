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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
	pluginapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/livestate"
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

// TODO: make lib for fakePlugin to use in other tests
type fakePlugin struct {
	pluginapi.PluginClient
	name           string
	syncStrategy   *deployment.DetermineStrategyResponse
	quickStages    []*model.PipelineStage
	pipelineStages []*model.PipelineStage
	rollbackStages []*model.PipelineStage
}

func (p *fakePlugin) Close() error { return nil }
func (p *fakePlugin) BuildQuickSyncStages(ctx context.Context, req *deployment.BuildQuickSyncStagesRequest, opts ...grpc.CallOption) (*deployment.BuildQuickSyncStagesResponse, error) {
	if req.Rollback {
		return &deployment.BuildQuickSyncStagesResponse{
			Stages: append(p.quickStages, p.rollbackStages...),
		}, nil
	}
	return &deployment.BuildQuickSyncStagesResponse{
		Stages: p.quickStages,
	}, nil
}
func (p *fakePlugin) BuildPipelineSyncStages(ctx context.Context, req *deployment.BuildPipelineSyncStagesRequest, opts ...grpc.CallOption) (*deployment.BuildPipelineSyncStagesResponse, error) {
	getIndex := func(stageName string) int32 {
		for _, s := range req.Stages {
			if s.Name == stageName {
				return s.Index
			}
		}
		return -1
	}

	for _, s := range p.pipelineStages {
		s.Index = getIndex(s.Name)
	}

	if req.Rollback {
		return &deployment.BuildPipelineSyncStagesResponse{
			Stages: append(p.pipelineStages, p.rollbackStages...),
		}, nil
	}
	return &deployment.BuildPipelineSyncStagesResponse{
		Stages: p.pipelineStages,
	}, nil
}
func (p *fakePlugin) DetermineStrategy(ctx context.Context, req *deployment.DetermineStrategyRequest, opts ...grpc.CallOption) (*deployment.DetermineStrategyResponse, error) {
	return p.syncStrategy, nil
}
func (p *fakePlugin) DetermineVersions(ctx context.Context, req *deployment.DetermineVersionsRequest, opts ...grpc.CallOption) (*deployment.DetermineVersionsResponse, error) {
	return &deployment.DetermineVersionsResponse{
		Versions: []*model.ArtifactVersion{},
	}, nil
}
func (p *fakePlugin) FetchDefinedStages(ctx context.Context, req *deployment.FetchDefinedStagesRequest, opts ...grpc.CallOption) (*deployment.FetchDefinedStagesResponse, error) {
	stages := make([]string, 0, len(p.quickStages)+len(p.pipelineStages)+len(p.rollbackStages))

	for _, s := range p.quickStages {
		stages = append(stages, s.Name)
	}
	for _, s := range p.pipelineStages {
		stages = append(stages, s.Name)
	}
	for _, s := range p.rollbackStages {
		stages = append(stages, s.Name)
	}
	return &deployment.FetchDefinedStagesResponse{
		Stages: stages,
	}, nil
}
func (p *fakePlugin) GetLivestate(ctx context.Context, in *livestate.GetLivestateRequest, opts ...grpc.CallOption) (*livestate.GetLivestateResponse, error) {
	return &livestate.GetLivestateResponse{
		ApplicationLiveState: &model.ApplicationLiveState{},
		SyncState:            &model.ApplicationSyncState{},
	}, nil
}
func (p *fakePlugin) Name() string {
	return p.name
}

type fakeAPILister struct {
	applicationLister
	apps []*model.Application
}

func (f *fakeAPILister) List() []*model.Application {
	return f.apps
}

func Test_reporter_flushSnapshots(t *testing.T) {
	// NOTE: This test only checks no panic

	gitClient, err := git.NewClient()
	require.NoError(t, err)

	pr := &reporter{
		snapshotFlushInterval: 1 * time.Minute,
		appLister: &fakeAPILister{
			apps: []*model.Application{
				{
					Id:   "app-id",
					Name: "app-name",
					GitPath: &model.ApplicationGitPath{
						Repo: &model.ApplicationGitRepository{
							Id:     "repo-id",
							Remote: "https://github.com/pipe-cd/examples.git",
							Branch: "master",
						},
						Path:           "kubernetes/canary",
						ConfigFilename: "app.pipecd.yaml",
					},
				},
			},
		},
		apiClient: &fakeAPIClient{},
		repoMap:   make(map[string]git.Repo, 0),
		pluginRegistry: func() plugin.PluginRegistry {
			r, err := plugin.NewPluginRegistry(
				context.Background(),
				[]plugin.Plugin{
					{
						Name: "k8s",
						Cli: &fakePlugin{
							pipelineStages: []*model.PipelineStage{
								{
									Name: "K8S_CANARY_ROLLOUT",
								},
								{
									Name: "K8S_CANARY_CLEAN",
								},
								{
									Name: "K8S_PRIMARY_ROLLOUT",
								},
							},
						},
					},
					{
						Name: "wait",
						Cli: &fakePlugin{
							pipelineStages: []*model.PipelineStage{
								{
									Name: "WAIT",
								},
							},
						},
					},
				},
			)

			require.NoError(t, err)
			return r
		}(),
		gitClient: gitClient,
		pipedConfig: &config.PipedSpec{
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

// createLocalTestRepo creates a local git repository with the required files for testing.
// It returns the path to the repository.
func createLocalTestRepo(t *testing.T, repoDir string, files map[string]string) {
	t.Helper()

	gitPath, err := exec.LookPath("git")
	require.NoError(t, err)

	require.NoError(t, os.MkdirAll(repoDir, 0o755))

	run := func(args ...string) {
		cmd := exec.Command(gitPath, args...)
		cmd.Dir = repoDir
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "git %v failed: %s", args, out)
	}

	run("init", "--initial-branch", "master")
	run("config", "user.email", "test@example.com")
	run("config", "user.name", "Test User")

	for path, content := range files {
		fullPath := filepath.Join(repoDir, path)
		require.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0o755))
		require.NoError(t, os.WriteFile(fullPath, []byte(content), 0o644))
	}

	run("add", ".")
	run("commit", "-m", "Initial commit")
}

func Test_reporter_flush(t *testing.T) {
	workingDir := t.TempDir()

	const (
		repoID     = "repo-id"
		appPath    = "kubernetes/canary"
		configFile = "app.pipecd.yaml"
	)

	appPipecdYAML := `apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  plugins:
    k8s: {}
`

	remoteDir := t.TempDir()
	repoDir := filepath.Join(remoteDir, repoID)
	createLocalTestRepo(t, repoDir, map[string]string{
		filepath.Join(appPath, configFile): appPipecdYAML,
	})

	gitClient, err := git.NewClient()
	require.NoError(t, err)

	cfgRepo := config.PipedRepository{
		RepoID: repoID,
		Remote: repoDir,
		Branch: "master",
	}

	repo, err := gitClient.Clone(context.Background(), cfgRepo.RepoID, cfgRepo.Remote, cfgRepo.Branch, fmt.Sprintf("%s/%s", workingDir, cfgRepo.RepoID))
	require.NoError(t, err)

	pr := &reporter{
		snapshotFlushInterval: 1 * time.Minute,
		apiClient:             &fakeAPIClient{},
		pluginRegistry: func() plugin.PluginRegistry {
			r, err := plugin.NewPluginRegistry(
				context.Background(),
				[]plugin.Plugin{
					{
						Name: "k8s",
						Cli: &fakePlugin{
							pipelineStages: []*model.PipelineStage{
								{
									Name: "K8S_CANARY_ROLLOUT",
								},
								{
									Name: "K8S_CANARY_CLEAN",
								},
								{
									Name: "K8S_PRIMARY_ROLLOUT",
								},
							},
						},
					},
					{
						Name: "wait",
						Cli: &fakePlugin{
							pipelineStages: []*model.PipelineStage{
								{
									Name: "WAIT",
								},
							},
						},
					},
				},
			)

			require.NoError(t, err)
			return r
		}(),
		gitClient: gitClient,
		pipedConfig: &config.PipedSpec{
			Repositories: []config.PipedRepository{
				{
					RepoID: repoID,
					Remote: repoDir,
					Branch: "master",
				},
			},
		},
		secretDecrypter: nil,
		workingDir:      workingDir,
		logger:          zaptest.NewLogger(t),
	}

	app := &model.Application{
		Id:   "app-id",
		Name: "app-name",
		GitPath: &model.ApplicationGitPath{
			Repo: &model.ApplicationGitRepository{
				Id:     repoID,
				Remote: repoDir,
				Branch: "master",
			},
			Path:           appPath,
			ConfigFilename: configFile,
		},
	}

	err = pr.flush(context.Background(), app, repo)
	assert.NoError(t, err)
}

func Benchmark_reporter_flushSnapshots(b *testing.B) {
	gitClient, err := git.NewClient()
	require.NoError(b, err)

	pr := &reporter{
		snapshotFlushInterval: 1 * time.Minute,
		appLister: &fakeAPILister{
			apps: func() []*model.Application {
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
							Path:           "kubernetes/canary",
							ConfigFilename: "app.pipecd.yaml",
						},
					})
				}

				return apps
			}(),
		},
		apiClient: &fakeAPIClient{},
		pluginRegistry: func() plugin.PluginRegistry {
			r, err := plugin.NewPluginRegistry(
				context.Background(),
				[]plugin.Plugin{
					{
						Name: "k8s",
						Cli: &fakePlugin{
							pipelineStages: []*model.PipelineStage{
								{
									Name: "K8S_CANARY_ROLLOUT",
								},
								{
									Name: "K8S_CANARY_CLEAN",
								},
								{
									Name: "K8S_PRIMARY_ROLLOUT",
								},
							},
						},
					},
					{
						Name: "wait",
						Cli: &fakePlugin{
							pipelineStages: []*model.PipelineStage{
								{
									Name: "WAIT",
								},
							},
						},
					},
				},
			)

			require.NoError(b, err)
			return r
		}(),
		gitClient: gitClient,
		repoMap:   make(map[string]git.Repo, 0),
		pipedConfig: &config.PipedSpec{
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

	for range b.N {
		pr.flushSnapshots(context.Background())
	}
}
