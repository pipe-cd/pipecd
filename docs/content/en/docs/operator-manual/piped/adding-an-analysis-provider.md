---
title: "Adding an analysis provider"
linkTitle: "Adding an analysis provider"
weight: 7
description: >
  This page describes how to add an analysis provider for doing deployment analysis.
---


You define the information needed to connect from your Piped to the Analysis Provider:
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

The full list of configurable fields are [here](/docs/operator-manual/piped/configuration-reference/#analysisprovider).

> TBA
