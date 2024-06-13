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

## How to authenticate piped at OpenTelemetry Collector

OpenTelemetry Collector has a customization feature that implements a custom authenticator.
[go.opentelemetry.io/collector/extension/auth on pkg.go.dev](https://pkg.go.dev/go.opentelemetry.io/collector/extension/auth@v0.102.1)
We can implement authentication by implementing this Client/Server interface, then using the Client at piped and the Server with a collector.

# Alternatives

Another way is to implement a custom client to send traces to the control plane. Then, the control plane sends them to the OpenTelemetry Collectors.
It's harder to maintain because we have to maintain not only the custom client but also the control plane proxy implementation.
With this RPC's proposed method, we can only maintain envoy config and authentication mechanisms.

# Unresolved questions

There is no plan for detailed implementations of custom authentication extensions for the OpenTelemetry Collector.
