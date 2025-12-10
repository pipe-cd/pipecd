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

| Plugin | Latest Version | Release Date | Status |
|--------|----------------|--------------|--------|
| **[Kubernetes](kubernetes/)** | v0.3.0 | 2025-09-26 | Stable |
| **[Terraform](terraform/)** | v0.2.1 | 2025-10-09 | Stable |
| **[Wait](wait/)** | v0.1.1 | 2025-10-09 | Stable |
| **[Wait Approval](waitapproval/)** | v0.2.0 | 2025-10-09 | Stable |
| **[Analysis](analysis/)** | v0.1.1 | 2025-09-03 | Stable |
| **[Script Run](scriptrun/)** | v0.1.0 | 2025-09-04 | Stable |
| **Cloud Run** | No releases | - | In Development |
| **Kubernetes Multi-Cluster** | No releases | - | In Development |

> **Note:** Version information above is as of October 2024. For the most up-to-date information, run the version script or check [GitHub Releases](https://github.com/pipe-cd/pipecd/releases).

### Finding Plugin Versions

To check for new plugin versions:

1. **GitHub Releases:** Visit [releases page](https://github.com/pipe-cd/pipecd/releases) and filter by plugin tags(e.g., `pkg/app/pipedv1/plugin/kubernetes/v*`)

2. **API:** Query the GitHub API for plugin-specific releases:
 curl -s https://api.github.com/repos/pipe-cd/pipecd/releases | jq -I | select(tag_name |  startswith("pkg/app/pipedv1/plugin/"))

{{< alert title="Note" >}}
Plugin architecture with independent versioning is currently in alpha. Full independent release cycles are planned for PipeCD v1.0.
{{< /alert >}}

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
