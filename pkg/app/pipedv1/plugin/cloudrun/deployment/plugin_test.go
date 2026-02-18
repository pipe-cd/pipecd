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

package deployment

import (
	"testing"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/stretchr/testify/assert"
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
					Name:               StageCloudRunSync,
					Description:        StageCloudRunSyncDescription,
					Rollback:           false,
					Metadata:           map[string]string{},
					AvailableOperation: sdk.ManualOperationNone,
				},
			},
		},
		{
			name:     "with rollback",
			rollback: true,
			expected: []sdk.QuickSyncStage{
				{
					Name:               StageCloudRunSync,
					Description:        StageCloudRunSyncDescription,
					Rollback:           false,
					Metadata:           map[string]string{},
					AvailableOperation: sdk.ManualOperationNone,
				},
				{
					Name:               StageRollback,
					Description:        StageRollbackDescription,
					Rollback:           true,
					Metadata:           map[string]string{},
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
					Name:  "CLOUDRUN_PROMOTE",
					Index: 0,
				},
				{
					Name:  "CLOUDRUN_PROMOTE",
					Index: 1,
				},
			},
			autoRollback: false,
			expected: []sdk.PipelineStage{
				{
					Name:               StageCloudRunPromote,
					Index:              0,
					Rollback:           false,
					Metadata:           map[string]string{},
					AvailableOperation: sdk.ManualOperationNone,
				},
				{
					Name:               StageCloudRunPromote,
					Index:              1,
					Rollback:           false,
					Metadata:           map[string]string{},
					AvailableOperation: sdk.ManualOperationNone,
				},
			},
		},
		{
			name: "with auto rollback",
			stages: []sdk.StageConfig{
				{
					Name:  "CLOUDRUN_PROMOTE",
					Index: 0,
				},
				{
					Name:  "CLOUDRUN_PROMOTE",
					Index: 1,
				},
			},
			autoRollback: true,
			expected: []sdk.PipelineStage{
				{
					Name:               StageCloudRunPromote,
					Index:              0,
					Rollback:           false,
					Metadata:           map[string]string{},
					AvailableOperation: sdk.ManualOperationNone,
				},
				{
					Name:               StageCloudRunPromote,
					Index:              1,
					Rollback:           false,
					Metadata:           map[string]string{},
					AvailableOperation: sdk.ManualOperationNone,
				},
				{
					Name:               StageRollback,
					Index:              0,
					Rollback:           true,
					Metadata:           map[string]string{},
					AvailableOperation: sdk.ManualOperationNone,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := buildPipelineStages(tt.stages, tt.autoRollback)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestImageVersionExtraction(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		image    string
		expected string
	}{
		{
			name:     "GCR Standard Tag",
			image:    "gcr.io/pipecd/webapp:v1.2.3",
			expected: "v1.2.3",
		},
		{
			name:     "Image Digest (SHA)",
			image:    "gcr.io/pipecd/webapp@sha256:45b23dee08af",
			expected: "sha256:45b23dee08af",
		},
		{
			name:     "Registry with Port",
			image:    "localhost:5000/pipecd/webapp:v1.0.0",
			expected: "v1.0.0",
		},
		{
			name:     "Untagged Image",
			image:    "gcr.io/pipecd/webapp",
			expected: "latest",
		},
		{
			name:     "Empty String",
			image:    "",
			expected: "unknown",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := ImageVersionExtraction(tc.image)
			assert.Equal(t, tc.expected, actual)
		})
	}
}