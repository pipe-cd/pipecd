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
	canaryVariant                   = "canary"
	addedCanaryResourcesMetadataKey = "canary-resources"
)

func (e *Executor) ensureCanaryRollout(ctx context.Context) model.StageStatus {
	var (
		commitHash = e.Deployment.Trigger.Commit.Hash
		options    = e.StageConfig.K8sCanaryRolloutStageOptions
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

	if len(manifests) == 0 {
		e.LogPersister.AppendError("This application has no Kubernetes manifests to handle")
		return model.StageStatus_STAGE_FAILURE
	}

	// Find and generate workload & service manifests for CANARY variant.
	canaryManifests, err := e.generateCanaryManifests(manifests, *options)
	if err != nil {
		e.LogPersister.AppendErrorf("Unable to generate manifests for CANARY variant (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Add builtin annotations for tracking application live state.
	e.addBuiltinAnnontations(canaryManifests, canaryVariant, commitHash)

	// Store added resource keys into metadata for cleaning later.
	addedResources := make([]string, 0, len(canaryManifests))
	for _, m := range canaryManifests {
		addedResources = append(addedResources, m.Key.String())
	}
	metadata := strings.Join(addedResources, ",")
	err = e.MetadataStore.Set(ctx, addedCanaryResourcesMetadataKey, metadata)
	if err != nil {
		e.LogPersister.AppendErrorf("Unable to save deployment metadata (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Start rolling out the resources for CANARY variant.
	e.LogPersister.AppendInfo("Start rolling out CANARY variant...")
	if err := e.applyManifests(ctx, canaryManifests); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}

	e.LogPersister.AppendSuccess("Successfully rolled out CANARY variant")
	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) ensureCanaryClean(ctx context.Context) model.StageStatus {
	value, ok := e.MetadataStore.Get(addedCanaryResourcesMetadataKey)
	if !ok {
		e.LogPersister.AppendError("Unable to determine the applied CANARY resources")
		return model.StageStatus_STAGE_FAILURE
	}

	resources := strings.Split(value, ",")
	if err := e.removeCanaryResources(ctx, resources); err != nil {
		e.LogPersister.AppendErrorf("Unable to remove canary resources: %v", err)
		return model.StageStatus_STAGE_FAILURE
	}
	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) removeCanaryResources(ctx context.Context, resources []string) error {
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
			e.LogPersister.AppendErrorf("Had an error while decoding CANARY resource key: %s, %v", r, err)
			continue
		}
		if key.IsWorkload() {
			workloadKeys = append(workloadKeys, key)
		} else {
			serviceKeys = append(serviceKeys, key)
		}
	}

	// We delete the service first to close all incoming connections.
	e.LogPersister.AppendInfo("Starting finding and deleting service resources of CANARY variant")
	if err := e.deleteResources(ctx, serviceKeys); err != nil {
		return err
	}

	// Next, delete all workloads.
	e.LogPersister.AppendInfo("Starting finding and deleting workload resources of CANARY variant")
	if err := e.deleteResources(ctx, workloadKeys); err != nil {
		return err
	}

	return nil
}

func (e *Executor) generateCanaryManifests(manifests []provider.Manifest, opts config.K8sCanaryRolloutStageOptions) ([]provider.Manifest, error) {
	suffix := canaryVariant
	if opts.Suffix != "" {
		suffix = opts.Suffix
	}

	workloads := findWorkloadManifests(manifests, e.config.Workloads)
	if len(workloads) == 0 {
		return nil, fmt.Errorf("unable to find any workload manifests for CANARY variant")
	}

	var canaryManifests []provider.Manifest

	// Find service manifests and duplicate them for CANARY variant.
	if opts.CreateService {
		serviceName := e.config.Service.Name
		services := findManifests(provider.KindService, serviceName, manifests)
		if len(services) == 0 {
			return nil, fmt.Errorf("unable to find any service for name=%q", serviceName)
		}
		// Because the loaded maninests are read-only
		// so we duplicate them to avoid updating the shared manifests data in cache.
		services = duplicateManifests(services, "")

		generatedServices, err := generateVariantServiceManifests(services, canaryVariant, suffix)
		if err != nil {
			return nil, err
		}
		canaryManifests = append(canaryManifests, generatedServices...)
	}

	// Find config map manifests and duplicate them for CANARY variant.
	configMaps := findConfigMapManifests(manifests)
	configMaps = duplicateManifests(configMaps, suffix)
	canaryManifests = append(canaryManifests, configMaps...)

	// Find secret manifests and duplicate them for CANARY variant.
	secrets := findSecretManifests(manifests)
	secrets = duplicateManifests(secrets, suffix)
	canaryManifests = append(canaryManifests, secrets...)

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
	generatedWorkloads, err := generateVariantWorkloadManifests(workloads, configMaps, secrets, canaryVariant, suffix, replicasCalculator)
	if err != nil {
		return nil, err
	}
	canaryManifests = append(canaryManifests, generatedWorkloads...)

	return canaryManifests, nil
}
