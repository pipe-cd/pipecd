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
		e.LogPersister.AppendError(fmt.Sprintf("Unabled to generate manifests for CANARY variant (%v)", err))
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
		e.LogPersister.AppendError(fmt.Sprintf("Unabled to save deployment metadata (%v)", err))
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
		e.LogPersister.AppendError("Unabled to determine the applied CANARY resources")
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
		e.LogPersister.AppendError(fmt.Sprintf("Unabled to delete resource %s (%v)", k, err))
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
		e.LogPersister.AppendError(fmt.Sprintf("Unabled to delete workload resource %s (%v)", k, err))
		//return model.StageStatus_STAGE_FAILURE
	}

	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) generateCanaryManifests(namespace string, manifests []provider.Manifest, opts config.K8sCanaryRolloutStageOptions) ([]provider.Manifest, error) {
	// List of default configurations.
	var (
		suffix           = canaryVariant
		workloadKind     = provider.KindDeployment
		workloadName     = ""
		workloadReplicas = 1
		foundWorkload    = false
		canaryManifests  []provider.Manifest
	)

	// Apply the specified configuration if they are present.
	if sc := e.config.CanaryVariant; sc != nil {
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
		if err := m.AddVariantLabel(canaryVariant); err != nil {
			return err
		}
		m.SetReplicas(workloadReplicas)
		canaryManifests = append(canaryManifests, m)
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
		return nil, fmt.Errorf("unabled to detect workload manifest for CANARY variant")
	}

	// TODO: Generate ConfigMap and Secret manifests and update Workload to use the copy.
	// TODO: Generate Service manifest for kubernetes CANARY variant.

	// Add labels to the generated canary manifests.
	for _, m := range canaryManifests {
		m.Key.Name = m.Key.Name + "-" + suffix
		if namespace != "" {
			m.SetNamespace(namespace)
			m.Key.Namespace = namespace
		}
		m.AddAnnotations(map[string]string{
			provider.LabelManagedBy:          provider.ManagedByPiped,
			provider.LabelPiped:              e.PipedConfig.PipedID,
			provider.LabelApplication:        e.Deployment.ApplicationId,
			provider.LabelVariant:            canaryVariant,
			provider.LabelOriginalAPIVersion: m.Key.APIVersion,
			provider.LabelResourceKey:        m.Key.String(),
			provider.LabelCommitHash:         e.Deployment.Trigger.Commit.Hash,
		})
	}
	return canaryManifests, nil
}
