---
title: "Headlamp plugin"
linkTitle: "Headlamp plugin"
weight: 999
description: >
  This page describes how to view and sync your PipeCD applications from Headlamp using the PipeCD plugin.
---

[Headlamp](https://headlamp.dev/) is a web UI for Kubernetes. The PipeCD plugin for Headlamp lets you see your PipeCD pipeds, applications and deployments inside Headlamp, next to the cluster resources they create. You can check the sync status of an application, follow a deployment stage by stage, read stage logs, and start a sync without switching to the PipeCD web console.

The plugin works with every platform PipeCD supports, including Kubernetes, Terraform, Cloud Run, Lambda and ECS.

## What you can do

| Page | What it shows |
|------|---------------|
| Overview | Counts of pipeds, applications and deployments, the change failure rate, a deployments chart for the last 24 hours, 7, 14 or 30 days, a table of pipeds, and the newest deployments. |
| Applications | All applications in the project with platform, sync status, last synced time and Git repository. You can search by name or repository, and filter by platform. |
| Application details | Description, Git repository and path, piped, platform provider, labels, and the 20 most recent deployments. This page also has the **Sync** button. |
| Deployment details | Status, commit, who triggered it, duration, and every stage of the pipeline. Open a stage to read its logs. |
| Settings | The control plane address and API key the plugin uses. |

Applications and deployments have an **Open in PipeCD** link to the same item in the web console.

## Prerequisites

- A running PipeCD control plane. See [Installation](../../installation/).
- Headlamp, either the desktop app or Headlamp running in a cluster.
- An API key for your project. A `READ_ONLY` key is enough to view everything. A `READ_WRITE` key is needed to sync.

Some features need a control plane newer than v0.58.0:

- **Quick Sync** and **Pipeline Sync**. An older control plane ignores the choice and runs its default sync.
- The full list of pipeds with their online or offline status on the Overview page. With an older control plane, the Overview only shows pipeds that have an application or a recent deployment, and it cannot show whether they are online.

## Installing the plugin

In Headlamp desktop, open the **Plugin Catalog**, search for **PipeCD**, and click **Install**. Reload Headlamp, and a **PipeCD** item appears in the sidebar.

For Headlamp running in a cluster, follow the Headlamp guide for [installing plugins in a cluster](https://headlamp.dev/docs/latest/installation/in-cluster/).

## Allowing Headlamp to call the control plane

The plugin calls the control plane's API from Headlamp, using gRPC-Web. The control plane already accepts gRPC-Web.

The Headlamp desktop app needs no other setup. When you open Headlamp in a web browser, for example Headlamp running in a cluster, the browser blocks requests to another origin unless the control plane allows them. In that case, turn on CORS in the control plane Helm values, and add the address you open Headlamp at:

```yaml
cors:
  enabled: true
  allowOrigins:
    - "https://headlamp.example.com"
```

Then upgrade the control plane with the new values. The origin must match exactly, including the scheme and port.

If this is not set, the connection test in the plugin fails in the browser even when the address and API key are correct.

## Configuring the plugin

1. Create an API key from the `Settings/API Key` tab on the web UI. Choose the `READ_ONLY` role if you only want to view, or `READ_WRITE` if you also want to sync. See [Authentication](../command-line-tool/#authentication) for more about API keys.

   ![Adding an API key in the PipeCD web console](/images/settings-api-key.png)

2. In Headlamp, go to **PipeCD > Settings**.

3. Enter the control plane address in **PipeCD Server URL**, for example `https://pipecd.example.com`, and the key in **API Key**.

4. Click **Test Connection**. When it succeeds, click **Save Settings**.

The address must start with `https://`, because the API key is sent with every request. `http://` is only accepted for `localhost` and `127.0.0.1`, for local testing.

Only the address is saved. The API key is kept in memory, so you need to enter it again after you reload or restart Headlamp. To remove the key without reloading, clear the **API Key** field and save.

![The plugin settings page after a successful connection test](/images/headlamp-plugin-settings.png)

## Viewing the overview

**PipeCD > Overview** is the first page of the PipeCD section. Use **Time range** to choose the last 24 hours, 7 days, 14 days or 30 days.

- **Pipeds** shows how many pipeds are online and offline. With an older control plane, it shows how many run in the cluster open in Headlamp instead.
- **Applications** shows how many applications are synced and out of sync.
- **Deployments** shows how many deployments ran in the time range, and about how many per day or per hour.
- **Change failure rate** is the share of finished deployments that failed.

The chart shows deployments per day, or per hour for the last 24 hours, colored by status. Only the newest 1,000 deployments are counted, and the page says so when it reaches that limit.

The pipeds table shows each piped's status, version, start time and number of applications. The **In this cluster** column looks for the piped's pod in the Kubernetes cluster open in Headlamp, and shows **Running** or **Not running** with the reason, for example `CrashLoopBackOff`. A piped that runs somewhere else shows `—`.

Click a row in **Recent deployments** to open that deployment.

![The Overview page with the stat cards, the deployments chart and the pipeds table](/images/headlamp-plugin-overview.png)

## Viewing applications

Open **PipeCD > Applications** to see every application in the project. The sync status uses the same colors as the web console: green for `Synced`, red for `Out of Sync`, and yellow for `Deploying` and `Invalid Config`.

![The Applications page with applications on Kubernetes and ECS](/images/headlamp-plugin-applications.png)

Click an application to see its details and recent deployments. In the **Triggered by** column, a deployment started with an API key shows `API key`, and you can hover over it to see the key ID. A deployment started by a Git push shows the commit author.

![An application details page with the Sync button and recent deployments](/images/headlamp-plugin-application-details.png)

## Viewing a deployment

Click a deployment to see its pipeline. Each stage shows its status and how long it ran. Click a stage to open its logs. While a deployment is running, opening a stage again loads its newest logs.

A rollback stage is only shown when it actually ran.

![A deployment with its K8S_SYNC stage open, showing the stage logs](/images/headlamp-plugin-deployment.png)

## Syncing an application

On the application details page, click **Sync** to start a deployment. The arrow next to it lets you choose how to sync:

| Option | Behavior |
|--------|----------|
| Sync | PipeCD decides whether to run a quick sync or the pipeline. |
| Quick Sync | Applies the changes directly, without running the pipeline. |
| Pipeline Sync | Runs the pipeline configured for the application. |

After you confirm, the plugin waits up to a minute for piped to pick up the sync, and then shows whether it started. If PipeCD does not answer in time, the plugin tells you instead of sending the sync again, because the sync may still run. Check the deployments list before you sync again.

Syncing needs a `READ_WRITE` API key. With a `READ_ONLY` key, the plugin shows a message that the key cannot sync, and everything else keeps working. **Sync** is turned off for a disabled application, as in the web console.

![The confirmation dialog shown before a sync](/images/headlamp-plugin-sync.png)

## Troubleshooting

| Problem | What to check |
|---------|---------------|
| Test Connection fails | Check the address and the API key. For Headlamp in a web browser, check that CORS is turned on and the Headlamp address is in `allowOrigins`. See [Allowing Headlamp to call the control plane](#allowing-headlamp-to-call-the-control-plane). |
| The page says PipeCD is not configured | Fill in and save **PipeCD > Settings**. After Headlamp restarts, enter the API key again. |
| Sync says the API key is read-only | Generate a `READ_WRITE` key and save it in the plugin settings. |
| Sync is requested but no deployment starts | Check that the piped for the application is running and connected. |
| Quick Sync or Pipeline Sync runs the default sync | The control plane is v0.58.0 or older. Upgrade it. |
| The Overview shows fewer pipeds than the web console | The control plane is v0.58.0 or older, so the plugin cannot list every piped. Upgrade it. |

## See also

- [Headlamp plugins](https://headlamp.dev/docs/latest/development/plugins/)
- [Plugin source code](https://github.com/headlamp-k8s/plugins/tree/main/pipecd)
