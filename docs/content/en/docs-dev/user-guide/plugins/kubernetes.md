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

Configure the Kubernetes plugin in your Piped configuration (v1):

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: dev
  pipedID: your-piped-id
  pipedKeyFile: /etc/piped-secret/piped-key
  apiAddress: your-control-plane:443
  git:
    sshKeyFile: /etc/piped-secret/ssh-key
  repositories:
    - repoId: examples
      remote: https://github.com/your-org/examples.git
      branch: master
  plugins:
    - name: kubernetes
      port: 7001
      url: https://github.com/pipe-cd/pipecd/releases/download/pkg/app/pipedv1/plugin/kubernetes/v0.3.0/kubernetes_linux_amd64
      deployTargets:
        - name: production
          config:
            masterURL: https://kubernetes.default.svc
            kubeConfigPath: /etc/kube/config
            kubectlVersion: 1.32.2
            appStateInformer:
              includeResources:
                - apiVersion: apps/v1
                  kind: Deployment
        - name: staging
          config:
            masterURL: https://staging-cluster.example.com
            kubeConfigPath: /etc/kube/staging-config
            kubectlVersion: 1.32.2
```

## Application Configuration

Define your deployment pipeline in `.pipe.yaml` (v1):

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-k8s-app
  labels:
    env: production
    team: backend
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
  plugins:
    kubernetes:
      input:
        kubectlVersion: 1.32.2
        manifests:
          - deployment.yaml
          - service.yaml
```

### With Helm Chart

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-helm-app
  labels:
    env: production
  pipeline:
    stages:
      - name: K8S_SYNC
  plugins:
    kubernetes:
      input:
        kubectlVersion: 1.32.2
        helmChart:
          repository: https://charts.example.com
          name: myapp
          version: v1.0.0
```

## Available Stages

- **K8S_SYNC:** Standard synchronization deployment
- **K8S_CANARY_ROLLOUT:** Deploy canary variant with traffic splitting
- **K8S_PRIMARY_ROLLOUT:** Promote canary to primary
- **K8S_CANARY_CLEAN:** Remove canary resources
- **K8S_BASELINE_ROLLOUT:** Deploy baseline for automated analysis
- **K8S_TRAFFIC_ROUTING:** Configure traffic routing rules
- **WAIT:** Wait for a specified duration
- **WAIT_APPROVAL:** Manual approval gate
- **ANALYSIS:** Automated deployment analysis

## Examples

### Quick Sync

Deploy immediately with standard sync:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-k8s-app
  labels:
    env: production
  pipeline:
    stages:
      - name: K8S_SYNC
  plugins:
    kubernetes:
      input:
        manifests:
          - deployment.yaml
```

### Canary with Analysis

Progressive canary deployment with automated analysis:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-k8s-app
  labels:
    env: production
    team: platform
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
  plugins:
    kubernetes:
      input:
        kubectlVersion: 1.32.2
        manifests:
          - deployment.yaml
          - service.yaml
```

### Blue-Green Deployment

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-k8s-app
  labels:
    env: production
  pipeline:
    stages:
      - name: K8S_BASELINE_ROLLOUT
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 100%
      - name: WAIT
        with:
          duration: 5m
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_CANARY_CLEAN
  plugins:
    kubernetes:
      input:
        manifests:
          - deployment.yaml
```

## Source Code

- [`pkg/app/pipedv1/plugin/kubernetes/`](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin/kubernetes)

## See Also

- [Configuration Reference](/docs-dev/user-guide/configuration-reference/)
- [Managing Applications](/docs-dev/user-guide/managing-application/)
- [Migrating to PipeCD V1](/docs-dev/migrating-from-v0-to-v1/)