---
title: "Configuring Cloud Run application"
linkTitle: "Cloud Run"
weight: 3
description: >
  Specific guide to configuring deployment for Cloud Run application.
---

Deploying a Cloud Run application requires a `service.yaml` file placing inside the application directory. That file contains the service specification used by Cloud Run as following: 

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

By default, when the [pipeline](../../../configuration-reference/#cloud-run-application) was not specified, PipeCD triggers a quick sync deployment for the merged pull request.
Quick sync for a Cloud Run deployment will roll out the new version and switch all traffic to it.

## Sync with the specified pipeline

The [pipeline](../../../configuration-reference/#cloud-run-application) field in the application configuration is used to customize the way to do the deployment.
You can add a manual approval before routing traffic to the new version or add an analysis stage the do some smoke tests against the new version before allowing them to receive the real traffic.

These are the provided stages for Cloud Run application you can use to build your pipeline:

- `CLOUDRUN_PROMOTE`
  - promote the new version to receive an amount of traffic

and other common stages:
- `WAIT`
- `WAIT_APPROVAL`
- `ANALYSIS`

See the description of each stage at [Customize application deployment](../../customizing-deployment/).

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

See [Configuration Reference](../../../configuration-reference/#cloud-run-application) for the full configuration.
