// Copyright 2022 The PipeCD Authors.
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

type Milestone struct {
	// Mark that our collector has handled all
	// deployment that was created before this value
	DeploymentCreatedAtMilestone int64 `json:"deployment_created_at_milestone"`
	// Mark that our collector has handled all deployment
	// that was completed before this value. This will be
	// used while calculating CHANGE_FAILURE_RATE.
	DeploymentCompletedAtMilestone int64 `json:"deployment_completed_at_milestone"`
}
