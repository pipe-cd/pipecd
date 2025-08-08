---
title: "Plugin Development Guide for PipeCD v1"
linkTitle: "Plugin Development"
weight: 20
description: >
  Comprehensive guide for developing PipeCD v1 plugins, including SDK usage, best practices, and examples
---

# Plugin Development Guide for PipeCD v1

This guide provides everything you need to develop plugins for PipeCD v1, from basic concepts to advanced implementation patterns.

## Overview

PipeCD v1 plugins are independent processes that communicate with the Piped core via gRPC. They can be written in any programming language that supports gRPC, though the official SDK is currently available in Go.

## Plugin Types and Interfaces

### 1. Deployment Plugins

Deployment plugins handle the complete deployment lifecycle for specific platforms:

```go
type DeploymentServiceServer interface {
    // Fetch stages that this plugin supports
    FetchDefinedStages(context.Context, *FetchDefinedStagesRequest) (*FetchDefinedStagesResponse, error)
    
    // Determine artifact versions for deployment
    DetermineVersions(context.Context, *DetermineVersionsRequest) (*DetermineVersionsResponse, error)
    
    // Determine deployment strategy
    DetermineStrategy(context.Context, *DetermineStrategyRequest) (*DetermineStrategyResponse, error)
    
    // Build deployment pipeline stages
    BuildPipelineSyncStages(context.Context, *BuildPipelineSyncStagesRequest) (*BuildPipelineSyncStagesResponse, error)
    
    // Build quick sync stages
    BuildQuickSyncStages(context.Context, *BuildQuickSyncStagesRequest) (*BuildQuickSyncStagesResponse, error)
    
    // Execute a specific stage
    ExecuteStage(context.Context, *ExecuteStageRequest) (*ExecuteStageResponse, error)
}
```

### 2. LiveState Plugins

LiveState plugins fetch and report the current state of deployed resources:

```go
type LivestateServiceServer interface {
    // Fetch the current state of resources
    FetchDefinedResources(context.Context, *FetchDefinedResourcesRequest) (*FetchDefinedResourcesResponse, error)
    
    // Get current resource state
    GetResources(context.Context, *GetResourcesRequest) (*GetResourcesResponse, error)
}
```

### 3. Stage Plugins

Stage plugins provide specific deployment stages (like Wait, Approval, Analysis):

```go
// Stage plugins implement the ExecuteStage method from DeploymentService
// and return their supported stages via FetchDefinedStages
```

## Getting Started with Plugin Development

### Prerequisites

- Go 1.21 or later (for Go plugins)
- Understanding of gRPC concepts
- Familiarity with your target platform/service

### 1. Project Setup

Create a new Go module for your plugin:

```bash
mkdir my-pipecd-plugin
cd my-pipecd-plugin
go mod init github.com/myorg/my-pipecd-plugin
```

Add the PipeCD plugin SDK dependency:

```bash
go get github.com/pipe-cd/piped-plugin-sdk-go
```

### 2. Basic Plugin Structure

```go
package main

import (
    "log"
    
    sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

func main() {
    plugin, err := sdk.NewPlugin(
        "1.0.0", // plugin version
        sdk.WithDeploymentPlugin(&MyDeploymentPlugin{}),
        // sdk.WithLivestatePlugin(&MyLivestatePlugin{}), // optional
    )
    if err != nil {
        log.Fatalln(err)
    }
    
    if err := plugin.Run(); err != nil {
        log.Fatalln(err)
    }
}
```

### 3. Implementing a Deployment Plugin

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
    "github.com/pipe-cd/pipecd/pkg/model"
)

type MyDeploymentPlugin struct{}

// FetchDefinedStages returns the stages this plugin supports
func (p *MyDeploymentPlugin) FetchDefinedStages(
    ctx context.Context,
    req *deployment.FetchDefinedStagesRequest,
) (*deployment.FetchDefinedStagesResponse, error) {
    return &deployment.FetchDefinedStagesResponse{
        Stages: []string{
            "MY_DEPLOY",
            "MY_ROLLBACK",
            "MY_CLEANUP",
        },
    }, nil
}

