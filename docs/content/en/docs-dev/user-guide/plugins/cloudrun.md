---
title: "Cloud Run"
linkTitle: "Cloud Run"
weight: 40
description: >
  Deploy serverless applications to Google Cloud Run
---

The Cloud Run plugin enables PipeCD to deploy serverless container applications to Google Cloud Run with progressive rollout capabilities.

## Features

- **Progressive rollouts:** Gradual traffic shifting
- **Revision management:** Automated revision handling
- **Traffic splitting:** Precise traffic control between revisions
- **Automated rollback:** Rollback on deployment failures
- **Multi-region support:** Deploy to multiple Cloud Run regions

## Piped Configuration

Configure the Cloud Run plugin in your Piped configuration:


```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  platforms:
    - name: cloudrun-default
      type: CLOUDRUN
      config:
        project: your-gcp-project
        region: us-central1
        credentialsFile: /etc/piped-secret/gcp-sa.json
```

## Application Configuration

Example `.pipe.yaml` for Cloud Run applications:


```yaml
apiVersion: pipecd.dev/v1beta1
kind: CloudRunApp
spec:
  input:
    serviceManifestFile: service.yaml
  pipeline:
    stages:
      - name: CLOUDRUN_PROMOTE
        with:
          percent: 10
      - name: WAIT
        with:
          duration: 5m
      - name: CLOUDRUN_PROMOTE
        with:
          percent: 100
```

## Available Stages

- **CLOUDRUN_SYNC:** Deploy new revision with 100% traffic
- **CLOUDRUN_PROMOTE:** Gradually promote new revision

## Examples

### Quick Sync

```yaml
apiVersion: pipecd.dev/v1beta1
kind: CloudRunApp
spec:
  pipeline:
    stages:
      - name: CLOUDRUN_SYNC
```

### Canary Rollout

```yaml
apiVersion: pipecd.dev/v1beta1
kind: CloudRunApp
spec:
  pipeline:
    stages:
      - name: CLOUDRUN_PROMOTE
        with:
          percent: 10
      - name: ANALYSIS
        with:
          duration: 10m
      - name: CLOUDRUN_PROMOTE
        with:
          percent: 100
```

## Source Code

- [`pkg/app/pipedv1/plugin/cloudrun/`](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin/cloudrun)

## See Also

- [Configuration Reference](/docs-dev/user-guide/configuration-reference/#cloudrun-application)
- [Managing Applications](/docs-dev/user-guide/managing-application/)
