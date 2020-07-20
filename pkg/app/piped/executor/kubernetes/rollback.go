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
	"strings"

	"github.com/pipe-cd/pipe/pkg/model"
)

func (e *Executor) ensureRollback(ctx context.Context) model.StageStatus {
	commitHash := e.Deployment.RunningCommitHash

	// Firstly, we reapply all manifests at running commit
	// to revert PRIMARY resources and TRAFFIC ROUTING resources.

	// Load the manifests at the specified commit.
	e.LogPersister.AppendInfof("Loading manifests at running commit %s for handling", commitHash)
	manifests, err := e.loadRunningManifests(ctx)
	if err != nil {
		e.LogPersister.AppendErrorf("Failed while loading running manifests (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.AppendSuccessf("Successfully loaded %d manifests", len(manifests))

	// Because the loaded maninests are read-only
	// we duplicate them to avoid updating the shared manifests data in cache.
	manifests = duplicateManifests(manifests, "")

	// Add builtin annotations for tracking application live state.
	e.addBuiltinAnnontations(manifests, primaryVariant, commitHash)

	// Start applying all manifests to add or update running resources.
	if err := e.applyManifests(ctx, manifests); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}

	var errs []error

	// Next we delete all resources of CANARY variant.
	e.LogPersister.AppendInfo("Start checking to ensure that the CANARY variant should be removed")
	if value, ok := e.MetadataStore.Get(addedCanaryResourcesMetadataKey); ok {
		resources := strings.Split(value, ",")
		if err := e.removeCanaryResources(ctx, resources); err != nil {
			errs = append(errs, err)
		}
	}

	// Then delete all resources of BASELINE variant.
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
