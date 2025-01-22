---
date: 2024-12-26
title: "The Gap Between CI/CD and How PipeCD Appeared to FulFill the Gap"
linkTitle: "The Gap Between CI/CD and How PipeCD Appeared to FulFill the Gap"
weight: 984
author: Khanh Tran ([@khanhtc1202](https://github.com/khanhtc1202))
categories: ["Announcement"]
tags: ["New Feature", "Feature Improvement", "Developer Experience"]
---

PipeCD became a CNCF Sandbox last year and has only just begun to attract the attention of the community, but it has a long history of development since its inception, more than four years ago.

In this article, I would like to share the fundamental CICD challenges that the PipeCD development and operations team at CyberAgent identified, how we solved them in the current version of PipeCD, and our future direction.

**Note:** The content may be change in the future since it's still under development.

## GitOps and CICD

Although there are many different definitions and understandings, there are some commonalities when it comes to related concepts such as GitOps and CICD:

- CI/CD separation is necessary
- CI belongs to the application source code scope, and the output of CI is the input of CD. CD belongs to the manifest source code scope, and the input of CD is the output of CI
- CD tools perform tasks based on GitOps principles


## EventWatcher as the bridge between CI and CD

Since the first version, which didn't have much user feedback, the development team has been pushing testing and using dogfooding techniques as a way to detect shortcomings in PipeCD's functionality.
Since these early days, the problem of effectively linking CI and CD has emerged.

![](/images/cicd-flow.png)

source: https://www.weave.works/blog/what-is-gitops-really

A common way to bypass the "immutability firewall" between CI and CD without paying too much attention to what CI is being used for is to rely on an artifact (container image) that is the common element between CI and CD. This approach is used by many popular CDs today, such as ArgoCD and FluxCD.

![](/images/image-watcher-flow.png)

In addition to advantages such as a generic interface (as it only depends on the container image itself), this approach also has limitations, such as:

- It needs to support many different container registries (Docker Hub, ECR, GCR, ACR, Harbor, etc.)
- It can only support container images as the artifacts, and not other artifacts such as Helm Charts, Terraform modules, Lambda source zips, etc.
- As the number of container images to be monitored increases, performance-related issues arise 

Because PipeCD's most important feature is to support various platforms, being able to only support container images as artifacts is not suitable for PipeCD. Therefore, we decided to develop a new solution to solve this problem: EventWatcher.

![](/images/pipecd-event-watcher-flow.png)

Instead of choosing a container image as the interface between CI and CD, PipeCD directly uses events sent from CI to PipeCD as triggers for the CD pipeline.

An event is considered a command to PipeCD (similar to triggering a deployment from the UI). It contains information such as:

- The event definition name and event label, are used to determine which application deployments should be triggered when piped processes this event
- Information about what data is used by piped when processing the event
- Information specifying paths to fields in the manifest file

For example, suppose you want to use PipeCD EventWatcher feature to auto trigger deployment for your Kubernetes application. You first need to update your application configuration manifest like this

```yaml

apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: helloworld
  eventWatcher:
    - matcher:
        name: event-test
      handler:
        type: GIT_UPDATE
        config:
          replacements:
            - file: deployment.yaml
              yamlField: $.spec.template.spec.containers[0].image
```

Then, in the last step of your CI pipeline, you need to use pipectl to send an event to PipeCD controlplane like this

```bash
$ pipectl event register --name=event-test --data=ghcr.io/pipecd/helloworld:v0.49.0
```

The piped which manage your Kubernetes application will be notified by the event, it will create a commit to your application manifest repo, which contains

```diff
  spec:
    template:
      spec:
        containers:
        - name: helloworld
-         image: ghcr.io/pipecd/helloworld:v0.48.0
+         image: ghcr.io/pipecd/helloworld:v0.49.0
```

And the piped will trigger a new deployment based on the changed commit.

EventWatcher helps you avoid performance issues as the number of container images you monitor grows. It supports a wide variety of artifact types (Helm charts, Terraform modules, Lambda source zips, etc.) while easily ensuring compatibility with a wide range of CIs.

## The issue with EventWatcher

EventWatcher solves the above problems, but it also has certain limitations, although they are not directly due to the design of EventWatcher.

Due to security requirements, projects usually do not store source code and manifests in the same repository, but separate them into two repositories.

![](/images/code-manifest-repo-separation.png)

A common request is that when users look at a pipeline deployment on PipeCD, how can they know what source code (code repo) changes are being deployed in this deployment?

If you use container images as an interface between CI and CD, you can take advantage of the LABEL feature of containers to store information about source code commit hashes, branches, etc. The OCI specification also has some predefined LABELs that you can use without defining additional keys.

```dockerfile
LABEL org.opencontainers.image.source="https://github.com/pipe-cd/pipecd"
LABEL org.opencontainers.image.revision="0123456789abcdefg"
```

However, although this approach is similar to the previous one, it is only suitable for artifacts that are container images, and at the same time, scanning container images to get information about commit hashes and displaying it in the UI incurs additional processing overhead.

At this point, we realized a strength in the EventWatcher design: sending events from CI to CD.

The essence of this problem is that when viewing a deployment, you need to know information in the code repository that is only relevant to the manifest repository. In the approach of sending events between CI and CD, you add the necessary information, such as the commit hash and PR URL, as part of the event. When it is piped to create a commit change to the manifest repository, this information is used to link the deployment with the commit trigger and commit the code change.

We use the [git trailer](https://git-scm.com/docs/git-interpret-trailers) feature to add this information to the manifest change commit. And finally, wrap up it as `--context` flag of the pipectl command.

```bash
pipectl event register --name=event-test --data=ghcr.io/pipecd/helloworld:v0.49.0 \
--context=Source-Commit-Hash= 6512b1f3fcba7342c11,Source-URL=https://github.com/pipe-cd/pipecd/commit/6512b1f3fcba7342c11,PipeCD-Official-Docs=https://pipecd.dev/
```

If you're using GitHub as your manifest repo remote storage, you can see it on GitHub UI like this

![](/images/pipecd-commit-context-github-ui.png)

With this feature released under [PipeCD v0.49.3](https://github.com/pipe-cd/pipecd/releases/tag/v0.49.3), users now can have a link from CD to CI for tracking back which application code change causes this deployment in CD pipeline.

## It's good, but it's not enough

After the release of the EventWatcher update, the community provided positive feedback about resolving the connection issues between CI and CD.
However, we notice some limitations when using PipeCD in very complex projects with many microservices in the same repository.

Each event sent to PipeCD to trigger a deployment can actually be shared to trigger multiple deployments. This is obvious, for example, changes related to protos or common shared libraries, such as logging and metrics libraries, lead to multiple deployments at once.

![](/images/eventwatcher-trigger-multi-deployments.png)

This becomes problematic when a development team includes many people working on the same repository and the number of changes affecting many services is large enough: PipeCD's current UI deployment list makes it difficult to determine exactly which deployments were triggered by which code change.

![](/images/problematic-deployment-list-ui.jpg)

Naively, we could just add a `root_cause_commit_hash` or `codebase_commit_hash` field to PipeCD’s deployment model and use this information to filter deployments triggered by the same repository code change in the PipeCD UI.

However, if users need to know more than just commit hash, this small change doesn't reliably achieve the goal. What users really want is a link from the code repository to the manifest repository, that the devs (who know the code change), can use to track which deployments shift those changes to the clusters.

So we looked back at the design of PipeCD and EventWatcher and took the thinking a step further.

We decided to add the concept of Deployment Trace to PipeCD, using the code repository commit hash as the root key. For each code repository commit hash, a Deployment Trace contains the following information:

- commit_hash: The commit hash of the code repository is used as the key
- code repository information: Includes code change author, URL, branches, commit messages, authors, etc.
- deployments: A list of deployments triggered by the key commit_hash. Each deployment contains the complete information of the PipeCD deployment (commit, path, pipeline in manifest repository, etc.)

This allows users to easily follow the information from the code repository -> manifest repository -> deployment without searching for information step by step. From the perspective of an application engineer making changes to the code repository, it is easy to see which deployment was triggered by a commit, what that deployment was, and how that deployment impacts the service.

We are summarizing [this GitHub issue](https://github.com/pipe-cd/pipecd/issues/5444) on the design and implementation of deployment tracing. Feel free to add any idea on the issue.


## Conclusion

In this article, we shared the issues we encountered when using PipeCD in a highly complex project, how we solved them through the design and implementation of EventWatcher, and our next steps with deployment tracing.

The deployment tracing feature is planned and will be implemented as soon as possible. It is expected to be released around March 2025.
The PipeCD project is an open-source project; anyone can participate in the development of the project, and we welcome any contributions from the community.

If you have any questions or opinions, feel free to talk in [CNCF Slack](https://cloud-native.slack.com/) > #pipecd or [PipeCD Community Meeting](https://docs.google.com/document/d/1AtE0CQYbUV5wLfvAcl9mo9MyTCH52BuU7AngVUvE7vg/edit).
