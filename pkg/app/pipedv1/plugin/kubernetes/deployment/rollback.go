// Copyright 2024 The PipeCD Authors.
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

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
)

func (a *DeploymentService) executeK8sRollbackStage(ctx context.Context, lp logpersister.StageLogPersister, input *deployment.ExecutePluginInput) model.StageStatus {
	if input.GetDeployment().GetRunningCommitHash() == "" {
		lp.Errorf("Unable to determine the last deployed commit to rollback. It seems this is the first deployment.")
		return model.StageStatus_STAGE_FAILURE
	}

	lp.Info("Start rolling back the deployment")

	cfg, err := config.DecodeYAML[*kubeconfig.KubernetesApplicationSpec](input.GetRunningDeploymentSource().GetApplicationConfig())
	if err != nil {
		lp.Errorf("Failed while decoding application config (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	lp.Infof("Loading manifests at commit %s for handling", input.GetDeployment().GetRunningCommitHash())

	// TODO: consider multiple multiTargets
	manifests, err := a.loadManifests(ctx, input.GetDeployment(), cfg.Spec, input.GetRunningDeploymentSource(), kubeconfig.KubernetesMultiTarget{})
	if err != nil {
		lp.Errorf("Failed while loading manifests (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	lp.Successf("Successfully loaded %d manifests", len(manifests))

	// Because the loaded manifests are read-only
	// we duplicate them to avoid updating the shared manifests data in cache.
	// TODO: implement duplicateManifests function

	// When addVariantLabelToSelector is true, ensure that all workloads
	// have the variant label in their selector.
	var (
		variantLabel   = cfg.Spec.VariantLabel.Key
		primaryVariant = cfg.Spec.VariantLabel.PrimaryValue
	)
	// TODO: handle cfg.Spec.QuickSync.AddVariantLabelToSelector

	// Add variant annotations to all manifests.
	for i := range manifests {
		manifests[i].AddAnnotations(map[string]string{
			variantLabel: primaryVariant,
		})
	}

	if err := annotateConfigHash(manifests); err != nil {
		lp.Errorf("Unable to set %q annotation into the workload manifest (%v)", provider.AnnotationConfigHash, err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Get the deploy target config.
	targets, err := input.GetDeployment().GetDeployTargets(a.pluginConfig.Name)
	if err != nil {
		lp.Errorf("Failed while finding deploy target config (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	deployTargetConfig, err := kubeconfig.FindDeployTarget(a.pluginConfig, targets[0]) // TODO: consider multiple targets
	if err != nil {
		lp.Errorf("Failed while unmarshalling deploy target config (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Get the kubectl tool path.
	kubectlPath, err := a.toolRegistry.Kubectl(ctx, cmp.Or(cfg.Spec.Input.KubectlVersion, deployTargetConfig.KubectlVersion, defaultKubectlVersion))
	if err != nil {
		lp.Errorf("Failed while getting kubectl tool (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Create the applier for the target cluster.
	applier := provider.NewApplier(provider.NewKubectl(kubectlPath), cfg.Spec.Input, deployTargetConfig, a.logger)

	// Start applying all manifests to add or update running resources.
	if err := applyManifests(ctx, applier, manifests, cfg.Spec.Input.Namespace, lp); err != nil {
		lp.Errorf("Failed while applying manifests (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	// TODO: implement prune resources
	// TODO: delete all resources of CANARY variant
	// TODO: delete all resources of BASELINE variant

	return model.StageStatus_STAGE_SUCCESS
}
