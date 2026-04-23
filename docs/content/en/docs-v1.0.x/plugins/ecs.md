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

### IAM permissions

ECS deployments involve three distinct IAM roles. Confusing them is a common source of permission errors.

**Deployment role**: the role assumed by Piped when calling AWS APIs to deploy your application. Configured via `roleARN` in the deploy target config, or resolved from the default credential chain. This role needs:

```json
{
  "Effect": "Allow",
  "Action": [
    "ecs:CreateService",
    "ecs:UpdateService",
    "ecs:DescribeServices",
    "ecs:RegisterTaskDefinition",
    "ecs:DescribeTaskDefinition",
    "ecs:CreateTaskSet",
    "ecs:UpdateServicePrimaryTaskSet",
    "ecs:DeleteTaskSet",
    "ecs:DescribeTaskSets",
    "ecs:ListTasks",
    "ecs:DescribeTasks",
    "ecs:RunTasks",
    "ecs:ListTagsForResource",
    "ecs:TagResource",
    "ecs:UntagResource",
    "elasticloadbalancing:DescribeListeners",
    "elasticloadbalancing:DescribeRules",
    "elasticloadbalancing:ModifyRule",
    "iam:PassRole"
  ],
  "Resource": "*"
}
```

`iam:PassRole` is required because AWS enforces it whenever a caller passes a role ARN to ECS APIs. Specifically, it is triggered by `RegisterTaskDefinition` (which receives `executionRoleArn` and `taskRoleArn` from the task definition file) and `CreateService` (which receives the service role ARN). In practice, scope this permission to the specific role ARNs used in your task and service definitions rather than `"Resource": "*"`.

**[Task execution role](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_execution_IAM_role.html)**: assumed by the ECS agent on your behalf to prepare the task environment before the container starts. It is not used by Piped or your application. Typical uses:

- Pull container images from Amazon ECR
- Publish container logs to CloudWatch Logs
- Retrieve secrets from Secrets Manager or SSM Parameter Store

This role is specified as `executionRoleArn` in your task definition file.

**[Task role](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-iam-roles.html)**: assumed by the application running inside the container. It is not used by Piped or the ECS agent. Grant only the permissions your application needs at runtime, for example reading from S3 or writing to DynamoDB.

This role is specified as `taskRoleArn` in your task definition file.

## Definition files

The ECS plugin reads two YAML or JSON files from your application directory to determine what to deploy.

### Task definition file

Specifies the container(s) to run: image, CPU, memory, ports, environment variables, and IAM roles. The file is parsed into the AWS SDK `types.TaskDefinition` struct (the response type), then the plugin forwards a subset of those fields to `RegisterTaskDefinition`. Field names must match the Go struct field names (case-insensitive), which correspond to the AWS API JSON field names in camelCase. Fields present in `RegisterTaskDefinitionInput` but absent from `types.TaskDefinition` cannot be set through this file.

```yaml
# taskdef.yaml
family: my-service
executionRoleArn: arn:aws:iam::XXXX:role/ECSTaskExecutionRole
taskRoleArn: arn:aws:iam::XXXX:role/ECSTaskRole # optional: grants your app AWS permissions
requiresCompatibilities:
  - FARGATE
networkMode: awsvpc
cpu: "256"
memory: "512"
containerDefinitions:
  - name: web
    image: public.ecr.aws/nginx/nginx:1.27
    portMappings:
      - containerPort: 80
    cpu: 100
    memory: 100
```

The default filename is `taskdef.json`. Override it with `taskDefinitionFile` in the application config.

> **Note:** The `tags` field in the task definition file is not currently forwarded to `RegisterTaskDefinition`. Support is planned for a future version.

### Service definition file

Specifies how the ECS service runs: cluster, desired task count, network configuration, and load balancer settings. The file is parsed into the AWS SDK `types.Service` struct (the response type), then the plugin forwards a subset of those fields to `CreateService` or `UpdateService`. Field names must match the Go struct field names (case-insensitive), which correspond to the AWS API JSON field names in camelCase. Fields present in `CreateServiceInput`/`UpdateServiceInput` but absent from `types.Service` cannot be set through this file.

