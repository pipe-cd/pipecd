---
title: "Install Piped"
linkTitle: "Install Piped"
weight: 3
description: >
  This page describes how you can run the `piped` binary connects your infrastructure to the PipeCD Control Plane.
---

Since Piped is a stateless agent, no database or storage is required to run. In addition, a Piped can interact with one or multiple platform providers, so the number of `piped`'s and where they should run is entirely up to your preference. For example, you can run your Pipeds in a Kubernetes cluster to deploy not just Kubernetes applications but your Terraform and Cloud Run applications as well.

In this guide, we will see how you can configure your `piped` agent and install it on different platforms.
