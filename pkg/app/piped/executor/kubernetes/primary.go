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

	provider "github.com/kapetaniosci/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/kapetaniosci/pipe/pkg/model"
)

func (e *Executor) ensurePrimaryUpdate(ctx context.Context) model.StageStatus {
	manifests, err := e.loadManifests(ctx)
	if err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Failed while loading manifests (%v)", err))
		return model.StageStatus_STAGE_FAILURE
	}

	if len(manifests) == 0 {
		e.LogPersister.AppendError("No kubernetes manifests to handle")
		return model.StageStatus_STAGE_FAILURE
	}

	for _, m := range manifests {
		// Add variant label to PRIMARY workload.
		// if m.Kind == "Deployment" {
		// 	if err := m.AddVariantLabel(primaryVariant); err != nil {
		// 		e.LogPersister.AppendError(fmt.Sprintf("Unabled to configure variant label for %s deployment (%v)", m.Name, err))
		// 		return model.StageStatus_STAGE_FAILURE
		// 	}
		// }

		m.AddAnnotations(map[string]string{
			provider.LabelManagedBy:          provider.ManagedByPiped,
			provider.LabelApplication:        e.Deployment.ApplicationId,
			provider.LabelVariant:            primaryVariant,
			provider.LabelOriginalAPIVersion: m.Key.APIVersion,
			provider.LabelResourceKey:        m.Key.String(),
			provider.LabelCommitHash:         e.Deployment.Trigger.Commit.Hash,
		})
	}

	e.LogPersister.AppendInfo(fmt.Sprintf("Updating %d primary resources", len(manifests)))
	if err = e.provider.ApplyManifests(ctx, manifests); err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Unabled to update primary variant (%v)", err))
		return model.StageStatus_STAGE_FAILURE
	}

	e.LogPersister.AppendSuccess(fmt.Sprintf("Successfully updated %d primary resources", len(manifests)))
	return model.StageStatus_STAGE_SUCCESS
}
