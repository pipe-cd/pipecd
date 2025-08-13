# Plugin Architecture

PipeCD v1 introduces a revolutionary plugin-based architecture that transforms PipeCD from a monolithic continuous delivery system into an extensible, plugin-based platform.

---

## Overview

PipeCD v1 transforms PipeCD into an extensible, plugin-based platform that can support any deployment target via custom plugins. This architecture enables PipeCD to achieve its vision of becoming "**The One CD for All {applications, platforms, operations}**".

### Key Benefits

- **Extensibility**: Add support for new deployment platforms through custom plugins
- **Modularity**: Plugins implement specific interfaces (Deployment, LiveState, Drift)
- **Community-Driven**: Plugin ecosystem developed and maintained by the community
- **Backward Compatibility**: Seamless coexistence with PipeCD v0 installations
- **Performance**: Efficient gRPC-based communication between piped core and plugins

### Current Status

- **Target Release**: February 2025
- **Plugin SDK**: Available and documented
- **Built-in Plugins**: Kubernetes, Terraform, Cloud Run, ECS, Lambda
- **Migration Support**: Tools and documentation available for v0 to v1 migration

### Architecture Changes

The plugin architecture introduces key changes:

- **Platform Providers** → **Deploy Targets**: Configuration moved to plugin-specific deploy targets
- **Application Kind** → **Labels**: Application types now specified as `metadata.labels.kind`
- **Plugin System**: gRPC-based plugins replace built-in platform implementations

For comprehensive information and detailed examples, refer to the complete documentation:

:::info
**Full Documentation**: [Plugin Architecture for PipeCD v1](/docs-dev/plugin-architecture/)
:::

## Quick Examples

### Piped Configuration (v1)
```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  plugins:
  - name: kubernetes
    port: 7003
    url: file:///usr/local/bin/kubernetes-plugin
    deployTargets:
    - name: production
      labels:
        env: prod
      config:
        masterURL: https://k8s-prod.company.com
        kubeConfigPath: /etc/kubeconfig
```

### Application Configuration (v1)
```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
metadata:
  labels:
    kind: KUBERNETES
spec:
  name: web-api
  deployTargets:
    - production
  plugins:
    kubernetes:
      input:
        namespace: production
        manifests:
          - k8s/
```
