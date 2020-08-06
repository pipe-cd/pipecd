---
title: "Configuring an application"
linkTitle: "Configuring an application"
weight: 1
description: >
  This page describes how to add and configure an application.
---

An application is a collect of resources and configurations that are managed together.
It represents the service which you are going to deploy. With PipeCD all application configurations and its deployment configuration (`.piped.yaml`) must be committed to a directory of a Git respository. That directory is called as application configuration directory.

Before deploying an application, the application must be registered from the web UI and a deployment configuration file (`.piped.yaml`) must be committed to the application configuration directory.
An application must belong to exactly one environment and can be handled by one registered `piped`. Currently, PipeCD supports the following application kinds:

- Kubernetes application
- Terraform application
- CloudRun application
- Lambda application

## Registering a new application from Web UI

Registering application helps PipedCD know where the application configuration is placing, what `piped` should handle the application as well as what cloud the application should be deployed to.

By clicking on `+ADD` button at the application list page, a popup at the right side will be revealed as the following:

![](/images/registering-an-application.png)
<p style="text-align: center;">
Popup for registering a new application from Web UI
</p>

| Field | Description | Required |
|-|-|-|-|
| Name | The application name | Yes |
| Kind | The application kind. Select one of these values: `Kubernetes`, `Terraform`, `CloudRun`, `Lambda` | Yes |
| Env | The environment this application should belongs to. Select one of the registered environments at `Settings/Environment` page.  | Yes |
| Piped | The piped that handles this application. Select one of the registered `piped`s at `Settings/Piped` page. | Yes |
| Repository | The Git repository contains application configuration and deployment configuration. Select one of the registered repositories in `piped` configuration. | Yes |
| Path | The relative path from the root of the Git repository to the directory containing application configuration and deployment configuration. | Yes |
| Config Filename | The name of deployment configuration file. Default is `.pipe.yaml`. | No |
| Cloud Provider | Where the application will be deployed to. Select one of the registered cloud providers in `piped` configuration. | Yes |

After filling all the above fields, click `Save` button to complete application registering.

## Adding a deployment configuration file `.piped.yaml`

After registering the application, one more step left is adding the deployment configuration file (`.pipe.yaml`) to the application configuration directory in Git repository.
This deployment configuration specifies how application should be deployed such as canary/bluegreen strategy, required manual approval...

> TBA

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
```
