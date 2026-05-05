# AWS ECS plugin

## Overview

ECS plugin supports the Deployment for AWS ECS.

> [!CAUTION]
> Currently, this is alpha status.

### Quick sync

Quick sync rolls out the new version and switches all traffic to it immediately.

It will be planned in one of the following cases:
- no pipeline was specified in the application configuration file
- the application is a standalone task (pipeline sync is not supported for standalone tasks)

For example, the application configuration below is missing the pipeline field. This means any pull request that touches the application will trigger a quick sync deployment.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
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

### Pipeline sync

You can configure the pipeline to enable a progressive deployment with a strategy like canary or blue-green.

ECS plugin defines two variants for each application: primary (aka stable) and canary.
- `primary` runs the current version of code and configuration.
- `canary` runs the proposed change of code or configuration.

Once the deployment is completed, only the `primary` variant should remain.

These are the provided stages for ECS plugin you can use to build your pipeline:

- `ECS_SYNC`: sync ECS service with given task definition (used for quick sync)
- `ECS_PRIMARY_ROLLOUT`: roll out the new task set as PRIMARY variant
- `ECS_CANARY_ROLLOUT`: roll out the new task set as CANARY variant (serves no traffic initially when using ELB)
- `ECS_CANARY_CLEAN`: clean up task sets of CANARY variant
- `ECS_TRAFFIC_ROUTING`: route traffic between PRIMARY and CANARY task sets
- `ECS_ROLLBACK`: rollback to the previous task set (automatically added when `autoRollback` is enabled)

## Directory Structure

```
ecs/
├── main.go         # Plugin entrypoint
├── config/         # Configuration types for piped, deploy targets, application, and stages
├── deployment/     # Stage execution and pipeline planning
├── provider/       # AWS ECS API client and resource operations
├── livestate/      # Live state fetcher (reports current state of ECS resources)
└── planpreview/    # Plan preview (shows diff before deployment)
```

## How to Build & Run

### Build

Build the ECS plugin binary using the `build/plugin` Makefile target from the repository root:

```bash
# Build all plugins (including ecs)
make build/plugin

# Build only the ECS plugin
make build/plugin PLUGINS=ecs
```

The binary will be placed at `~/.piped/plugins/ecs` by default.

### Run

Configure your piped to load the ECS plugin by adding it to the piped config:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  plugins:
  - name: ecs
    port: 7003 # any unused port
    url: file:///home/<user>/.piped/plugins/ecs
    deployTargets:
      - name: production
        config:
          region: us-east-1
```

Then start piped as usual. The plugin process will be launched automatically by piped on the specified port.

> [!NOTE]
> If you build the plugin manually with `go build` and place the binary at a custom path, piped may use a previously cached binary instead. Pass `--force-plugin-redownload` when starting piped to ensure it picks up your locally built binary:
> ```bash
> piped --force-plugin-redownload ...
> ```

## How to Test

### Unit tests

Unit tests use mocks for all AWS API calls and do not require real AWS credentials. Run them from the repository root:

```bash
go test ./pkg/app/pipedv1/plugin/ecs/...
```

### Integration tests

To test against a real AWS environment, build and run piped with the plugin pointing to a real ECS cluster. The plugin resolves credentials in the following order:

1. `credentialsFile` + `profile` in deploy target config: explicit credentials file and profile
2. `tokenFile` + `roleARN` in deploy target config: OIDC web identity role assumption
3. AWS default credential chain: env vars, `~/.aws/credentials`, instance metadata

The simplest approach for local development is to rely on the default credential chain by setting environment variables:

```bash
export AWS_ACCESS_KEY_ID=<your-access-key-id>
export AWS_SECRET_ACCESS_KEY=<your-secret-access-key>
export AWS_REGION=<your-region>
```

## Plugin Configuration

### Piped Config

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  plugins:
  - name: ecs
    port: 7003 # any unused port
    url: file:///path/to/.piped/plugins/ecs # or remoteUrl(TBD)
    deployTargets:
      - name: production
        config:
          region: us-east-1
```

