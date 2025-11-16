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

Configure the Cloud Run plugin in your Piped configuration (v1):
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
    - name: cloudrun
      port: 7003
      url: https://github.com/pipe-cd/pipecd/releases/download/pkg/app/pipedv1/plugin/cloudrun/v0.1.0/cloudrun_linux_amd64
      deployTargets:
        - name: production
          config:
            project: your-gcp-project
            region: us-central1
            credentialsFile: /etc/piped-secret/gcp-sa.json
        - name: staging
          config:
            project: your-gcp-project
            region: us-west1
            credentialsFile: /etc/piped-secret/gcp-sa.json
```

## Application Configuration

Example `.pipe.yaml` for Cloud Run applications (v1):

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-cloudrun-app
  labels:
    env: production
    team: backend
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
  plugins:
    cloudrun:
      input:
        serviceManifestFile: service.yaml
```

### Cloud Run Service Manifest

Example `service.yaml`:

```yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: my-service
spec:
  template:
    spec:
      containers:
        - image: gcr.io/your-project/your-image:v1.0.0
          ports:
            - containerPort: 8080
          env:
            - name: ENV
              value: production
```

## Available Stages

- **CLOUDRUN_SYNC:** Deploy new revision with 100% traffic
- **CLOUDRUN_PROMOTE:** Gradually promote new revision
- **WAIT:** Wait for a specified duration
- **WAIT_APPROVAL:** Manual approval gate
- **ANALYSIS:** Automated deployment analysis

## Examples

### Quick Sync

Deploy immediately with 100% traffic:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-cloudrun-app
  labels:
    env: production
  pipeline:
    stages:
      - name: CLOUDRUN_SYNC
  plugins:
    cloudrun:
      input:
        serviceManifestFile: service.yaml

```

### Canary Rollout

Progressive rollout with automated analysis:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-cloudrun-app
  labels:
    env: production
    team: backend
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
          percent: 50
      - name: WAIT
        with:
          duration: 5m
      - name: CLOUDRUN_PROMOTE
        with:
          percent: 100
  plugins:
    cloudrun:
      input:
        serviceManifestFile: service.yaml
```

### Multi-Stage Canary with Approval

Gradual rollout with manual approval:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-cloudrun-app
  labels:
    env: production
    team: backend
  pipeline:
    stages:
      - name: CLOUDRUN_PROMOTE
        with:
          percent: 25
      - name: WAIT
        with:
          duration: 10m
      - name: WAIT_APPROVAL
      - name: CLOUDRUN_PROMOTE
        with:
          percent: 100
  plugins:
    cloudrun:
      input:
        serviceManifestFile: service.yaml
```

## Source Code

- [`pkg/app/pipedv1/plugin/cloudrun/`](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin/cloudrun)

## See Also

- [Configuration Reference](/docs-dev/user-guide/configuration-reference/)
- [Managing Applications](/docs-dev/user-guide/managing-application/)
- [Migrating to PipeCD V1](/docs-dev/migrating-from-v0-to-v1/)


## Source Code

- [`pkg/app/pipedv1/plugin/cloudrun/`](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin/cloudrun)

## See Also

- [Configuration Reference](/docs-dev/user-guide/configuration-reference/)
- [Managing Applications](/docs-dev/user-guide/managing-application/)
- [Migrating to PipeCD V1](/docs-dev/migrating-from-v0-to-v1/)
