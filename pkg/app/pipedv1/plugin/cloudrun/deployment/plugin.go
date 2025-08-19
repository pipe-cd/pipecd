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
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/cloudrunservice/config"
)

type Plugin struct{}

const (
	// StageCloudRunSync does quick sync by rolling out the new version
	// and switching all traffic to it.
	StageCloudRunSync = "CLOUDRUN_SYNC"
	// StageCloudRunPromote promotes the new version to receive amount of traffic.
	StageCloudRunPromote = "CLOUDRUN_PROMOTE"
	// StageRollback the legacy generic rollback stage name
	StageRollback = "ROLLBACK"

	StageCloudRunSyncDescription = "Deploy the new version and configure all traffic to it"
	StageRollbackDescription     = "Rollback the deployment"

	DefaultServiceManifestFilename = "service.yaml"
)

func (p *Plugin) FetchDefinedStages() []string {
	return []string{StageCloudRunSync, StageCloudRunPromote, StageRollback}
}

func (p *Plugin) BuildPipelineSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	return &sdk.BuildPipelineSyncStagesResponse{
		Stages: buildPipelineStages(input.Request.Stages, input.Request.Rollback),
	}, nil
}

func (p *Plugin) ExecuteStage(ctx context.Context, _ *sdk.ConfigNone, dts []*sdk.DeployTarget[config.CloudRunDeployTargetConfig], input *sdk.ExecuteStageInput[config.CloudRunApplicationSpec]) (*sdk.ExecuteStageResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (p *Plugin) DetermineVersions(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineVersionsInput[config.CloudRunApplicationSpec]) (*sdk.DetermineVersionsResponse, error) {
	sm, err := loadServiceManifest(input.Request.DeploymentSource.ApplicationDirectory, input.Request.DeploymentSource.ApplicationConfigFilename)
	if err != nil {
		return nil, err
	}
	versions, err := findArtifactVersions(sm)
	if err != nil {
		return nil, err
	}
	return &sdk.DetermineVersionsResponse{
		Versions: versions,
	}, nil
}

func (p *Plugin) DetermineStrategy(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineStrategyInput[config.CloudRunApplicationSpec]) (*sdk.DetermineStrategyResponse, error) {
	runningServiceManifest, err1 := loadServiceManifest(input.Request.RunningDeploymentSource.ApplicationDirectory, input.Request.RunningDeploymentSource.ApplicationConfigFilename)
	targetServiceManifest, err2 := loadServiceManifest(input.Request.TargetDeploymentSource.ApplicationDirectory, input.Request.TargetDeploymentSource.ApplicationConfigFilename)
	if err1 == nil && err2 == nil {
		oldVersion, err1 := findArtifactVersions(runningServiceManifest)
		newVersion, err2 := findArtifactVersions(targetServiceManifest)
		if err1 == nil && err2 == nil {
			return &sdk.DetermineStrategyResponse{
				Strategy: sdk.SyncStrategyPipelineSync,
				Summary:  fmt.Sprintf("Sync with pipeline to update image from %s to %s", oldVersion, newVersion),
			}, nil
		}
	}
	return &sdk.DetermineStrategyResponse{
		Strategy: sdk.SyncStrategyPipelineSync,
		Summary:  "Sync with the specified pipeline",
	}, nil
}

func (p *Plugin) BuildQuickSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	return &sdk.BuildQuickSyncStagesResponse{
		Stages: buildQuickSyncPipeline(input.Request.Rollback),
	}, nil
}

func buildQuickSyncPipeline(autoRollback bool) []sdk.QuickSyncStage {
	out := make([]sdk.QuickSyncStage, 0, 2)
	out = append(out, sdk.QuickSyncStage{
		Name:               StageCloudRunSync,
		Description:        StageCloudRunSyncDescription,
		Rollback:           false,
		Metadata:           map[string]string{},
		AvailableOperation: sdk.ManualOperationNone,
	})
	if autoRollback {
		out = append(out, sdk.QuickSyncStage{
			Name:               StageRollback,
			Description:        StageRollbackDescription,
			Rollback:           true,
			Metadata:           map[string]string{},
			AvailableOperation: sdk.ManualOperationNone,
		})
	}
	return out
}

func buildPipelineStages(stages []sdk.StageConfig, autoRollback bool) []sdk.PipelineStage {
	out := make([]sdk.PipelineStage, 0, len(stages)+1)
	for _, stage := range stages {
		out = append(out, sdk.PipelineStage{
			Name:               stage.Name,
			Index:              stage.Index,
			Rollback:           false,
			Metadata:           map[string]string{},
			AvailableOperation: sdk.ManualOperationNone,
		})
	}
	if autoRollback {
		out = append(out, sdk.PipelineStage{
			Name: StageRollback,
			Index: slices.MinFunc(stages, func(a, b sdk.StageConfig) int {
				return a.Index - b.Index
			}).Index,
			Rollback:           true,
			Metadata:           map[string]string{},
			AvailableOperation: sdk.ManualOperationNone,
		})
	}
	return out
}

type ServiceManifest struct {
	Name string
	u    *unstructured.Unstructured
}

func loadServiceManifest(appDir, serviceManifestFile string) (ServiceManifest, error) {
	if serviceManifestFile == "" {
		serviceManifestFile = DefaultServiceManifestFilename
	}
	path := filepath.Join(appDir, serviceManifestFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return ServiceManifest{}, err
	}
	return parseServiceManifest(data)
}

func parseServiceManifest(data []byte) (ServiceManifest, error) {
	var obj unstructured.Unstructured
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return ServiceManifest{}, err
	}

	return ServiceManifest{
		Name: obj.GetName(),
		u:    &obj,
	}, nil
}

func findArtifactVersions(sm ServiceManifest) ([]sdk.ArtifactVersion, error) {
	containers, ok, err := unstructured.NestedSlice(sm.u.Object, "spec", "template", "spec", "containers")
	if err != nil {
		return nil, err
	}
	if !ok || len(containers) == 0 {
		return nil, fmt.Errorf("spec.template.spec.containers was missing")
	}

	container, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&containers[0])
	if err != nil {
		return nil, fmt.Errorf("invalid container format")
	}

	image, ok, err := unstructured.NestedString(container, "image")
	if err != nil {
		return nil, err
	}
	if !ok || image == "" {
		return nil, fmt.Errorf("image was missing")
	}
	name, tag := parseContainerImage(image)

	return []sdk.ArtifactVersion{
		{
			Version: tag,
			Name:    name,
			URL:     image,
		},
	}, nil
}

func parseContainerImage(image string) (name, tag string) {
	parts := strings.Split(image, ":")
	if len(parts) == 2 {
		tag = parts[1]
	}
	paths := strings.Split(parts[0], "/")
	name = paths[len(paths)-1]
	return
}
