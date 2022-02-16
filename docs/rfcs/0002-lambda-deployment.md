- Start Date: 2020-12-21
- Target Version: 1.0.0

# Summary

This RFC proposes adding a new service deployment from PipeCD: AWS Lambda deployment. Similar to the current [Google clould run deployment](https://pipecd.dev/docs/feature-status/#cloudrun-deployment).

# Motivation

PipeCD aims to support wide range of deployable services, currently [Terraform deployment](https://pipecd.dev/docs/feature-status/#terraform-deployment) and [Cloud Run deployment](https://pipecd.dev/docs/feature-status/#cloudrun-deployment) are supported. Lambda deployment is the current missing piece for PipeCD's purpose.

# Detailed design

Note:
- Since the `Lambda function as container image` feature was available on AWS recently, we plan to implement PipeCD Lambda deployment based on that update.
- To simplify the implementation, the initial release of this feature will only support container images from ECR ( the AWS container registry ).

### Usage

The deployment configuration is used to customize the way to do the deployment. In case of Lambda function deployment, current common stage options (`WAIT`, `WAIT_APPROVAL`, `ANALYSIS`) are all inherited, besides with the stages for Lambda deployment itself: `LAMBDA_CANARY_ROLLOUT`, `LAMBDA_TRAFFIC_ROUTING`.

```yaml
# https://docs.aws.amazon.com/lambda/latest/dg/configuration-versions.html
apiVersion: pipecd.dev/v1beta1
kind: LambdaApp
spec:
  input:
    image: xxxxxxxxxxxx.dkr.ecr.ap-northeast-1.amazonaws.com/demoapp:v0.0.1
  pipeline:
    stages:
      # Deploy workloads of the new version.
      # But this is still receiving no traffic.
      - name: LAMBDA_CANARY_ROLLOUT
      # Change the traffic routing state where
      # the new version will receive the specified percentage of traffic.
      # This is known as multi-phase canary strategy.
      - name: LAMBDA_TRAFFIC_ROUTING
        with:
          newVersion: 10
      # Optional: We can also add an ANALYSIS stage to verify the new version.
      # If this stage finds any not good metrics of the new version,
      # a rollback process to the previous version will be executed.
      - name: ANALYSIS
      # Change the traffic routing state where
      # thre new version will receive 100% of the traffic.
      - name: LAMBDA_TRAFFIC_ROUTING
        with:
          newVersion: 100
```

### Architecture

Just as current Cloud Run but under `pkg/cloudprovider/lambda` package.

# Alternatives

Lambda function could also be deployed via source code by 2 steps method:
1. Configure piped to be able to clone the source of Lambda function ( on the same repo which deployment be handled by piped or via remote git repo ).
2. Compress the source code and deploy with [SAM](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-command-reference.html) or using [aws-sdk](https://github.com/aws/aws-lambda-go) ( piped handle the deployment ).

The deployment configuration sample as bellow:

In case of source code for Lambda function is on the same repo handled by Piped
```yaml
apiVersion: pipecd.dev/v1beta1
kind: LambdaApp
spec:
  input:
    # Lambda code sourced from the same Git repository.
    path: lambdas/helloworld
```

In case of using difference git repo for lambda function source versioning
```yaml
apiVersion: pipecd.dev/v1beta1
kind: LambdaApp
spec:
  input:
    # Lambda code sourced from another Git repository.
    git: git@github.com:org/source-repo.git
    path: lambdas/helloworld
    ref: v1.0.0
```

# Unresolved questions

Currently, PipeCD team plans to implement this feature using the newly arrived `Lambda function as container images` feature of AWS ( as mentioned on #Detailed Design ). Since this feature is new ( even on AWS ), not all languages are supported to be deployed as container images currently.
