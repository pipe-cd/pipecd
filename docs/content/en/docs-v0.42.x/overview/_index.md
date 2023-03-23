---
title: "Overview"
linkTitle: "Overview"
weight: 1
description: >
  Overview about PipeCD.
---

![](/images/pipecd-explanation.png)
<p style="text-align: center;">
PipeCD - a Gitops style continuous delivery solution
</p>

## What Is PipeCD?

{{% pageinfo %}}
PipeCD provides a unified continuous delivery solution for multiple application kinds on multi-cloud that empowers engineers to deploy faster with more confidence, a GitOps tool that enables doing deployment operations by pull request on Git.
{{% /pageinfo %}}

## Why PipeCD?

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

For a good understanding of the PipeCD's components, see the [Concepts](../concepts) page.

If you are an **operator** wanting to install and configure PipeCD for other developers.
- [Quickstart](../quickstart/)
- [Managing Control Plane](../user-guide/managing-controlplane/)
- [Managing Piped](../user-guide/managing-piped/)

If you are a **user** using PipeCD to deploy your application/infrastructure:
- [User Guide](../user-guide/)
- [Examples](../user-guide/examples)

If you want to be a **contributor**:
- [Contributor Guide](../contribution-guidelines/)
