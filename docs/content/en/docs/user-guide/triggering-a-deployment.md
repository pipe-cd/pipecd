---
title: "Triggering a deployment"
linkTitle: "Triggering a deployment"
weight: 3
description: >
  This page describes when a deployment is triggered automatically and how to manually trigger a deployment.
---

PipeCD is using Git as a single source of truth, all application resources and infrastructure changes should be done by making a pull request to Git.
The mission of the deployments is syncing all running resources/components of applications in the cluster to the state specified in the newest commit.
So by default, when a new merged pull request touches an application, a new deployment for that application will be triggered to sync the application to the state specified in the newest merged commit.

A pull request (commit) is considered as touching an application whenever its changes include:
- one or more files inside the application directory
- one or more files inside one of the [dependencies](/docs/user-guide/configuration-reference/#kubernetesdeploymentinput) of the application

After a new deployment was triggered, it will be queued to handle by the appropriate `piped`. And at this time the deployment pipeline was not decided yet.
`piped` schedules all deployments of applications to ensure that for each application only one deployment will be executed at the same time.
When no deployment of an application is running, `piped` picks one queueing deployment for that application to plan the deploying pipeline.
`piped` plans the deploying pipeline based on the deployment configuration and the diff between the running state and the specified state in the newest commit.
For example:

- when the merged pull request updated a Deployment's container image or updated a mounting ConfigMap or Secret, `piped` planner will decide that the deployment should use the specified pipeline to do a progressive deployment.
- when the merged pull request just updated the `replicas` number, `piped` planner will decide to use a quick sync to scale the resources.

You can force `piped` planer to decide to use the [QuickSync](docs/concepts/#quick-sync) or the specified pipeline based on the commit message by configuring [CommitMatcher](/docs/user-guide/configuration-reference/#commitmatcher) in the deployment configuration.

After being planned, the deployment will be executed as the decided pipeline. The deployment execution including the state of each stage as well as their logs can be viewed in realtime at the deployment details page.

![](/images/deployment-details.png)
<p style="text-align: center;">
A Running Deployment at the Deployment Details Page
</p>

As explained above, by default all deployments will be triggered automatically by checking the merged commits but you also can manually trigger a new deployment from web UI.
By clicking on `SYNC` button at the application details page, a new deployment for that application will be triggered to sync the application to be the state specified at the newest commit of the master branch (default branch).

![](/images/application-details.png)
<p style="text-align: center;">
Application Details Page
</p>

