---
title: "Concepts"
linkTitle: "Concepts"
weight: 2
description: >
  This page describes several core concepts in PipeCD.
---

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

### Environment
>Deprecated: Please use Labels instead.

An environment is a logical group of applications of a project. A project can have multiple environments.
Each application must belong to one and only one environment. While each piped must belong to at least one environment.

### Deployment

A deployment is a process that does transition from the current state (running state) to the desired state (specified state in Git) of a specific application.
When the deployment is success, it means the running state is synced with the desired state specified in the target commit.

### Application Configuration

A yaml file that contains configuration data to define how to deploy the application.
Each application requires one application configuration file at application directory in the Git repository.
The default file name is `app.pipecd.yaml`.

### Application Directory

A directory in Git repository containing application configuration file and application manifests.
Each application must have one application directory.

### Quick Sync

Quick sync is a fast way to sync application to the state specified in a Git commit without any progressive strategy or manual approving. Its pipeline contains just only one predefined `SYNC` stage. For examples:
- quick sync a Kubernetes application is just applying all manifests
- quick sync a Terraform application is automatically applying all detected changes
- quick sync a CloudRun/Lambda application is rolling out the new version and routing all traffic to it

### Pipeline

A list of stages specified by user in the application configuration file that tells `piped` how the application should be deployed. If the pipeline is not specified, the application will be deployed by Quick Sync way.

### Stage

A temporary middle state between current state and desired state of a deployment process.

### Cloud Provider

PipeCD supports multiple clouds and multiple kinds of applications.
Cloud Provider defines which cloud and where application should be deployed to.

Currently, PipeCD is supporting these five cloud providers: `KUBERNETES`, `ECS`, `TERRAFORM`, `CLOUDRUN`, `LAMBDA`.

### Analysis Provider
An external product that provides metrics/logs to evaluate deployments, such as `Prometheus`, `Datadog`, `Stackdriver`, `CloudWatch` and so on.
It is mainly used in the [Automated deployment analysis](/docs/user-guide/automated-deployment-analysis/) context.
