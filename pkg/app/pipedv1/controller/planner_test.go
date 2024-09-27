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

	"github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
	pluginapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
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
		plugins        []pluginapi.PluginClient
		cfg            *config.GenericApplicationSpec
		wantErr        bool
		expectedStages []*model.PipelineStage
	}{
		{
			name: "only one plugin",
			plugins: []pluginapi.PluginClient{
				&fakePlugin{
					quickStages: []*model.PipelineStage{
						{
							Id: "plugin-1-stage-1",
						},
					},
					rollbackStages: []*model.PipelineStage{
						{
							Id:       "plugin-1-rollback",
							Rollback: true,
						},
					},
				},
			},
			cfg: &config.GenericApplicationSpec{
				Planner: config.DeploymentPlanner{
					AutoRollback: pointerBool(true),
				},
			},
			wantErr: false,
			expectedStages: []*model.PipelineStage{
				{
					Id: "plugin-1-stage-1",
				},
				{
					Id:       "plugin-1-rollback",
					Rollback: true,
				},
			},
		},
		{
			name: "multi plugins",
			plugins: []pluginapi.PluginClient{
				&fakePlugin{
					quickStages: []*model.PipelineStage{
						{
							Id:    "plugin-1-stage-1",
							Index: 0,
						},
					},
					rollbackStages: []*model.PipelineStage{
						{
							Id:       "plugin-1-rollback",
							Index:    0,
							Rollback: true,
						},
					},
				},
				&fakePlugin{
					quickStages: []*model.PipelineStage{
						{
							Id:    "plugin-2-stage-1",
							Index: 1,
						},
					},
					rollbackStages: []*model.PipelineStage{
						{
							Id:       "plugin-2-rollback",
							Index:    1,
							Rollback: true,
						},
					},
				},
			},
			cfg: &config.GenericApplicationSpec{
				Planner: config.DeploymentPlanner{
					AutoRollback: pointerBool(true),
				},
			},
			wantErr: false,
			expectedStages: []*model.PipelineStage{
				{
					Id:    "plugin-1-stage-1",
					Index: 0,
				},
				{
					Id:    "plugin-2-stage-1",
					Index: 1,
				},
				{
					Id:       "plugin-1-rollback",
					Index:    0,
					Rollback: true,
				},
				{
					Id:       "plugin-2-rollback",
					Index:    1,
					Rollback: true,
				},
			},
		},
		{
			name: "multi plugins - no rollback",
			plugins: []pluginapi.PluginClient{
				&fakePlugin{
					quickStages: []*model.PipelineStage{
						{
							Id:    "plugin-1-stage-1",
							Index: 0,
						},
					},
					rollbackStages: []*model.PipelineStage{
						{
							Id:       "plugin-1-rollback",
							Rollback: true,
						},
					},
				},
				&fakePlugin{
					quickStages: []*model.PipelineStage{
						{
							Id:    "plugin-2-stage-1",
							Index: 1,
						},
					},
					rollbackStages: []*model.PipelineStage{
						{
							Id:       "plugin-2-rollback",
							Rollback: true,
						},
					},
				},
			},
			cfg: &config.GenericApplicationSpec{
				Planner: config.DeploymentPlanner{
					AutoRollback: pointerBool(false),
				},
			},
			wantErr: false,
			expectedStages: []*model.PipelineStage{
				{
					Id:    "plugin-1-stage-1",
					Index: 0,
				},
				{
					Id:    "plugin-2-stage-1",
					Index: 1,
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			planner := &planner{
				plugins: tc.plugins,
			}
			stages, err := planner.buildQuickSyncStages(context.TODO(), tc.cfg)
			require.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.expectedStages, stages)
		})
	}
}

