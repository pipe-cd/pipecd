---
title: "FAQ"
linkTitle: "FAQ"
weight: 9
description: >
  List of frequently asked questions.
---

If you have any other questions, please feel free to create the issue in the [pipe-cd/pipecd](https://github.com/pipe-cd/pipecd/issues/new/choose) repository or contact us on [Cloud Native Slack](https://slack.cncf.io) (channel [#pipecd](https://app.slack.com/client/T08PSQ7BQ/C01B27F9T0X)).

### 1. What kind of application (platform provider) will be supported?

Currently, PipeCD can be used to deploy `Kubernetes`, `ECS`, `Terraform`, `CloudRun`, `Lambda` applications.

In the near future we also want to support `Crossplane`...

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
- Shows the delivery performance  insights

### 5. What are the differences between PipeCD and ArgoCD?

- Not just Kubernetes applications, PipeCD also provides a unified interface for other cloud services (GCP CloudRun, AWS Lamda...) and Terraform
- One tool for both GitOps sync and progressive deployment
- Don't need another CRD or changing the existing manifests for doing Canary/BlueGreen. PipeCD just uses the standard Kubernetes deployment object
- Easier and safer to operate multi-tenancy, multi-cluster for multiple teams (even some teams are running in a private/restricted network)
- Has built-in secrets management
- Shows the delivery performance  insights

### 6. What should I do if I lost my Piped key?

You can create a new Piped key. Go to the `Piped` tab at `Settings` page, and click the vertical ellipsis of the Piped that you would like to create the new Piped key. Don't forget deleting the old Key, too.

### 7. What is the strong point if PipeCD is used only for Kubernetes?

- Simple interface, easy to understand no extra CRD required
- Easy to install, upgrade, and manage (both the ControlPlane and the agent Piped)
- Not strict depend on any Kubernetes API, not being part of issues for your Kubernetes cluster versioning upgrade
- Easy to interact with any CI; Plan preview feature gives you an early look at what will be changed in your cluster even before manifests update
- Insights show metrics like lead time, deployment frequency, MTTR, and change failure rate to measure delivery performance

### 8. Is it open source?

Yes, PipeCD is fully open source project with APACHE LICENSE, VERSION 2.0!!

### 9. How should I investigate high CPU usage or memory usage in piped, or when OOM occurs?

If you're noticing high CPU usage, memory usage, or facing OOM issues in Piped, you can use the built-in support for `pprof`, a tool for visualization and analysis of profiling data.  
`pprof` can help you identify the parts of your application that are consuming the most resources. For more detailed information and examples of how to use `pprof` in Piped, please refer to our [Using Pprof in Piped guide](../managing-piped/using-pprof-in-piped).