// DetermineVersions extracts version information from the deployment
func (p *MyDeploymentPlugin) DetermineVersions(
    ctx context.Context,
    req *deployment.DetermineVersionsRequest,
) (*deployment.DetermineVersionsResponse, error) {
    // Implement version detection logic
    // This could parse manifests, container images, etc.
    versions := []*model.ArtifactVersion{
        {
            Kind:    model.ArtifactVersion_CONTAINER_IMAGE,
            Version: "v1.2.3",
        },
    }
    
    return &deployment.DetermineVersionsResponse{
        Versions: versions,
    }, nil
}

// DetermineStrategy decides the deployment strategy
func (p *MyDeploymentPlugin) DetermineStrategy(
    ctx context.Context,
    req *deployment.DetermineStrategyRequest,
) (*deployment.DetermineStrategyResponse, error) {
    // Implement strategy logic
    // Could be based on configuration, file changes, etc.
    return &deployment.DetermineStrategyResponse{
        SyncStrategy: model.SyncStrategy_PIPELINE,
        Summary:      "Pipeline deployment with progressive rollout",
    }, nil
}

// BuildPipelineSyncStages creates the deployment pipeline
func (p *MyDeploymentPlugin) BuildPipelineSyncStages(
    ctx context.Context,
    req *deployment.BuildPipelineSyncStagesRequest,
) (*deployment.BuildPipelineSyncStagesResponse, error) {
    var stages []*model.PipelineStage
    
    for _, stageConfig := range req.Stages {
        stage := &model.PipelineStage{
            Id:      stageConfig.Index,
            Name:    stageConfig.Name,
            Desc:    stageConfig.Desc,
            Index:   stageConfig.Index,
            Timeout: stageConfig.Timeout,
            Status:  model.StageStatus_STAGE_NOT_STARTED_YET,
        }
        stages = append(stages, stage)
    }
    
    return &deployment.BuildPipelineSyncStagesResponse{
        Stages: stages,
    }, nil
}

// BuildQuickSyncStages creates quick deployment stages
func (p *MyDeploymentPlugin) BuildQuickSyncStages(
    ctx context.Context,
    req *deployment.BuildQuickSyncStagesRequest,
) (*deployment.BuildQuickSyncStagesResponse, error) {
    // Quick sync usually has just one stage
    stages := []*model.PipelineStage{
        {
            Id:     0,
            Name:   "MY_QUICK_SYNC",
            Desc:   "Quick deployment without progressive rollout",
            Index:  0,
            Status: model.StageStatus_STAGE_NOT_STARTED_YET,
        },
    }
    
    return &deployment.BuildQuickSyncStagesResponse{
        Stages: stages,
    }, nil
}

// ExecuteStage performs the actual deployment work
func (p *MyDeploymentPlugin) ExecuteStage(
    ctx context.Context,
    req *deployment.ExecuteStageRequest,
) (*deployment.ExecuteStageResponse, error) {
    stageName := req.Stage.Name
    
    switch stageName {
    case "MY_DEPLOY":
        return p.executeDeploy(ctx, req)
    case "MY_ROLLBACK":
        return p.executeRollback(ctx, req)
    case "MY_CLEANUP":
        return p.executeCleanup(ctx, req)
    case "MY_QUICK_SYNC":
        return p.executeQuickSync(ctx, req)
    default:
        return &deployment.ExecuteStageResponse{
            Status:  model.StageStatus_STAGE_FAILURE,
            Message: fmt.Sprintf("Unknown stage: %s", stageName),
        }, nil
    }
}

func (p *MyDeploymentPlugin) executeDeploy(
    ctx context.Context,
    req *deployment.ExecuteStageRequest,
) (*deployment.ExecuteStageResponse, error) {
    // Implement your deployment logic here
    // Access deployment info: req.Deployment
    // Access stage config: req.StageConfig
    // Access source code: req.TargetDeploymentSource
    
    // Example: Deploy to your platform
    err := p.deployToMyPlatform(ctx, req)
    if err != nil {
        return &deployment.ExecuteStageResponse{
            Status:  model.StageStatus_STAGE_FAILURE,
            Message: fmt.Sprintf("Deployment failed: %v", err),
        }, nil
    }
    
    return &deployment.ExecuteStageResponse{
        Status:  model.StageStatus_STAGE_SUCCESS,
        Message: "Deployment completed successfully",
    }, nil
}