func TestBuildPipelineSyncStages(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name           string
		plugins        []pluginapi.PluginClient
		cfg            *config.GenericApplicationSpec
		wantErr        bool
		expectedStages []*model.PipelineStage
	}{
		{
			name: "only one plugin",
			plugins: []pluginapi.PluginClient{
				&fakePlugin{
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
			plugins: []pluginapi.PluginClient{
				&fakePlugin{
					pipelineStages: []*model.PipelineStage{
						{
							Id:   "plugin-1-stage-1",
							Name: "plugin-1-stage-1",
						},
						{
							Id:   "plugin-1-stage-2",
							Name: "plugin-1-stage-2",
						},
						{
							Id:   "plugin-1-stage-3",
							Name: "plugin-1-stage-3",
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
				&fakePlugin{
					pipelineStages: []*model.PipelineStage{
						{
							Id:   "plugin-2-stage-1",
							Name: "plugin-2-stage-1",
						},
						{
							Id:   "plugin-2-stage-2",
							Name: "plugin-2-stage-2",
						},
					},
				},
			},
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
			plugins: []pluginapi.PluginClient{
				&fakePlugin{
					pipelineStages: []*model.PipelineStage{
						{
							Id:   "plugin-1-stage-1",
							Name: "plugin-1-stage-1",
						},
						{
							Id:   "plugin-1-stage-2",
							Name: "plugin-1-stage-2",
						},
						{
							Id:   "plugin-1-stage-3",
							Name: "plugin-1-stage-3",
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
				&fakePlugin{
					pipelineStages: []*model.PipelineStage{
						{
							Id:   "plugin-2-stage-1",
							Name: "plugin-2-stage-1",
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
			stageBasedPluginMap := make(map[string]pluginapi.PluginClient)
			for _, p := range tc.plugins {
				stages, _ := p.FetchDefinedStages(context.TODO(), &deployment.FetchDefinedStagesRequest{})
				for _, s := range stages.Stages {
					stageBasedPluginMap[s] = p
				}
			}
			planner := &planner{
				plugins:              tc.plugins,
				stageBasedPluginsMap: stageBasedPluginMap,
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
		cfg            *config.GenericApplicationSpec
		deployment     *model.Deployment
		wantErr        bool
		expectedOutput *plannerOutput
	}{
		{
			name:          "quick sync strategy triggered by web console",
			isFirstDeploy: false,
			plugins: []pluginapi.PluginClient{
				&fakePlugin{
					quickStages: []*model.PipelineStage{
						{
							Id:      "plugin-1-stage-1",
							Visible: true,
						},
					},
				},
			},
			cfg: &config.GenericApplicationSpec{
				Planner: config.DeploymentPlanner{
					AutoRollback: pointerBool(true),
				},
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
			plugins: []pluginapi.PluginClient{
				&fakePlugin{
					pipelineStages: []*model.PipelineStage{
						{
							Id:      "plugin-1-stage-1",
							Name:    "plugin-1-stage-1",
							Visible: true,
						},
					},
				},
			},
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
			plugins: []pluginapi.PluginClient{
				&fakePlugin{
					quickStages: []*model.PipelineStage{
						{
							Id:      "plugin-1-stage-1",
							Visible: true,
						},
					},
				},
			},
			cfg: &config.GenericApplicationSpec{
				Planner: config.DeploymentPlanner{
					AutoRollback: pointerBool(true),
				},
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
			plugins: []pluginapi.PluginClient{
				&fakePlugin{
					pipelineStages: []*model.PipelineStage{
						{
							Id:      "plugin-1-stage-1",
							Name:    "plugin-1-stage-1",
							Visible: true,
						},
					},
				},
			},
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
			plugins: []pluginapi.PluginClient{
				&fakePlugin{
					quickStages: []*model.PipelineStage{
						{
							Id:      "plugin-1-stage-1",
							Visible: true,
						},
					},
				},
			},
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
			stageBasedPluginMap := make(map[string]pluginapi.PluginClient)
			for _, p := range tc.plugins {
				stages, _ := p.FetchDefinedStages(context.TODO(), &deployment.FetchDefinedStagesRequest{})
				for _, s := range stages.Stages {
					stageBasedPluginMap[s] = p
				}
			}
			planner := &planner{
				plugins:                      tc.plugins,
				stageBasedPluginsMap:         stageBasedPluginMap,
				deployment:                   tc.deployment,
				lastSuccessfulCommitHash:     "",
				lastSuccessfulConfigFilename: "",
				workingDir:                   "",
				apiClient:                    nil,
				gitClient:                    nil,
				notifier:                     nil,
				logger:                       zap.NewNop(),
				nowFunc:                      func() time.Time { return time.Now() },
			}

			if !tc.isFirstDeploy {
				planner.lastSuccessfulCommitHash = "123"
			}

			runningDS := &model.DeploymentSource{}

			jsonBytes, err := json.Marshal(tc.cfg)
			require.NoError(t, err)
			targetDS := &model.DeploymentSource{
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
