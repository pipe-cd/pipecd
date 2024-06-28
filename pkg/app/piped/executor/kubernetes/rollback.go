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
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type rollbackExecutor struct {
	executor.Input

	appDir string
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
	case model.StageScriptRunRollback:
		status = e.ensureScriptRunRollback(ctx)
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

	appCfg := ds.ApplicationConfig.KubernetesApplicationSpec
	if appCfg == nil {
		e.LogPersister.Error("Malformed application configuration: missing KubernetesApplicationSpec")
		return model.StageStatus_STAGE_FAILURE
	}

	if appCfg.Input.HelmChart != nil {
		chartRepoName := appCfg.Input.HelmChart.Repository
		if chartRepoName != "" {
			appCfg.Input.HelmChart.Insecure = e.PipedConfig.IsInsecureChartRepository(chartRepoName)
		}
	}

	e.appDir = ds.AppDir

	loader := provider.NewLoader(e.Deployment.ApplicationName, ds.AppDir, ds.RepoDir, e.Deployment.GitPath.ConfigFilename, appCfg.Input, e.GitClient, e.Logger)
	e.Logger.Info("start executing kubernetes stage",
		zap.String("stage-name", e.Stage.Name),
		zap.String("app-dir", ds.AppDir),
	)

	// Firstly, we reapply all manifests at running commit
	// to revert PRIMARY resources and TRAFFIC ROUTING resources.

	// Load the manifests at the specified commit.
	e.LogPersister.Infof("Loading manifests at running commit %s for handling", e.Deployment.RunningCommitHash)
	manifests, err := loadManifests(
		ctx,
		e.Deployment.ApplicationId,
		e.Deployment.RunningCommitHash,
		e.AppManifestsCache,
		loader,
		e.Logger,
	)
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
	var (
		variantLabel   = appCfg.VariantLabel.Key
		primaryVariant = appCfg.VariantLabel.PrimaryValue
	)
	if appCfg.QuickSync.AddVariantLabelToSelector {
		workloads := findWorkloadManifests(manifests, appCfg.Workloads)
		for _, m := range workloads {
			if err := ensureVariantSelectorInWorkload(m, variantLabel, primaryVariant); err != nil {
				e.LogPersister.Errorf("Unable to check/set %q in selector of workload %s (%v)", variantLabel+": "+primaryVariant, m.Key.ReadableString(), err)
				return model.StageStatus_STAGE_FAILURE
			}
		}
	}

	// Add builtin annotations for tracking application live state.
	addBuiltinAnnotations(
		manifests,
		variantLabel,
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

	ag, err := newApplierGroup(e.Deployment.PlatformProvider, *appCfg, e.PipedConfig, e.Logger)
	if err != nil {
		e.LogPersister.Error(err.Error())
		return model.StageStatus_STAGE_FAILURE
	}

	// Start applying all manifests to add or update running resources.
	if err := applyManifests(ctx, ag, manifests, appCfg.Input.Namespace, e.LogPersister); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}

	var errs []error

	// Next we delete all resources of CANARY variant.
	e.LogPersister.Info("Start checking to ensure that the CANARY variant should be removed")
	if value, ok := e.MetadataStore.Shared().Get(addedCanaryResourcesMetadataKey); ok {
		resources := strings.Split(value, ",")
		if err := removeCanaryResources(ctx, ag, resources, e.LogPersister); err != nil {
			errs = append(errs, err)
		}
	}

	// Then delete all resources of BASELINE variant.
	e.LogPersister.Info("Start checking to ensure that the BASELINE variant should be removed")
	if value, ok := e.MetadataStore.Shared().Get(addedBaselineResourcesMetadataKey); ok {
		resources := strings.Split(value, ",")
		if err := removeBaselineResources(ctx, ag, resources, e.LogPersister); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return model.StageStatus_STAGE_FAILURE
	}
	return model.StageStatus_STAGE_SUCCESS
}

func (e *rollbackExecutor) ensureScriptRunRollback(ctx context.Context) model.StageStatus {
	e.LogPersister.Info("Runnnig commands for rollback...")

	onRollback, ok := e.Stage.Metadata["onRollback"]
	if !ok {
		e.LogPersister.Error("onRollback metadata is missing")
		return model.StageStatus_STAGE_FAILURE
	}

	if onRollback == "" {
		e.LogPersister.Info("No commands to run")
		return model.StageStatus_STAGE_SUCCESS
	}

	envStr, ok := e.Stage.Metadata["env"]
	env := make(map[string]string, 0)
	if ok {
		_ = json.Unmarshal([]byte(envStr), &env)
	}

	for _, v := range strings.Split(onRollback, "\n") {
		if v != "" {
			e.LogPersister.Infof("   %s", v)
		}
	}

	defaultEnvs := map[string]string{
		"DEPLOYMENT_ID":  e.Deployment.Id,
		"APPLICATION_ID": e.Deployment.ApplicationId,
	}

	envs := make([]string, 0, len(env)+len(defaultEnvs))
	for key, value := range defaultEnvs {
		envs = append(envs, "SR_"+key+"="+value)
	}
	for key, value := range env {
		envs = append(envs, key+"="+value)
	}

	cmd := exec.Command("/bin/sh", "-l", "-c", onRollback)
	cmd.Dir = e.appDir
	cmd.Env = append(os.Environ(), envs...)
	cmd.Stdout = e.LogPersister
	cmd.Stderr = e.LogPersister
	if err := cmd.Run(); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}
	return model.StageStatus_STAGE_SUCCESS
}
