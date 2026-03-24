// Copyright 2026 The PipeCD Authors.
package deployment

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/cloudrunservice/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/cloudrunservice/provider"
)

func determineVersions(ctx context.Context, input *sdk.DetermineVersionsInput[config.CloudRunApplicationSpec]) (*sdk.DetermineVersionsResponse, error) {
	manifestPath := "service.yaml"
	if input.Request.DeploymentSource.ApplicationConfig != nil && input.Request.DeploymentSource.ApplicationConfig.Spec != nil && input.Request.DeploymentSource.ApplicationConfig.Spec.Input.ServiceManifestFile != "" {
		manifestPath = input.Request.DeploymentSource.ApplicationConfig.Spec.Input.ServiceManifestFile
	}
	path := filepath.Join(input.Request.DeploymentSource.ApplicationDirectory, manifestPath)

	sm, err := provider.LoadServiceManifest(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load service manifest: %w", err)
	}

	images, err := sm.ExtractImages()
	if err != nil {
		return nil, fmt.Errorf("failed to extract images: %w", err)
	}

	versions := make([]sdk.ArtifactVersion, 0, len(images))
	for _, image := range images {
		name, tag := parseContainerImage(image)
		versions = append(versions, sdk.ArtifactVersion{
			Version: tag,
			Name:    name,
			URL:     image,
		})
	}

	return &sdk.DetermineVersionsResponse{
		Versions: versions,
	}, nil
}

func parseContainerImage(image string) (name, tag string) {
	lastColon := strings.LastIndex(image, ":")
	lastSlash := strings.LastIndex(image, "/")

	if lastColon > lastSlash {
		tag = image[lastColon+1:]
		imageWithoutTag := image[:lastColon]
		paths := strings.Split(imageWithoutTag, "/")
		name = paths[len(paths)-1]
	} else {
		paths := strings.Split(image, "/")
		name = paths[len(paths)-1]
		tag = "latest"
	}
	return
}

func determineStrategy(ctx context.Context, input *sdk.DetermineStrategyInput[config.CloudRunApplicationSpec]) (*sdk.DetermineStrategyResponse, error) {
	if input.Request.RunningDeploymentSource.ApplicationDirectory == "" {
		return &sdk.DetermineStrategyResponse{
			Strategy: sdk.SyncStrategyQuickSync,
			Summary:  "First time deployment. Quick sync will be used",
		}, nil
	}

	return &sdk.DetermineStrategyResponse{
		Strategy: sdk.SyncStrategyPipelineSync,
		Summary:  "Updating existing deployment. Pipeline sync will be used",
	}, nil
}
