---
title: "Implementation: Empty Implementation to Satisfy the Interface"
weight: 10
description: >
  Scaffolding the plugin struct and satisfying the interface with stub methods.
---

Before implementing the logic for each method, let's establish a skeleton structure in `main.go`. This guarantees our method signatures match the expected interface and allows the project to compile.

Add the following code to `main.go` (or a separate file under `package main`):

```go
package main

import (
	"context"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go/pkg/plugin/sdk"
)

// Ensure that plugin implements the sdk.DeploymentPlugin interface.
// If any methods are missing or have incorrect signatures, this will trigger a compile-time error.
var _ sdk.DeploymentPlugin[config, deployTargetConfig, applicationConfig] = plugin{}

// plugin provides the implementation of the DeploymentPlugin interface.
type plugin struct{}

// FetchDefinedStages returns the list of stages supported by this plugin.
func (plugin) FetchDefinedStages() []string {
	panic("unimplemented")
}

// DetermineVersions determines the version of the resource being deployed.
func (plugin) DetermineVersions(ctx context.Context, cfg *config, input *sdk.DetermineVersionsInput[applicationConfig]) (*sdk.DetermineVersionsResponse, error) {
	panic("unimplemented")
}

// DetermineStrategy decides whether to execute a Quick Sync or a Pipeline Sync.
func (plugin) DetermineStrategy(ctx context.Context, cfg *config, input *sdk.DetermineStrategyInput[applicationConfig]) (*sdk.DetermineStrategyResponse, error) {
	panic("unimplemented")
}

// BuildPipelineSyncStages constructs the pipeline stages for a Pipeline Sync.
func (plugin) BuildPipelineSyncStages(ctx context.Context, cfg *config, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	panic("unimplemented")
}

// BuildQuickSyncStages constructs the stages for a Quick Sync.
func (plugin) BuildQuickSyncStages(ctx context.Context, cfg *config, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	panic("unimplemented")
}

// ExecuteStage runs the actual sync logic for a specific stage.
func (plugin) ExecuteStage(ctx context.Context, cfg *config, targets []*sdk.DeployTarget[deployTargetConfig], input *sdk.ExecuteStageInput[applicationConfig]) (*sdk.ExecuteStageResponse, error) {
	panic("unimplemented")
}
```
