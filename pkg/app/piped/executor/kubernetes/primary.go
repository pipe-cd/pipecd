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

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

const (
	primaryVariant = "primary"
)

func (e *Executor) ensurePrimaryUpdate(ctx context.Context) model.StageStatus {
	manifests, err := e.loadManifests(ctx)
	if err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Failed while loading manifests (%v)", err))
		return model.StageStatus_STAGE_FAILURE
	}

	if len(manifests) == 0 {
		e.LogPersister.AppendError("There are no kubernetes manifests to handle")
		return model.StageStatus_STAGE_FAILURE
	}

	primaryManifests, err := e.generatePrimaryManifests(e.config.Input.Namespace, e.Deployment.Trigger.Commit.Hash, manifests)
	if err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Unable to generate manifests for PRIMARY variant (%v)", err))
		return model.StageStatus_STAGE_FAILURE
	}

	e.LogPersister.AppendInfo(fmt.Sprintf("Applying %d primary resources", len(primaryManifests)))
	for _, m := range primaryManifests {
		if err = e.provider.ApplyManifest(ctx, m); err != nil {
			e.LogPersister.AppendError(fmt.Sprintf("Failed to apply manifest: %s (%v)", m.Key.ReadableString(), err))
			return model.StageStatus_STAGE_FAILURE
		}
		e.LogPersister.AppendSuccess(fmt.Sprintf("- applied manifest: %s", m.Key.ReadableString()))
	}

	e.LogPersister.AppendSuccess(fmt.Sprintf("Successfully applied %d primary resources", len(primaryManifests)))
	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) rollbackPrimary(ctx context.Context) error {
	manifests, err := e.loadRunningManifests(ctx)
	if err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Failed while loading running manifests (%v)", err))
		return err
	}

	if len(manifests) == 0 {
		e.LogPersister.AppendError("This application has no running Kubernetes manifests to handle")
		return err
	}

	primaryManifests, err := e.generatePrimaryManifests(e.config.Input.Namespace, e.Deployment.RunningCommitHash, manifests)
	if err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Unable to generate manifests for PRIMARY variant (%v)", err))
		return err
	}

	// Start rolling out the resources for PRIMARY variant.
	e.LogPersister.AppendInfo("Start rolling back PRIMARY variant...")
	for _, m := range primaryManifests {
		if err = e.provider.ApplyManifest(ctx, m); err != nil {
			e.LogPersister.AppendError(fmt.Sprintf("Failed to apply manifest: %s (%v)", m.Key.ReadableString(), err))
			return err
		}
		e.LogPersister.AppendSuccess(fmt.Sprintf("- applied manifest: %s", m.Key.ReadableString()))
	}

	e.LogPersister.AppendSuccess("Successfully rolled back PRIMARY variant")
	return nil
}

func (e *Executor) generatePrimaryManifests(namespace, commitHash string, manifests []provider.Manifest) ([]provider.Manifest, error) {
	var (
		serviceName      string
		generateService  bool
		primaryManifests = make([]provider.Manifest, 0, len(manifests)+1)
		suffix           = primaryVariant
	)

	// Apply the specified configuration if they are present.
	if sc := e.config.PrimaryVariant; sc != nil {
		var ok bool
		if sc.Suffix != "" {
			suffix = sc.Suffix
		}
		generateService = sc.Service.Create

		_, serviceName, ok = config.ParseVariantResourceReference(sc.Service.Reference)
		if !ok {
			return nil, fmt.Errorf("malformed service reference: %s", sc.Service.Reference)
		}
	}

	primaryManifests = append(primaryManifests, manifests...)

	// Find service manifests and duplicate them for PRIMARY variant.
	if generateService {
		services := findManifests(provider.KindService, serviceName, manifests)
		generatedServices, err := generateServiceManifests(services, primaryVariant, suffix)
		if err != nil {
			return nil, err
		}
		primaryManifests = append(primaryManifests, generatedServices...)
	}

	// Add labels to the generated primary manifests.
	for _, m := range primaryManifests {
		if namespace != "" {
			m.SetNamespace(namespace)
			m.Key.Namespace = namespace
		}
		m.AddAnnotations(e.builtinAnnotations(m, primaryVariant, commitHash))
	}
	return primaryManifests, nil
}
