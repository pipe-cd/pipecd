---
title: "Setup, plugin types, and project scaffold"
linkTitle: "Setup, plugin types, and project scaffold"
weight: 2
description: >
  Create the project, add the plugin SDK, and learn which plugin type you build.
---

In this chapter you create the project for the file plugin, add the plugin SDK, and learn where the plugin fits among the SDK's plugin types. By the end you have a Go module with the SDK in place, ready for the code you write in the following chapters.

## Set up the project

The plugin is a normal Go module. Create a directory for it, initialize a Git repository, and create the Go module:

```bash
git init pipecd-plugin-file
cd pipecd-plugin-file
go mod init github.com/<YOUR_USERNAME>/pipecd-plugin-file
```

Replace `<YOUR_USERNAME>` with your GitHub account name, or use any module path you prefer. The plugin builds with Go 1.26 or later, matching the version the official plugins use.

Commit as you go. This tutorial does not point out every commit from here on, but small, frequent commits make it easy to retrace your own steps.

## Add the plugin SDK

Plugins are built with the official plugin SDK ([`github.com/pipe-cd/piped-plugin-sdk-go`](https://pkg.go.dev/github.com/pipe-cd/piped-plugin-sdk-go)). Add it to the module:

```bash
go get github.com/pipe-cd/piped-plugin-sdk-go@v0.4.0
```

The SDK provides the plugin server, the interfaces you implement, and the request and response types that `piped` sends and expects. Apart from the SDK, the file plugin uses only the Go standard library.

## Plugin types

`piped` does not define separate kinds of plugins on its own. For convenience, the SDK groups plugins by the interface they implement:

- **StagePlugin** provides stages that are useful during a deployment but has nothing of its own to deploy. The `wait` plugin, which pauses a pipeline for a set time, is a StagePlugin.
- **DeploymentPlugin** has something to deploy and syncs it. The `kubernetes` plugin is a DeploymentPlugin. A DeploymentPlugin also provides everything a StagePlugin does.
- **LivestatePlugin** reports the live state of deployed resources, so the web UI can show the difference between what is running and what is defined in Git. It is often implemented alongside a DeploymentPlugin.

The file plugin treats copying files as its deployment, so it is a **DeploymentPlugin**.

## The DeploymentPlugin interface

A DeploymentPlugin has three type parameters. They let the SDK decode configuration into types that you define:

- **Config** is configuration shared across the plugin, written in the `piped` configuration.
- **DeployTargetConfig** is configuration for a single deploy target, such as the connection details for a cluster.
- **ApplicationConfigSpec** is per-application configuration, such as the files an application deploys.

The file plugin needs neither plugin-wide nor deploy-target configuration, so its Config and DeployTargetConfig are empty. You define all three types in the next chapter.

To satisfy the DeploymentPlugin interface, you implement the following methods:

```go
FetchDefinedStages() []string
DetermineVersions(context.Context, *Config, *DetermineVersionsInput[ApplicationConfigSpec]) (*DetermineVersionsResponse, error)
DetermineStrategy(context.Context, *Config, *DetermineStrategyInput[ApplicationConfigSpec]) (*DetermineStrategyResponse, error)
BuildPipelineSyncStages(context.Context, *Config, *BuildPipelineSyncStagesInput) (*BuildPipelineSyncStagesResponse, error)
BuildQuickSyncStages(context.Context, *Config, *BuildQuickSyncStagesInput) (*BuildQuickSyncStagesResponse, error)
ExecuteStage(context.Context, *Config, []*DeployTarget[DeployTargetConfig], *ExecuteStageInput[ApplicationConfigSpec]) (*ExecuteStageResponse, error)
```

You implement these across the next several chapters, starting from the top. For now the project is set up and the SDK is in place, so the next chapter defines the configuration types and writes an empty implementation that satisfies this interface.
