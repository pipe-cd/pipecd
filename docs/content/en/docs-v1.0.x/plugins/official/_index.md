---
title: "Official Plugins"
linkTitle: "Official Plugins"
weight: 10
description: >
  Configuration guides for the official PipeCD plugins.
---

The PipeCD maintainers develop and maintain the following plugins. Each plugin is versioned and released independently. You can download the plugin binaries from the [releases page](https://github.com/pipe-cd/pipecd/releases).

## Deployment plugins

| Plugin | Description |
|--------|-------------|
| [Kubernetes](kubernetes/) | Deploys applications to a Kubernetes cluster. Supports quick sync and pipeline sync with canary, baseline, and blue-green strategies. |
| [Kubernetes multi-cluster](kubernetes-multicluster/) | Deploys a single application to multiple Kubernetes clusters with one pipeline. |
| [Terraform](terraform/) | Applies infrastructure changes by running `terraform plan` and `terraform apply` in a pipeline. |
| [Amazon ECS](ecs/) | Deploys applications to Amazon ECS. |

## Stage plugins

| Plugin | Stage | Description |
|--------|-------|-------------|
| [Wait](wait/) | `WAIT` | Waits for a specified duration before continuing the pipeline. |
| [Wait approval](wait-approval/) | `WAIT_APPROVAL` | Pauses the pipeline until a user approves the deployment. |
| [Analysis](analysis/) | `ANALYSIS` | Evaluates the deployment by querying metrics, logs, or HTTP endpoints. |
| [Script run](script-run/) | `SCRIPT_RUN` | Runs arbitrary commands as a pipeline stage. |
