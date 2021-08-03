---
title: "Configuration reference"
linkTitle: "Configuration reference"
weight: 22
description: >
  This page describes all configurable fields in the deployment configuration and analysis template.
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
| input | [KubernetesDeploymentInput](/docs/user-guide/configuration-reference/#kubernetesdeploymentinput) | Input for Kubernetes deployment such as kubectl version, helm version, manifests filter... | No |
| commitMatcher | [CommitMatcher](/docs/user-guide/configuration-reference/#commitmatcher) | Forcibly use QuickSync or Pipeline when commit message matched the specified pattern. | No |
| quickSync | [KubernetesQuickSync](/docs/user-guide/configuration-reference/#kubernetesquicksync) | Configuration for quick sync. | No |
| pipeline | [Pipeline](/docs/user-guide/configuration-reference/#pipeline) | Pipeline for deploying progressively. | No |
| service | [KubernetesService](/docs/user-guide/configuration-reference/#kubernetesservice) | Which Kubernetes resource should be considered as the Service of application. Empty means the first Service resource will be used. | No |
| workloads | [][KubernetesWorkload](/docs/user-guide/configuration-reference/#kubernetesworkload) | Which Kubernetes resources should be considered as the Workloads of application. Empty means all Deployment resources. | No |
| trafficRouting | [KubernetesTrafficRouting](/docs/user-guide/configuration-reference/#kubernetestrafficrouting) | How to change traffic routing percentages. | No |
| sealedSecrets | [][SealedSecretMapping](/docs/user-guide/configuration-reference/#sealedsecretmapping) | The list of sealed secrets should be decrypted. | No |
| triggerPaths | []string | List of directories or files where their changes will trigger the deployment. Regular expression can be used. | No |
| timeout | duration | The maximum length of time to execute deployment before giving up. Default is 6h. | No |

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
| input | [TerraformDeploymentInput](/docs/user-guide/configuration-reference/#terraformdeploymentinput) | Input for Terraform deployment such as terraform version, workspace... | No |
| quickSync | [TerraformQuickSync](/docs/user-guide/configuration-reference/#terraformquicksync) | Configuration for quick sync. | No |
| pipeline | [Pipeline](/docs/user-guide/configuration-reference/#pipeline) | Pipeline for deploying progressively. | No |
| sealedSecrets | [][SealedSecretMapping](/docs/user-guide/configuration-reference/#sealedsecretmapping) | The list of sealed secrets should be decrypted. | No |
| triggerPaths | []string | List of directories or files where their changes will trigger the deployment. Regular expression can be used. | No |
| timeout | duration | The maximum length of time to execute deployment before giving up. Default is 6h. | No |

## CloudRun application

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
| input | [CloudRunDeploymentInput](/docs/user-guide/configuration-reference/#cloudrundeploymentinput) | Input for CloudRun deployment such as docker image... | No |
| quickSync | [CloudRunQuickSync](/docs/user-guide/configuration-reference/#cloudrunquicksync) | Configuration for quick sync. | No |
| pipeline | [Pipeline](/docs/user-guide/configuration-reference/#pipeline) | Pipeline for deploying progressively. | No |
| triggerPaths | []string | List of directories or files where their changes will trigger the deployment. Regular expression can be used. | No |
| sealedSecrets | [][SealedSecretMapping](/docs/user-guide/configuration-reference/#sealedsecretmapping) | The list of sealed secrets should be decrypted. | No |
| timeout | duration | The maximum length of time to execute deployment before giving up. Default is 6h. | No |

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
| quickSync | [LambdaQuickSync](/docs/user-guide/configuration-reference/#lambdaquicksync) | Configuration for quick sync. | No |
| pipeline | [Pipeline](/docs/user-guide/configuration-reference/#pipeline) | Pipeline for deploying progressively. | No |
| triggerPaths | []string | List of directories or files where their changes will trigger the deployment. Regular expression can be used. | No |
| sealedSecrets | [][SealedSecretMapping](/docs/user-guide/configuration-reference/#sealedsecretmapping) | The list of sealed secrets should be decrypted. | No |
| timeout | duration | The maximum length of time to execute deployment before giving up. Default is 6h. | No |

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
| input | [ECSDeploymentInput](#ecsdeploymentinput) | Input for ECS deployment such as TaskDefinition, Service... | Yes |
| quickSync | [ECSQuickSync](/docs/user-guide/configuration-reference/#ecsquicksync) | Configuration for quick sync. | No |
| pipeline | [Pipeline](/docs/user-guide/configuration-reference/#pipeline) | Pipeline for deploying progressively. | No |
| triggerPaths | []string | List of directories or files where their changes will trigger the deployment. Regular expression can be used. | No |
| sealedSecrets | [][SealedSecretMapping](/docs/user-guide/configuration-reference/#sealedsecretmapping) | The list of sealed secrets should be decrypted. | No |
| timeout | duration | The maximum length of time to execute deployment before giving up. Default is 6h. | No |

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
| metrics | map[string][AnalysisMetrics](/docs/user-guide/configuration-reference/#analysismetrics) | Template for metrics. | No |

## Event Watcher Configuration

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
| replacements | [][EventWatcherReplacement](/docs/user-guide/configuration-reference/#eventwatcherreplacement) | List of places where will be replaced when the new event matches. | Yes |

## EventWatcherReplacement

| Field | Type | Description | Required |
|-|-|-|-|
| file | string | The path to the file to be updated. | Yes |
| yamlField | string | The yaml path to the field to be updated. It requires to start with `$` which represents the root element. e.g. `$.foo.bar[0].baz`. | Yes |

## CommitMatcher

| Field | Type | Description | Required |
|-|-|-|-|
| quickSync | string | Regular expression string to forcibly do QuickSync when it matches the commit message. | No |
| pipeline | string | Regular expression string to forcibly do Pipeline when it matches the commit message. | No |

## SealedSecretMapping

| Field | Type | Description | Required |
|-|-|-|-|
| path | string | Relative path from the application directory to the sealed secret file. | Yes |
| outFilename | string | The filename for the decrypted secret. Empty means the same name with the sealed secret file. | No |
| outDir | string | The directory name where to put the decrypted secret. Empty means the same directory with the sealed secret file. | No |

## Pipeline

| Field | Type | Description | Required |
|-|-|-|-|
| stages | [][PipelineStage](/docs/user-guide/configuration-reference/#pipelinestage) | List of deployment pipeline stages. | No |

## PipelineStage

| Field | Type | Description | Required |
|-|-|-|-|
| id | string | The unique ID of the stage. | No |
| name | string | One of the provided stage names. | Yes |
| desc | string | The description about the stage. | No |
| timeout | duration | The maximum time the stage can be taken to run. | No |
| with | [StageOptions](/docs/user-guide/configuration-reference/#stageoptions) | Specific configuration for the stage. This must be one of these [StageOptions](/docs/user-guide/configuration-reference/#stageoptions). | No |

## KubernetesDeploymentInput

| Field | Type | Description | Required |
|-|-|-|-|
| manifests | []string | List of manifest files in the application directory used to deploy. Empty means all manifest files in the directory will be used. | No |
| kubectlVersion | string | Version of kubectl will be used. Empty means the [default version](https://github.com/pipe-cd/pipe/blob/master/dockers/piped-base/install-kubectl.sh#L34) will be used. | No |
| kustomizeVersion | string | Version of kustomize will be used. Empty means the [default version](https://github.com/pipe-cd/pipe/blob/master/dockers/piped-base/install-kustomize.sh#L34) will be used. | No |
| kustomizeOptions | map[string]string | List of options that should be used by Kustomize commands. | No |
| helmVersion | string | Version of helm will be used. Empty means the [default version](https://github.com/pipe-cd/pipe/blob/master/dockers/piped-base/install-helm.sh#L35) will be used. | No |
| helmChart | [HelmChart](/docs/user-guide/configuration-reference/#helmchart) | Where to fetch helm chart. | No |
| helmOptions | [HelmOptions](/docs/user-guide/configuration-reference/#helmoptions) | Configurable parameters for helm commands. | No |
| namespace | string | The namespace where manifests will be applied. | No |
| autoRollback | bool | Automatically reverts all deployment changes on failure. Default is `true`. | No |

## HelmChart

| Field | Type | Description | Required |
|-|-|-|-|
| gitRemote | string | Git remote address where the chart is placing. Empty means the same repository. | No |
| ref | string | The commit SHA or tag value. Only valid when gitRemote is not empty. | No |
| path | string | Relative path from the repository root to the chart directory. | No |
| repository | string | The name of a registered Helm Chart Repository. | No |
| name | string | The chart name. | No |
| version | string | The chart version. | No |

## HelmOptions

| Field | Type | Description | Required |
|-|-|-|-|
| releaseName | string | The release name of helm deployment. By default, the release name is equal to the application name. | No |
| valueFiles | []string | List of value files should be loaded. | No |
| setFiles | map[string]string | List of file path for values. | No |

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
| istio | [IstioTrafficRouting](/docs/user-guide/configuration-reference/#istiotrafficrouting)| Istio configuration when the method is `istio`. | No |

## IstioTrafficRouting

| Field | Type | Description | Required |
|-|-|-|-|
| editableRoutes | []string | List of routes in the VirtualService that can be changed to update traffic routing. Empty means all routes should be updated. | No |
| host | string | The service host. | No |
| virtualService | [IstioVirtualService](/docs/user-guide/configuration-reference/#istiovirtualservice) | The reference to VirtualService manifest. Empty means the first VirtualService resource will be used. | No |

## IstioVirtualService

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
| autoRollback | bool | Automatically reverts all changes from all stages when one of them failed. | No |

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

## LambdaQuickSync

| Field | Type | Description | Required |
|-|-|-|-|

## ECSDeploymentInput

| Field | Type | Description | Required |
|-|-|-|-|
| serviceDefinitionFile | string | The path ECS Service configuration file. Allow file in both `yaml` and `json` format. The default value is `service.json`. | No |
| taskDefinitionFile | string | The path to ECS TaskDefinition configuration file. Allow file in both `yaml` and `json` format. The default value is `taskdef.json`. | No |
| targetGroups | [ECSTargetGroupInput](#ecstargetgroupinput) | The target groups configuration, will be used to routing traffic to created task sets. | Yes |

### ECSTargetGroupInput

| Field | Type | Description | Required |
|-|-|-|-|
| primary | ECSTargetGroupObject | The PRIMARY target group, will be used to register the PRIMARY ECS task set. | Yes |
| canary | ECSTargetGroupObject | The CANARY target group, will be used to register the CANARY ECS task set if exist. It's required to enable PipeCD to perform the multi-stage deployment. | No |

Note: You can get examples for those object from [here](/docs/examples/#ecs-applications).

## ECSQuickSync

| Field | Type | Description | Required |
|-|-|-|-|

## AnalysisMetrics

| Field | Type | Description | Required |
|-|-|-|-|
| provider | string | The unique name of provider defined in the Piped Configuration. | Yes |
| query | string | A query performed against the [Analysis Provider](/docs/concepts/#analysis-provider). | Yes |
| expected | [AnalysisExpected](/docs/user-guide/configuration-reference/#analysisexpected) | The expected query result. | Yes |
| interval | duration | Run a query at specified intervals. | Yes |
| failureLimit | int | Acceptable number of failures. e.g. If 1 is set, the `ANALYSIS` stage will end with failure after two queries results failed. Defaults to 1. | No |
| skipOnNoData | bool | If true, it considers as a success when no data returned from the analysis provider. Defaults to false. | No |
| timeout | duration | How long after which the query times out. | No |
| template | [AnalysisTemplateRef](/docs/user-guide/configuration-reference/#analysistemplateref) | Reference to the template to be used. | No |


## AnalysisLog

| Field | Type | Description | Required |
|-|-|-|-|

## AnalysisHttp

| Field | Type | Description | Required |
|-|-|-|-|

## AnalysisExpected

| Field | Type | Description | Required |
|-|-|-|-|
| min | float64 | Failure, if the query result is less than this value. | No |
| max | float64 | Failure, if the query result is larger than this value. | No |

## AnalysisTemplateRef

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The template name to refer. | Yes |
| args | map[string]string | The arguments for custom-args. | No |

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
| patches | [][KubernetesResourcePatch](/docs/user-guide/configuration-reference/#kubernetesresourcepatch) | List of patches used to customize manifests for CANARY variant. | No |

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
This stage routes traffic with the method specified in [KubernetesTrafficRouting](https://pipecd.dev/docs/user-guide/configuration-reference/#kubernetestrafficrouting).
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
| metrics | [][AnalysisMetrics](/docs/user-guide/configuration-reference/#analysismetrics) | Configuration for analysis by metrics. | No |

## PipeCD rich defined types

### Percentage
A wrapper of type `int` to represent percentage data. Basically, you can pass `10` or `"10"` or `10%` and they will be treated as `10%` in PipeCD.

### KubernetesResourcePatch

| Field | Type | Description | Required |
|-|-|-|-|
| target | [KubernetesResourcePatchTarget](/docs/user-guide/configuration-reference/#kubernetesresourcepatchtarget) | Which manifest, which field will be the target of patch operations. | Yes |
| yamlOps | [][KubernetesResourcePatchYAMLOp](/docs/user-guide/configuration-reference/#kubernetesresourcepatchyamlop) | List of yaml operations should be applied to the above target. | No |

### KubernetesResourcePatchTarget

| Field | Type | Description | Required |
|-|-|-|-|
| kind | string | The resource kind. e.g. `ConfigMap` | Yes |
| name | string | The resource name. e.g. `config-map-name` | Yes |
| field | string | A string field whose content will be the target of patch operations. Empty means the whole manifest will be the target. e.g. `$.data.envoy-config` | No |

### KubernetesResourcePatchYAMLOp

| Field | Type | Description | Required |
|-|-|-|-|
| op | string | The operation type. This must be one of `replace`, `add`, `remove`. Default is `replace`. | No |
| path | string | The path string pointing to the manipulated field. e.g. `$.foo.array[0].bar` | Yes |
| value | string | The value string whose content will be used as new value for the field. | No |
