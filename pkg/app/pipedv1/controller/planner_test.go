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

package controller

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
	pluginapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/registry"
)

type fakePlugin struct {
	pluginapi.PluginClient
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
	getIndex := func(stageID string) int32 {
		for _, s := range req.Stages {
			if s.Id == stageID {
				return s.Index
			}
		}
		return -1
	}

	for _, s := range p.pipelineStages {
		s.Index = getIndex(s.Id)
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

func pointerBool(b bool) *bool {
	return &b
}

func TestBuildQuickSyncStages(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name           string
		pluginRegistry registry.PluginRegistry
		cfg            *config.GenericApplicationSpec
		wantErr        bool
		expectedStages []*model.PipelineStage
	}{
		{
			name: "only one plugin",
			pluginRegistry: func() registry.PluginRegistry {
				pr, err := registry.NewPluginRegistry(context.TODO(), []registry.Plugin{
					{
						Name: "plugin-1",
						Cli: &fakePlugin{
							quickStages: []*model.PipelineStage{
								{
									Id:   "plugin-1-stage-1",
									Name: "plugin-1-stage-1",
								},
							},
							rollbackStages: []*model.PipelineStage{
								{
									Id:       "plugin-1-rollback",
									Name:     "plugin-1-rollback",
									Rollback: true,
								},
							},
						},
					},
				})
				require.NoError(t, err)

				return pr
			}(),
			cfg: &config.GenericApplicationSpec{
				Planner: config.DeploymentPlanner{
					AutoRollback: pointerBool(true),
				},
				Plugins: []string{"plugin-1"},
			},
			wantErr: false,
			expectedStages: []*model.PipelineStage{
				{
					Id:   "plugin-1-stage-1",
					Name: "plugin-1-stage-1",
				},
				{
					Id:       "plugin-1-rollback",
					Name:     "plugin-1-rollback",
					Rollback: true,
				},
			},
		},
		{
			name: "multi plugins",
			pluginRegistry: func() registry.PluginRegistry {
				pr, err := registry.NewPluginRegistry(context.TODO(), []registry.Plugin{
					{
						Name: "plugin-1",
						Cli: &fakePlugin{
							quickStages: []*model.PipelineStage{
								{
									Id:   "plugin-1-stage-1",
									Name: "plugin-1-stage-1",
								},
							},
							rollbackStages: []*model.PipelineStage{
								{
									Id:       "plugin-1-rollback",
									Name:     "plugin-1-rollback",
									Rollback: true,
								},
							},
						},
					},
					{
						Name: "plugin-2",
						Cli: &fakePlugin{
							quickStages: []*model.PipelineStage{
								{
									Id:   "plugin-2-stage-1",
									Name: "plugin-2-stage-1",
								},
							},
							rollbackStages: []*model.PipelineStage{
								{
									Id:       "plugin-2-rollback",
									Name:     "plugin-2-rollback",
									Rollback: true,
								},
							},
						},
					},
				})
				require.NoError(t, err)

				return pr
			}(),
			cfg: &config.GenericApplicationSpec{
				Planner: config.DeploymentPlanner{
					AutoRollback: pointerBool(true),
				},
				Plugins: []string{"plugin-1", "plugin-2"},
			},
			wantErr: false,
			expectedStages: []*model.PipelineStage{
				{
					Id:   "plugin-1-stage-1",
					Name: "plugin-1-stage-1",
				},
				{
					Id:   "plugin-2-stage-1",
					Name: "plugin-2-stage-1",
				},
				{
					Id:       "plugin-1-rollback",
					Name:     "plugin-1-rollback",
					Rollback: true,
				},
				{
					Id:       "plugin-2-rollback",
					Name:     "plugin-2-rollback",
					Rollback: true,
				},
			},
		},
		{
			name: "multi plugins - no rollback",
			pluginRegistry: func() registry.PluginRegistry {
				pr, err := registry.NewPluginRegistry(context.TODO(), []registry.Plugin{
					{
						Name: "plugin-1",
						Cli: &fakePlugin{
							quickStages: []*model.PipelineStage{
								{
									Id:   "plugin-1-stage-1",
									Name: "plugin-1-stage-1",
								},
							},
							rollbackStages: []*model.PipelineStage{
								{
									Id:       "plugin-1-rollback",
									Name:     "plugin-1-rollback",
									Rollback: true,
								},
							},
						},
					},
					{
						Name: "plugin-2",
						Cli: &fakePlugin{
							quickStages: []*model.PipelineStage{
								{
									Id:   "plugin-2-stage-1",
									Name: "plugin-2-stage-1",
								},
							},
							rollbackStages: []*model.PipelineStage{
								{
									Id:       "plugin-2-rollback",
									Name:     "plugin-2-rollback",
									Rollback: true,
								},
							},
						},
					},
				})
				require.NoError(t, err)

				return pr
			}(),
			cfg: &config.GenericApplicationSpec{
				Planner: config.DeploymentPlanner{
					AutoRollback: pointerBool(false),
				},
				Plugins: []string{"plugin-1", "plugin-2"},
			},
			wantErr: false,
			expectedStages: []*model.PipelineStage{
				{
					Id:   "plugin-1-stage-1",
					Name: "plugin-1-stage-1",
				},
				{
					Id:   "plugin-2-stage-1",
					Name: "plugin-2-stage-1",
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			planner := &planner{
				pluginRegistry: tc.pluginRegistry,
			}
			stages, err := planner.buildQuickSyncStages(context.TODO(), tc.cfg)
			require.Equal(t, tc.wantErr, err != nil)
			assert.ElementsMatch(t, tc.expectedStages, stages)
		})
	}
}

func TestBuildPipelineSyncStages(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name           string
		pluginRegistry registry.PluginRegistry
		cfg            *config.GenericApplicationSpec
		wantErr        bool
		expectedStages []*model.PipelineStage
	}{
		{
			name: "only one plugin",
			pluginRegistry: func() registry.PluginRegistry {
				pr, err := registry.NewPluginRegistry(context.TODO(), []registry.Plugin{
					{
						Name: "plugin-1",
						Cli: &fakePlugin{
							pipelineStages: []*model.PipelineStage{
								{
									Id:    "plugin-1-stage-1",
									Index: 0,
									Name:  "plugin-1-stage-1",
								},
								{
									Id:    "plugin-1-stage-2",
									Index: 1,
									Name:  "plugin-1-stage-2",
								},
							},
							rollbackStages: []*model.PipelineStage{
								{
									Id:       "plugin-1-rollback",
									Index:    0,
									Name:     "plugin-1-rollback",
									Rollback: true,
								},
							},
						},
					},
				})
				require.NoError(t, err)

				return pr
			}(),
			cfg: &config.GenericApplicationSpec{
				Planner: config.DeploymentPlanner{
					AutoRollback: pointerBool(true),
				},
				Pipeline: &config.DeploymentPipeline{
					Stages: []config.PipelineStage{
						{
							ID:   "plugin-1-stage-1",
							Name: "plugin-1-stage-1",
						},
						{
							ID:   "plugin-1-stage-2",
							Name: "plugin-1-stage-2",
						},
					},
				},
			},
			wantErr: false,
			expectedStages: []*model.PipelineStage{
				{
					Id:    "plugin-1-stage-1",
					Name:  "plugin-1-stage-1",
					Index: 0,
				},
				{
					Id:       "plugin-1-stage-2",
					Name:     "plugin-1-stage-2",
					Index:    1,
					Requires: []string{"plugin-1-stage-1"},
				},
				{
					Id:       "plugin-1-rollback",
					Name:     "plugin-1-rollback",
					Index:    0,
					Rollback: true,
				},
			},
		},
		{
			name: "multi plugins single rollback",
			pluginRegistry: func() registry.PluginRegistry {
				pr, err := registry.NewPluginRegistry(context.TODO(), []registry.Plugin{
					{
						Name: "plugin-1",
						Cli: &fakePlugin{
							pipelineStages: []*model.PipelineStage{
								{
									Id:    "plugin-1-stage-1",
									Index: 0,
									Name:  "plugin-1-stage-1",
								},
								{
									Id:    "plugin-1-stage-2",
									Index: 1,
									Name:  "plugin-1-stage-2",
								},
								{
									Id:    "plugin-1-stage-3",
									Index: 2,
									Name:  "plugin-1-stage-3",
								},
							},
							rollbackStages: []*model.PipelineStage{
								{
									Id:       "plugin-1-rollback",
									Index:    0,
									Name:     "plugin-1-rollback",
									Rollback: true,
								},
							},
						},
					},
					{
						Name: "plugin-2",
						Cli: &fakePlugin{
							pipelineStages: []*model.PipelineStage{
								{
									Id:    "plugin-2-stage-1",
									Index: 0,
									Name:  "plugin-2-stage-1",
								},
								{
									Id:    "plugin-2-stage-2",
									Index: 1,
									Name:  "plugin-2-stage-2",
								},
							},
						},
					},
				})
				require.NoError(t, err)

				return pr
			}(),
			cfg: &config.GenericApplicationSpec{
				Planner: config.DeploymentPlanner{
					AutoRollback: pointerBool(true),
				},
				Pipeline: &config.DeploymentPipeline{
					Stages: []config.PipelineStage{
						{
							ID:   "plugin-1-stage-1",
							Name: "plugin-1-stage-1",
						},
						{
							ID:   "plugin-1-stage-2",
							Name: "plugin-1-stage-2",
						},
						{
							ID:   "plugin-2-stage-1",
							Name: "plugin-2-stage-1",
						},
						{
							ID:   "plugin-1-stage-3",
							Name: "plugin-1-stage-3",
						},
						{
							ID:   "plugin-2-stage-2",
							Name: "plugin-2-stage-2",
						},
					},
				},
			},
			wantErr: false,
			expectedStages: []*model.PipelineStage{
				{
					Id:    "plugin-1-stage-1",
					Name:  "plugin-1-stage-1",
					Index: 0,
				},
				{
					Id:       "plugin-1-stage-2",
					Name:     "plugin-1-stage-2",
					Index:    1,
					Requires: []string{"plugin-1-stage-1"},
				},
				{
					Id:       "plugin-2-stage-1",
					Name:     "plugin-2-stage-1",
					Index:    2,
					Requires: []string{"plugin-1-stage-2"},
				},
				{
					Id:       "plugin-1-stage-3",
					Name:     "plugin-1-stage-3",
					Index:    3,
					Requires: []string{"plugin-2-stage-1"},
				},
				{
					Id:       "plugin-2-stage-2",
					Name:     "plugin-2-stage-2",
					Index:    4,
					Requires: []string{"plugin-1-stage-3"},
				},
				{
					Id:       "plugin-1-rollback",
					Name:     "plugin-1-rollback",
					Index:    0,
					Rollback: true,
				},
			},
		},
		{
			name: "multi plugins multi rollback",
			pluginRegistry: func() registry.PluginRegistry {
				pr, err := registry.NewPluginRegistry(context.TODO(), []registry.Plugin{
					{
						Name: "plugin-1",
						Cli: &fakePlugin{
							pipelineStages: []*model.PipelineStage{
								{
									Id:    "plugin-1-stage-1",
									Index: 0,
									Name:  "plugin-1-stage-1",
								},
								{
									Id:    "plugin-1-stage-2",
									Index: 1,
									Name:  "plugin-1-stage-2",
								},
								{
									Id:    "plugin-1-stage-3",
									Index: 2,
									Name:  "plugin-1-stage-3",
								},
							},
							rollbackStages: []*model.PipelineStage{
								{
									Id:       "plugin-1-rollback",
									Index:    0,
									Name:     "plugin-1-rollback",
									Rollback: true,
								},
							},
						},
					},
					{
						Name: "plugin-2",
						Cli: &fakePlugin{
							pipelineStages: []*model.PipelineStage{
								{
									Id:    "plugin-2-stage-1",
									Index: 0,
									Name:  "plugin-2-stage-1",
								},
							},
							rollbackStages: []*model.PipelineStage{
								{
									Id:       "plugin-2-rollback",
									Index:    2,
									Name:     "plugin-2-rollback",
									Rollback: true,
								},
							},
						},
					},
				})
				require.NoError(t, err)

				return pr
			}(),
			cfg: &config.GenericApplicationSpec{
				Planner: config.DeploymentPlanner{
					AutoRollback: pointerBool(true),
				},
				Pipeline: &config.DeploymentPipeline{
					Stages: []config.PipelineStage{
						{
							ID:   "plugin-1-stage-1",
							Name: "plugin-1-stage-1",
						},
						{
							ID:   "plugin-1-stage-2",
							Name: "plugin-1-stage-2",
						},
						{
							ID:   "plugin-2-stage-1",
							Name: "plugin-2-stage-1",
						},
						{
							ID:   "plugin-1-stage-3",
							Name: "plugin-1-stage-3",
						},
					},
				},
			},
			wantErr: false,
			expectedStages: []*model.PipelineStage{
				{
					Id:    "plugin-1-stage-1",
					Name:  "plugin-1-stage-1",
					Index: 0,
				},
				{
					Id:       "plugin-1-stage-2",
					Name:     "plugin-1-stage-2",
					Index:    1,
					Requires: []string{"plugin-1-stage-1"},
				},
				{
					Id:       "plugin-2-stage-1",
					Name:     "plugin-2-stage-1",
					Index:    2,
					Requires: []string{"plugin-1-stage-2"},
				},
				{
					Id:       "plugin-1-stage-3",
					Name:     "plugin-1-stage-3",
					Index:    3,
					Requires: []string{"plugin-2-stage-1"},
				},
				{
					Id:       "plugin-1-rollback",
					Index:    0,
					Name:     "plugin-1-rollback",
					Rollback: true,
				},
				{
					Id:       "plugin-2-rollback",
					Index:    2,
					Name:     "plugin-2-rollback",
					Rollback: true,
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			planner := &planner{
				pluginRegistry: tc.pluginRegistry,
			}
			stages, err := planner.buildPipelineSyncStages(context.TODO(), tc.cfg)
			require.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.expectedStages, stages)
		})
	}
}

func TestPlanner_BuildPlan(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name           string
		isFirstDeploy  bool
		plugins        []pluginapi.PluginClient
		pluginRegistry registry.PluginRegistry
		cfg            *config.GenericApplicationSpec
		deployment     *model.Deployment
		wantErr        bool
		expectedOutput *plannerOutput
	}{
		{
			name:          "quick sync strategy triggered by web console",
			isFirstDeploy: false,
			pluginRegistry: func() registry.PluginRegistry {
				pr, err := registry.NewPluginRegistry(context.TODO(), []registry.Plugin{
					{
						Name: "plugin-1",
						Cli: &fakePlugin{
							quickStages: []*model.PipelineStage{
								{
									Id:      "plugin-1-stage-1",
									Name:    "plugin-1-stage-1",
									Visible: true,
								},
							},
						},
					},
				})
				require.NoError(t, err)

				return pr
			}(),
			cfg: &config.GenericApplicationSpec{
				Planner: config.DeploymentPlanner{
					AutoRollback: pointerBool(true),
				},
				Plugins: []string{"plugin-1"},
			},
			deployment: &model.Deployment{
				Trigger: &model.DeploymentTrigger{
					SyncStrategy:    model.SyncStrategy_QUICK_SYNC,
					StrategySummary: "Triggered by web console",
				},
			},
			wantErr: false,
			expectedOutput: &plannerOutput{
				SyncStrategy: model.SyncStrategy_QUICK_SYNC,
				Summary:      "Triggered by web console",
				Stages: []*model.PipelineStage{
					{
						Id:      "plugin-1-stage-1",
						Name:    "plugin-1-stage-1",
						Visible: true,
					},
				},
				Versions: []*model.ArtifactVersion{
					{
						Kind:    model.ArtifactVersion_UNKNOWN,
						Version: versionUnknown,
					},
				},
			},
		},
		{
			name:          "pipeline sync strategy triggered by web console",
			isFirstDeploy: false,
			pluginRegistry: func() registry.PluginRegistry {
				pr, err := registry.NewPluginRegistry(context.TODO(), []registry.Plugin{
					{
						Name: "plugin-1",
						Cli: &fakePlugin{
							pipelineStages: []*model.PipelineStage{
								{
									Id:      "plugin-1-stage-1",
									Index:   0,
									Name:    "plugin-1-stage-1",
									Visible: true,
								},
							},
						},
					},
				})
				require.NoError(t, err)

				return pr
			}(),
			cfg: &config.GenericApplicationSpec{
				Planner: config.DeploymentPlanner{
					AutoRollback: pointerBool(true),
				},
				Pipeline: &config.DeploymentPipeline{
					Stages: []config.PipelineStage{
						{
							ID:   "plugin-1-stage-1",
							Name: "plugin-1-stage-1",
						},
					},
				},
			},
			deployment: &model.Deployment{
				Trigger: &model.DeploymentTrigger{
					SyncStrategy:    model.SyncStrategy_PIPELINE,
					StrategySummary: "Triggered by web console",
				},
			},
			wantErr: false,
			expectedOutput: &plannerOutput{
				SyncStrategy: model.SyncStrategy_PIPELINE,
				Summary:      "Triggered by web console",
				Stages: []*model.PipelineStage{
					{
						Id:      "plugin-1-stage-1",
						Name:    "plugin-1-stage-1",
						Index:   0,
						Visible: true,
					},
				},
				Versions: []*model.ArtifactVersion{
					{
						Kind:    model.ArtifactVersion_UNKNOWN,
						Version: versionUnknown,
					},
				},
			},
		},
		{
			name:          "quick sync due to no pipeline configured",
			isFirstDeploy: false,
			pluginRegistry: func() registry.PluginRegistry {
				pr, err := registry.NewPluginRegistry(context.TODO(), []registry.Plugin{
					{
						Name: "plugin-1",
						Cli: &fakePlugin{
							quickStages: []*model.PipelineStage{
								{
									Id:      "plugin-1-stage-1",
									Name:    "plugin-1-stage-1",
									Visible: true,
								},
							},
						},
					},
				})
				require.NoError(t, err)

				return pr
			}(),
			cfg: &config.GenericApplicationSpec{
				Planner: config.DeploymentPlanner{
					AutoRollback: pointerBool(true),
				},
				Plugins: []string{"plugin-1"},
			},
			deployment: &model.Deployment{
				Trigger: &model.DeploymentTrigger{},
			},
			wantErr: false,
			expectedOutput: &plannerOutput{
				SyncStrategy: model.SyncStrategy_QUICK_SYNC,
				Summary:      "Quick sync due to the pipeline was not configured",
				Stages: []*model.PipelineStage{
					{
						Id:      "plugin-1-stage-1",
						Name:    "plugin-1-stage-1",
						Visible: true,
					},
				},
				Versions: []*model.ArtifactVersion{
					{
						Kind:    model.ArtifactVersion_UNKNOWN,
						Version: versionUnknown,
					},
				},
			},
		},
		{
			name:          "pipeline sync due to alwaysUsePipeline",
			isFirstDeploy: false,
			pluginRegistry: func() registry.PluginRegistry {
				pr, err := registry.NewPluginRegistry(context.TODO(), []registry.Plugin{
					{
						Name: "plugin-1",
						Cli: &fakePlugin{
							pipelineStages: []*model.PipelineStage{
								{
									Id:      "plugin-1-stage-1",
									Index:   0,
									Name:    "plugin-1-stage-1",
									Visible: true,
								},
							},
						},
					},
				})
				require.NoError(t, err)

				return pr
			}(),
			cfg: &config.GenericApplicationSpec{
				Planner: config.DeploymentPlanner{
					AlwaysUsePipeline: true,
					AutoRollback:      pointerBool(true),
				},
				Pipeline: &config.DeploymentPipeline{
					Stages: []config.PipelineStage{
						{
							ID:   "plugin-1-stage-1",
							Name: "plugin-1-stage-1",
						},
					},
				},
			},
			deployment: &model.Deployment{
				Trigger: &model.DeploymentTrigger{},
			},
			wantErr: false,
			expectedOutput: &plannerOutput{
				SyncStrategy: model.SyncStrategy_PIPELINE,
				Summary:      "Sync with the specified pipeline (alwaysUsePipeline was set)",
				Stages: []*model.PipelineStage{
					{
						Id:      "plugin-1-stage-1",
						Name:    "plugin-1-stage-1",
						Index:   0,
						Visible: true,
					},
				},
				Versions: []*model.ArtifactVersion{
					{
						Kind:    model.ArtifactVersion_UNKNOWN,
						Version: versionUnknown,
					},
				},
			},
		},
		{
			name:          "quick sync due to first deployment",
			isFirstDeploy: true,
			pluginRegistry: func() registry.PluginRegistry {
				pr, err := registry.NewPluginRegistry(context.TODO(), []registry.Plugin{
					{
						Name: "plugin-1",
						Cli: &fakePlugin{
							quickStages: []*model.PipelineStage{
								{
									Id:      "plugin-1-stage-1",
									Name:    "plugin-1-stage-1",
									Visible: true,
								},
							},
						},
					},
				})
				require.NoError(t, err)

				return pr
			}(),
			cfg: &config.GenericApplicationSpec{
				Planner: config.DeploymentPlanner{
					AutoRollback: pointerBool(true),
				},
				Pipeline: &config.DeploymentPipeline{
					Stages: []config.PipelineStage{
						{
							ID:   "plugin-1-stage-1",
							Name: "plugin-1-stage-1",
						},
					},
				},
			},
			deployment: &model.Deployment{
				Trigger: &model.DeploymentTrigger{},
			},
			wantErr: false,
			expectedOutput: &plannerOutput{
				SyncStrategy: model.SyncStrategy_QUICK_SYNC,
				Summary:      "Quick sync, it seems this is the first deployment of the application",
				Stages: []*model.PipelineStage{
					{
						Id:      "plugin-1-stage-1",
						Name:    "plugin-1-stage-1",
						Visible: true,
					},
				},
				Versions: []*model.ArtifactVersion{
					{
						Kind:    model.ArtifactVersion_UNKNOWN,
						Version: versionUnknown,
					},
				},
			},
		},
		{
			name:          "pipeline sync determined by plugin",
			isFirstDeploy: false,
			plugins: []pluginapi.PluginClient{
				&fakePlugin{
					syncStrategy: &deployment.DetermineStrategyResponse{
						SyncStrategy: model.SyncStrategy_PIPELINE,
						Summary:      "determined by plugin",
					},
					pipelineStages: []*model.PipelineStage{
						{
							Id:      "plugin-1-stage-1",
							Name:    "plugin-1-stage-1",
							Visible: true,
						},
					},
					quickStages: []*model.PipelineStage{
						{
							Id:      "plugin-1-quick-stage-1",
							Visible: true,
						},
					},
				},
			},
			pluginRegistry: func() registry.PluginRegistry {
				pr, err := registry.NewPluginRegistry(context.TODO(), []registry.Plugin{
					{
						Name: "plugin-1",
						Cli: &fakePlugin{
							syncStrategy: &deployment.DetermineStrategyResponse{
								SyncStrategy: model.SyncStrategy_PIPELINE,
								Summary:      "determined by plugin",
							},
							pipelineStages: []*model.PipelineStage{
								{
									Id:      "plugin-1-stage-1",
									Index:   0,
									Name:    "plugin-1-stage-1",
									Visible: true,
								},
							},
							quickStages: []*model.PipelineStage{
								{
									Id:      "plugin-1-quick-stage-1",
									Visible: true,
								},
							},
						},
					},
				})
				require.NoError(t, err)

				return pr
			}(),
			cfg: &config.GenericApplicationSpec{
				Planner: config.DeploymentPlanner{
					AutoRollback: pointerBool(true),
				},
				Pipeline: &config.DeploymentPipeline{
					Stages: []config.PipelineStage{
						{
							ID:   "plugin-1-stage-1",
							Name: "plugin-1-stage-1",
						},
					},
				},
			},
			deployment: &model.Deployment{
				Trigger: &model.DeploymentTrigger{},
			},
			wantErr: false,
			expectedOutput: &plannerOutput{
				SyncStrategy: model.SyncStrategy_PIPELINE,
				Summary:      "determined by plugin",
				Stages: []*model.PipelineStage{
					{
						Id:      "plugin-1-stage-1",
						Name:    "plugin-1-stage-1",
						Index:   0,
						Visible: true,
					},
				},
				Versions: []*model.ArtifactVersion{
					{
						Kind:    model.ArtifactVersion_UNKNOWN,
						Version: versionUnknown,
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			planner := &planner{
				pluginRegistry:               tc.pluginRegistry,
				deployment:                   tc.deployment,
				lastSuccessfulCommitHash:     "",
				lastSuccessfulConfigFilename: "",
				workingDir:                   "",
				apiClient:                    nil,
				gitClient:                    nil,
				notifier:                     nil,
				logger:                       zap.NewNop(),
				nowFunc:                      time.Now,
			}

			if !tc.isFirstDeploy {
				planner.lastSuccessfulCommitHash = "123"
			}

			runningDS := &deployment.DeploymentSource{}

			type genericConfig struct {
				Kind       config.Kind `json:"kind"`
				APIVersion string      `json:"apiVersion,omitempty"`
				Spec       any         `json:"spec"`
			}

			jsonBytes, err := json.Marshal(genericConfig{
				Kind:       config.KindApplication,
				APIVersion: config.VersionV1Beta1,
				Spec:       tc.cfg,
			})

			require.NoError(t, err)
			targetDS := &deployment.DeploymentSource{
				ApplicationConfig: jsonBytes,
			}
			out, err := planner.buildPlan(context.TODO(), runningDS, targetDS)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.expectedOutput, out)
		})
	}
}
