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

package planpreview

import (
	"context"
	"fmt"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/piped-plugin-sdk-go/diff"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/toolregistry"
)

var (
	_ sdk.PlanPreviewPlugin[kubeconfig.KubernetesPluginConfig, kubeconfig.KubernetesDeployTargetConfig, kubeconfig.KubernetesApplicationSpec] = (*Plugin)(nil)
)

type Plugin struct{}

func (p *Plugin) GetPlanPreview(ctx context.Context, _ *kubeconfig.KubernetesPluginConfig, dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], input *sdk.GetPlanPreviewInput[kubeconfig.KubernetesApplicationSpec]) (*sdk.GetPlanPreviewResponse, error) {
	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())
	loader := provider.NewLoader(toolRegistry)

	var (
		oldManifests []provider.Manifest
		newManifests []provider.Manifest
		err          error
	)

	// load manifests from target ds
	targetDS := input.Request.TargetDeploymentSource
	targetAppCfg, err := targetDS.AppConfig()
	if err != nil {
		return nil, err
	}
	tagetSpec := targetAppCfg.Spec

	newManifests, err = loader.LoadManifests(ctx, provider.LoaderInput{
		PipedID:          input.Request.PipedID,
		AppID:            input.Request.ApplicationID,
		CommitHash:       targetDS.CommitHash,
		AppName:          input.Request.ApplicationName,
		AppDir:           targetDS.ApplicationDirectory,
		ConfigFilename:   targetDS.ApplicationConfigFilename,
		Manifests:        tagetSpec.Input.Manifests,
		Namespace:        tagetSpec.Input.Namespace,
		KustomizeVersion: tagetSpec.Input.KustomizeVersion,
		KustomizeOptions: tagetSpec.Input.KustomizeOptions,
		HelmVersion:      tagetSpec.Input.HelmVersion,
		HelmChart:        tagetSpec.Input.HelmChart,
		HelmOptions:      tagetSpec.Input.HelmOptions,
		Logger:           input.Logger,
	})
	if err != nil {
		return nil, err
	}

	runningDS := input.Request.RunningDeploymentSource
	if runningDS.CommitHash != "" {
		runningAppCfg, err := runningDS.AppConfig()
		if err != nil {
			return nil, err
		}
		runningSpec := runningAppCfg.Spec
		oldManifests, err = loader.LoadManifests(ctx, provider.LoaderInput{
			PipedID:          input.Request.PipedID,
			AppID:            input.Request.ApplicationID,
			CommitHash:       runningDS.CommitHash,
			AppName:          input.Request.ApplicationName,
			AppDir:           runningDS.ApplicationDirectory,
			ConfigFilename:   runningDS.ApplicationConfigFilename,
			Manifests:        runningSpec.Input.Manifests,
			Namespace:        runningSpec.Input.Namespace,
			KustomizeVersion: runningSpec.Input.KustomizeVersion,
			KustomizeOptions: runningSpec.Input.KustomizeOptions,
			HelmVersion:      runningSpec.Input.HelmVersion,
			HelmChart:        runningSpec.Input.HelmChart,
			HelmOptions:      runningSpec.Input.HelmOptions,
			Logger:           input.Logger,
		})
		if err != nil {
			return nil, err
		}
	}

	// diff
	result, err := provider.DiffList(
		oldManifests,
		newManifests,
		input.Logger,
		diff.WithEquateEmpty(),
		diff.WithCompareNumberAndNumericString(),
	)
	if err != nil {
		return nil, err
	}

	return toResponse(result, dts[0].Name), nil
}

func toResponse(result *provider.DiffListResult, deployTarget string) *sdk.GetPlanPreviewResponse {
	if result.NoChanges() {
		return &sdk.GetPlanPreviewResponse{
			Results: []sdk.PlanPreviewResult{
				{
					DeployTarget: deployTarget,
					NoChange:     true,
					Summary:      "No changes were detected",
					DiffLanguage: "diff",
				},
			},
		}
	}

	details := result.Render(provider.DiffRenderOptions{
		MaskSecret:     true,
		UseDiffCommand: true,
	})

	// return result
	return &sdk.GetPlanPreviewResponse{
		Results: []sdk.PlanPreviewResult{
			{
				DeployTarget: deployTarget,
				NoChange:     false,
				Summary:      fmt.Sprintf("%d added manifests, %d changed manifests, %d deleted manifests", len(result.Adds), len(result.Changes), len(result.Deletes)),
				DiffLanguage: "diff",
				Details:      []byte(details),
			},
		},
	}
}
