---
title: "Overview"
linkTitle: "Overview"
weight: 1
description: >
  Overview about deploying Kubernetes application.
---

> TBA


### Kubernetes Application Variant

Each Kubernetes application can has 3 variants: primary (aka stable), baseline, canary.
- `primary` runs the current version of code and configuration.
- `baseline` runs the same version of code and configuration as the primary variant. (Creating a brand new baseline workload ensures that the metrics produced are free of any effects caused by long-running processes.)
- `canary` runs the proposed changed of code or configuration.
