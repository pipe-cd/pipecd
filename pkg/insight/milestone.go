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

package insight

import "github.com/pipe-cd/pipe/pkg/model"

type Milestone struct {
	DeploymentCreatedAtMilestone   int64 // Mark that our collector has handled all deployment that was created before this value
	DeploymentCompletedAtMilestone int64 // Mark that our collector has handled all deployment that was completed before this value. This will be used while calculating CHANGE_FAILURE_RATE.
}

func (i *Milestone) updateMilestone(m int64, kind model.InsightMetricsKind) {
	switch kind {
	case model.InsightMetricsKind_CHANGE_FAILURE_RATE:
		i.DeploymentCompletedAtMilestone = m
	default:
		i.DeploymentCreatedAtMilestone = m
	}
}
