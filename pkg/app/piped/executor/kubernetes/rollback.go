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
	"strings"

	"go.uber.org/zap"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
	"github.com/pipe-cd/pipe/pkg/model"
)

type rollbackExecutor struct {
	executor.Input
}

func (e *rollbackExecutor) Execute(sig executor.StopSignal) model.StageStatus {
	var (
		ctx            = sig.Context()
		originalStatus = e.Stage.Status
		status         model.StageStatus
	)

	switch model.Stage(e.Stage.Name) {
	case model.StageRollback:
		status = e.ensureRollback(ctx)

	default:
		e.LogPersister.Errorf("Unsupported stage %s for kubernetes application", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}

	return executor.DetermineStageStatus(sig.Signal(), originalStatus, status)
}

func (e *rollbackExecutor) ensureRollback(ctx context.Context) model.StageStatus {
	// There is nothing to do if this is the first deployment.
	if e.Deployment.RunningCommitHash == "" {
		e.LogPersister.Errorf("Unable to determine the last deployed commit to rollback. It seems this is the first deployment.")
		return model.StageStatus_STAGE_FAILURE
	}

	ds, err := e.RunningDSP.Get(ctx, e.LogPersister)
	if err != nil {
		e.LogPersister.Errorf("Failed to prepare running deploy source data (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	deployCfg := ds.DeploymentConfig.KubernetesDeploymentSpec
	if deployCfg == nil {
		e.LogPersister.Error("Malformed deployment configuration: missing KubernetesDeploymentSpec")
		return model.StageStatus_STAGE_FAILURE
	}

	if deployCfg.Input.HelmChart != nil {
		chartRepoName := deployCfg.Input.HelmChart.Repository
		if chartRepoName != "" {
			deployCfg.Input.HelmChart.Insecure = e.PipedConfig.IsInsecureChartRepository(chartRepoName)
		}
	}

	p := provider.NewProvider(e.Deployment.ApplicationName, ds.AppDir, ds.RepoDir, e.Deployment.GitPath.ConfigFilename, deployCfg.Input, e.GitClient, e.Logger)
	e.Logger.Info("start executing kubernetes stage",
		zap.String("stage-name", e.Stage.Name),
		zap.String("app-dir", ds.AppDir),
	)

	// Firstly, we reapply all manifests at running commit
	// to revert PRIMARY resources and TRAFFIC ROUTING resources.

	// Load the manifests at the specified commit.
	e.LogPersister.Infof("Loading manifests at running commit %s for handling", e.Deployment.RunningCommitHash)
	manifests, err := loadManifests(ctx, e.Deployment.ApplicationId, e.Deployment.RunningCommitHash, e.AppManifestsCache, p, e.Logger)
	if err != nil {
		e.LogPersister.Errorf("Failed while loading running manifests (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.Successf("Successfully loaded %d manifests", len(manifests))

	// Because the loaded manifests are read-only
	// we duplicate them to avoid updating the shared manifests data in cache.
	manifests = duplicateManifests(manifests, "")

	// When addVariantLabelToSelector is true, ensure that all workloads
	// have the variant label in their selector.
	if deployCfg.QuickSync.AddVariantLabelToSelector {
		workloads := findWorkloadManifests(manifests, deployCfg.Workloads)
		for _, m := range workloads {
			if err := ensureVariantSelectorInWorkload(m, primaryVariant); err != nil {
				e.LogPersister.Errorf("Unable to check/set %q in selector of workload %s (%v)", variantLabel+": "+primaryVariant, m.Key.ReadableLogString(), err)
				return model.StageStatus_STAGE_FAILURE
			}
		}
	}

	// Add builtin annotations for tracking application live state.
	addBuiltinAnnontations(
		manifests,
		primaryVariant,
		e.Deployment.RunningCommitHash,
		e.PipedConfig.PipedID,
		e.Deployment.ApplicationId,
	)

	// Add config-hash annotation to the workloads.
	if err := annotateConfigHash(manifests); err != nil {
		e.LogPersister.Errorf("Unable to set %q annotation into the workload manifest (%v)", provider.AnnotationConfigHash, err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Start applying all manifests to add or update running resources.
	if err := applyManifests(ctx, p, manifests, deployCfg.Input.Namespace, e.LogPersister); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}

	var errs []error

	// Next we delete all resources of CANARY variant.
	e.LogPersister.Info("Start checking to ensure that the CANARY variant should be removed")
	if value, ok := e.MetadataStore.Shared().Get(addedCanaryResourcesMetadataKey); ok {
		resources := strings.Split(value, ",")
		if err := removeCanaryResources(ctx, p, resources, e.LogPersister); err != nil {
			errs = append(errs, err)
		}
	}

	// Then delete all resources of BASELINE variant.
	e.LogPersister.Info("Start checking to ensure that the BASELINE variant should be removed")
	if value, ok := e.MetadataStore.Shared().Get(addedBaselineResourcesMetadataKey); ok {
		resources := strings.Split(value, ",")
		if err := removeBaselineResources(ctx, p, resources, e.LogPersister); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return model.StageStatus_STAGE_FAILURE
	}
	return model.StageStatus_STAGE_SUCCESS
}
