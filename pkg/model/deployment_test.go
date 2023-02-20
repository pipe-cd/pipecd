// Copyright 2023 The PipeCD Authors.
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

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeployment_ContainTags(t *testing.T) {
	testcases := []struct {
		name       string
		deployment *Deployment
		labels     map[string]string
		want       bool
	}{
		{
			name:       "all given tags aren't contained",
			deployment: &Deployment{Labels: map[string]string{"key1": "value1"}},
			labels: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			want: false,
		},
		{
			name: "a label is contained",
			deployment: &Deployment{Labels: map[string]string{
				"key1": "value1",
				"key2": "value2",
			}},
			labels: map[string]string{
				"key1": "value1",
			},
			want: true,
		},
		{
			name: "all tags are contained",
			deployment: &Deployment{Labels: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			}},
			labels: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			want: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.deployment.ContainLabels(tc.labels)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDeployment_StageMap(t *testing.T) {
	testcases := []struct {
		name       string
		deployment *Deployment
		want       map[string]*PipelineStage
	}{
		{
			name: "ok",
			deployment: &Deployment{
				Stages: []*PipelineStage{
					{
						Id: "stage1",
					},
					{
						Id: "stage2",
					},
				},
			},
			want: map[string]*PipelineStage{
				"stage1": {
					Id: "stage1",
				},
				"stage2": {
					Id: "stage2",
				},
			},
		},
		{
			name:       "no stages",
			deployment: &Deployment{},
			want:       map[string]*PipelineStage{},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.deployment.StageMap()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStageStatus_IsCompleted(t *testing.T) {
	testcases := []struct {
		name   string
		status StageStatus
		want   bool
	}{
		{
			name:   "running",
			status: StageStatus_STAGE_RUNNING,
			want:   false,
		},
		{
			name:   "success",
			status: StageStatus_STAGE_SUCCESS,
			want:   true,
		},
		{
			name:   "failure",
			status: StageStatus_STAGE_FAILURE,
			want:   true,
		},
		{
			name:   "cancelled",
			status: StageStatus_STAGE_CANCELLED,
			want:   true,
		},
		{
			name:   "skipped",
			status: StageStatus_STAGE_SKIPPED,
			want:   true,
		},
		{
			name:   "exited",
			status: StageStatus_STAGE_EXITED,
			want:   true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.status.IsCompleted()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDeploymentStatus_IsCompleted(t *testing.T) {
	testcases := []struct {
		name   string
		status DeploymentStatus
		want   bool
	}{
		{
			name:   "pending",
			status: DeploymentStatus_DEPLOYMENT_PENDING,
			want:   false,
		},
		{
			name:   "planned",
			status: DeploymentStatus_DEPLOYMENT_PLANNED,
			want:   false,
		},
		{
			name:   "running",
			status: DeploymentStatus_DEPLOYMENT_RUNNING,
			want:   false,
		},
		{
			name:   "rolling back",
			status: DeploymentStatus_DEPLOYMENT_ROLLING_BACK,
			want:   false,
		},
		{
			name:   "success",
			status: DeploymentStatus_DEPLOYMENT_SUCCESS,
			want:   true,
		},
		{
			name:   "failure",
			status: DeploymentStatus_DEPLOYMENT_FAILURE,
			want:   true,
		},
		{
			name:   "cancelled",
			status: DeploymentStatus_DEPLOYMENT_CANCELLED,
			want:   true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.status.IsCompleted()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDeployment_Stage(t *testing.T) {
	testcases := []struct {
		name       string
		id         string
		deployment *Deployment
		want       *PipelineStage
		exists     bool
	}{
		{
			name: "ok",
			deployment: &Deployment{
				Stages: []*PipelineStage{
					{
						Id: "id",
					},
				},
			},
			id:     "id",
			want:   &PipelineStage{Id: "id"},
			exists: true,
		},
		{
			name: "not found",
			deployment: &Deployment{
				Stages: []*PipelineStage{
					{
						Id: "id",
					},
				},
			},
			id:     "foo",
			want:   nil,
			exists: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := tc.deployment.Stage(tc.id)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.exists, ok)
		})
	}
}
