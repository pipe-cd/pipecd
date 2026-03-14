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
	"github.com/pipe-cd/pipecd/pkg/yamlprocessor"
)

func (p *Plugin) executeK8sMultiCanaryRolloutStage(ctx context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()

	cfg, err := input.Request.TargetDeploymentSource.AppConfig()
	if err != nil {
		lp.Errorf("Failed while decoding application config (%v)", err.Error())
		return sdk.StageStatusFailure
	}

	var stageCfg kubeconfig.K8sCanaryRolloutStageOptions
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

	// If no multi-targets are specified, roll out canary to all deploy targets.
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
			lp.Infof("Start canary rollout for target %s", tc.deployTarget.Name)
			status := p.canaryRollout(ctx, input, tc.deployTarget, tc.multiTarget, stageCfg)
			if status == sdk.StageStatusFailure {
				return fmt.Errorf("failed to canary rollout for target %s", tc.deployTarget.Name)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		lp.Errorf("Failed while rolling out canary (%v)", err)
		return sdk.StageStatusFailure
	}

	return sdk.StageStatusSuccess
}

func (p *Plugin) canaryRollout(
	ctx context.Context,
	input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec],
	dt *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig],
	multiTarget *kubeconfig.KubernetesMultiTarget,
	stageCfg kubeconfig.K8sCanaryRolloutStageOptions,
) sdk.StageStatus {
	lp := input.Client.LogPersister()

	cfg, err := input.Request.TargetDeploymentSource.AppConfig()
	if err != nil {
		lp.Errorf("Failed while loading application config (%v)", err)
		return sdk.StageStatusFailure
	}

	var (
		appCfg        = cfg.Spec
		variantLabel  = appCfg.VariantLabel.Key
		canaryVariant = appCfg.VariantLabel.CanaryValue
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

	// Because the loaded manifests are read-only
	// we duplicate them to avoid updating the shared manifests data in cache.
	manifests = provider.DeepCopyManifests(manifests)

	// Patches the manifests if needed.
	if len(stageCfg.Patches) > 0 {
		lp.Info("Patching manifests before generating for CANARY variant")
		manifests, err = patchManifests(manifests, stageCfg.Patches, patchManifest)
		if err != nil {
			lp.Errorf("Failed while patching manifests (%v)", err)
			return sdk.StageStatusFailure
		}
	}

	// Find and generate workload & service manifests for CANARY variant.
	canaryManifests, err := generateCanaryManifests(appCfg, manifests, stageCfg, variantLabel, canaryVariant)
	if err != nil {
		lp.Errorf("Unable to generate manifests for CANARY variant (%v)", err)
		return sdk.StageStatusFailure
	}

	addVariantLabelsAndAnnotations(canaryManifests, variantLabel, canaryVariant)

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

	lp.Info("Start rolling out CANARY variant...")
	if err := applyManifests(ctx, applier, canaryManifests, appCfg.Input.Namespace, lp); err != nil {
		lp.Errorf("Failed while applying canary manifests (%v)", err)
		return sdk.StageStatusFailure
	}

	lp.Success("Successfully rolled out CANARY variant")
	return sdk.StageStatusSuccess
}

func generateCanaryManifests(appCfg *kubeconfig.KubernetesApplicationSpec, manifests []provider.Manifest, opts kubeconfig.K8sCanaryRolloutStageOptions, variantLabel, variant string) ([]provider.Manifest, error) {
	suffix := variant
	if opts.Suffix != "" {
		suffix = opts.Suffix
	}

	workloads := findWorkloadManifests(manifests, appCfg.Workloads)
	if len(workloads) == 0 {
		return nil, fmt.Errorf("unable to find any workload manifests for CANARY variant")
	}

	var canaryManifests []provider.Manifest

	// Find service manifests and duplicate them for CANARY variant.
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
		canaryManifests = append(canaryManifests, generatedServices...)
	}

	// Find config map manifests and duplicate them for CANARY variant.
	configMaps := findConfigMapManifests(manifests)
	canaryConfigMaps := duplicateManifests(configMaps, suffix)
	canaryManifests = append(canaryManifests, canaryConfigMaps...)

	// Find secret manifests and duplicate them for CANARY variant.
	secrets := findSecretManifests(manifests)
	canarySecrets := duplicateManifests(secrets, suffix)
	canaryManifests = append(canaryManifests, canarySecrets...)

	// Generate new workload manifests for CANARY variant.
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
	canaryManifests = append(canaryManifests, generatedWorkloads...)

	return canaryManifests, nil
}

type patcher func(m provider.Manifest, cfg kubeconfig.K8sResourcePatch) (*provider.Manifest, error)

