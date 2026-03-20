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
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/toolregistry"
)

func (p *Plugin) executeK8sMultiPrimaryRolloutStage(ctx context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()

	cfg, err := input.Request.TargetDeploymentSource.AppConfig()
	if err != nil {
		lp.Errorf("Failed while decoding application config (%v)", err.Error())
		return sdk.StageStatusFailure
	}

	var stageCfg kubeconfig.K8sPrimaryRolloutStageOptions
	if len(input.Request.StageConfig) > 0 {
		if err := json.Unmarshal(input.Request.StageConfig, &stageCfg); err != nil {
			lp.Errorf("Failed while unmarshalling stage config (%v)", err)
			return sdk.StageStatusFailure
		}
	}

	type targetConfig struct {
		deployTarget *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]
		multiTarget  *kubeconfig.KubernetesMultiTarget
	}

	deployTargetMap := make(map[string]*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig])
	targetConfigs := make([]targetConfig, 0, len(dts))

	for _, target := range dts {
		deployTargetMap[target.Name] = target
	}

	// If no multi-targets are specified, roll out primary to all deploy targets.
	if len(cfg.Spec.Input.MultiTargets) == 0 {
		for _, dt := range dts {
			targetConfigs = append(targetConfigs, targetConfig{
				deployTarget: dt,
				multiTarget:  nil,
			})
		}
	} else {
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
		eg.Go(func() error {
			lp.Infof("Start primary rollout for target %s", tc.deployTarget.Name)
			status := p.primaryRollout(ctx, input, tc.deployTarget, tc.multiTarget, stageCfg)
			if status == sdk.StageStatusFailure {
				return fmt.Errorf("failed to primary rollout for target %s", tc.deployTarget.Name)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		lp.Errorf("Failed while rolling out primary (%v)", err)
		return sdk.StageStatusFailure
	}

	return sdk.StageStatusSuccess
}

func (p *Plugin) primaryRollout(
	ctx context.Context,
	input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec],
	dt *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig],
	multiTarget *kubeconfig.KubernetesMultiTarget,
	stageCfg kubeconfig.K8sPrimaryRolloutStageOptions,
) sdk.StageStatus {
	lp := input.Client.LogPersister()

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

	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())
	loader := provider.NewLoader(toolRegistry)

	lp.Infof("Loading manifests at commit %s for handling", input.Request.TargetDeploymentSource.CommitHash)
	manifests, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.TargetDeploymentSource, loader, input.Logger, multiTarget)
	if err != nil {
		lp.Errorf("Failed while loading manifests (%v)", err)
		return sdk.StageStatusFailure
	}
	lp.Successf("Successfully loaded %d manifests", len(manifests))

	if len(manifests) == 0 {
		lp.Error("This application has no Kubernetes manifests to handle")
		return sdk.StageStatusFailure
	}

	// Generate the manifests for applying.
	lp.Info("Start generating manifests for PRIMARY variant")
	primaryManifests, err := generatePrimaryManifests(appCfg, manifests, stageCfg, variantLabel, primaryVariant)
	if err != nil {
		lp.Errorf("Unable to generate manifests for PRIMARY variant (%v)", err)
		return sdk.StageStatusFailure
	}
	lp.Successf("Successfully generated %d manifests for PRIMARY variant", len(primaryManifests))

	addVariantLabelsAndAnnotations(primaryManifests, variantLabel, primaryVariant)

	if err := annotateConfigHash(primaryManifests); err != nil {
		lp.Errorf("Unable to set %q annotation into the workload manifest (%v)", provider.AnnotationConfigHash, err)
		return sdk.StageStatusFailure
	}

	deployTargetConfig := dt.Config

	// Resolve kubectl version: multiTarget > spec > deployTarget
	kubectlVersion := cmp.Or(appCfg.Input.KubectlVersion, deployTargetConfig.KubectlVersion)
	if multiTarget != nil {
		kubectlVersion = cmp.Or(multiTarget.KubectlVersion, kubectlVersion)
	}

	kubectlPath, err := toolRegistry.Kubectl(ctx, kubectlVersion)
	if err != nil {
		lp.Errorf("Failed while getting kubectl tool (%v)", err)
		return sdk.StageStatusFailure
	}

	kubectl := provider.NewKubectl(kubectlPath)
	applier := provider.NewApplier(kubectl, appCfg.Input, deployTargetConfig, input.Logger)

	lp.Info("Start rolling out PRIMARY variant...")
	if err := applyManifests(ctx, applier, primaryManifests, appCfg.Input.Namespace, lp); err != nil {
		lp.Errorf("Failed while applying manifests (%v)", err)
		return sdk.StageStatusFailure
	}

	if !stageCfg.Prune {
		lp.Info("Resource GC was skipped because prune was not configured")
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
	deletedCount := deleteResources(ctx, lp, applier, removeKeys)
	lp.Successf("Successfully deleted %d resources", deletedCount)

	return sdk.StageStatusSuccess
}

// generatePrimaryManifests generates manifests for the PRIMARY variant.
// It deep-copies the input manifests, adds the variant label to workload selectors
// if requested, and generates a variant Service manifest if requested.
func generatePrimaryManifests(appCfg *kubeconfig.KubernetesApplicationSpec, manifests []provider.Manifest, stageCfg kubeconfig.K8sPrimaryRolloutStageOptions, variantLabel, variant string) ([]provider.Manifest, error) {
	suffix := variant
	if stageCfg.Suffix != "" {
		suffix = stageCfg.Suffix
	}

	primaryManifests := provider.DeepCopyManifests(manifests)

	// Add the variant label to workload selectors if requested.
	if stageCfg.AddVariantLabelToSelector {
		workloads := findWorkloadManifests(primaryManifests, nil)
		for _, m := range workloads {
			if err := ensureVariantSelectorInWorkload(m, variantLabel, variant); err != nil {
				return nil, fmt.Errorf("unable to check/set %q in selector of workload %s (%w)", variantLabel+": "+variant, m.Key().ReadableString(), err)
			}
		}
	}

	// Generate Service manifests for the PRIMARY variant if requested.
	if stageCfg.CreateService {
		serviceName := appCfg.Service.Name
		services := findManifests(provider.KindService, serviceName, primaryManifests)
		if len(services) == 0 {
			return nil, fmt.Errorf("unable to find any service for PRIMARY variant")
		}
		// Deep-copy the services to avoid mutating the shared primaryManifests slice entries.
		services = provider.DeepCopyManifests(services)

		generatedServices, err := generateVariantServiceManifests(services, variantLabel, variant, suffix)
		if err != nil {
			return nil, fmt.Errorf("failed to generate service manifests: %w", err)
		}
		primaryManifests = append(primaryManifests, generatedServices...)
	}

	return primaryManifests, nil
}
