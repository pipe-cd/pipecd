---
title: "Automated deployment analysis"
linkTitle: "Automated deployment analysis"
weight: 8
description: >
  This page describes how to configure Automated Deployment Analysis feature.
---

>NOTE: This feature is currently alpha status.

Automated Deployment Analysis (ADA) lets you automate the verification of the deployment process by analyzing the metrics data, log entries, and the responses of the configured HTTP requests.
ADA is available as a [Stage](/docs/concepts/#stage) in the pipeline specified in the deployment configuration file.

ADA does the analysis by periodically performing queries against the [Analysis Provider](/docs/concepts/#analysis-provider) and evaluating the results to know the impact of the deployment. Then based on these evaluating results, the deployment can be rolled back immediately to minimize any negative impacts.

![](/images/deployment-analysis-stage.png)
<p style="text-align: center;">
Automatic rollback based on the analysis result
</p>

### Prerequisites
Before enabling ADA inside the pipeline, all required Analysis Providers must be configured in the Piped Configuration according to [this guide](/docs/operator-manual/piped/adding-an-analysis-provider/).

### Configuration
All you have to do is appending one or some `ANALYSIS` stages to your deployment pipeline:
```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  pipeline:
    stages:
      - name: ANALYSIS
        with:
          duration: 30m
          metrics:
            - provider: prometheus-dev
              interval: 5m
              query: grpc_request_error_percentage
              expected:
                max: 10
```

In the `provider` field, put the name of provider you have set to Piped configuration in the [Prerequisites](/docs/user-guide/automated-deployment-analysis/#prerequisites) section.

The `ANALYSIS` stage will continue to run for the period specified in the `duration` field.
In the meantime, Piped sends the given `query` to the Analysis Provider at each specified `interval`.

For each query, it checks if the result is within the expected range. If it's not expected, this `ANALYSIS` stage will fail (typically the rollback stage will be started).
You can change the acceptable number of failures by setting the `failureLimit` field.

The full list of configurable `ANALYSIS` stage fields are [here](/docs/user-guide/configuration-reference/#analysisstageoptions).

The canonical use case for this stage is to determine if your canary deployment should proceed. See more the [example](https://github.com/pipe-cd/examples/blob/master/kubernetes/analysis-by-metrics/.pipe.yaml).

### [Optional] Analysis Template
Analysis Templating is a feature that allows you to define some shared analysis configurations to be used by multiple applications. These templates must be placed at the `.pipe` directory at the root of the Git repository. Any application in that Git repository can use to the defined template by specifying the name of the template in the deployment configuration file.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: AnalysisTemplate
spec:
  metrics:
    http_error_rate:
      interval: 5m
      provider: prometheus-dev
      expected:
        max: 0
      query: |
        sum without(status) (rate(http_requests_total{status=~"5.*", job="{{ .App.Name }}"}[1m]))
        /
        sum without(status) (rate(http_requests_total{job="{{ .App.Name }}"}[1m]))
```


An `ANALYSIS` stage can reference a template with `template` field:
```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  pipeline:
    stages:
      - name: ANALYSIS
        with:
          duration: 30m
          metrics:
            - template:
                name: http_error_rate
```

Analysis Template uses the [Go templating engine](https://golang.org/pkg/text/template/) which only replaces values. This allows deployment-specific data to be embedded in the analysis template.

The available built-in args are:

| Property | Type | Description |
|-|-|-|
| App.Name | string | Application Name. |
| K8s.Namespace | string | The Kubernetes namespace where manifests will be applied. |

Also, custom args is supported. Custom args placeholders can be defined as `{{ .Args.<name> }}`.


See [here](https://github.com/pipe-cd/examples/blob/master/.pipe/analysis-template.yaml) for more examples.
And the full list of configurable `AnalysisTemplate` fields are [here](/docs/user-guide/configuration-reference/#analysis-template-configuration).
