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

	"golang.org/x/sync/errgroup"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/toolregistry"
)

func (p *Plugin) executeK8sMultiBaselineRolloutStage(ctx context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()

	cfg, err := input.Request.TargetDeploymentSource.AppConfig()
	if err != nil {
		lp.Errorf("Failed while decoding application config (%v)", err.Error())
		return sdk.StageStatusFailure
	}

	var stageCfg kubeconfig.K8sBaselineRolloutStageOptions
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

	deployTargetMap := make(map[string]*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], 0)
	targetConfigs := make([]targetConfig, 0, len(dts))

	for _, target := range dts {
		deployTargetMap[target.Name] = target
	}

	// If no multi-targets are specified, roll out baseline to all deploy targets.
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
			lp.Infof("Start baseline rollout for target %s", tc.deployTarget.Name)
			status := p.baselineRollout(ctx, input, tc.deployTarget, tc.multiTarget, stageCfg)
			if status == sdk.StageStatusFailure {
				return fmt.Errorf("failed to baseline rollout for target %s", tc.deployTarget.Name)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		lp.Errorf("Failed while rolling out baseline (%v)", err)
		return sdk.StageStatusFailure
	}

	return sdk.StageStatusSuccess
}

func (p *Plugin) baselineRollout(
	ctx context.Context,
	input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec],
	dt *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig],
	multiTarget *kubeconfig.KubernetesMultiTarget,
	stageCfg kubeconfig.K8sBaselineRolloutStageOptions,
) sdk.StageStatus {
	lp := input.Client.LogPersister()

	cfg, err := input.Request.TargetDeploymentSource.AppConfig()
	if err != nil {
		lp.Errorf("Failed while loading application config (%v)", err)
		return sdk.StageStatusFailure
	}

	var (
		appCfg          = cfg.Spec
		variantLabel    = appCfg.VariantLabel.Key
		baselineVariant = appCfg.VariantLabel.BaselineValue
	)

	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())
	loader := provider.NewLoader(toolRegistry)

	// Baseline uses the RUNNING deployment source (current live version), not the target.
	lp.Infof("Loading manifests at commit %s for handling", input.Request.RunningDeploymentSource.CommitHash)
	manifests, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.RunningDeploymentSource, loader, input.Logger, multiTarget)
	if err != nil {
		lp.Errorf("Failed while loading manifests (%v)", err)
		return sdk.StageStatusFailure
	}
	lp.Successf("Successfully loaded %d manifests", len(manifests))

	if len(manifests) == 0 {
		lp.Error("This application has no Kubernetes manifests to handle")
		return sdk.StageStatusFailure
	}

	// Because the loaded manifests are read-only
	// we duplicate them to avoid updating the shared manifests data in cache.
	manifests = provider.DeepCopyManifests(manifests)

	// Find and generate workload & service manifests for BASELINE variant.
	baselineManifests, err := generateBaselineManifests(appCfg, manifests, stageCfg, variantLabel, baselineVariant)
	if err != nil {
		lp.Errorf("Unable to generate manifests for BASELINE variant (%v)", err)
		return sdk.StageStatusFailure
	}

	addVariantLabelsAndAnnotations(baselineManifests, variantLabel, baselineVariant)

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

	lp.Info("Start rolling out BASELINE variant...")
	if err := applyManifests(ctx, applier, baselineManifests, appCfg.Input.Namespace, lp); err != nil {
		lp.Errorf("Failed while applying baseline manifests (%v)", err)
		return sdk.StageStatusFailure
	}

	lp.Success("Successfully rolled out BASELINE variant")
	return sdk.StageStatusSuccess
}

