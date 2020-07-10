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
	"errors"
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
	canaryOptions := e.StageConfig.K8sCanaryRolloutStageOptions
	if canaryOptions == nil {
		e.LogPersister.AppendError(fmt.Sprintf("Malformed configuration for stage %s", e.Stage.Name))
		return model.StageStatus_STAGE_FAILURE
	}

	manifests, err := e.loadManifests(ctx)
	if err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Failed while loading manifests (%v)", err))
		return model.StageStatus_STAGE_FAILURE
	}

	if len(manifests) == 0 {
		e.LogPersister.AppendError("This application has no Kubernetes manifests to handle")
		return model.StageStatus_STAGE_FAILURE
	}

	canaryManifests, err := e.generateCanaryManifests(e.config.Input.Namespace, manifests, *canaryOptions)
	if err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Unable to generate manifests for CANARY variant (%v)", err))
		return model.StageStatus_STAGE_FAILURE
	}

	// Store added resource keys into metadata for cleaning later.
	addedResources := make([]string, 0, len(canaryManifests))
	for _, m := range canaryManifests {
		addedResources = append(addedResources, m.Key.String())
	}
	metadata := strings.Join(addedResources, ",")
	err = e.MetadataStore.Set(ctx, addedCanaryResourcesMetadataKey, metadata)
	if err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Unable to save deployment metadata (%v)", err))
		return model.StageStatus_STAGE_FAILURE
	}

	// Start rolling out the resources for CANARY variant.
	e.LogPersister.AppendInfo("Start rolling out CANARY variant...")
	for _, m := range canaryManifests {
		if err = e.provider.ApplyManifest(ctx, m); err != nil {
			e.LogPersister.AppendError(fmt.Sprintf("Failed to apply manifest: %s (%v)", m.Key.ReadableString(), err))
			return model.StageStatus_STAGE_FAILURE
		}
		e.LogPersister.AppendSuccess(fmt.Sprintf("- applied manifest: %s", m.Key.ReadableString()))
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
			e.LogPersister.AppendError(fmt.Sprintf("Had an error while decoding CANARY resource key: %s, %v", r, err))
			continue
		}
		if key.IsWorkload() {
			workloadKeys = append(workloadKeys, key)
		} else {
			serviceKeys = append(serviceKeys, key)
		}
	}

	// We delete the service first to close all incoming connections.
	for _, k := range serviceKeys {
		err := e.provider.Delete(ctx, k)
		if err == nil {
			e.LogPersister.AppendInfo(fmt.Sprintf("Deleted resource %s", k))
			continue
		}
		if errors.Is(err, provider.ErrNotFound) {
			e.LogPersister.AppendInfo(fmt.Sprintf("No resource %s to delete", k))
			continue
		}
		e.LogPersister.AppendError(fmt.Sprintf("Unable to delete resource %s (%v)", k, err))
		//return model.StageStatus_STAGE_FAILURE
	}

	// Next, delete all workloads.
	for _, k := range workloadKeys {
		err := e.provider.Delete(ctx, k)
		if err == nil {
			e.LogPersister.AppendInfo(fmt.Sprintf("Deleted workload resource %s", k))
			continue
		}
		if errors.Is(err, provider.ErrNotFound) {
			e.LogPersister.AppendInfo(fmt.Sprintf("No worload resource %s to delete", k))
			continue
		}
		e.LogPersister.AppendError(fmt.Sprintf("Unable to delete workload resource %s (%v)", k, err))
		//return model.StageStatus_STAGE_FAILURE
	}

	return nil
}

func (e *Executor) generateCanaryManifests(namespace string, manifests []provider.Manifest, opts config.K8sCanaryRolloutStageOptions) ([]provider.Manifest, error) {
	var (
		workloadKind, workloadName string
		serviceName                string
		generateService            bool
		canaryManifests            []provider.Manifest
		suffix                     = canaryVariant
	)

	// Apply the specified configuration if they are present.
	if sc := e.config.CanaryVariant; sc != nil {
		var ok bool
		if sc.Suffix != "" {
			suffix = sc.Suffix
		}
		generateService = sc.Service.Create

		workloadKind, workloadName, ok = config.ParseVariantResourceReference(sc.Workload.Reference)
		if !ok {
			return nil, fmt.Errorf("malformed workload reference: %s", sc.Workload.Reference)
		}
		if workloadKind == "" {
			workloadKind = provider.KindDeployment
		}

		_, serviceName, ok = config.ParseVariantResourceReference(sc.Service.Reference)
		if !ok {
			return nil, fmt.Errorf("malformed service reference: %s", sc.Workload.Reference)
		}
	}

	workloads := findManifests(workloadKind, workloadName, manifests)
	if len(workloads) == 0 {
		return nil, fmt.Errorf("unable to find any workload manifests for CANARY variant")
	}

	// Find service manifests and duplicate them for CANARY variant.
	if generateService {
		services := findManifests(provider.KindService, serviceName, manifests)
		generatedServices, err := generateServiceManifests(services, canaryVariant, suffix)
		if err != nil {
			return nil, err
		}
		canaryManifests = append(canaryManifests, generatedServices...)
	}

	// Find config map manifests and duplicate them for CANARY variant.
	configmaps := findConfigMapManifests(manifests)
	for _, m := range configmaps {
		m = m.Duplicate(m.Key.Name + "-" + suffix)
		canaryManifests = append(canaryManifests, m)
	}

	// Find secret manifests and duplicate them for CANARY variant.
	secrets := findSecretManifests(manifests)
	for _, m := range secrets {
		m = m.Duplicate(m.Key.Name + "-" + suffix)
		canaryManifests = append(canaryManifests, m)
	}

	// Generate new workload manifests for CANARY variant.
	// The generated ones will mount to the new ConfigMaps and Secrets.
	replicasCalculator := func(cur *int32) int32 {
		if cur == nil {
			return 1
		}
		num := opts.Replicas.Calculate(int(*cur), 1)
		return int32(num)
	}
	generatedWorkloads, err := generateWorkloadManifests(workloads, configmaps, secrets, canaryVariant, suffix, replicasCalculator)
	if err != nil {
		return nil, err
	}
	canaryManifests = append(canaryManifests, generatedWorkloads...)

	// Add labels to the generated canary manifests.
	for _, m := range canaryManifests {
		if namespace != "" {
			m.SetNamespace(namespace)
			m.Key.Namespace = namespace
		}
		m.AddAnnotations(e.builtinAnnotations(m, canaryVariant, e.Deployment.Trigger.Commit.Hash))
	}
	return canaryManifests, nil
}
