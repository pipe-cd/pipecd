---
title: "BlueGreen deployment for ECS app"
linkTitle: "BlueGreen ECS app"
weight: 6
description: >
  How to enable blue-green deployment for ECS application.
---

Similar to [canary deployment](../ecs-app-canary/), PipeCD allows you to enable and automate the blue-green deployment strategy for your ECS application, also based on ECS's external deployment controller and TaskSets.

In both canary and blue-green strategies, the old version and the new version of the application get deployed at the same time. But while the canary strategy slowly routes traffic to the new version, the blue-green strategy quickly routes all traffic to one version or the other.

In this guide, we will show you how to configure the application configuration file to apply the blue-green strategy.

Complete source code for this example is hosted in [pipe-cd/examples](https://github.com/pipe-cd/examples/tree/master/ecs/bluegreen) repository.

## Before you begin

- Add a new ECS application by following the instructions in [this guide](../../managing-application/adding-an-application/)
- Ensure your [Service definition](https://github.com/pipe-cd/examples/blob/master/ecs/bluegreen/servicedef.yaml) uses the `EXTERNAL` deployment controller, which is required for PipeCD to manage TaskSets directly

``` yaml
deploymentController:
  type: EXTERNAL
```

- Ensure you have a `primary` and a `canary` target group already linked to your load balancer's listener rules. PipeCD only controls the traffic weight of listener rules that already reference these target groups; it does not create or link them for you.

``` yaml
targetGroups:
  primary:
    targetGroupArn: arn:aws:elasticloadbalancing:ap-northeast-1:XXXX:targetgroup/ecs-tg-blue/YYYY
    containerName: web
    containerPort: 80
  canary:
    targetGroupArn: arn:aws:elasticloadbalancing:ap-northeast-1:XXXX:targetgroup/ecs-tg-green/ZZZZ
    containerName: web
    containerPort: 80
```

## Enabling blue-green strategy

- Add the following application configuration file into the application directory in Git.

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: ECSApp
spec:
  input:
    serviceDefinitionFile: servicedef.yaml
    taskDefinitionFile: taskdef.yaml
    targetGroups:
      primary:
        targetGroupArn: arn:aws:elasticloadbalancing:ap-northeast-1:XXXX:targetgroup/ecs-tg-blue/YYYY
        containerName: web
        containerPort: 80
      canary:
        targetGroupArn: arn:aws:elasticloadbalancing:ap-northeast-1:XXXX:targetgroup/ecs-tg-green/ZZZZ
        containerName: web
        containerPort: 80
  pipeline:
    stages:
      - name: ECS_CANARY_ROLLOUT
        with:
          scale: 100
      - name: ECS_TRAFFIC_ROUTING
        with:
          canary: 100
      - name: WAIT
        with:
          duration: 150s
      - name: ECS_PRIMARY_ROLLOUT
      - name: ECS_TRAFFIC_ROUTING
        with:
          primary: 100
      - name: ECS_CANARY_CLEAN
```

- Send a PR to update the task definition (e.g. the container image) and merge it to trigger a new deployment. PipeCD will plan the deployment with the specified blue-green strategy.

- Now you have an automated blue-green deployment for your ECS application!

## Understanding what happened

In this example, you configured the application to switch all traffic from the old version to the new version at once, instead of gradually shifting it.

- Stage 1: `ECS_CANARY_ROLLOUT` creates a new TaskSet running the new version, scaled to 100% of the current TaskSet's count — matching capacity 1:1 rather than a small fraction. It is registered to the canary target group, but at this time it still handles no traffic.

- Stage 2: `ECS_TRAFFIC_ROUTING` updates the listener rule weights so that 100% of traffic is routed to the canary target group at once. Traffic is now fully served by the new version, while the old version keeps running unused in the background in case a rollback is needed.

- Stage 3: `WAIT` holds the deployment for the configured duration, giving you time to verify the new version under full production traffic. You can replace this with a `WAIT_APPROVAL` stage for a manual gate, or add an `ANALYSIS` stage to automate the check.

- Stage 4: `ECS_PRIMARY_ROLLOUT` creates a new TaskSet running the new version, registered to the same primary target group as the original current TaskSet.

- Stage 5: `ECS_TRAFFIC_ROUTING` routes 100% of traffic back to the primary target group.

- Stage 6: `ECS_CANARY_CLEAN` promotes the new TaskSet to `PRIMARY` status and deletes both the old current TaskSet and the canary TaskSet, returning the service to the same shape it started in — just running the new version.

If any stage in this pipeline fails, PipeCD automatically rolls back the deployment: traffic is routed back to the original current TaskSet, and any TaskSets created during the deployment are removed.
