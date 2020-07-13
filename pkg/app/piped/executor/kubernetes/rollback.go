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

package kubernetes

import (
	"context"
	"fmt"
	"strings"

	"github.com/pipe-cd/pipe/pkg/model"
)

func (e *Executor) ensureRollback(ctx context.Context) model.StageStatus {
	// 1. Revert PRIMARY resources.
	e.LogPersister.AppendInfo(fmt.Sprintf("Start checking to ensure that PRIMARY resources match to commit %s", e.Deployment.RunningCommitHash))
	if err := e.rollbackPrimary(ctx); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}

	// 2. Ensure that all traffics are routed to the PRIMARY variant.
	e.LogPersister.AppendInfo("Start checking to ensure that all traffics being routed to the PRIMARY variant")
	if err := e.rollbackTraffic(ctx); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}

	var errs []error

	// 3. Delete all resources of CANARY variant.
	e.LogPersister.AppendInfo("Start checking to ensure that the CANARY variant should be removed")
	if value, ok := e.MetadataStore.Get(addedCanaryResourcesMetadataKey); ok {
		resources := strings.Split(value, ",")
		if err := e.removeCanaryResources(ctx, resources); err != nil {
			errs = append(errs, err)
		}
	}

	// 4. Delete all resources of BASELINE variant.
	e.LogPersister.AppendInfo("Start checking to ensure that the BASELINE variant should be removed")
	if value, ok := e.MetadataStore.Get(addedBaselineResourcesMetadataKey); ok {
		resources := strings.Split(value, ",")
		if err := e.removeBaselineResources(ctx, resources); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return model.StageStatus_STAGE_FAILURE
	}
	return model.StageStatus_STAGE_SUCCESS
}
