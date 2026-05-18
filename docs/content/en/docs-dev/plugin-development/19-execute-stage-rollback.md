---
title: "ExecuteStage Implementation: ROLLBACK Stage"
weight: 19
description: >
  Implementing the ROLLBACK stage to restore local files in case of failure.
---

The final step in our sync engine is the `FILE_ROLLBACK` stage.

Fortunately, rollback is structurally identical to the sync operation. The only difference is that instead of syncing the target directory with the _newly requested_ deployment source (`input.Request.TargetDeploymentSource`), we sync it with the _previously successful_ deployment source (`input.Request.RunningDeploymentSource`).

Update `executeStageRollback` as follows:

```go
func (plugin) executeStageRollback(ctx context.Context, input *sdk.ExecuteStageInput[applicationConfig]) (*sdk.ExecuteStageResponse, error) {
	lp := input.Client.LogPersister()

	lp.Info("Restoring files to target directory (rolling back)...")
	if err := copyFiles(
		input.Request.RunningDeploymentSource.ApplicationConfig.Spec.Path,
		os.DirFS(input.Request.RunningDeploymentSource.ApplicationDirectory),
		map[string]struct{}{
			input.Request.RunningDeploymentSource.ApplicationConfigFilename: {},
		},
	); err != nil {
		return nil, fmt.Errorf("error copying files during rollback: %w", err)
	}

	lp.Info("Removing newer files from target directory...")
	if err := removeFiles(
		input.Request.RunningDeploymentSource.ApplicationConfig.Spec.Path,
		os.DirFS(input.Request.RunningDeploymentSource.ApplicationDirectory),
		map[string]struct{}{
			input.Request.RunningDeploymentSource.ApplicationConfigFilename: {},
		},
	); err != nil {
		return nil, fmt.Errorf("error removing files during rollback: %w", err)
	}

	lp.Success("File rollback completed successfully")
	return &sdk.ExecuteStageResponse{
		Status: sdk.StageStatusSuccess,
	}, nil
}
```
