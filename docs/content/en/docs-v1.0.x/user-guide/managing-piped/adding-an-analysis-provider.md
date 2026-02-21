---
title: "Adding an analysis provider"
linkTitle: "Adding analysis provider"
weight: 6
description: >
  This page describes how to add an Analysis Provider to analyize the metrics of your deployment.
---

To enable [Automated deployment analysis](../../managing-application/customizing-deployment/automated-deployment-analysis/) feature, you have to set the needed information for Piped to connect to the [Analysis Provider](../../../concepts/#analysis-provider).

Currently, PipeCD supports the following providers:

- [Prometheus](https://prometheus.io/)
- [Datadog](https://datadoghq.com/)

## Prometheus

Piped queries the [range query endpoint](https://prometheus.io/docs/prometheus/latest/querying/api/#range-queries) to obtain metrics used to evaluate the deployment.

You need to define the Prometheus server address so that it can be accessed by your `piped`.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  plugins:
    - name: analysis
      port:
      url:
      config:
        analysisProviders:
          - name: prometheus-dev
            type: PROMETHEUS
            config:
              address: https://your-prometheus.dev
```

To know more, see the full list of [configurable fields](configuration-reference/#analysisproviderdatadogconfig).

## Datadog

Piped queries the [MetricsApi.QueryMetrics](https://docs.datadoghq.com/api/latest/metrics/#query-timeseries-points) endpoint to obtain metrics used to evaluate the deployment.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  plugins:
    - name: analysis
      port:
      url:
      config:
        analysisProviders:
          - name: datadog-dev
            type: DATADOG
            config: 
              apiKeyFile: /etc/piped-secret/datadog-api-key
              applicationKeyFile: /etc/piped-secret/datadog-application-key
```

To know more, see the full list of [configurable fields](configuration-reference/#analysisproviderdatadogconfig).

If you choose `Helm` as the installation method, we recommend using `--set-file` to mount the key files while performing the [upgrading process](../../../installation/install-piped/installing-on-kubernetes/#in-the-cluster-wide-mode).

```bash
--set-file secret.data.datadog-api-key={PATH_TO_API_KEY_FILE} \
--set-file secret.data.datadog-application-key={PATH_TO_APPLICATION_KEY_FILE}
```
