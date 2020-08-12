---
title: "Kubernetes"
linkTitle: "Kubernetes"
weight: 2
description: >
  Specific guides for Kubernetes applications.
---

> TBA

### Kubernetes Application Variant

Each Kubernetes application can has 3 variants: primary (aka stable), baseline, canary.
- `primary` runs the current version of code and configuration.
- `baseline` runs the same version of code and configuration as the primary variant. (Creating a brand new baseline workload ensures that the metrics produced are free of any effects caused by long-running processes.)
- `canary` runs the proposed changed of code or configuration.


Kubernetes Stages:

- `K8S_PRIMARY_ROLLOUT`
- `K8S_CANARY_ROLLOUT`
- `K8S_CANARY_CLEAN`
- `K8S_BASELINE_ROLLOUT`
- `K8S_BASELINE_CLEAN`
- `K8S_TRAFFIC_ROUTING`

Common Stages:
- `WAIT`
- `WAIT_APPROVAL`
- `ANALYSIS`
