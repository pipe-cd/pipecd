---
title: "ECS Plugin"
linkTitle: "ECS"
weight: 2
description: >
  Specific guide to configuring deployment for Amazon ECS application using the ECS plugin.
---

> **Note:** The ECS plugin is currently in **alpha** status.

The ECS plugin enables PipeCD to deploy applications to Amazon ECS. It manages deployments via ECS task sets using the [`EXTERNAL` deployment controller](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/deployment-type-external.html), which allows PipeCD to control traffic routing and progressive delivery directly.

There are two main deployment modes:

- **Service**: An application that runs continuously, optionally placed behind a load balancer. Requires a `TaskDefinition` file and a `Service` definition file.
- **Standalone task**: A one-time or periodic batch job. Requires only a `TaskDefinition` file.

## Prerequisites

### Piped configuration

Add the ECS plugin to your piped configuration:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  plugins:
    - name: ecs
      port: 7003
      url: file:///path/to/.piped/plugins/ecs
      deployTargets:
        - name: my-ecs-target
          config:
            region: us-east-1
```

See [ECSDeployTargetConfig](#ecsdeploytargetconfig) for all available options including cross-account role assumption and OIDC-based authentication.

### Service definition requirements

The ECS plugin manages deployments via task sets. Your service definition file **must** specify the `EXTERNAL` deployment controller:

```yaml
deploymentController:
  type: EXTERNAL
```

Without this, the plugin cannot create or manage task sets for your service.

## Quick sync

By default, when no `pipeline` is specified in the application configuration, PipeCD triggers a **quick sync** for any merged pull request. Quick sync registers the new task definition, creates/updates the ECS service, promotes a new primary task set, waits for stability, and removes old task sets. All traffic is switched to the new version immediately.

> Standalone tasks always use quick sync. Pipeline sync is not supported for standalone tasks.

### Service with load balancer

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-service
  plugin: ecs
  deployTarget: my-ecs-target
  plugins:
    ecs:
      input:
        serviceDefinitionFile: servicedef.yaml
        taskDefinitionFile: taskdef.yaml
        targetGroups:
          primary:
            targetGroupArn: arn:aws:elasticloadbalancing:ap-northeast-1:XXXX:targetgroup/ecs-lb/YYYY
            containerName: web
            containerPort: 80
```

### Standalone task (batch job)

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-batch-job
  plugin: ecs
  deployTarget: my-ecs-target
  plugins:
    ecs:
      input:
        taskDefinitionFile: taskdef.yaml
        runStandaloneTask: true
        clusterArn: arn:aws:ecs:ap-northeast-1:XXXX:cluster/my-cluster
        launchType: FARGATE
        awsvpcConfiguration:
          assignPublicIp: ENABLED
          subnets:
            - subnet-YYYY
          securityGroups:
            - sg-YYYY
```

### Service discovery (no load balancer)

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-internal-service
  plugin: ecs
  deployTarget: my-ecs-target
  plugins:
    ecs:
      input:
        serviceDefinitionFile: servicedef.yaml
        taskDefinitionFile: taskdef.yaml
        accessType: SERVICE_DISCOVERY
```

## Pipeline sync (progressive delivery)

Add a `pipeline` field to the application configuration to define a custom deployment strategy using the following ECS stages:

| Stage | Description |
|---|---|
| `ECS_PRIMARY_ROLLOUT` | Roll out the new task set as the PRIMARY variant |
| `ECS_CANARY_ROLLOUT` | Roll out the new task set as the CANARY variant (receives no traffic initially) |
| `ECS_TRAFFIC_ROUTING` | Route traffic between PRIMARY and CANARY task sets via ALB listener weights. Requires `accessType: ELB` and both `primary` and `canary` target groups to be configured. |
| `ECS_CANARY_CLEAN` | Clean up all task sets belonging to the CANARY variant |

Common stages also available:

| Stage | Description |
|---|---|
| `WAIT` | Wait for a specified duration |
| `WAIT_APPROVAL` | Wait for a manual approval |
| `ANALYSIS` | Run analysis to validate metrics before proceeding |

`ECS_ROLLBACK` is automatically added when `autoRollback: true` is set (the default).

> **Note:** `ECS_SYNC` is the stage used internally by quick sync. It is not intended for use in a custom pipeline.

