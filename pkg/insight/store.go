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

package insight

import (
	"context"
)

// Storage for applications.
type (
	ApplicationGetter interface {
		GetApplications(ctx context.Context, projectID string) (*ProjectApplicationData, error)
	}

	ApplicationStore interface {
		ApplicationGetter
		PutApplications(ctx context.Context, projectID string, as *ProjectApplicationData) error
	}
)

// Storage for completed deployments.
type (
	CompletedDeploymentGetter interface {
		GetMilestone(ctx context.Context) (*Milestone, error)
		ListCompletedDeployments(ctx context.Context, projectID string, from, to int64) ([]*DeploymentData, error)
	}

	CompletedDeploymentStore interface {
		CompletedDeploymentGetter
		PutMilestone(ctx context.Context, m *Milestone) error
		PutCompletedDeployments(ctx context.Context, projectID string, ds []*DeploymentData) error
	}
)

type (
	Getter interface {
		ApplicationGetter
		CompletedDeploymentGetter
	}

	Store interface {
		ApplicationStore
		CompletedDeploymentStore
	}
)

type Milestone struct {
	// Mark that our collector has accumulated all deployments
	// that was completed before this value.
	DeploymentCompletedAtMilestone int64 `json:"deployment_completed_at_milestone"`
}
