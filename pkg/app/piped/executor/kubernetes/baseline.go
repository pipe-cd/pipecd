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
	baselineVariant                   = "baseline"
	addedBaselineResourcesMetadataKey = "baseline-resources"
)

func (e *Executor) ensureBaselineRollout(ctx context.Context) model.StageStatus {
	baselineOptions := e.StageConfig.K8sBaselineRolloutStageOptions
	if baselineOptions == nil {
		e.LogPersister.AppendError(fmt.Sprintf("Malformed configuration for stage %s", e.Stage.Name))
		return model.StageStatus_STAGE_FAILURE
	}

	manifests, err := e.loadRunningManifests(ctx)
	if err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Failed while loading running manifests (%v)", err))
		return model.StageStatus_STAGE_FAILURE
	}

	if len(manifests) == 0 {
		e.LogPersister.AppendError("This application has no running Kubernetes manifests to handle")
		return model.StageStatus_STAGE_FAILURE
	}

	baselineManifests, err := e.generateBaselineManifests(e.config.Input.Namespace, manifests, *baselineOptions)
	if err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Unable to generate manifests for BASELINE variant (%v)", err))
		return model.StageStatus_STAGE_FAILURE
	}

	// Store added resource keys into metadata for cleaning later.
	addedResources := make([]string, 0, len(baselineManifests))
	for _, m := range baselineManifests {
		addedResources = append(addedResources, m.Key.String())
	}
	metadata := strings.Join(addedResources, ",")
	err = e.MetadataStore.Set(ctx, addedBaselineResourcesMetadataKey, metadata)
	if err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Unable to save deployment metadata (%v)", err))
		return model.StageStatus_STAGE_FAILURE
	}

	// Start rolling out the resources for BASELINE variant.
	e.LogPersister.AppendInfo("Start rolling out BASELINE variant...")
	for _, m := range baselineManifests {
		if err = e.provider.ApplyManifest(ctx, m); err != nil {
			e.LogPersister.AppendError(fmt.Sprintf("Failed to apply manifest: %s (%v)", m.Key.ReadableString(), err))
			return model.StageStatus_STAGE_FAILURE
		}
		e.LogPersister.AppendSuccess(fmt.Sprintf("- applied manifest: %s", m.Key.ReadableString()))
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
			e.LogPersister.AppendError(fmt.Sprintf("Had an error while decoding BASELINE resource key: %s, %v", r, err))
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

func (e *Executor) generateBaselineManifests(namespace string, manifests []provider.Manifest, opts config.K8sBaselineRolloutStageOptions) ([]provider.Manifest, error) {
	// List of default configurations.
	var (
		suffix            = baselineVariant
		workloadCfg       *config.K8sWorkload
		serviceCfg        *config.K8sService
		baselineManifests []provider.Manifest
	)

	// Apply the specified configuration if they are present.
	if sc := e.config.BaselineVariant; sc != nil {
		if sc.Suffix != "" {
			suffix = sc.Suffix
		}
		workloadCfg = sc.Workload
		serviceCfg = sc.Service
	}

	workloads := findWorkloadManifests(workloadCfg, manifests)
	if len(workloads) == 0 {
		return nil, fmt.Errorf("unable to find any workload manifests for BASELINE variant")
	}

	// Find service manifests and duplicate them for BASELINE variant.
	services := findServiceManifests(serviceCfg, manifests)
	generatedServices, err := generateServiceManifests(services, baselineVariant, suffix)
	if err != nil {
		return nil, err
	}
	baselineManifests = append(baselineManifests, generatedServices...)

	// Generate new workload manifests for VANARY variant.
	// The generated ones will mount to the new ConfigMaps and Secrets.
	replicasCalculator := func(cur *int32) int32 {
		if cur == nil {
			return 1
		}
		num := opts.Replicas.Calculate(int(*cur), 1)
		return int32(num)
	}
	generatedWorkloads, err := generateWorkloadManifests(workloads, nil, nil, baselineVariant, suffix, replicasCalculator)
	if err != nil {
		return nil, err
	}
	baselineManifests = append(baselineManifests, generatedWorkloads...)

	// Add labels to the generated baseline manifests.
	for _, m := range baselineManifests {
		if namespace != "" {
			m.SetNamespace(namespace)
			m.Key.Namespace = namespace
		}
		m.AddAnnotations(e.builtinAnnotations(m, baselineVariant, e.Deployment.RunningCommitHash))
	}
	return baselineManifests, nil
}
