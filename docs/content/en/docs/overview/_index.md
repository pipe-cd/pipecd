---
title: "Overview"
linkTitle: "Overview"
weight: 1
description: >
  Overview about PipeCD.
---

{{% pageinfo %}}
PipeCD provides a Continuous Delivery solution for Declarative Kubernetes Application, Serverless Application and Infrastructure.
{{% /pageinfo %}}

## What Is PipeCD?

> TBA: Introduce your project, including what it does or lets you do, why you would use it, and its primary goal (and how it achieves it). This should be similar to your README description, though you can go into a little more detail here if you want.

![](/images/architecture-overview.png)

## Why PipeCD?

**Powerful**
- Unifed Deployment System: kubernetes (plain-yaml, helm, kustomize), terraform, lambda, cloudrun...
- Progressive Deployment Strategies: canary, bluegreen, rolling update
- Automated Analysis: by metrics, log, smoke test...
- Automated Rollback
- Automated Configuration Drift Detection
- Insights shows Delivery Perfomance
- Support Webhook and Slack notifications

**Easy to Use**
- Operations by Pull Request: scale, rollout, rollback by PR
- Realtime Visualization of application state
- Deployment Pipeline to see what is happenning
- Intuitive UI

**Easy to Operate**
- Two seperate components: single binary `piped` and `control-plane`
- `piped` can be run in a Kubernetes cluster, a single VM or even a local machine
- Easy to operate multi-tenancy, multi-cluster

**Safety and Security**
- Support single sign-on (SSO) and role-based access control (RBAC)
- Your credentials are not exposed outside your cluster and not saved in control-plane

**Open Source**

- Released as an Open Source project
- Under APACHE 2.0 license, see [LICENSE](https://github.com/pipe-cd/pipe/blob/master/LICENSE)

## Where should I go next?

If you are an **operator** who are wanting to install and configure PipeCD for other developers.
- [Quickstart](/docs/quickstart/)
- [Operator Manual](/docs/operator-manual/)

If you are an **user** who are using PipeCD to deploy your application/infrastructure:
- [User Guide](/docs/user-guide/)
- [Examples](/docs/user-guide/examples)

If you want to be a **contributor**:
- [Contributor Guide](/docs/contributor-guide/)
