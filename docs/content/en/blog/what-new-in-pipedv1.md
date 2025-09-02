---
date: 2025-09-02
title: "What is new in pipedv1 (plugin-arch piped)"
linkTitle: "what is new in pipedv1"
weight: 979
description: ""
author: Khanh Tran ([@khanhtc1202](https://github.com/khanhtc1202))
categories: ["Announcement"]
tags: ["Plugin", "New Feature"]
---

Since the alpha release from Jun 2025, plugin-arch piped (aka. pipedv1) is closer to first offical release, we are working on v1.0.0-rc0 and will be released in the coming days.

In this article I would like to share some improvements that have been fixed in the new version pipedv1 compared to its predecessor piped.

## Overview

While developing pipedv1, we always considered ensuring backward compatibility, meaning that the features available in piped will be mostly preserved in pipedv1. This ensures seamlessness between using piped and pipedv1, reducing the risk of problems when switching to pipedv1.

Changes in pipedv1 are of 2 types:

- Performance improvements or addressing some features requested in pipedv0 but could not be met due to technical reasons.

- Removal of some features that are obsolete or no longer prioritized

Fundamental changes like pipedv1 being able to support different platforms and making build pipeline deployment more powerful by turning stage executors into plugins have been described in [this blog](https://pipecd.dev/blog/2024/11/28/overview-of-the-plan-for-pluginnable-pipecd/). This post will focus on smaller changes related to the actual usage of pipecd.

## New in pipedv1

### Pipeline planning and configuration

Creating pipelines based on different deployment strategies for various specific applications is a strong point that makes PipeCD stand out. With pipedv1, this feature is further enhanced by solving some old limitations in the current piped.

Since pipedv1, you can create a pipeline consisting of only one SYNC stage.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: test-app
  pipeline:
    stages:
      - name: K8S_SYNC
        with:
          ...
```

Or optionally, create a pipeline that includes the SYNC stage combined with other common stages, for example.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: test-app
  pipeline:
    stages:
      - name: WAIT
        with:
          duration: 30s
      - name: K8S_SYNC
      - name: WAIT
        with:
          duration: 30s
```

### Stages configurations

Regarding stages (units of the pipecd pipeline), a feature that has been requested by many users is `skipOn`, which allows skipping a stage in the application deployment pipeline under certain conditions.

In the current version of piped, skipOn is only supported in limited ways for some stages (like Analysis Stage - ref: [docs](https://pipecd.dev/docs-v0.53.x/user-guide/configuration-reference/#analysisstageoptions)), in pipedv1, all stages support this feature.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: test-app
  pipeline:
    stages:
      - name: WAIT
        skipOn:
          paths:
          - '*/canary'
        with:
          duration: 30s
```

An unstable feature in piped is stage timeout, in pipedv1 version, all stages also support timeout settings. The default stage timeout will be 6 hours.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: test-app
  pipeline:
    stages:
      - name: WAIT_APPROVAL
        timeout: 1h # Default is 6h
        with:
          approvers:
          - khanhtc1202
```

## Deprecated / Dis-supported features

### "Weird" flag in piped execution

A long-time "weird" flag for piped binary execution if you don't use PipeCD for a Kubernetes application is `--enable-default-kubernetes-cloud-provider`, which is basically legacy due to the first implementation of PipeCD focusing on supporting only Kubernetes. The default value is `false`, but if you set that flag to `true`, piped execution requires connecting to the Kubernetes cluster even before any application deployment is actually triggered. In pipedv1, nothing like that flag exists in piped execution because pipedv1 supports whatever platform its plugins support equally.

### Kubernetes templating feature

Related to Kubernetes support features of PipeCD. Currently, piped supports the Helm Git Remote Chart feature (ref [piped config helm chart repository](https://pipecd.dev/docs-v0.53.x/user-guide/managing-piped/configuration-reference/#chartrepository)), which pulls the chart file stored directly on Git to the local and builds/templates as for a Local Chart. With the emergence and standardization of OCI, storing and sharing Helm Charts has become easier via the [Helm registry](https://helm.sh/docs/helm/helm_registry/). Therefore, in pipedv1, we decided to stop supporting this feature.

### Analysis stage query templating feature

Current piped support building queries used while evaluating metrics with deployment-specific data to be embedded in the analysis template (ref: [analysis templating docs](https://pipecd.dev/docs-v0.53.x/user-guide/managing-application/customizing-deployment/automated-deployment-analysis/#optional-analysis-template)).

From pipedv1, along with built-in and custom args are supported with placeholders as `{{ .App }}` and `{{ .AppCustomArgs }}` respectively, Kubernetes-specific built-in args like `{{ .K8s.Namespace }}` will be marked as deprecated and unsupported after several releases. The corresponding usage for the Kubernetes Namespace use case is changed to `{{ .AppCustomArgs.k8sNamespace }}`.

Here are some changes you might notice when switching to pipedv1. We ensure a certain level of backward compatibility between piped and pipedv1. The improvements are all aimed at making pipedv1 support more platforms and making it easier to build pipelines based on plugins.

The official documentation for pipedv1 is still being prepared, and the experimental release of pipedv1 will be available on the [official pipecd repo release tab](https://github.com/pipe-cd/pipecd/releases), along with some built-in plugins in the next few days. Thanks for your attention, cheer 🍻

