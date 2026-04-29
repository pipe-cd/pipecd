---
title: "Cancelling a deployment"
linkTitle: "Cancelling a deployment"
weight: 5
description: >
  Learn how to cancel a running deployment using the Control Plane.
---

A running deployment can be cancelled from web UI at the deployment details page.

## Canceling a Deployment

When you cancel an active deployment, PipeCD can automatically execute a rollback to the previous stable version if rollback is enabled in your application configuration.

### Rollback Options

You have two ways to control rollback behavior when canceling a deployment:

1. **Automatic Rollback**: If enabled in the application configuration, rollback executes automatically after cancellation.

2. **Manual Selection**: Click the dropdown arrow (`▼`) next to the `CANCEL` button in the web UI to explicitly choose whether to rollback or not.

![A screenshot demonstrating how to cancel a deployment](/images/cancel-deployment.png)
<p style="text-align: center;">
Cancel a Deployment from web UI
</p>
