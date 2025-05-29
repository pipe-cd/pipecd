---
date: 2025-04-28
title: "Introduction to EventWatcher: Connecting CI to PipeCD"
linkTitle: "Introduction to EventWatcher"
weight: 983
author: Tetsuya Kikuchi ([@t-kikuc](https://github.com/t-kikuc))
categories: ["Introduction"]
tags: ["Feature Guide"]
---

This article explains the basics of "**EventWatcher**", a crucial feature in practical PipeCD usage.

This article is intended for those who:
- "Don't know how to integrate CI with PipeCD"
- "Have found the EventWatcher feature but can't grasp how it works"

_This article is based on PipeCD v0.51.2 (latest at the time of writing)._

## Background: Why EventWatcher is needed

Basically, PipeCD is a CD tool that continuously deploys specified config/manifests.

![](/images/eventwatcher-only-cd.drawio.png)

Deployment typically requires manifest changes. So, how can we "**deploy using a new image (etc.) after CI completion**"? It's annoying to update the manifest repo manually each time.

![](/images/eventwatcher-problem.drawio.png)

This is where the EventWatcher feature comes into play!

## What is EventWatcher?

EventWatcher is a feature in PipeCD that seamlessly connects CI and CD. It updates the manifest repo based on events from CI, triggering CD.

https://pipecd.dev/docs-v0.51.x/user-guide/event-watcher/

### Mechanism

The overall picture of EventWatcher looks like the diagram below.

![](/images/eventwatcher-overview.drawio.png)

The EventWatcher feature itself handles steps 4-6:

- 1-3. Develop new app code and store the results (container images, etc.) in a container registry
- **4.** Publish an event in CI to pass the new image URI to PipeCD
- **5.** Piped detects the event
- **6.** Update the manifest repo using the data provided by the event (image URI, etc.)

- 7-8. Automatic deployment occurs through the standard deployment flow of PipeCD

## Usage

You need to setup three areas:

#### 1. Piped Configuration

1. Which manifest repo events to handle
2. Commit permissions for the manifest repo

https://pipecd.dev/docs-v0.51.x/user-guide/managing-piped/configuring-event-watcher/

#### 2. app.pipecd.yaml Configuration

1. Which events to handle for the Application (`matcher`)
2. What actions to take for matched events (`handler`)

(Example) Configuration to update the `spec.template.spec.containers[0].image` field in `deployment.yaml` when the `helloworld-image-update` event is triggered:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: helloworld
  eventWatcher:
    - matcher: # Which events to handle
        name: helloworld-image-update
      handler: # What actions to take
        type: GIT_UPDATE
        config:
          replacements:
            - file: deployment.yaml
              yamlField: $.spec.template.spec.containers[0].image
```

Full configuration reference:
https://pipecd.dev/docs-v0.51.x/user-guide/configuration-reference/#eventwatcher

#### 3. Event Triggering via pipectl or GitHub Actions

1. Event name (set in `app.pipecd.yaml` as `matcher.name`)
2. New value (image URI, etc.)

(Example) Command to trigger an event named `helloworld-image-update` to change the image URI to `ghcr.io/xxx/helloworld:v0.2.0`:

```sh
pipectl event register \
    --address={CONTROL_PLANE_API_ADDRESS} \
    --api-key={API_KEY} \
    --name=helloworld-image-update \
    --data=ghcr.io/xxx/helloworld:v0.2.0
```

If you use GitHub Actions for your CI, using PipeCD's official `actions-event-register` is recommended (configurations are the same):

https://github.com/marketplace/actions/pipecd-register-event

For example, triggering an event named `helloworld-image-update` to change the image URI to `ghcr.io/xxx/helloworld:v0.2.0`:
```yaml
      - uses: pipe-cd/actions-event-register@v1.2.0
        with:
          api-address: ${{ secrets.API_ADDRESS }}
          api-key: ${{ secrets.API_KEY }}
          event-name: helloworld-image-update
          data: ghcr.io/xxx/helloworld:v0.2.0
```

You can also handle events differently by environment using the `--labels` option, even for events with the same name.

##### Advanced: `--contexts` Option

Adding `--contexts` to `pipectl event register` allows you to include various information in the manifest repo commit. For example, passing the **"commit hash and URL of the app repo that triggered the event"** makes it easier to track "which app repo changes caused this event" in the manifest repo. This is particularly useful when deployments occur frequently.

https://pipecd.dev/docs-v0.51.x/user-guide/event-watcher/#optional-using-contexts

## Appendix: Code Reading

The EventWatcher code is consolidated in one file (though it's quite long). Remembering that "EventWatcher runs continuously as a goroutine from Piped startup" makes it slightly easier to read.

https://github.com/pipe-cd/pipecd/blob/v0.50.2/pkg/app/piped/eventwatcher/eventwatcher.go

## Conclusion

Using EventWatcher enables seamless integration between CI and CD. There are various advanced applications, and it would be great to see more examples being shared.

Additionally, **Deployment Traces** feature enhances the integration of CI and PipeCD. For more details, please see this article:

https://pipecd.dev/blog/2024/12/26/the-gap-between-ci/cd-and-how-pipecd-appeared-to-fulfill-the-gap/
