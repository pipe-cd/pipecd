---
title: "Metrics"
linkTitle: "Metrics"
weight: 7
description: >
  This page describes how to enable monitoring system for collecting PipeCD' metrics.
---

> WIP

PipeCD comes with a monitoring system including Prometheus, Alertmanager, and Grafana.
This page walks you through how to set up and use them.

## Enable monitoring system
To enable monitoring system for PipeCD, you first need to set the following value to `helm install` when [installing](/docs/operator-manual/control-plane/installation/#3-preparing-control-plane-configuration-file-and-installing).

```
--set monitoring.enabled=true
```


### Grafana dashboard
If you've already enabled monitoring system in the previous section, you can access Grafana using port forwarding:

```
kubectl port-forward -n {NAMESPACE} svc/{PIPECD_RELEASE_NAME}-grafana 3000:80
```

### Alert notifications
If you want to send alert notifications to external services like Slack, you need to set an alertmanager configuration file.

For example, let's say you use Slack as a receiver. Create `values.yaml` and put the following configuration to there.

```yaml
prometheus:
  alertmanagerFiles:
    alertmanager.yml:
      global:
        slack_api_url: {YOUR_WEBHOOK_URL}
      route:
        receiver: slack-notifications
      receivers:
        - name: slack-notifications
          slack_configs:
            - channel: '#your-channel'
```

And give it to the `helm install` command when [installing](/docs/operator-manual/control-plane/installation/#3-preparing-control-plane-configuration-file-and-installing).

```
--values=values.yaml
```