### Canary deployment

Gradually shifts a small fraction of traffic to the new version first, observes it, then promotes to primary. Requires **two ALB target groups** (primary and canary) both attached to the same listener.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-service
  plugin: ecs
  deployTarget: my-ecs-target
  plugins:
    ecs:
      input:
        serviceDefinitionFile: servicedef.yaml
        taskDefinitionFile: taskdef.yaml
        targetGroups:
          primary:
            targetGroupArn: arn:aws:elasticloadbalancing:ap-northeast-1:XXXX:targetgroup/ecs-primary/YYYY
            containerName: web
            containerPort: 80
          canary:
            targetGroupArn: arn:aws:elasticloadbalancing:ap-northeast-1:XXXX:targetgroup/ecs-canary/ZZZZ
            containerName: web
            containerPort: 80
  pipeline:
    stages:
      # 1. Deploy the canary task set at 30% of the primary workload capacity.
      #    At this point canary is running but receives no traffic yet.
      - name: ECS_CANARY_ROLLOUT
        with:
          scale: 30

      # 2. Route 10% of traffic to canary, 90% stays on primary.
      - name: ECS_TRAFFIC_ROUTING
        with:
          canary: 10

      # 3. (Optional) Run automated analysis while canary is receiving 10% traffic.
      #    Piped queries the analysis provider every 5 minutes for 10 minutes total.
      #    If the HTTP error rate exceeds 1%, the stage fails and ECS_ROLLBACK triggers.
      #    Remove this stage if you do not have an analysis provider configured.
      - name: ANALYSIS
        with:
          duration: 10m
          metrics:
            - strategy: THRESHOLD
              provider: my-prometheus # name of the analysis provider in piped config
              interval: 5m
              expected:
                max: 0.01 # fail if error rate exceeds 1%
              query: |
                sum(rate(http_requests_total{job="my-service",status=~"5.*"}[5m]))
                /
                sum(rate(http_requests_total{job="my-service"}[5m]))

      # 4. Increase canary traffic to 50%.
      - name: ECS_TRAFFIC_ROUTING
        with:
          canary: 50

      # 5. Wait for manual approval before completing the rollout.
      - name: WAIT_APPROVAL
        with:
          timeout: 1h

      # 6. Promote: update the PRIMARY task set to the new version.
      #    Primary still receives 50% traffic during rollout.
      - name: ECS_PRIMARY_ROLLOUT

      # 7. Route 100% of traffic back to primary (new version).
      - name: ECS_TRAFFIC_ROUTING
        with:
          primary: 100

      # 8. Remove the canary task set.
      - name: ECS_CANARY_CLEAN
```

### Blue/Green deployment

Deploys the new version (green) at full scale alongside the old version (blue), then cuts over all traffic atomically after approval. Requires **two ALB target groups** (primary and canary) both attached to the same listener.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-service
  plugin: ecs
  deployTarget: my-ecs-target
  plugins:
    ecs:
      input:
        serviceDefinitionFile: servicedef.yaml
        taskDefinitionFile: taskdef.yaml
        targetGroups:
          primary:
            targetGroupArn: arn:aws:elasticloadbalancing:ap-northeast-1:XXXX:targetgroup/ecs-blue/YYYY
            containerName: web
            containerPort: 80
          canary:
            targetGroupArn: arn:aws:elasticloadbalancing:ap-northeast-1:XXXX:targetgroup/ecs-green/ZZZZ
            containerName: web
            containerPort: 80
  pipeline:
    stages:
      # 1. Deploy the green task set at 100% capacity alongside the existing blue.
      #    Green receives no traffic yet.
      - name: ECS_CANARY_ROLLOUT
        with:
          scale: 100

      # 2. Wait for approval before cutting over traffic.
      - name: WAIT_APPROVAL
        with:
          timeout: 1h

      # 3. Cut over: send 100% of traffic to green (canary target group).
      #    Blue tasks are still running as a fallback.
      - name: ECS_TRAFFIC_ROUTING
        with:
          canary: 100

      # 4. Update the PRIMARY task set to the new version.
      - name: ECS_PRIMARY_ROLLOUT

      # 5. Switch traffic back to the primary target group (now running the new version).
      - name: ECS_TRAFFIC_ROUTING
        with:
          primary: 100

      # 6. Remove the green (canary) task set.
      - name: ECS_CANARY_CLEAN
```

