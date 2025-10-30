---
title: "Wait"
linkTitle: "Wait"
weight: 60
description: >
  Add wait stages to pipelines
---

The Wait plugin enables PipeCD to add pause stages in deployment pipelines, useful for waiting between deployment phases, rate limiting, or coordinating with external processes.

## Features

- **Duration-based Waiting:** Pause deployment for specified duration
- **Simple Configuration:** Minimal configuration required
- **Pipeline Control:** Coordinate deployment phases
- **Rate Limiting:** Space out deployment operations
- **Integration Friendly:** Works with other plugins

## Piped Configuration

Configure the Wait plugin in your Piped configuration:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: dev
  pipedID: xxx
  plugins:
    - name: wait
      port: 7002
      url: https://github.com/pipe-cd/pipecd/releases/download/...
```

## Application Configuration

Add wait stages in `.pipe.yaml`:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-app
  labels:
    env: production
  pipeline:
    stages:
      - name: WAIT
        with:
          duration: 30s
  plugins: {}
```

## Available Stages

- **WAIT:** Wait for specified duration before proceeding to next stage

## Stage Configuration

### WAIT

Pause deployment for specified duration.

```yaml
- name: WAIT
  with:
    duration: 5m
```

## Duration Format

Duration values support standard Go duration format:
- Seconds: `30s`, `1m`
- Minutes: `5m`, `10m`
- Hours: `1h`
- Combined: `1h30m`, `2h15m30s`

## Examples

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
          replicas: 10%
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

### Canary Deployment with Wait

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
      - name: WAIT
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

### Progressive Rollout with Multiple Waits

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: web-service
  labels:
    env: production
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 20%
      - name: WAIT
        with:
          duration: 5m
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 50%
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

- [Wait Plugin](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin/wait)

## See Also

- [Configuration Reference](/docs-dev/user-guide/configuration-reference/)
