---
title: "Concepts"
linkTitle: "Concepts"
weight: 2
description: >
  This page describes several core concepts in PipeCD.
---

![](/images/architecture-overview.png)
<p style="text-align: center;">
Component Architecture
</p>

### Piped

`piped` is a single binary component you run as an agent in your cluster, your local network to handle the deployment tasks.
It can be run inside a Kubernetes cluster by simply starting a Pod or a Deployment.
This component is designed to be stateless, so it can also be run in a single VM or even your local machine.

### Control Plane

A centralized component managing deployment data and provides gPRC API for connecting `piped`s as well as all web-functionalities of PipeCD such as
authentication, showing deployment list/details, application list/details, delivery insights...

### Project

A project is a logical group of applications to be managed by a group of users.
Each project can have multiple `piped` instances from different clouds or environments.

There are three types of project roles:

- **Viewer** has only permissions of viewing to deployment and application in the project.
- **Editor** has all viewer permissions, plus permissions for actions that modify state such as manually trigger/cancel the deployment.
- **Admin** has all editor permissions, plus permissions for managing project data, managing project `piped`.

### Application

A collect of resources (containers, services, infrastructure components...) and configurations that are managed together.
PipeCD supports multiple kinds of applications such as `KUBERNETES`, `TERRAFORM`, `ECS`, `CLOUDRUN`, `LAMBDA`...

### Application Configuration

A YAML file that contains information to define and configure application.
Each application requires one file at application directory stored in the Git repository.
The default file name is `app.pipecd.yaml`.

### Application Directory

A directory in Git repository containing application configuration file and application manifests.
Each application must have one application directory.

### Deployment

A deployment is a process that does transition from the current state (running state) to the desired state (specified state in Git) of a specific application.
When the deployment is success, it means the running state is being synced with the desired state specified in the target commit.

### Sync Strategy

There are 3 strategies that PipeCD supports while syncing your application state with its configuration stored in Git. Which are:
- Quick Sync: a fast way to make the running application state as same as its Git stored configuration. The generated pipeline contains only one predefined `SYNC` stage.
- Pipeline Sync: sync the running application state with its Git stored configuration through a pipeline defined in its application configuration.
- Auto Sync: depends on your defined application configuration, `piped` will decide the best way to sync your application state with its Git stored configuration.

### Platform Provider

Note: The previous name of this concept was Cloud Provider.

PipeCD supports multiple platforms and multiple kinds of applications.
Platform Provider defines which platform, cloud and where application should be deployed to.

Currently, PipeCD is supporting these five platform providers: `KUBERNETES`, `ECS`, `TERRAFORM`, `CLOUDRUN`, `LAMBDA`.

### Analysis Provider
An external product that provides metrics/logs to evaluate deployments, such as `Prometheus`, `Datadog`, `Stackdriver`, `CloudWatch` and so on.
It is mainly used in the [Automated deployment analysis](../user-guide/managing-application/customizing-deployment/automated-deployment-analysis/) context.
