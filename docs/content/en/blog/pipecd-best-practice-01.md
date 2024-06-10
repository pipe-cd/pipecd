---
date: 2021-12-29
title: "PipeCD best practice 01 - operate your own PipeCD cluster"
linkTitle: "PipeCD best practice 01"
weight: 996
description: "This blog is a part of PipeCD best practice series, a guideline for you to operate your own PipeCD cluster."
author: Khanh Tran ([@khanhtc1202](https://twitter.com/khanhtc1202))
categories: ["Practice"]
tags: ["Control Plane"]
---

PipeCD is an open source project, you can freely use the released versions of PipeCD to create and operate a continuous delivery system for your service or company. [Quickstart](/docs/quickstart/) docs - a complete guide for you to install required components of PipeCD and deploy a simple application via PipeCD is a good point to get started. However, for the sake of the simplicity of that tutorial, some points to keep in mind when operating a PipeCD cluster have been omitted. You will review the PipeCD architecture again and get some tips on how to operate a PipeCD cluster in this post.

### Again, review the PipeCD architecture

![images](/images/architecture-overview-with-roles.png)
<p style="text-align: center;">
Component Architecture
</p>

> Note: Please refer to [concepts](/docs/concepts/) docs for definitions of PipeCD components such as Control Plane, Piped, Application, etc.

At a glance, PipeCD is composed of 2 components, the Control Plane and the Piped(s). As the figure above, depending on your task/role, you need to work with different components of PipeCD.

1. As a operator - platform team, most of your time working with PipeCD is to operate the PipeCD Control Plane. You may also need to cooperate with the product team, helping them install Piped(s) to their applications running cluster depending on your company team's structure. The space surrounded by a <span style="color: green;">green border</span> shows the operators working area.
2. As a developer - product team, most of your time working with PipeCD is to deploy, manage and observe your applications via application configuration files and PipeCD web console. You may also need to cooperate with the platform team to install Piped(s) to your applications running cluster if necessary. The space surrounded by a <span style="color: blue;">blue border</span> shows the developers working area.

### Tips to operate PipeCD cluster

1. Only a single Control Plane needs to be installed and will be operated by the operators (platform team). The installed Control Plane doesn't have to be in the same cluster as your applications running cluster or the Piped installed cluster. Just make sure the Piped(s) can connect to the Control Plane via outbound requests is enough.

2. Developers from the product team may need to self-install piped into their applications running cluster if the platform team does not manage those credentials.\
\
Once the Piped - the runner is installed successfully, developers only need to care about [app configuration](/docs/user-guide/adding-an-application/) - which defines the pipeline so that piped(s) can use them to deploy your applications. Interaction with the applications is mainly done via the PipeCD web console and the configuration files stored on Git.

> As an __operator__, PipeCD [control-plane](/docs/operator-manual/control-plane/) and [piped](/docs/operator-manual/piped/)(s) are what you should care about. As a __developer__, you should care about [piped](/docs/operator-manual/piped/) which installed in your applications running cluster and your [applications' PipeCD configurations](/docs/user-guide/adding-an-application/).

3. It is possible to use a single piped to manage deployments of all the applications to the clusters other than the one in which piped is being installed. However, as security best practices and avoid generating cluster credentials, __we highly recommend running each Piped inside each cluster and that Piped will only manage the applications on that cluster__. In case you want to deploy applications that are required to be installed across clusters/regions/environments, and want to follow this rule of thumb, please refer to the [deployment-chain](/docs/user-guide/deployment-chain/) feature.

4. Credentials used by piped(s) while deploying your applications should be managed by the developers of the production team themselves. You can refer to the [secret management](/docs/user-guide/secret-management/) feature supported by PipeCD for this purpose.
