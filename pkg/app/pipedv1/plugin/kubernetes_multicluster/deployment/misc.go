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
	"fmt"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/provider"
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

// deleteVariantResources finds and deletes all live resources labeled with the given variant.
// It deletes in order: Services → Workloads → Others → Cluster-scoped resources.
func deleteVariantResources(ctx context.Context, lp sdk.StageLogPersister, kubectl *provider.Kubectl, kubeConfig string, applier *provider.Applier, applicationID, variantLabel, variant string) error {
	namespacedLiveResources, clusterScopedLiveResources, err := provider.GetLiveResources(ctx, kubectl, kubeConfig, applicationID, fmt.Sprintf("%s=%s", variantLabel, variant))
	if err != nil {
		return err
	}

	services := make([]provider.ResourceKey, 0, len(namespacedLiveResources))
	workloads := make([]provider.ResourceKey, 0, len(namespacedLiveResources))
	others := make([]provider.ResourceKey, 0, len(namespacedLiveResources))
	clusterScoped := make([]provider.ResourceKey, 0, len(clusterScopedLiveResources))

	for _, r := range namespacedLiveResources {
		switch {
		case r.IsService():
			services = append(services, r.Key())
		case r.IsWorkload():
			workloads = append(workloads, r.Key())
		default:
			others = append(others, r.Key())
		}
	}

	for _, r := range clusterScopedLiveResources {
		clusterScoped = append(clusterScoped, r.Key())
	}

	var deletedCount int
	deletedCount += deleteResources(ctx, lp, applier, services)
	deletedCount += deleteResources(ctx, lp, applier, workloads)
	deletedCount += deleteResources(ctx, lp, applier, others)
	deletedCount += deleteResources(ctx, lp, applier, clusterScoped)
	lp.Successf("Successfully deleted %d resources", deletedCount)

	return nil
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
