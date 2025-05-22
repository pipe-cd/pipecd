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

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

func (p *Plugin) executeK8sPrimaryRolloutStage(ctx context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()
	lp.Info("Start primary rollout")

	// Get the deploy target config.
	if len(dts) == 0 {
		lp.Error("No deploy target was found")
		return sdk.StageStatusFailure
	}
	deployTargetConfig := dts[0].Config

	cfg, err := input.Request.TargetDeploymentSource.AppConfig()
	if err != nil {
		lp.Errorf("Failed while loading application config (%v)", err)
		return sdk.StageStatusFailure
	}

	var (
		appCfg         = cfg.Spec
		variantLabel   = appCfg.VariantLabel.Key
		primaryVariant = appCfg.VariantLabel.PrimaryValue
	)

	var stageCfg kubeconfig.K8sPrimaryRolloutStageOptions
	if err := json.Unmarshal(input.Request.StageConfig, &stageCfg); err != nil {
		lp.Errorf("Failed while unmarshalling stage config (%v)", err)
		return sdk.StageStatusFailure
	}

	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())
	loader := provider.NewLoader(toolRegistry)

	lp.Infof("Loading manifests at commit %s for handling", input.Request.TargetDeploymentSource.CommitHash)
	manifests, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.TargetDeploymentSource, loader)
	if err != nil {
		lp.Errorf("Failed while loading manifests (%v)", err)
		return sdk.StageStatusFailure
	}
	lp.Successf("Successfully loaded %d manifests", len(manifests))

	var primaryManifests []provider.Manifest
	routingMethod := kubeconfig.DetermineKubernetesTrafficRoutingMethod(appCfg.TrafficRouting)
	switch routingMethod {
	case kubeconfig.KubernetesTrafficRoutingMethodPodSelector:
		primaryManifests = manifests
	case kubeconfig.KubernetesTrafficRoutingMethodIstio:
		// TODO: support routing by Istio
		lp.Errorf("Traffic routing method %v is not yet implemented", routingMethod)
		return sdk.StageStatusFailure
	default:
		lp.Errorf("Traffic routing method %v is not supported", routingMethod)
		return sdk.StageStatusFailure
	}

	// Check if the variant selector is in the workloads.
	if !stageCfg.AddVariantLabelToSelector &&
		routingMethod == kubeconfig.KubernetesTrafficRoutingMethodPodSelector &&
		cfg.HasStage(StageK8sTrafficRouting) {
		workloads := findWorkloadManifests(primaryManifests, appCfg.Workloads)
		var invalid bool
		for _, m := range workloads {
			if err := checkVariantSelectorInWorkload(m, variantLabel, primaryVariant); err != nil {
				invalid = true
			}
		}
		if invalid {
			lp.Errorf("Missing %q in selector of workload", variantLabel+": "+primaryVariant)
			return sdk.StageStatusFailure
		}
	}

	// Generate the manifests for applying.
	lp.Infof("Start generating manifests for PRIMARY variant")
	if primaryManifests, err = generatePrimaryManifests(primaryManifests, stageCfg, variantLabel, primaryVariant); err != nil {
		lp.Errorf("Unable to generate manifests for PRIMARY variant (%v)", err)
		return sdk.StageStatusFailure
	}
	lp.Successf("Successfully generated %d manifests for PRIMARY variant", len(primaryManifests))

	addVariantLabelsAndAnnotations(primaryManifests, variantLabel, primaryVariant)

	if err := annotateConfigHash(primaryManifests); err != nil {
		lp.Errorf("Unable to set %q annotation into the workload manifest (%v)", provider.AnnotationConfigHash, err)
		return sdk.StageStatusFailure
	}

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
	if err := applyManifests(ctx, applier, primaryManifests, cfg.Spec.Input.Namespace, lp); err != nil {
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

	// Find the running resources that are not defined in Git.
	lp.Info("Start finding all running PRIMARY resources but no longer defined in Git")
	namespacedLiveResources, clusterScopedLiveResources, err := provider.GetLiveResources(ctx, kubectl, deployTargetConfig.KubeConfigPath, input.Request.Deployment.ApplicationID, fmt.Sprintf("%s=%s", variantLabel, primaryVariant))
	if err != nil {
		lp.Errorf("Failed while getting live resources (%v)", err)
		return sdk.StageStatusFailure
	}

	if len(namespacedLiveResources)+len(clusterScopedLiveResources) == 0 {
		lp.Info("There is no data about live resource so no resource will be removed")
		return sdk.StageStatusSuccess
	}

	lp.Successf("Successfully loaded %d live resources", len(namespacedLiveResources)+len(clusterScopedLiveResources))

	removeKeys := provider.FindRemoveResources(primaryManifests, namespacedLiveResources, clusterScopedLiveResources)
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

// generatePrimaryManifests generates manifests for the PRIMARY variant.
// It duplicates the input manifests, adds the variant label to workloads if needed,
// and generates Service manifests with a name suffix and variant selector if requested.
func generatePrimaryManifests(manifests []provider.Manifest, stageCfg kubeconfig.K8sPrimaryRolloutStageOptions, variantLabel, variant string) ([]provider.Manifest, error) {
	suffix := variant
	if stageCfg.Suffix != "" {
		suffix = stageCfg.Suffix
	}

	primaryManifests := provider.DeepCopyManifests(manifests)

	// Add the variant label to workload selectors if requested.
	if stageCfg.AddVariantLabelToSelector {
		workloads := findWorkloadManifests(primaryManifests, nil) // All Deployments if refs is nil.
		for _, m := range workloads {
			if err := ensureVariantSelectorInWorkload(m, variantLabel, variant); err != nil {
				return nil, fmt.Errorf("unable to check/set %q in selector of workload %s (%w)", variantLabel+": "+variant, m.Key().ReadableString(), err)
			}
		}
	}

	// Generate Service manifests for the PRIMARY variant if requested.
	if stageCfg.CreateService {
		services := findManifests("Service", "", primaryManifests)
		if len(services) == 0 {
			return nil, fmt.Errorf("unable to find any service for PRIMARY variant")
		}
		generatedServices, err := generateVariantServiceManifests(services, variantLabel, variant, suffix)
		if err != nil {
			return nil, fmt.Errorf("failed to generate service manifests: %w", err)
		}
		primaryManifests = append(primaryManifests, generatedServices...)
	}

	return primaryManifests, nil
}
