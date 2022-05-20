---
title: "Adding a helm chart registry"
linkTitle: "Adding helm chart registry"
weight: 6
description: >
  This page describes how to add a new Helm chart registry.
---

TODO: write

``` yaml
# piped configuration file
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  chartRegistries:
    - type: OCI
      address: example.com
      username: sample-username
      password: sample-password
```
