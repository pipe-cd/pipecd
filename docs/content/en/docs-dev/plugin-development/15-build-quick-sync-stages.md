---
title: "Implementation: BuildQuickSyncStages"
weight: 15
description: >
  Implementing BuildQuickSyncStages to configure the fallback sync pipeline.
---

`BuildQuickSyncStages` is called when `QuickSync` is selected. It returns a default set of execution stages defined entirely by the plugin, bypassing the user's pipeline.

### Differences from BuildPipelineSyncStages

1. **No `Index` Field**: In a Quick Sync, stages are executed concurrently or without a rigid order, so `Index` is not required.
2. **`Description` Requirement**: Since there is no user-defined pipeline configuration, the plugin must supply a human-readable `Description` for each stage.

For our `file` plugin's Quick Sync, we will execute a single `FILE_SYNC` stage, along with a `FILE_ROLLBACK` stage if rollback is requested.

Update `BuildQuickSyncStages` as follows:

```go
func (plugin) BuildQuickSyncStages(_ context.Context, _ *config, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	stages := make([]sdk.QuickSyncStage, 0, 2)

	stages = append(stages, sdk.QuickSyncStage{
		Name:        stageSync,
		Description: "Synchronize local files from the Git repository",
	})

	if input.Request.Rollback {
		stages = append(stages, sdk.QuickSyncStage{
			Name:        stageRollback,
			Description: "Rollback local files to the previous successful version",
			Rollback:    true,
		})
	}

	return &sdk.BuildQuickSyncStagesResponse{
		Stages: stages,
	}, nil
}
```
