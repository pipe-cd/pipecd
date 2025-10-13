---
title: "Mirgrating from v0 to v1"
linkTitle: "Mirgrating from v0 to v1"
weight: 9
description: >
  Documentation on migrating PipeCD v0 deployments to PipeCD v1
---

This page explains how to migrate your existing PipeCD system to **piped v1**, the new plugin-based architecture that brings modularity and extensibility to PipeCD.

## Overview

PipeCD v1 (internally referred to as *pipedv1*) introduces a **pluggable architecture** that allows developers to add and maintain custom deployment and operational plugins without modifying the core system of PipeCD.

Migration from v0 is designed to be **safe** and **incremental**, allowing you to switch between piped and pipedv1 during the process with minimal disruption.

## Components

| Component | Description | Compatibility |
|------------|--------------|----------------|
| **Control Plane** | Manages projects, deployments, and applications. | Supports both piped and pipedv1 concurrently. |
| **Piped** | Manages the actual deployment and syncing of applications. | Backward compatible - You can switch between versions safely. |

---

## Prerequisites

Before you start, ensure that:

- You are running PipeCD **v0.54.0-rc1 or later**.
- You have the **latest Control Plane** installed.
- You have **pipectl v0.54.0-rc1 or later**.
- You have access to your Control Plane with **API write permissions**.

> **Note:** If you’re new to the plugin architecture, read:
> - [Overview of the plan for plugin-enabled PipeCD](https://pipecd.dev/blog/2024/11/28/overview-of-the-plan-for-pluginnable-pipecd/)
> - [What’s new in pipedv1](https://pipecd.dev/blog/2025/09/02/what-is-new-in-pipedv1-plugin-arch-piped/)

---

## Migration Process Overview

The migration flow involves the following steps:

1. Update `piped` and `pipectl` binaries.  
2. Convert application configurations to the v1 format.  
3. Update the application model in the Control Plane database.  
4. Update piped configuration for v1 (plugins).  
5. Deploy and verify pipedv1.  
6. Optionally, switch back to the old piped.

---

## 1. Update piped and pipectl

Install or upgrade to `piped` and `pipectl` **v0.54.0-rc1 or newer**:

```bash
# Example for upgrading pipectl
curl -Lo ./pipectl https://github.com/pipe-cd/pipecd/releases/download/v0.54.0-rc1/pipectl_<OS>_<ARCH>
chmod +x ./pipectl
mv ./pipectl /usr/local/bin/





## What Is PipeCD?

{{% pageinfo %}}
PipeCD provides a unified continuous delivery solution for multiple application kinds on multi-cloud that empowers engineers to deploy faster with more confidence, a GitOps tool that enables doing deployment operations by pull request on Git.
{{% /pageinfo %}}

## Why PipeCD?

- Simple, unified and easy to use but powerful pipeline definition to construct your deployment
- Same deployment interface to deploy applications of any platform, including Kubernetes, Terraform, GCP Cloud Run, AWS Lambda, AWS ECS
- No CRD or applications' manifest changes are required; Only need a pipeline definition along with your application manifests
- No deployment credentials are exposed or required outside the application cluster
- Built-in deployment analysis as part of the deployment pipeline to measure impact based on metrics, logs, emitted requests
- Easy to interact with any CI; The CI tests and builds artifacts, PipeCD takes the rest
- Insights show metrics like lead time, deployment frequency, MTTR and change failure rate to measure delivery performance
- Designed to manage thousands of cross-platform applications in multi-cloud for company scale but also work well for small projects

## PipeCD's Characteristics in detail

**Visibility**
- Deployment pipeline UI shows clarify what is happening
- Separate logs viewer for each individual deployment
- Realtime visualization of application state
- Deployment notifications to slack, webhook endpoints
- Insights show metrics like lead time, deployment frequency, MTTR and change failure rate to measure delivery performance

**Automation**
- Automated deployment analysis to measure deployment impact based on metrics, logs, emitted requests
- Automatically roll back to the previous state as soon as analysis or a pipeline stage fails
- Automatically detect configuration drift to notify and render the changes
- Automatically trigger a new deployment when a defined event has occurred (e.g. container image pushed, helm chart published, etc)

**Safety and Security**
- Support single sign-on and role-based access control
- Credentials are not exposed outside the cluster and not saved in the Control Plane
- Piped makes only outbound requests and can run inside a restricted network
- Built-in secrets management

**Multi-provider & Multi-Tenancy**
- Support multiple application kinds on multi-cloud including Kubernetes, Terraform, Cloud Run, AWS Lambda, Amazon ECS
- Support multiple analysis providers including Prometheus, Datadog, Stackdriver, and more
- Easy to operate multi-cluster, multi-tenancy by separating Control Plane and Piped

**Open Source**

- Released as an Open Source project
- Under APACHE 2.0 license, see [LICENSE](https://github.com/pipe-cd/pipecd/blob/master/LICENSE)

## Where should I go next?

For a good understanding of the PipeCD's components.
- [Concepts](../concepts): describes each components.
- [FAQ](../faq): describes the difference between PipeCD and other tools.

If you are an **operator** wanting to install and configure PipeCD for other developers.
- [Quickstart](../quickstart/)
- [Managing Control Plane](../user-guide/managing-controlplane/)
- [Managing Piped](../user-guide/managing-piped/)

If you are a **user** using PipeCD to deploy your application/infrastructure:
- [User Guide](../user-guide/)
- [Examples](../user-guide/examples)

If you want to be a **contributor**:
- [Contributor Guide](../contribution-guidelines/)
