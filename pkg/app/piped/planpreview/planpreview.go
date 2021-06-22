// Copyright 2021 The PipeCD Authors.
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

package planpreview

import (
	"context"

	"github.com/pipe-cd/pipe/pkg/model"
)

type Planner interface {
	Plan(ctx context.Context, repoID, branch, commit string) ([]*model.ApplicationPlanPreviewResult, error)
}

type planner struct {
}

func (p *planner) Plan(ctx context.Context, repoID, branch, commit string) ([]*model.ApplicationPlanPreviewResult, error) {
	// TODO: Implement Plan functionality.

	// 1. List all applications placing in that repository.
	// 2. Fetch the source code at the specified branch commit.
	// 3. Determine the list of applications that will be triggered.
	//    - Based on the changed files between 2 commits: head commit and mostRecentlyTriggeredCommit
	// 4. For each application:
	//    4.1. Start a planner to check what/why strategy will be used
	//    4.2. Check what resources should be added, deleted and modified
	//         - Terraform app: used terraform plan command
	//         - Kubernetes app: calculate the diff of resources at head commit and mostRecentlySuccessfulCommit

	return nil, nil
}
