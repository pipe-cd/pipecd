---
title: "FAQ"
linkTitle: "FAQ"
weight: 9
description: >
  List of frequently asked questions.
---

If you have any other questions, please feel free to create the issue in the [pipe-cd/pipe](https://github.com/pipe-cd/pipe/issues/new/choose) repository or contact us on [Slack](https://cloud-native.slack.com/archives/C01B27F9T0X).

### 1. What kind of application (cloud provider) will be supported?

Currently, PipeCD can be used to deploy `Kubernetes`, `Terraform`, `CloudRun`, `Lambda` applications.

In the near future we also want to support `ECS`, `Crossplane`...

### 2. What kind of templating methods for Kubernetes application will be supported?

Currently, PipeCD is supporting `Helm` and `Kustomize` as templating method for Kubernetes applications.

### 3. Istio is supported now?

Yes, you can use PipeCD for both mesh (Istio, SMI) applications and non-mesh applications.
