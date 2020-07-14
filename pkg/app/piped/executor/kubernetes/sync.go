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
	"time"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipe/pkg/model"
)

func (e *Executor) ensureSync(ctx context.Context, commitHash string, manifestsLoader func(ctx context.Context) ([]provider.Manifest, error)) model.StageStatus {
	// Load the manifests at the specified commit.
	e.LogPersister.AppendInfof("Loading manifests at commit %s for handling", commitHash)
	manifests, err := manifestsLoader(ctx)
	if err != nil {
		e.LogPersister.AppendErrorf("Failed while loading running manifests (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.AppendSuccessf("Successfully loaded %d manifests", len(manifests))

	// Generate the manifests for applying.
	applyManifests := e.generateSyncManifests(commitHash, manifests)

	// Start applying all manifests to add or update running resources.
	e.LogPersister.AppendInfof("Start applying %d manifests", len(applyManifests))
	for _, m := range applyManifests {
		if err := e.provider.ApplyManifest(ctx, m); err != nil {
			e.LogPersister.AppendErrorf("Failed to apply manifest: %s (%v)", m.Key.ReadableString(), err)
			return model.StageStatus_STAGE_FAILURE
		}
		e.LogPersister.AppendSuccessf("- applied manifest: %s", m.Key.ReadableString())
	}
	e.LogPersister.AppendSuccessf("Successfully applied %d manifests", len(applyManifests))

	// Wait for all applied manifests to be stable.
	// In theory, we don't need to wait for them to be stable before going to the next step
	// but waiting for a while reduces the number of Kubernetes changes in a short time.
	e.LogPersister.AppendInfo("Waiting for the applied manifests to be stable")
	select {
	case <-time.After(15 * time.Second):
		break
	case <-ctx.Done():
		break
	}

	// Find the running resources that are not defined in Git for removing.
	e.LogPersister.AppendInfo("Start finding all running resources but no longer defined in Git")
	liveResources, ok := e.AppLiveResourceLister.ListKubernetesResources()
	if !ok {
		e.LogPersister.AppendInfo("There is no data about live resource so no resource will be removed")
		return model.StageStatus_STAGE_SUCCESS
	}

	var (
		applyKeys  = make(map[provider.ResourceKey]struct{}, len(applyManifests))
		removeKeys = make([]provider.ResourceKey, 0)
	)
	for _, m := range applyManifests {
		key := m.Key
		key.Namespace = ""
		applyKeys[key] = struct{}{}
	}
	for _, r := range liveResources {
		key := provider.ResourceKey{
			Name:       r.Name,
			APIVersion: r.ApiVersion,
			Kind:       r.Kind,
		}
		if _, ok := applyKeys[key]; ok {
			continue
		}
		key.Namespace = r.Namespace
		removeKeys = append(removeKeys, key)
	}
	if len(removeKeys) == 0 {
		e.LogPersister.AppendInfo("There are no live resources should be removed")
		return model.StageStatus_STAGE_SUCCESS
	}

	// Start deleting all running resources that are not defined in Git.
	e.LogPersister.AppendInfof("Start deleting %d resources", len(removeKeys))
	for _, k := range removeKeys {
		if err := e.provider.Delete(ctx, k); err != nil {
			e.LogPersister.AppendErrorf("Failed to delete resource: %s (%v)", k.ReadableString(), err)
			return model.StageStatus_STAGE_FAILURE
		}
		e.LogPersister.AppendSuccessf("- deleted resource: %s", k.ReadableString())
	}
	e.LogPersister.AppendSuccessf("Successfully deleted %d resources", len(removeKeys))

	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) generateSyncManifests(commitHash string, manifests []provider.Manifest) []provider.Manifest {
	out := make([]provider.Manifest, 0, len(manifests))

	for _, manifest := range manifests {
		// Because the loaded maninests are read-only
		// so we duplicate them to avoid updating the shared manifests data in cache.
		m := manifest.Duplicate(manifest.Key.Name)
		// Add predefined annotation to the manifest.
		m.AddAnnotations(e.builtinAnnotations(m, primaryVariant, commitHash))
		out = append(out, manifest)
	}
	return out
}
