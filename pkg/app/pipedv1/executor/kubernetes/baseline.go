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
	"fmt"
	"strings"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/executor"
	provider "github.com/pipe-cd/pipecd/pkg/app/pipedv1/platformprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	addedBaselineResourcesMetadataKey = "baseline-resources"
)

func (e *deployExecutor) ensureBaselineRollout(ctx context.Context) model.StageStatus {
	var (
		runningCommit   = e.Deployment.RunningCommitHash
		options         = e.StageConfig.K8sBaselineRolloutStageOptions
		variantLabel    = e.appCfg.VariantLabel.Key
		baselineVariant = e.appCfg.VariantLabel.BaselineValue
	)
	if options == nil {
		e.LogPersister.Errorf("Malformed configuration for stage %s", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}

	// Load running manifests at the most successful deployed commit.
	e.LogPersister.Infof("Loading running manifests at commit %s for handling", runningCommit)
	manifests, err := e.loadRunningManifests(ctx)
	if err != nil {
		e.LogPersister.Errorf("Failed while loading running manifests (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.Successf("Successfully loaded %d manifests", len(manifests))

	if len(manifests) == 0 {
		e.LogPersister.Error("This application has no running Kubernetes manifests to handle")
		return model.StageStatus_STAGE_FAILURE
	}

	baselineManifests, err := e.generateBaselineManifests(manifests, *options, variantLabel, baselineVariant)
	if err != nil {
		e.LogPersister.Errorf("Unable to generate manifests for BASELINE variant (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Add builtin annotations for tracking application live state.
	addBuiltinAnnotations(
		baselineManifests,
		variantLabel,
		baselineVariant,
		runningCommit,
		e.PipedConfig.PipedID,
		e.Deployment.ApplicationId,
	)

	// Store added resource keys into metadata for cleaning later.
	addedResources := make([]string, 0, len(baselineManifests))
	for _, m := range baselineManifests {
		addedResources = append(addedResources, m.Key.String())
	}
	metadata := strings.Join(addedResources, ",")
	err = e.MetadataStore.Shared().Put(ctx, addedBaselineResourcesMetadataKey, metadata)
	if err != nil {
		e.LogPersister.Errorf("Unable to save deployment metadata (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Start rolling out the resources for BASELINE variant.
	e.LogPersister.Info("Start rolling out BASELINE variant...")
	if err := applyManifests(ctx, e.applierGetter, baselineManifests, e.appCfg.Input.Namespace, e.LogPersister); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}

	e.LogPersister.Success("Successfully rolled out BASELINE variant")
	return model.StageStatus_STAGE_SUCCESS
}

func (e *deployExecutor) ensureBaselineClean(ctx context.Context) model.StageStatus {
	value, ok := e.MetadataStore.Shared().Get(addedBaselineResourcesMetadataKey)
	if !ok {
		e.LogPersister.Error("Unable to determine the applied BASELINE resources")
		return model.StageStatus_STAGE_FAILURE
	}

	resources := strings.Split(value, ",")
	if err := removeBaselineResources(ctx, e.applierGetter, resources, e.LogPersister); err != nil {
		e.LogPersister.Errorf("Unable to remove baseline resources: %v", err)
		return model.StageStatus_STAGE_FAILURE
	}
	return model.StageStatus_STAGE_SUCCESS
}

func (e *deployExecutor) generateBaselineManifests(manifests []provider.Manifest, opts config.K8sBaselineRolloutStageOptions, variantLabel, variant string) ([]provider.Manifest, error) {
	suffix := variant
	if opts.Suffix != "" {
		suffix = opts.Suffix
	}

	workloads := findWorkloadManifests(manifests, e.appCfg.Workloads)
	if len(workloads) == 0 {
		return nil, fmt.Errorf("unable to find any workload manifests for BASELINE variant")
	}

	var baselineManifests []provider.Manifest

	// Find service manifests and duplicate them for BASELINE variant.
	if opts.CreateService {
		serviceName := e.appCfg.Service.Name
		services := findManifests(provider.KindService, serviceName, manifests)
		if len(services) == 0 {
			return nil, fmt.Errorf("unable to find any service for name=%q", serviceName)
		}
		// Because the loaded manifests are read-only
		// so we duplicate them to avoid updating the shared manifests data in cache.
		services = duplicateManifests(services, "")

		generatedServices, err := generateVariantServiceManifests(services, variantLabel, variant, suffix)
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
	generatedWorkloads, err := generateVariantWorkloadManifests(workloads, nil, nil, variantLabel, variant, suffix, replicasCalculator)
	if err != nil {
		return nil, err
	}
	baselineManifests = append(baselineManifests, generatedWorkloads...)

	return baselineManifests, nil
}

func removeBaselineResources(ctx context.Context, ag applierGetter, resources []string, lp executor.LogPersister) error {
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
			lp.Errorf("Had an error while decoding BASELINE resource key: %s, %v", r, err)
			continue
		}
		if key.IsWorkload() {
			workloadKeys = append(workloadKeys, key)
		} else {
			serviceKeys = append(serviceKeys, key)
		}
	}

	// We delete the service first to close all incoming connections.
	lp.Info("Starting finding and deleting service resources of BASELINE variant")
	if err := deleteResources(ctx, ag, serviceKeys, lp); err != nil {
		return err
	}

	// Next, delete all workloads.
	lp.Info("Starting finding and deleting workload resources of BASELINE variant")
	if err := deleteResources(ctx, ag, workloadKeys, lp); err != nil {
		return err
	}

	return nil
}
