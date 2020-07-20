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
	"strings"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

const (
	baselineVariant                   = "baseline"
	addedBaselineResourcesMetadataKey = "baseline-resources"
)

func (e *Executor) ensureBaselineRollout(ctx context.Context) model.StageStatus {
	var (
		commitHash = e.Deployment.RunningCommitHash
		options    = e.StageConfig.K8sBaselineRolloutStageOptions
	)
	if options == nil {
		e.LogPersister.AppendErrorf("Malformed configuration for stage %s", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}

	// Load running manifests at the most successful deployed commit.
	e.LogPersister.AppendInfof("Loading running manifests at commit %s for handling", commitHash)
	manifests, err := e.loadRunningManifests(ctx)
	if err != nil {
		e.LogPersister.AppendErrorf("Failed while loading running manifests (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.AppendSuccessf("Successfully loaded %d manifests", len(manifests))

	if len(manifests) == 0 {
		e.LogPersister.AppendError("This application has no running Kubernetes manifests to handle")
		return model.StageStatus_STAGE_FAILURE
	}

	baselineManifests, err := e.generateBaselineManifests(manifests, *options)
	if err != nil {
		e.LogPersister.AppendErrorf("Unable to generate manifests for BASELINE variant (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Add builtin annotations for tracking application live state.
	e.addBuiltinAnnontations(baselineManifests, baselineVariant, commitHash)

	// Store added resource keys into metadata for cleaning later.
	addedResources := make([]string, 0, len(baselineManifests))
	for _, m := range baselineManifests {
		addedResources = append(addedResources, m.Key.String())
	}
	metadata := strings.Join(addedResources, ",")
	err = e.MetadataStore.Set(ctx, addedBaselineResourcesMetadataKey, metadata)
	if err != nil {
		e.LogPersister.AppendErrorf("Unable to save deployment metadata (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Start rolling out the resources for BASELINE variant.
	e.LogPersister.AppendInfo("Start rolling out BASELINE variant...")
	if err := e.applyManifests(ctx, baselineManifests); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}

	e.LogPersister.AppendSuccess("Successfully rolled out BASELINE variant")
	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) ensureBaselineClean(ctx context.Context) model.StageStatus {
	value, ok := e.MetadataStore.Get(addedBaselineResourcesMetadataKey)
	if !ok {
		e.LogPersister.AppendError("Unable to determine the applied BASELINE resources")
		return model.StageStatus_STAGE_FAILURE
	}

	resources := strings.Split(value, ",")
	if err := e.removeBaselineResources(ctx, resources); err != nil {
		e.LogPersister.AppendErrorf("Unable to remove baseline resources: %v", err)
		return model.StageStatus_STAGE_FAILURE
	}
	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) removeBaselineResources(ctx context.Context, resources []string) error {
	if len(resources) == 0 {
		return nil
	}

	var (
		workloadKeys = make([]provider.ResourceKey, 0)
		serviceKeys  = make([]provider.ResourceKey, 0)
	)
	for _, r := range resources {
		key, err := provider.DecodeResourceKey(r)
		if err != nil {
			e.LogPersister.AppendErrorf("Had an error while decoding BASELINE resource key: %s, %v", r, err)
			continue
		}
		if key.IsWorkload() {
			workloadKeys = append(workloadKeys, key)
		} else {
			serviceKeys = append(serviceKeys, key)
		}
	}

	// We delete the service first to close all incoming connections.
	e.LogPersister.AppendInfo("Starting finding and deleting service resources of BASELINE variant")
	if err := e.deleteResources(ctx, serviceKeys); err != nil {
		return err
	}

	// Next, delete all workloads.
	e.LogPersister.AppendInfo("Starting finding and deleting workload resources of BASELINE variant")
	if err := e.deleteResources(ctx, workloadKeys); err != nil {
		return err
	}

	return nil
}

func (e *Executor) generateBaselineManifests(manifests []provider.Manifest, opts config.K8sBaselineRolloutStageOptions) ([]provider.Manifest, error) {
	var (
		workloadKind, workloadName string
		serviceName                string
		generateService            bool
		baselineManifests          []provider.Manifest
		suffix                     = baselineVariant
	)

	// Apply the specified configuration if they are present.
	if sc := e.config.BaselineVariant; sc != nil {
		var ok bool
		if sc.Suffix != "" {
			suffix = sc.Suffix
		}
		generateService = sc.Service.Create

		workloadKind, workloadName, ok = config.ParseVariantResourceReference(sc.Workload.Reference)
		if !ok {
			return nil, fmt.Errorf("malformed workload reference: %s", sc.Workload.Reference)
		}

		_, serviceName, ok = config.ParseVariantResourceReference(sc.Service.Reference)
		if !ok {
			return nil, fmt.Errorf("malformed service reference: %s", sc.Service.Reference)
		}
	}
	if workloadKind == "" {
		workloadKind = provider.KindDeployment
	}

	workloads := findManifests(workloadKind, workloadName, manifests)
	if len(workloads) == 0 {
		return nil, fmt.Errorf("unable to find any workload manifests for BASELINE variant")
	}

	// Find service manifests and duplicate them for BASELINE variant.
	if generateService {
		services := findManifests(provider.KindService, serviceName, manifests)
		if len(services) == 0 {
			return nil, fmt.Errorf("unable to find any service for name=%q", serviceName)
		}
		// Because the loaded maninests are read-only
		// so we duplicate them to avoid updating the shared manifests data in cache.
		services = duplicateManifests(services, "")

		generatedServices, err := generateVariantServiceManifests(services, baselineVariant, suffix)
		if err != nil {
			return nil, err
		}
		baselineManifests = append(baselineManifests, generatedServices...)
	}

	// Generate new workload manifests for VANARY variant.
	// The generated ones will mount to the new ConfigMaps and Secrets.
	replicasCalculator := func(cur *int32) int32 {
		if cur == nil {
			return 1
		}
		num := opts.Replicas.Calculate(int(*cur), 1)
		return int32(num)
	}
	generatedWorkloads, err := generateVariantWorkloadManifests(workloads, nil, nil, baselineVariant, suffix, replicasCalculator)
	if err != nil {
		return nil, err
	}
	baselineManifests = append(baselineManifests, generatedWorkloads...)

	return baselineManifests, nil
}
