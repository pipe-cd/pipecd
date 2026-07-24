---
title: "Wait plugin"
linkTitle: "Wait"
weight: 40
description: >
  Pause the pipeline for a fixed duration.
---

The `wait` plugin provides the `WAIT` stage, which pauses the pipeline for a fixed duration before continuing to the next stage. Because it is a stage plugin, `WAIT` can be added to any deployment pipeline.

## Prerequisites

Register the plugin in the piped configuration:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  # ...
  plugins:
    - name: wait
      port: 7004
      url: file:///path/to/plugin/binary  # or an https:// release URL
```

## The WAIT stage

Add a `WAIT` stage and set how long to pause under `with.duration`. For example, wait 5 minutes between a canary rollout and the primary rollout:

```yaml
pipeline:
  stages:
    - name: K8S_CANARY_ROLLOUT
      with:
        replicas: 10%
    - name: WAIT
      with:
        duration: 5m
    - name: K8S_PRIMARY_ROLLOUT
```

## Configuration reference

### WAIT stage options

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| duration | duration | How long to pause before continuing (e.g. `5m`). | Yes |
