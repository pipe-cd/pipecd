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
	"fmt"

	"google.golang.org/protobuf/proto"
)

const (
	MetadataKeyDeploymentNotification = "DeploymentNotification"
)

var notCompletedDeploymentStatuses = []DeploymentStatus{
	DeploymentStatus_DEPLOYMENT_PENDING,
	DeploymentStatus_DEPLOYMENT_PLANNED,
	DeploymentStatus_DEPLOYMENT_RUNNING,
	DeploymentStatus_DEPLOYMENT_ROLLING_BACK,
}

// IsCompleted checks whether the deployment is at a completion state.
func (d DeploymentStatus) IsCompleted() bool {
	return d == DeploymentStatus_DEPLOYMENT_SUCCESS ||
		d == DeploymentStatus_DEPLOYMENT_FAILURE ||
		d == DeploymentStatus_DEPLOYMENT_CANCELLED
}

// IsCompleted checks whether the stage is at a completion state.
func (s StageStatus) IsCompleted() bool {
	return s == StageStatus_STAGE_SUCCESS ||
		s == StageStatus_STAGE_FAILURE ||
		s == StageStatus_STAGE_CANCELLED ||
		s == StageStatus_STAGE_SKIPPED ||
		s == StageStatus_STAGE_EXITED
}

// CanUpdateDeploymentStatus checks whether the deployment can transit to the given status.
func CanUpdateDeploymentStatus(cur, next DeploymentStatus) bool {
	switch next {
	case DeploymentStatus_DEPLOYMENT_PENDING:
		return cur <= DeploymentStatus_DEPLOYMENT_PENDING
	case DeploymentStatus_DEPLOYMENT_PLANNED:
		return cur <= DeploymentStatus_DEPLOYMENT_PLANNED
	case DeploymentStatus_DEPLOYMENT_RUNNING:
		return cur <= DeploymentStatus_DEPLOYMENT_RUNNING
	case DeploymentStatus_DEPLOYMENT_ROLLING_BACK:
		return cur <= DeploymentStatus_DEPLOYMENT_ROLLING_BACK
	case DeploymentStatus_DEPLOYMENT_SUCCESS:
		return cur <= DeploymentStatus_DEPLOYMENT_ROLLING_BACK
	case DeploymentStatus_DEPLOYMENT_FAILURE:
		return cur <= DeploymentStatus_DEPLOYMENT_ROLLING_BACK
	case DeploymentStatus_DEPLOYMENT_CANCELLED:
		return cur <= DeploymentStatus_DEPLOYMENT_ROLLING_BACK
	}
	return false
}

// GetNotCompletedDeploymentStatuses returns DeploymentStatus slice that is not completed.
func GetNotCompletedDeploymentStatuses() []DeploymentStatus {
	return notCompletedDeploymentStatuses
}

// CanUpdateStageStatus checks whether the stage can transit to the given status.
func CanUpdateStageStatus(cur, next StageStatus) bool {
	switch next {
	case StageStatus_STAGE_NOT_STARTED_YET:
		return cur <= StageStatus_STAGE_NOT_STARTED_YET
	case StageStatus_STAGE_RUNNING:
		return cur <= StageStatus_STAGE_RUNNING
	case StageStatus_STAGE_SUCCESS:
		return cur <= StageStatus_STAGE_RUNNING
	case StageStatus_STAGE_FAILURE:
		return cur <= StageStatus_STAGE_RUNNING
	case StageStatus_STAGE_CANCELLED:
		return cur <= StageStatus_STAGE_RUNNING
	}
	return false
}

// StageMap returns the map of id and the stage.
func (d *Deployment) StageMap() map[string]*PipelineStage {
	stage := make(map[string]*PipelineStage, len(d.Stages))
	for _, s := range d.Stages {
		stage[s.Id] = s
	}
	return stage
}

