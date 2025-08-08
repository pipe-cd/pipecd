---
title: "PipeCD v1 Plugin Architecture Overview"
linkTitle: "PipeCD v1 Overview"
weight: 10
description: >
  Complete overview of PipeCD v1 plugin architecture, migration guide, and key changes from v0
---

# PipeCD v1: Plugin Architecture Overview

PipeCD v1 introduces a revolutionary plugin-based architecture that fundamentally transforms how deployments are executed and managed. This document provides a comprehensive overview of the new architecture, its benefits, and migration considerations.

## What's New in PipeCD v1?

### Plugin-Based Architecture

The most significant change in PipeCD v1 is the shift from a monolithic Piped agent to a plugin-based system:

- **v0 Architecture**: Piped agent directly executes deployments for all supported platforms
- **v1 Architecture**: Plugins execute deployments while Piped core orchestrates the process

![Plugin Architecture Diagram](/images/plugin-arch-overview.png)

### Key Benefits

1. **Extensibility**: Deploy to any platform by developing or using community plugins
2. **Flexibility**: Choose different deployment strategies for different applications
3. **Community-Driven**: Leverage plugins developed by the community
4. **Maintainability**: Smaller, focused plugin codebases are easier to maintain
5. **Multi-Language Support**: Plugins can be written in any language that supports gRPC

## Architecture Components

### Piped Core (v1)

The Piped core in v1 is responsible for:
- Orchestrating deployment workflows
- Managing plugin lifecycle
- Handling Git operations
- Coordinating communication between plugins
- Managing configuration and secrets

### Plugins

Plugins are independent processes that handle platform-specific operations:
- Run as gRPC servers
- Implement specific interfaces (Deployment, LiveState)
- Can be written in any programming language
- Loaded dynamically at Piped startup

### Communication Protocol

- **Protocol**: gRPC over local network
- **Security**: Local communication within Piped environment
- **Interfaces**: Standardized plugin APIs for consistency

## Supported Plugin Types

### 1. Deployment Plugins
Handle application deployment logic:
- **Kubernetes**: Kubernetes deployments with progressive delivery
- **Terraform**: Infrastructure as Code deployments
- **ECS**: Amazon ECS service deployments
- **Lambda**: AWS Lambda function deployments
- **CloudRun**: Google Cloud Run service deployments

### 2. Stage Plugins
Provide specialized deployment stages:
- **Wait**: Time-based waiting stages
- **WaitApproval**: Manual approval gates
- **ScriptRun**: Custom script execution
- **Analysis**: Automated deployment analysis

### 3. Custom Plugins
Community and custom-developed plugins for:
- Platform-specific deployments
- Custom analysis and validation
- Integration with third-party tools

## Configuration Structure

### Piped Configuration

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: my-project
  pipedID: my-piped-id
  pipedKeyData: <base64-encoded-key>
  apiAddress: control-plane.example.com:8080
  
  plugins:
    - name: kubernetes
      port: 7001
      url: https://github.com/pipe-cd/pipecd/releases/download/v1.0.0/plugin_kubernetes_linux_amd64
      deployTargets:
        - name: dev-cluster
          config:
            masterURL: https://dev-k8s.example.com
            kubeConfigPath: /path/to/kubeconfig
        - name: prod-cluster
          config:
            masterURL: https://prod-k8s.example.com
            kubeConfigPath: /path/to/prod-kubeconfig
    
    - name: wait
      port: 7002
      url: https://github.com/pipe-cd/pipecd/releases/download/v1.0.0/plugin_wait_linux_amd64
```

### Application Configuration

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
metadata:
  labels:
    kind: KUBERNETES
spec:
  name: my-app
  deployTargets:
    - dev-cluster
  
  plugins:
    kubernetes:
      input:
        manifests:
          - deployment.yaml
          - service.yaml
        namespace: my-app
      service:
        kind: Service
        name: my-app
  
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 10%
      - name: WAIT
        with:
          duration: 30s
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_CANARY_CLEAN
```

## Migration from v0 to v1

### Compatibility Period

- **v0 Support**: Maintained until end of 2025
- **Coexistence**: Control plane supports both v0 and v1 Piped agents
- **Migration Timeline**: Flexible, project-by-project migration

### Key Changes

#### Platform Providers → Deploy Targets
```yaml
# v0 Configuration
platformProviders:
  - name: dev-k8s
    type: KUBERNETES
    config:
      masterURL: https://dev-k8s.example.com

# v1 Configuration
plugins:
  - name: kubernetes
    deployTargets:
      - name: dev-k8s
        config:
          masterURL: https://dev-k8s.example.com
```

