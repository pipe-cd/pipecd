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

package deployment

import (
	"testing"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func Test_buildQuickSyncPipeline(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		rollback bool
		expected []sdk.QuickSyncStage
	}{
		{
			name:     "without rollback",
			rollback: false,
			expected: []sdk.QuickSyncStage{
				{
					Name:               StageK8sSync,
					Description:        StageDescriptionK8sSync,
					Rollback:           false,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
			},
		},
		{
			name:     "with rollback",
			rollback: true,
			expected: []sdk.QuickSyncStage{
				{
					Name:               StageK8sSync,
					Description:        StageDescriptionK8sSync,
					Rollback:           false,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
				{
					Name:               StageK8sRollback,
					Description:        StageDescriptionK8sRollback,
					Rollback:           true,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := buildQuickSyncPipeline(tt.rollback)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_buildPipelineStages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		stages       []sdk.StageConfig
		autoRollback bool
		expected     []sdk.PipelineStage
	}{
		{
			name: "without auto rollback",
			stages: []sdk.StageConfig{
				{
					Name:  "K8S_CANARY_ROLLOUT",
					Index: 0,
				},
				{
					Name:  "K8S_PRIMARY_ROLLOUT",
					Index: 1,
				},
			},
			autoRollback: false,
			expected: []sdk.PipelineStage{
				{
					Name:               "K8S_CANARY_ROLLOUT",
					Index:              0,
					Rollback:           false,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
				{
					Name:               "K8S_PRIMARY_ROLLOUT",
					Index:              1,
					Rollback:           false,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
			},
		},
		{
			name: "with auto rollback",
			stages: []sdk.StageConfig{
				{
					Name:  "K8S_CANARY_ROLLOUT",
					Index: 0,
				},
				{
					Name:  "K8S_PRIMARY_ROLLOUT",
					Index: 1,
				},
			},
			autoRollback: true,
			expected: []sdk.PipelineStage{
				{
					Name:               "K8S_CANARY_ROLLOUT",
					Index:              0,
					Rollback:           false,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
				{
					Name:               "K8S_PRIMARY_ROLLOUT",
					Index:              1,
					Rollback:           false,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
				{
					Name:               "K8S_ROLLBACK",
					Index:              0,
					Rollback:           true,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
			},
		},
		{
			name: "with traffic routing stage with all primary",
			stages: []sdk.StageConfig{
				{
					Name:  "K8S_TRAFFIC_ROUTING",
					Index: 0,
					Config: []byte(`{
						"all": "primary"
					}`),
				},
			},
			autoRollback: false,
			expected: []sdk.PipelineStage{
				{
					Name:     "K8S_TRAFFIC_ROUTING",
					Index:    0,
					Rollback: false,
					Metadata: map[string]string{
						sdk.MetadataKeyStageDisplay: "Primary: 100%, Canary: 0%, Baseline: 0%",
					},
					AvailableOperation: sdk.ManualOperationNone,
				},
			},
		},
		{
			name: "with traffic routing stage with all canary",
			stages: []sdk.StageConfig{
				{
					Name:  "K8S_TRAFFIC_ROUTING",
					Index: 0,
					Config: []byte(`{
						"all": "canary"
					}`),
				},
			},
			autoRollback: false,
			expected: []sdk.PipelineStage{
				{
					Name:     "K8S_TRAFFIC_ROUTING",
					Index:    0,
					Rollback: false,
					Metadata: map[string]string{
						sdk.MetadataKeyStageDisplay: "Primary: 0%, Canary: 100%, Baseline: 0%",
					},
					AvailableOperation: sdk.ManualOperationNone,
				},
			},
		},
		{
			name: "with traffic routing stage with all baseline",
			stages: []sdk.StageConfig{
				{
					Name:  "K8S_TRAFFIC_ROUTING",
					Index: 0,
					Config: []byte(`{
						"all": "baseline"
					}`),
				},
			},
			autoRollback: false,
			expected: []sdk.PipelineStage{
				{
					Name:     "K8S_TRAFFIC_ROUTING",
					Index:    0,
					Rollback: false,
					Metadata: map[string]string{
						sdk.MetadataKeyStageDisplay: "Primary: 0%, Canary: 0%, Baseline: 100%",
					},
					AvailableOperation: sdk.ManualOperationNone,
				},
			},
		},
		{
			name: "with traffic routing stage with primary and canary",
			stages: []sdk.StageConfig{
				{
					Name:  "K8S_TRAFFIC_ROUTING",
					Index: 0,
					Config: []byte(`{
						"primary": 50,
						"canary": 50
					}`),
				},
			},
			autoRollback: false,
			expected: []sdk.PipelineStage{
				{
					Name:     "K8S_TRAFFIC_ROUTING",
					Index:    0,
					Rollback: false,
					Metadata: map[string]string{
						sdk.MetadataKeyStageDisplay: "Primary: 50%, Canary: 50%, Baseline: 0%",
					},
					AvailableOperation: sdk.ManualOperationNone,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual, err := buildPipelineStages(&sdk.BuildPipelineSyncStagesInput{
				Request: sdk.BuildPipelineSyncStagesRequest{
					Stages:   tt.stages,
					Rollback: tt.autoRollback,
				},
				Logger: zaptest.NewLogger(t),
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