| Field | Type | Description | Required |
|-|-|-|-|
| deployTargets | [][DeployTargetConfig](#deploytargetconfig) | The config for the destinations to deploy applications | Yes |

#### DeployTargetConfig

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The name of the deploy target. | Yes |
| labels | map[string]string | The labels of the deploy target. | No |
| config | [ECSDeployTargetConfig](#ecsdeploytargetconfig) | The configuration of the deploy target for ECS plugin. | No |

##### ECSDeployTargetConfig

| Field | Type | Description | Required |
|-|-|-|-|
| region | string | The AWS region where the ECS cluster is located (e.g., `us-west-2`). | Yes |
| profile | string | The AWS profile to use from the credentials file. If empty, uses the default profile. | No |
| credentialsFile | string | The path to the AWS shared credentials file (e.g., `~/.aws/credentials`). If empty, uses the default location. | No |
| roleARN | string | The IAM role ARN to assume when accessing AWS resources. Required when assuming a role across accounts. | No |
| tokenFile | string | The path to the OIDC token file for web identity federation. Required when `roleARN` is set for OIDC-based authentication. | No |

### Application Config

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  plugins:
    ecs: # same name as the one defined in `spec.plugins[].name`
      input:
        ...
      quickSync:
        ...
  pipeline:
    stages:
      ...
```

| Field | Type | Description | Required |
|-|-|-|-|
| input | [ECSDeploymentInput](#ecsdeploymentinput) | Input for ECS deployment such as task definition file, service definition file, target groups... | No |
| quickSync | [ECSSyncStageOptions](#ecssyncstageoptions) | Options for the quick sync. | No |

#### ECSDeploymentInput

| Field | Type | Description | Required |
|-|-|-|-|
| taskDefinitionFile | string | The path to the task definition file (YAML/JSON). Default is `taskdef.json`. | No |
| serviceDefinitionFile | string | The path to the service definition file (YAML/JSON). Required for service-based deployments. | No |
| runStandaloneTask | bool | Whether to run the task as a standalone task without creating/updating an ECS service. Only Quick Sync is supported for standalone tasks. Default is `false`. | No |
| clusterArn | string | The ARN of the ECS cluster where the task and service will be deployed. | No |
| launchType | string | The launch type on which to run the task. Valid values: `EC2`, `FARGATE`. Default is `FARGATE`. | No |
| accessType | string | How the ECS service is accessed. Valid values: `ELB`, `SERVICE_DISCOVERY`. Default is `ELB`. | No |
| awsvpcConfiguration | [ECSVpcConfiguration](#ecsvpcconfiguration) | The VPC configuration for running ECS tasks. Required for `FARGATE` launch type or `EC2` with `awsvpc` network mode. | No |
| targetGroups | [ECSTargetGroups](#ecstargetgroups) | The load balancer target groups for the ECS service. Required when `accessType` is `ELB`. | No |

#### ECSVpcConfiguration

| Field | Type | Description | Required |
|-|-|-|-|
| subnets | []string | List of VPC subnet IDs where tasks will be launched. Maximum 16 subnets. | Yes |
| assignPublicIp | string | Whether to assign a public IP address to the task's ENI. Valid values: `ENABLED`, `DISABLED`. | No |
| securityGroups | []string | List of security group IDs associated with the task's ENI. Maximum 5 security groups. If not specified, the default security group for the VPC will be used. | No |

#### ECSTargetGroups

| Field | Type | Description | Required |
|-|-|-|-|
| primary | [ECSTargetGroup](#ecstargetgroup) | The target group for the primary service. | No |
| canary | [ECSTargetGroup](#ecstargetgroup) | The target group for the canary service. Required to enable canary/blue-green deployment strategy. | No |

##### ECSTargetGroup

| Field | Type | Description | Required |
|-|-|-|-|
| targetGroupArn | string | The ARN of the target group. | No |
| containerName | string | The name of the container to associate with the target group. | No |
| containerPort | int | The port on the container to associate with the target group. | No |

### Stage Config

```yaml
pipeline:
  stages:
    - name: ECS_CANARY_ROLLOUT
      with:
        ...
    - name: ECS_PRIMARY_ROLLOUT
    - name: ECS_TRAFFIC_ROUTING
      with:
        ...
    - name: ECS_CANARY_CLEAN
```

#### `ECS_SYNC`

| Field | Type | Description | Required |
|-|-|-|-|
| recreate | bool | Whether to recreate the service. Enabling this will stop all running tasks before creating a new task set. Default is `false`. | No |

#### `ECS_PRIMARY_ROLLOUT`

No configuration options.

#### `ECS_CANARY_ROLLOUT`

| Field | Type | Description | Required |
|-|-|-|-|
| scale | float | The percentage of tasks to run as canary relative to the current primary workload (0-100). | No |

#### `ECS_CANARY_CLEAN`

No configuration options.

#### `ECS_ROLLBACK`

No configuration options. This stage is automatically appended to the pipeline when `autoRollback` is enabled in the application configuration. It restores the ECS service to the previous task set.

#### `ECS_TRAFFIC_ROUTING`

| Field | Type | Description | Required |
|-|-|-|-|
| primary | int | The percentage of traffic to route to the primary variant. If set, canary receives `100 - primary` percent. | No |
| canary | int | The percentage of traffic to route to the canary variant. If set, primary receives `100 - canary` percent. | No |
