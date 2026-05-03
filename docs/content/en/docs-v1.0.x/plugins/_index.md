---
title: "Plugins"
linkTitle: "Plugins"
weight: 2
description: >
  Extend PipeCD's deployment capabilities with a modular, gRPC-based plugin architecture.
---

# Introduction to PipeCD Plugins

In **PipeCD v1**, the core engine is designed to be platform-agnostic. This means that instead of building every deployment provider (like Kubernetes, Terraform, or Lambda) directly into the core code, PipeCD uses a **modular plugin architecture**.

## What is a PipeCD Plugin?

A plugin is a standalone binary that implements the `pipedv1/plugin` interface. It communicates with the **Piped agent** over **gRPC**. This allows the community to build and maintain support for new deployment platforms without modifying the PipeCD core.

## Why the Plugin-First Approach?

- **Decoupling**: Platform-specific logic is isolated from the core GitOps engine.
- **Extensibility**: Easily add support for custom or proprietary deployment workflows.
- **Stability**: Core updates don't break platform-specific logic, and vice versa.

---

## 📘 Plugin Development Book (In Progress)

The Plugin Development Book is currently being translated from Japanese into English
and will be hosted here within PipeCD docs. It guides contributors through building
a complete PipeCD v1 plugin from scratch — from project setup to testing with `piped`.

### What the Book Covers

| Chapter | Topic |
|---------|-------|
| 01 | Introduction |
| 02 | Plugin functionality |
| 03 | Technology selection |
| 04 | Project initialization |
| 05 | Adding dependencies |
| 06 | Plugin types |
| 07 | First step in plugin implementation |
| 08 | Satisfying the `DeploymentPlugin` interface |
| 09 | Defining configuration types |
| 10 | Temporarily satisfying the interface |
| 11 | Implementing `FetchDefinedStages` |
| 12 | Implementing `DetermineVersions` |
| 13 | Implementing `DetermineStrategy` |
| 14 | Implementing `BuildPipelineSyncStages` |
| 15 | Implementing `BuildQuickSyncStages` |
| 16 | Introducing `ExecuteStage` |
| 17 | `ExecuteStage`: DIFF stage |
| 18 | `ExecuteStage`: SYNC stage |
| 19 | `ExecuteStage`: ROLLBACK stage |
| 20 | Modifying `main` |
| 21 | Testing with `piped` |
| 22 | Conclusion |

> **Note:** This book was originally written based on PipeCD's Alpha Release (June 2025).
> Always refer to the latest PipeCD documentation alongside it.

Follow the progress or contribute: [tracking issue #6679](https://github.com/pipe-cd/pipecd/issues/6679)

> **Interested in contributing?** Join the discussion in the `#pipecd` CNCF Slack channel!