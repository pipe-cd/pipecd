---
title: "Configuration"
linkTitle: "Configuration"
weight: 30
description: >
  Defining and parsing plugin configuration.
---

Most plugins need input from the user. For example, a "Wait" plugin needs to know how long to wait.

In PipeCD, configuration is passed through the `app.pipecd.yaml` file in the Git repository.

## 1. Define the Options Struct

In your plugin, define a Go struct that matches the YAML structure of your stage configuration. Use JSON tags for mapping.

```go
package main

import (
    "encoding/json"
    "fmt"
)

// MyStageOptions defines the fields allowed in the pipeline stage.
type MyStageOptions struct {
    Message string `json:"message"`
    Repeat  int    `json:"repeat"`
}

// Helper function to decode raw JSON from Piped
func decodeOptions(data []byte) (MyStageOptions, error) {
    var opts MyStageOptions
    if err := json.Unmarshal(data, &opts); err != nil {
        return opts, err
    }
    return opts, nil
}
```

## 2. Access Configuration in `ExecuteStage`

Inside the `ExecuteStage` method, you can access the raw configuration from the input request and decode it.

```go
func (p *myPlugin) ExecuteStage(ctx context.Context, _ sdk.ConfigNone, _ sdk.DeployTargetsNone, input *sdk.ExecuteStageInput[struct{}]) (*sdk.ExecuteStageResponse, error) {
    lp := input.LogPersister

    // Decode the configuration
    opts, err := decodeOptions(input.Request.StageConfig)
    if err != nil {
        lp.Errorf("Failed to decode config: %v", err)
        return &sdk.ExecuteStageResponse{Status: sdk.StageStatusFailure}, nil
    }

    // Use the options
    for i := 0; i < opts.Repeat; i++ {
        lp.Infof("Message %d: %s", i+1, opts.Message)
    }

    return &sdk.ExecuteStageResponse{Status: sdk.StageStatusSuccess}, nil
}
```

## 3. Usage in `app.pipecd.yaml`

Once your plugin is registered (which we'll cover in the next chapter), users can use it in their application configuration like this:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  pipeline:
    stages:
      - name: MY_CUSTOM_STAGE
        with:
          message: "Hello from Git!"
          repeat: 3
```

The fields under `with` are passed as `StageConfig` to your plugin.

---

> [!NOTE]
> You can also define **Plugin-wide configuration** (set in the `piped` config) and **Deploy Target configuration**. These are passed as the second and third arguments to `ExecuteStage` respectively.
