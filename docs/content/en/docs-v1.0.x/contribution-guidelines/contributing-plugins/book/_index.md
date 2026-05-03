---
title: "Plugin Development Book"
linkTitle: "Plugin Development Book"
weight: 1
description: >
  A hands-on guide to building your first PipeCD plugin.
---

Welcome to the Plugin Development Book! This guide is designed to take you from a basic understanding of PipeCD to building and testing your own custom plugin.

PipeCD v1 is built with extensibility in mind. By creating a plugin, you can add support for new platforms, custom deployment strategies, or specialized automation tasks.

## What we will build

In this book, we will walk through the creation of a simple **Stage Plugin**. This type of plugin is used to execute specific steps (stages) within a deployment pipeline.

By the end of this guide, you will have a working plugin that can:
1. Be registered with a Piped agent.
2. Execute a custom stage defined in an application's `app.pipecd.yaml`.
3. Report logs and status back to PipeCD.

## Chapters

1. **[Architecture](./01-architecture/)**: Understand how Piped and Plugins communicate via gRPC.
2. **[Your First Plugin](./02-first-stage-plugin/)**: Set up the project structure and implement the basic interface.
3. **[Configuration](./03-config/)**: Learn how to pass parameters from Git to your plugin.
4. **[Testing and Debugging](./04-testing/)**: Run your plugin locally and verify its behavior.

---

> [!TIP]
> This book focuses on the Go SDK. While plugins can be written in any language that supports gRPC, using the official [piped-plugin-sdk-go](https://github.com/pipe-cd/piped-plugin-sdk-go) is highly recommended as it handles most of the boilerplate for you.
