---
title: "Configuring ECS application"
linkTitle: "ECS"
weight: 5
description: >
  Specific guide to configuring deployment for Amazon ECS application.
---

Deploying an Amazon ECS application requires `TaskDefinition` and `Service` configuration files placing inside the application directory. Those files contain all configuration for [ECS TaskDefinition](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definitions.html) object and [ECS Service](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs_services.html) object, and will be used by Piped agent while deploy your application/service to ECS cluster.

If you're not familiar with ECS, you can get examples for those files from [here](../../../../examples/#ecs-applications).

## Quick sync

By default, when the [pipeline](../../../configuration-reference/#ecs-application) was not specified, PipeCD triggers a quick sync deployment for the merged pull request.
Quick sync for an ECS deployment will roll out the new version and switch all traffic to it immediately.

## Sync with the specified pipeline

The [pipeline](../../../configuration-reference/#ecs-application) field in the application configuration is used to customize the way to do the deployment.
You can add a manual approval before routing traffic to the new version or add an analysis stage the do some smoke tests against the new version before allowing them to receive the real traffic.

These are the provided stages for ECS application you can use to build your pipeline:

- `ECS_CANARY_ROLLOUT`
  - deploy workloads of the new version as CANARY variant, but it is still receiving no traffic.
- `ECS_PRIMARY_ROLLOUT`
  - deploy workloads of the new version as PRIMARY variant, but it is still receiving no traffic.
- `ECS_TRAFFIC_ROUTING`
  - routing traffic to the specified variants.
- `ECS_CANARY_CLEAN`
  - destroy all workloads of CANARY variant.

and other common stages:
- `WAIT`
- `WAIT_APPROVAL`
- `ANALYSIS`

See the description of each stage at [Customize application deployment](../../customizing-deployment/).

Here is an example that rolls out the new version gradually:

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: ECSApp
spec:
  input:
    # Path to Service configuration file in Yaml/JSON format.
    # Default is `service.json`
    serviceDefinitionFile: servicedef.yaml
    # Path to TaskDefinition configuration file in Yaml/JSON format.
    # Default is `taskdef.json`
    taskDefinitionFile: taskdef.yaml
    targetGroups:
      primary:
        targetGroupArn: arn:aws:elasticloadbalancing:ap-northeast-1:XXXX:targetgroup/ecs-canary-blue/YYYY
        containerName: web
        containerPort: 80
      canary:
        targetGroupArn: arn:aws:elasticloadbalancing:ap-northeast-1:XXXX:targetgroup/ecs-canary-green/YYYY
        containerName: web
        containerPort: 80
  pipeline:
    stages:
      # Deploy the workloads of CANARY variant, the number of workload
      # for CANARY variant is equal to 30% of PRIMARY's workload.
      # But this is still receiving no traffic.
      - name: ECS_CANARY_ROLLOUT
        with:
          scale: 30
      # Change the traffic routing state where
      # the CANARY workloads will receive the specified percentage of traffic.
      # This is known as multi-phase canary strategy.
      - name: ECS_TRAFFIC_ROUTING
        with:
          canary: 20
      # Optional: We can also add an ANALYSIS stage to verify the new version.
      # If this stage finds any not good metrics of the new version,
      # a rollback process to the previous version will be executed.
      - name: ANALYSIS
      # Update the workload of PRIMARY variant to the new version.
      - name: ECS_PRIMARY_ROLLOUT
      # Change the traffic routing state where
      # the PRIMARY workloads will receive 100% of the traffic.
      - name: ECS_TRAFFIC_ROUTING
        with:
          primary: 100
      # Destroy all workloads of CANARY variant.
      - name: ECS_CANARY_CLEAN
```

## Reference

See [Configuration Reference](../../../configuration-reference/#ecs-application) for the full configuration.
