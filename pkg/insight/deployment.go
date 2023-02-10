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

package insight

import (
	"github.com/pipe-cd/pipecd/pkg/model"
)

type DeploymentData struct {
	Id                string            `json:"id"`
	AppID             string            `json:"app_id"`
	Labels            map[string]string `json:"labels"`
	StartedAt         int64             `json:"started_at"`
	CompletedAt       int64             `json:"completed_at"`
	CompleteStatus    string            `json:"complete_status"`
	RollbackStartedAt int64             `json:"rollback_started_at"`
}

func BuildDeploymentData(d *model.Deployment) DeploymentData {
	var rollbackStartedAt int64
	if s, ok := d.FindRollbackStage(); ok {
		rollbackStartedAt = s.CreatedAt
	}

	return DeploymentData{
		Id:                d.Id,
		AppID:             d.ApplicationId,
		Labels:            d.Labels,
		StartedAt:         d.CreatedAt,
		CompletedAt:       d.CompletedAt,
		RollbackStartedAt: rollbackStartedAt,
		CompleteStatus:    d.Status.String(),
	}
}

func (d *DeploymentData) ContainLabels(labels map[string]string) bool {
	if len(labels) == 0 {
		return true
	}

	for k, v := range labels {
		if dv := d.Labels[k]; dv != v {
			return false
		}
	}
	return true
}
