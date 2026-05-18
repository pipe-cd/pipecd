---
title: "Satisfying the DeploymentPlugin Interface"
weight: 8
description: >
  Analyzing the Go interface methods required to implement a DeploymentPlugin.
---

To register a plugin as a [DeploymentPlugin](https://pkg.go.dev/github.com/pipe-cd/piped-plugin-sdk-go#DeploymentPlugin) via `sdk.NewPlugin`, we must implement its interface methods.

First, let's understand the three type parameters required by `DeploymentPlugin`:

1. **`Config`**: The global configuration common to all instances of this plugin. This is defined in Piped's global configuration file.
2. **`DeployTargetConfig`**: The configuration unique to each deployment target. For example, in a Kubernetes plugin, this represents the target cluster credentials and context.
3. **`ApplicationConfigSpec`**: The configuration specified per application in `app.pipecd.yaml`. For example, a list of target manifest files or resource configurations.

### Interface Methods

To fully implement a `DeploymentPlugin` (including the embedded `StagePlugin` methods), we must define the following methods on our plugin struct:

- `FetchDefinedStages() []string`
- `DetermineVersions(context.Context, *Config, *DetermineVersionsInput[ApplicationConfigSpec]) (*DetermineVersionsResponse, error)`
- `DetermineStrategy(context.Context, *Config, *DetermineStrategyInput[ApplicationConfigSpec]) (*DetermineStrategyResponse, error)`
- `BuildPipelineSyncStages(context.Context, *Config, *BuildPipelineSyncStagesInput) (*BuildPipelineSyncStagesResponse, error)`
- `BuildQuickSyncStages(context.Context, *Config, *BuildQuickSyncStagesInput) (*BuildQuickSyncStagesResponse, error)`
- `ExecuteStage(context.Context, *Config, []*DeployTarget[DeployTargetConfig], *ExecuteStageInput[ApplicationConfigSpec]) (*ExecuteStageResponse, error)`

We will implement each of these methods step-by-step.
