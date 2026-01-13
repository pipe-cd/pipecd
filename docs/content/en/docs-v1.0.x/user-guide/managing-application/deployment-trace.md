---
title: "Deployment Trace"
linkTitle: "Deployment Trace"
weight: 992
description: >
  Understand the deployment trace feature and how you can link application commits to the deployments they trigger.
---

Deployment Trace links application code changes to the resulting deployments that reflect those code changes in the PipeCD Web UI.

When using PipeCD [Event Watcher](./event-watcher) to trigger deployments, you can attach commit information as event data. PipeCD uses that information to create links between your application code commits and the triggered deployments.

![Deployment Trace feature](/images/deployment-trace-ui.png)

## Usage

Use the `pipectl event register` command:

```bash
  pipectl event register [flags]
```

Flags:

```
      --commit-author string      The author of commit that triggers the event.
      --commit-hash string        The commit hash that triggers the event.
      --commit-message string     The commit message that triggers the event.
      --commit-timestamp int      The timestamp of commit that triggers the event.
      --commit-title string       The title of commit that triggers the event.
      --commit-url string         The URL of commit that triggers the event.
```

>**Note:** Attach at least `commit-hash` and `commit-url` as event data to use the Deployment Trace feature.

## GitHub Actions

If you are using GitHub Actions in your CI workflow, use [actions-event-register](https://github.com/marketplace/actions/pipecd-register-event) to register events without any installation.
