---
title: "Configuring event watcher"
linkTitle: "Configuring event watcher"
weight: 6
description: >
  This page describes how to configure piped to enable event watcher.
---

To enable [EventWatcher](/docs/user-guide/event-watcher/), you have to configure your piped at first.

### [require] Grant write permission
The [SSH key used by Piped](/docs/operator-manual/piped/configuration-reference/#git) must be a key with write-access because Event watcher automates the deployment flow by commitng and pushing to your git repository.

### [optional] Settings for watcher
The Piped's behavior can be finely controlled by setting the `eventWatcher` field.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  eventWatcher:
    checkInterval: 5m
    gitRepos:
      - repoId: repo-1
        commitMessage: Update values by Event watcher
        includes:
          - event-watcher-dev.yaml
          - event-watcher-stg.yaml
```

If multiple Pipeds handle a single repository, you can prevent conflicts by splitting into the multiple EventWatcher files and setting `includes/excludes` to specify the files that should be monitored by this Piped.
`excludes` is prioritized if both `includes` and `excludes` are given.

The full list of configurable fields are [here](/docs/operator-manual/piped/configuration-reference/#eventwatcher).

### [optional] Settings for git user
By default, every git commit uses `piped` as a username and `pipecd.dev@gmail.com` as an email. You can change it with the [git](/docs/operator-manual/piped/configuration-reference/#git) field.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  git:
    username: foo
    email: foo@example.com
```
