---
title: "Adding an application"
linkTitle: "Adding an application"
weight: 1
description: >
  This page describes how to add a new application.
---

An application is a collect of resources and configurations that are managed together.
It represents the service which you are going to deploy. With PipeCD, all application's manifests and its deployment configuration (`.piped.yaml`) must be committed into a directory of a Git respository. That directory is called as application configuration directory.

Before deploying an application, the application must be registered from the web UI and a deployment configuration file (`.piped.yaml`) must be committed to the application configuration directory.
An application must belong to exactly one environment and can be handled by one of the registered `piped`s. Currently, PipeCD is supporting the following kinds of application:

- Kubernetes application
- Terraform application
- CloudRun application
- Lambda application

## Registering a new application from Web UI

Registering application helps PipedCD know the basic information about that application, where the application configuration is placing, what `piped` should handle it as well as what cloud the application should be deployed to.

By clicking on `+ADD` button at the application list page, a popup will be revealed from the right side as below:

![](/images/registering-an-application.png)
<p style="text-align: center;">
Popup for registering a new application from Web UI
</p>

After filling all the required fields, click `Save` button to complete the application registering.

Here are the list of fields in the register form:

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

## Adding deployment configuration file

After registering the application, one more step left is adding the deployment configuration file (`.pipe.yaml`) for that application into the application configuration directory in Git repository.

While registering application helps PipeCD know the basic information about application, the deployment configuration file is used by `piped`, and it helps `piped` know how the application should be deployed, such as doing canary/bluegreen strategy or requiring a manual approval...
That deployment configuration file is in `YAML` format as below:

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: ApplicationKind
spec:
  ...
```

- `kind` is the application kind. As explained before, supporting kinds of application are: `Kubernetes`, `Terrform`, `CloudRun`, `Lambda`.
- `spec` is the specific configuration for each application kind.

After clicking on the `Save` button at the previous step, the popup will be changed to allow you fill your deployment configuration. You can also choose one of the prepared templates.

![](/images/adding-deployment-configuration-file.png)
<p style="text-align: center;">
Popup for registering a new application from Web UI
</p>

<br/>

The [next section](/docs/user-guide/configuring-deployment/) guides you how to configure the deployment for each specific application kinds.