```yaml
# servicedef.yaml
serviceName: my-service
clusterArn: arn:aws:ecs:ap-northeast-1:XXXX:cluster/my-cluster
desiredCount: 2
deploymentController:
  type: EXTERNAL # required by the ECS plugin
launchType: FARGATE
networkConfiguration:
  awsvpcConfiguration:
    assignPublicIp: ENABLED
    subnets:
      - subnet-YYYY
    securityGroups:
      - sg-YYYY
schedulingStrategy: REPLICA
```

> **Note:** Due to current parsing limitations, not all fields are forwarded to `CreateService` and `UpdateService`. The fields listed below are not yet supported:
>
> - `CreateService`: `capacityProviderStrategy`, `serviceConnectConfiguration`
> - `UpdateService`: `deploymentConfiguration`, `healthCheckGracePeriodSeconds`, `placementConstraints`, `platformVersion`, `capacityProviderStrategy`, `serviceConnectConfiguration`, `serviceRegistries`, `volumeConfigurations`
>
> Full field support is planned alongside support for the ECS deployment controller in a future version.

**Fields you should not include:**

- `loadBalancers`: PipeCD injects the target group configuration from `targetGroups` in the app config at deploy time.
- `desiredCount`: omit this when using Auto Scaling, otherwise PipeCD will reconcile it back to the value in the file on every deployment.
- Tags managed by PipeCD (`pipecd/managed-by`, `pipecd/commit-hash`, etc.): these are stamped automatically and will be overwritten.

The service definition file is not needed for standalone tasks.

### Service definition requirements

The ECS plugin manages deployments via task sets. Your service definition file **must** specify the `EXTERNAL` deployment controller:

```yaml
deploymentController:
  type: EXTERNAL
```

Without this, the plugin cannot create or manage task sets for your service.

> **Note:** Support for the ECS deployment controller (`type: ECS`) is planned for a future version.


## Choosing a deployment strategy

When a `pipeline` is configured, the ECS plugin compares the container images in the target task definition against those in the currently running deployment to decide whether to execute the pipeline or skip it with a quick sync:

| Condition | Strategy |
|---|---|
| No running deployment source available | Pipeline sync |
| Running task definition cannot be loaded | Pipeline sync (fallback) |
| One or more container images added, removed, or changed | Pipeline sync |
| No container image changes | Quick sync |

The comparison checks each container by name. If the image name differs, the full change is reported; if only the tag or digest differs, the version change is reported. All detected changes appear in the deployment summary in the UI.

> Non-image changes (environment variables, CPU/memory, etc.) do not affect strategy selection. If you update only those fields, the plugin will choose quick sync even when a pipeline is configured.

**Use quick sync when:**
- You want a simple, immediate rollout with no traffic splitting.
- The environment is non-critical (e.g. dev, staging) where gradual rollout adds no value.

> Standalone tasks always use quick sync. Pipeline sync is not supported for standalone tasks.

**Use pipeline sync when:**
- You want to gradually shift traffic using canary or blue-green strategies.
- You need manual approval gates before full rollout.
- You want automated analysis (metrics, logs) to validate the new version before promoting it.

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

## Migrating from v0

This section covers the steps specific to ECS applications. For the general migration process (updating `pipectl`, migrating the Control Plane database, installing pipedv1), follow the [Migrate to PipeCD v1](../../migrating-from-v0-to-v1/) guide first.

### 1. Update the piped configuration

In v0, ECS was configured as a platform provider with `type: ECS`. In v1, it is a plugin.

**v0 piped config:**

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  platformProviders:
    - name: ecs-dev
      type: ECS
      config:
        region: us-east-1
        credentialsFile: ~/.aws/credentials
        profile: default
        roleARN: arn:aws:iam::XXXX:role/deployment-role
        tokenFile: /var/run/secrets/token
```

**v1 piped config:**

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  plugins:
    - name: ecs
      port: 7003
      url: https://github.com/pipe-cd/pipecd/releases/download/<version>/ecs-plugin
      deployTargets:
        - name: ecs-dev
          config:
            region: us-east-1
            credentialsFile: ~/.aws/credentials
            profile: default
            roleARN: arn:aws:iam::XXXX:role/deployment-role
            tokenFile: /var/run/secrets/token
```

