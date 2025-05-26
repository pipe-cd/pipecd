// Copyright 2025 The PipeCD Authors.
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

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/provider"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

func ensureVariantSelectorInWorkload(m provider.Manifest, variantLabel, variant string) error {
	variantMap := map[string]string{
		variantLabel: variant,
	}
	if err := m.AddStringMapValues(variantMap, "spec", "selector", "matchLabels"); err != nil {
		return err
	}
	return m.AddStringMapValues(variantMap, "spec", "template", "metadata", "labels")
}

// addVariantLabelsAndAnnotations adds the variant label and annotation to the given manifests.
func addVariantLabelsAndAnnotations(m []provider.Manifest, variantLabel, variant string) {
	for _, m := range m {
		m.AddLabels(map[string]string{
			variantLabel: variant,
		})
		m.AddAnnotations(map[string]string{
			variantLabel: variant,
		})
	}
}

// deleteResources deletes the given resources.
// It returns the number of deleted resources.
func deleteResources(ctx context.Context, lp sdk.StageLogPersister, applier *provider.Applier, keys []provider.ResourceKey) int {
	var deletedCount int

	for _, k := range keys {
		if err := applier.Delete(ctx, k); err != nil {
			if errors.Is(err, provider.ErrNotFound) {
				lp.Infof("Specified resource does not exist, so skip deleting the resource: %s (%v)", k.ReadableString(), err)
				continue
			}
			lp.Errorf("Failed while deleting resource %s (%v)", k.ReadableString(), err)
			continue // continue to delete other resources
		}
		deletedCount++
		lp.Successf("- deleted resource: %s", k.ReadableString())
	}

	return deletedCount
}
