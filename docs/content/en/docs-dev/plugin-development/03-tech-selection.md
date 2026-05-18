---
title: "Technology Selection"
weight: 3
description: >
  Selecting the language and development tools for our plugin.
---

Next, let's select our development language and tools.

When developing plugins for PipeCD, **Go** is the most natural choice because PipeCD officially provides the [piped-plugin-sdk-go](https://pkg.go.dev/github.com/pipe-cd/piped-plugin-sdk-go). Therefore, we will build our plugin using Go.

To keep things educational and focused on understanding the core SDK, we will use the standard library as much as possible alongside the official SDK, avoiding any third-party dependencies.

The version of the SDK used in this book is `v0.0.0-20250619080234-1ee9423d23c1`. Since pluggable architecture is continuously evolving, always check the latest official documentation when developing production-ready plugins.
