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

import "path/filepath"

// IsCompleted checks whether the deployment is at a completion state.
func (d *Deployment) IsCompleted() bool {
	if d.Status.String() == "" {
		return false
	}

	switch d.Status {
	case DeploymentStatus_DEPLOYMENT_SUCCESS:
		return true
	case DeploymentStatus_DEPLOYMENT_FAILURE:
		return true
	case DeploymentStatus_DEPLOYMENT_CANCELLED:
		return true
	}
	return false
}

// CanUpdateStatus checks whether the deployment can transit to the given status.
func (d *Deployment) CanUpdateStatus(status DeploymentStatus) bool {
	switch status {
	case DeploymentStatus_DEPLOYMENT_TRIGGERED:
		return d.Status <= DeploymentStatus_DEPLOYMENT_TRIGGERED
	case DeploymentStatus_DEPLOYMENT_PENDING:
		return d.Status <= DeploymentStatus_DEPLOYMENT_PENDING
	case DeploymentStatus_DEPLOYMENT_RUNNING:
		return d.Status <= DeploymentStatus_DEPLOYMENT_RUNNING
	case DeploymentStatus_DEPLOYMENT_SUCCESS:
		return d.Status <= DeploymentStatus_DEPLOYMENT_RUNNING
	case DeploymentStatus_DEPLOYMENT_FAILURE:
		return d.Status <= DeploymentStatus_DEPLOYMENT_RUNNING
	case DeploymentStatus_DEPLOYMENT_CANCELLED:
		return d.Status <= DeploymentStatus_DEPLOYMENT_RUNNING
	}
	return false
}

// GetDeploymentConfigFilePath returns the path to deployment configuration directory.
func (d *Deployment) GetDeploymentConfigFilePath(filename string) string {
	if path := d.GitPath.ConfigPath; path != "" {
		return path
	}
	return filepath.Join(d.GitPath.Path, filename)
}
