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
		e.LogPersister.Errorf("Malformed configuration for stage %s", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}

	// Load the manifests at the triggered commit.
	e.LogPersister.Infof("Loading manifests at trigered commit %s for handling", commitHash)
	manifests, err := e.loadManifests(ctx)
	if err != nil {
		e.LogPersister.Errorf("Failed while loading manifests (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.Successf("Successfully loaded %d manifests", len(manifests))

	var primaryManifests []provider.Manifest
	if config.DetermineKubernetesTrafficRoutingMethod(e.config.TrafficRouting) == config.KubernetesTrafficRoutingMethodPodSelector {
		primaryManifests = manifests
	} else {
		// Find traffic routing manifests and filter out it from primary manifests.
		trafficRoutingManifests, err := findTrafficRoutingManifests(manifests, e.config.Service.Name, e.config.TrafficRouting)
		if err != nil {
			e.LogPersister.Errorf("Failed while finding traffic routing manifest: (%v)", err)
			return model.StageStatus_STAGE_FAILURE
		}
		if len(trafficRoutingManifests) > 0 {
			primaryManifests = make([]provider.Manifest, 0, len(manifests)-1)
			for _, m := range manifests {
				if m.Key == trafficRoutingManifests[0].Key {
					continue
				}
				primaryManifests = append(primaryManifests, m)
			}
		}
	}

	// Check the variant selector in the workloads if the addVariantLabelToSelector is false.
	if !options.AddVariantLabelToSelector {
		workloads := findWorkloadManifests(primaryManifests, e.config.Workloads)
		var invalid bool
		for _, m := range workloads {
			if err := checkVariantSelectorInWorkload(m, primaryVariant); err != nil {
				invalid = true
				e.LogPersister.Errorf("Missing %q in selector of workload %s (%v)", variantLabel+": "+primaryVariant, m.Key.ReadableString(), err)
			}
		}
		if invalid {
			return model.StageStatus_STAGE_FAILURE
		}
	}

	// Generate the manifests for applying.
	e.LogPersister.Info("Start generating manifests for PRIMARY variant")
	applyManifests, err := e.generatePrimaryManifests(primaryManifests, *options)
	if err != nil {
		e.LogPersister.Errorf("Unable to generate manifests for PRIMARY variant (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.Successf("Successfully generated %d manifests for PRIMARY variant", len(applyManifests))

	// Add builtin annotations for tracking application live state.
	e.addBuiltinAnnontations(applyManifests, primaryVariant, commitHash)

	// Start applying all manifests to add or update running resources.
	e.LogPersister.Info("Start rolling out PRIMARY variant...")
	if err := e.applyManifests(ctx, applyManifests); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.Success("Successfully rolled out PRIMARY variant")

	if !options.Prune {
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

	// Find the running resources that are not defined in Git.
	e.LogPersister.Info("Start finding all running PRIMARY resources but no longer defined in Git")
	runningManifests, err := e.loadRunningManifests(ctx)
	if err != nil {
		e.LogPersister.Errorf("Failed while loading running manifests (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	removeKeys := findRemoveManifests(runningManifests, manifests, e.config.Input.Namespace)
	if len(removeKeys) == 0 {
		e.LogPersister.Info("There are no live resources should be removed")
		return model.StageStatus_STAGE_SUCCESS
	}
	e.LogPersister.Infof("Found %d live resources that are no longer defined in Git", len(removeKeys))

	// Start deleting all running resources that are not defined in Git.
	e.LogPersister.Infof("Start deleting %d resources", len(removeKeys))
	if err := e.deleteResources(ctx, removeKeys); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}

	return model.StageStatus_STAGE_SUCCESS
}

func findRemoveManifests(prevs []provider.Manifest, curs []provider.Manifest, namespace string) []provider.ResourceKey {
	var (
		keys       = make(map[provider.ResourceKey]struct{}, len(curs))
		removeKeys = make([]provider.ResourceKey, 0)
	)
	for _, m := range curs {
		keys[m.Key] = struct{}{}
	}
	for _, m := range prevs {
		key := m.Key
		if _, ok := keys[key]; ok {
			continue
		}
		if key.Namespace == "" {
			key.Namespace = namespace
		}
		removeKeys = append(removeKeys, key)
	}
	return removeKeys
}

func (e *Executor) generatePrimaryManifests(manifests []provider.Manifest, opts config.K8sPrimaryRolloutStageOptions) ([]provider.Manifest, error) {
	suffix := primaryVariant
	if opts.Suffix != "" {
		suffix = opts.Suffix
	}

	// Because the loaded maninests are read-only
	// we duplicate them to avoid updating the shared manifests data in cache.
	manifests = duplicateManifests(manifests, "")

	// When addVariantLabelToSelector is true, ensure that all workloads
	// have the variant label in their selector.
	if opts.AddVariantLabelToSelector {
		workloads := findWorkloadManifests(manifests, e.config.Workloads)
		for _, m := range workloads {
			if err := ensureVariantSelectorInWorkload(m, primaryVariant); err != nil {
				return nil, fmt.Errorf("unable to check/set %q in selector of workload %s (%v)", variantLabel+": "+primaryVariant, m.Key.ReadableString(), err)
			}
		}
	}

	// Find service manifests and duplicate them for PRIMARY variant.
	if opts.CreateService {
		serviceName := e.config.Service.Name
		services := findManifests(provider.KindService, serviceName, manifests)
		if len(services) == 0 {
			return nil, fmt.Errorf("unable to find any service for name=%q", serviceName)
		}
		services = duplicateManifests(services, "")

		generatedServices, err := generateVariantServiceManifests(services, primaryVariant, suffix)
		if err != nil {
			return nil, err
		}
		manifests = append(manifests, generatedServices...)
	}

	return manifests, nil
}
