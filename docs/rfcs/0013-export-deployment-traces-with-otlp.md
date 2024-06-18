- Start Date: 2024-06-13
- Target Version: 0.49.0

# Summary

This RFC proposes a new feature to export the OpenTelemetry Trace with deployment and stage spans.

# Motivation

We want to investigate deployment performances, such as the time to complete each deployment or stage.
We have `deployment_status` metrics for now, but this metric has one metric for each deployment status, and the count of metrics variety only glows larger.
When we save all histories of such metrics, we need more storage as time goes on since the completed deployments have a 0/1 value for each status.
Conversely, if we save the histories only while it's running, we have to fill the metrics with the last value or zero to get some statistics for deployments.

So, I propose to separate the metrics to get some statistics and the traces to get deployment performances.
This RFC proposes the latter, the traces, to get deployment performances.

# Detailed design

## Architecture

Collect spans at piped
→ send to pipecd-gateway envoy
→ proxies to OpenTelemetry Collector
→ send anywhere with exporters

## How to authorize piped at OpenTelemetry Collector

Envoy has a feature that authorizes incoming requests with an external authorizer.
[document](https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/filters/http/ext_authz/v3/ext_authz.proto)
We can use this filter by implementing [Authorization Service](https://github.com/envoyproxy/envoy/blob/d79f6e8d453ee260e9094093b8dd31af0056e67b/api/envoy/service/auth/v3/external_auth.proto#L29-L34) and some configuration to tell Envoy to use this service.
Then, the OpenTelemetry Collector only receives authorized requests.

At the piped, we configure the OpenTelemetry gRPC exporter with the authorization header, the same as piped requests to control plane requests.

## Traces planned to collect

I plan to collect traces/spans tagged with these values.
These are not secret values, but help investigate deployment performance and problems.

- project ID
- piped ID
- application ID
- application Kind
- deployment ID
- stage Name
- stage ID

## How to use these traces

We can send traces anywhere from OpenTelemetry Collector and use any hosting to collect/view them.
This section contains sample views of traces collected with Jaeger and its usage.

P.S.
The sample images in this section are from Jaeger UI. These traces are sent directly to Jaeger from Piped, so it's not implemented the way this RFC proposes. But it's enough to take sample images.

### How to configure to send traces you want
Users can set the Helm values to configure the OpenTelemetry Collector.
This sample configuration sends traces to the OTLP receiver running at `otlp.example.com:4317`.

```yaml
opentelemetry-collector:
    exporters:
      otlp:
        endpoint: otlp.example.com:4317

    service:
      pipelines:
        traces:
          exporters:
            - otlp
```

### Detail view of deployment trace
In Jaeger UI, we can see a detailed view of a deployment.
With this view, we can inspect which stage takes much time.
In this case, it's QuickSync, so there is one and only one stage, and it takes about 20 seconds.

![detail view of deployment trace](./assets/0013-jaeger-trace-detail.png)

### Timeline view of multiple deployment traces
In Jaeger UI, we can see multiple traces in one graph.
Each point in this graph represents the duration of the trace and the time it occurred.
With this graph, we can see the performances of multiple deployments.
In this case, many deployments occurred at the same time, and there are performance impacts. At the leftmost point of the graph, there is a deployment with a duration below 25 seconds. At the rightmost point, there are many deployments with durations over 30 seconds.

![timeline view of multiple deployment traces](./assets/0013-jaeger-trace-timeview.png)

# Alternatives

Another way is to implement a custom client to send traces to the control plane. Then, the control plane sends them to the OpenTelemetry Collectors.
It's harder to maintain because we have to maintain not only the custom client but also the control plane proxy implementation.
With this RPC's proposed method, we can only maintain envoy config and authentication mechanisms.
