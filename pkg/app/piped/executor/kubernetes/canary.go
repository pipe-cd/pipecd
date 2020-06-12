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
	"github.com/pipe-cd/pipe/pkg/model"
)

func (e *Executor) ensureCanaryRollout(ctx context.Context) model.StageStatus {
	manifests, err := e.loadManifests(ctx)
	if err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Failed while loading manifests (%v)", err))
		return model.StageStatus_STAGE_FAILURE
	}

	if len(manifests) == 0 {
		e.LogPersister.AppendError("No kubernetes manifests to handle")
		return model.StageStatus_STAGE_FAILURE
	}

	stageManifests, err := e.generateStageManifests(ctx, manifests)
	if err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Unabled to generate manifests for CANARY variant (%v)", err))
		return model.StageStatus_STAGE_FAILURE
	}

	// Store will adding resource keys into metadata for cleaning later.
	addedResources := make([]string, 0, len(stageManifests))
	for _, m := range stageManifests {
		addedResources = append(addedResources, m.Key.String())
	}
	metadata := strings.Join(addedResources, ",")
	err = e.MetadataStore.Set(ctx, metadataKeyAddedStageResources, metadata)
	if err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Unabled to save deployment metadata (%v)", err))
	}

	e.LogPersister.AppendInfo("Rolling out CANARY variant")
	if err = e.provider.ApplyManifests(ctx, stageManifests); err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Unabled to rollout CANARY variant (%v)", err))
		return model.StageStatus_STAGE_FAILURE
	}

	e.LogPersister.AppendSuccess("Successfully rolled out CANARY variant")
	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) ensureCanaryClean(ctx context.Context) model.StageStatus {
	value, ok := e.MetadataStore.Get(metadataKeyAddedStageResources)
	if !ok {
		// We have to re-render manifests to check stage resources.
		value = ""
	}
	var (
		resources    = strings.Split(value, ",")
		workloadKeys = make([]provider.ResourceKey, 0)
		serviceKeys  = make([]provider.ResourceKey, 0)
	)
	for _, r := range resources {
		key, _ := provider.DecodeResourceKey(r)
		switch key.Kind {
		case "Deployment", "ReplicaSet", "DaemonSet", "Pod":
			workloadKeys = append(workloadKeys, key)
		default:
			serviceKeys = append(serviceKeys, key)
		}
	}

	// We delete the service first to close all incoming connections.
	for _, k := range serviceKeys {
		if err := e.provider.Delete(ctx, k); err != nil {
			e.LogPersister.AppendError(fmt.Sprintf("Unabled to delete resource %s (%v)", k, err))
			continue
		}
		e.LogPersister.AppendInfo(fmt.Sprintf("Deleted resource %s", k))
	}

	// Next, delete all workloads.
	for _, k := range workloadKeys {
		if err := e.provider.Delete(ctx, k); err != nil {
			e.LogPersister.AppendError(fmt.Sprintf("Unabled to delete workload resource %s (%v)", k, err))
			continue
		}
		e.LogPersister.AppendInfo(fmt.Sprintf("Deleted workload resource %s", k))
	}

	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) generateStageManifests(ctx context.Context, manifests []provider.Manifest) ([]provider.Manifest, error) {
	// List of default configurations.
	var (
		suffix           = canaryVariant
		workloadKind     = "Deployment"
		workloadName     = ""
		workloadReplicas = 1
		foundWorkload    = false
		stageManifests   []provider.Manifest
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
		stageManifests = append(stageManifests, m)
		foundWorkload = true
		return nil
	}

	for _, m := range manifests {
		if err := findWorkload(m); err != nil {
			return nil, err
		}
	}

	if !foundWorkload {
		return nil, fmt.Errorf("unabled to detect workload manifest for CANARY variant")
	}

	for _, m := range stageManifests {
		m.Key.Name = m.Key.Name + "-" + suffix
		m.AddAnnotations(map[string]string{
			provider.LabelManagedBy:          provider.ManagedByPiped,
			provider.LabelApplication:        e.Deployment.ApplicationId,
			provider.LabelVariant:            canaryVariant,
			provider.LabelOriginalAPIVersion: m.Key.APIVersion,
			provider.LabelResourceKey:        m.Key.String(),
			provider.LabelCommitHash:         e.Deployment.Trigger.Commit.Hash,
		})
	}
	return stageManifests, nil
}