The credential fields (`region`, `credentialsFile`, `profile`, `roleARN`, `tokenFile`) are identical. They move from `platformProviders[].config` to `plugins[].deployTargets[].config`.

### 2. Update the application configuration

The application config format changes in two ways: the `kind` field and the location of ECS-specific fields.

**v0 app config:**

```yaml
apiVersion: pipecd.dev/v1beta1
kind: ECSApp
spec:
  name: my-service
  labels:
    env: production
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
  quickSync:
    recreate: false
  pipeline:
    stages:
      - name: ECS_CANARY_ROLLOUT
        with:
          scale: 30
      - name: ECS_TRAFFIC_ROUTING
        with:
          canary: 20
      - name: ECS_PRIMARY_ROLLOUT
      - name: ECS_TRAFFIC_ROUTING
        with:
          primary: 100
      - name: ECS_CANARY_CLEAN
```

**v1 app config:**

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-service
  labels:
    env: production
  plugin: ecs
  deployTarget: ecs-dev
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
      quickSync:
        recreate: false
  pipeline:
    stages:
      - name: ECS_CANARY_ROLLOUT
        with:
          scale: 30
      - name: ECS_TRAFFIC_ROUTING
        with:
          canary: 20
      - name: ECS_PRIMARY_ROLLOUT
      - name: ECS_TRAFFIC_ROUTING
        with:
          primary: 100
      - name: ECS_CANARY_CLEAN
```

Summary of field changes:

| v0 | v1 |
|---|---|
| `kind: ECSApp` | `kind: Application` |
| (no field) | `spec.plugin: ecs` |
| (no field) | `spec.deployTarget: <deploy-target-name>` |
| `spec.input.*` | `spec.plugins.ecs.input.*` |
| `spec.quickSync.*` | `spec.plugins.ecs.quickSync.*` |
| `spec.pipeline.*` | `spec.pipeline.*` (unchanged) |
| `spec.name`, `spec.labels` | unchanged |

### 3. Review behavioral changes

The v1 plugin includes several fixes to behaviors that could cause service disruptions in v0. Read the [Changes from v0](#changes-from-v0) section to understand what changed before running your first v1 deployment.

## Changes from v0

This section describes behavioral differences between the legacy ECS provider in PipeCD v0 and the ECS plugin in PipeCD v1. If you are migrating from v0, review these changes before deploying.

### ECS_PRIMARY_ROLLOUT: only the old PRIMARY task set is removed

**v0 behavior:** `ECS_PRIMARY_ROLLOUT` deleted all ACTIVE task sets (including the CANARY task set) before promoting the new PRIMARY. This created a window where no task set was serving traffic, causing HTTP 503 errors on the load balancer: [issue #4710](https://github.com/pipe-cd/pipecd/issues/4710).

**v1 behavior:** `ECS_PRIMARY_ROLLOUT` records the current PRIMARY task set ARN *before* creating the new one, promotes the new task set to PRIMARY, then deletes only the old PRIMARY. Any CANARY task set created by an earlier `ECS_CANARY_ROLLOUT` stage remains intact and continues to serve its share of traffic throughout the rollout window. The CANARY is cleaned up separately by `ECS_CANARY_CLEAN`.

This change eliminates the 503 window that existed in v0 during canary deployments.

### Rollback: ELB weights are restored before task sets are modified

**v0 behavior:** During rollback, task sets were recreated first. This left a window where the ALB listener was still sending a fraction of traffic to the canary target group even though its tasks were being deleted, resulting in 503 errors for requests hitting that target group.

**v1 behavior:** The `ECS_ROLLBACK` stage restores ALB listener weights to `100% primary / 0% canary` *before* touching any task set. Listener ARNs and the canary target group ARN are persisted in deployment metadata by `ECS_TRAFFIC_ROUTING`, so rollback can look them up without an additional AWS API call. Only after the listener is safe does rollback create the new task set, promote it to PRIMARY, and delete the remaining task sets.

### Drift detection: commit hash tag instead of config field comparison

**v0 behavior:** Drift detection compared the fields of the live ECS resources against what was declared in the Git definition files. This produced false OUT_OF_SYNC signals whenever AWS mutated a field outside of PipeCD's control, for example when Auto Scaling adjusted `desiredCount`, or when AWS updated internal service metadata.

**v1 behavior:** The plugin stamps a `pipecd/commit-hash` tag onto each PRIMARY task set at the end of every deployment. Drift detection compares this tag value against the current Git commit hash:

- If they match: **SYNCED**
- If they differ: **OUT_OF_SYNC**, with the deployed and expected hashes shown as the reason

AWS-side mutations that do not change the deployed commit (such as Auto Scaling adjustments) no longer trigger drift alerts.

## Application live state

The ECS plugin continuously monitors your application's live state and displays it on the PipeCD UI. Piped periodically polls AWS and compares the running state against the commit hash declared in Git.

### Sync status

| Status | Meaning |
|---|---|
| **SYNCED** | The PRIMARY task set on AWS is tagged with the current Git commit hash. |
| **OUT\_OF\_SYNC** | The commit hash tag on the PRIMARY task set does not match the current Git commit, the service was not found, or no PRIMARY task set exists. |
| **UNKNOWN** | The application is a standalone task (no persistent service to observe), or live state could not be fetched due to a transient error. |
| **INVALID\_CONFIG** | `serviceDefinitionFile` is missing from the application configuration. |

### Resource tree

The UI displays the ECS resources as a hierarchy. The structure depends on the deployment controller type configured in the service definition.

**EXTERNAL controller** (required by this plugin for canary/blue-green):

```
ECS:Service
└── ECS:TaskSet  (one per variant: PRIMARY, CANARY, etc.)
    └── ECS:Task
