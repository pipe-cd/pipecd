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
	"time"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

const (
	primaryVariant = "primary"
)

func (e *Executor) ensurePrimaryRollout(ctx context.Context) model.StageStatus {
	var (
		commitHash = e.Deployment.Trigger.Commit.Hash
		options    = e.StageConfig.K8sPrimaryRolloutStageOptions
	)
	if options == nil {
		e.LogPersister.AppendErrorf("Malformed configuration for stage %s", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}

	// Load the manifests at the triggered commit.
	e.LogPersister.AppendInfof("Loading manifests at commit %s for handling", commitHash)
	manifests, err := e.loadManifests(ctx)
	if err != nil {
		e.LogPersister.AppendErrorf("Failed while loading manifests (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.AppendSuccessf("Successfully loaded %d manifests", len(manifests))

	// Generate the manifests for applying.
	e.LogPersister.AppendInfo("Start generating manifests for PRIMARY variant")
	applyManifests, err := e.generatePrimaryManifests(commitHash, manifests)
	if err != nil {
		e.LogPersister.AppendErrorf("Unable to generate manifests for PRIMARY variant (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.AppendSuccessf("Successfully generated %d manifests for PRIMARY variant", len(applyManifests))

	// Start applying all manifests to add or update running resources.
	e.LogPersister.AppendInfof("Start applying %d primary resources", len(applyManifests))
	for _, m := range applyManifests {
		if err = e.provider.ApplyManifest(ctx, m); err != nil {
			e.LogPersister.AppendErrorf("Failed to apply manifest: %s (%v)", m.Key.ReadableString(), err)
			return model.StageStatus_STAGE_FAILURE
		}
		e.LogPersister.AppendSuccessf("- applied manifest: %s", m.Key.ReadableString())
	}
	e.LogPersister.AppendSuccessf("Successfully applied %d primary resources", len(applyManifests))

	if !options.Prune {
		e.LogPersister.AppendInfo("Resource GC was skipped because sync.prune was not configured")
		return model.StageStatus_STAGE_SUCCESS
	}

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

	// Find the running resources that are not defined in Git.
	e.LogPersister.AppendInfo("Start finding all running PRIMARY resources but no longer defined in Git")
	runningManifests, err := e.loadRunningManifests(ctx)
	if err != nil {
		e.LogPersister.AppendErrorf("Failed while loading running manifests (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	var (
		applyKeys  = make(map[provider.ResourceKey]struct{}, len(applyManifests))
		removeKeys = make([]provider.ResourceKey, 0)
	)
	for _, m := range applyManifests {
		applyKeys[m.Key] = struct{}{}
	}
	for _, m := range runningManifests {
		key := m.Key
		if _, ok := applyKeys[key]; ok {
			continue
		}
		if key.Namespace == "" {
			key.Namespace = e.config.Input.Namespace
		}
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

func (e *Executor) generatePrimaryManifests(commitHash string, manifests []provider.Manifest) ([]provider.Manifest, error) {
	var (
		serviceName      string
		generateService  bool
		primaryManifests = make([]provider.Manifest, 0, len(manifests)+1)
		suffix           = primaryVariant
	)

	// Apply the specified configuration if they are present.
	if sc := e.config.PrimaryVariant; sc != nil {
		var ok bool
		if sc.Suffix != "" {
			suffix = sc.Suffix
		}
		generateService = sc.Service.Create

		_, serviceName, ok = config.ParseVariantResourceReference(sc.Service.Reference)
		if !ok {
			return nil, fmt.Errorf("malformed service reference: %s", sc.Service.Reference)
		}
	}

	for _, m := range manifests {
		// Because the loaded maninests are read-only
		// so we duplicate them to avoid updating the shared manifests data in cache.
		primaryManifests = append(primaryManifests, m.Duplicate(m.Key.Name))
	}

	// Find service manifests and duplicate them for PRIMARY variant.
	if generateService {
		services := findManifests(provider.KindService, serviceName, manifests)
		if len(services) == 0 {
			return nil, fmt.Errorf("unable to find any service for name=%q", serviceName)
		}

		// Because the loaded maninests are read-only
		// so we duplicate them to avoid updating the shared manifests data in cache.
		duplicates := make([]provider.Manifest, 0, len(services))
		for _, m := range services {
			duplicates = append(duplicates, m.Duplicate(m.Key.Name))
		}

		generatedServices, err := generateServiceManifests(duplicates, primaryVariant, suffix)
		if err != nil {
			return nil, err
		}
		primaryManifests = append(primaryManifests, generatedServices...)
	}

	// TODO: Find out traffic-routing manfiests and keep them as the previously configured routing.

	// Add predefined annotations to the generated manifests.
	for _, m := range primaryManifests {
		m.AddAnnotations(e.builtinAnnotations(m, primaryVariant, commitHash))
	}
	return primaryManifests, nil
}
