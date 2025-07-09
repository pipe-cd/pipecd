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
	"context"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
)

func (p *Plugin) executeK8sTrafficRoutingStage(ctx context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()
	lp.Info("Start routing the traffic")

	cfg, err := input.Request.TargetDeploymentSource.AppConfig()
	if err != nil {
		lp.Errorf("Failed while loading application config (%v)", err)
		return sdk.StageStatusFailure
	}

	switch kubeconfig.DetermineKubernetesTrafficRoutingMethod(cfg.Spec.TrafficRouting) {
	case kubeconfig.KubernetesTrafficRoutingMethodPodSelector:
		return p.executeK8sTrafficRoutingStagePodSelector(ctx, input, dts, cfg)
	case kubeconfig.KubernetesTrafficRoutingMethodIstio:
		return p.executeK8sTrafficRoutingStageIstio(ctx, input, dts, cfg)
	default:
		lp.Errorf("Unknown traffic routing method: %s", cfg.Spec.TrafficRouting.Method)
		return sdk.StageStatusFailure
	}
}

func (p *Plugin) executeK8sTrafficRoutingStagePodSelector(_ context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], _ []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], _ *sdk.ApplicationConfig[kubeconfig.KubernetesApplicationSpec]) sdk.StageStatus {
	lp := input.Client.LogPersister()
	lp.Error("Traffic routing by PodSelector is not yet implemented")
	return sdk.StageStatusFailure
}

func (p *Plugin) executeK8sTrafficRoutingStageIstio(_ context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], _ []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], _ *sdk.ApplicationConfig[kubeconfig.KubernetesApplicationSpec]) sdk.StageStatus {
	lp := input.Client.LogPersister()
	lp.Error("Traffic routing by Istio is not yet implemented")
	return sdk.StageStatusFailure
}
