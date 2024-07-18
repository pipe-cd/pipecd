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

package kubernetes

import (
	"context"
	"time"

	provider "github.com/pipe-cd/pipecd/pkg/app/pipedv1/platformprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func (e *deployExecutor) ensureSync(ctx context.Context) model.StageStatus {
	// Load the manifests at the specified commit.
	e.LogPersister.Infof("Loading manifests at commit %s for handling", e.commit)
	manifests, err := loadManifests(
		ctx,
		e.Deployment.ApplicationId,
		e.commit,
		e.AppManifestsCache,
		e.loader,
		e.Logger,
	)
	if err != nil {
		e.LogPersister.Errorf("Failed while loading manifests (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.Successf("Successfully loaded %d manifests", len(manifests))

	// Because the loaded manifests are read-only
	// we duplicate them to avoid updating the shared manifests data in cache.
	manifests = duplicateManifests(manifests, "")

	// When addVariantLabelToSelector is true, ensure that all workloads
	// have the variant label in their selector.
	var (
		variantLabel   = e.appCfg.VariantLabel.Key
		primaryVariant = e.appCfg.VariantLabel.PrimaryValue
	)
	if e.appCfg.QuickSync.AddVariantLabelToSelector {
		workloads := findWorkloadManifests(manifests, e.appCfg.Workloads)
		for _, m := range workloads {
			if err := ensureVariantSelectorInWorkload(m, variantLabel, primaryVariant); err != nil {
				e.LogPersister.Errorf("Unable to check/set %q in selector of workload %s (%v)", variantLabel+": "+primaryVariant, m.Key.ReadableString(), err)
				return model.StageStatus_STAGE_FAILURE
			}
		}
	}

	// Add builtin annotations for tracking application live state.
	addBuiltinAnnotations(
		manifests,
		variantLabel,
		primaryVariant,
		e.commit,
		e.PipedConfig.PipedID,
		e.Deployment.ApplicationId,
	)

	// Add config-hash annotation to the workloads.
	if err := annotateConfigHash(manifests); err != nil {
		e.LogPersister.Errorf("Unable to set %q annotation into the workload manifest (%v)", provider.AnnotationConfigHash, err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Start applying all manifests to add or update running resources.
	if err := applyManifests(ctx, e.applierGetter, manifests, e.appCfg.Input.Namespace, e.LogPersister); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}

	if !e.appCfg.QuickSync.Prune {
		e.LogPersister.Info("Resource GC was skipped because sync.prune was not configured")
		return model.StageStatus_STAGE_SUCCESS
	}

	// Wait for all applied manifests to be stable.
	// In theory, we don't need to wait for them to be stable before going to the next step
	// but waiting for a while reduces the number of Kubernetes changes in a short time.
	e.LogPersister.Info("Waiting for the applied manifests to be stable")
	select {
	case <-time.After(15 * time.Second):
		break
	case <-ctx.Done():
		break
	}

	// Find the running resources that are not defined in Git for removing.
	e.LogPersister.Info("Start finding all running resources but no longer defined in Git")
	liveResources, ok := e.AppLiveResourceLister.ListKubernetesResources()
	if !ok {
		e.LogPersister.Info("There is no data about live resource so no resource will be removed")
		return model.StageStatus_STAGE_SUCCESS
	}
	e.LogPersister.Successf("Successfully loaded %d live resources", len(liveResources))
	for _, m := range liveResources {
		e.LogPersister.Successf("- loaded live resource: %s", m.Key.ReadableString())
	}

	removeKeys := findRemoveResources(manifests, liveResources)
	if len(removeKeys) == 0 {
		e.LogPersister.Info("There are no live resources should be removed")
		return model.StageStatus_STAGE_SUCCESS
	}
	e.LogPersister.Infof("Found %d live resources that are no longer defined in Git", len(removeKeys))

	// Start deleting all running resources that are not defined in Git.
	if err := deleteResources(ctx, e.applierGetter, removeKeys, e.LogPersister); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}

	return model.StageStatus_STAGE_SUCCESS
}

func findRemoveResources(manifests []provider.Manifest, liveResources []provider.Manifest) []provider.ResourceKey {
	var (
		keys       = make(map[provider.ResourceKey]struct{}, len(manifests))
		removeKeys = make([]provider.ResourceKey, 0)
	)
	for _, m := range manifests {
		key := m.Key
		key.Namespace = ""
		keys[key] = struct{}{}
	}
	for _, m := range liveResources {
		key := m.Key
		key.Namespace = ""
		if _, ok := keys[key]; ok {
			continue
		}
		key.Namespace = m.Key.Namespace
		removeKeys = append(removeKeys, key)
	}
	return removeKeys
}
