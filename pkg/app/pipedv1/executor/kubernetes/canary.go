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
	addedCanaryResourcesMetadataKey = "canary-resources"
)

func (e *deployExecutor) ensureCanaryRollout(ctx context.Context) model.StageStatus {
	var (
		options       = e.StageConfig.K8sCanaryRolloutStageOptions
		variantLabel  = e.appCfg.VariantLabel.Key
		canaryVariant = e.appCfg.VariantLabel.CanaryValue
	)
	if options == nil {
		e.LogPersister.Errorf("Malformed configuration for stage %s", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}

	// Load the manifests at the triggered commit.
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

	if len(manifests) == 0 {
		e.LogPersister.Error("This application has no Kubernetes manifests to handle")
		return model.StageStatus_STAGE_FAILURE
	}

	// Patches the manifests if needed.
	if len(options.Patches) > 0 {
		e.LogPersister.Info("Patching manifests before generating for CANARY variant")
		manifests, err = patchManifests(manifests, options.Patches, patchManifest)
		if err != nil {
			e.LogPersister.Errorf("Failed while patching manifests (%v)", err)
			return model.StageStatus_STAGE_FAILURE
		}
	}

	// Find and generate workload & service manifests for CANARY variant.
	canaryManifests, err := e.generateCanaryManifests(manifests, *options, variantLabel, canaryVariant)
	if err != nil {
		e.LogPersister.Errorf("Unable to generate manifests for CANARY variant (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Add builtin annotations for tracking application live state.
	addBuiltinAnnotations(
		canaryManifests,
		variantLabel,
		canaryVariant,
		e.commit,
		e.PipedConfig.PipedID,
		e.Deployment.ApplicationId,
	)

	// Store added resource keys into metadata for cleaning later.
	addedResources := make([]string, 0, len(canaryManifests))
	for _, m := range canaryManifests {
		addedResources = append(addedResources, m.Key.String())
	}
	metadata := strings.Join(addedResources, ",")
	err = e.MetadataStore.Shared().Put(ctx, addedCanaryResourcesMetadataKey, metadata)
	if err != nil {
		e.LogPersister.Errorf("Unable to save deployment metadata (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Start rolling out the resources for CANARY variant.
	e.LogPersister.Info("Start rolling out CANARY variant...")
	if err := applyManifests(ctx, e.applierGetter, canaryManifests, e.appCfg.Input.Namespace, e.LogPersister); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}

	e.LogPersister.Success("Successfully rolled out CANARY variant")
	return model.StageStatus_STAGE_SUCCESS
}

func (e *deployExecutor) ensureCanaryClean(ctx context.Context) model.StageStatus {
	value, ok := e.MetadataStore.Shared().Get(addedCanaryResourcesMetadataKey)
	if !ok {
		e.LogPersister.Error("Unable to determine the applied CANARY resources")
		return model.StageStatus_STAGE_FAILURE
	}

	resources := strings.Split(value, ",")
	if err := removeCanaryResources(ctx, e.applierGetter, resources, e.LogPersister); err != nil {
		e.LogPersister.Errorf("Unable to remove canary resources: %v", err)
		return model.StageStatus_STAGE_FAILURE
	}
	return model.StageStatus_STAGE_SUCCESS
}

func (e *deployExecutor) generateCanaryManifests(manifests []provider.Manifest, opts config.K8sCanaryRolloutStageOptions, variantLabel, variant string) ([]provider.Manifest, error) {
	suffix := variant
	if opts.Suffix != "" {
		suffix = opts.Suffix
	}

	workloads := findWorkloadManifests(manifests, e.appCfg.Workloads)
	if len(workloads) == 0 {
		return nil, fmt.Errorf("unable to find any workload manifests for CANARY variant")
	}

	var canaryManifests []provider.Manifest

	// Find service manifests and duplicate them for CANARY variant.
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
		canaryManifests = append(canaryManifests, generatedServices...)
	}

	// Find config map manifests and duplicate them for CANARY variant.
	configMaps := findConfigMapManifests(manifests)
	canaryConfigMaps := duplicateManifests(configMaps, suffix)
	canaryManifests = append(canaryManifests, canaryConfigMaps...)

	// Find secret manifests and duplicate them for CANARY variant.
	secrets := findSecretManifests(manifests)
	canarySecrets := duplicateManifests(secrets, suffix)
	canaryManifests = append(canaryManifests, canarySecrets...)

	// Generate new workload manifests for CANARY variant.
	// The generated ones will mount to the new ConfigMaps and Secrets.
	replicasCalculator := func(cur *int32) int32 {
		if cur == nil {
			return 1
		}
		num := opts.Replicas.Calculate(int(*cur), 1)
		return int32(num)
	}
	// We don't need to duplicate the workload manifests
	// because generateVariantWorkloadManifests function is already making a duplicate while decoding.
	// workloads = duplicateManifests(workloads, suffix)
	generatedWorkloads, err := generateVariantWorkloadManifests(workloads, configMaps, secrets, variantLabel, variant, suffix, replicasCalculator)
	if err != nil {
		return nil, err
	}
	canaryManifests = append(canaryManifests, generatedWorkloads...)

	return canaryManifests, nil
}

func removeCanaryResources(ctx context.Context, ag applierGetter, resources []string, lp executor.LogPersister) error {
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
			lp.Errorf("Had an error while decoding CANARY resource key: %s, %v", r, err)
			continue
		}
		if key.IsWorkload() {
			workloadKeys = append(workloadKeys, key)
		} else {
			serviceKeys = append(serviceKeys, key)
		}
	}

	// We delete the service first to close all incoming connections.
	lp.Info("Starting finding and deleting service resources of CANARY variant")
	if err := deleteResources(ctx, ag, serviceKeys, lp); err != nil {
		return err
	}

	// Next, delete all workloads.
	lp.Info("Starting finding and deleting workload resources of CANARY variant")
	if err := deleteResources(ctx, ag, workloadKeys, lp); err != nil {
		return err
	}

	return nil
}
