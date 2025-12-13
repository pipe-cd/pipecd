---
title: "Adding an application"
linkTitle: "Adding an application"
weight: 1
description: >
  This page describes how to add a new application.
---

An application is a collection of resources and configurations that are managed together.
It represents the service which you are going to deploy. With PipeCD, all application's manifests and its application configuration (`app.pipecd.yaml`) must be committed into a directory of a Git repository. That directory is called as application directory.

Each application can be handled by one and only one `piped`. Currently, PipeCD is supporting 5 kinds of application: Kubernetes, Terraform, CloudRun, Lambda, ECS.

> Note: Be sure your application manifests repository is listed in [Piped managing repositories configuration](../managing-piped/configuration-reference/#gitrepository:~:text=No-,repositories,-%5B%5DRepository).

Before deploying an application, it must be registered to help PipeCD know:

- where the application configuration is placed
- which `piped` should handle it and which platform the application should be deployed to

Through the web console, you can register a new application in one of the following ways:

- Picking from a list of unused apps suggested by Pipeds while scanning Git repositories (Recommended)
- Manually configuring application information

(If you prefer to use [`pipectl`](../../command-line-tool/#adding-a-new-application) command-line tool, see its usage for the details.)

## Picking from a list of unused apps suggested by Pipeds

You have to __prepare a configuration file__ which contains your application configuration and store that file in the Git repository which your Piped is watching first to enable adding a new application this way.

The application configuration file name must be suffixed by `.pipecd.yaml` because Piped periodically checks for files with this suffix.

{{< tabpane >}}
{{< tab lang="yaml" header="KubernetesApp" >}}
# For application's configuration in detail for KubernetesApp, please visit
# https://pipecd.dev/docs/user-guide/managing-application/defining-app-configuration/kubernetes/

apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: foo
  labels:
    team: bar
{{< /tab >}}
{{< tab lang="yaml" header="TerraformApp" >}}
# For application's configuration in detail for TerraformApp, please visit
# https://pipecd.dev/docs/user-guide/managing-application/defining-app-configuration/terraform/

apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: foo
  labels:
    team: bar
{{< /tab >}}
{{< /tabpane >}}

To define your application deployment pipeline which contains the guideline to show Piped how to deploy your application, please visit [Defining app configuration](../defining-app-configuration/).

Go to the PipeCD web console on application list page, click the `+ADD` button at the top left corner of the application list page and then go to the `ADD FROM GIT` tab.

Select the Piped and Platform Provider that you deploy to, once the Piped that's watching your Git repository catches the new unregistered application configuration file, it will be listed up in this panel. Click `ADD` to complete the registration.

![Registering an Application from Suggestions: PipeCD v1](/images/add-from-suggestions-v1.png)
<p style="text-align: center;">
Registering an Application from Suggestions
</p>

## Updating an application

Regardless of which method you used to register the application, the web console can only be used to disable/enable/delete the application, besides the adding operation. All updates on application information must be done via the application configuration file stored in Git as a single source of truth.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: new-name
  labels:
    team: new-team
```

Refer to [configuration reference](../configuration-reference/) to see the full list of configurable fields.
