---
title: "Analysis plugin"
linkTitle: "Analysis"
weight: 30
description: >
  Evaluate a deployment by analyzing metrics, logs, and HTTP responses.
---

The `analysis` plugin provides the `ANALYSIS` stage, which evaluates a running deployment for a defined period by querying metrics, logs, or HTTP endpoints. If the results fall outside the configured expectations, the stage fails and the deployment is rolled back. Because it is a stage plugin, `ANALYSIS` can be added to any deployment pipeline (for example between a canary rollout and a primary rollout).

Analysis uses providers that are configured once in the **piped** configuration. Each `ANALYSIS` stage then references a provider by name.

Provider support by analysis type:

| Analysis type | Providers |
|---------------|-----------|
| Metrics | Prometheus, Datadog |
| Logs | Stackdriver |
| HTTP | None (the stage queries the URL directly) |

## Prerequisites

1. **Register the plugin and its providers in the piped configuration.** Add an `analysis` plugin block and list your providers under `config.analysisProviders`:

   ```yaml
   apiVersion: pipecd.dev/v1beta1
   kind: Piped
   spec:
     # ...
     plugins:
       - name: analysis
         port: 7003
         url: file:///path/to/plugin/binary  # or an https:// release URL
         config:
           analysisProviders:
             - name: prometheus-dev
               type: PROMETHEUS
               config:
                 address: https://your-prometheus.dev
   ```

