---
title: "Better connect between CI and CD with Deployment Trace"
linkTitle: "Deployment Trace"
weight: 992
description: >
  A helper fulfill the gap between CI and CD.
---

You are a developer who works with application code change, and don't know what deployment is triggered by your commit on PipeCD UI? This feature is for you.

If you're using PipeCD [Event Watcher](./event-watcher) to trigger the deployment for your code change, you can attach information of the triggered commit as the event data, PipeCD will use that information and helps you to make a link between your application code commit and the triggered deployments that reflect your code change.

![](/images/deployment-trace-ui.png)

## Usage

Via `pipectl event register` command

```bash
Usage:
  pipectl event register [flags]

Flags:
      --commit-author string      The author of commit that triggers the event.
      --commit-hash string        The commit hash that triggers the event.
      --commit-message string     The message of commit that triggers the event.
      --commit-timestamp int      The timestamp of commit that triggers the event.
      --commit-title string       The title of commit that triggers the event.
      --commit-url string         The URL of commit that triggers the event.
```

Note: You have to attach at least `commit-hash` and `commit-url` as the event data in order to use the Deployment Trace feature.

## Github Actions
If you're using Github Actions in your CI workflow, [actions-event-register](https://github.com/marketplace/actions/pipecd-register-event) is for you!
With it, you can easily register events without any installation.

