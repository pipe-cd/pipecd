// Copyright 2025 The PipeCD Authors.
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

package deployment

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

func (p *Plugin) executeK8sMultiSyncStage(ctx context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()

	cfg, err := input.Request.TargetDeploymentSource.AppConfig()
	if err != nil {
		lp.Errorf("Failed while decoding application config (%v)", err.Error())
		return sdk.StageStatusFailure
	}

	type targetConfig struct {
		deployTarget *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]
		multiTarget  *kubeconfig.KubernetesMultiTarget
	}

	deployTargetMap := make(map[string]*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], 0)
	targetConfigs := make([]targetConfig, 0, len(dts))

	// prevent the deployment when its deployTarget is not found in the piped config
	for _, target := range dts {
		deployTargetMap[target.Name] = target
	}

	// If no multi-targets are specified, sync to all deploy targets.
	if len(cfg.Spec.Input.MultiTargets) == 0 {
		for _, dt := range dts {
			targetConfigs = append(targetConfigs, targetConfig{
				deployTarget: dt,
				multiTarget:  nil,
			})
		}
	} else {
		// Sync to the specified multi-targets.
		for _, multiTarget := range cfg.Spec.Input.MultiTargets {
			dt, ok := deployTargetMap[multiTarget.Target.Name]
			if !ok {
				lp.Infof("Ignore multi target '%s': not matched any deployTarget", multiTarget.Target.Name)
				continue
			}

			targetConfigs = append(targetConfigs, targetConfig{
				deployTarget: dt,
				multiTarget:  &multiTarget,
			})
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, tc := range targetConfigs {
		// Start syncing the deployment to the target.
		eg.Go(func() error {
			lp.Infof("Start syncing the deployment to the target %s", tc.deployTarget.Name)
			status := p.sync(ctx, input, tc.deployTarget, tc.multiTarget)
			if status == sdk.StageStatusFailure {
				return fmt.Errorf("failed to sync the deployment to the target %s", tc.deployTarget.Name)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		lp.Errorf("Failed while syncing the deployment (%v)", err)
		return sdk.StageStatusFailure
	}

	return sdk.StageStatusSuccess
}

func (p *Plugin) sync(
	ctx context.Context,
	input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec],
	dt *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig],
	multiTarget *kubeconfig.KubernetesMultiTarget,
) sdk.StageStatus {
	lp := input.Client.LogPersister()
	lp.Info("Start syncing the deployment")

	cfg, err := input.Request.TargetDeploymentSource.AppConfig()
	if err != nil {
		lp.Errorf("Failed while loading application config (%v)", err)
		return sdk.StageStatusFailure
	}

	var stageCfg kubeconfig.K8sSyncStageOptions
	if len(input.Request.StageConfig) > 0 {
		// TODO: this is a temporary solution to support the stage options specified under "with"
		// When the stage options under "with" are empty, we cannot detect whether the stage is a quick sync stage or not.
		// So we have to add a new field to the sdk.ExecuteStageRequest or sdk.Deployment to indicate that the deployment is a quick sync strategy or in a pipeline sync strategy.
		if err := json.Unmarshal(input.Request.StageConfig, &stageCfg); err != nil {
			lp.Errorf("Failed while unmarshalling stage config (%v)", err)
			return sdk.StageStatusFailure
		}
	} else {
		stageCfg = cfg.Spec.QuickSync
	}

	// TODO: find the way to hold the tool registry and loader in the plugin.
	// Currently, we create them every time the stage is executed beucause we can't pass input.Client.toolRegistry to the plugin when starting the plugin.
	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())
	loader := provider.NewLoader(toolRegistry)

	lp.Infof("Loading manifests at commit %s for handling", input.Request.TargetDeploymentSource.CommitHash)
	manifests, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.TargetDeploymentSource, loader, multiTarget)
	if err != nil {
		lp.Errorf("Failed while loading manifests (%v)", err)
		return sdk.StageStatusFailure
	}
	lp.Successf("Successfully loaded %d manifests", len(manifests))

	// Because the loaded manifests are read-only
	// we duplicate them to avoid updating the shared manifests data in cache.
	manifests = provider.DeepCopyManifests(manifests)

	// When addVariantLabelToSelector is true, ensure that all workloads
	// have the variant label in their selector.
	var (
		variantLabel   = cfg.Spec.VariantLabel.Key
		primaryVariant = cfg.Spec.VariantLabel.PrimaryValue
	)

	if stageCfg.AddVariantLabelToSelector {
		workloads := findWorkloadManifests(manifests, cfg.Spec.Workloads)
		for _, m := range workloads {
			if err := ensureVariantSelectorInWorkload(m, variantLabel, primaryVariant); err != nil {
				lp.Errorf("Unable to check/set %q in selector of workload %s (%v)", variantLabel+": "+primaryVariant, m.Key().ReadableString(), err)
				return sdk.StageStatusFailure
			}
		}
	}

	addVariantLabelsAndAnnotations(manifests, variantLabel, primaryVariant)

	if err := annotateConfigHash(manifests); err != nil {
		lp.Errorf("Unable to set %q annotation into the workload manifest (%v)", provider.AnnotationConfigHash, err)
		return sdk.StageStatusFailure
	}

	deployTargetConfig := dt.Config

	// Get the kubectl tool path.
	kubectlPath, err := toolRegistry.Kubectl(ctx, cmp.Or(cfg.Spec.Input.KubectlVersion, deployTargetConfig.KubectlVersion))
	if err != nil {
		lp.Errorf("Failed while getting kubectl tool (%v)", err)
		return sdk.StageStatusFailure
	}

	// Create the kubectl wrapper for the target cluster.
	kubectl := provider.NewKubectl(kubectlPath)

	// Create the applier for the target cluster.
	applier := provider.NewApplier(kubectl, cfg.Spec.Input, deployTargetConfig, input.Logger)

	// Start applying all manifests to add or update running resources.
	// TODO: use applyManifests instead of applyManifestsSDK
	if err := applyManifests(ctx, applier, manifests, cfg.Spec.Input.Namespace, lp); err != nil {
		lp.Errorf("Failed while applying manifests (%v)", err)
		return sdk.StageStatusSuccess
	}

	if !stageCfg.Prune {
		lp.Info("Resource GC was skipped because sync.prune was not configured")
		return sdk.StageStatusSuccess
	}

	// Wait for all applied manifests to be stable.
	// In theory, we don't need to wait for them to be stable before going to the next step
	// but waiting for a while reduces the number of Kubernetes changes in a short time.
	lp.Info("Waiting for the applied manifests to be stable")
	select {
	case <-time.After(15 * time.Second):
		break
	case <-ctx.Done():
		break
	}

	lp.Info("Start finding all running resources but no longer defined in Git")

	namespacedLiveResources, clusterScopedLiveResources, err := provider.GetLiveResources(ctx, kubectl, deployTargetConfig.KubeConfigPath, input.Request.Deployment.ApplicationID)
	if err != nil {
		lp.Errorf("Failed while getting live resources (%v)", err)
		return sdk.StageStatusFailure
	}

	if len(namespacedLiveResources)+len(clusterScopedLiveResources) == 0 {
		lp.Info("There is no data about live resource so no resource will be removed")
		return sdk.StageStatusSuccess
	}

	lp.Successf("Successfully loaded %d live resources", len(namespacedLiveResources)+len(clusterScopedLiveResources))

	removeKeys := provider.FindRemoveResources(manifests, namespacedLiveResources, clusterScopedLiveResources)
	if len(removeKeys) == 0 {
		lp.Info("There are no live resources should be removed")
		return sdk.StageStatusSuccess
	}

	lp.Infof("Start pruning %d resources", len(removeKeys))
	var deletedCount int
	for _, key := range removeKeys {
		if err := kubectl.Delete(ctx, deployTargetConfig.KubeConfigPath, key.Namespace(), key); err != nil {
			if errors.Is(err, provider.ErrNotFound) {
				lp.Infof("Specified resource does not exist, so skip deleting the resource: %s (%v)", key.ReadableString(), err)
				continue
			}
			lp.Errorf("Failed while deleting resource %s (%v)", key.ReadableString(), err)
			continue // continue to delete other resources
		}
		deletedCount++
		lp.Successf("- deleted resource: %s", key.ReadableString())
	}

	lp.Successf("Successfully deleted %d resources", deletedCount)

	return sdk.StageStatusSuccess
}
