---
title: "Configuring image watcher"
linkTitle: "Configuring image watcher"
weight: 6
description: >
  This page describes how to configure piped to enable image watcher.
---

To enable [ImageWatcher](/docs/user-guide/image-watcher/), you have to configure your piped at first.

## Prerequisites
The [SSH key Piped use](/docs/operator-manual/piped/configuration-reference/#git) must be a key with write-access because Image watcher automates the deployment flow by commitng and pushing to your git repository.

## Adding an image provider
Define arbitrary number of [image providers](/docs/concepts#image-provider) which is information needed to connect from your Piped to the container registry.
It will run a pull operation every 5 minutes by default. This interval can be set in the `imageWatcher` field touch upon later.
Also, we plan to provide a FAKE image provider mentioned below to avoid the rate limit.

Currently, PipeCD is supporting:
- [Google Container Registry (GCR)](https://cloud.google.com/container-registry)
- [Amazon Elastic Container Registry (ECR)](https://aws.amazon.com/ecr)

### GCR
Append the `GCR` image provider to the Piped configuration file as:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  imageProviders:
    - name: my-gcr
      type: GCR
      config:
        serviceAccountFile: /etc/piped-secret/gcr-service-account.json
```

For public repositories, no configuration is required.

If you want to watch private repository, you should set up authentication.
A [service account](https://cloud.google.com/compute/docs/access/service-accounts) is the only authentication way currently available.
You give the path to the json file of service account with the required `roles/storage.objectViewer` role.

The full list of GCR fields are [here](/docs/operator-manual/piped/configuration-reference/#imageprovidergcrconfig).

### ECR

>NOTE: Currently, it supports only ECR private repositories.

Append the `ECR` image provider to the Piped configuration file as:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  imageProviders:
    - name: my-ecr
      type: ECR
      config:
        region: ap-northeast-1
        credentialsFile: /etc/piped-secret/aws-credentials
        profile: user1
```

The only required field is `region`.

You will generally need your AWS credentials to authenticate with ECR. Piped provides multiple methods of loading these credentials.
It attempts to retrieve credentials in the following order:
1. From the environment variables. Available environment variables are `AWS_ACCESS_KEY_ID` or `AWS_ACCESS_KEY` and `AWS_SECRET_ACCESS_KEY` or `AWS_SECRET_KEY`
1. From the given credentials file.
1. From the EC2 Instance Role

Hence, you don't have to set `credentialsFile` if you use the environment variables or the EC2 Instance Role. Keep in mind the IAM role/user that you use with your Piped must possess the IAM policy permission for `ecr:DescribeImages`.

The full list of ECR fields are [here](/docs/operator-manual/piped/configuration-reference/#imageproviderecrconfig).

### DockerHub

>TBA

### FAKE

>TBA: We plan to provide a FAKE image provider to deal with the rate limit from the container registry.
>
>The FAKE container registry is deployed at the control-plane, and you can store the metadata about newly updated images whenever you want to (e.g. on your CI).


## [optional] Settings for watcher
The Piped's behavior can be finely controlled by setting the `imageWatcher` field.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  imageWatcher:
    checkInterval: 5m
    gitRepos:
      - repoId: foo
        commitMessage: Update image
        includes:
          - imagewatcher-dev.yaml
          - imagewatcher-stg.yaml
```

If multiple Pipeds handle a single repository, you can prevent conflicts by splitting into the multiple ImageWatcher files and setting `includes/excludes` to specify the files that should be monitored by this Piped.
`excludes` is prioritized if both `includes` and `excludes` are given.

The full list of configurable fields are [here](/docs/operator-manual/piped/configuration-reference/#imagewatcher).

## [optional] Settings for git user
By default, every git commit uses `piped` as a username and `pipecd.dev@gmail.com` as an email. You can change it with the [git](/docs/operator-manual/piped/configuration-reference/#git) field.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  git:
    username: foo
    email: foo@example.com
```
