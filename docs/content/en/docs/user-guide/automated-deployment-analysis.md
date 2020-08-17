---
title: "Automated Deployment Analysis"
linkTitle: "Automated Deployment Analysis"
weight: 11
description: >
  This page describes how to configure Automated Deployment Analysis feature.
---

Automated Deployment Analysis (ADA) lets you automate the rollbacks by analyzing the metrics and logs.
ADA is available as a [Stage](/docs/concepts/#stage) in your Pipeline.

The ADA analyzes the new version by periodically performing queries against the [Analysis Provider](/docs/concepts/#analysis-provider) and evaluating the results.

### Getting started
All you have to do is modify Deployment and Piped configuration.

First up, you define the information needed to connect from your Piped to the Analysis Provider:
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
The full list of configurable Analysis Providers fields are [here](/docs/operator-manual/piped/configuration-reference/#analysisprovider).

And then append the `ANALYSIS` stage to your deployment pipeline. The canonical use case for that stage is to determine if your canary deployment should proceed:
```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 10%
      - name: ANALYSIS
        with:
          duration: 10m
          metrics:
            - provider: prometheus-dev
              interval: 1m
              failureLimit: 1
              expected:
                max: 10
              query: grpc_request_error_percentage
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_CANARY_CLEAN
```
The full list of configurable `ANALYSIS` stage fields are [here](/docs/user-guide/configuration-reference/#analysisstageoptions).

### Analysis Template
With Analysis Template, you can re-use the analysis configuration.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: AnalysisTemplate
spec:
  metrics:
    grpc_error_rate_percentage:
      interval: 1m
      provider: prometheus-dev
      failureLimit: 1
      expected:
        max: 10
      query: |
        100 - sum(
            rate(
                grpc_server_handled_total{
                  grpc_code!="OK",
                  kubernetes_namespace="{{ .K8s.Namespace }}",
                  kubernetes_pod_name=~"{{ .App.Name }}-[0-9a-zA-Z]+(-[0-9a-zA-Z]+)"
                }[{{ .Args.interval }}]
            )
        )
        /
        sum(
            rate(
                grpc_server_started_total{
                  kubernetes_namespace="{{ .K8s.Namespace }}",
                  kubernetes_pod_name=~"{{ .App.Name }}-[0-9a-zA-Z]+(-[0-9a-zA-Z]+)"
                }[{{ .Args.interval }}]
            )
        ) * 100
```

Place it under the `.pipe/` directory, as in [the example](https://github.com/pipe-cd/examples/blob/master/.pipe/analysis-template.yaml). The full list of configurable `AnalysisTemplate` fields are [here](/docs/user-guide/configuration-reference/#analysis-template-configuration).

An `ANALYSIS` stage can reference a template with `template` field:
```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  pipeline:
    stages:
      - name: ANALYSIS
        with:
          duration: 10m
          metrics:
            - template:
                name: grpc_error_rate_percentage
                args:
                  interval: 1m
```

Analysis Template uses the [Go templating engine](https://golang.org/pkg/text/template/) which only replaces values. This allows deployment-specific data to be embedded in the analysis template.

The available built-in args are:

| Property | Type | Description |
|-|-|-|
| App.Name | string | Application Name. |
| K8s.Namespace | string | The Kubernetes namespace where manifests will be applied. |

Also, custom args is supported. Custom args placeholders can be defined as `{{ .Args.<name> }}`

### Supported Providers

##### Prometheus
```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  analysisProviders:
    - name: prometheus-name
      type: PROMETHEUS
      config:
        address: https://your-prometheus.dev
```

##### Datadog (comming soon)

