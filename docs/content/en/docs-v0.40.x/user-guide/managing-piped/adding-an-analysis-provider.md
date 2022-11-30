---
title: "Adding an analysis provider"
linkTitle: "Adding analysis provider"
weight: 6
description: >
  This page describes how to add an analysis provider for doing deployment analysis.
---

To enable [Automated deployment analysis](../../managing-application/customizing-deployment/automated-deployment-analysis/) feature, you have to set the needed information for Piped to connect to the [Analysis Provider](../../../concepts/#analysis-provider).

Currently, PipeCD supports the following providers:
- [Prometheus](https://prometheus.io/)
- [Datadog](https://datadoghq.com/)


## Prometheus
Piped queries the [range query endpoint](https://prometheus.io/docs/prometheus/latest/querying/api/#range-queries) to obtain metrics used to evaluate the deployment.

You need to define the Prometheus server address accessible to Piped.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  analysisProviders:
    - name: prometheus-dev
      type: PROMETHEUS
      config:
        address: https://your-prometheus.dev
```
The full list of configurable fields are [here](../configuration-reference/#analysisproviderprometheusconfig).

## Datadog
Piped queries the [MetricsApi.QueryMetrics](https://docs.datadoghq.com/api/latest/metrics/#query-timeseries-points) endpoint to obtain metrics used to evaluate the deployment.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  analysisProviders:
    - name: datadog-dev
      type: DATADOG
      config:
        apiKeyFile: /etc/piped-secret/datadog-api-key
        applicationKeyFile: /etc/piped-secret/datadog-application-key
```

The full list of configurable fields are [here](../configuration-reference/#analysisproviderdatadogconfig).

If you choose `Helm` as the installation method, we recommend using `--set-file` to mount the key files while performing the [upgrading process](../../../installation/install-piped/installing-on-kubernetes/#in-the-cluster-wide-mode).

```console
--set-file secret.data.datadog-api-key={PATH_TO_API_KEY_FILE} \
--set-file secret.data.datadog-application-key={PATH_TO_APPLICATION_KEY_FILE}
```
