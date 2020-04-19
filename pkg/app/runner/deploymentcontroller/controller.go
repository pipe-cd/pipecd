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

// Package deploymentcontroller provides a runner component
// that managing all of the Deployment CRDs.
// This manages a pool of DeploymentExecutors.
// Whenever a new Deployment CRD is created, this created a new DeploymentExecutor
// for that DeploymentCRD to handle the deployment.
// The DeploymentExecutor will update the deployment state back to its Deployment CRD.
package deploymentcontroller

import (
	"context"
	"time"
)

type DeploymentController struct {
	gracePeriod time.Duration
}

// NewController creates a new instance for DeploymentController.
func NewController(gracePeriod time.Duration) *DeploymentController {
	return &DeploymentController{
		gracePeriod: gracePeriod,
	}
}

// Run starts running DeploymentController until the specified context
// has done. This also waits for its cleaning up before returning.'
func (t *DeploymentController) Run(ctx context.Context) error {
	return nil
}
