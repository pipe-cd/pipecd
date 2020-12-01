---
title: "FAQ"
linkTitle: "FAQ"
weight: 9
description: >
  List of frequently asked questions.
---

If you have any other questions, please feel free to create the issue in the [pipe-cd/pipe](https://github.com/pipe-cd/pipe/issues/new/choose) repository or contact us on [Slack](https://join.slack.com/t/cloud-native/shared_invite/zt-fyy3b8up-qHeDNVqbz1j8HDY6g1cY4w).

### 1. What kind of application (cloud provider) will be supported?

Currently, PipeCD can be used to deploy `Kubernetes`, `Terraform`, `CloudRun`, `Lambda` applications.

In the near future we also want to support `ECS`, `Crossplane`...

### 2. What kind of templating methods for Kubernetes application will be supported?

Currently, PipeCD is supporting `Helm` and `Kustomize` as templating method for Kubernetes applications.

### 3. Istio is supported now?

Yes, you can use PipeCD for both mesh (Istio, SMI) applications and non-mesh applications.

### 4. What are the differences between PipeCD and FluxCD?

- Not just Kubernetes applications, PipeCD also provides a unified interface for other cloud services (CloudRun, AWS Lamda...) and Terraform
- One tool for both GitOps sync and progressive deployment
- Supports multiple Git repositories
- Has web UI for better visibility
    - Log viewer for each deployment
    - Visualization of application component/state in realtime
    - Show configuration drift in realtime
- Also supports Canary and BlueGreen for non-mesh applications
- Has built-in secrets management
- Supports gradual rollout of a single app to multiple clusters
- Shows the delivery performance  insights

### 5. What are the differences between PipeCD and ArgoCD?

- Not just Kubernetes applications, PipeCD also provides a unified interface for other cloud services (CloudRun, AWS Lamda...) and Terraform
- One tool for both GitOps sync and progressive deployment
- Don't need another CRD or changing the existing manifests for doing Canary/BlueGreen. PipeCD just uses the standard Kubernetes deployment object
- Easier and safer to operate multi-tenancy, multi-cluster for multiple teams (even some teams are running in a private/restricted network)
- Has built-in secrets management
- Supports gradual rollout of a single app to multiple clusters
- Shows the delivery performance  insights
