---
title: "Introduction"
linkTitle: "Introduction"
weight: 1
description: >
  What you build in this tutorial and how the chapters fit together.
---

You start with an empty directory and finish with a simple plugin that `piped` can load and run.

## What you build

The plugin copies a directory of files from a Git repository onto the machine where `piped` runs. This is a common way to deploy to virtual machines and bare-metal hosts.

The plugin implements two deployment stages:

- `FILE_DIFF` compares the files in the Git repository against the files already on the machine and prints the difference to the deployment log, without changing anything.
- `FILE_SYNC` applies the change. It copies the files that exist in the Git repository and removes any files on the machine that are no longer in Git.

The plugin stays small and uses only the Go standard library and the plugin SDK.

## Before you start

You should be comfortable with Go and know basic gRPC, since a plugin runs as a gRPC server that `piped` talks to. You do not need to know PipeCD's internals. If you are new to plugins, read the [Plugins concepts](/docs-v1.0.x/concepts/#plugins) first and then the [plugin architecture blog](https://pipecd.dev/blog/2024/11/28/overview-of-the-plan-for-pluginnable-pipecd/) which explains the design.

## How the chapters fit together

Each chapter builds on the previous one:

- Chapters 2 and 3 set up the project and define the configuration types.
- Chapters 4 and 5 implement the methods that plan a deployment.
- Chapters 6 and 7 implement the methods that run each stage.
- Chapters 8 and 9 connect the plugin to `piped`, run it, and cover where to go next.

## A note on versions

The plugin SDK ([`github.com/pipe-cd/piped-plugin-sdk-go`](https://pkg.go.dev/github.com/pipe-cd/piped-plugin-sdk-go)) still changes between releases. The code here tracks the SDK version this documentation targets. If a snippet does not compile, trust the [`pipecd`](https://github.com/pipe-cd/pipecd) source and SDK reference over the text.