#### Application Kind → Labels
```yaml
# v0 Configuration
spec:
  kind: KUBERNETES

# v1 Configuration
metadata:
  labels:
    kind: KUBERNETES
```

### Migration Steps

1. **Update Control Plane**: Ensure v0.52.0 or later
2. **Prepare v1 Configuration**: Convert existing configs to v1 format
3. **Test v1 Piped**: Deploy v1 Piped alongside existing v0 Piped
4. **Migrate Applications**: Update application configurations incrementally
5. **Complete Migration**: Remove v0 Piped agents

## Plugin Lifecycle

### Startup Process

1. **Configuration Loading**: Piped reads plugin configurations
2. **Plugin Download**: Fetch plugin binaries from specified URLs
3. **Plugin Launch**: Start plugins as gRPC servers
4. **Health Checks**: Verify plugin connectivity and readiness
5. **Registration**: Register available stages and capabilities

### Runtime Management

- **Health Monitoring**: Continuous health checks for plugins
- **Failure Recovery**: Automatic plugin restart on failure
- **Resource Management**: Monitor plugin resource usage
- **Graceful Shutdown**: Clean plugin termination on Piped shutdown

## Development and Community

### Official Plugins

Maintained by the PipeCD team:
- Kubernetes deployment plugin
- Terraform deployment plugin
- Wait and WaitApproval stage plugins
- ScriptRun execution plugin

### Community Plugins

Third-party plugins available via:
- [Community Plugins Repository](https://github.com/pipe-cd/community-plugins)
- GitHub releases and registries
- Custom repositories

### Plugin Development

- **SDK Available**: Go SDK for plugin development
- **Multi-Language Support**: gRPC enables any language
- **Documentation**: Comprehensive developer guides
- **Examples**: Reference implementations available

## Performance and Scalability

### Resource Usage

- **Memory**: Plugins run in separate processes, isolated memory usage
- **CPU**: Parallel plugin execution for better performance
- **Network**: Local gRPC communication, minimal overhead

### Scaling Considerations

- **Plugin Limits**: Practical limits based on available ports and resources
- **Concurrent Deployments**: Multiple plugins can execute simultaneously
- **Resource Isolation**: Plugin failures don't affect Piped core

## Security Considerations

### Plugin Security

- **Local Communication**: gRPC over localhost, no external exposure
- **Process Isolation**: Plugins run in separate processes
- **Permission Model**: Plugins inherit Piped permissions
- **Binary Verification**: Plugin integrity checks (future enhancement)

### Access Control

- **Deploy Target Isolation**: Plugins can only access configured targets
- **Configuration Scope**: Plugin-specific configuration isolation
- **Secret Management**: Centralized secret handling by Piped core

## Troubleshooting

### Common Issues

1. **Plugin Download Failures**: Check network connectivity and URL validity
2. **Port Conflicts**: Ensure unique ports for each plugin
3. **Configuration Errors**: Validate plugin-specific configuration
4. **Health Check Failures**: Verify plugin startup and gRPC server status

### Debugging Tools

- **Plugin Logs**: Individual plugin log files
- **Health Endpoints**: Plugin health check endpoints
- **Configuration Validation**: Built-in configuration validation
- **Debug Mode**: Verbose logging for troubleshooting

## Future Roadmap

### Planned Features

- **Plugin Registry**: Centralized plugin discovery and management
- **Binary Verification**: Plugin signature verification
- **Hot Reload**: Dynamic plugin loading without Piped restart
- **Resource Limits**: Plugin resource consumption controls

### Community Integration

- **Plugin Marketplace**: Community plugin sharing platform
- **Certification Program**: Verified plugin certification
- **Developer Tools**: Enhanced SDK and development tools
- **Integration Guides**: Platform-specific integration guides

## Getting Started

### Quick Start

1. **Download v1 Piped**: Get the latest v1 Piped binary
2. **Configure Plugins**: Set up basic plugin configuration
3. **Deploy Sample App**: Try the hello-world example
4. **Explore Examples**: Review example configurations and plugins

### Next Steps

- [Plugin Development Guide](../plugin-development/)
- [Configuration Reference](../configuration/)
- [Migration Guide](../migration/)
- [Troubleshooting Guide](../troubleshooting/)

## Support and Community

- **Documentation**: [pipecd.dev/docs-v1](https://pipecd.dev/docs-v1)
- **Slack Channel**: [#pipecd](https://cloud-native.slack.com/archives/C01B27F9T0X)
- **GitHub Issues**: [pipe-cd/pipecd](https://github.com/pipe-cd/pipecd/issues)
- **Community Meetings**: Monthly development and user meetings
