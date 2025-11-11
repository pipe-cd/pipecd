---
title: "Install Piped"
linkTitle: "Install Piped"
weight: 3
description: >
  This page describes how you can run the `piped` binary connects your infrastructure to the PipeCD Control Plane.
---
 PipeCD V1 introduces a plugin-based architecture. A plugin is responsible for the implementation and the logic of a deployment. For example, A kubernetes plugin handles deployment to Kubernetes clusters, An ECS plugin handles deployment to Amazon ECS Services and so on.

 In this installation guide, we will see how you can configure your `piped` agent to connect to different plugins.

 >**NOTE:**
 >If you are using a PipeCD V0.x.x, and want to switch to PipeCD V1, see [Migrating from V0 to V1](../../migrating-from-v0-to-v1/_index.md).