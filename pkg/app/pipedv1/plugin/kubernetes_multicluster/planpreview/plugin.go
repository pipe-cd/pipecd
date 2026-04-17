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

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/toolregistry"
)

var (
	_ sdk.PlanPreviewPlugin[sdk.ConfigNone, kubeconfig.KubernetesDeployTargetConfig, kubeconfig.KubernetesApplicationSpec] = (*Plugin)(nil)
)

// Plugin implements the sdk.PlanPreviewPlugin interface for the kubernetes_multicluster plugin.
type Plugin struct{}

// GetPlanPreview returns the plan preview result showing what will change across all deploy targets.
func (p *Plugin) GetPlanPreview(ctx context.Context, _ *sdk.ConfigNone, dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], input *sdk.GetPlanPreviewInput[kubeconfig.KubernetesApplicationSpec]) (*sdk.GetPlanPreviewResponse, error) {
	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())
	loader := provider.NewLoader(toolRegistry)

	targetDS := input.Request.TargetDeploymentSource
	targetAppCfg, err := targetDS.AppConfig()
	if err != nil {
		return nil, err
	}
	targetSpec := targetAppCfg.Spec

	runningDS := input.Request.RunningDeploymentSource

	multiTargets := targetSpec.Input.MultiTargets

	// Single-target fallback: no multiTargets configured — load manifests once and return one result.
	if len(multiTargets) == 0 {
		newManifests, err := loadManifests(ctx, loader, input, &targetDS, targetSpec, nil)
		if err != nil {
			return nil, err
		}

		var oldManifests []provider.Manifest
		if runningDS.CommitHash != "" {
			runningAppCfg, err := runningDS.AppConfig()
			if err != nil {
				return nil, err
			}
			oldManifests, err = loadManifests(ctx, loader, input, &runningDS, runningAppCfg.Spec, nil)
			if err != nil {
				return nil, err
			}
		}

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

		deployTargetName := ""
		if len(dts) > 0 {
			deployTargetName = dts[0].Name
		}
		return &sdk.GetPlanPreviewResponse{
			Results: []sdk.PlanPreviewResult{toResult(result, deployTargetName)},
		}, nil
	}

	// Multi-target: produce one PlanPreviewResult per deploy target.
	results := make([]sdk.PlanPreviewResult, 0, len(dts))
	for _, dt := range dts {
		// Find the matching KubernetesMultiTarget config for this deploy target.
		var mt *kubeconfig.KubernetesMultiTarget
		for i := range multiTargets {
			if multiTargets[i].Target.Name == dt.Name {
				mt = &multiTargets[i]
				break
			}
		}

		newManifests, err := loadManifests(ctx, loader, input, &targetDS, targetSpec, mt)
		if err != nil {
			results = append(results, sdk.PlanPreviewResult{
				DeployTarget: dt.Name,
				NoChange:     false,
				Summary:      fmt.Sprintf("Failed to load target manifests: %v", err),
				DiffLanguage: "diff",
			})
			continue
		}

		var oldManifests []provider.Manifest
		if runningDS.CommitHash != "" {
			runningAppCfg, err := runningDS.AppConfig()
			if err != nil {
				return nil, err
			}
			oldManifests, err = loadManifests(ctx, loader, input, &runningDS, runningAppCfg.Spec, mt)
			if err != nil {
				results = append(results, sdk.PlanPreviewResult{
					DeployTarget: dt.Name,
					NoChange:     false,
					Summary:      fmt.Sprintf("Failed to load running manifests: %v", err),
					DiffLanguage: "diff",
				})
				continue
			}
		}

		result, err := provider.DiffList(
			oldManifests,
			newManifests,
			input.Logger,
			diff.WithEquateEmpty(),
			diff.WithCompareNumberAndNumericString(),
		)
		if err != nil {
			results = append(results, sdk.PlanPreviewResult{
				DeployTarget: dt.Name,
				NoChange:     false,
				Summary:      fmt.Sprintf("Failed to diff manifests: %v", err),
				DiffLanguage: "diff",
			})
			continue
		}

		results = append(results, toResult(result, dt.Name))
	}

	return &sdk.GetPlanPreviewResponse{Results: results}, nil
}

// loadManifests loads manifests from the given deployment source, optionally overriding
// the manifest paths from the multiTarget config.
func loadManifests(ctx context.Context, loader *provider.Loader, input *sdk.GetPlanPreviewInput[kubeconfig.KubernetesApplicationSpec], ds *sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec], spec *kubeconfig.KubernetesApplicationSpec, mt *kubeconfig.KubernetesMultiTarget) ([]provider.Manifest, error) {
	manifestPaths := spec.Input.Manifests
	if mt != nil && len(mt.Manifests) > 0 {
		manifestPaths = mt.Manifests
	}

	return loader.LoadManifests(ctx, provider.LoaderInput{
		PipedID:          input.Request.PipedID,
		AppID:            input.Request.ApplicationID,
		CommitHash:       ds.CommitHash,
		AppName:          input.Request.ApplicationName,
		AppDir:           ds.ApplicationDirectory,
		ConfigFilename:   ds.ApplicationConfigFilename,
		Manifests:        manifestPaths,
		Namespace:        spec.Input.Namespace,
		KustomizeVersion: spec.Input.KustomizeVersion,
		KustomizeOptions: spec.Input.KustomizeOptions,
		HelmVersion:      spec.Input.HelmVersion,
		HelmChart:        spec.Input.HelmChart,
		HelmOptions:      spec.Input.HelmOptions,
		Logger:           input.Logger,
	})
}

// toResult converts a DiffListResult into a PlanPreviewResult for the given deploy target.
func toResult(result *provider.DiffListResult, deployTarget string) sdk.PlanPreviewResult {
	if result.NoChanges() {
		return sdk.PlanPreviewResult{
			DeployTarget: deployTarget,
			NoChange:     true,
			Summary:      "No changes were detected",
			DiffLanguage: "diff",
		}
	}

	details := result.Render(provider.DiffRenderOptions{
		MaskSecret: true,
	})

	return sdk.PlanPreviewResult{
		DeployTarget: deployTarget,
		NoChange:     false,
		Summary:      fmt.Sprintf("%d added manifests, %d changed manifests, %d deleted manifests", len(result.Adds), len(result.Changes), len(result.Deletes)),
		DiffLanguage: "diff",
		Details:      []byte(details),
	}
}
