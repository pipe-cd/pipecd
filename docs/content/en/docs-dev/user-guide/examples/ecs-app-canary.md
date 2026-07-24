---
title: "Canary deployment for ECS app"
linkTitle: "Canary ECS app"
weight: 5
description: >
  How to enable canary deployment for ECS application.
---

> Canary release is a technique to reduce the risk of introducing a new software version in production by slowly rolling out the change to a small subset of users before rolling it out to the entire infrastructure and making it available to everybody.
> -- [martinfowler.com/canaryrelease](https://martinfowler.com/bliki/CanaryRelease.html)

PipeCD enables the canary strategy for ECS applications by using ECS' [external deployment controller](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/deployment-type-external.html) together with TaskSets. Instead of a single rolling Service update, PipeCD manages separate TaskSets for the current version and the new version, and shifts traffic between them by controlling the weights of the target groups attached to your load balancer's listener rules.

In this guide, we will show you how to configure the application configuration file to roll out a `CANARY` TaskSet, send a portion of traffic to it, and then, after a hold period, promote the new version to `PRIMARY` and clean up the intermediate resources.

Complete source code for this example is hosted in [pipe-cd/examples](https://github.com/pipe-cd/examples/tree/master/ecs/canary) repository.

## Before you begin

- Add a new ECS application by following the instructions in [this guide](../../managing-application/adding-an-application/)
- Ensure your [Service definition](https://github.com/pipe-cd/examples/blob/master/ecs/canary/servicedef.yaml) uses the `EXTERNAL` deployment controller, which is required for PipeCD to manage TaskSets directly

``` yaml
deploymentController:
  type: EXTERNAL
```

- Ensure you have a `primary` and a `canary` target group already linked to your load balancer's listener rules. PipeCD only controls the traffic weight of listener rules that already reference these target groups; it does not create or link them for you.

``` yaml
targetGroups:
  primary:
    targetGroupArn: arn:aws:elasticloadbalancing:ap-northeast-1:XXXX:targetgroup/ecs-canary-blue/YYYY
    containerName: web
    containerPort: 80
  canary:
    targetGroupArn: arn:aws:elasticloadbalancing:ap-northeast-1:XXXX:targetgroup/ecs-canary-green/ZZZZ
    containerName: web
    containerPort: 80
```

## Enabling canary strategy

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
        targetGroupArn: arn:aws:elasticloadbalancing:ap-northeast-1:XXXX:targetgroup/ecs-canary-blue/YYYY
        containerName: web
        containerPort: 80
      canary:
        targetGroupArn: arn:aws:elasticloadbalancing:ap-northeast-1:XXXX:targetgroup/ecs-canary-green/ZZZZ
        containerName: web
        containerPort: 80
  pipeline:
    stages:
      - name: ECS_CANARY_ROLLOUT
        with:
          scale: 30
      - name: ECS_TRAFFIC_ROUTING
        with:
          canary: 20
      - name: WAIT
        with:
          duration: 150s
      - name: ECS_PRIMARY_ROLLOUT
      - name: ECS_TRAFFIC_ROUTING
        with:
          primary: 100
      - name: ECS_CANARY_CLEAN
```

- Send a PR to update the task definition (e.g. the container image) and merge it to trigger a new deployment. PipeCD will plan the deployment with the specified canary strategy.

- Now you have an automated canary deployment for your ECS application!

## Understanding what happened

Throughout this pipeline, PipeCD works with up to three ECS TaskSets at once: the **current** TaskSet that was already serving traffic, a temporary **canary** TaskSet, and a **new** TaskSet that eventually replaces the current one.

- Stage 1: `ECS_CANARY_ROLLOUT` creates the canary TaskSet running the new version of the task definition. Its task count is scaled to 30% of the current TaskSet's count. It is registered to the canary target group, but at this time it still handles no traffic.

- Stage 2: `ECS_TRAFFIC_ROUTING` updates the listener rule weights so that 20% of traffic is routed to the canary target group, while the remaining 80% keeps going to the primary target group (still served by the current TaskSet). This lets you observe how the new version behaves under real traffic before committing to it fully.

- Stage 3: `WAIT` holds the deployment for the configured duration, giving you time to check metrics, logs, or dashboards for the canary version. You can replace this with a `WAIT_APPROVAL` stage for a manual gate, or add an `ANALYSIS` stage to automate the check.

- Stage 4: `ECS_PRIMARY_ROLLOUT` creates the new TaskSet, running the new version and registered to the same primary target group as the current TaskSet. At this point, the current and new TaskSets briefly share that target group.

- Stage 5: `ECS_TRAFFIC_ROUTING` updates the listener rule again, routing 100% of traffic back to the primary target group and away from the canary target group.

- Stage 6: `ECS_CANARY_CLEAN` promotes the new TaskSet to `PRIMARY` status and deletes both the old current TaskSet and the canary TaskSet, returning the service to the same shape it started in — just running the new version.

If any stage in this pipeline fails, PipeCD automatically rolls back the deployment: traffic is routed back to the original current TaskSet, and any TaskSets created during the deployment are removed.