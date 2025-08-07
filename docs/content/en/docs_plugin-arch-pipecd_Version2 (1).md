# PipeCD v1 (Plugin Architecture) Documentation

## Table of Contents

- [Introduction](#introduction)
- [Motivation](#motivation)
- [Architecture Overview](#architecture-overview)
- [Plugin Lifecycle](#plugin-lifecycle)
- [Configuration Example](#configuration-example)
- [Developing Plugins](#developing-plugins)
- [Protocol and Interfaces](#protocol-and-interfaces)
- [Migration from Pipedv0](#migration-from-pipedv0)
- [Advantages](#advantages)
- [References & Further Reading](#references--further-reading)

---

## Introduction

PipeCD v1 introduces a **plugin architecture** that enables support for more platforms and custom deployment behaviors. In this model, plugins are external actors that execute deployments, providing extensibility and allowing the community to implement and share new deployment strategies.

See: [Overview of the Plan for Pluginnable PipeCD](https://pipecd.dev/blog/2024/11/28/overview-of-the-plan-for-pluginnable-pipecd/)

---

## Motivation

- PipeCD aims to be a unified, progressive delivery tool for any application platform.
- Supports Kubernetes, ECS, Terraform, Lambda, Cloud Run out of the box.
- The plugin model allows support for **any platform** (e.g., Azure, Cloudflare, CloudFormation/CDK) via external plugins.
- Enables multiple versions/implementations of deployment logic for the same platform.

See RFC: [`0015-pipecd-plugin-arch-meta.md`](https://github.com/pipe-cd/pipecd/blob/master/docs/rfcs/0015-pipecd-plugin-arch-meta.md)

---

## Architecture Overview

In PipeCD v1 (plugin-arch):
- The **Piped Agent** loads and manages plugins.
- **Plugins** are binaries (not built-in), fetched from external sources (e.g., GitHub).
- **Control Plane** orchestrates deployments via Pipeds and plugins.

**Comparison:**
- _Previous (v0)_: Piped directly handled deployments for each platform.
- _Plugin-arch (v1)_: Piped delegates deployment execution to plugins via gRPC.

![Current and Pluginnable PipeCD mechanism](docs/static/images/plugin-intro-mechanism-new.drawio.png)

---

## Plugin Lifecycle

- Plugins are external binaries (can be built by anyone).
- On startup, Piped loads configured plugins, launches each as a gRPC server.
- During deployment, Piped communicates with plugins via gRPC.

**Plugin Source Example:**

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  plugin:
    - name: k8s_plugin
      sourceURL: https://github.com/org/k8s-plugin
      port: 8081
      deployTargets:
        - name: dev
          labels:
            env: dev
          config: # plugin-specific
```

See RFC: [`0015-pipecd-plugin-arch-meta.md`](https://github.com/pipe-cd/pipecd/blob/master/docs/rfcs/0015-pipecd-plugin-arch-meta.md)

---

## Configuration Example

A plugin is registered in the Piped config YAML under the `plugin:` section.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  plugin:
    - name: custom_plugin
      sourceURL: https://github.com/myorg/custom-plugin
      port: 8082
```

- Plugins can be loaded from HTTP URLs or local paths.
- Multiple plugins can be loaded; each can handle different platforms/stages.

---

## Developing Plugins

**Anyone can develop a plugin.**

- Main task: implement deployment logic for the platform/stage.
- Piped core manages deployment flows (you don't need to implement GitOps mechanics).
- Plugins can be written in any language supporting gRPC.
- Plugins can be published independently; no need to merge with PipeCD repo.

Plugin developer resources (coming soon):
- [Plugin Development Guide](https://pipecd.dev/blog/2024/11/28/overview-of-the-plan-for-pluginnable-pipecd/#how-to-develop-a-plugin)
- Example plugins: Kubernetes, CloudRun, Terraform, ECS, Lambda.

---

## Protocol and Interfaces

**Protocol:**  
- Piped communicates with plugins using **gRPC** (chosen for maintainability; consistent with other PipeCD components).

**Plugin interface requirements:**  
- Implement gRPC service to receive deployment instructions, report status/results.
- Provide hooks for additional features (Plan Preview, Drift Detection, Analysis stage).

See RFC section: [Protocol](https://github.com/pipe-cd/pipecd/blob/master/docs/rfcs/0015-pipecd-plugin-arch-meta.md#the-protocol)

---

## Migration from Pipedv0

- Pipedv0 (built-in platform support) will be supported until end of 2025.
- Single control plane supports both v0 and v1 during transition.
- Migration involves updating Piped configuration and application config to register plugins instead of platform/kind.
- Official plugins will be provided for existing platforms (Kubernetes, etc.).

See RFC section: [Migration process](https://github.com/pipe-cd/pipecd/blob/master/docs/rfcs/0015-pipecd-plugin-arch-meta.md#migration-process)

---

## Advantages

- **Platform Extensibility:** Deploy to any platform via plugin.
- **Custom Stages:** Add custom behaviors (e.g., analysis, jobs) via plugins.
- **Community Contributions:** Use plugins by others, or publish your own.
- **Multiple implementations:** Use built-in or custom plugins for same platform.

See: [Advantages of Pluginnable PipeCD](https://github.com/pipe-cd/pipecd/blob/master/docs/content/en/blog/plugin-intro.md#advantages-of-pluginnable-pipecd)

---

## References & Further Reading

- [RFC: Plugin Architecture Meta](https://github.com/pipe-cd/pipecd/blob/master/docs/rfcs/0015-pipecd-plugin-arch-meta.md)
- [Blog: Plugin Architecture Alpha](https://github.com/pipe-cd/pipecd/blob/master/docs/content/en/blog/plugin-arch-piped-alpha.md)
- [Blog: Plugin Architecture Intro](https://github.com/pipe-cd/pipecd/blob/master/docs/content/en/blog/plugin-intro.md)
- [Pipedv1 README](https://github.com/pipe-cd/pipecd/blob/master/cmd/pipedv1/README.md)
- [PipeCD Website](https://pipecd.dev)
- [PipeCD Community](https://cloud-native.slack.com/archives/C01B27F9T0X)

---

_This document is based on official RFCs, blog posts, and code from the PipeCD repository. For updates, join the PipeCD community or see [issue #6077](https://github.com/pipe-cd/pipecd/issues/6077)._