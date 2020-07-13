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

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipe/pkg/model"
)

func (e *Executor) ensureSync(ctx context.Context, commitHash string, manifestsLoader func(ctx context.Context) ([]provider.Manifest, error)) model.StageStatus {
	// Load the manifests at the specified commit.
	e.LogPersister.AppendInfo(fmt.Sprintf("Loading manifests at commit %s for handling", commitHash))
	manifests, err := manifestsLoader(ctx)
	if err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Failed while loading running manifests (%v)", err))
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.AppendSuccess(fmt.Sprintf("Successfully loaded %d manifests", len(manifests)))

	// Generate the manifests for applying.
	applyManifests := e.generateSyncManifests(e.config.Input.Namespace, commitHash, manifests)

	// Start applying all manifests to add or update running resources.
	e.LogPersister.AppendInfo(fmt.Sprintf("Start applying %d manifests", len(applyManifests)))
	for _, m := range applyManifests {
		if err := e.provider.ApplyManifest(ctx, m); err != nil {
			e.LogPersister.AppendError(fmt.Sprintf("Failed to apply manifest: %s (%v)", m.Key.ReadableString(), err))
			return model.StageStatus_STAGE_FAILURE
		}
		e.LogPersister.AppendSuccess(fmt.Sprintf("- applied manifest: %s", m.Key.ReadableString()))
	}
	e.LogPersister.AppendSuccess(fmt.Sprintf("Successfully applied %d manifests", len(applyManifests)))

	// TODO: Wait for all applied manifests to be ready.
	e.LogPersister.AppendInfo("Waiting for the applied manifests to be ready")

	// TODO: Find and remove the running resources that are not defined in Git.
	e.LogPersister.AppendInfo("Start finding and removing all running resources but not in Git")

	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) generateSyncManifests(namespace, commitHash string, manifests []provider.Manifest) []provider.Manifest {
	out := make([]provider.Manifest, 0, len(manifests))

	for _, manifest := range manifests {
		// Because the loaded maninests are read-only
		// so we duplicate them to avoid updating the shared manifests data in cache.
		m := manifest.Duplicate(manifest.Key.Name)
		if namespace != "" {
			m.SetNamespace(namespace)
			m.Key.Namespace = namespace
		}
		// Add predefined annotation to the manifest.
		m.AddAnnotations(e.builtinAnnotations(m, primaryVariant, commitHash))
		out = append(out, manifest)
	}
	return out
}
