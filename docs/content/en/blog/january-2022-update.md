---
date: 2022-01-05
title: "January 2022 update"
linkTitle: "January 2022 update"
weight: 997
description: "Development status update to recap what happened in December"
author: Le Van Nghia ([@nghialv](https://twitter.com/nghialv2607))
categories: ["Announcement"]
tags: ["News"]
---

_Published by the PipeCD dev team every month, this update will provide you with news and updates about the project! Please click [here](/blog/2021/11/01/november-2021-update/) if you want to see the last status update._

### Happy New Year
---
First of all, PipeCD team would like to wish you all a very Happy New Year. 🥳

2021 has been a great year for both PipeCD project and PipeCD team. Many features have been introduced by more than 2000 pull requests from 33 contributors. We would like to thank all those wonderful contributors.

![](/images/january-2022-contributor-list.png)

Stepping into 2022, PipeCD team looks forward to contributing even more to make PipeCD project better and more useful for many users in our OSS community.

### What's changed
---

Since the last report, PipeCD team has introduced 2 releases ([v0.22.0](https://github.com/pipe-cd/pipecd/releases/tag/v0.22.0), [v0.23.0](https://github.com/pipe-cd/pipecd/releases/tag/v0.23.0)). These releases bring many updates to help PipeCD becomes more stable as well as introduce some interesting features. In this blog post, we will recap those new features. For all other changes, please check out each release note to see.

#### Deployment chain

As you know, each application in PipeCD contains a group of resources that are managed together and be deployed to a single cloud provider such as a Kubernetes cluster. The application works independently of each other since there is no connection between them.

But in some use cases, you might want to build a more complex deployment flow where the connection between applications is required. That feature has been requested from a number of PipeCD users.

By [v0.23.0](https://github.com/pipe-cd/pipecd/releases/tag/v0.23.0), PipeCD team has brought that feature to reality. It will allow you to roll out applications to multiple clusters/regions gradually or to promote applications across environments. Those deployment strategies can be done simply by specifying a chain of applications in the `postSync` field of the first application.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  ...
  postSync:
    chain:
      applications:
        - name: second-app
          kind: KUBERNETES
        - name: third-app
          kind: CLOUDRUN
```

Currently, in-chain applications can be specified only by `name`, `kind` but finding applications by `label` will also be supported in the near future as well. Since PipeCD is supporting multiple application kinds (such as Kubernetes, Terraform, CloudRun...) so applications from different kinds are allowed to be included in a chain.

For more details, please check out its [documentation](https://pipecd.dev/docs/user-guide/deployment-chain/) page.

#### Environment is no longer required and will be replaced by Label

Environment was changed to be optional while registering Piped or Application. You still can use it from the web console or specify the environment for an application via the `envName` field in the application configuration file, but we are planning to completely remove the environment concept in the near future. As an alternative, a new `Label` concept has been introduced, and `Environment` can be imagined as a particular label.

An application can contain one or more labels in its application configuration file as below:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: TerraformApp
spec:
  name: simple-app
  envName: production
  labels:
    #env: production
    service: backend
  ...
```

Filtering application or deployment by the label is being implemeted currently. We hope that it will be shipped in the next release.

#### Simplifying the way of registering an application

Before [v0.23.0](https://github.com/pipe-cd/pipecd/releases/tag/v0.23.0), while registering a new application user has to configure the application information manually via a registration form. For now, most of the information can be specified in the application configuration file. And Piped agents automatically find the un-registered applications to suggest to users. It means that all you have to do is just a few clicks on the web console.

![](/images/legacy-registering-an-application-from-git.png)
<p style="text-align: center;">
Picking from the suggested application list to register
</p>

Please note that to be suggested, the application configuration file must be suffixed by `.pipecd.yaml`. And the default name of that file has been changed to `app.pipecd.yaml` instead of `.pipe.yaml` as before.

#### Showing the connection status of Pipeds

The connection status of Piped agents to the control plane is shown on the Settings page in real-time. This gives the operator a quick look at the current status of a particular Piped agent is.

![](/images/january-2022-piped-connection-status.png)

#### More controls on triggering the deployment

By default, when a new merged pull request touches an application, a new deployment for that application will be triggered to execute the sync process.

From `v0.22.0`, a new [`trigger`](https://pipecd.dev/docs/user-guide/configuration-reference/#deploymenttrigger) field has been added to the application configuration file to allow users to customize the triggering condition. For example, using `onOutOfSync` to enable the ability to attempt to resolve `OUT_OF_SYNC` state whenever a configuration drift has been detected.

The deployment triggering of each application can be configured with the following fields:

* `onCommit`: controls triggering new deployment when newly added Git commits touched the application.
* `onCommand`: controls triggering new deployment when received a new `SYNC` command from the web console or `pipectl`
* `onOutOfSync`: controls triggering new deployment when application is at `OUT_OF_SYNC` state. Enabling this will force Piped to always attempt to keep the application as synced as possible.
* `onChain`: controls triggering new deployment when the application is counted as a node of a deployment chain. (this configuration is added from `v0.23.0`)

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  trigger:
    onCommit:
      disabled: false
    onCommand:
      disabled: false
    onOutOfSync:
      disabled: true
    onChain:
      disabled: true
```

### What are next
---

The team continues actively working on improving the PipeCD product. Besides fixing the reported issues, enhancing the existing features, here are some new features the team is currently working on:

- Reduce the maintenance cost of the control plane by supporting using file storage (such as GCS, S3, Minio) for both data store and file store. It means **no database** is required to run the control plane.
- Automated configuration drift detection for Cloud Run application

If you have any features want to request or find out a problem, please let us know by creating issues to the [pipe-cd/pipecd](https://github.com/pipe-cd/pipecd/issues) repository.


---
*Follow us on Twitter to keep track of all the latest news: https://twitter.com/pipecd_dev*
