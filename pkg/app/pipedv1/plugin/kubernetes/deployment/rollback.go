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

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/toolregistry"
)

func (p *Plugin) executeK8sRollbackStage(ctx context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()

	if input.Request.RunningDeploymentSource.CommitHash == "" {
		lp.Errorf("Unable to determine the last deployed commit to rollback. It seems this is the first deployment.")
		return sdk.StageStatusFailure
	}

	lp.Info("Start rolling back the deployment")

	cfg, err := input.Request.RunningDeploymentSource.AppConfig()
	if err != nil {
		lp.Errorf("Failed while loading application config (%v)", err)
		return sdk.StageStatusFailure
	}

	lp.Infof("Loading manifests at commit %s for handling", input.Request.RunningDeploymentSource.CommitHash)
	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())
	manifests, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.RunningDeploymentSource, provider.NewLoader(toolRegistry), input.Logger)
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
		variantLabel    = cfg.Spec.VariantLabel.Key
		primaryVariant  = cfg.Spec.VariantLabel.PrimaryValue
		baselineVariant = cfg.Spec.VariantLabel.BaselineValue
		canaryVariant   = cfg.Spec.VariantLabel.CanaryValue
	)
	// TODO: Consider other fields to configure whether to add a variant label to the selector
	// because the rollback stage is executed in both quick sync and pipeline sync strategies.
	if cfg.Spec.QuickSync.AddVariantLabelToSelector {
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

	// Get the deploy target config.
	if len(dts) == 0 {
		lp.Error("No deploy target was found")
		return sdk.StageStatusFailure
	}
	deployTargetConfig := dts[0].Config

	// Get the kubectl tool path.
	kubectlPath, err := toolRegistry.Kubectl(ctx, cmp.Or(cfg.Spec.Input.KubectlVersion, deployTargetConfig.KubectlVersion))
	if err != nil {
		lp.Errorf("Failed while getting kubectl tool (%v)", err)
		return sdk.StageStatusFailure
	}

	kubectl := provider.NewKubectl(kubectlPath)

	// Create the applier for the target cluster.
	applier := provider.NewApplier(kubectl, cfg.Spec.Input, deployTargetConfig, input.Logger)

	// Start applying all manifests to add or update running resources.
	if err := applyManifests(ctx, applier, manifests, cfg.Spec.Input.Namespace, lp); err != nil {
		lp.Errorf("Failed while applying manifests (%v)", err)
		return sdk.StageStatusFailure
	}

	var failed bool

	// TODO: prune resources which doesn't exist in the running manifests but exists in the target manifests.
	// This occurs when the user adds a new resource and failed the deployment pipeline.
	// This feature is not implemented in pipedv0, but it's nice to have it in this plugin.

	lp.Info("Start removing CANARY variant resources if exists")
	if err := deleteVariantResources(ctx, lp, kubectl, deployTargetConfig.KubeConfigPath, applier, input.Request.Deployment.ApplicationID, variantLabel, canaryVariant); err != nil {
		lp.Errorf("Failed while deleting variant resources (%v)", err)
		failed = true
	}

	lp.Info("Start removing BASELINE variant resources if exists")
	if err := deleteVariantResources(ctx, lp, kubectl, deployTargetConfig.KubeConfigPath, applier, input.Request.Deployment.ApplicationID, variantLabel, baselineVariant); err != nil {
		lp.Errorf("Failed while deleting variant resources (%v)", err)
		failed = true
	}

	if failed {
		return sdk.StageStatusFailure
	}

	return sdk.StageStatusSuccess
}
