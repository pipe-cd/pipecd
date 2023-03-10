---
title: "Configuring Lambda application"
linkTitle: "Lambda"
weight: 4
description: >
  Specific guide to configuring deployment for Lambda application.
---

Deploying a Lambda application requires a `function.yaml` file placing inside the application directory. That file contains values to be used to deploy Lambda function on your AWS cluster.
Currently, Piped supports deploying all types of Lambda deployment packages:
- container images (called [container image as Lambda function](https://aws.amazon.com/blogs/aws/new-for-aws-lambda-container-image-support/))
- `.zip` file archives (which stored in AWS S3)

Besides, Piped also supports deploying your Lambda function __directly from the function source code__ which is stored in a remote git repository.

#### Deploy container image as Lambda function

A sample `function.yaml` file for container image as Lambda function used deployment as follows:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: LambdaFunction
spec:
  name: SimpleFunction
  image: ecr.ap-northeast-1.amazonaws.com/lambda-test:v0.0.1
  role: arn:aws:iam::76xxxxxxx:role/lambda-role
  # The amount of memory available to the Lambda application 
  # at runtime. The value can be any multiple of 1 MB.
  memory: 512
  # Timeout of the Lambda application, the value must
  # in between 1 to 900 seconds.
  timeout: 30
  tags:
    app: simple
  environments:
    FOO: bar
```

Except the `tags` and the `environments` field, all others are required fields for the deployment to run.

The `role` value represents the service role (for your Lambda function to run), not for Piped agent to deploy your Lambda application. To be able to pull container images from AWS ECR, besides policies to run as usual, you need to add `Lambda.ElasticContainerRegistry` __read__ permission to your Lambda function service role.

The `environments` field represents environment variables that can be accessed by your Lambda application at runtime. __In case of no value set for this field, all environment variables for the deploying Lambda application will be revoked__, so make sure you set all currently required environment variables of your running Lambda application on `function.yaml` if you migrate your app to PipeCD deployment.

#### Deploy .zip file archives as Lambda function

It's recommended to use container image as Lambda function due to its simplicity, but as mentioned above, below is a sample `function.yaml` file for Lambda which uses zip packing source code stored in AWS S3.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: LambdaFunction
spec:
  name: SimpleZipPackingS3Function
  role: arn:aws:iam::76xxxxxxx:role/lambda-role
  # --- 3 next lines allow Piped to determine your Lambda function code stored in AWS S3.
  s3Bucket: pipecd-sample-lambda
  s3Key: pipecd-sample-src
  s3ObjectVersion: 1pTK9_v0Kd7I8Sk4n6abzCL
  # ---
  handler: app.lambdaHandler
  runtime: nodejs14.x
  memory: 512
  timeout: 30
  environments:
    FOO: bar
  tags:
    app: simple-zip-s3
```

Value for the `runtime` field should be listed in [AWS Lambda runtimes official docs](https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html). All other fields setting are remained as in the case of using [container image as Lambda function](#deploy-container-image-as-lambda-function) pattern.

#### Deploy source code directly as Lambda function

In case you don’t have a separated CI pipeline that provides artifacts (such as container image, built zip files) as its outputs and want to set up a simple pipeline to deploy the Lambda function directly from your source code, this deployment package is for you.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: LambdaFunction
spec:
  name: SimpleCanaryZipFunction
  role: arn:aws:iam::76xxxxxxx:role/lambda-role
  # source configuration use to determine the source code of your Lambda function.
  source:
    # git remote address where the source code is placing.
    git: git@github.com:username/lambda-function-code.git
    # the commit SHA or tag for remote git. Use branch name means you will always use
    # the latest code of that branch as Lambda function code which is NOT recommended.
    ref: dede7cdea5bbd3fdbcc4674bfcd2b2f9e0579603
    # relative path from the repository root directory to the function code directory.
    path: hello-world
  handler: app.lambdaHandler
  runtime: nodejs14.x
  memory: 128
  timeout: 5
  tags:
    app: canary-zip
```

All other fields setting are remained as in the case of using [.zip archives as Lambda function](#deploy-zip-file-archives-as-lambda-function) pattern.

## Quick sync

By default, when the [pipeline](../../../configuration-reference/#lambda-application) was not specified, PipeCD triggers a quick sync deployment for the merged pull request.
Quick sync for a Lambda deployment will roll out the new version and switch all traffic to it.

## Sync with the specified pipeline

The [pipeline](../../../configuration-reference/#lambda-application) field in the application configuration is used to customize the way to do the deployment.
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

See the description of each stage at [Customize application deployment](../../customizing-deployment/).

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
      - name: LAMBDA_PROMOTE
        with:
          percent: 50
      - name: WAIT
        with:
          duration: 10m
      # Promote new version to receive all traffic.
      - name: LAMBDA_PROMOTE
        with:
          percent: 100
```

## Reference

See [Configuration Reference](../../../configuration-reference/#lambda-application) for the full configuration.
