---
title: "Understanding Plugin Types"
weight: 6
description: >
  Deep dive into the types of plugins supported by the PipeCD SDK.
---

While PipeCD's core agent `Piped` interacts with plugins transparently, the official Go SDK categorizes plugins into distinct types to simplify development.

### 1. [StagePlugin](https://pkg.go.dev/github.com/pipe-cd/piped-plugin-sdk-go#StagePlugin)

A utility plugin that does not manage resources but provides utility stages to be executed within a deployment pipeline. An example is a plugin that provides a `WAIT` stage to pause the deployment pipeline for a defined duration.

### 2. [DeploymentPlugin](https://pkg.go.dev/github.com/pipe-cd/piped-plugin-sdk-go#DeploymentPlugin)

A plugin that manages actual target resources and performs state synchronization (syncing, diffing, rolling back). A Kubernetes plugin is a prime example.

In addition to core sync capabilities, a `DeploymentPlugin` can also implement the [LiveStatePlugin](https://pkg.go.dev/github.com/pipe-cd/piped-plugin-sdk-go#LiveStatePlugin) interface to report the live status of active resources on the target platform and render differences in the PipeCD Web UI.

When starting a plugin project, choose the type that fits your goals and implement its required interface. Since our `file` plugin manages local files as target resources, it will be implemented as a **`DeploymentPlugin`**.
