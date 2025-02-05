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

package deployment

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
)

func (a *DeploymentService) executeK8sSyncStage(ctx context.Context, lp logpersister.StageLogPersister, input *deployment.ExecutePluginInput) model.StageStatus {
	eg, ctx := errgroup.WithContext(ctx)

	cfg, err := config.DecodeYAML[*kubeconfig.KubernetesApplicationSpec](input.GetTargetDeploymentSource().GetApplicationConfig())
	if err != nil {
		lp.Errorf("Failed while decoding application config (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	type targetConfig struct {
		deployTarget kubeconfig.KubernetesDeployTargetConfig
		multiTarget  *kubeconfig.KubernetesMultiTarget
	}
	deployTargetMap := make(map[string]kubeconfig.KubernetesDeployTargetConfig, 0)
	targetConfigs := make([]targetConfig, 0, len(input.GetDeployment().GetDeployTargets()))

	// prevent the deployment when its deployTarget is not found in the piped config
	for _, target := range input.GetDeployment().GetDeployTargets() {
		dt, err := kubeconfig.FindDeployTarget(a.pluginConfig, target)
		if err != nil {
			lp.Errorf("Failed while finding deploy target (%v)", err)
			return model.StageStatus_STAGE_FAILURE
		}

		deployTargetMap[dt.Name] = dt
	}

	// If no multi-targets are specified, sync to all deploy targets.
	if len(cfg.Spec.Input.MultiTargets) == 0 {
		for _, dt := range deployTargetMap {
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

	for _, tc := range targetConfigs {
		// Start syncing the deployment to the target.
		eg.Go(func() error {
			lp.Infof("Start syncing the deployment to the target %s", tc.deployTarget.Name)
			status := a.sync(ctx, lp, input, tc.deployTarget, tc.multiTarget)
			if status != model.StageStatus_STAGE_SUCCESS {
				return fmt.Errorf("failed to sync the deployment to the target %s", tc.deployTarget.Name)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		lp.Errorf("Failed while syncing the deployment (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	return model.StageStatus_STAGE_SUCCESS
}

func (a *DeploymentService) sync(ctx context.Context, lp logpersister.StageLogPersister, input *deployment.ExecutePluginInput, deployTargetConfig kubeconfig.KubernetesDeployTargetConfig, multiTarget *kubeconfig.KubernetesMultiTarget) model.StageStatus {
	lp.Infof("Start syncing the deployment")

	cfg, err := config.DecodeYAML[*kubeconfig.KubernetesApplicationSpec](input.GetTargetDeploymentSource().GetApplicationConfig())
	if err != nil {
		lp.Errorf("Failed while decoding application config (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	lp.Infof("Loading manifests at commit %s for handling", input.GetDeployment().GetTrigger().GetCommit().GetHash())
	manifests, err := a.loadManifests(ctx, input.GetDeployment(), cfg.Spec, input.GetTargetDeploymentSource(), multiTarget)
	if err != nil {
		lp.Errorf("Failed while loading manifests (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	lp.Successf("Successfully loaded %d manifests", len(manifests))

	// Because the loaded manifests are read-only
	// we duplicate them to avoid updating the shared manifests data in cache.
	// TODO: implement duplicateManifests function

	// When addVariantLabelToSelector is true, ensure that all workloads
	// have the variant label in their selector.
	var (
		variantLabel   = cfg.Spec.VariantLabel.Key
		primaryVariant = cfg.Spec.VariantLabel.PrimaryValue
	)
	// TODO: handle cfg.Spec.QuickSync.AddVariantLabelToSelector

	// Add variant annotations to all manifests.
	for i := range manifests {
		manifests[i].AddLabels(map[string]string{
			variantLabel: primaryVariant,
		})
		manifests[i].AddAnnotations(map[string]string{
			variantLabel: primaryVariant,
		})
	}

	if err := annotateConfigHash(manifests); err != nil {
		lp.Errorf("Unable to set %q annotation into the workload manifest (%v)", provider.AnnotationConfigHash, err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Get the deploy target config.
	targets, err := input.GetDeployment().GetDeployTargets(a.pluginConfig.Name)
	if err != nil {
		lp.Errorf("Failed while finding deploy target config (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	deployTargetConfig, err := kubeconfig.FindDeployTarget(a.pluginConfig, targets[0]) // TODO: consider multiple targets
	if err != nil {
		lp.Errorf("Failed while unmarshalling deploy target config (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Get the kubectl tool path.
	kubectlVersions := []string{cfg.Spec.Input.KubectlVersion, deployTargetConfig.KubectlVersion, defaultKubectlVersion}
	// If multi-target is specified, use the kubectl version specified in it.
	if multiTarget != nil {
		kubectlVersions = append([]string{multiTarget.KubectlVersion}, kubectlVersions...)
	}
	kubectlPath, err := a.toolRegistry.Kubectl(ctx, cmp.Or(kubectlVersions...))
	if err != nil {
		lp.Errorf("Failed while getting kubectl tool (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Create the kubectl wrapper for the target cluster.
	kubectl := provider.NewKubectl(kubectlPath)

	// Create the applier for the target cluster.
	applier := provider.NewApplier(kubectl, cfg.Spec.Input, deployTargetConfig, a.logger)

	// Start applying all manifests to add or update running resources.
	if err := applyManifests(ctx, applier, manifests, cfg.Spec.Input.Namespace, lp); err != nil {
		lp.Errorf("Failed while applying manifests (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	// TODO: treat the stage options specified under "with"
	if !cfg.Spec.QuickSync.Prune {
		lp.Info("Resource GC was skipped because sync.prune was not configured")
		return model.StageStatus_STAGE_SUCCESS
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

	namespacedLiveResources, clusterScopedLiveResources, err := provider.GetLiveResources(ctx, kubectl, deployTargetConfig.KubeConfigPath, input.GetDeployment().GetApplicationId())
	if err != nil {
		lp.Errorf("Failed while getting live resources (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	if len(namespacedLiveResources)+len(clusterScopedLiveResources) == 0 {
		lp.Info("There is no data about live resource so no resource will be removed")
		return model.StageStatus_STAGE_SUCCESS
	}

	lp.Successf("Successfully loaded %d live resources", len(namespacedLiveResources)+len(clusterScopedLiveResources))

	removeKeys := provider.FindRemoveResources(manifests, namespacedLiveResources, clusterScopedLiveResources)
	if len(removeKeys) == 0 {
		lp.Info("There are no live resources should be removed")
		return model.StageStatus_STAGE_SUCCESS
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
	return model.StageStatus_STAGE_SUCCESS
}
