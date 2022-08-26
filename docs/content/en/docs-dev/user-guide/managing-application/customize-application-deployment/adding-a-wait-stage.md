---
title: "Adding a wait stage"
linkTitle: "Wait stage"
weight: 1
description: >
  This page describes how to add a WAIT stage.
---

In addition to waiting for approvals from someones, the deployment pipeline can be configured to wait an amount of time before continuing.
This can be done by adding the `WAIT` stage into the pipeline. This stage has one configurable field is `duration` to configure how long should be waited.

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
      - name: WAIT
        with:
          duration: 5m
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_CANARY_CLEAN
```

![](/images/deployment-wait-stage.png)
<p style="text-align: center;">
Deployment with a WAIT stage
</p>