## Notes

- When using an ALB, all listener rules that reference the configured target groups are controlled by PipeCD. Attach the target groups to your listener rules before deploying.
- [Service Connect](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/service-connect.html) uses the `ECS` deployment controller, which is incompatible with the `EXTERNAL` controller required by this plugin. Canary and blue/green deployments are not available when using Service Connect.
- When using [Service Auto Scaling](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/service-auto-scaling.html), omit `desiredCount` from the service definition file to prevent PipeCD from reconciling it on every deployment.

## Configuration reference

### ECSDeployTargetConfig

Configured under `plugins[].deployTargets[].config` in the piped configuration.

| Field | Type | Description | Required |
|---|---|---|---|
| region | string | AWS region (e.g. `us-east-1`). | Yes |
| profile | string | AWS credentials profile name. Defaults to the `default` profile. | No |
| credentialsFile | string | Path to the AWS shared credentials file. Defaults to `~/.aws/credentials`. | No |
| roleARN | string | IAM role ARN to assume for cross-account access. | No |
| tokenFile | string | Path to the OIDC token file. Required when `roleARN` is set for OIDC-based authentication. | No |

### ECSDeploymentInput

Configured under `plugins.ecs.input` in the application configuration.

| Field | Type | Default | Description |
|---|---|---|---|
| taskDefinitionFile | string | `taskdef.json` | Path to the task definition file (YAML/JSON). |
| serviceDefinitionFile | string | | Path to the service definition file (YAML/JSON). Omit for standalone tasks. |
| runStandaloneTask | bool | `false` | When `true`, runs the task directly without creating or updating an ECS service. Only quick sync is supported. |
| clusterArn | string | | ARN of the ECS cluster. |
| launchType | string | `FARGATE` | Task launch type. Valid values: `EC2`, `FARGATE`. |
| accessType | string | `ELB` | How the service is accessed. Valid values: `ELB`, `SERVICE_DISCOVERY`. |
| awsvpcConfiguration | [ECSVpcConfiguration](#ecsvpcconfiguration) | | VPC configuration. Required for `FARGATE` or `EC2` with `awsvpc` network mode. |
| targetGroups | [ECSTargetGroups](#ecstargetgroups) | | Load balancer target groups. Required when `accessType` is `ELB`. |

### ECSVpcConfiguration

| Field | Type | Description | Required |
|---|---|---|---|
| subnets | []string | VPC subnet IDs where tasks will be launched. Maximum 16. | Yes |
| assignPublicIp | string | Assign a public IP to the task's ENI. Valid values: `ENABLED`, `DISABLED`. | No |
| securityGroups | []string | Security group IDs for the task's ENI. Maximum 5. Uses the VPC default security group if omitted. | No |

### ECSTargetGroups

| Field | Type | Description |
|---|---|---|
| primary | [ECSTargetGroup](#ecstargetgroup) | Target group for the primary variant. |
| canary | [ECSTargetGroup](#ecstargetgroup) | Target group for the canary variant. Required for canary/blue-green deployments. |

### ECSTargetGroup

| Field | Type | Description |
|---|---|---|
| targetGroupArn | string | ARN of the target group. |
| containerName | string | Name of the container to associate with the target group. |
| containerPort | int | Port on the container to associate with the target group. |

### ECSSyncStageOptions (quickSync)

Configured under `plugins.ecs.quickSync`.

| Field | Type | Default | Description |
|---|---|---|---|
| recreate | bool | `false` | When `true`, stops all running tasks before creating the new task set. Guarantees a clean restart but causes brief downtime. |

### ECS_CANARY_ROLLOUT options

| Field | Type | Description |
|---|---|---|
| scale | float | Percentage of the primary workload capacity to deploy as canary (0–100). For example, with 10 primary tasks and `scale: 30`, the canary will have 3 tasks. |

### ECS_TRAFFIC_ROUTING options

Specify either `primary` or `canary`. The other variant automatically receives the remainder.

| Field | Type | Description |
|---|---|---|
| primary | int | Percentage of traffic routed to the primary variant. |
| canary | int | Percentage of traffic routed to the canary variant. |
