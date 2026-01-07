---
title: "Official Plugins"
linkTitle: "Plugins"
weight: 60
description: >
  Official PipeCD plugins for deploying to multiple platforms
---

PipeCD supports multiple deployment platforms through official plugins. Each plugin implements platform-specific deployment logic and integrates with the Piped agent.

## Available Plugins

| Plugin | Description | Status |
|--------|-------------|--------|
| **[Kubernetes](kubernetes/)** | Deploy applications to Kubernetes clusters | Stable |
| **[Kubernetes Multi-Cluster](kubernetes-multicluster/)** | Deploy to multiple Kubernetes clusters | Stable |
| **[Terraform](terraform/)** | Manage infrastructure with Terraform | Stable |
| **[Cloud Run](cloudrun/)** | Deploy to Google Cloud Run | Stable |
| **[Analysis](analysis/)** | Automated deployment analysis | Stable |
| **[Script Run](scriptrun/)** | Execute custom scripts during deployment | Stable |
| **[Wait](wait/)** | Add wait stages to pipelines | Stable |
| **[Wait Approval](waitapproval/)** | Manual approval gates in pipelines | Stable |


### Current Plugin Versions

Currently, all official plugins are **bundled with PipeCD** and share the same version as the Piped agent.

| Plugin | Version | Status |
|--------|---------|--------|
| **Kubernetes** | Bundled (v0.50.x) | Stable |
| **Terraform** | Bundled (v0.50.x) | Stable |
| **Cloud Run** | Bundled (v0.50.x) | Stable |
| **Analysis** | Bundled (v0.50.x) | Stable |
| **Script Run** | Bundled (v0.50.x) | Stable |
| **Wait** | Bundled (v0.50.x) | Stable |
| **Wait Approval** | Bundled (v0.50.x) | Stable |
| **Kubernetes Multi-Cluster** | Bundled (v0.50.x) | Stable |

> **Note:** Independent versioning for plugins is planned for future releases. For now, please refer to the [PipeCD Release Notes](https://github.com/pipe-cd/pipecd/releases) for updates on all components.

### Finding Plugin Versions

Since plugins are currently bundled:
1. Check the [PipeCD Release Notes](https://github.com/pipe-cd/pipecd/releases).
2. The plugin version corresponds to the PipeCD version (e.g., specific Piped version).

## Plugin Architecture

With the new pluggable architecture, PipeCD plugins:

- Run as separate gRPC servers
- Communicate with the Piped agent via the plugin SDK
- Implement the standard plugin interface
- Can be developed and deployed independently (planned for v1.0)

For more information, see:
- [Plugin Architecture Blog Post](/blog/plugin-arch-piped-alpha/)
- [Plugin Introduction](/blog/plugin-intro/)

## Getting Started

To use a plugin:

1. **Configure the plugin** in your Piped configuration file
2. **Create an application** that uses the plugin's platform kind
3. **Define the deployment pipeline** in your application's `.pipe.yaml`

See individual plugin documentation below for specific configuration options and examples.

## Plugin Development

Plugin source code is located in the PipeCD repository:
- **Path:** `pkg/app/pipedv1/plugin/<plugin-name>/`
- **SDK:** `pkg/plugin/sdk/`

For contributing to plugins or developing custom plugins, see the [Contributor Guide](https://github.com/pipe-cd/pipecd/blob/master/CONTRIBUTING.md).