func patchManifests(manifests []provider.Manifest, patches []kubeconfig.K8sResourcePatch, patcher patcher) ([]provider.Manifest, error) {
	if len(patches) == 0 {
		return manifests, nil
	}

	out := make([]provider.Manifest, len(manifests))
	copy(out, manifests)

	for _, p := range patches {
		target := -1
		for i, m := range out {
			if m.Key().Kind() != p.Target.Kind {
				continue
			}
			if m.Key().Name() != p.Target.Name {
				continue
			}
			target = i
			break
		}
		if target < 0 {
			return nil, fmt.Errorf("no manifest matches the given patch: kind=%s, name=%s", p.Target.Kind, p.Target.Name)
		}
		patched, err := patcher(out[target], p)
		if err != nil {
			return nil, fmt.Errorf("failed to patch manifest: %s, error: %w", out[target].Key(), err)
		}
		out[target] = *patched
	}

	return out, nil
}

func patchManifest(m provider.Manifest, patch kubeconfig.K8sResourcePatch) (*provider.Manifest, error) {
	if len(patch.Ops) == 0 {
		return &m, nil
	}

	fullBytes, err := m.YamlBytes()
	if err != nil {
		return nil, err
	}

	process := func(bytes []byte) ([]byte, error) {
		proc, err := yamlprocessor.NewProcessor(bytes)
		if err != nil {
			return nil, err
		}

		for _, o := range patch.Ops {
			switch o.Op {
			case kubeconfig.K8sResourcePatchOpYAMLReplace:
				if err := proc.ReplaceString(o.Path, o.Value); err != nil {
					return nil, fmt.Errorf("failed to replace value at path: %s, error: %w", o.Path, err)
				}
			default:
				return nil, fmt.Errorf("%s operation is not supported currently", o.Op)
			}
		}

		return proc.Bytes(), nil
	}

	buildManifest := func(bytes []byte) (*provider.Manifest, error) {
		manifests, err := provider.ParseManifests(string(bytes))
		if err != nil {
			return nil, err
		}
		if len(manifests) != 1 {
			return nil, fmt.Errorf("unexpected number of manifests, expected 1, got %d", len(manifests))
		}
		return &manifests[0], nil
	}

	root := patch.Target.DocumentRoot
	if root == "" {
		out, err := process(fullBytes)
		if err != nil {
			return nil, err
		}
		return buildManifest(out)
	}

	proc, err := yamlprocessor.NewProcessor(fullBytes)
	if err != nil {
		return nil, err
	}

	v, err := proc.GetValue(root)
	if err != nil {
		return nil, err
	}
	sv, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("the value for the specified root %s must be a string", root)
	}

	out, err := process([]byte(sv))
	if err != nil {
		return nil, err
	}

	if err := proc.ReplaceString(root, string(out)); err != nil {
		return nil, err
	}

	return buildManifest(proc.Bytes())
}

func (p *Plugin) executeK8sMultiCanaryCleanStage(ctx context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
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
			targetConfigs = append(targetConfigs, targetConfig{deployTarget: dt})
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, tc := range targetConfigs {
		eg.Go(func() error {
			lp.Infof("Start cleaning CANARY variant on target %s", tc.deployTarget.Name)
			if err := p.canaryClean(ctx, input, tc.deployTarget, cfg); err != nil {
				return fmt.Errorf("failed to clean CANARY variant on target %s: %w", tc.deployTarget.Name, err)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		lp.Errorf("Failed while cleaning CANARY variant (%v)", err)
		return sdk.StageStatusFailure
	}

	return sdk.StageStatusSuccess
}

func (p *Plugin) canaryClean(
	ctx context.Context,
	input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec],
	dt *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig],
	cfg *sdk.ApplicationConfig[kubeconfig.KubernetesApplicationSpec],
) error {
	lp := input.Client.LogPersister()

	var (
		appCfg        = cfg.Spec
		variantLabel  = appCfg.VariantLabel.Key
		canaryVariant = appCfg.VariantLabel.CanaryValue
	)

	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())

	kubectlPath, err := toolRegistry.Kubectl(ctx, cmp.Or(appCfg.Input.KubectlVersion, dt.Config.KubectlVersion))
	if err != nil {
		return fmt.Errorf("failed while getting kubectl tool: %w", err)
	}

	kubectl := provider.NewKubectl(kubectlPath)
	applier := provider.NewApplier(kubectl, appCfg.Input, dt.Config, input.Logger)

	if err := deleteVariantResources(ctx, lp, kubectl, dt.Config.KubeConfigPath, applier, input.Request.Deployment.ApplicationID, variantLabel, canaryVariant); err != nil {
		return fmt.Errorf("unable to remove canary resources: %w", err)
	}

	return nil
}
