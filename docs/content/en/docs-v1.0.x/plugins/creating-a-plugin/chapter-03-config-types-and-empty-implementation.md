---
title: "Config types and an empty implementation"
linkTitle: "Config types and an empty implementation"
weight: 3
description: >
  Define the plugin's configuration types and a compiling, empty DeploymentPlugin.
---

In this chapter you define the configuration types the plugin uses and add an empty implementation of the DeploymentPlugin interface. The methods do nothing yet, but the project compiles and the plugin starts, which gives you a foundation to fill in one method at a time.

## Define the configuration types

Recall the three type parameters of a DeploymentPlugin: Config, DeployTargetConfig, and ApplicationConfigSpec. The file plugin needs only one of them.

- It has no plugin-wide configuration, so it uses `sdk.ConfigNone`, the SDK's type for an empty configuration, in place of a Config type.
- It has no deploy targets, so its DeployTargetConfig is an empty struct.
- It deploys to a directory that is chosen per application, so its application configuration holds a single `Path` field.

Create `plugin.go` with the configuration types:

```go
package main

import (
	"context"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

type (
	// deployTargetConfig is empty because the file plugin has no deploy targets.
	deployTargetConfig struct{}

	// applicationConfig holds the per-application settings. The file plugin
	// deploys to a directory chosen per application, so it defines Path.
	applicationConfig struct {
		// Path is the directory the application deploys to.
		Path string `json:"path"`
	}
)
```

Both `context` and `sdk` are used by the implementation you add next.

## Add an empty implementation

Define a `plugin` type and give it a method for every method in the DeploymentPlugin interface. Each method calls `panic` for now. This lets you confirm that the method set and the signatures are correct before writing any real logic.

Add the following to `plugin.go`:

```go
// plugin provides the DeploymentPlugin implementation for the file plugin.
type plugin struct{}

// Verify at compile time that plugin implements sdk.DeploymentPlugin.
// A missing method or a wrong signature fails the build here.
var _ sdk.DeploymentPlugin[sdk.ConfigNone, deployTargetConfig, applicationConfig] = (*plugin)(nil)

func (p *plugin) FetchDefinedStages() []string {
	panic("unimplemented")
}

func (p *plugin) DetermineVersions(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineVersionsInput[applicationConfig]) (*sdk.DetermineVersionsResponse, error) {
	panic("unimplemented")
}

func (p *plugin) DetermineStrategy(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineStrategyInput[applicationConfig]) (*sdk.DetermineStrategyResponse, error) {
	panic("unimplemented")
}

func (p *plugin) BuildPipelineSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	panic("unimplemented")
}

func (p *plugin) BuildQuickSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	panic("unimplemented")
}

func (p *plugin) ExecuteStage(ctx context.Context, _ *sdk.ConfigNone, dts []*sdk.DeployTarget[deployTargetConfig], input *sdk.ExecuteStageInput[applicationConfig]) (*sdk.ExecuteStageResponse, error) {
	panic("unimplemented")
}
```

The `var _ sdk.DeploymentPlugin[...] = (*plugin)(nil)` line is a compile-time check. If a method is missing or its signature is wrong, the build fails on this line rather than somewhere less obvious.

## Register the plugin

Now that a plugin type exists, register it so `piped` can load it. Create `main.go`:

```go
package main

import (
	"log"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

func main() {
	p, err := sdk.NewPlugin(
		"0.0.1",
		sdk.WithDeploymentPlugin(&plugin{}),
	)
	if err != nil {
		log.Fatalln(err)
	}

	if err := p.Run(); err != nil {
		log.Fatalln(err)
	}
}
```

`sdk.NewPlugin` creates the plugin server, `sdk.WithDeploymentPlugin` registers your implementation, and `p.Run()` starts the server. The first argument, `"0.0.1"`, is the plugin's own version.

Build the project to confirm everything fits together:

```bash
go build ./...
```

The project now compiles and the plugin starts, though every stage still panics. In the next chapter you replace the first `panic` with a real implementation of `FetchDefinedStages`.
