---
title: "Lambda"
linkTitle: "Lambda"
weight: 4
description: >
  Specific guide for configuring Lambda deployment.
---

Deploying a Lambda application requires a `function.yaml` file placing inside the application directory. That file contains values to be used to deploy Lambda function on your AWS cluster.
Currently, only container image built source is supported by piped deployment. For more information about container images as function, read [this post on AWS blog](https://aws.amazon.com/blogs/aws/new-for-aws-lambda-container-image-support/).

A sample `function.yaml` file as following:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: LambdaFunction
spec:
  name: SimpleFunction
  role: arn:aws:iam::76xxxxxxx:role/lambda-role
  image: ecr.ap-northeast-1.amazonaws.com/lambda-test:v0.0.1
  tags:
    app: simple
```

Except the `tags` field, all others are required fields for the deployment to run.

## Quick sync

By default, when the [pipeline](/docs/user-guide/configuration-reference/#lambda-application) was not specified, PipeCD triggers a quick sync deployment for the merged pull request.
Quick sync for a Lambda deployment will roll out the new version and switch all traffic to it.

## Sync with the specified pipeline

The [pipeline](/docs/user-guide/configuration-reference/#lambda-application) field in the deployment configuration is used to customize the way to do the deployment.
You can add a manual approval before routing traffic to the new version or add an analysis stage the do some smoke tests against the new version before allowing them to receive the real traffic.

These are the provided stages for Lambda application you can use to build your pipeline:

- `LAMBDA_CANARY_ROLLOUT`
  - deploy workloads of the new version, but it is still receiving no traffic.
- `LAMBDA_PROMOTE`
  - promote the new version to receive an amount of traffic.

and other common stages:
- `WAIT`
- `WAIT_APPROVAL`
- `ANALYSIS`

See the description of each stage at [Configuration Reference](/docs/user-guide/configuration-reference/#stageoptions).

Here is an example that rolls out the new version gradually:

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: LambdaApp
spec:
  pipeline:
    stages:
      # Deploy workloads of the new version.
      # But this is still receiving no traffic.
      - name: LAMBDA_CANARY_ROLLOUT
      # Promote new version to receive 10% of traffic.
      - name: LAMBDA_PROMOTE
        with:
          percent: 10
      - name: WAIT
        with:
          duration: 10m
      # Promote new version to receive 50% of traffic.
      - name: CLOUDRUN_PROMOTE
        with:
          percent: 50
      - name: WAIT
        with:
          duration: 10m
      # Promote new version to receive all traffic.
      - name: CLOUDRUN_PROMOTE
        with:
          percent: 100
```

## Reference

See [Configuration Reference](/docs/user-guide/configuration-reference/#lambda-application) for the full configuration.
