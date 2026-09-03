---
title: "Next steps: testing and publishing your plugin"
linkTitle: "Next steps: testing and publishing"
weight: 9
description: >
  Test the plugin further, publish it, and review what you built.
---

Your plugin is complete, and you have watched it run a real deployment through a local `piped`. This final chapter covers where to go from here: testing the plugin more thoroughly, publishing it so others can use it, and a short recap of what you built.

## Test your plugin more thoroughly

The tests you wrote cover the file helpers, which hold the real logic. The stage methods themselves were left untested because they mostly wire those helpers together and write logs. As a plugin grows, its methods take on more logic, and testing them pays off.

Methods such as `DetermineVersions` and the stage builders read from the deployment source, so testing them means constructing an `ApplicationConfig`. The SDK provides a helper for exactly this:

```go
func LoadApplicationConfigForTest[Spec any](t *testing.T, filename string, pluginName string) *ApplicationConfig[Spec]
```

It loads an application configuration from a YAML file, so you can keep example `app.pipecd.yaml` files as test fixtures and load them in your tests rather than building the structs by hand. The real plugins in the `pipecd` repository, under [`pkg/app/pipedv1/plugin`](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin), show how the maintainers test theirs. The `kubernetes`, `terraform`, and `wait` plugins are good references, and `wait` is the smallest.

## Publish your plugin

A plugin is a binary that `piped` downloads, so sharing one means making its binary available. The `url` field in the `piped` configuration accepts an `https://` release URL as well as a local `file://` path, so you can build a release binary, attach it to a Git release, and point other `piped` instances at it.

To share a plugin with the wider community, PipeCD has a dedicated home for community-built plugins: the [community-plugins repository](https://github.com/pipe-cd/community-plugins). Follow its contribution guide to add yours. Publishing there makes the plugin discoverable to other PipeCD users and puts it alongside the other community plugins.

## What we covered

Starting from an empty directory, you built a deployment plugin that:

- defines its configuration types and registers with `piped` through the SDK,
- reports the deployed version and the stages it provides,
- plans a deployment for both pipeline sync and quick sync, and
- runs the `FILE_DIFF`, `FILE_SYNC`, and `FILE_ROLLBACK` stages against a real `piped`.

That covers every method of the DeploymentPlugin interface, backed by tested helpers - a complete example to build your own plugin from.
