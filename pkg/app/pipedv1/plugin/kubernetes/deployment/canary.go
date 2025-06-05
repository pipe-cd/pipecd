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

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/deployment/yamlprocessor"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/toolregistry"
)

func (p *Plugin) executeK8sCanaryRolloutStage(ctx context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()
	lp.Info("Start canary rollout")

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
		appCfg        = cfg.Spec
		variantLabel  = appCfg.VariantLabel.Key
		canaryVariant = appCfg.VariantLabel.CanaryValue
	)

	var stageCfg kubeconfig.K8sCanaryRolloutStageOptions
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

	if len(manifests) == 0 {
		lp.Error("This application has no Kubernetes manifests to handle")
		return sdk.StageStatusFailure
	}

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
	canaryManifests, err := p.generateCanaryManifests(appCfg, manifests, stageCfg, variantLabel, canaryVariant)
	if err != nil {
		lp.Errorf("Unable to generate manifests for CANARY variant (%v)", err)
		return sdk.StageStatusFailure
	}

	addVariantLabelsAndAnnotations(canaryManifests, variantLabel, canaryVariant)

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

	// Start rolling out the resources for CANARY variant.
	lp.Info("Start rolling out CANARY variant...")
	if err := applyManifests(ctx, applier, canaryManifests, appCfg.Input.Namespace, lp); err != nil {
		return sdk.StageStatusFailure
	}

	lp.Success("Successfully rolled out CANARY variant")
	return sdk.StageStatusSuccess
}

func (p *Plugin) generateCanaryManifests(appCfg *kubeconfig.KubernetesApplicationSpec, manifests []provider.Manifest, opts kubeconfig.K8sCanaryRolloutStageOptions, variantLabel, variant string) ([]provider.Manifest, error) {
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
		// Because the loaded manifests are read-only
		// so we duplicate them to avoid updating the shared manifests data in cache.
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
	// The generated ones will mount to the new ConfigMaps and Secrets.
	replicasCalculator := func(cur *int32) int32 {
		if cur == nil {
			return 1
		}
		num := opts.Replicas.Calculate(int(*cur), 1)
		return int32(num)
	}
	// We don't need to duplicate the workload manifests
	// because generateVariantWorkloadManifests function is already making a duplicate while decoding.
	// workloads = duplicateManifests(workloads, suffix)
	generatedWorkloads, err := generateVariantWorkloadManifests(workloads, configMaps, secrets, variantLabel, variant, suffix, replicasCalculator)
	if err != nil {
		return nil, err
	}
	canaryManifests = append(canaryManifests, generatedWorkloads...)

	return canaryManifests, nil
}

func (p *Plugin) executeK8sCanaryCleanStage(ctx context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()
	lp.Info("Start canary clean")

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
		appCfg        = cfg.Spec
		variantLabel  = appCfg.VariantLabel.Key
		canaryVariant = appCfg.VariantLabel.CanaryValue
	)

	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())

	// Get the kubectl tool path.
	kubectlPath, err := toolRegistry.Kubectl(ctx, cmp.Or(appCfg.Input.KubectlVersion, deployTargetConfig.KubectlVersion))
	if err != nil {
		lp.Errorf("Failed while getting kubectl tool (%v)", err)
		return sdk.StageStatusFailure
	}

	// Create the kubectl wrapper for the target cluster.
	kubectl := provider.NewKubectl(kubectlPath)

	// Create the applier for the target cluster.
	applier := provider.NewApplier(kubectl, appCfg.Input, deployTargetConfig, input.Logger)

	if err := deleteVariantResources(ctx, lp, kubectl, deployTargetConfig.KubeConfigPath, applier, input.Request.Deployment.ApplicationID, variantLabel, canaryVariant); err != nil {
		lp.Errorf("Unable to remove canary resources: (%v)", err)
		return sdk.StageStatusFailure
	}

	lp.Success("Successfully cleaned CANARY variant")
	return sdk.StageStatusSuccess
}

func findConfigMapManifests(manifests []provider.Manifest) []provider.Manifest {
	out := make([]provider.Manifest, 0, len(manifests))
	for _, m := range manifests {
		if !m.IsConfigMap() {
			continue
		}
		out = append(out, m)
	}
	return out
}

func findSecretManifests(manifests []provider.Manifest) []provider.Manifest {
	out := make([]provider.Manifest, 0, len(manifests))
	for _, m := range manifests {
		if !m.IsSecret() {
			continue
		}
		out = append(out, m)
	}
	return out
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
		p, err := yamlprocessor.NewProcessor(bytes)
		if err != nil {
			return nil, err
		}

		for _, o := range patch.Ops {
			switch o.Op {
			case kubeconfig.K8sResourcePatchOpYAMLReplace:
				if err := p.ReplaceString(o.Path, o.Value); err != nil {
					return nil, fmt.Errorf("failed to replace value at path: %s, error: %w", o.Path, err)
				}
			default:
				// TODO: Support more patch operation for K8sCanaryRolloutStageOptions.
				return nil, fmt.Errorf("%s operation is not supported currently", o.Op)
			}
		}

		return p.Bytes(), nil
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

	// When the target is the whole manifest,
	// just pass full bytes to process and build a new manifest based on the returned data.
	root := patch.Target.DocumentRoot
	if root == "" {
		out, err := process(fullBytes)
		if err != nil {
			return nil, err
		}
		return buildManifest(out)
	}

	// When the target is a manifest field specified by documentRoot,
	// we have to extract that field value as a string.
	p, err := yamlprocessor.NewProcessor(fullBytes)
	if err != nil {
		return nil, err
	}

	v, err := p.GetValue(root)
	if err != nil {
		return nil, err
	}
	sv, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("the value for the specified root %s must be a string", root)
	}

	// And process that field data.
	out, err := process([]byte(sv))
	if err != nil {
		return nil, err
	}

	// Then rewrite the new data into the specified root.
	if err := p.ReplaceString(root, string(out)); err != nil {
		return nil, err
	}

	return buildManifest(p.Bytes())
}
