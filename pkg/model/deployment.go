// Copyright 2020 The PipeCD Authors.
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

// IsCompletedDeployment checks whether the deployment is at a completion state.
func IsCompletedDeployment(status DeploymentStatus) bool {
	if status.String() == "" {
		return false
	}

	switch status {
	case DeploymentStatus_DEPLOYMENT_SUCCESS:
		return true
	case DeploymentStatus_DEPLOYMENT_FAILURE:
		return true
	case DeploymentStatus_DEPLOYMENT_CANCELLED:
		return true
	}
	return false
}

// IsCompletedStage checks whether the stage is at a completion state.
func IsCompletedStage(status StageStatus) bool {
	if status.String() == "" {
		return false
	}

	switch status {
	case StageStatus_STAGE_SUCCESS:
		return true
	case StageStatus_STAGE_FAILURE:
		return true
	case StageStatus_STAGE_CANCELLED:
		return true
	}
	return false
}

// CanUpdateDeploymentStatus checks whether the deployment can transit to the given status.
func CanUpdateDeploymentStatus(cur, next DeploymentStatus) bool {
	switch next {
	case DeploymentStatus_DEPLOYMENT_TRIGGERED:
		return cur <= DeploymentStatus_DEPLOYMENT_TRIGGERED
	case DeploymentStatus_DEPLOYMENT_PENDING:
		return cur <= DeploymentStatus_DEPLOYMENT_PENDING
	case DeploymentStatus_DEPLOYMENT_RUNNING:
		return cur <= DeploymentStatus_DEPLOYMENT_RUNNING
	case DeploymentStatus_DEPLOYMENT_SUCCESS:
		return cur <= DeploymentStatus_DEPLOYMENT_RUNNING
	case DeploymentStatus_DEPLOYMENT_FAILURE:
		return cur <= DeploymentStatus_DEPLOYMENT_RUNNING
	case DeploymentStatus_DEPLOYMENT_CANCELLED:
		return cur <= DeploymentStatus_DEPLOYMENT_RUNNING
	}
	return false
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
