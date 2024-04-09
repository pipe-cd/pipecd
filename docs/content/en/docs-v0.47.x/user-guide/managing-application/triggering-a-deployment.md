---
title: "Triggering a deployment"
linkTitle: "Triggering a deployment"
weight: 4
description: >
  This page describes when a deployment is triggered automatically and how to manually trigger a deployment.
---

PipeCD uses Git as a single source of truth; all application resources are defined declaratively and immutably in Git. Whenever a developer wants to update the application or infrastructure, they will send a pull request to that Git repository to propose the change. The state defined in Git is the desired state for the application and infrastructure running in the cluster. 

PipeCD applies the proposed changes to running resources in the cluster by triggering needed deployments for applications. The deployment mission is syncing all running resources of the application in the cluster to the state specified in the newest commit in Git.

By default, when a new merged pull request touches an application, a new deployment for that application will be triggered to execute the sync process. But users can configure the application to control when a new deployment should be triggered or not. For example, using [`onOutOfSync`](#trigger-configuration) to enable the ability to attempt to resolve `OUT_OF_SYNC` state whenever a configuration drift has been detected. 

### Trigger configuration

Configuration for the trigger used to determine whether we trigger a new deployment. There are several configurable types:
- `onCommit`: Controls triggering new deployment when new Git commits touched the application.
- `onCommand`: Controls triggering new deployment when received a new `SYNC` command.
- `onOutOfSync`: Controls triggering new deployment when application is at `OUT_OF_SYNC` state.
- `onChain`: Controls triggering new deployment when the application is counted as a node of some chains.

See [Configuration Reference](../../configuration-reference/#deploymenttrigger) for the full configuration.

After a new deployment was triggered, it will be queued to handle by the appropriate `piped`. And at this time the deployment pipeline was not decided yet.
`piped` schedules all deployments of applications to ensure that for each application only one deployment will be executed at the same time.
When no deployment of an application is running, `piped` picks queueing one to plan the deploying pipeline.
`piped` plans the deploying pipeline based on the application configuration and the diff between the running state and the specified state in the newest commit.
For example:

- when the merged pull request updated a Deployment's container image or updated a mounting ConfigMap or Secret, `piped` planner will decide that the deployment should use the specified pipeline to do a progressive deployment.
- when the merged pull request just updated the `replicas` number, `piped` planner will decide to use a quick sync to scale the resources.

You can force `piped` planner to decide to use the [QuickSync](../../../concepts/#sync-strategy) or the specified pipeline based on the commit message by configuring [CommitMatcher](../../configuration-reference/#commitmatcher) in the application configuration.

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