func (p *MyDeploymentPlugin) deployToMyPlatform(
    ctx context.Context,
    req *deployment.ExecuteStageRequest,
) error {
    // Your platform-specific deployment logic
    // This could involve:
    // - API calls to your platform
    // - Executing CLI commands
    // - Processing manifests/configurations
    // - Waiting for deployment completion
    return nil
}

// Implement other execution methods...
func (p *MyDeploymentPlugin) executeRollback(ctx context.Context, req *deployment.ExecuteStageRequest) (*deployment.ExecuteStageResponse, error) {
    // Rollback implementation
    return &deployment.ExecuteStageResponse{Status: model.StageStatus_STAGE_SUCCESS}, nil
}

func (p *MyDeploymentPlugin) executeCleanup(ctx context.Context, req *deployment.ExecuteStageRequest) (*deployment.ExecuteStageResponse, error) {
    // Cleanup implementation
    return &deployment.ExecuteStageResponse{Status: model.StageStatus_STAGE_SUCCESS}, nil
}

func (p *MyDeploymentPlugin) executeQuickSync(ctx context.Context, req *deployment.ExecuteStageRequest) (*deployment.ExecuteStageResponse, error) {
    // Quick sync implementation
    return &deployment.ExecuteStageResponse{Status: model.StageStatus_STAGE_SUCCESS}, nil
}
```

## Advanced Plugin Development

### Configuration Handling

Plugins can receive configuration from both Piped config and Application config:

```go
// In the plugin configuration (Piped config)
type PluginConfig struct {
    APIEndpoint string `json:"apiEndpoint"`
    APIKey      string `json:"apiKey"`
    Timeout     string `json:"timeout"`
}

// In the application configuration
type AppConfig struct {
    ServiceName string            `json:"serviceName"`
    Environment string            `json:"environment"`
    Resources   map[string]string `json:"resources"`
}

func (p *MyDeploymentPlugin) parseConfigs(req *deployment.ExecuteStageRequest) (*PluginConfig, *AppConfig, error) {
    // Parse plugin configuration from Piped config
    pluginConfig := &PluginConfig{}
    if err := json.Unmarshal(req.PluginConfig, pluginConfig); err != nil {
        return nil, nil, fmt.Errorf("failed to parse plugin config: %w", err)
    }
    
    // Parse application configuration
    appConfig := &AppConfig{}
    if err := json.Unmarshal(req.TargetDeploymentSource.ApplicationConfig, appConfig); err != nil {
        return nil, nil, fmt.Errorf("failed to parse app config: %w", err)
    }
    
    return pluginConfig, appConfig, nil
}
```

### Error Handling and Logging

```go
import (
    "log/slog"
    "os"
)

func (p *MyDeploymentPlugin) ExecuteStage(
    ctx context.Context,
    req *deployment.ExecuteStageRequest,
) (*deployment.ExecuteStageResponse, error) {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
        "deployment_id", req.Deployment.Id,
        "stage_name", req.Stage.Name,
        "app_name", req.Deployment.ApplicationName,
    )
    
    logger.Info("Starting stage execution")
    
    // Your implementation here
    result, err := p.performDeployment(ctx, req)
    if err != nil {
        logger.Error("Stage execution failed", "error", err)
        return &deployment.ExecuteStageResponse{
            Status:  model.StageStatus_STAGE_FAILURE,
            Message: fmt.Sprintf("Stage failed: %v", err),
        }, nil
    }
    
    logger.Info("Stage execution completed successfully")
    return result, nil
}
```

### Working with Deploy Targets

```go
func (p *MyDeploymentPlugin) getDeployTarget(deploymentId string) (*DeployTarget, error) {
    // Deploy targets are configured in Piped config
    // You'll need to implement target selection logic
    // based on your plugin's requirements
    
    // Example: Select target based on deployment labels
    for _, target := range p.deployTargets {
        if p.matchesTarget(target, deploymentId) {
            return target, nil
        }
    }
    
    return nil, fmt.Errorf("no suitable deploy target found")
}
```

### Implementing LiveState Plugin

```go
type MyLivestatePlugin struct{}

