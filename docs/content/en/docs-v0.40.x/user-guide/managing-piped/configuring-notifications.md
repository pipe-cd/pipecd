---
title: "Configuring notifications"
linkTitle: "Configuring notifications"
weight: 8
description: >
  This page describes how to configure piped to send notifications to external services.
---

PipeCD events (deployment triggered, planned, completed, analysis result, piped started...) can be sent to external services like Slack or a Webhook service. While forwarding those events to a chat service helps developers have a quick and convenient way to know the deployment's current status, forwarding to a Webhook service may be useful for triggering other related tasks like CI jobs.

PipeCD events are emitted and sent by the `piped` component. So all the needed configurations can be specified in the `piped` configuration file.
Notification configuration including:
- a list of `Route`s which used to match events and decide where the event should be sent to
- a list of `Receiver`s which used to know how to send events to the external service

[Notification Route](configuration-reference/#notificationroute) matches events based on their metadata like `name`, `group`, `app`, `labels`.
Below is the list of supporting event names and their groups.

| Event | Group | Supported | Description |
|-|-|-|-|
| DEPLOYMENT_TRIGGERED | DEPLOYMENT | <p style="text-align: center;"><input type="checkbox" checked disabled></p> |  |
| DEPLOYMENT_PLANNED | DEPLOYMENT | <p style="text-align: center;"><input type="checkbox" checked disabled></p> |  |
| DEPLOYMENT_APPROVED | DEPLOYMENT | <p style="text-align: center;"><input type="checkbox" checked disabled></p> |  |
| DEPLOYMENT_WAIT_APPROVAL | DEPLOYMENT | <p style="text-align: center;"><input type="checkbox" checked disabled></p> |  |
| DEPLOYMENT_ROLLING_BACK | DEPLOYMENT | <p style="text-align: center;"><input type="checkbox" disabled></p> | PipeCD sends a notification when a deployment is completed, while it does not send a notification when a deployment status changes to DEPLOYMENT_ROLLING_BACK because it is not a completion status. See [#4547](https://github.com/pipe-cd/pipecd/issues/4547) |
| DEPLOYMENT_SUCCEEDED | DEPLOYMENT | <p style="text-align: center;"><input type="checkbox" checked disabled></p> |  |
| DEPLOYMENT_FAILED | DEPLOYMENT | <p style="text-align: center;"><input type="checkbox" checked disabled></p> |  |
| DEPLOYMENT_CANCELLED | DEPLOYMENT | <p style="text-align: center;"><input type="checkbox" checked disabled></p> |  |
| DEPLOYMENT_TRIGGER_FAILED | DEPLOYMENT | <p style="text-align: center;"><input type="checkbox" checked disabled></p> |  |
| APPLICATION_SYNCED | APPLICATION_SYNC | <p style="text-align: center;"><input type="checkbox" disabled></p> |  |
| APPLICATION_OUT_OF_SYNC | APPLICATION_SYNC | <p style="text-align: center;"><input type="checkbox" disabled></p> |  |
| APPLICATION_HEALTHY | APPLICATION_HEALTH | <p style="text-align: center;"><input type="checkbox" disabled></p> |  |
| APPLICATION_UNHEALTHY | APPLICATION_HEALTH | <p style="text-align: center;"><input type="checkbox" disabled></p> |  |
| PIPED_STARTED | PIPED | <p style="text-align: center;"><input type="checkbox" checked  disabled></p> |  |
| PIPED_STOPPED | PIPED | <p style="text-align: center;"><input type="checkbox" checked disabled></p> |  |

### Sending notifications to Slack

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  notifications:
    routes:
      # Sending all event which contains labels `env: dev` to dev-slack-channel.
      - name: dev-slack
        labels:
          env: dev
        receiver: dev-slack-channel
      # Only sending deployment started and completed events which contains
      # labels `env: prod` and `team: pipecd` to prod-slack-channel.
      - name: prod-slack
        events:
          - DEPLOYMENT_TRIGGERED
          - DEPLOYMENT_SUCCEEDED
        labels:
          env: prod
          team: pipecd
        receiver: prod-slack-channel
    receivers:
      - name: dev-slack-channel
        slack:
          hookURL: https://slack.com/dev
      - name: prod-slack-channel
        slack:
          hookURL: https://slack.com/prod
```


![](/images/slack-notification-deployment.png)
<p style="text-align: center;">
Deployment was triggered, planned and completed successfully
</p>

![](/images/slack-notification-piped-started.png)
<p style="text-align: center;">
A piped has been started
</p>


For detailed configuration, please check the [configuration reference for Notifications](configuration-reference/#notifications) section.

### Sending notifications to external services via webhook

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  notifications:
    routes:
      # Sending all events an external service.
      - name: all-events-to-a-external-service
        receiver: a-webhook-service
    receivers:
      - name: a-webhook-service
        webhook:
          url: {WEBHOOK_SERVICE_URL}
          signatureValue: {RANDOM_SIGNATURE_STRING}
```

For detailed configuration, please check the [configuration reference for NotificationReceiverWebhook](configuration-reference/#notificationreceiverwebhook) section.
