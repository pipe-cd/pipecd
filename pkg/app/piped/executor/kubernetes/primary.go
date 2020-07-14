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

func (e *Executor) ensurePrimaryRollout(ctx context.Context) model.StageStatus {
	commitHash := e.Deployment.Trigger.Commit.Hash

	// Load the manifests at the triggered commit.
	e.LogPersister.AppendInfof("Loading manifests at commit %s for handling", commitHash)
	manifests, err := e.loadManifests(ctx)
	if err != nil {
		e.LogPersister.AppendErrorf("Failed while loading manifests (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.AppendSuccessf("Successfully loaded %d manifests", len(manifests))

	// Generate the manifests for applying.
	e.LogPersister.AppendInfo("Start generating manifests for PRIMARY variant")
	primaryManifests, err := e.generatePrimaryManifests(commitHash, manifests)
	if err != nil {
		e.LogPersister.AppendErrorf("Unable to generate manifests for PRIMARY variant (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.AppendSuccessf("Successfully generated %d manifests for PRIMARY variant", len(primaryManifests))

	// Start applying all manifests to add or update running resources.
	e.LogPersister.AppendInfof("Start applying %d primary resources", len(primaryManifests))
	for _, m := range primaryManifests {
		if err = e.provider.ApplyManifest(ctx, m); err != nil {
			e.LogPersister.AppendErrorf("Failed to apply manifest: %s (%v)", m.Key.ReadableString(), err)
			return model.StageStatus_STAGE_FAILURE
		}
		e.LogPersister.AppendSuccessf("- applied manifest: %s", m.Key.ReadableString())
	}
	e.LogPersister.AppendSuccessf("Successfully applied %d primary resources", len(primaryManifests))

	// TODO: Wait for all applied manifests to be ready.
	e.LogPersister.AppendInfo("Waiting for the applied manifests to be ready")

	// TODO: Find and remove the running resources that are not defined in Git.
	e.LogPersister.AppendInfo("Start finding and removing all running PRIMARY resources but not in Git")

	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) generatePrimaryManifests(commitHash string, manifests []provider.Manifest) ([]provider.Manifest, error) {
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

	for _, m := range manifests {
		// Because the loaded maninests are read-only
		// so we duplicate them to avoid updating the shared manifests data in cache.
		primaryManifests = append(primaryManifests, m.Duplicate(m.Key.Name))
	}

	// Find service manifests and duplicate them for PRIMARY variant.
	if generateService {
		services := findManifests(provider.KindService, serviceName, manifests)
		if len(services) == 0 {
			return nil, fmt.Errorf("unable to find any service for name=%q", serviceName)
		}

		// Because the loaded maninests are read-only
		// so we duplicate them to avoid updating the shared manifests data in cache.
		duplicates := make([]provider.Manifest, 0, len(services))
		for _, m := range services {
			duplicates = append(duplicates, m.Duplicate(m.Key.Name))
		}

		generatedServices, err := generateServiceManifests(duplicates, primaryVariant, suffix)
		if err != nil {
			return nil, err
		}
		primaryManifests = append(primaryManifests, generatedServices...)
	}

	// TODO: Find out traffic-routing manfiests and keep them as the previously configured routing.

	// Add predefined annotations to the generated manifests.
	for _, m := range primaryManifests {
		m.AddAnnotations(e.builtinAnnotations(m, primaryVariant, commitHash))
	}
	return primaryManifests, nil
}