2. **Add an `ANALYSIS` stage** to a deployment pipeline and reference the provider by name (see [The ANALYSIS stage](#the-analysis-stage)).

## Analysis providers

Providers are defined under the plugin's `config.analysisProviders` in the **piped** configuration. Each entry has a `name`, a `type`, and a type-specific `config`.

### Prometheus

Used for metrics analysis. `piped` queries the Prometheus [range query endpoint](https://prometheus.io/docs/prometheus/latest/querying/api/#range-queries).

```yaml
analysisProviders:
  - name: prometheus-dev
    type: PROMETHEUS
    config:
      address: https://your-prometheus.dev
```

See [AnalysisProviderPrometheusConfig](#analysisproviderprometheusconfig) for all fields.

### Datadog

Used for metrics analysis. `piped` queries the Datadog [query timeseries endpoint](https://docs.datadoghq.com/api/latest/metrics/#query-timeseries-points).

```yaml
analysisProviders:
  - name: datadog-dev
    type: DATADOG
    config:
      apiKeyFile: /etc/piped-secret/datadog-api-key
      applicationKeyFile: /etc/piped-secret/datadog-application-key
```

See [AnalysisProviderDatadogConfig](#analysisproviderdatadogconfig) for all fields.

If you install `piped` with Helm, use `--set-file` to mount the key files during the [upgrade process](../installation/install-piped/installing-on-kubernetes/#in-the-cluster-wide-mode):

```bash
--set-file secret.data.datadog-api-key={PATH_TO_API_KEY_FILE} \
--set-file secret.data.datadog-application-key={PATH_TO_APPLICATION_KEY_FILE}
```

### Stackdriver

Used for log analysis.

```yaml
analysisProviders:
  - name: stackdriver-dev
    type: STACKDRIVER
    config:
      serviceAccountFile: /etc/piped-secret/gcp-service-account.json
```

See [AnalysisProviderStackdriverConfig](#analysisproviderstackdriverconfig) for all fields.

## The ANALYSIS stage

Add an `ANALYSIS` stage to a pipeline and configure it under `with`. The `duration` field is required and sets how long the analysis runs. Within that window, the stage can run three kinds of checks, each a list:

- `metrics` - query a metrics provider (Prometheus or Datadog).
- `logs` - query a log provider (Stackdriver).
- `https` - send HTTP requests and check the response.

Example: run a canary, then evaluate its error rate against a threshold before promoting to primary.

```yaml
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
            strategy: THRESHOLD
            query: |
              sum(rate(http_requests_total{job="my-app",status=~"5.."}[1m]))
              / sum(rate(http_requests_total{job="my-app"}[1m]))
            interval: 1m
            expected:
              max: 0.01
    - name: K8S_PRIMARY_ROLLOUT
```

### Metrics strategies

Metrics analysis supports four strategies, set per metrics check with the `strategy` field:

- `THRESHOLD` (default) - compare each query result against the `expected` range (`min`/`max`).
- `PREVIOUS` - compare the result against the same query from the previous successful deployment.
- `CANARY_BASELINE` - compare the canary variant against the baseline variant.
- `CANARY_PRIMARY` - compare the canary variant against the primary variant. (Not recommended, since primary may serve production traffic.)

For `PREVIOUS`, `CANARY_BASELINE`, and `CANARY_PRIMARY`, use `deviation` to set which direction counts as a failure (`LOW`, `HIGH`, or `EITHER`).

## Configuration reference

### AnalysisStageOptions

Options for the `ANALYSIS` stage, set under `with`.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| duration | duration | How long the analysis runs (e.g. `10m`). | Yes |
| restartThreshold | int | Allowed number of pod restarts during the analysis. | No |
| metrics | [][AnalysisMetrics](#analysismetrics) | Metrics checks to run. | No |
| logs | [][AnalysisLog](#analysislog) | Log checks to run. | No |
| https | [][AnalysisHTTP](#analysishttp) | HTTP checks to run. | No |

### AnalysisMetrics

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| strategy | string | One of `THRESHOLD`, `PREVIOUS`, `CANARY_BASELINE`, `CANARY_PRIMARY`. | No (default `THRESHOLD`) |
| provider | string | Name of a metrics provider defined in the piped configuration. | Yes |
| query | string | Query run against the provider. | Yes |
| expected | [AnalysisExpected](#analysisexpected) | Expected result range. Required for the `THRESHOLD` strategy. | For `THRESHOLD` |
| interval | duration | How often the query runs. | Yes |
| failureLimit | int | Number of failed checks tolerated before the analysis fails. | No (default `0`) |
| skipOnNoData | bool | Treat "no data returned" as a success. | No (default `false`) |
| timeout | duration | Query timeout. | No (default `30s`) |
| deviation | string | Failure direction for non-threshold strategies: `LOW`, `HIGH`, or `EITHER`. | No (default `EITHER`) |
| canaryArgs | map[string]string | Template args for the canary query, referenced as `{{ .VariantArgs.xxx }}`. | No |
| baselineArgs | map[string]string | Template args for the baseline query. | No |
| primaryArgs | map[string]string | Template args for the primary query. | No |

### AnalysisExpected

At least one of `min` or `max` is required.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| min | float | Minimum acceptable value. | No |
| max | float | Maximum acceptable value. | No |

### AnalysisLog

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| provider | string | Name of a log provider defined in the piped configuration. | Yes |
| query | string | Query run against the provider. | Yes |
| interval | duration | How often the query runs. | Yes |
| failureLimit | int | Number of failed checks tolerated before the analysis fails. | No (default `0`) |
| skipOnNoData | bool | Treat "no data returned" as a success. | No (default `false`) |
| timeout | duration | Query timeout. | No |

### AnalysisHTTP

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| url | string | URL to send the request to. | Yes |
| method | string | HTTP method. | No |
| headers | [][AnalysisHTTPHeader](#analysishttpheader) | Request headers. | No |
| expectedCode | int | Expected HTTP status code. | No |
| expectedResponse | string | Expected response body. | No |
| interval | duration | How often the request is sent. | Yes |
| failureLimit | int | Number of failed checks tolerated before the analysis fails. | No (default `0`) |
| skipOnNoData | bool | Treat "no data returned" as a success. | No (default `false`) |
| timeout | duration | Request timeout. | No |

### AnalysisHTTPHeader

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| key | string | Header name. | Yes |
| value | string | Header value. | Yes |

### AnalysisProviderPrometheusConfig

Configured under `config.analysisProviders[].config` in the **piped** configuration.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| address | string | Address of the Prometheus server. | Yes |
| usernameFile | string | Path to a file containing the username for basic auth. | No |
| passwordFile | string | Path to a file containing the password for basic auth. | No |

### AnalysisProviderDatadogConfig

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| address | string | Datadog API server address. One of `datadoghq.com`, `us3.datadoghq.com`, `datadoghq.eu`, `ddog-gov.com`. | No (default `datadoghq.com`) |
| apiKeyFile | string | Path to the API key file. Mutually exclusive with `apiKeyData`. | Yes (or `apiKeyData`) |
| applicationKeyFile | string | Path to the application key file. Mutually exclusive with `applicationKeyData`. | Yes (or `applicationKeyData`) |
| apiKeyData | string | Base64-encoded API key. Mutually exclusive with `apiKeyFile`. | No |
| applicationKeyData | string | Base64-encoded application key. Mutually exclusive with `applicationKeyFile`. | No |

### AnalysisProviderStackdriverConfig

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| serviceAccountFile | string | Path to the GCP service account file. | Yes |

### AnalysisApplicationSpec

The `spec` of an application using analysis shares the [common application fields](../user-guide/managing-application/configuration-reference/) and adds the following under `plugins.analysis`:

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| appCustomArgs | map[string]string | Custom arguments populated into queries, referenced as `{{ .AppCustomArgs.xxx }}`. | No |
