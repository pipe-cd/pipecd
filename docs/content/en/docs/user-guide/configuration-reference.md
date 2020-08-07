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
| commitMatcher | [CommitMatcher](/docs/user-guide/configuration-reference/#commitmatcher) | Forcibly use QuickSync or Pipeline when commit message matched the specified format. | No |
| quickSync | [KubernetesQuickSync](/docs/user-guide/configuration-reference/#kubernetesquicksync) | Configuration for quick sync. | No |
| pipeline | [Pipeline](/docs/user-guide/configuration-reference/#pipeline) | Pipeline for deploying progressively. | No |
| service | [KubernetesService](/docs/user-guide/configuration-reference/#kubernetesservice) | Which Kubernetes resource should be considered as the Service of application. Empty means the first Service resource will be used. | No |
| workloads | [][KubernetesWorkload](/docs/user-guide/configuration-reference/#kubernetesworkload) | Which Kubernetes resources should be consiedered as the Workloads of application. Empty means all Deployment resources. | No |
| trafficRouting | [KubernetesTrafficRouting](/docs/user-guide/configuration-reference/#kubernetestrafficrouting) | How to change traffic routing percentages. | No |

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

## Lambda application

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: LambdaApp
spec:
  input:
  pipeline:
  ...
```

| Field | Type | Description | Required |
|-|-|-|-|
| input | [CloudRunDeploymentInput](/docs/user-guide/configuration-reference/#cloudrundeploymentinput) | Input for Lambda deployment such as where to fetch source code... | No |
| quickSync | [CloudRunQuickSync](/docs/user-guide/configuration-reference/#cloudrunquicksync) | Configuration for quick sync. | No |
| pipeline | [Pipeline](/docs/user-guide/configuration-reference/#pipeline) | Pipeline for deploying progressively. | No |

## Analysis Template Configuration

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: AnalysisTemplate
spec:
  metrics:
```

| Field | Type | Description | Required |
|-|-|-|-|

## CommitMatcher

| Field | Type | Description | Required |
|-|-|-|-|
| quickSync | string | Regular expression string to forcibly do QuickSync when it matches the commit message. | No |
| pipeline | string | Regular expression string to forcibly do Pipeline when it matches the commit message. | No |

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
| manifests | []string | List of manifest files in the application configuration directory used to deploy. Empty means all manifest files in the directory will be used. | No |
| kubectlVersion | string | Version of kubectl will be used. Empty means the [default version](https://github.com/pipe-cd/pipe/blob/master/dockers/piped-base/install-kubectl.sh#L34) will be used. | No |
| kustomizeVersion | string | Version of kustomize will be used. Empty means the [default version](https://github.com/pipe-cd/pipe/blob/master/dockers/piped-base/install-kustomize.sh#L34) will be used. | No |
| helmVersion | string | Version of helm will be used. Empty means the [default version](https://github.com/pipe-cd/pipe/blob/master/dockers/piped-base/install-helm.sh#L35) will be used. | No |
| helmChart | [HelmChart](/docs/user-guide/configuration-reference/#helmchart) | Where to fetch helm chart. | No |
| helmOptions | [HelmOptions](/docs/user-guide/configuration-reference/#helmoptions) | Configurable parameters for helm commands. | No |
| namespace | string | The namespace where manifests will be applied. | No |
| autoRollback | bool | Automatically reverts all deployment changes on failure. Default is `true`. | No |
| dependencies | []string | List of directories where their changes will trigger the deployment. | No |

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
| releaseName | string | The release name of helm deployment. By default the release name is equal to the application name. | No |
| valueFiles | []string | List of value files should be loaded. | No |
| setFiles | []string | List of file path for values. | No |

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
| method | string | Which traffic routing method will be used. Avaiable values are `istio`, `smi`, `podselector`. Default is `podselector`. | No |
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

## TerraformQuickSync

| Field | Type | Description | Required |
|-|-|-|-|

## CloudRunDeploymentInput

| Field | Type | Description | Required |
|-|-|-|-|

## CloudRunQuickSync

| Field | Type | Description | Required |
|-|-|-|-|

## LambdaDeploymentInput

| Field | Type | Description | Required |
|-|-|-|-|

## LambdaQuickSync

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
| replicas | Replicas | How many pods for CANARY workloads. Default is `1` pod. | No |
| suffix | string | Suffix that should be used when naming the CANARY variant's resources. Default is `canary`. | No |
| createService | bool | Whether the CANARY service should be created. Default is `false`. | No |

### KubernetesCanaryCleanStageOptions

| Field | Type | Description | Required |
|-|-|-|-|
| | | | |

### KubernetesBaselineRolloutStageOptions

| Field | Type | Description | Required |
|-|-|-|-|
| replicas | Replicas | How many pods for BASELINE workloads. Default is `1` pod. | No |
| suffix | string | Suffix that should be used when naming the BASELINE variant's resources. Default is `baseline`. | No |
| createService | bool | Whether the BASELINE service should be created. Default is `false`. | No |

### KubernetesBaselineCleanStageOptions

| Field | Type | Description | Required |
|-|-|-|-|
| | | | |

### KubernetesTrafficRoutingStageOptions

| Field | Type | Description | Required |
|-|-|-|-|
| all | string | Which variant should receive all traffic. Avaiable values are "primary", "canary", "baseline". Default is `primary`. | No |
| primary | int | The percentage of traffic should be routed to PRIMARY variant. | No |
| canary | int | The percentage of traffic should be routed to CANARY variant. | No |
| baseline | int | The percentage of traffic should be routed to BASELINE variant. | No |

### TerraformPlanStageOptions

| Field | Type | Description | Required |
|-|-|-|-|

### TerraformApplyStageOptions

| Field | Type | Description | Required |
|-|-|-|-|

### CloudRunCanaryRolloutStageOptions

| Field | Type | Description | Required |
|-|-|-|-|

### CloudRunTrafficRoutingStageOptions

| Field | Type | Description | Required |
|-|-|-|-|

### LambdaCanaryRolloutStageOptions

| Field | Type | Description | Required |
|-|-|-|-|

### LambdaTrafficRoutingStageOptions

| Field | Type | Description | Required |
|-|-|-|-|
