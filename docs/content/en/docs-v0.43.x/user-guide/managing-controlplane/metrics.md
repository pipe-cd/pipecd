---
title: "Metrics"
linkTitle: "Metrics"
weight: 5
description: >
  This page describes how to enable monitoring system for collecting PipeCD' metrics.
---

PipeCD comes with a monitoring system including Prometheus, Alertmanager, and Grafana.
This page walks you through how to set up and use them.

## Monitoring overview

![](/images/metrics-architecture.png)
<p style="text-align: center;">
Monitoring Architecture
</p>

The piped agent collects its own monitoring information and periodically sends metrics to the Control plane, which publishes all of the resource usage and cluster information it collects and the metrics sent by the piped agent on its http server. When the PipeCD monitoring feature is turned on, Prometheus, Alertmanager and Grafana are deployed with the Control plane, and Prometheus retrieves metrics information from the Control plane's http server.

The piped agent not only sends the collected monitoring metrics to the Control plane, but also provides them through its internal http server, so that the developer managing the piped agent can also get metrics directly from the piped agent and monitor them with their custom monitoring service.

## Enable monitoring system
To enable monitoring system for PipeCD, you first need to set the following value to `helm install` when [installing](../../../installation/install-controlplane/#2-preparing-control-plane-configuration-file-and-installing).

```
--set monitoring.enabled=true
```

## Dashboards
If you've already enabled monitoring system in the previous section, you can access Grafana using port forwarding:

```
kubectl port-forward -n {NAMESPACE} svc/{PIPECD_RELEASE_NAME}-grafana 3000:80
```

#### Control Plane dashboards
There are three dashboards related to Control Plane:
- Overview - usage stats of PipeCD
- Incoming Requests - gRPC and HTTP requests stats to check for any negative impact on users
- Go - processes stats of PipeCD components

#### Piped dashboards
Visualize the metrics of Piped registered in the Control plane.
- Overview - usage stats of piped agents
- Process - resource usage of piped agent
- Go - processes stats of piped agents.

#### Cluster dashboards
Because cluster dashboards tracks cluster-wide metrics, defaults to disable. You can enable it with:

```
--monitoring.clusterStats=true
```

There are three dashboards that track metrics for:
- Node - nodes stats within the Kubernetes cluster where PipeCD runs on
- Pod - stats for pods that make PipeCD up
- Prometheus - stats for Prometheus itself

## Alert notifications
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

And give it to the `helm install` command when [installing](../../../installation/install-controlplane/#2-preparing-control-plane-configuration-file-and-installing).

```
--values=values.yaml
```

See [here](https://prometheus.io/docs/alerting/latest/configuration/) for more details on AlertManager's configuration.
