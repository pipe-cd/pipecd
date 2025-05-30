---
title: "Tool Registry"
description: "How to manage external tools in your PipeCD plugin"
weight: 2
---

The Tool Registry helps you manage external tools (like Terraform, OpenTofu, Helm, etc.) that your plugin needs to function. It handles downloading, installing, and managing these tools automatically.

## Basic Implementation

Here's how to implement a tool registry for your plugin:

```go
// toolregistry/registry.go
package toolregistry

import (
    "context"
    "fmt"
    "os/exec"
    "path/filepath"

    "github.com/pipe-cd/pipecd/pkg/plugin/toolregistry"
)

// Registry manages the lifecycle of your tool
type Registry struct {
    client toolregistry.ToolRegistry
}

// NewRegistry creates a new registry instance
func NewRegistry(client toolregistry.ToolRegistry) *Registry {
    return &Registry{
        client: client,
    }
}

// GetTool downloads and installs the tool if needed, returns the path to the executable
func (r *Registry) GetTool(ctx context.Context, version string) (string, error) {
    return r.client.InstallTool(ctx, version)
}
```

## Installation Script

Create an installation script that defines how to download and install your tool:

```go
// toolregistry/scripts.go
package toolregistry

const installScript = `
#!/bin/sh
set -e

# Download the tool
curl -L "https://github.com/your-tool/releases/download/v{{ .Version }}/tool_{{ .Version }}_linux_amd64.zip" -o tool.zip

# Install it
unzip tool.zip
mv tool {{ .OutPath }}
chmod +x {{ .OutPath }}

# Clean up
rm tool.zip
`
```

## Using the Registry

Here's how to use the registry in your plugin:

```go
// In your plugin code
type Plugin struct {
    // Note: tool registry is created when needed, not stored as a field
}

func (p *Plugin) Execute(ctx context.Context, input *plugin.Input) error {
    // Create tool registry when needed
    toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())
    
    // Get the tool path with version
    // Each plugin should define its own tool methods (e.g. GetTool, GetCustomTool)
    toolPath, err := toolRegistry.GetTool(ctx, "1.0.0")
    if err != nil {
        return fmt.Errorf("failed to get tool: %w", err)
    }

    // Use the tool
    cmd := exec.CommandContext(ctx, toolPath, "your-command")
    cmd.Dir = input.WorkingDir
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("command failed: %w", err)
    }

    return nil
}
```

### Key Points

1. **Tool Registry Creation**
   - Create the tool registry when needed in your execution methods
   - Don't store it as a field in your plugin struct
   - Use `toolregistry.NewRegistry(input.Client.ToolRegistry())`

2. **Tool Methods**
   - Define specific methods for each tool your plugin needs
   - Example:
   ```go
   type Registry struct {
       client client
   }

   func (r *Registry) GetTool(ctx context.Context, version string) (string, error) {
       return r.client.InstallTool(ctx, "tool-name", version, installScript)
   }
   ```

3. **Version Handling**
   - Consider providing default versions for your tools
   - Handle version fallbacks gracefully
   - Example:
   ```go
   const defaultToolVersion = "1.0.0"

   func (r *Registry) GetTool(ctx context.Context, version string) (string, error) {
       actualVersion := version
       if actualVersion == "" {
           actualVersion = defaultToolVersion
       }
       return r.client.InstallTool(ctx, "tool-name", actualVersion, installScript)
   }
   ```


4. **Tool Usage**
   - Use the tool path with `exec.CommandContext`


Remember that each plugin might have different tool requirements and should implement its own specific tool registry methods accordingly.

## Example: OpenTofu Plugin

Here's how we use the Tool Registry in our OpenTofu plugin:

```go
// Get OpenTofu binary
tofuPath, err := p.toolRegistry.GetTool(ctx, version)
if err != nil {
    return fmt.Errorf("failed to get OpenTofu binary: %w", err)
}

// Use OpenTofu
cmd := exec.CommandContext(ctx, tofuPath, "init")
cmd.Dir = workDir
if err := cmd.Run(); err != nil {
    return fmt.Errorf("failed to run command: %w", err)
}
```

## Next Steps

- Check [Example Plugins](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin)
- Explore [Plugin SDK](https://pkg.go.dev/github.com/pipe-cd/pipecd/pkg/plugin/sdk)