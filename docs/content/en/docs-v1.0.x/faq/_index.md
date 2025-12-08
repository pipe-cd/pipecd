---
title: "FAQ"
linkTitle: "FAQ"
weight: 10
description: >
  List of frequently asked questions.
---

We have answered some of the most frequently asked questions below. If you have any other questions, please feel free to create the issue in the [pipe-cd/pipecd](https://github.com/pipe-cd/pipecd/issues/new/choose) repository or contact us on [Cloud Native Slack](https://slack.cncf.io) (channel [#pipecd](https://app.slack.com/client/T08PSQ7BQ/C01B27F9T0X)).

### 1. What is PipeCD V1? How is it different from the PipeCD Versions so far?

PipeCD V1 introduces a plugin-based architecture where each application deployment is managed by a 'plugin', created specifically for that application. This replaces the concept of Platform Providers from earlier versions of PipeCD. This change makes PipeCD versatile, allowing users to create custom plugins to deploy the application of their choice.

### 2. What kind of applications will be supported in PipeCD V1?

Since PipeCD V1 introduces a plugin architecture, you can now deploy any application using plugins.

Check out the latest releases on GitHub for the list of available plugins. Additionally, we also have a Community Plugins Repository for plugins made by the PipeCD community. As of now, the official plugins maintained by the PipeCD Maintainers are Kubernetes, Terraform, Analysis, ScriptRun, Wait, and WaitApproval.

The broader plan in the future releases is to add plugins for Amazon ECS and GCP Cloud Run, which will be maintained by PipeCD, while plugins for other applications will go in the Community Plugins Repository.

### 3. What kind of templating methods for Kubernetes application will be supported?

Currently, PipeCD is supporting `Helm` and `Kustomize` as templating method for Kubernetes applications.

### 4. Is Istio is supported now?

Yes, you can use PipeCD for both mesh (Istio, SMI) applications and non-mesh applications.

### 5. What are the differences between PipeCD and FluxCD?

- Apart from Kubernetes applications, PipeCD also provides a unified interface for other cloud services (GCP Cloud Run, AWS ECS, AWS Lambda and more). Starting PipeCD V1, users can use PipeCD with even more applications by creating custom plugins for their deployments.
Here are some standout features of PipeCD when compared to Flux:

- One tool for both GitOps sync and progressive deployment
- Supports multiple Git repositories
- Has web UI for better visibility
  - Log viewer for each deployment
  - Visualization of application component/state in realtime
  - Show configuration drift in realtime
- Also supports Canary and BlueGreen for non-mesh applications
- Has built-in secrets management
- Shows the delivery performance insights

### 6. What are the differences between PipeCD and ArgoCD?

- Apart from Kubernetes applications, PipeCD also provides a unified interface for other cloud services (GCP Cloud Run, AWS ECS, AWS Lambda and more). Starting PipeCD V1, users can use PipeCD with even more applications by creating custom plugins for their deployments.
Here are some standout features of PipeCD when compared to ArgoCD:

- One tool for both GitOps sync and progressive deployment
- Don't need another CRD or changing the existing manifests for doing Canary/BlueGreen. PipeCD just uses the standard Kubernetes deployment object
- Easier and safer to operate multi-tenancy, multi-cluster for multiple teams (even some teams are running in a private/restricted network)
- Has built-in secrets management
- Shows the delivery performance  insights

### 7. What should I do if I lose my Piped key?

You can create a new Piped key. Go to the `Piped` tab at `Settings` page, and click the vertical ellipsis of the Piped that you would like to create the new Piped key. Don't forget deleting the old Key, too.

### 8. What is the strong point if PipeCD is used only for Kubernetes?

- Simple interface, easy to understand no extra CRD required
- Easy to install, upgrade, and manage (both the ControlPlane and the agent Piped)
- Not strict depend on any Kubernetes API, not being part of issues for your Kubernetes cluster versioning upgrade
- Easy to interact with any CI; Plan preview feature gives you an early look at what will be changed in your cluster even before manifests update
- Insights show metrics like lead time, deployment frequency, MTTR, and change failure rate to measure delivery performance

### 9. Is PipeCD open source?

Yes, PipeCD is fully open source project with APACHE LICENSE, VERSION 2.0

### 10. How should I investigate high CPU usage or memory usage in piped, or when OOM occurs?

If you're noticing high CPU usage, memory usage, or facing OOM issues in Piped, you can use the built-in support for `pprof`, a tool for visualization and analysis of profiling data.  
`pprof` can help you identify the parts of your application that are consuming the most resources. For more detailed information and examples of how to use `pprof` in Piped, please see [Using Pprof in Piped](../managing-piped/using-pprof-in-piped).
