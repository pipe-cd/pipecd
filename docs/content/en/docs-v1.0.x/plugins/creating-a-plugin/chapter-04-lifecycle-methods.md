---
title: "Lifecycle methods: stages, versions, and strategy"
linkTitle: "Lifecycle methods: stages, versions, and strategy"
weight: 4
description: >
  Implement `FetchDefinedStages`, `DetermineVersions`, and `DetermineStrategy`.
---

In this chapter you replace the first three empty methods with real ones. Together they tell `piped` which stages the plugin provides, which versions it is deploying, and which deployment strategy to use. None of them run a deployment yet; they give `piped` the information it needs before a deployment starts.

## Declare the stages

`FetchDefinedStages` returns the names of the stages the plugin provides. The file plugin defines three: `FILE_DIFF` and `FILE_SYNC` for the two deployment stages, and `FILE_ROLLBACK` for the stage that runs when a deployment fails and needs to roll back.

Define the stage names as constants, because other methods refer to them too, and return them from `FetchDefinedStages`:

```go
const (
	stageDiff     = "FILE_DIFF"
	stageSync     = "FILE_SYNC"
	stageRollback = "FILE_ROLLBACK"
)

func (p *plugin) FetchDefinedStages() []string {
	return []string{stageDiff, stageSync, stageRollback}
}
```

## Report the deployed version

`DetermineVersions` reports the versions of the artifacts a deployment uses, computed from the deployment source for each deployment. The file plugin uses the Git commit hash from the deployment source as its version and does not set a name or URL:

```go
func (p *plugin) DetermineVersions(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineVersionsInput[applicationConfig]) (*sdk.DetermineVersionsResponse, error) {
	return &sdk.DetermineVersionsResponse{
		Versions: []sdk.ArtifactVersion{{Version: input.Request.DeploymentSource.CommitHash}},
	}, nil
}
```

## Choose a deployment strategy

`DetermineStrategy` decides how a deployment runs. It is only called when piped's built-in checks - pipeline length, whether this is the first deployment, and similar - do not determine a strategy on their own. There are two strategies:

- **Pipeline sync** runs the deployment pipeline the user defines in the application configuration.
- **Quick sync** skips the user's pipeline and runs a fixed set of stages the plugin defines, which is useful for small changes that do not need a full pipeline.

A plugin only needs to implement `DetermineStrategy` when it has its own logic for choosing between the two. The file plugin does not, so it returns `nil, nil`. Returning `nil, nil` means no custom logic: the SDK answers piped on the plugin's behalf with pipeline sync.

```go
func (p *plugin) DetermineStrategy(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineStrategyInput[applicationConfig]) (*sdk.DetermineStrategyResponse, error) {
	return nil, nil
}
```

Build the project again to confirm the three methods compile:

```bash
go build ./...
```

The plugin now reports its stages, version, and strategy. In the next chapter you build the stage lists that these strategies run, starting with the pipeline sync stages.
