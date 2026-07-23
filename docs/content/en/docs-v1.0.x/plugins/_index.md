---
title: "Plugins"
linkTitle: "Plugins"
weight: 2
description: >
  This section describes the plugins available for PipeCD v1 and how to use them.
---

> **Note:**
> The Plugins section is a work in progress. More plugin docs are on the way. Happy PipeCDing!

| Plugin | Description | Status |
|--------|-------------|--------|
| [Kubernetes multi-cluster](kubernetes-multicluster/) | Deploy a single application to multiple Kubernetes clusters with one pipeline. | Alpha |
| [ECS](ecs/) | Deploy applications to Amazon ECS using the `EXTERNAL` deployment controller | Alpha |

This section contains configuration guides for the official PipeCD plugins.

- [Kubernetes](./kubernetes/)
- [Terraform](./terraform/)
- [Analysis](./analysis/)

In PipeCD v1, plugins handle deployments. `piped` runs each configured plugin as a separate process and communicates with it over gRPC, so which platforms your `piped` can deploy to depends on which plugins you configure. See more about [plugins](../concepts/#plugins).

There are two types of plugins:

- **Deployment plugins**: handle the deployment for a specific platform such as Kubernetes or Terraform.
- **Stage plugins**: provide pipeline stages that can be used with any deployment plugin, such as `WAIT` or `ANALYSIS`.

## Official plugins

The PipeCD maintainers develop and maintain the following plugins. Each plugin is versioned and released independently. You can download the plugin binaries from the [releases page](https://github.com/pipe-cd/pipecd/releases).

### Deployment plugins

| Plugin | Description |
|--------|-------------|
| Kubernetes | Deploys applications to a Kubernetes cluster. Supports quick sync and pipeline sync with canary, baseline, and blue-green strategies. |
| Kubernetes multi-cluster | Deploys a single application to multiple Kubernetes clusters with one pipeline. |
| Terraform | Applies infrastructure changes by running `terraform plan` and `terraform apply` in a pipeline. |
| Amazon ECS | Deploys applications to Amazon ECS. |

### Stage plugins

| Plugin | Stage | Description |
|--------|-------|-------------|
| Wait | `WAIT` | Waits for a specified duration before continuing the pipeline. |
| Wait approval | `WAIT_APPROVAL` | Pauses the pipeline until a user approves the deployment. |
| Analysis | `ANALYSIS` | Evaluates the deployment by querying metrics, logs, or HTTP endpoints. |
| Script run | `SCRIPT_RUN` | Runs arbitrary commands as a pipeline stage. |

## Community plugins

The PipeCD community maintains additional plugins in the [community-plugins repository](https://github.com/pipe-cd/community-plugins). Visit the repository for a list of available plugins and their documentation.

## Using a plugin

To add a plugin to your `piped` and register deploy targets, see [Configuring a plugin](../user-guide/managing-piped/configuring-a-plugin/).

## Writing your own plugin

Anyone can develop a plugin for PipeCD. See the [plugin development guide](../contribution-guidelines/contributing-plugins/) to get started.