```

### Resource metadata

Each resource type exposes the following metadata fields in the UI:

**ECS:Service**

| Field | Description |
|---|---|
| status | Service status (`ACTIVE`, `DRAINING`, `INACTIVE`) |
| runningCount | Number of tasks currently running |
| desiredCount | Number of tasks the service intends to maintain |
| pendingCount | Number of tasks in the process of starting |

**ECS:TaskSet** (EXTERNAL controller only)

| Field | Description |
|---|---|
| status | Task set status (`PRIMARY`, `ACTIVE`, `DRAINING`) |
| taskDefinition | ARN of the task definition revision running in this task set |
| scale | Traffic weight assigned to this task set (e.g. `30%`) |
| commit | Git commit hash that created this task set |

**ECS:Task**

| Field | Description |
|---|---|
| lastStatus | Last known task status (`RUNNING`, `PENDING`, `STOPPED`) |
| startedAt | Timestamp when the task transitioned to RUNNING |


### Health status

| Resource | HEALTHY condition | UNHEALTHY condition |
|---|---|---|
| Service | Status is `ACTIVE`, `pendingCount` is 0, and `runningCount ≥ desiredCount` | Status is `DRAINING`/`INACTIVE`, or tasks are pending |
| TaskSet | Stability status is `STEADY_STATE` | Stability status is `STABILIZING` |
| Deployment | Rollout state is `COMPLETED` | Rollout state is `FAILED` or `IN_PROGRESS` |
| Task | Last status is `RUNNING` and container health checks pass | Last status is `PENDING`/`STOPPED`, or container health checks are failing |

## Plan preview

Before a deployment is triggered, PipeCD can show a preview of what will change. The ECS plugin compares the **running** definition files (from the last deployed commit) against the **target** definition files (from the incoming commit) and displays a unified diff.

### What is shown

| Item | Description |
|---|---|
| Task definition diff | Unified diff between the running `taskdef` file and the target `taskdef` file |
| Service definition diff | Unified diff between the running `servicedef` file and the target `servicedef` file (only shown when `serviceDefinitionFile` is set) |
| Summary | A one-line summary: `task definition changed`, `service definition changed`, or `No changes were detected` |

### Example output

When only the container image in the task definition is updated, the plan preview looks like:

```diff
--- taskdef (running)
+++ taskdef (target)
@@ -3,7 +3,7 @@
 containerDefinitions:
   - name: web
-    image: public.ecr.aws/nginx/nginx:1.25
+    image: public.ecr.aws/nginx/nginx:1.27
     cpu: 100
     memory: 100
     portMappings:
```

If no files changed between the two commits, the summary shows `No changes were detected` and no diff is displayed.

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
