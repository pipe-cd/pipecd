---
title: "CloudRun"
linkTitle: "CloudRun"
weight: 3
description: >
  Specific guide for configuring CloudRun deployment.
---

Deploying a CloudRun application requires a `service.yaml` file placing inside the application directory. That file contains the service specification used by CloudRun as following: 

``` yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: SERVICE_NAME
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/maxScale: '5'
    spec:
      containerConcurrency: 80
      containers:
      - args:
        - server
        image: gcr.io/pipecd/helloworld:v0.5
        ports:
        - containerPort: 9085
        resources:
          limits:
            cpu: 1000m
            memory: 128Mi
```

## Quick sync

By default, when the [pipeline](/docs/user-guide/configuration-reference/#cloudrun-application) was not specified, PipeCD triggers a quick sync deployment for the merged pull request.
Quick sync for a Cloud Run deployment will roll out the new version and switch all traffic to it.

## Sync with the specified pipeline

The [pipeline](/docs/user-guide/configuration-reference/#cloudrun-application) field in the deployment configuration is used to customize the way to do the deployment.
You can add a manual approval before routing traffic to the new version or add an analysis stage the do some smoke tests against the new version before allowing them to receive the real traffic.

These are the provided stages for CloudRun application you can use to build your pipeline:

- `CLOUDRUN_PROMOTE`
  - promote the new version to receive an amount of traffic

and other common stages:
- `WAIT`
- `WAIT_APPROVAL`
- `ANALYSIS`

See the description of each stage at [Configuration Reference](/docs/user-guide/configuration-reference/#stageoptions).

Here is an example that rolls out the new version gradually:

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: CloudRunApp
spec:
  pipeline:
    stages:
      # Promote new version to receive 10% of traffic.
      - name: CLOUDRUN_PROMOTE
        with:
          percent: 10
      - name: WAIT
        with:
          duration: 10m
      # Promote new version to receive 50% of traffic.
      - name: CLOUDRUN_PROMOTE
        with:
          percent: 50
      - name: WAIT
        with:
          duration: 10m
      # Promote new version to receive all traffic.
      - name: CLOUDRUN_PROMOTE
        with:
          percent: 100
```

## Reference

See [Configuration Reference](/docs/user-guide/configuration-reference/#cloudrun-application) for the full configuration.
