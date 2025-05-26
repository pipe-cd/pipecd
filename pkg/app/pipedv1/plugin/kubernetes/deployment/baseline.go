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

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

func (p *Plugin) executeK8sBaselineRolloutStage(ctx context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()
	lp.Info("Start baseline rollout")

	// Get the deploy target config.
	if len(dts) == 0 {
		lp.Error("No deploy target was found")
		return sdk.StageStatusFailure
	}
	deployTargetConfig := dts[0].Config

	cfg, err := input.Request.RunningDeploymentSource.AppConfig()
	if err != nil {
		lp.Errorf("Failed while loading application config (%v)", err)
		return sdk.StageStatusFailure
	}

	var (
		appCfg          = cfg.Spec
		variantLabel    = appCfg.VariantLabel.Key
		baselineVariant = appCfg.VariantLabel.BaselineValue
	)

	var stageCfg kubeconfig.K8sBaselineRolloutStageOptions
	if err := json.Unmarshal(input.Request.StageConfig, &stageCfg); err != nil {
		lp.Errorf("Failed while unmarshalling stage config (%v)", err)
		return sdk.StageStatusFailure
	}

	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())
	loader := provider.NewLoader(toolRegistry)

	lp.Infof("Loading manifests at commit %s for handling", input.Request.RunningDeploymentSource.CommitHash)
	manifests, err := p.loadManifests(ctx, &input.Request.Deployment, appCfg, &input.Request.RunningDeploymentSource, loader)
	if err != nil {
		lp.Errorf("Failed while loading manifests (%v)", err)
		return sdk.StageStatusFailure
	}
	lp.Successf("Successfully loaded %d manifests", len(manifests))

	if len(manifests) == 0 {
		lp.Error("This application has no running Kubernetes manifests to handle")
		return sdk.StageStatusFailure
	}

	baselineManifests, err := generateBaselineManifests(appCfg, manifests, stageCfg, variantLabel, baselineVariant)
	if err != nil {
		lp.Errorf("Unable to generate manifests for BASELINE variant (%v)", err)
		return sdk.StageStatusFailure
	}

	addVariantLabelsAndAnnotations(baselineManifests, variantLabel, baselineVariant)

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

	lp.Infof("Start rolling out BASELINE variant...")
	if err := applyManifests(ctx, applier, baselineManifests, appCfg.Input.Namespace, lp); err != nil {
		return sdk.StageStatusFailure
	}

	lp.Success("Successfully rolled out BASELINE variant")
	return sdk.StageStatusSuccess
}

func (p *Plugin) executeK8sBaselineCleanStage(_ context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], _ []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
	input.Client.LogPersister().Error("Baseline clean is not yet implemented")
	return sdk.StageStatusFailure
}

func generateBaselineManifests(appCfg *kubeconfig.KubernetesApplicationSpec, manifests []provider.Manifest, stageCfg kubeconfig.K8sBaselineRolloutStageOptions, variantLabel, variant string) ([]provider.Manifest, error) {
	suffix := variant
	if stageCfg.Suffix != "" {
		suffix = stageCfg.Suffix
	}

	workloads := findWorkloadManifests(manifests, appCfg.Workloads)
	if len(workloads) == 0 {
		return nil, fmt.Errorf("unable to find any workload manifests for BASELINE variant")
	}

	var baselineManifests []provider.Manifest

	// Find service manifests and duplicate them for BASELINE variant.
	if stageCfg.CreateService {
		serviceName := appCfg.Service.Name
		services := findManifests(provider.KindService, serviceName, manifests)
		if len(services) == 0 {
			return nil, fmt.Errorf("unable to find any service for name=%q", serviceName)
		}
		// Because the loaded manifests are read-only
		// so we duplicate them to avoid updating the shared manifests data in cache.
		services = provider.DeepCopyManifests(services)

		generatedServices, err := generateVariantServiceManifests(services, variantLabel, variant, suffix)
		if err != nil {
			return nil, err
		}
		baselineManifests = append(baselineManifests, generatedServices...)
	}

	// Generate new workload manifests for VANARY variant.
	// The generated ones will mount to the new ConfigMaps and Secrets.
	replicasCalculator := func(cur *int32) int32 {
		if cur == nil {
			return 1
		}
		num := stageCfg.Replicas.Calculate(int(*cur), 1)
		return int32(num)
	}
	generatedWorkloads, err := generateVariantWorkloadManifests(workloads, nil, nil, variantLabel, variant, suffix, replicasCalculator)
	if err != nil {
		return nil, err
	}
	baselineManifests = append(baselineManifests, generatedWorkloads...)

	return baselineManifests, nil
}
