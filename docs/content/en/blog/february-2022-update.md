---
date: 2022-02-10
title: "February 2022 update"
linkTitle: "February 2022 update"
weight: 995
description: "Development status update to recap what happened in January"
author: Le Van Nghia ([@nghialv](https://twitter.com/nghialv2607))
---

_Published by the PipeCD dev team every month, this update will provide you with news and updates about the project! Please click [here](/blog/2022/01/05/january-2022-update/) if you want to see the last status update._

### What's changed
---

Since the last report, PipeCD team has introduced 4 releases ([v0.24.0](https://github.com/pipe-cd/pipecd/releases/tag/v0.24.0), [v0.24.5](https://github.com/pipe-cd/pipecd/releases/tag/v0.24.5), [v0.25.0](https://github.com/pipe-cd/pipecd/releases/tag/v0.25.0), [v0.25.1](https://github.com/pipe-cd/pipecd/releases/tag/v0.25.1)). This blog post walks you through their notable changes. For all other changes, please check out each release note.

#### Introducing Label mechanism

Label concept has been introduced. Labels are key/value pairs that are attached to PipeCD resources such as application, deployment... As a developer, you can use labels to add identifying attributes to them. We believe that supporting the Label will give us more flexibility on grouping, filtering as well as managing access control on PipeCD resources.

Application labels can be specified via the `labels` field in the application configuration file while deployment labels are inherited from its application.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: example
  labels:
    env: prod
    team: product
```

From `v0.24.0`, application labels are shown on the application list and detail page. All filter forms also have a new label box to allow filtering resources by their labels.

#### Environment and Deployment Configuration File were deprecated

Since the Label was introduced, the Environment becomes a subset of the Label concept. So we decided to deprecate the Environment concept from `v0.24.0`. Even though it is still displayed on the web console, it will be removed completely from the next release. So make sure to switch to using Label instead.

The deployment configuration file (`.pipe.yaml`) has been deprecated as well. Please use the application configuration file instead, it has the same format but with some new fields such as `name`, `description`, and `labels`. By using this new config file, all application information can be stored and managed in your Git repository. The new config file also gives PipeCD the ability to detect unregistered applications in Git repository to suggest users to add on the control plane.

Basically, a deployment configuration file can be migrated to be an application configuration file just by adding the `name` field.

```yaml
# Old deployment configuration file
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  ...
```

```yaml
# New application configuration file
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: example # Add this field.
  labels:
    env: prod
    team: product
  ...
```

If you are having many applications, you can use this [pipectl command](/docs-v0.25.x/user-guide/command-line-tool/#migrating-deployment-configuration-files-to-application-configuration-files) to migrate a bunch of your files seamlessly.

#### Event list page on web console

The event concept was designed to help external services such as CI systems to be able to interact with PipeCD. For instance, you can use the [Event Watcher](/docs/user-guide/event-watcher/
) feature to trigger a new deployment of a given application when a new version of its container image or Helm chart or Terraform package has been just published.

Checking the status of each sent event could be helpful for the developer while troubleshooting, so from `v0.25.0` that page was added to PipeCD’s web console.

![](/images/event-list-page.png)

#### Automated configuration drift detection for Cloud Run application

Configuration Drift is a phenomenon where running resources of service become more and more different from the definitions in Git as time goes on, due to manual ad-hoc changes and updates. As PipeCD is using Git as a single source of truth, all application resources and infrastructure changes should be done by making a pull request to Git. Whenever a configuration drift occurs it should be notified to the developers and be fixed.

PipeCD includes the Configuration Drift Detection feature, which automatically detects the configuration drift and shows the result in the application details web page as well as sends the notifications to the developers.

Before `v0.25.0`, this feature could only be used for Kubernetes apps, but now we're happy to announce that Cloud Run apps can also have that feature.

![](/images/cloud-run-out-of-sync.png)

### What are next
---

The team continues actively working on improving the PipeCD product. Besides fixing the reported issues, enhancing the existing features, here are some new features the team is currently working on:

- Realtime application state for Cloud Run application
- RBAC for PipeCD resources based on the Label mechanism
- Enable running control-plane without database. It means file storage (such as GCS, S3, Minio) can be used for both data store and file store.

If you have any features want to request or find out a problem, please let us know by creating issues to the [pipe-cd/pipecd](https://github.com/pipe-cd/pipecd/issues) repository.


---
*Follow us on Twitter to keep track of all the latest news: https://twitter.com/pipecd_dev*
