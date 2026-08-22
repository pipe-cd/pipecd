---
title: "Configuring event watcher"
linkTitle: "Configuring event watcher"
weight: 7
description: >
  This page describes how to configure `piped` to enable event watcher.
---

To enable [Event watcher](../managing-application/event-watcher/), you have to configure your `piped` at first.

## Grant write permission

The [SSH key used by `piped`](../configuration-reference/#git) must be a key with write-access because `piped` needs to commit and push to your git repository when any incoming event matches.

### Specify Git repositories to be observed

`piped` watches events only for the Git repositories specified in the `gitRepos` list.
You need to add all repositories you want to enable Event watcher.

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

### [optional] Specify Event watcher files `piped` will use

>NOTE: This way is valid only for defining events using [.pipe/](../managing-application/event-watcher/#use-the-pipe-directory).

If multiple `piped` instances handle a single repository, you can prevent conflicts by splitting into multiple Event watcher files and setting `includes/excludes` to specify the files that should be monitored by this `piped`.

Say for instance, if you only want the `piped` to use the Event watcher files under `.pipe/dev/`:

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

See the full list of [configurable fields for event watcher](../configuration-reference/#eventwatcher) for more information.

### **OPTIONAL** Settings for git user

By default, every git commit uses `piped` as the username and **pipecd.dev@gmail.com** as the email. You can change it with the [git](../configuration-reference/#git) field.

For example:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  git:
    username: foo
    email: foo@example.com
```
