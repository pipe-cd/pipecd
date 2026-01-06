---
title: "Rolling back a deployment"
linkTitle: "Rolling back a deployment"
weight: 6
description: >
  This page describes when a deployment is rolled back automatically and how to manually roll back a deployment.
---

You can automate rollbacks by enabling the `autoRollback` field in your application configuration. When `autoRollback` is enabled, the deployment will be rolled back if any of the following conditions are met:

- A deployment pipeline stage fails
- An analysis stage determines that the deployment has a negative impact
- An error occurs while deploying

When a rollback is triggered, PipeCD adds a new `ROLLBACK` stage to the deployment pipeline and reverts all applied changes.

![Screenshot of rolling back a deployment](/images/rolled-back-deployment.png)
<p style="text-align: center;">
A deployment was rolled back
</p>

Alternatively, you can manually roll back a deployment from the Web UI by clicking the `Cancel with Rollback` button.
