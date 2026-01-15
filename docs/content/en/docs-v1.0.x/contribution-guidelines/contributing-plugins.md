---
title: "Contributing to plugins"
linkTitle: "Contributing to plugins"
weight: 4
description: >
  This page describes how to contribute plugins for piped.
---

PipeCD's plugin architecture allows anyone to extend piped's capabilities by creating custom plugins. This guide explains how to develop and contribute plugins.

## Understanding the plugin architecture

In PipeCD v1, plugins are the actors that execute deployments on behalf of piped. Instead of piped directly deploying to platforms, plugins handle platform-specific logic while piped's core controls deployment flows.

**Key concepts:**

- **Plugins** run as gRPC servers, launched and managed by piped
- **Deploy targets** define where a plugin deploys (e.g., a Kubernetes cluster)
- Plugins can be **official** (maintained by PipeCD team) or **community-contributed**

For a detailed overview, see the [Plugin Architecture blog post](https://pipecd.dev/blog/2024/11/28/overview-of-the-plan-for-pluginnable-pipecd/).

## Plugin types

Plugins can implement one or more of these interfaces:

| Interface | Purpose |
|-----------|---------|
| **Deployment** | Plan and execute deployment stages |
| **LiveState** | Fetch and build the state of live resources |
| **Drift** | Calculate drift between live and git-source manifests |

For example:
- A Kubernetes plugin implements all three interfaces
- A Wait stage plugin only implements the Deployment interface

## Where plugins live

- **Official plugins**: Located in `/pkg/app/pipedv1/plugin/` in the [pipecd repository](https://github.com/pipe-cd/pipecd)
- **Community plugins**: Located in the [pipe-cd/community-plugins](https://github.com/pipe-cd/community-plugins) repository

## Getting started

### Prerequisites

- [Go 1.24 or later](https://go.dev/)
- Understanding of gRPC
- Familiarity with the platform you're building a plugin for

### Study existing plugins

Before creating a new plugin, study the existing ones:

| Plugin | Complexity | Good for learning |
|--------|------------|-------------------|
| [wait](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin/wait) | Simple | Basic plugin structure |
| [waitapproval](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin/waitapproval) | Simple | Stage-only plugin |
| [kubernetes](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin/kubernetes) | Complex | Full-featured plugin |
| [terraform](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin/terraform) | Complex | Infrastructure as Code plugin |

Community plugins:

| Plugin | Description |
|--------|-------------|
| [opentofu](https://github.com/pipe-cd/community-plugins/tree/main/plugins/opentofu) | OpenTofu deployment plugin |

### Plugin structure

A minimal plugin needs:

```
your-plugin/
├── go.mod
├── go.sum
├── main.go           # Entry point, starts gRPC server
├── plugin.go         # Implements plugin interfaces
├── config/           # Plugin-specific configuration
│   └── application.go
└── README.md         # Documentation
```

### Plugin configuration

Plugins are configured in the piped config. See the [piped installation guide](/docs-v1.0.x/installation/install-piped/) for configuration examples:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  plugins:
    - name: your-plugin
      port: 7001              # Any unused port
      url: <PLUGIN_URL>
      deployTargets:          # Optional, depends on plugin
        - name: target1
          config:
            # Plugin-specific config
```

## Contributing to official plugins

1. **Open an issue** first to discuss your plugin idea with maintainers
2. **Fork and clone** the [pipecd repository](https://github.com/pipe-cd/pipecd)
3. **Create your plugin** under `/pkg/app/pipedv1/plugin/your-plugin/`
4. **Write tests** — see existing plugins for patterns
5. **Add a README** documenting configuration and usage
6. **Submit a PR** linking to the discussion issue

### Build and test

```bash
# Build all plugins
make build/plugin

# Run tests
make test/go

# Run piped locally with your plugin
make run/piped CONFIG_FILE=piped-config.yaml EXPERIMENTAL=true INSECURE=true
```

## Contributing to community plugins

The [community-plugins repository](https://github.com/pipe-cd/community-plugins) welcomes plugins that may not fit in the official repo.

1. **Fork** the community-plugins repository
2. **Create your plugin** following the structure above
3. **Submit a PR** with documentation

## Resources

- [Plugin Architecture RFC](https://github.com/pipe-cd/pipecd/blob/master/docs/rfcs/0015-pipecd-plugin-arch-meta.md)
- [Plugin Concepts](/docs-v1.0.x/concepts/#plugins)
- [Installing piped](/docs-v1.0.x/installation/install-piped/)
- [Plugin Alpha Release blog](https://pipecd.dev/blog/2025/06/16/plugin-architecture-piped-alpha-version-has-been-released/)
- [#pipecd Slack channel](https://cloud-native.slack.com/) for questions

Thank you for contributing to PipeCD plugins!
