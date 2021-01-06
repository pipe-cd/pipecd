---
title: "Image watcher"
linkTitle: "Image watcher"
weight: 12
description: >
  Watching container image changes and automatically deploying the new images.
---

Image watcher automatically triggers a new Deployment when a new image tag is pushed to your container registry.

The canonical deployment flow with PipeCD is:

1. CI server pushes the generated image to the container registry after app-repo updated.
1. You update the config-repo manually.

It is the User's responsibility to automate these steps to be done in a series of actions, while it is quite a bit of painful.
Image watcher lets you automate this workflow by continuously performing `git push` to your config-repo.
That is, it frees you from the hassle of manually updating config-repo every time.

## Prerequisites
Before configuring ImageWatcher, all required Image providers must be configured in the Piped Configuration according to [this guide](/docs/operator-manual/piped/configuring-image-watcher/).

## Configuration

Prepare ImageWatcher files placed at the `.pipe/` directory at the root of the Git repository.
In that files, you define what image should be watched and what file should be updated.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: ImageWatcher
spec:
  targets:
    - image: gcr.io/pipecd/helloworld
      provider: my-gcr
      filePath: helloworld/deployment.yaml
      field: $.spec.template.spec.containers[0].image
```

Image watcher periodically compares the latest tag of the following two images:
- a given `image` in a given `provider`
- an image defined at a given `field` in a given `filePath`

And then pushes them to the git repository if there are any deviations.
Note that it uses only pure git push, does not use features that depend on Git hosting services, such as pull requests.


### Examples
Suppose there is the above ImageWatcher file placed at the `.pipe/` directory at the root of the Git repository, and the below kubernetes manifest is placed at `helloworld/deployment.yaml`.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  ...
spec:
  ...
  template:
    spec:
      containers:
      - name: helloworld
        image: gcr.io/pipecd/helloworld:0.1.0
```

Then let's say you pushed the new tag `0.2.0` to the image provider with name `my-gcr`. Piped will create a commit as shown below:

```diff
     spec:
       containers:
       - name: helloworld
-        image: gcr.io/pipecd/helloworld:0.1.0
+        image: gcr.io/pipecd/helloworld:0.2.0
```

See [here](https://github.com/pipe-cd/examples/tree/master/.pipe) for more examples.

The full list of configurable `ImageWatcher` fields are [here](/docs/user-guide/configuration-reference/#image-watcher-configuration).

>ProTip: If multiple Pipeds handle a single repository, you can prevent conflicts by splitting into the multiple files and specifying includes/excludes in the Piped config.
