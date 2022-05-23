---
title: "Adding a helm chart registry"
linkTitle: "Adding helm chart registry"
weight: 6
description: >
  This page describes how to add a new Helm chart registry.
---

A Helm chart [registry](https://helm.sh/docs/topics/registries/) is a mechanism enabled by default in Helm 3.8.0 and later that allows the OCI registry to be used for storage and distribution of Helm charts.

Before an application can be configured to use a chart from a registry, that registry must be enabled in the related `piped` by adding the [ChartRegistry](/docs/operator-manual/piped/configuration-reference/#chartregistry) struct to the piped configuration file if authentication is enabled at the registry.

``` yaml
# piped configuration file
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  chartRegistries:
    - type: OCI
      address: registry.example.com
      username: sample-username
      password: sample-password
```
