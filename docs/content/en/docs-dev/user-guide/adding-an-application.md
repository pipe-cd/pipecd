---
title: "Adding an application"
linkTitle: "Adding an application"
weight: 1
description: >
  This page describes how to add a new application.
---

An application is a collect of resources and configurations that are managed together.
It represents the service which you are going to deploy. With PipeCD, all application's manifests and its application configuration (`app.pipecd.yaml`) must be committed into a directory of a Git repository. That directory is called as application directory.

Before deploying an application, it must be registered via the web console.
Registering application helps PipeCD know where the application configuration is placed, which `piped` should handle it as well as which cloud the application should be deployed to.

Each application can be handled by one and only one `piped`. Currently, PipeCD is supporting the following application kinds:

- Kubernetes application
- Terraform application
- CloudRun application
- Lambda application
- ECS application

There are two ways to register an application:
- Scanning the unused application configuration files in Git to add (recommended)
- Manually configure application information

## From the application configuration in your Git repository (recommended)
In this way, you define all information in the application configuration defined in the Git repository and use it as a single source of truth.

It starts with creating an application configuration file as following and pushing it to the Git repository watched by a Piped with version v0.23.0 or higher.
The file name must be suffixed by `.pipecd.yaml` because Piped periodically checks for files with this suffix.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: foo
```

Visit [here](/user-guide/configuration-reference/) for supported fields.

After waiting for a while (it depends on the Piped's setting), click the `+ADD` button at the top left corner of the application list page and then go to the `ADD FROM GIT` tab.
Select the Piped and Cloud Provider that you deploy to, and the application you have just selected should appear as a candidate.
Click `ADD` to complete the registration.

![](/images/registering-an-application-from-git.png)
<p style="text-align: center;">
</p>

## From the web UI
In this way, you set the necessary information on the web.
By clicking on `+ADD` button at the application list page, a popup will be revealed from the right side as below:

![](/images/registering-an-application-from-web.png)
<p style="text-align: center;">
</p>

After filling all the required fields, click `Save` button to complete the application registering.

Here are the list of fields in the register form:

| Field | Description | Required |
|-|-|-|-|
| Name | The application name | Yes |
| Kind | The application kind. Select one of these values: `Kubernetes`, `Terraform`, `CloudRun`, `Lambda` and `ECS`. | Yes |
| Env | The environment this application should belongs to. Select one of the registered environments at `Settings/Environment` page.  | No |
| Piped | The piped that handles this application. Select one of the registered `piped`s at `Settings/Piped` page. | Yes |
| Repository | The Git repository contains application configuration and application configuration. Select one of the registered repositories in `piped` configuration. | Yes |
| Path | The relative path from the root of the Git repository to the directory containing application configuration and application configuration. Use `./` means repository root. | Yes |
| Config Filename | The name of application configuration file. Default is `app.pipecd.yaml`. | No |
| Cloud Provider | Where the application will be deployed to. Select one of the registered cloud providers in `piped` configuration. | Yes |

### Adding application configuration file

After registering the application, one more step left is adding the application configuration file for that application into the application directory in Git repository.

Adding application configuration file helps `piped` know how the application should be deployed, such as doing canary/blue-green strategy or requiring a manual approval...
That application configuration file is in `YAML` format as below:

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: ApplicationKind
spec:
  ...
```

- `kind` is the application kind. As explained before, supporting kinds of application are: `Kubernetes`, `Terrform`, `CloudRun`, `Lambda` and `ECS`.
- `spec` is the specific configuration for each application kind.

Please refer [pipecd/examples](/docs/user-guide/examples/) for the deployments being supported.

The [next section](/docs/user-guide/configuring-deployment/) guides you how to configure the deployment for each specific application kinds.

## Updating an application
Regardless of which method you used to register the application, the web console can only be used to disable/enable/delete the application, besides the adding operation. All updates on application information must be done via the application configuration file stored in Git.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: new-name
```