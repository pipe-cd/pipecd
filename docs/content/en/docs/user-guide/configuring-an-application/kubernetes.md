---
title: "Kubernetes"
linkTitle: "Kubernetes"
weight: 1
description: >
  Specific guide for configuring Kubernetes applications.
---

> TBA

## Quick Sync

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: Kubernetes
spec:
  input:
    helmChart:
      repository: pipecd
      name: helloworld
      version: v0.3.0
```

## Sync with Specified Pipeline

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


### Canary

> TBA

Canary deployment for a non-mesh Kubernetes application is done by manipulating pod selector in the Service resource.

### Canary with Istio

> TBA

### Canary with SMI

> TBA

### BlueGreen

> TBA

### BlueGreen with Istio

> TBA

### BlueGreen with SMI

> TBA

## Manifest Templating

> TBA

### Plain YAML

> TBA

### Helm

> TBA

### Kustomize

> TBA

## Reference

Go to [Configuration Reference](/docs/user-guide/configuration-reference/#kubernetes-application) to see the full configuration.
