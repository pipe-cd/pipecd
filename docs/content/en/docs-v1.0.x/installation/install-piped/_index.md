---
title: "Install piped"
linkTitle: "Install piped"
weight: 3
description: >
  This page describes how you can run the `piped` binary that connects your infrastructure to the PipeCD Control Plane.
---

Since `piped` is a stateless agent, no database or storage is required to run. In addition, a `piped` can interact with one or multiple plugins, so the number of `piped` instances and where they should run is entirely up to your preference. For example, you can run your `piped` instances in a Kubernetes cluster to deploy not just Kubernetes applications but your Terraform and Cloud Run applications as well.

In this guide, we will see how you can configure your `piped` agent and install it on different platforms.
