---
title: "Configuring event watcher"
linkTitle: "Configuring event watcher"
weight: 7
description: >
  This page describes how to configure piped to enable event watcher.
---

To enable [EventWatcher](../../event-watcher/), you have to configure your piped at first.

### Grant write permission
The [SSH key used by Piped](../configuration-reference/#git) must be a key with write-access because piped needs to commit and push to your git repository when any incoming event matches.

### Specify Git repositories to be observed
Piped watches events only for the Git repositories specified in the `gitRepos` list.
You need to add all repositories you want to enable Eventwatcher.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  eventWatcher:
    gitRepos:
      - repoId: repo-1
      - repoId: repo-2
      - repoId: repo-3
```

### [optional] Specify Eventwatcher files Piped will use
>NOTE: This way is valid only for defining events using [.pipe/](../../event-watcher/#use-the-pipe-directory).

If multiple Pipeds handle a single repository, you can prevent conflicts by splitting into the multiple EventWatcher files and setting `includes/excludes` to specify the files that should be monitored by this Piped.

Say for instance, if you only want the Piped to use the Eventwatcher files under `.pipe/dev/`:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  eventWatcher:
    gitRepos:
      - repoId: repo-1
        commitMessage: Update values by Event watcher
        includes:
          - dev/*.yaml
```

`excludes` is prioritized if both `includes` and `excludes` are given.

The full list of configurable fields are [here](../configuration-reference/#eventwatcher).

### [optional] Settings for git user
By default, every git commit uses `piped` as a username and `pipecd.dev@gmail.com` as an email. You can change it with the [git](../configuration-reference/#git) field.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  git:
    username: foo
    email: foo@example.com
```
