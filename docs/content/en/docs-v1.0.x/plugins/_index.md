---
title: "Plugins"
linkTitle: "Plugins"
weight: 2
description: >
  This section describes the plugins available for PipeCD v1 and how to use them.
---

In PipeCD v1, plugins handle deployments. `piped` runs each configured plugin as a separate process and communicates with it over gRPC, so which platforms your `piped` can deploy to depends on which plugins you configure. See more about [plugins](../concepts/#plugins).

There are two types of plugins:

- **Deployment plugins**: handle the deployment for a specific platform such as Kubernetes or Terraform.
- **Stage plugins**: provide pipeline stages that can be used with any deployment plugin, such as `WAIT` or `ANALYSIS`.

## Official plugins

The PipeCD maintainers develop and maintain a set of official plugins, each versioned and released independently. See [Official Plugins](official/) for the full list and configuration guides.

## Community plugins

The PipeCD community maintains additional plugins in the [community-plugins repository](https://github.com/pipe-cd/community-plugins). Visit the repository for a list of available plugins and their documentation.

## Using a plugin

To add a plugin to your `piped` and register deploy targets, see [Configuring a plugin](../user-guide/managing-piped/configuring-a-plugin/).

## Writing your own plugin

Anyone can develop a plugin for PipeCD. See the [plugin development guide](../contribution-guidelines/contributing-plugins/) to get started.

