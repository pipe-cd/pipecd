---
title: "Introduction"
weight: 1
description: >
  Introduction to the PipeCD Pluggable Architecture and the goal of this book.
---

### About this Book

PipeCD is an open-source, GitOps-style Continuous Delivery platform. Recently, the PipeCD project introduced a new **Plugin Architecture** in its Alpha release to allow developers to flexibly extend PipeCD's deployment capabilities to support arbitrary platforms and tools.

In this book, we will learn how to build a custom PipeCD plugin from scratch by implementing a plugin named `file`. This plugin will allow PipeCD to manage and sync files on a local file system.

### Reference Links

Here are some useful reference links regarding PipeCD and its Plugin Architecture:

- [Official PipeCD Website](https://pipecd.dev/)
- [Official Plugin SDK for Go](https://pkg.go.dev/github.com/pipe-cd/piped-plugin-sdk-go)
- [Official Blog Post: Overview of the Plan for Pluggable PipeCD](https://pipecd.dev/blog/2024/11/28/overview-of-the-plan-for-pluginnable-pipecd/)
- [Japanese Article: Explaining the Plugin Architecture](https://zenn.dev/cadp/articles/pipecd-plugin-intro)

### Note on SDK Stability

This book is based on the PipeCD Plugin SDK version `v0.0.0-20250619080234-1ee9423d23c1` released during the Alpha stage. As pluggable architecture is continuously evolving, we recommend verifying Go interfaces, RPC definitions, and the `pipedv1` architecture against the latest default branch in the [`pipecd`](https://github.com/pipe-cd/pipecd) repository when implementing production-grade plugins.
