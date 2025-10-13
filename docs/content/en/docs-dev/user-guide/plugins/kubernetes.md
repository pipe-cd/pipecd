---
title: "Kubernetes"
linkTitle: "Kubernetes"
weight: 10
description: >
  Deploy applications to Kubernetes clusters
---

The Kubernetes plugin enables PipeCD to deploy and manage applications on Kubernetes clusters with support for progressive delivery strategies.

## Features

- **Progressive delivery:** Canary, blue-green, and analysis-based deployments
- **Automated rollback:** Automatic rollback on deployment failures
- **Live state sync:** Real-time synchronization with cluster state
- **Drift detection:** Detect configuration drift from Git
- **Multi-manifest support:** Helm charts, Kustomize, and plain YAML manifests
- **Traffic management:** Integration with service meshes (Istio, SMI)

## Piped Configuration

Configure the Kubernetes plugin in your Piped configuration:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  platforms:
    - name: kubernetes-default
      type: KUBERNETES
      config:
        masterURL: https://kubernetes.default.svc
        kubeConfigPath: /etc/kube/config
        appStateInformer:
          includeResources:
            - apiVersion: apps/v1
              kind: Deployment
```

## Application Configuration

Define your deployment pipeline in `.pipe.yaml`:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  input:
    helmChart:
      repository: https://charts.example.com
      name: myapp
      version: v1.0.0
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 10%
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 50%
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_CANARY_CLEAN
```

## Available Stages

- **K8S_SYNC:** Standard synchronization deployment
- **K8S_CANARY_ROLLOUT:** Deploy canary variant with traffic splitting
- **K8S_PRIMARY_ROLLOUT:** Promote canary to primary
- **K8S_CANARY_CLEAN:** Remove canary resources
- **K8S_BASELINE_ROLLOUT:** Deploy baseline for automated analysis
- **K8S_TRAFFIC_ROUTING:** Configure traffic routing rules

## Examples

### Quick Sync

```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  pipeline:
    stages:
      - name: K8S_SYNC
```

### Canary with Analysis

```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 20%
      - name: ANALYSIS
        with:
          duration: 10m
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_CANARY_CLEAN
```

## Source Code

- [`pkg/app/pipedv1/plugin/kubernetes/`](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin/kubernetes)

## See Also

- [Configuration Reference](/docs-dev/user-guide/configuration-reference/#kubernetes-application)
- [Managing Applications](/docs-dev/user-guide/managing-application/)
- [Examples](/docs-dev/user-guide/examples/)