func (p *MyLivestatePlugin) FetchDefinedResources(
    ctx context.Context,
    req *livestate.FetchDefinedResourcesRequest,
) (*livestate.FetchDefinedResourcesResponse, error) {
    // Return the types of resources this plugin can monitor
    return &livestate.FetchDefinedResourcesResponse{
        DefinedResources: []*livestate.DefinedResource{
            {
                Name: "MyService",
                Kind: "Service",
            },
            {
                Name: "MyDeployment", 
                Kind: "Deployment",
            },
        },
    }, nil
}

func (p *MyLivestatePlugin) GetResources(
    ctx context.Context,
    req *livestate.GetResourcesRequest,
) (*livestate.GetResourcesResponse, error) {
    // Fetch current state from your platform
    resources, err := p.fetchCurrentState(ctx, req)
    if err != nil {
        return nil, err
    }
    
    return &livestate.GetResourcesResponse{
        Resources: resources,
    }, nil
}

func (p *MyLivestatePlugin) fetchCurrentState(
    ctx context.Context,
    req *livestate.GetResourcesRequest,
) ([]*livestate.Resource, error) {
    // Implement platform-specific state fetching
    // This could involve API calls, CLI commands, etc.
    return nil, nil
}
```

## Stage-Only Plugins

For simple stage plugins (like Wait, Approval), you only need to implement the Deployment interface:

```go
type WaitPlugin struct{}

func (p *WaitPlugin) FetchDefinedStages(
    ctx context.Context,
    req *deployment.FetchDefinedStagesRequest,
) (*deployment.FetchDefinedStagesResponse, error) {
    return &deployment.FetchDefinedStagesResponse{
        Stages: []string{"WAIT"},
    }, nil
}

func (p *WaitPlugin) ExecuteStage(
    ctx context.Context,
    req *deployment.ExecuteStageRequest,
) (*deployment.ExecuteStageResponse, error) {
    if req.Stage.Name != "WAIT" {
        return &deployment.ExecuteStageResponse{
            Status:  model.StageStatus_STAGE_FAILURE,
            Message: "Unsupported stage",
        }, nil
    }
    
    // Parse wait duration from stage config
    config := &WaitConfig{}
    if err := json.Unmarshal(req.StageConfig, config); err != nil {
        return &deployment.ExecuteStageResponse{
            Status:  model.StageStatus_STAGE_FAILURE,
            Message: fmt.Sprintf("Invalid config: %v", err),
        }, nil
    }
    
    duration, err := time.ParseDuration(config.Duration)
    if err != nil {
        return &deployment.ExecuteStageResponse{
            Status:  model.StageStatus_STAGE_FAILURE,
            Message: fmt.Sprintf("Invalid duration: %v", err),
        }, nil
    }
    
    // Wait for the specified duration
    select {
    case <-ctx.Done():
        return &deployment.ExecuteStageResponse{
            Status:  model.StageStatus_STAGE_CANCELLED,
            Message: "Wait cancelled",
        }, nil
    case <-time.After(duration):
        return &deployment.ExecuteStageResponse{
            Status:  model.StageStatus_STAGE_SUCCESS,
            Message: fmt.Sprintf("Waited for %v", duration),
        }, nil
    }
}

type WaitConfig struct {
    Duration string `json:"duration"`
}

// Implement other required methods with appropriate responses
func (p *WaitPlugin) DetermineVersions(ctx context.Context, req *deployment.DetermineVersionsRequest) (*deployment.DetermineVersionsResponse, error) {
    return &deployment.DetermineVersionsResponse{}, nil
}

