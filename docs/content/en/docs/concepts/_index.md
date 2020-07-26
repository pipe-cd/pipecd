---
title: "Concepts"
linkTitle: "Concepts"
weight: 2
description: >
  This page describes serveral core concepts in PipeCD.
---

### Piped

`piped` (the ’d’ is for daemon) is a single binary component that you run in your cluster, your local network to handle the deployment tasks.
It can be run inside a Kubernetes cluster by simply starting a Pod or a Deployment.
This component is designed to be stateless so it can also be run in a single VM or even your local machine.

### Control Plane

A centralized component that manages deployment data and provides gPRC API for connecting `piped`s as well as all web-functionalities of PipeCD such as authentication, viewing deployment list/details, viewing application list/details, viewing delivery insights...

### Project

A project is a logical group of applications to be managed by a group of users.
Each project can have multiple `piped` instances from different clouds or environments.

There are three types of project roles:

- **Viewer** has only view permissions to deployment and application in the project.
- **Editor** has all viewer permissions, plus permissions for actions that modify state such as manually trigger/cancel the deployment.
- **Admin** has all user permissions, plus permissions for managing project data, managing project `piped`.

### Application

A collect of resources (containers, services...) and configuration that are managed together.

### Pipeline

A sequence of stages provided by PipeCD. A pipeline processes a transition from the current state (running state) to the desired state (specified state in Git) of a specific application.

### Stage

A temporary middle state between current state and desired state of a deployment process.

### Deployment Configuration

A `.pipe.yaml` yaml file that contains configuration data to define how to deploy the application. Each application has one deployment configuration file in Git repository at application directory.

### Cloud Provider

PipeCD supports multiple clouds and multiple kinds of applications.
Cloud Provider defines which cloud and where application should be deployed to.

Currently, PipeCD is supporting these 5 cloud providers: `KUBERNETES`, `TERRAFORM`, `CLOUDRUN`, `LAMBDA`, `ECS`.

### Analysis Provider

PipeCD supports multiple methods to automate the analysis process of your deployments. It can be by using metrics, logs or by checking the configured http requests.
Analysis Provider defines where to get those metrics/log data, like `Prometheus`, `Datadog`, `Stackdriver`, `CloudWatch`, and so on.

### Kubernetes Application Variant

Each Kubernetes application can has 3 variants: primary (aka stable), baseline, canary.
- `primary` runs the current version of code and configuration.
- `baseline` runs the same version of code and configuration as the primary variant. (Creating a brand new baseline workload ensures that the metrics produced are free of any effects caused by long-running processes.)
- `canary` runs the proposed changed of code or configuration.
