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
	"github.com/pipe-cd/pipe/pkg/model"
)

const (
	baselineVariant                   = "baseline"
	addedBaselineResourcesMetadataKey = "baseline-resources"
)

func (e *Executor) ensureBaselineRollout(ctx context.Context) model.StageStatus {
	manifests, err := e.loadRunningManifests(ctx)
	if err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Failed while loading running manifests (%v)", err))
		return model.StageStatus_STAGE_FAILURE
	}

	if len(manifests) == 0 {
		e.LogPersister.AppendError("This application has no running Kubernetes manifests to handle")
		return model.StageStatus_STAGE_FAILURE
	}

	baselineManifests, err := e.generateBaselineManifests(e.config.Input.Namespace, manifests)
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

	var (
		resources    = strings.Split(value, ",")
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

	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) generateBaselineManifests(namespace string, manifests []provider.Manifest) ([]provider.Manifest, error) {
	// List of default configurations.
	var (
		suffix            = baselineVariant
		workloadKind      = provider.KindDeployment
		workloadName      = ""
		workloadReplicas  = 1
		foundWorkload     = false
		baselineManifests []provider.Manifest
	)

	// Apply the specified configuration if they are present.
	if sc := e.config.BaselineVariant; sc != nil {
		if sc.Suffix != "" {
			suffix = sc.Suffix
		}
		if sc.Workload.Kind != "" {
			workloadKind = sc.Workload.Kind
		}
		if sc.Workload.Name != "" {
			workloadName = sc.Workload.Name
		}
	}

	findWorkload := func(m provider.Manifest) error {
		if m.Key.Kind != workloadKind {
			return nil
		}
		if workloadName != "" && m.Key.Name != workloadName {
			return nil
		}
		m = m.Duplicate(m.Key.Name + "-" + suffix)
		if err := m.AddVariantLabel(baselineVariant); err != nil {
			return err
		}
		// TODO: Load baseline replicas number from configuration.
		m.SetReplicas(workloadReplicas)
		baselineManifests = append(baselineManifests, m)
		foundWorkload = true
		return nil
	}

	for _, m := range manifests {
		if err := findWorkload(m); err != nil {
			return nil, err
		}
		if foundWorkload {
			break
		}
	}

	if !foundWorkload {
		return nil, fmt.Errorf("unable to detect workload manifest for BASELINE variant")
	}

	// TODO: Generate Service manifest for kubernetes BASELINE variant.

	// Add labels to the generated baseline manifests.
	for _, m := range baselineManifests {
		m.Key.Name = m.Key.Name + "-" + suffix
		if namespace != "" {
			m.SetNamespace(namespace)
			m.Key.Namespace = namespace
		}
		m.AddAnnotations(map[string]string{
			provider.LabelManagedBy:          provider.ManagedByPiped,
			provider.LabelPiped:              e.PipedConfig.PipedID,
			provider.LabelApplication:        e.Deployment.ApplicationId,
			provider.LabelVariant:            baselineVariant,
			provider.LabelOriginalAPIVersion: m.Key.APIVersion,
			provider.LabelResourceKey:        m.Key.String(),
			provider.LabelCommitHash:         e.Deployment.Trigger.Commit.Hash,
		})
	}
	return baselineManifests, nil
}
