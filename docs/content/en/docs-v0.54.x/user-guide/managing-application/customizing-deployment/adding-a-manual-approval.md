---
title: "Adding a manual approval stage"
linkTitle: "Manual approval stage"
weight: 2
description: >
  This page describes how to add a manual approval stage.
---

While deploying an application to production environments, some teams require manual approvals before continuing.
The manual approval stage enables you to control when the deployment is allowed to continue by requiring a specific person or team to approve.
This stage is named by `WAIT_APPROVAL` and you can add it to your pipeline before some stages should be approved before they can be executed.

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
      - name: WAIT_APPROVAL
        with:
          timeout: 6h
          approvers:
            - user-abc
      - name: K8S_PRIMARY_ROLLOUT
```

As above example, the deployment requires an approval from `user-abc` before `K8S_PRIMARY_ROLLOUT` stage can be executed.

The value of user ID in the `approvers` list depends on your [SSO configuration](../../../managing-controlplane/auth/), it must be GitHub's user ID if your SSO was configured to use GitHub provider, it must be Gmail account if your SSO was configured to use Google provider.

In case the `approvers` field was not configured, anyone in the project who has `Editor` or `Admin` role can approve the deployment pipeline.

Also, it will end with failure when the time specified in `timeout` has elapsed. Default is `6h`.

![](/images/deployment-wait-approval-stage.png)
<p style="text-align: center;">
Deployment with a WAIT_APPROVAL stage
</p>
