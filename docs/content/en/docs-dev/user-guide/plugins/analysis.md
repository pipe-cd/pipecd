---
title: "Analysis"
linkTitle: "Analysis"
weight: 40
description: >
  Automated deployment analysis
---

The Analysis plugin enables PipeCD to automate deployment analysis by analyzing metrics and logs to determine if a deployment is successful or should be rolled back.

## Features

- **Metrics Analysis:** Analyze deployment impact using metrics
- **Log Analysis:** Inspect logs for errors and issues
- **Multiple Providers:** Support for Prometheus, Datadog, and more
- **Custom Queries:** Define custom analysis queries
- **Automated Decision:** Automatically approve/reject based on analysis results
- **Failure Detection:** Detect errors and anomalies during deployment

## Piped Configuration

Configure the Analysis plugin in your Piped configuration:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: dev
  pipedID: xxx
  plugins:
    - name: analysis
      port: 7004
      url: https://github.com/pipe-cd/pipecd/releases/download/...
```

## Application Configuration

Define analysis in `.pipe.yaml`:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-app
  labels:
    env: production
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
        manifests:
          - deployment.yaml
```

## Available Stages

- **ANALYSIS:** Run automated analysis on deployment metrics and logs

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-app
 
```

## Stage Configuration

### ANALYSIS

Analyze deployment metrics for specified duration.

```yaml
- name: ANALYSIS
  with:
    duration: 10m
```

## Examples

### Canary Deployment with Analysis

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: api-service
  labels:
    env: production
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 10%
      - name: ANALYSIS
        with:
          duration: 15m
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_CANARY_CLEAN
  plugins:
    kubernetes:
      input:
        manifests:
          - deployment.yaml
          - service.yaml
```

## Source Code

- [Analysis Plugin](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin/analysis)

## See Also

- [Configuration Reference](/docs-dev/user-guide/configuration-reference/)
- [Metrics Configuration](/)