func (p *WaitPlugin) DetermineStrategy(ctx context.Context, req *deployment.DetermineStrategyRequest) (*deployment.DetermineStrategyResponse, error) {
    return &deployment.DetermineStrategyResponse{
        Unsupported: true,
    }, nil
}

func (p *WaitPlugin) BuildPipelineSyncStages(ctx context.Context, req *deployment.BuildPipelineSyncStagesRequest) (*deployment.BuildPipelineSyncStagesResponse, error) {
    return &deployment.BuildPipelineSyncStagesResponse{}, nil
}

func (p *WaitPlugin) BuildQuickSyncStages(ctx context.Context, req *deployment.BuildQuickSyncStagesRequest) (*deployment.BuildQuickSyncStagesResponse, error) {
    return &deployment.BuildQuickSyncStagesResponse{}, nil
}
```

## Testing Your Plugin

### Unit Testing

```go
package main

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
    "github.com/pipe-cd/pipecd/pkg/model"
)

func TestWaitPlugin_ExecuteStage(t *testing.T) {
    plugin := &WaitPlugin{}
    
    stageConfig := `{"duration": "100ms"}`
    
    req := &deployment.ExecuteStageRequest{
        Stage: &model.PipelineStage{
            Name: "WAIT",
        },
        StageConfig: []byte(stageConfig),
    }
    
    start := time.Now()
    resp, err := plugin.ExecuteStage(context.Background(), req)
    elapsed := time.Since(start)
    
    assert.NoError(t, err)
    assert.Equal(t, model.StageStatus_STAGE_SUCCESS, resp.Status)
    assert.True(t, elapsed >= 100*time.Millisecond)
    assert.True(t, elapsed < 200*time.Millisecond)
}

func TestWaitPlugin_FetchDefinedStages(t *testing.T) {
    plugin := &WaitPlugin{}
    
    resp, err := plugin.FetchDefinedStages(context.Background(), &deployment.FetchDefinedStagesRequest{})
    
    assert.NoError(t, err)
    assert.Equal(t, []string{"WAIT"}, resp.Stages)
}
```

### Integration Testing

```go
func TestPluginIntegration(t *testing.T) {
    // Start plugin as gRPC server
    plugin, err := sdk.NewPlugin(
        "test",
        sdk.WithDeploymentPlugin(&MyDeploymentPlugin{}),
    )
    assert.NoError(t, err)
    
    // Start server in background
    go func() {
        err := plugin.Run()
        assert.NoError(t, err)
    }()
    
    // Give server time to start
    time.Sleep(100 * time.Millisecond)
    
    // Create client and test
    conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
    assert.NoError(t, err)
    defer conn.Close()
    
    client := deployment.NewDeploymentServiceClient(conn)
    
    // Test your plugin methods
    resp, err := client.FetchDefinedStages(context.Background(), &deployment.FetchDefinedStagesRequest{})
    assert.NoError(t, err)
    assert.NotEmpty(t, resp.Stages)
}
```

## Building and Packaging

### Build Script

Create a `Makefile`:

```makefile
PLUGIN_NAME := my-plugin
VERSION := v1.0.0
BUILD_DIR := ./bin

.PHONY: build
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(PLUGIN_NAME) .

.PHONY: build-all
build-all:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(PLUGIN_NAME)_linux_amd64 .
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(PLUGIN_NAME)_linux_arm64 .
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(PLUGIN_NAME)_darwin_amd64 .
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(PLUGIN_NAME)_darwin_arm64 .
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(PLUGIN_NAME)_windows_amd64.exe .

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
```

### Distribution

1. **GitHub Releases**: Upload binaries to GitHub releases
2. **Container Images**: Package plugins in container images  
3. **Plugin Registry**: Submit to PipeCD plugin registry (future)

## Best Practices

### 1. Error Handling

- Always return proper status codes (SUCCESS, FAILURE, CANCELLED)
- Provide meaningful error messages
- Use structured logging
- Handle context cancellation gracefully

### 2. Configuration

- Validate all configuration at startup
- Provide clear error messages for invalid configuration
- Support environment variable substitution
- Document all configuration options

### 3. Resource Management

- Clean up resources on context cancellation
- Implement proper timeouts
- Monitor resource usage
- Use connection pooling for external APIs

### 4. Security

- Validate all inputs
- Use secure communication protocols
- Handle secrets securely
- Follow principle of least privilege

### 5. Performance

- Implement proper caching where appropriate
- Use streaming for large data transfers
- Optimize for common use cases
- Monitor and log performance metrics

### 6. Compatibility

- Version your plugin APIs
- Maintain backward compatibility when possible
- Test with different PipeCD versions
- Document version requirements

## Multi-Language Plugin Development

While the official SDK is Go-based, you can develop plugins in any language that supports gRPC:

### Python Plugin Example

```python
import grpc
from concurrent import futures
import time

