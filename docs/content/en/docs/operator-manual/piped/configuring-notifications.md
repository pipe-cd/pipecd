---
title: "Configuring notifications"
linkTitle: "Configuring notifications"
weight: 7
description: >
  This page describes how to configure piped to send notications to external services.
---

> TBA

### Sending notifications to Slack

![](/images/notification-slack-deployment-planned.png)
<p style="text-align: center;">
Deployment was planned
</p>

![](/images/notification-slack-deployment-completed-successfully.png)
<p style="text-align: center;">
Deployment was completed successfully
</p>

![](/images/notification-slack-piped-started.png)
<p style="text-align: center;">
A piped has been started
</p>

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  notifications:
    routes:
      - name: dev-slack
        envs:
          - dev
        receiver: dev-slack-channel
      - name: prod-slack
        events:
          - DEPLOYMENT_STARTED
          - DEPLOYMENT_COMPLETED
        envs:
          - dev
        receiver: prod-slack-channel
      - name: all-events-to-ci
        receiver: ci-webhook
    receivers:
      - name: dev-slack-channel
        slack:
          hookURL: https://slack.com/dev
      - name: prod-slack-channel
        slack:
          hookURL: https://slack.com/prod
      - name: ci-webhook
        webhook:
          url: https://pipecd.dev/dev-hook
```

### Sending notifications to webhook endpoints
