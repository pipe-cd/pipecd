---
title: "Rolling back a deployment"
linkTitle: "Rolling back a deployment"
weight: 6
description: >
  This page describes when a deployment is rolled back automatically and how to manually rollback a deployment.
---

Rolling back a deployment can be automated by enabling the `autoRollback` field in the application configuration of the application. When `autoRollback` is enabled, the deployment will be rolled back if any of the following conditions are met:
- a stage of the deployment pipeline was failed
- an analysis stage determined that the deployment had a negative impact
- any error occurs while deploying

When the rolling back process is triggered, a new `ROLLBACK` stage will be added to the deployment pipeline and it reverts all the applied changes.

![](/images/rolled-back-deployment.png)
<p style="text-align: center;">
A deployment was rolled back
</p>

Alternatively, manually rolling back a running deployment can be done from web UI by clicking on `Cancel with rollback` button.
