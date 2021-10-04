---
title: "Configuring notifications"
linkTitle: "Configuring notifications"
weight: 7
description: >
  This page describes how to configure piped to send notifications to external services.
---

PipeCD events (deployment triggered, planned, completed, analysis result, piped started...) can be sent to external services like Slack or a Webhook service. While forwarding those events to a chat service helps developers have a quick and convenient way to know the deployment's current status, forwarding to a Webhook service may be useful for triggering other related tasks like CI jobs.

PipeCD events are emitted and sent by the `piped` component. So all the needed configurations can be specified in the `piped` configuration file.
Notification configuration including:
- a list of `Route`s which used to match events and decide where the event should be sent to
- a list of `Receiver`s which used to know how to send events to the external service

[Notification Route](/docs/operator-manual/piped/configuration-reference/#notificationroute) matches events based on their metadata like `name`, `group`, `env`, `app`.
Below is the list of supporting event names and their groups.(Events of `APPLICATION~` are not supported yet.)

| Event | Group |
|-|-|
| DEPLOYMENT_TRIGGERED | DEPLOYMENT |
| DEPLOYMENT_PLANNED | DEPLOYMENT |
| DEPLOYMENT_APPROVED | DEPLOYMENT |
| DEPLOYMENT_WAIT_APPROVAL | DEPLOYMENT |
| DEPLOYMENT_ROLLING_BACK | DEPLOYMENT |
| DEPLOYMENT_SUCCEEDED | DEPLOYMENT |
| DEPLOYMENT_FAILED | DEPLOYMENT |
| DEPLOYMENT_CANCELLED | DEPLOYMENT |
| APPLICATION_SYNCED | APPLICATION_SYNC |
| APPLICATION_OUT_OF_SYNC | APPLICATION_SYNC |
| APPLICATION_HEALTHY | APPLICATION_HEALTH |
| APPLICATION_UNHEALTHY | APPLICATION_HEALTH |
| PIPED_STARTED | PIPED |
| PIPED_STOPPED | PIPED |

### Sending notifications to Slack

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  notifications:
    routes:
      # Sending all event from development environment to dev-slack-channel.
      - name: dev-slack
        envs:
          - dev
        receiver: dev-slack-channel
      # Only sending deployment started and completed events to prod-slack-channel.
      - name: prod-slack
        events:
          - DEPLOYMENT_STARTED
          - DEPLOYMENT_COMPLETED
        envs:
          - dev
        receiver: prod-slack-channel
      # Sending all events a CI service.
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


![](/images/slack-notification-deployment.png)
<p style="text-align: center;">
Deployment was triggered, planned and completed successfully
</p>

![](/images/slack-notification-piped-started.png)
<p style="text-align: center;">
A piped has been started
</p>


For detailed configuration, please check the [configuration reference](/docs/operator-manual/piped/configuration-reference/#notifications) section.

### Sending notifications to webhook endpoints

> TBA