// Stage returns the stage by stage id.
func (d *Deployment) Stage(id string) (*PipelineStage, bool) {
	for _, s := range d.Stages {
		if s.Id == id {
			return s, true
		}
	}
	return nil, false
}

// IsSkippable checks whether skippable or not.
func (p *PipelineStage) IsSkippable() bool {
	return p.Name == StageAnalysis.String()
}

// CommitHash returns the hash value of trigger commit.
func (d *Deployment) CommitHash() string {
	return d.Trigger.Commit.Hash
}

func (d *Deployment) TriggeredBy() string {
	if d.Trigger.Commander != "" {
		return d.Trigger.Commander
	}
	return d.Trigger.Commit.Author
}

// Clone returns a deep copy of the deployment.
func (d *Deployment) Clone() *Deployment {
	msg := proto.Clone(d)
	return msg.(*Deployment)
}

// TriggerBefore reports whether the deployment is triggered before the given one.
func (d *Deployment) TriggerBefore(other *Deployment) bool {
	if d.Trigger.Commit.CreatedAt > other.Trigger.Commit.CreatedAt {
		return false
	}
	if d.Trigger.Commit.CreatedAt == other.Trigger.Commit.CreatedAt && d.Trigger.Timestamp > other.Trigger.Timestamp {
		return false
	}
	return true

}

// FindRollbackStage finds the rollback stage in stage list.
func (d *Deployment) FindRollbackStage() (*PipelineStage, bool) {
	for i := len(d.Stages) - 1; i >= 0; i-- {
		if d.Stages[i].Name == StageRollback.String() || d.Stages[i].Name == StageCustomSyncRollback.String() {
			return d.Stages[i], true
		}
	}
	return nil, false
}

func (d *Deployment) FindRollbackStages() ([]*PipelineStage, bool) {
	rollbackStages := make([]*PipelineStage, 0, len(d.Stages))
	for i, stage := range d.Stages {
		if d.Stages[i].Rollback {
			rollbackStages = append(rollbackStages, stage)
		}
		// TODO: Remove this condition after all deployments have the rollback field.
		// Probally when all piped are updated to v1.
		if d.Stages[i].Name == StageRollback.String() || d.Stages[i].Name == StageScriptRunRollback.String() {
			rollbackStages = append(rollbackStages, stage)
		}
	}
	return rollbackStages, len(rollbackStages) > 0
}

// DeploymentStatusesFromStrings converts a list of strings to list of DeploymentStatus.
func DeploymentStatusesFromStrings(statuses []string) ([]DeploymentStatus, error) {
	out := make([]DeploymentStatus, 0, len(statuses))
	for _, s := range statuses {
		status, ok := DeploymentStatus_value[s]
		if !ok {
			return nil, fmt.Errorf("invalid status %s", s)
		}
		out = append(out, DeploymentStatus(status))
	}
	return out, nil
}

// DeploymentStatusStrings returns a list of available deployment statuses in string.
func DeploymentStatusStrings() []string {
	out := make([]string, 0, len(DeploymentStatus_value))
	for s := range DeploymentStatus_value {
		out = append(out, s)
	}
	return out
}

// ContainLabels checks if it has all the given labels.
func (d *Deployment) ContainLabels(labels map[string]string) bool {
	if len(d.Labels) < len(labels) {
		return false
	}

	for k, v := range labels {
		value, ok := d.Labels[k]
		if !ok {
			return false
		}
		if value != v {
			return false
		}
	}
	return true
}

// IsInChainDeployment returns true if the current deployment belongs
// to a deployment chain.
func (d *Deployment) IsInChainDeployment() bool {
	return d.DeploymentChainId != ""
}

func (d *Deployment) SetUpdatedAt(t int64) {
	d.UpdatedAt = t
}

// Implement sort.Interface for PipelineStages.
type PipelineStages []*PipelineStage

func (p PipelineStages) Len() int {
	return len(p)
}

func (p PipelineStages) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PipelineStages) Less(i, j int) bool {
	return p[i].Index < p[j].Index
}
