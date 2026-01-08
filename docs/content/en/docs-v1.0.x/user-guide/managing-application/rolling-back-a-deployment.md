---
title: "Rolling back a deployment"
linkTitle: "Rolling back a deployment"
weight: 6
description: >
  This page describes when a deployment is rolled back automatically and how to manually roll back a deployment.
---

Rollbacks allow you to revert your application to a previous stable state when something goes wrong during deployment. PipeCD supports both automatic and manual rollbacks, giving you flexibility in handling deployment failures.

## Automatic rollback
You can automate rollbacks by enabling the `autoRollback` field in your application configuration:

```bash
spec:
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
      - name: K8S_PRIMARY_ROLLOUT
  autoRollback: true
```

When `autoRollback` is enabled, the deployment will be rolled back if any of the following conditions are met:

- A deployment pipeline stage fails
- An analysis stage determines that the deployment has a negative impact
- An error occurs while deploying

If a rollback is triggered, PipeCD adds a new `ROLLBACK` stage to the deployment pipeline and reverts all applied changes.

## Manual rollback

If you need to roll back a deployment that has already completed, or want to intervene during a deployment which is in-progress, you can trigger a rollback manually from the Web UI by clicking the `Cancel with Rollback` button.

![Screenshot of rolling back a deployment](/images/rolled-back-deployment.png)
<p style="text-align: center;">
A deployment bring rolled back
</p>
