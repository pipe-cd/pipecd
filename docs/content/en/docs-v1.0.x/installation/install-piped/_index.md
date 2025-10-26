---
title: "Install Piped"
linkTitle: "Install Piped"
weight: 3
description: >
  This page describes how you can run the `piped` binary connects your infrastructure to the PipeCD Control Plane.
---
 PipeCD V1 introduces a plugin-based architecture. In this model, platforms have been replaced with plugins. A plugin is responsible for the implementation and the logic of a deployment. For example, A kubernetes plugin handles deployment to Kubernetes clusters, An ECS plugin handles deployment to Amazon ECS Services and so on. The plugin-based architecture also unlocks the possibilities of creating custom plugins for various other platforms.

 If you are using a PipeCD V0.x.x, and want to switch to PipeCD V1, see [Migrating from V0 to V1]().

