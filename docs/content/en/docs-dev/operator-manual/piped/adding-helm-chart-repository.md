---
title: "Adding a helm chart repository"
linkTitle: "Adding helm chart repo"
weight: 5
description: >
  This page describes how to add a new Helm chart repository.
---

PipeCD supports Kubernetes applications that are using Helm for templating and packaging. In addition to being able to deploy a Helm chart that is sourced from the same Git repository (`local chart`) or from a different Git repository (`remote git chart`), an application can use a chart sourced from a Helm chart repository.

A Helm [chart repository](https://helm.sh/docs/topics/chart_repository/) is a location backed by an HTTP server where packaged charts can be stored and shared. Before an application can be configured to use a chart from a Helm chart repository, that chart repository must be enabled in the related `piped` by adding the [ChartRepository](/docs/operator-manual/piped/configuration-reference/#chartrepository) struct to the piped configuration file.

``` yaml
# piped configuration file
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  chartRepositories:
    - name: pipecd
      address: https://charts.pipecd.dev
```

For example, the above snippet enables the official chart repository of PipeCD project. After that, you can configure the Kubernetes application to load a chart from that chart repository for executing the deployment.

``` yaml
# .pipe.yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  input:
    # Helm chart sourced from a Helm Chart Repository.
    helmChart:
      repository: pipecd
      name: helloworld
      version: v0.5.0
```

In case the chart repository is backed by HTTP basic authentication, the username and password strings are required in [configuration](/docs/operator-manual/piped/configuration-reference/#chartrepository).
