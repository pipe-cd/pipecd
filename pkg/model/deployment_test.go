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

package model

import (
	"sort"
	"testing"
	"time"

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

func TestCanUpdateDeploymentStatus(t *testing.T) {
	tests := []struct {
		name string
		cur  DeploymentStatus
		next DeploymentStatus
		want bool
	}{
		{
			name: "can update from PENDING to PLANNED",
			cur:  DeploymentStatus_DEPLOYMENT_PENDING,
			next: DeploymentStatus_DEPLOYMENT_PLANNED,
			want: true,
		},
		{
			name: "cannot update from PLANNED to PENDING",
			cur:  DeploymentStatus_DEPLOYMENT_PLANNED,
			next: DeploymentStatus_DEPLOYMENT_PENDING,
			want: false,
		},
		{
			name: "can update from RUNNING to ROLLING_BACK",
			cur:  DeploymentStatus_DEPLOYMENT_RUNNING,
			next: DeploymentStatus_DEPLOYMENT_ROLLING_BACK,
			want: true,
		},
		{
			name: "cannot update from ROLLING_BACK to RUNNING",
			cur:  DeploymentStatus_DEPLOYMENT_ROLLING_BACK,
			next: DeploymentStatus_DEPLOYMENT_RUNNING,
			want: false,
		},
		{
			name: "can update from ROLLING_BACK to SUCCESS",
			cur:  DeploymentStatus_DEPLOYMENT_ROLLING_BACK,
			next: DeploymentStatus_DEPLOYMENT_SUCCESS,
			want: true,
		},
		{
			name: "cannot update from SUCCESS to ROLLING_BACK",
			cur:  DeploymentStatus_DEPLOYMENT_SUCCESS,
			next: DeploymentStatus_DEPLOYMENT_ROLLING_BACK,
			want: false,
		},
		{
			name: "can update from ROLLING_BACK to FAILURE",
			cur:  DeploymentStatus_DEPLOYMENT_ROLLING_BACK,
			next: DeploymentStatus_DEPLOYMENT_FAILURE,
			want: true,
		},
		{
			name: "cannot update from FAILURE to ROLLING_BACK",
			cur:  DeploymentStatus_DEPLOYMENT_FAILURE,
			next: DeploymentStatus_DEPLOYMENT_ROLLING_BACK,
			want: false,
		},
		{
			name: "can update from ROLLING_BACK to CANCELLED",
			cur:  DeploymentStatus_DEPLOYMENT_ROLLING_BACK,
			next: DeploymentStatus_DEPLOYMENT_CANCELLED,
			want: true,
		},
		{
			name: "cannot update from CANCELLED to ROLLING_BACK",
			cur:  DeploymentStatus_DEPLOYMENT_CANCELLED,
			next: DeploymentStatus_DEPLOYMENT_ROLLING_BACK,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanUpdateDeploymentStatus(tt.cur, tt.next)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCanUpdateStageStatus(t *testing.T) {
	tests := []struct {
		name string
		cur  StageStatus
		next StageStatus
		want bool
	}{
		{
			name: "can update from NOT_STARTED_YET to RUNNING",
			cur:  StageStatus_STAGE_NOT_STARTED_YET,
			next: StageStatus_STAGE_RUNNING,
			want: true,
		},
		{
			name: "can update from RUNNING to SUCCESS",
			cur:  StageStatus_STAGE_RUNNING,
			next: StageStatus_STAGE_SUCCESS,
			want: true,
		},
		{
			name: "can update from RUNNING to FAILURE",
			cur:  StageStatus_STAGE_RUNNING,
			next: StageStatus_STAGE_FAILURE,
			want: true,
		},
		{
			name: "can update from RUNNING to CANCELLED",
			cur:  StageStatus_STAGE_RUNNING,
			next: StageStatus_STAGE_CANCELLED,
			want: true,
		},
		{
			name: "cannot update from SUCCESS to RUNNING",
			cur:  StageStatus_STAGE_SUCCESS,
			next: StageStatus_STAGE_RUNNING,
			want: false,
		},
		{
			name: "cannot update from FAILURE to RUNNING",
			cur:  StageStatus_STAGE_FAILURE,
			next: StageStatus_STAGE_RUNNING,
			want: false,
		},
		{
			name: "cannot update from CANCELLED to RUNNING",
			cur:  StageStatus_STAGE_CANCELLED,
			next: StageStatus_STAGE_RUNNING,
			want: false,
		},
		{
			name: "cannot update from SUCCESS to FAILURE",
			cur:  StageStatus_STAGE_SUCCESS,
			next: StageStatus_STAGE_FAILURE,
			want: false,
		},
		{
			name: "cannot update from SUCCESS to CANCELLED",
			cur:  StageStatus_STAGE_SUCCESS,
			next: StageStatus_STAGE_CANCELLED,
			want: false,
		},
		{
			name: "cannot update from FAILURE to SUCCESS",
			cur:  StageStatus_STAGE_FAILURE,
			next: StageStatus_STAGE_SUCCESS,
			want: false,
		},
		{
			name: "cannot update from FAILURE to CANCELLED",
			cur:  StageStatus_STAGE_FAILURE,
			next: StageStatus_STAGE_CANCELLED,
			want: false,
		},
		{
			name: "cannot update from CANCELLED to SUCCESS",
			cur:  StageStatus_STAGE_CANCELLED,
			next: StageStatus_STAGE_SUCCESS,
			want: false,
		},
		{
			name: "cannot update from CANCELLED to FAILURE",
			cur:  StageStatus_STAGE_CANCELLED,
			next: StageStatus_STAGE_FAILURE,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanUpdateStageStatus(tt.cur, tt.next)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTriggeredBy(t *testing.T) {
	tests := []struct {
		name    string
		trigger DeploymentTrigger
		want    string
	}{
		{
			name: "returns commander name if set",
			trigger: DeploymentTrigger{
				Commander: "Alice",
				Commit: &Commit{
					Author: "Bob",
				},
			},
			want: "Alice",
		},
		{
			name: "returns commit author name if commander not set",
			trigger: DeploymentTrigger{
				Commander: "",
				Commit: &Commit{
					Author: "Bob",
				},
			},
			want: "Bob",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Deployment{
				Trigger: &tt.trigger,
			}
			got := d.TriggeredBy()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTriggerBefore(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name  string
		d     Deployment
		other Deployment
		want  bool
	}{
		{
			name: "returns true if d trigger is before other trigger",
			d: Deployment{
				Trigger: &DeploymentTrigger{
					Commit: &Commit{
						CreatedAt: now.Add(-time.Hour).Unix(),
					},
					Timestamp: now.Add(-time.Minute).Unix(),
				},
			},
			other: Deployment{
				Trigger: &DeploymentTrigger{
					Commit: &Commit{
						CreatedAt: now.Unix(),
					},
					Timestamp: now.Unix(),
				},
			},
			want: true,
		},
		{
			name: "returns false if d trigger is after other trigger",
			d: Deployment{
				Trigger: &DeploymentTrigger{
					Commit: &Commit{
						CreatedAt: now.Unix(),
					},
					Timestamp: now.Unix(),
				},
			},
			other: Deployment{
				Trigger: &DeploymentTrigger{
					Commit: &Commit{
						CreatedAt: now.Add(-time.Hour).Unix(),
					},
					Timestamp: now.Add(-time.Minute).Unix(),
				},
			},
			want: false,
		},
		{
			name: "returns true if d trigger is same as other trigger",
			d: Deployment{
				Trigger: &DeploymentTrigger{
					Commit: &Commit{
						CreatedAt: now.Unix(),
					},
					Timestamp: now.Unix(),
				},
			},
			other: Deployment{
				Trigger: &DeploymentTrigger{
					Commit: &Commit{
						CreatedAt: now.Unix(),
					},
					Timestamp: now.Unix(),
				},
			},
			want: true,
		},
		{
			name: "returns false if d trigger is same as other trigger",
			d: Deployment{
				Trigger: &DeploymentTrigger{
					Commit: &Commit{
						CreatedAt: now.Unix(),
					},
					Timestamp: now.Add(time.Minute).Unix(),
				},
			},
			other: Deployment{
				Trigger: &DeploymentTrigger{
					Commit: &Commit{
						CreatedAt: now.Unix(),
					},
					Timestamp: now.Unix(),
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.TriggerBefore(&tt.other)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFindRollbackStage(t *testing.T) {
	tests := []struct {
		name           string
		stages         []*PipelineStage
		wantStage      *PipelineStage
		wantStageFound bool
	}{
		{
			name: "found",
			stages: []*PipelineStage{
				{Name: StageK8sSync.String()},
				{Name: StageRollback.String()},
			},
			wantStage:      &PipelineStage{Name: StageRollback.String()},
			wantStageFound: true,
		},
		{
			name: "not found",
			stages: []*PipelineStage{
				{Name: StageK8sSync.String()},
			},
			wantStage:      nil,
			wantStageFound: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			d := &Deployment{
				Stages: tt.stages,
			}
			stage, found := d.FindRollbackStage()
			assert.Equal(t, tt.wantStage, stage)
			assert.Equal(t, tt.wantStageFound, found)
		})
	}
}

func TestFindRollbackStags(t *testing.T) {
	tests := []struct {
		name           string
		stages         []*PipelineStage
		wantStages     []*PipelineStage
		wantStageFound bool
	}{
		{
			name: "found",
			stages: []*PipelineStage{
				{Name: StageK8sSync.String()},
				{Name: StageRollback.String()},
				{Name: StageScriptRunRollback.String()},
			},
			wantStages: []*PipelineStage{
				{Name: StageRollback.String()},
				{Name: StageScriptRunRollback.String()},
			},
			wantStageFound: true,
		},
		{
			name: "found based on rollback field",
			stages: []*PipelineStage{
				{Name: "Some-plugin-stage-name", Rollback: true},
			},
			wantStages: []*PipelineStage{
				{Name: "Some-plugin-stage-name", Rollback: true},
			},
			wantStageFound: true,
		},
		{
			name: "not found",
			stages: []*PipelineStage{
				{Name: StageK8sSync.String()},
			},
			wantStages:     []*PipelineStage{},
			wantStageFound: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			d := &Deployment{
				Stages: tt.stages,
			}
			stages, found := d.FindRollbackStages()
			assert.Equal(t, tt.wantStages, stages)
			assert.Equal(t, tt.wantStageFound, found)
		})
	}
}

func TestSortPipelineStagesByIndex(t *testing.T) {
	stages := []*PipelineStage{
		{Index: 2},
		{Index: 1},
		{Index: 4},
		{Index: 3},
	}
	sort.Sort(PipelineStages(stages))
	assert.Equal(t, []*PipelineStage{
		{Index: 1},
		{Index: 2},
		{Index: 3},
		{Index: 4},
	}, stages)
}
