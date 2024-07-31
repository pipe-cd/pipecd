---
title: "Configuration reference"
linkTitle: "Configuration reference"
weight: 11
description: >
  This page describes all configurable fields in the application configuration and analysis template.
---

## Kubernetes Application

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  input:
  pipeline:
  ...
```

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The application name. | Yes (if you want to create PipeCD application through the application configuration file) |
| labels | map[string]string | Additional attributes to identify applications. | No |
| description | string | Notes on the Application. | No |
| input | [KubernetesDeploymentInput](#kubernetesdeploymentinput) | Input for Kubernetes deployment such as kubectl version, helm version, manifests filter... | No |
| trigger | [DeploymentTrigger](#deploymenttrigger) | Configuration for trigger used to determine should we trigger a new deployment or not. | No |
| planner | [DeploymentPlanner](#deploymentplanner) | Configuration for planner used while planning deployment. | No |
| commitMatcher | [CommitMatcher](#commitmatcher) | Forcibly use QuickSync or Pipeline when commit message matched the specified pattern. | No |
| quickSync | [KubernetesQuickSync](#kubernetesquicksync) | Configuration for quick sync. | No |
| pipeline | [Pipeline](#pipeline) | Pipeline for deploying progressively. | No |
| service | [KubernetesService](#kubernetesservice) | Which Kubernetes resource should be considered as the Service of application. Empty means the first Service resource will be used. | No |
| workloads | [][KubernetesWorkload](#kubernetesworkload) | Which Kubernetes resources should be considered as the Workloads of application. Empty means all Deployment resources. | No |
| trafficRouting | [KubernetesTrafficRouting](#kubernetestrafficrouting) | How to change traffic routing percentages. | No |
| encryption | [SecretEncryption](#secretencryption) | List of encrypted secrets and targets that should be decrypted before using. | No |
| attachment | [Attachment](#attachment) | List of attachment sources and targets that should be attached to manifests before using. | No |
| timeout | duration | The maximum length of time to execute deployment before giving up. Default is 6h. | No |
| notification | [DeploymentNotification](#deploymentnotification) | Additional configuration used while sending notification to external services. | No |
| postSync | [PostSync](#postsync) | Additional configuration used as extra actions once the deployment is triggered. | No |
| variantLabel | [KubernetesVariantLabel](#kubernetesvariantlabel) | The label will be configured to variant manifests used to distinguish them. | No |
| eventWatcher | [][EventWatcher](#eventwatcher) | List of configurations for event watcher. | No |
| driftDetection | [DriftDetection](#driftdetection) | Configuration for drift detection. | No |

### Annotations

Kubernetes resources can be managed by some annotations provided by PipeCD.

| Annotation key | Target resource(s) | Possible values | Description |
|-|-|-|-|
| `pipecd.dev/ignore-drift-detection` | any | "true" | Whether the drift detection should ignore this resource. |
| `pipecd.dev/server-side-apply` | any | "true" | Use server side apply instead of client side apply. |

## Terraform application

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: TerraformApp
spec:
  input:
  pipeline:
  ...
```

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The application name. | Yes if you set the application through the application configuration file |
| labels | map[string]string | Additional attributes to identify applications. | No |
| description | string | Notes on the Application. | No |
| input | [TerraformDeploymentInput](#terraformdeploymentinput) | Input for Terraform deployment such as terraform version, workspace... | No |
| trigger | [DeploymentTrigger](#deploymenttrigger) | Configuration for trigger used to determine should we trigger a new deployment or not. | No |
| planner | [DeploymentPlanner](#deploymentplanner) | Configuration for planner used while planning deployment. | No |
| quickSync | [TerraformQuickSync](#terraformquicksync) | Configuration for quick sync. | No |
| pipeline | [Pipeline](#pipeline) | Pipeline for deploying progressively. | No |
| encryption | [SecretEncryption](#secretencryption) | List of encrypted secrets and targets that should be decrypted before using. | No |
| attachment | [Attachment](#attachment) | List of attachment sources and targets that should be attached to manifests before using. | No |
| timeout | duration | The maximum length of time to execute deployment before giving up. Default is 6h. | No |
| notification | [DeploymentNotification](#deploymentnotification) | Additional configuration used while sending notification to external services. | No |
| postSync | [PostSync](#postsync) | Additional configuration used as extra actions once the deployment is triggered. | No |
| eventWatcher | [][EventWatcher](#eventwatcher) | List of configurations for event watcher. | No |

## Cloud Run application

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: CloudRunApp
spec:
  input:
  pipeline:
  ...
```

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The application name. | Yes if you set the application through the application configuration file |
| labels | map[string]string | Additional attributes to identify applications. | No |
| description | string | Notes on the Application. | No |
| input | [CloudRunDeploymentInput](#cloudrundeploymentinput) | Input for Cloud Run deployment such as docker image... | No |
| trigger | [DeploymentTrigger](#deploymenttrigger) | Configuration for trigger used to determine should we trigger a new deployment or not. | No |
| planner | [DeploymentPlanner](#deploymentplanner) | Configuration for planner used while planning deployment. | No |
| quickSync | [CloudRunQuickSync](#cloudrunquicksync) | Configuration for quick sync. | No |
| pipeline | [Pipeline](#pipeline) | Pipeline for deploying progressively. | No |
| encryption | [SecretEncryption](#secretencryption) | List of encrypted secrets and targets that should be decrypted before using. | No |
| attachment | [Attachment](#attachment) | List of attachment sources and targets that should be attached to manifests before using. | No |
| timeout | duration | The maximum length of time to execute deployment before giving up. Default is 6h. | No |
| notification | [DeploymentNotification](#deploymentnotification) | Additional configuration used while sending notification to external services. | No |
| postSync | [PostSync](#postsync) | Additional configuration used as extra actions once the deployment is triggered. | No |
| eventWatcher | [][EventWatcher](#eventwatcher) | List of configurations for event watcher. | No |

## Lambda application

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: LambdaApp
spec:
  pipeline:
  ...
```

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The application name. | Yes if you set the application through the application configuration file |
| labels | map[string]string | Additional attributes to identify applications. | No |
| description | string | Notes on the Application. | No |
| input | [LambdaDeploymentInput](#lambdadeploymentinput) | Input for Lambda deployment such as path to function manifest file... | No |
| architectures | []string| Specific architecture for which a function supports (Default x86_64). | No |
| trigger | [DeploymentTrigger](#deploymenttrigger) | Configuration for trigger used to determine should we trigger a new deployment or not. | No |
| planner | [DeploymentPlanner](#deploymentplanner) | Configuration for planner used while planning deployment. | No |
| quickSync | [LambdaQuickSync](#lambdaquicksync) | Configuration for quick sync. | No |
| pipeline | [Pipeline](#pipeline) | Pipeline for deploying progressively. | No |
| encryption | [SecretEncryption](#secretencryption) | List of encrypted secrets and targets that should be decrypted before using. | No |
| attachment | [Attachment](#attachment) | List of attachment sources and targets that should be attached to manifests before using. | No |
| timeout | duration | The maximum length of time to execute deployment before giving up. Default is 6h. | No |
| notification | [DeploymentNotification](#deploymentnotification) | Additional configuration used while sending notification to external services. | No |
| postSync | [PostSync](#postsync) | Additional configuration used as extra actions once the deployment is triggered. | No |
| eventWatcher | [][EventWatcher](#eventwatcher) | List of configurations for event watcher. | No |

## ECS application

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: ECSApp
spec:
  input:
  pipeline:
  ...
```

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The application name. | Yes if you set the application through the application configuration file |
| labels | map[string]string | Additional attributes to identify applications. | No |
| description | string | Notes on the Application. | No |
| input | [ECSDeploymentInput](#ecsdeploymentinput) | Input for ECS deployment such as path to TaskDefinition, Service... | No |
| trigger | [DeploymentTrigger](#deploymenttrigger) | Configuration for trigger used to determine should we trigger a new deployment or not. | No |
| planner | [DeploymentPlanner](#deploymentplanner) | Configuration for planner used while planning deployment. | No |
| quickSync | [ECSQuickSync](#ecsquicksync) | Configuration for quick sync. | No |
| pipeline | [Pipeline](#pipeline) | Pipeline for deploying progressively. | No |
| encryption | [SecretEncryption](#secretencryption) | List of encrypted secrets and targets that should be decrypted before using. | No |
| attachment | [Attachment](#attachment) | List of attachment sources and targets that should be attached to manifests before using. | No |
| timeout | duration | The maximum length of time to execute deployment before giving up. Default is 6h. | No |
| notification | [DeploymentNotification](#deploymentnotification) | Additional configuration used while sending notification to external services. | No |
| postSync | [PostSync](#postsync) | Additional configuration used as extra actions once the deployment is triggered. | No |
| eventWatcher | [][EventWatcher](#eventwatcher) | List of configurations for event watcher. | No |

## Analysis Template Configuration

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: AnalysisTemplate
spec:
  metrics:
    grpc_error_rate_percentage:
      interval: 1m
      provider: prometheus-dev
      failureLimit: 1
      expected:
        max: 10
      query: awesome_query
```

| Field | Type | Description | Required |
|-|-|-|-|
| metrics | map[string][AnalysisMetrics](#analysismetrics) | Template for metrics. | No |

## Event Watcher Configuration (deprecated)

```yaml
apiVersion: pipecd.dev/v1beta1
kind: EventWatcher
spec:
  events:
    - name: helloworld-image-update
      replacements:
        - file: helloworld/deployment.yaml
          yamlField: $.spec.template.spec.containers[0].image
```

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The event name. | Yes |
| labels | map[string]string | Additional attributes of event. This can make an event definition unique even if the one with the same name exists. | No |
| replacements | [][EventWatcherReplacement](#eventwatcherreplacement) | List of places where will be replaced when the new event matches. | Yes |

### EventWatcherReplacement
One of `yamlField` or `regex` is required.

| Field | Type | Description | Required |
|-|-|-|-|
| file | string | The relative path from the repository root to the file to be updated. | Yes |
| yamlField | string | The yaml path to the field to be updated. It requires to start with `$` which represents the root element. e.g. `$.foo.bar[0].baz`. | No |
| regex | string | The regex string that specify what should be replaced. The only first capturing group enclosed by `()` will be replaced with the new value. e.g. `host.xz/foo/bar:(v[0-9].[0-9].[0-9])` | No |

## CommitMatcher

| Field | Type | Description | Required |
|-|-|-|-|
| quickSync | string | Regular expression string to forcibly do QuickSync when it matches the commit message. | No |
| pipeline | string | Regular expression string to forcibly do Pipeline when it matches the commit message. | No |

## SecretEncryption

| Field | Type | Description | Required |
|-|-|-|-|
| encryptedSecrets | map[string]string | List of encrypted secrets. | No |
| decryptionTargets | []string | List of files to be decrypted before using. | No |

## Attachment

| Field | Type | Description | Required |
|-|-|-|-|
| sources | map[string]string | List of attaching files with key is its refer name. | No |
| targets | []string | List of files which should contain the attachments. | No |

## DeploymentPlanner

| Field | Type | Description | Required |
|-|-|-|-|
| alwaysUsePipeline | bool | Always use the defined pipeline to deploy the application in all deployments. Default is `false`. | No |

## DeploymentTrigger

| Field | Type | Description | Required |
|-|-|-|-|
| onCommit | [OnCommit](#oncommit) | Controls triggering new deployment when new Git commits touched the application. | No |
| onCommand | [OnCommand](#oncommand) | Controls triggering new deployment when received a new `SYNC` command. | No |
| onOutOfSync | [OnOutOfSync](#onoutofsync) | Controls triggering new deployment when application is at `OUT_OF_SYNC` state. | No |
| onChain | [OnChain](#onchain) | Controls triggering new deployment when the application is counted as a node of some chains. | No |

### OnCommit

| Field | Type | Description | Required |
|-|-|-|-|
| disabled | bool | Whether to exclude application from triggering target when new Git commits touched it. Default is `false`. | No |
| paths | []string | List of directories or files where any changes of them will be considered as touching the application. Regular expression can be used. Empty means watching all changes under the application directory. | No |
| ignores | []string | List of directories or files where any changes of them will NOT be considered as touching the application. Regular expression can be used. This config has a higher priority compare to `paths`. | No |

### OnCommand

| Field | Type | Description | Required |
|-|-|-|-|
| disabled | bool | Whether to exclude application from triggering target when received a new `SYNC` command. Default is `false`. | No |

### OnOutOfSync

| Field | Type | Description | Required |
|-|-|-|-|
| disabled | bool | Whether to exclude application from triggering target when application is at `OUT_OF_SYNC` state. Default is `true`. | No |
| minWindow | duration | Minimum amount of time must be elapsed since the last deployment. This can be used to avoid triggering unnecessary continuous deployments based on `OUT_OF_SYNC` status. Default is `5m`. | No |

### OnChain

| Field | Type | Description | Required |
|-|-|-|-|
| disabled | bool | Whether to exclude application from triggering target when application is counted as a node of some chains. Default is `true`. | No |

## Pipeline

| Field | Type | Description | Required |
|-|-|-|-|
| stages | [][PipelineStage](#pipelinestage) | List of deployment pipeline stages. | No |

### PipelineStage

| Field | Type | Description | Required |
|-|-|-|-|
| id | string | The unique ID of the stage. | No |
| name | string | One of the provided stage names. | Yes |
| desc | string | The description about the stage. | No |
| timeout | duration | The maximum time the stage can be taken to run. | No |
| with | [StageOptions](#stageoptions) | Specific configuration for the stage. This must be one of these [StageOptions](#stageoptions). | No |

## DeploymentNotification

| Field | Type | Description | Required |
|-|-|-|-|
| mentions | [][NotificationMention](#notificationmention) | List of users to be notified for each event. | No |

### NotificationMention

| Field | Type | Description | Required |
|-|-|-|-|
| event | string | The event to be notified to users. | Yes |
| slack | []string | List of user IDs for mentioning in Slack. See [here](https://api.slack.com/reference/surfaces/formatting#mentioning-users) for more information on how to check them. | No |

## KubernetesDeploymentInput

| Field | Type | Description | Required |
|-|-|-|-|
| manifests | []string | List of manifest files in the application directory used to deploy. Empty means all manifest files in the directory will be used. | No |
| kubectlVersion | string | Version of kubectl will be used. Empty means the version set on [piped config](../managing-piped/configuration-reference/#platformproviderkubernetesconfig) or [default version](https://github.com/pipe-cd/pipecd/blob/master/tool/piped-base/install-kubectl.sh#L24) will be used. | No |
| kustomizeVersion | string | Version of kustomize will be used. Empty means the [default version](https://github.com/pipe-cd/pipecd/blob/master/tool/piped-base/install-kustomize.sh#L24) will be used. | No |
| kustomizeOptions | map[string]string | List of options that should be used by Kustomize commands. | No |
| helmVersion | string | Version of helm will be used. Empty means the [default version](https://github.com/pipe-cd/pipecd/blob/master/tool/piped-base/install-helm.sh#L24) will be used. | No |
| helmChart | [HelmChart](#helmchart) | Where to fetch helm chart. | No |
| helmOptions | [HelmOptions](#helmoptions) | Configurable parameters for helm commands. | No |
| namespace | string | The namespace where manifests will be applied. | No |
| autoRollback | bool | Automatically reverts all deployment changes on failure. Default is `true`. | No |
| autoCreateNamespace | bool | Automatically create a new namespace if it does not exist. Default is `false`. | No |

### HelmChart

| Field | Type | Description | Required |
|-|-|-|-|
| gitRemote | string | Git remote address where the chart is placing. Empty means the same repository. | No |
| ref | string | The commit SHA or tag value. Only valid when gitRemote is not empty. | No |
| path | string | Relative path from the repository root to the chart directory. | No |
| repository | string | The name of a registered Helm Chart Repository. | No |
| name | string | The chart name. | No |
| version | string | The chart version. | No |

### HelmOptions

| Field | Type | Description | Required |
|-|-|-|-|
| releaseName | string | The release name of helm deployment. By default, the release name is equal to the application name. | No |
| setValues | map[string]string | List of values. | No |
| valueFiles | []string | List of value files should be loaded. Only local files stored under the application directory or remote files served at the http(s) endpoint are allowed. | No |
| setFiles | map[string]string | List of file path for values. | No |
| apiVersions | []string | Kubernetes api versions used for Capabilities.APIVersions. | No |
| kubeVersion | string | Kubernetes version used for Capabilities.KubeVersion. | No |

## KubernetesVariantLabel

| Field | Type | Description | Required |
|-|-|-|-|
| key | string | The key of the label. Default is `pipecd.dev/variant`. | No |
| primaryValue | string | The label value for PRIMARY variant. Default is `primary`. | No |
| canaryValue | string | The label value for CANARY variant. Default is `canary`. | No |
| baselineValue | string | The label value for BASELINE variant. Default is `baseline`. | No |

## KubernetesQuickSync

| Field | Type | Description | Required |
|-|-|-|-|
| addVariantLabelToSelector | bool | Whether the PRIMARY variant label should be added to manifests if they were missing. Default is `false`. | No |
| prune | bool | Whether the resources that are no longer defined in Git should be removed or not. Default is `false` | No |

## KubernetesService

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The name of Service manifest. | No |

## KubernetesWorkload

| Field | Type | Description | Required |
|-|-|-|-|
| kind | string | The kind name of workload manifests. Currently, only `Deployment` is supported. In the future, we also want to support `ReplicationController`, `DaemonSet`, `StatefulSet`. | No |
| name | string | The name of workload manifest. | No |

## KubernetesTrafficRouting

| Field | Type | Description | Required |
|-|-|-|-|
| method | string | Which traffic routing method will be used. Available values are `istio`, `smi`, `podselector`. Default is `podselector`. | No |
| istio | [IstioTrafficRouting](#istiotrafficrouting)| Istio configuration when the method is `istio`. | No |

### IstioTrafficRouting

| Field | Type | Description | Required |
|-|-|-|-|
| editableRoutes | []string | List of routes in the VirtualService that can be changed to update traffic routing. Empty means all routes should be updated. | No |
| host | string | The service host. | No |
| virtualService | [IstioVirtualService](#istiovirtualservice) | The reference to VirtualService manifest. Empty means the first VirtualService resource will be used. | No |

#### IstioVirtualService

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The name of VirtualService manifest. | No |

## TerraformDeploymentInput

| Field | Type | Description | Required |
|-|-|-|-|
| workspace | string | The terraform workspace name. Empty means `default` workspace. | No |
| terraformVersion | string | The version of terraform should be used. Empty means the pre-installed version will be used. | No |
| vars | []string | List of variables that will be set directly on terraform commands with `-var` flag. The variable must be formatted by `key=value`. | No |
| varFiles | []string | List of variable files that will be set on terraform commands with `-var-file` flag. | No |
| commandFlags | [TerraformCommandFlags](#terraformcommandflags) | List of additional flags will be used while executing terraform commands. | No |
| commandEnvs | [TerraformCommandEnvs](#terraformcommandenvs) | List of additional environment variables will be used while executing terraform commands. | No |
| autoRollback | bool | Automatically reverts all changes from all stages when one of them failed. | No |

### TerraformCommandFlags

| Field | Type | Description | Required |
|-|-|-|-|
| shared | []string | List of additional flags used for all Terraform commands. | No |
| init | []string | List of additional flags used for Terraform `init` command. | No |
| plan | []string | List of additional flags used for Terraform `plan` command. | No |
| apply | []string | List of additional flags used for Terraform `apply` command. | No |

### TerraformCommandEnvs

| Field | Type | Description | Required |
|-|-|-|-|
| shared | []string | List of additional environment variables used for all Terraform commands. | No |
| init | []string | List of additional environment variables used for Terraform `init` command. | No |
| plan | []string | List of additional environment variables used for Terraform `plan` command. | No |
| apply | []string | List of additional environment variables used for Terraform `apply` command. | No |

## TerraformQuickSync

| Field | Type | Description | Required |
|-|-|-|-|
| retries | int | How many times to retry applying terraform changes. Default is `0`. | No |

## CloudRunDeploymentInput

| Field | Type | Description | Required |
|-|-|-|-|
| serviceManifestFile | string | The name of service manifest file placing in application directory. Default is `service.yaml`. | No |
| autoRollback | bool | Automatically reverts to the previous state when the deployment is failed. Default is `true`. | No |

## CloudRunQuickSync

| Field | Type | Description | Required |
|-|-|-|-|

## LambdaDeploymentInput

| Field | Type | Description | Required |
|-|-|-|-|
| functionManifestFile | string | The name of function manifest file placing in application directory. Default is `function.yaml`. | No |
| autoRollback | bool | Automatically reverts to the previous state when the deployment is failed. Default is `true`. | No |

### Specific function.yaml

One of `image`, `s3Bucket`, or `source` is required.

- If you use `s3Bucket`, `s3Key` and `s3ObjectVersion` are required.

- If you use `s3Bucket` or `source`, `handler` and `runtime` are required.

See [Configuring Lambda application](../managing-application/defining-app-configuration/lambda) for more details.

| Field            | Type             | Description                        | Required |
|------------------|------------------|------------------------------------|----------|
| name             | string           | Name of the Lambda function        | Yes      |
| role             | string           | IAM role ARN                       | Yes      |
| image            | string           | URI of the container image         | No      |
| s3Bucket         | string           | S3 bucket name for code package   | No      |
| s3Key            | string           | S3 key for code package            | No      |
| s3ObjectVersion  | string           | S3 object version for code package | No      |
| source       | [Source](#source)       | Git settings                | No      |
| handler          | string           | Lambda function handler            | No      |
| runtime          | string           | Runtime environment                | No      |
| architectures    | [][Architecture](#architecture)   | Supported architectures            | No       |
| ephemeralStorage | [EphemeralStorage](#ephemeralstorage)| Ephemeral storage configuration    | No       |
| memory           | int32            | Memory allocation (in MB)          | Yes      |
| timeout          | int32            | Function timeout (in seconds)      | Yes      |
| tags             | map[string]string| Key-value pairs for tags           | No       |
| environments     | map[string]string| Environment variables              | No       |
| vpcConfig        | [VPCConfig](#vpcconfig)       | VPC configuration                  | No       |
| layers        | []string       | ARNs of [layers](https://docs.aws.amazon.com/lambda/latest/dg/chapter-layers.html) to depend on                | No       |

#### Source

| Field | Type   | Description              | Required |
|-------|--------|--------------------------|----------|
| git   | string | Git repository URL       | Yes      |
| ref   | string | Git branch/tag/reference| Yes      |
| path  | string | Path within the repository | Yes    |

#### Architecture

| Field | Type   | Description            | Required |
|-------|--------|------------------------|----------|
| name  | string | Name of the architecture | Yes     |

#### EphemeralStorage

| Field | Type  | Description                  | Required |
|-------|-------|------------------------------|----------|
| size  | int32 | Size of the ephemeral storage| Yes       |

#### VPCConfig

| Field           | Type     | Description                 | Required |
|-----------------|----------|-----------------------------|----------|
| securityGroupIds| []string | List of security group IDs  | No       |
| subnetIds       | []string | List of subnet IDs          | No       |


## LambdaQuickSync

| Field | Type | Description | Required |
|-|-|-|-|

## ECSDeploymentInput

| Field | Type | Description | Required |
|-|-|-|-|
| serviceDefinitionFile | string | The path ECS Service configuration file. Allow file in both `yaml` and `json` format. The default value is `service.json`. See [here](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/service_definition_parameters.html) for parameters.| No |
| taskDefinitionFile | string | The path to ECS TaskDefinition configuration file. Allow file in both `yaml` and `json` format. The default value is `taskdef.json`. See [here](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html) for parameters. | No |
| targetGroups | [ECSTargetGroupInput](#ecstargetgroupinput) | The target groups configuration, will be used to routing traffic to created task sets. | Yes (if you want to perform progressive delivery) |
| runStandaloneTask | bool | Run standalone tasks during deployments. About standalone task, see [here](https://docs.aws.amazon.com/AmazonECS/latest/userguide/ecs_run_task-v2.html). The default value is `true`. |
| accessType | string | How the ECS service is accessed. One of `ELB` or `SERVICE_DISCOVERY`. See examples [here](https://github.com/pipe-cd/examples/tree/master/ecs/servicediscovery/simple). The default value is `ELB`. |

### ECSTargetGroupInput

| Field | Type | Description | Required |
|-|-|-|-|
| primary | [ECSTargetGroupObject](#ecstargetgroupobject) | The PRIMARY target group, will be used to register the PRIMARY ECS task set. | Yes |
| canary | [ECSTargetGroupObject](#ecstargetgroupobject) | The CANARY target group, will be used to register the CANARY ECS task set if exist. It's required to enable PipeCD to perform the multi-stage deployment. | No |

#### ECSTargetGroupObject

| Field | Type | Description | Required |
|-|-|-|-|
| targetGroupArn | string | The name of the container (as it appears in a container definition) to　associate with the load balancer | Yes |
| containerName | string | The full Amazon Resource Name (ARN) of the Elastic Load Balancing target group or groups associated with a service or task set. | Yes |
| containerPort | int | The port on the container to associate with the load balancer. | Yes |
| LoadBalancerName | string | The name of the load balancer to associate with the Amazon ECS service or task set. | No |

Note: The available values are identical to those found in the aws-sdk-go-v2 Types.LoadBalancer. For more details, please refer to [this link](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ecs/types#LoadBalancer) .

## ECSQuickSync

| Field | Type | Description | Required |
|-|-|-|-|
| recreate | bool | Whether to delete old tasksets before creating new ones or not. Default to false. | No |

## AnalysisMetrics

| Field | Type | Description | Required |
|-|-|-|-|
| provider | string | The unique name of provider defined in the Piped Configuration. | Yes |
| strategy | string | The strategy name. One of `THRESHOLD` or `PREVIOUS` or `CANARY_BASELINE` or `CANARY_PRIMARY` is available. Defaults to `THRESHOLD`. | No |
| query | string | A query performed against the [Analysis Provider](../../concepts/#analysis-provider). The stage will be skipped if no data points were returned. | Yes |
| expected | [AnalysisExpected](#analysisexpected) | The statically defined expected query result. This field is ignored if there was no data point as a result of the query. | Yes if the strategy is `THRESHOLD` |
| interval | duration | Run a query at specified intervals. | Yes |
| failureLimit | int | Acceptable number of failures. e.g. If 1 is set, the `ANALYSIS` stage will end with failure after two queries results failed. Defaults to 1. | No |
| skipOnNoData | bool | If true, it considers as a success when no data returned from the analysis provider. Defaults to false. | No |
| deviation | string | The stage fails on deviation in the specified direction. One of `LOW` or `HIGH` or `EITHER` is available. This can be used only for `PREVIOUS`, `CANARY_BASELINE` or `CANARY_PRIMARY`. Defaults to `EITHER`. | No |
| baselineArgs | map[string][string] | The custom arguments to be populated for the Baseline query. They can be reffered as `{{ .VariantCustomArgs.xxx }}`. | No |
| canaryArgs | map[string][string] | The custom arguments to be populated for the Canary query. They can be reffered as `{{ .VariantCustomArgs.xxx }}`. | No |
| primaryArgs | map[string][string] | The custom arguments to be populated for the Primary query. They can be reffered as `{{ .VariantCustomArgs.xxx }}`. | No |
| timeout | duration | How long after which the query times out. | No |
| template | [AnalysisTemplateRef](#analysistemplateref) | Reference to the template to be used. | No |


### AnalysisExpected

| Field | Type | Description | Required |
|-|-|-|-|
| min | float64 | Failure, if the query result is less than this value. | No |
| max | float64 | Failure, if the query result is larger than this value. | No |

### AnalysisTemplateRef

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The template name to refer. | Yes |
| appArgs | map[string]string | The arguments for custom-args. | No |

## AnalysisLog

| Field | Type | Description | Required |
|-|-|-|-|

## AnalysisHttp

| Field | Type | Description | Required |
|-|-|-|-|

## StageOptions

### KubernetesPrimaryRolloutStageOptions

| Field | Type | Description | Required |
|-|-|-|-|
| suffix | string | Suffix that should be used when naming the PRIMARY variant's resources. Default is `primary`. | No |
| createService | bool | Whether the PRIMARY service should be created. Default is `false`. | No |
| addVariantLabelToSelector | bool | Whether the PRIMARY variant label should be added to manifests if they were missing. Default is `false`. | No |
| prune | bool | Whether the resources that are no longer defined in Git should be removed or not. Default is `false` | No |

### KubernetesCanaryRolloutStageOptions

| Field | Type | Description | Required |
|-|-|-|-|
| replicas | int | How many pods for CANARY workloads. Default is `1` pod. Alternatively, can be specified a string suffixed by "%" to indicate a percentage value compared to the pod number of PRIMARY | No |
| suffix | string | Suffix that should be used when naming the CANARY variant's resources. Default is `canary`. | No |
| createService | bool | Whether the CANARY service should be created. Default is `false`. | No |
| patches | [][KubernetesResourcePatch](#kubernetesresourcepatch) | List of patches used to customize manifests for CANARY variant. | No |

### KubernetesCanaryCleanStageOptions

| Field | Type | Description | Required |
|-|-|-|-|
| | | | |

### KubernetesBaselineRolloutStageOptions

| Field | Type | Description | Required |
|-|-|-|-|
| replicas | int | How many pods for BASELINE workloads. Default is `1` pod. Alternatively, can be specified a string suffixed by "%" to indicate a percentage value compared to the pod number of PRIMARY | No |
| suffix | string | Suffix that should be used when naming the BASELINE variant's resources. Default is `baseline`. | No |
| createService | bool | Whether the BASELINE service should be created. Default is `false`. | No |

### KubernetesBaselineCleanStageOptions

| Field | Type | Description | Required |
|-|-|-|-|
| | | | |

### KubernetesTrafficRoutingStageOptions
This stage routes traffic with the method specified in [KubernetesTrafficRouting](#kubernetestrafficrouting).
When using `podselector` method as a traffic routing method, routing is done by updating the Service selector.
Therefore, note that all traffic will be routed to the primary if the the primary variant's service is rolled out by running the `K8S_PRIMARY_ROLLOUT` stage.

| Field | Type | Description | Required |
|-|-|-|-|
| all | string | Which variant should receive all traffic. Available values are "primary", "canary", "baseline". Default is `primary`. | No |
| primary | [Percentage](#percentage) | The percentage of traffic should be routed to PRIMARY variant. | No |
| canary | [Percentage](#percentage) | The percentage of traffic should be routed to CANARY variant. | No |
| baseline | [Percentage](#percentage) | The percentage of traffic should be routed to BASELINE variant. | No |

### TerraformPlanStageOptions

| Field | Type | Description | Required |
|-|-|-|-|
| exitOnNoChanges | bool | Whether exiting the pipeline when the result has no changes | No |

### TerraformApplyStageOptions

| Field | Type | Description | Required |
|-|-|-|-|
| retries | int | How many times to retry applying terraform changes. Default is `0`. | No |

### CloudRunPromoteStageOptions

| Field | Type | Description | Required |
|-|-|-|-|
| percent | [Percentage](#percentage) | Percentage of traffic should be routed to the new version. | No |

### LambdaCanaryRolloutStageOptions

| Field | Type | Description | Required |
|-|-|-|-|

### LambdaPromoteStageOptions

| Field | Type | Description | Required |
|-|-|-|-|
| percent | [Percentage](#percentage) | Percentage of traffic should be routed to the new version. | No |

### ECSPrimaryRolloutStageOptions

| Field | Type | Description | Required |
|-|-|-|-|

### ECSCanaryRolloutStageOptions

| Field | Type | Description | Required |
|-|-|-|-|
| scale | [Percentage](#percentage) | The percentage of workloads should be rolled out as CANARY variant's workload. | Yes |

### ECSTrafficRoutingStageOptions

| Field | Type | Description | Required |
|-|-|-|-|
| primary | [Percentage](#percentage) | The percentage of traffic should be routed to PRIMARY variant. | No |
| canary | [Percentage](#percentage) | The percentage of traffic should be routed to CANARY variant. | No |

Note: By default, the sum of traffic is rounded to 100. If both `primary` and `canary` numbers are not set, the PRIMARY variant will receive 100% while the CANARY variant will receive 0% of the traffic.

### AnalysisStageOptions

| Field | Type | Description | Required |
|-|-|-|-|
| duration | duration | Maximum time to perform the analysis. | Yes |
| metrics | [][AnalysisMetrics](#analysismetrics) | Configuration for analysis by metrics. | No |

### WaitStageOptions

| Field | Type | Description | Required |
|-|-|-|-|
| duration | duration | Time to wait. | Yes |

### WaitApprovalStageOptions

| Field | Type | Description | Required |
|-|-|-|-|
| timeout | duration | The maximum length of time to wait before giving up. Default is 6h. | No |
| approvers | []string | List of username who has permission to approve. | Yes |
| minApproverNum | int | Number of minimum needed approvals to make this stage complete. Default is 1. | No |

### CustomSyncStageOptions (deprecated)
| Field | Type | Description | Required |
|-|-|-|-|
| timeout | duration | The maximum time the stage can be taken to run. Default is `6h`| No |
| envs | map[string]string | Environment variables used with scripts. | No |
| run | string | Script run on this stage. | Yes |

### ScriptRunStageOptions
| Field | Type | Description | Required |
|-|-|-|-|
| run | string | Script run on this stage. | Yes |
| env | map[string]string | Environment variables used with scripts. | No |
| timeout | duration | The maximum time the stage can be taken to run. Default is `6h`| No |

## PostSync

| Field | Type | Description | Required |
|-|-|-|-|
| chain | [DeploymentChain](#deploymentchain) | Deployment chain configuration, used to determine and build deployments that should be triggered once the current deployment is triggered. | No |

### DeploymentChain

| Field | Type | Description | Required |
|-|-|-|-|
| applications | [][DeploymentChainApplication](#deploymentchainapplication) | The list of applications which should be triggered once deployment of this application rolled out successfully. | Yes |

#### DeploymentChainApplication

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The name of PipeCD application, note that application name is not unique in PipeCD datastore | No |
| kind | string | The kind of the PipeCD application, which should be triggered as a node in deployment chain. The value will be one of: KUBERNETES, TERRAFORM, CLOUDRUN, LAMBDA, ECS. | No |

## EventWatcher

| Field | Type | Description | Required |
|-|-|-|-|
| matcher | [EventWatcherMatcher](#eventwatchermatcher) | Which event will be handled. | Yes |
| handler | [EventWatcherHandler](#eventwatcherhandler) | What to do for the event which matched by the above matcher. | Yes |

### EventWatcherMatcher

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The event name. | Yes |
| labels | map[string]string | Additional attributes of event. This can make an event definition unique even if the one with the same name exists. | No |

### EventWatcherHandler

| Field | Type | Description | Required |
|-|-|-|-|
| type | string | The handler type. Currently, only `GIT_UPDATE` is supported. | Yes |
| config | [EventWatcherHandlerConfig](#eventwatcherhandlerconfig) | Configuration for the event watcher handler. | Yes |

### EventWatcherHandlerConfig

| Field | Type | Description | Required |
|-|-|-|-|
| commitMessage | string | The commit message used to push after replacing values. Default message is used if not given. | No |
| makePullRequest | bool | Whether to create a new branch or not when commit changes in event watcher. Default is `false`. | No |
| replacements | [][EventWatcherReplacement](#eventwatcherreplacement) | List of places where will be replaced when the new event matches. | Yes |

## DriftDetection

| Field | Type | Description | Required |
|-|-|-|-|
| ignoreFields | []string | List of fields path in manifests, which its diff should be ignored. | No |

## PipeCD rich defined types

### Percentage
A wrapper of type `int` to represent percentage data. Basically, you can pass `10` or `"10"` or `10%` and they will be treated as `10%` in PipeCD.

### KubernetesResourcePatch

| Field | Type | Description | Required |
|-|-|-|-|
| target | [KubernetesResourcePatchTarget](#kubernetesresourcepatchtarget) | Which manifest, which field will be the target of patch operations. | Yes |
| ops | [][KubernetesResourcePatchOp](#kubernetesresourcepatchop) | List of operations should be applied to the above target. | No |

### KubernetesResourcePatchTarget

| Field | Type | Description | Required |
|-|-|-|-|
| kind | string | The resource kind. e.g. `ConfigMap` | Yes |
| name | string | The resource name. e.g. `config-map-name` | Yes |
| documentRoot | string | In case you want to manipulate the YAML or JSON data specified in a field of the manfiest, specify that field's path. The string value of that field will be used as input for the patch operations. Otherwise, the whole manifest will be the target of patch operations. e.g. `$.data.envoy-config` | No |

### KubernetesResourcePatchOp

| Field | Type | Description | Required |
|-|-|-|-|
| op | string | The operation type. This must be one of `yaml-replace`, `yaml-add`, `yaml-remove`, `json-replace`, `text-regex`. Default is `yaml-replace`. | No |
| path | string | The path string pointing to the manipulated field. For yaml operations it looks like `$.foo.array[0].bar`. | No |
| value | string | The value string whose content will be used as new value for the field. | No |
