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
	"fmt"

	"golang.org/x/sync/errgroup"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/toolregistry"
)

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

	// Resolve which deploy targets to operate on.
	// If no multiTargets are configured at app level, operate on all deploy targets.
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
