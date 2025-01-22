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

package deployment

import (
	"context"
	"errors"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
)

func applyManifests(ctx context.Context, applier applier, manifests []provider.Manifest, namespace string, lp logpersister.StageLogPersister) error {
	if namespace == "" {
		lp.Infof("Start applying %d manifests", len(manifests))
	} else {
		lp.Infof("Start applying %d manifests to %q namespace", len(manifests), namespace)
	}

	for _, m := range manifests {
		// The force annotation has higher priority, so we need to check the annotation in the following order:
		// 1. force-sync-by-replace
		// 2. sync-by-replace
		// 3. others
		if annotation := m.GetAnnotations()[provider.LabelForceSyncReplace]; annotation == provider.UseReplaceEnabled {
			// Always try to replace first and create if it fails due to resource not found error.
			// This is because we cannot know whether resource already exists before executing command.
			err := applier.ForceReplaceManifest(ctx, m)
			if errors.Is(err, provider.ErrNotFound) {
				lp.Infof("Specified resource does not exist, so create the resource: %s (%w)", m.Key().ReadableString(), err)
				err = applier.CreateManifest(ctx, m)
			}
			if err != nil {
				lp.Errorf("Failed to forcefully replace or create manifest: %s (%w)", m.Key().ReadableString(), err)
				return err
			}
			lp.Successf("- forcefully replaced or created manifest: %s", m.Key().ReadableString())
			continue
		}

		if annotation := m.GetAnnotations()[provider.LabelSyncReplace]; annotation == provider.UseReplaceEnabled {
			// Always try to replace first and create if it fails due to resource not found error.
			// This is because we cannot know whether resource already exists before executing command.
			err := applier.ReplaceManifest(ctx, m)
			if errors.Is(err, provider.ErrNotFound) {
				lp.Infof("Specified resource does not exist, so create the resource: %s (%w)", m.Key().ReadableString(), err)
				err = applier.CreateManifest(ctx, m)
			}
			if err != nil {
				lp.Errorf("Failed to replace or create manifest: %s (%w)", m.Key().ReadableString(), err)
				return err
			}
			lp.Successf("- replaced or created manifest: %s", m.Key().ReadableString())
			continue
		}

		if err := applier.ApplyManifest(ctx, m); err != nil {
			lp.Errorf("Failed to apply manifest: %s (%w)", m.Key().ReadableString(), err)
			return err
		}
		lp.Successf("- applied manifest: %s", m.Key().ReadableString())
		continue

	}
	lp.Successf("Successfully applied %d manifests", len(manifests))
	return nil
}
