---
date: 2022-06-01
title: "May 2022 update"
linkTitle: "May 2022 update"
weight: 993
description: "Development status update to recap what happened last few months"
author: Khanh Tran ([@khanhtc1202](https://twitter.com/khanhtc1202))
categories: ["Announcement"]
tags: ["News"]
---

_Published by the PipeCD dev team every few months, this update will provide you with news and updates about the project! Please click [here](/blog/2022/02/10/february-2022-update/) if you want to see the last status update._

### What's changed
---

It's been a while from the last update about the development of our PipeCD project, but you can be assured that our PipeCD development is going well. Since the last report, PipeCD team has introduced 5 major releases (from [v0.27.0](https://github.com/pipe-cd/pipecd/releases/tag/v0.27.0) to the latest [v0.32.3](https://github.com/pipe-cd/pipecd/releases/tag/v0.32.3)) with several UI improvements and some interesting features that I can't wait to show you now. This blog post walks you through some notable changes. For all other changes, please check out each release note.

#### Cloud Run live state

Until now, you may notice that only application of kind Kubernetes has a part in its application detail page named live-state which shows you the situation of components of the current running PipeCD application. From version [v0.27.0](https://github.com/pipe-cd/pipecd/releases/tag/v0.27.0), we made that feature available for applications of kind Cloud Run too.

![](/images/cloudrun-live-state.png)

#### Make analysis stage skippable

One of the outstanding features PipeCD has is [Automated deployment analysis (ADA)](/docs/user-guide/automated-deployment-analysis/), which helps you evaluate the impact of the current deployment right middle its running. Though it’s necessary to have that kind of stage for the stability of your service, sometimes you may want to skip the stage for some kind of quick/hot fixes. In that case, this feature is for you.

![](/images/analysis-skippable.png)

#### Enable to load Helm chart from OCI registry

With the release of Helm [version v3.8.0](https://helm.sh/blog/storing-charts-in-oci/), Helm is able to store and work with charts in container registries, as an alternative to Helm repositories. PipeCD adopts that feature right a way and for now, you're able to fetch Helm chart from wherever OCI registries and use it to deploy your services. You can find more about how to use this feature in detail in this [docs](/docs/operator-manual/piped/adding-helm-chart-repository-or-registry/#adding-helm-chart-registry).

#### The very first touchable FileDB

It’s been a while since the first time we made a [tweet](https://twitter.com/nghialv2607/status/1480712569535209472) about how PipeCD aims to remove the database as its dependencies to make the installation easier and compact, today we want to make the first achievement on the road to __#nodatabase__ world.\
For now, you can have a quick shot of using the new FileDB “database” of PipeCD with the quickstart example easily on your local machine by changing the `quickstart/control-plane-values.yaml` file contains __config.data__ section to something like below:

```yaml
    apiVersion: "pipecd.dev/v1beta1"
    kind: ControlPlane
    spec:
      datastore:
        type: FILEDB
      filestore:
        type: MINIO
        config:
          endpoint: http://pipecd-minio:9000
          bucket: quickstart
          accessKeyFile: /etc/pipecd-secret/minio-access-key
          secretKeyFile: /etc/pipecd-secret/minio-secret-key
          autoCreateBucket: true
      projects:
        - id: quickstart
          staticAdmin:
            username: hello-pipecd
            passwordHash: "$2a$10$ye96mUqUqTnjUqgwQJbJzel/LJibRhUnmzyypACkvrTSnQpVFZ7qK" # bcrypt value of "hello-pipecd"
```

Please follow the [development guideline](/docs/contribution-guidelines/development/#how-to-run-control-plane-locally) or [quickstart guideline](/docs/quickstart/) if you don't know how to start/install PipeCD control plane on your local machine.

#### The PipeCD play is available

For all users who want to have a glance at how and what PipeCD can give you in real-life usecases, this [PipeCD play environment](https://play.pipecd.dev) is for you. We have a blog about that live demo, feel free to check it from [here](/blog/2022/04/12/the-pipecd-play-environment-is-here/).

And many other UI improvements are updated to the PipeCD web console, let's check it out on our [play environment](https://play.pipecd.dev).

### We're opening more and more

PipeCD is an OSS project! And we want to make not just the source code open, but the whole deployment process should be opened as well in order to get more contribution and support from the open community. We moved almost all of our CI flow to the open Github Actions platform, and as the result, everyone who makes contributions to our PipeCD repository can easily follow/access the CI stages via Github Actions UI. At the time I wrote this blog, PipeCD CI flow had been opened and can be accessed via this [link](https://github.com/pipe-cd/pipecd/actions).

### What are next
---

The team continues actively working on improving the PipeCD product. Besides fixing the reported issues, enhancing the existing features, here are some new features you may catch for few next releases:

- Piped management via web console: restartable and configuration checkable.
- RBAC for PipeCD resources based on the Label mechanism

If you have any features want to request or find out a problem, please let us know by creating issues to the [pipe-cd/pipecd](https://github.com/pipe-cd/pipecd/issues) repository.

---
*Follow us on Twitter to keep track of all the latest news: https://twitter.com/pipecd_dev*

Happy PipeCD-ing 👋
