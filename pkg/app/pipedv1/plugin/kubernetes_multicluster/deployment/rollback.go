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
	"sync"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/toolregistry"
)

func (p *Plugin) executeK8sMultiRollbackStage(ctx context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()

	cfg, err := input.Request.RunningDeploymentSource.AppConfig()
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

	// If no multi-targets are specified, rollback all deploy targets.
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

	type result struct {
		target string
		status sdk.StageStatus
	}

	results := make([]result, 0, len(targetConfigs))
	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	for _, tc := range targetConfigs {
		wg.Add(1)
		go func() {
			lp.Infof("Start rollbacking the deployment for the target %s", tc.deployTarget.Name)

			status := p.rollback(ctx, input, tc.deployTarget, tc.multiTarget)
			mu.Lock()
			results = append(results, result{
				target: tc.deployTarget.Name,
				status: status,
			})
			mu.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()

	finalStatus := sdk.StageStatusFailure
	for _, result := range results {
		// success at least one of the rollback succeed
		if result.status != sdk.StageStatusFailure {
			finalStatus = sdk.StageStatusSuccess
		}

		lp.Infof("target %s, status: %s", result.target, result.status.String())
	}

	return finalStatus
}

func (p *Plugin) rollback(ctx context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], dt *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], multiTarget *kubeconfig.KubernetesMultiTarget) sdk.StageStatus {
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
	// TODO: consider multiTarget later
	manifests, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.RunningDeploymentSource, provider.NewLoader(toolRegistry), input.Logger, multiTarget)
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
	deployTargetConfig := dt.Config

	kubectlVersions := []string{cfg.Spec.Input.KubectlVersion, deployTargetConfig.KubectlVersion}
	// If multi-target is specified, use the kubectl version specified in it.
	if multiTarget != nil {
		kubectlVersions = append([]string{multiTarget.KubectlVersion}, kubectlVersions...)
	}

	// Get the kubectl tool path.
	kubectlPath, err := toolRegistry.Kubectl(ctx, cmp.Or(kubectlVersions...))
	if err != nil {
		lp.Errorf("Failed while getting kubectl tool (%v)", err)
		return sdk.StageStatusFailure
	}

	// Create the applier for the target cluster.
	applier := provider.NewApplier(provider.NewKubectl(kubectlPath), cfg.Spec.Input, deployTargetConfig, input.Logger)

	// Start applying all manifests to add or update running resources.
	if err := applyManifests(ctx, applier, manifests, cfg.Spec.Input.Namespace, lp); err != nil {
		lp.Errorf("Failed while applying manifests (%v)", err)
		return sdk.StageStatusFailure
	}

	// TODO: implement prune resources
	// TODO: delete all resources of CANARY variant
	// TODO: delete all resources of BASELINE variant

	return sdk.StageStatusSuccess
}
