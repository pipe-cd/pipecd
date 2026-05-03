---
title: "Your First Plugin"
linkTitle: "Your First Plugin"
weight: 20
description: >
  Setting up and building a minimal stage plugin.
---

Let's build a simple "Hello World" plugin that executes a custom stage.

## 1. Initialize the Project

Create a new directory for your plugin and initialize the Go module:

```bash
mkdir my-plugin
cd my-plugin
go mod init github.com/your-username/my-plugin
```

Add the PipeCD Plugin SDK as a dependency:

```bash
go get github.com/pipe-cd/piped-plugin-sdk-go
```

## 2. The Main Entry Point

Create a `main.go` file. This file will use the SDK to start a gRPC server and register your plugin implementation.

```go
package main

import (
    "log"
    sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

func main() {
    // Create a new plugin instance and register our implementation
    p, err := sdk.NewPlugin("0.1.0", sdk.WithStagePlugin(&myPlugin{}))
    if err != nil {
        log.Fatalln(err)
    }

    // Run the gRPC server
    if err := p.Run(); err != nil {
        log.Fatalln(err)
    }
}
```

## 3. Implementing the Interface

Now, create `plugin.go` to implement the required methods for a `StagePlugin`.

```go
package main

import (
    "context"
    "fmt"
    sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

type myPlugin struct{}

// FetchDefinedStages returns the names of stages this plugin handles.
func (p *myPlugin) FetchDefinedStages() []string {
    return []string{"HELLO_WORLD"}
}

// BuildPipelineSyncStages prepares the stages for execution.
func (p *myPlugin) BuildPipelineSyncStages(ctx context.Context, _ sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
    stages := make([]sdk.PipelineStage, 0, len(input.Request.Stages))
    for _, rs := range input.Request.Stages {
        stages = append(stages, sdk.PipelineStage{
            Index: rs.Index,
            Name:  rs.Name,
        })
    }
    return &sdk.BuildPipelineSyncStagesResponse{Stages: stages}, nil
}

// ExecuteStage is where the actual work happens.
func (p *myPlugin) ExecuteStage(ctx context.Context, _ sdk.ConfigNone, _ sdk.DeployTargetsNone, input *sdk.ExecuteStageInput[struct{}]) (*sdk.ExecuteStageResponse, error) {
    // Access the log persister to send logs back to PipeCD
    lp := input.LogPersister
    lp.Infof("Hello from the plugin! Executing stage: %s", input.Stage.Name)

    return &sdk.ExecuteStageResponse{
        Status: sdk.StageStatusSuccess,
    }, nil
}
```

## 4. Understanding the methods

- **`FetchDefinedStages`**: This tells `piped` which stage names in the YAML should be routed to this plugin.
- **`BuildPipelineSyncStages`**: This is called when planning the deployment. It allows you to transform raw stage requests into structured `PipelineStage` objects.
- **`ExecuteStage`**: The core execution logic. You receive the stage data and can use the `LogPersister` to provide feedback to the user via the Web UI.

In the next chapter, we'll see how to add configuration options to our plugin.
