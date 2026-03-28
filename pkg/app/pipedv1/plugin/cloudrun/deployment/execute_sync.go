// Copyright 2026 The PipeCD Authors.
package deployment

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/cloudrunservice/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/cloudrunservice/provider"
)

func executeSync(ctx context.Context, targets []*sdk.DeployTarget[config.CloudRunDeployTargetConfig], input *sdk.ExecuteStageInput[config.CloudRunApplicationSpec]) (*sdk.ExecuteStageResponse, error) {
	if len(targets) == 0 {
		return nil, fmt.Errorf("deploy target is not configured")
	}
	target := targets[0]

	manifestPath := "service.yaml"
	if input.Request.TargetDeploymentSource.ApplicationConfig != nil && input.Request.TargetDeploymentSource.ApplicationConfig.Spec != nil && input.Request.TargetDeploymentSource.ApplicationConfig.Spec.Input.ServiceManifestFile != "" {
		manifestPath = input.Request.TargetDeploymentSource.ApplicationConfig.Spec.Input.ServiceManifestFile
	}
	path := filepath.Join(input.Request.TargetDeploymentSource.ApplicationDirectory, manifestPath)
	sm, err := provider.LoadServiceManifest(path)
	if err != nil {
		return &sdk.ExecuteStageResponse{
			Status: sdk.StageStatusFailure,
		}, fmt.Errorf("failed to load service manifest: %w", err)
	}

	client, err := provider.NewClient(ctx, target.Config.Project, target.Config.Region, target.Config.CredentialsFile, input.Logger)
	if err != nil {
		return &sdk.ExecuteStageResponse{
			Status: sdk.StageStatusFailure,
		}, fmt.Errorf("failed to create cloud run client: %w", err)
	}

	commitHash := input.Request.TargetDeploymentSource.CommitHash
	if len(commitHash) > 7 {
		commitHash = commitHash[:7]
	}
	revisionName := fmt.Sprintf("%s-%s", sm.Name(), commitHash)

	if err := sm.SetRevision(revisionName); err != nil {
		return nil, fmt.Errorf("failed to set revision name: %w", err)
	}
	if err := sm.UpdateAllTraffic(revisionName); err != nil {
		return nil, fmt.Errorf("failed to configure traffic mapping: %w", err)
	}
	sm.AddLabels(map[string]string{
		"pipecd-dev-managed-by": "piped",
		"pipecd-dev-application": input.Request.Deployment.ApplicationID,
	})

	input.Logger.Info(fmt.Sprintf("applying service manifest for %s...", sm.Name()))
	_, err = client.GetService(ctx, sm.Name())
	if err == provider.ErrServiceNotFound {
		_, err = client.Create(ctx, sm)
	} else if err == nil {
		_, err = client.Update(ctx, sm)
	}

	if err != nil {
		return &sdk.ExecuteStageResponse{
			Status: sdk.StageStatusFailure,
		}, fmt.Errorf("failed to apply service manifest: %w", err)
	}

	input.Logger.Info(fmt.Sprintf("waiting for revision %s to be ready...", revisionName))
	time.Sleep(2 * time.Second)

	return &sdk.ExecuteStageResponse{
		Status: sdk.StageStatusSuccess,
	}, nil
}
