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

	"github.com/pipe-cd/pipe/pkg/model"
)

func (e *Executor) ensureRollback(ctx context.Context) model.StageStatus {
	// 1. Revert workloads of PRIMARY variant.
	if state := e.rollbackPrimary(ctx); state != model.StageStatus_STAGE_SUCCESS {
		return state
	}

	// 2. Ensure that all traffics are routed to the PRIMARY variant.
	if state := e.rollbackTraffic(ctx); state != model.StageStatus_STAGE_SUCCESS {
		return state
	}

	// 3. Delete workloads of CANARY variant.
	if state := e.ensureCanaryClean(ctx); state != model.StageStatus_STAGE_SUCCESS {
		return state
	}

	// 4. Delete worloads of BASELINE variant.
	if state := e.ensureBaselineClean(ctx); state != model.StageStatus_STAGE_SUCCESS {
		return state
	}

	return model.StageStatus_STAGE_SUCCESS
}