# Generate from PipeCD proto files
import deployment_pb2
import deployment_pb2_grpc

class MyPythonPlugin(deployment_pb2_grpc.DeploymentServiceServicer):
    def FetchDefinedStages(self, request, context):
        return deployment_pb2.FetchDefinedStagesResponse(
            stages=["PYTHON_DEPLOY", "PYTHON_CLEANUP"]
        )
    
    def ExecuteStage(self, request, context):
        if request.stage.name == "PYTHON_DEPLOY":
            # Your deployment logic here
            return deployment_pb2.ExecuteStageResponse(
                status=deployment_pb2.STAGE_SUCCESS,
                message="Python deployment completed"
            )
        # Handle other stages...

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    deployment_pb2_grpc.add_DeploymentServiceServicer_to_server(
        MyPythonPlugin(), server
    )
    listen_addr = '0.0.0.0:8080'
    server.add_insecure_port(listen_addr)
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()
```

### Node.js Plugin Example

```javascript
const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');

// Load proto definition
const packageDefinition = protoLoader.loadSync('deployment.proto', {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
});

const deployment = grpc.loadPackageDefinition(packageDefinition).deployment;

class MyNodePlugin {
    fetchDefinedStages(call, callback) {
        callback(null, {
            stages: ['NODE_DEPLOY', 'NODE_CLEANUP']
        });
    }
    
    executeStage(call, callback) {
        const stageName = call.request.stage.name;
        
        if (stageName === 'NODE_DEPLOY') {
            // Your deployment logic here
            callback(null, {
                status: 'STAGE_SUCCESS',
                message: 'Node.js deployment completed'
            });
        }
        // Handle other stages...
    }
    
    // Implement other required methods...
}

const server = new grpc.Server();
server.addService(deployment.DeploymentService.service, new MyNodePlugin());
server.bindAsync('0.0.0.0:8080', grpc.ServerCredentials.createInsecure(), () => {
    server.start();
    console.log('Plugin server started on port 8080');
});
```

## Community and Support

### Contributing to the Ecosystem

1. **Share Your Plugins**: Contribute to the community plugins repository
2. **Documentation**: Write guides and tutorials
3. **Examples**: Provide real-world usage examples
4. **Feedback**: Report issues and suggest improvements

### Getting Help

- **Documentation**: [pipecd.dev/docs-v1](https://pipecd.dev/docs-v1)
- **Slack**: [#pipecd-plugin-dev](https://cloud-native.slack.com/)
- **GitHub Discussions**: [pipe-cd/pipecd](https://github.com/pipe-cd/pipecd/discussions)
- **Office Hours**: Monthly community calls

### Example Plugins

Study these example plugins for reference:

- **Kubernetes Plugin**: Full deployment plugin with progressive delivery
- **Wait Plugin**: Simple stage plugin implementation
- **ScriptRun Plugin**: Execute arbitrary commands and scripts
- **Community Examples**: Various platform integrations

## Conclusion

Plugin development for PipeCD v1 opens up endless possibilities for extending the platform's capabilities. Whether you're integrating with a new platform, adding custom deployment stages, or implementing specialized analysis tools, the plugin architecture provides a robust foundation for building powerful extensions.

Start with simple stage plugins to learn the concepts, then progress to full deployment plugins as you become more comfortable with the architecture. The community is here to help, and we're excited to see what you'll build!
