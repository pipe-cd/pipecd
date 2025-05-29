---
title: "Configuring Kubernetes application"
linkTitle: "Kubernetes"
weight: 1
description: >
  Specific guide to configuring deployment for Kubernetes application.
---

Based on the application configuration and the pull request changes, PipeCD plans how to execute the deployment: doing quick sync or doing progressive sync with the specified pipeline.

Note:

You can generate an application config file easily and interactively by [`pipectl init`](../../command-line-tool.md#generating-an-application-config-apppipecdyaml).


## Quick sync

Quick sync is a fast way to sync application to the state specified in the target Git commit without any progressive strategy. It just applies all the defined manifiests to sync the application.
The quick sync will be planned in one of the following cases:
- no pipeline was specified in the application configuration file
- [pipeline](../../../configuration-reference/#pipeline) was specified but the PR did not make any changes on workload (e.g. Deployment's pod template) or config (e.g. ConfigMap, Secret)

For example, the application configuration as below is missing the pipeline field. This means any pull request touches the application will trigger a quick sync deployment.

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  input:
    helmChart:
      repository: pipecd
      name: helloworld
      version: v0.3.0
```

In another case, even when the pipeline was specified, a PR that just changes the Deployment's replicas number for scaling will also trigger a quick sync deployment.

## Sync with the specified pipeline

The `pipeline` field in the application configuration is used to customize the way to do deployment by specifying and configuring the execution stages. You may want to configure those stages to enable a progressive deployment with a strategy like canary, blue-green, a manual approval, an analysis stage.

To enable customization, PipeCD defines three variants for each Kubernetes application: primary (aka stable), baseline and canary.
- `primary` runs the current version of code and configuration.
- `baseline` runs the same version of code and configuration as the primary variant. (Creating a brand-new baseline workload ensures that the metrics produced are free of any effects caused by long-running processes.)
- `canary` runs the proposed change of code or configuration.

Depending on the configured pipeline, any variants can exist and receive the traffic during the deployment process but once the deployment is completed, only the `primary` variant should be remained.

These are the provided stages for Kubernetes application you can use to build your pipeline:

- `K8S_PRIMARY_ROLLOUT`
  - update the primary resources to the state defined in the target commit
- `K8S_CANARY_ROLLOUT`
  - generate canary resources based on the definition of the primary resource in the target commit and apply them
- `K8S_CANARY_CLEAN`
  - remove all canary resources
- `K8S_BASELINE_ROLLOUT`
  - generate baseline resources based on the definition of the primary resource in the target commit and apply them
- `K8S_BASELINE_CLEAN`
  - remove all baseline resources
- `K8S_TRAFFIC_ROUTING`
  - split traffic between variants

and other common stages:
- `WAIT`
- `WAIT_APPROVAL`
- `ANALYSIS`

See the description of each stage at [Customize application deployment](../../customizing-deployment/).

## Manifest Templating

In addition to plain-YAML, PipeCD also supports Helm and Kustomize for templating application manifests.

A helm chart can be loaded from:
- the same git repository with the application directory, we call as a `local chart`

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  input:
    helmChart:
      path: ../../local/helm-charts/helloworld
```

- a different git repository, we call as a `remote git chart`

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  input:
    helmChart:
      gitRemote: git@github.com:pipe-cd/manifests.git
      ref: v0.5.0
      path: manifests/helloworld
```

- a Helm chart repository, we call as a `remote chart`

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  input:
    helmChart:
      repository: pipecd
      name: helloworld
      version: v0.5.0
```

A kustomize base can be loaded from:
- the same git repository with the application directory, we call as a `local base`
- a different git repository, we call as a `remote base`

See [Examples](../../../examples/#kubernetes-applications) for more specific.

## Reference

See [Configuration Reference](../../../configuration-reference/#kubernetes-application) for the full configuration.
