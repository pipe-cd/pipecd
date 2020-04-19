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

// Package deploymenttrigger provides a runner component
// that detects a list of application should be synced
// and then trigger their deployments by applying new
// Deployment-CRDs for them.
// Until V1, we detect based on the new merged commit and its changes.
// But in the next versions, we also want to enable the ability to detect
// based on the diff between the repo state (desired state) and cluster state (actual state).
package deploymenttrigger

import "time"

type DeploymentTrigger struct {
}

func NewTrigger() *DeploymentTrigger {
	return &DeploymentTrigger{}
}

func (t *DeploymentTrigger) Run() error {
	return nil
}

func (t *DeploymentTrigger) Stop(timeout time.Duration) error {
	return nil
}