func generateBaselineManifests(appCfg *kubeconfig.KubernetesApplicationSpec, manifests []provider.Manifest, opts kubeconfig.K8sBaselineRolloutStageOptions, variantLabel, variant string) ([]provider.Manifest, error) {
	suffix := variant
	if opts.Suffix != "" {
		suffix = opts.Suffix
	}

	workloads := findWorkloadManifests(manifests, appCfg.Workloads)
	if len(workloads) == 0 {
		return nil, fmt.Errorf("unable to find any workload manifests for BASELINE variant")
	}

	var baselineManifests []provider.Manifest

	// Find service manifests and duplicate them for BASELINE variant.
	if opts.CreateService {
		serviceName := appCfg.Service.Name
		services := findManifests(provider.KindService, serviceName, manifests)
		if len(services) == 0 {
			return nil, fmt.Errorf("unable to find any service for name=%q", serviceName)
		}
		// Duplicate them to avoid updating the shared manifests data in cache.
		services = duplicateManifests(services, "")

		generatedServices, err := generateVariantServiceManifests(services, variantLabel, variant, suffix)
		if err != nil {
			return nil, err
		}
		baselineManifests = append(baselineManifests, generatedServices...)
	}

	// Find config map manifests and duplicate them for BASELINE variant.
	configMaps := findConfigMapManifests(manifests)
	baselineConfigMaps := duplicateManifests(configMaps, suffix)
	baselineManifests = append(baselineManifests, baselineConfigMaps...)

	// Find secret manifests and duplicate them for BASELINE variant.
	secrets := findSecretManifests(manifests)
	baselineSecrets := duplicateManifests(secrets, suffix)
	baselineManifests = append(baselineManifests, baselineSecrets...)

	// Generate new workload manifests for BASELINE variant.
	replicasCalculator := func(cur *int32) int32 {
		if cur == nil {
			return 1
		}
		num := opts.Replicas.Calculate(int(*cur), 1)
		return int32(num)
	}
	generatedWorkloads, err := generateVariantWorkloadManifests(workloads, configMaps, secrets, variantLabel, variant, suffix, replicasCalculator)
	if err != nil {
		return nil, err
	}
	baselineManifests = append(baselineManifests, generatedWorkloads...)

	return baselineManifests, nil
}

func (p *Plugin) executeK8sMultiBaselineCleanStage(ctx context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()

	cfg, err := input.Request.TargetDeploymentSource.AppConfig()
	if err != nil {
		lp.Errorf("Failed while decoding application config (%v)", err)
		return sdk.StageStatusFailure
	}

	deployTargetMap := make(map[string]*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], len(dts))
	for _, dt := range dts {
		deployTargetMap[dt.Name] = dt
	}

	type targetConfig struct {
		deployTarget *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]
		multiTarget  *kubeconfig.KubernetesMultiTarget
	}

	targetConfigs := make([]targetConfig, 0, len(dts))
	if len(cfg.Spec.Input.MultiTargets) == 0 {
		for _, dt := range dts {
			targetConfigs = append(targetConfigs, targetConfig{deployTarget: dt})
		}
	} else {
		for _, mt := range cfg.Spec.Input.MultiTargets {
			dt, ok := deployTargetMap[mt.Target.Name]
			if !ok {
				lp.Infof("Ignore multi target '%s': not matched any deployTarget", mt.Target.Name)
				continue
			}
			targetConfigs = append(targetConfigs, targetConfig{deployTarget: dt, multiTarget: &mt})
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, tc := range targetConfigs {
		eg.Go(func() error {
			lp.Infof("Start cleaning BASELINE variant on target %s", tc.deployTarget.Name)
			if err := p.baselineClean(ctx, input, tc.deployTarget, tc.multiTarget, cfg); err != nil {
				return fmt.Errorf("failed to clean BASELINE variant on target %s: %w", tc.deployTarget.Name, err)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		lp.Errorf("Failed while cleaning BASELINE variant (%v)", err)
		return sdk.StageStatusFailure
	}

	return sdk.StageStatusSuccess
}

func (p *Plugin) baselineClean(
	ctx context.Context,
	input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec],
	dt *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig],
	multiTarget *kubeconfig.KubernetesMultiTarget,
	cfg *sdk.ApplicationConfig[kubeconfig.KubernetesApplicationSpec],
) error {
	lp := input.Client.LogPersister()

	var (
		appCfg          = cfg.Spec
		variantLabel    = appCfg.VariantLabel.Key
		baselineVariant = appCfg.VariantLabel.BaselineValue
	)

	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())

	// Resolve kubectl version: multiTarget > spec > deployTarget
	kubectlVersion := cmp.Or(appCfg.Input.KubectlVersion, dt.Config.KubectlVersion)
	if multiTarget != nil {
		kubectlVersion = cmp.Or(multiTarget.KubectlVersion, kubectlVersion)
	}

	kubectlPath, err := toolRegistry.Kubectl(ctx, kubectlVersion)
	if err != nil {
		return fmt.Errorf("failed while getting kubectl tool: %w", err)
	}

	kubectl := provider.NewKubectl(kubectlPath)
	applier := provider.NewApplier(kubectl, appCfg.Input, dt.Config, input.Logger)

	if err := deleteVariantResources(ctx, lp, kubectl, dt.Config.KubeConfigPath, applier, input.Request.Deployment.ApplicationID, variantLabel, baselineVariant); err != nil {
		return fmt.Errorf("unable to remove baseline resources: %w", err)
	}

	return nil
}
