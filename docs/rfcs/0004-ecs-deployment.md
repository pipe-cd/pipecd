- Start Date: (fill me in with today's date, YYYY-MM-DD)
- Target Version: (1.x / 2.x)

# Summary

This RFC proposes adding a new service deployment from PipeCD: AWS ECS deployment.

# Motivation

Why are we doing this? What do we expect?

# Detailed design

Note:
- Since the `Lambda function as container image` feature was available on AWS recently, we plan to implement PipeCD Lambda deployment based on that update.
- To simplify the implementation, the initial release of this feature will only support container images from ECR ( the AWS container registry ).

## Usage

The deployment configuration is used to customize the way to do the deployment. In the case of AWS ECS deployment, current common stages (`WAIT`, `WAIT_APPROVAL`, `ANALYSIS`) are all inherited, besides with the stages for ECS deployment `ECS_SYNC`.

```yaml
# https://docs.aws.amazon.com/lambda/latest/dg/configuration-versions.html
apiVersion: pipecd.dev/v1beta1
kind: ECSApp
spec:
  input:
    name: Sample
    taskDefinition: path/to/taskdef.json
    image: xxxxxxxxxxxx.dkr.ecr.ap-northeast-1.amazonaws.com/demoapp:v0.0.1
```

## Architecture

Just as current Lambda but under `pkg/cloudprovider/ecs` package.

# Alternatives

What other designs have been considered? What is the impact of not doing this?

# Unresolved questions

What parts of the design are still TBD?
