---
date: 2022-08-22
title: "PipeCD best practice 02 - PipeCD as PipeCD's CD"
linkTitle: "PipeCD best practice 02"
weight: 992
description: "This blog is a part of PipeCD best practice series, a guideline for you to operate your own PipeCD cluster."
author: Khanh Tran ([@khanhtc1202](https://twitter.com/khanhtc1202))
---

As we mentioned in the [May 2022 update](/blog/2022/06/01/may-2022-update/#were-opening-more-and-more) blog post, the PipeCD dev team moves our PipeCD's CI flow from our internal tools to use the Github Actions as a step to make the development of PipeCD more open as its OSS vision. So, how about the CD for PipeCD?

Good news, we're developing a powerful CD platform as well - __the PipeCD__! In this blog, I will show you a real life usecase of using PipeCD, which is using PipeCD as itself CD platform. Hope you can get some ideas for your own usecase of using PipeCD.

### The idea

> __Eating your own dog food__ or "__dogfooding__" is the practice of using one's own products or services. This can be a way for an organization to test its products in real-world usage using product management techniques (via [wiki](https://en.wikipedia.org/wiki/Eating_your_own_dog_food)).


It's a good idea to use your own product (at least in the context of software production), which helps the developers have an early look at any not-good points and problems of their product. Even though you can have test, dev, staging, and so on separated environments to test your own product in some ways that are well-defined, using your own product still brings some new looks and ideas to improve and keep the product alive.

PipeCD team adopts that idea, we design our CD flow for PipeCD as below.

![image](/images/pipecd-cd-architecture.png)
<p style="text-align: center;">
PipeCD's CD architecture
</p>

### The PipeCD's CD architecture

We are operating 5 clusters in total for different purposes, I will walk you through all of them.

##### 1. The production cluster and its PipeCD controlplane

This cluster host the PipeCD controlplane for our internal company use. The PipeCD dev team operates this cluster and other teams from our company only need to install Piped and use that for their own applications deployments.

PipeCD uses `project` as a logical context to separate applications of different teams and different purposes. You can host only one PipeCD controlplane for company-wide use like us with this controlplane.

Based on that idea, aside with projects of other company teams, we also use this controlplane to manage our PipeCD’s applications' deployments under a PipeCD project named `platform`.

##### 2. The dog cluster

In PipeCD model, we need a Piped which actually runs the deployment and a controlplane to manage that Piped. Basically, we can not use a Piped to self-deploy its managing controlplane, that is why this dog cluster and its running on controlplane exists.

In the dog cluster, we host a PipeCD controlplane, which manages exactly only one Piped, that deploys the Production controlplane we talked about above.

> If you ask, which Piped deploys this dog controlplane? It's the point interesting thing comes. The Piped which is registered to the Production controlplane is the one that deploys this dog cluster controlplane. Production and dog deploy each other!

##### 3. The dev, internal play and external play clusters

All three other clusters are used for the development of PipeCD which are:
- The dev cluster: run the dev use controlplane and piped, which used by PipeCD dev team.
- The internal play cluster: run the internal play controlplane and piped. Basically, it's a demonstrate environment to show the PipeCD users some early access versions which not yet been published.
- The external play cluster: run the external play controlplane and piped, we published it as [play.pipecd.dev](https://play.pipecd.dev).

These clusters: dog, dev, internal play, external play, and their running controlplane deployment are all deployed by a Piped which is managed by the Production controlplane under `project: platform` project.

#### Some notes you can learn from PipeCD's self CD architecture

- It's best practice to store the applications' configuration of different projects in separated repositories. But in our case, for simplicity, we stored our all PipeCD applications in the same place. We can do it with confidence because the applications' configurations in case of using PipeCD are all encrypted or be a part of the Piped configuration, not in the applications' configuration itself. In case of different teams using the same controlplane, those teams should have their own applications' configuration repo to reduce the miss.
- As mentioned in [PipeCD best practice 01 blog](/blog/2021/12/29/pipecd-best-practice-01-operate-your-own-pipecd-cluster/#:~:text=we%20highly%20recommend%20running%20each%20Piped%20inside%20each%20cluster%20and%20that%20Piped%20will%20only%20manage%20the%20applications%20on%20that%20cluster), we recommend running each Piped inside each cluster and that Piped should only manage the applications on that cluster to reduce the risk of credentials creating/sharing, you still can actually use one single Piped to deploy whatever you want. In case of PipeCD self CD architecture, the Piped of `project: platform` actually does not run in the applications cluster, the reason is at the time that Piped deploys Terraform to build dev and play cluster, those clusters have not existed yet. You can follow the PipeCD pattern but be sure that you do it on purpose.


That's all ;)
