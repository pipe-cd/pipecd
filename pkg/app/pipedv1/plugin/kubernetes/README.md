# Kubernetes plugin

## Overview

Kubernetes plugin supports the Deployment for Kubernetes.

> [!CAUTION] 
> Currently, this is alpha status.

### Quick sync

Quick sync just applies all the defined manifiests to sync the application.

It will be planned in one of the following cases:
- no pipeline was specified in the application configuration file
- `pipeline` was specified but the PR did not make any changes on workload (e.g. Deployment's pod template) or config (e.g. ConfigMap, Secret)

For example, the application configuration as below is missing the pipeline field. This means any pull request touches the application will trigger a quick sync deployment.

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  input:
    helmChart:
      name: helloworld
      path: /path/to/chart
      version: v0.3.0
```

In another case, even when the pipeline was specified, a PR that just changes the Deployment's replicas number for scaling will also trigger a quick sync deployment.

### Pipeline sync

You can configure the pipeline to enable a progressive deployment with a strategy like canary, blue-green.

To enable customization, Kubernetes plugin defines three variants for each application: primary (aka stable), baseline and canary.
- `primary` runs the current version of code and configuration.
- `baseline` runs the same version of code and configuration as the primary variant. (Creating a brand-new baseline workload ensures that the metrics produced are free of any effects caused by long-running processes.)
- `canary` runs the proposed change of code or configuration.

Depending on the configured pipeline, any variants can exist and receive the traffic during the deployment process but once the deployment is completed, only the `primary` variant should be remained.

These are the provided stages for Kubernetes plugin you can use to build your pipeline:

- `K8S_SYNC`
  - sync application to the state specified in the target Git commit without any progressive strategy
- `K8S_PRIMARY_ROLLOUT`
  - update the primary resources to the state defined in the target commit
- `K8S_CANARY_ROLLOUT`
  - generate canary resources based on the definition of the primary resource in the target commit and apply them
- `K8S_CANARY_CLEAN`
  - remove all canary resources
- `K8S_BASELINE_ROLLOUT`
  - generate baseline resources based on the definition of the primary resource in the target commit and apply them
- `K8S_BASELINE_CLEAN`
  - remove all baseline resources

## Plugin Configuration

### Piped Config

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  plugins:
  - name: kubernetes
    port: 7002 # any unused port
    url: file:///path/to/.piped/plugins/kubernetes # or remoteUrl(TBD)
    deployTargets:
      ...
```

| Field | Type | Description | Required |
|-|-|-|-|
| deployTargets | [][DeployTargetConfig](#DeployTargetConfig) | The config for the destinations to deploy applications | Yes |

#### DeployTargetConfig

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The name of the deploy target. | Yes |
| labels | string | The labes of the deploy target. | No |
| config | [KubernetesDeployTargetConfig](#KubernetesDeployTargetConfig) | The configuration of the deploy target for k8s plugin. | No |

##### KubernetesDeployTargetConfig

| Field | Type | Description | Required |
|-|-|-|-|
| masterURL | string | The master URL of the kubernetes cluster. Empty means in-cluster. | No |
| kubectlVersion | string | Version of kubectl which will be used to connect to your cluster. Empty means the [default version](https://github.com/pipe-cd/pipecd/blob/master/pkg/app/pipedv1/plugin/kubernetes/toolregistry/registry.go#L25) will be used. | No |
| kubeConfigPath | string | The path to the kubeconfig file. Empty means in-cluster. | No |

### Application Config

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
...
  plugins:
    kubernetes: # same name as the one defined in `spec.plugins[].name`
      input:
      service:
      workloads:
      variantLabel:
```

| Field | Type | Description | Required |
|-|-|-|-|
| input | [KubernetesDeploymentInput](#kubernetesdeploymentinput) | Input for Kubernetes deployment such as kubectl version, helm version, manifests filter... | No |
| service | [K8sResourceReference](#K8sResourceReference) | Which Kubernetes resource should be considered as the Service of application. Empty means the first Service resource will be used. | No |
| workloads | [][K8sResourceReference](#K8sResourceReference) | Which Kubernetes resources should be considered as the Workloads of application. Empty means all Deployment resources. | No |
| variantLabel | [KubernetesVariantLabel](#kubernetesvariantlabel) | The label will be configured to variant manifests used to distinguish them. | No |

#### KubernetesDeploymentInput

| Field | Type | Description | Required |
|-|-|-|-|
| manifests | []string | List of manifest files in the application directory used to deploy. Empty means all manifest files in the directory will be used. | No |
| kubectlVersion | string | Version of kubectl will be used. Empty means the version set on [piped config](#KubernetesDeployTargetConfig) or [default version](https://github.com/pipe-cd/pipecd/blob/master/pkg/app/pipedv1/plugin/kubernetes/toolregistry/registry.go#L25) will be used. | No |
| kustomizeVersion | string | Version of kustomize will be used. Empty means the [default version](https://github.com/pipe-cd/pipecd/blob/master/pkg/app/pipedv1/plugin/kubernetes/toolregistry/registry.go#L26) will be used. | No |
| kustomizeOptions | map[string]string | List of options that should be used by Kustomize commands. | No |
| helmVersion | string | Version of helm will be used. Empty means the [default version](https://github.com/pipe-cd/pipecd/blob/master/pkg/app/pipedv1/plugin/kubernetes/toolregistry/registry.go#L27) will be used. | No |
| helmChart | [HelmChart](#helmchart) | Where to fetch helm chart. | No |
| helmOptions | [HelmOptions](#helmoptions) | Configurable parameters for helm commands. | No |
| namespace | string | The namespace where manifests will be applied. | No |
| autoCreateNamespace | bool | Automatically create a new namespace if it does not exist. Default is `false`. | No |

##### HelmChart

| Field | Type | Description | Required |
|-|-|-|-|
| path | string | Relative path from the repository root to the chart directory. | No |
| name | string | The chart name. | No |
| version | string | The chart version. | No |

##### HelmOptions

| Field | Type | Description | Required |
|-|-|-|-|
| releaseName | string | The release name of helm deployment. By default, the release name is equal to the application name. | No |
| setValues | map[string]string | List of values. | No |
| valueFiles | []string | List of value files should be loaded. Only local files stored under the application directory or remote files served at the http(s) endpoint are allowed. | No |
| setFiles | map[string]string | List of file path for values. | No |
| apiVersions | []string | Kubernetes api versions used for Capabilities.APIVersions. | No |
| kubeVersion | string | Kubernetes version used for Capabilities.KubeVersion. | No |

#### K8sResourceReference

| Field | Type | Description | Required |
|-|-|-|-|
| kind | string | The kind name of resources. | No |
| name | string | The name of resources. | No |

#### KubernetesVariantLabel

| Field | Type | Description | Required |
|-|-|-|-|
| key | string | The key of the label. Default is `pipecd.dev/variant`. | No |
| primaryValue | string | The label value for PRIMARY variant. Default is `primary`. | No |
| canaryValue | string | The label value for CANARY variant. Default is `canary`. | No |
| baselineValue | string | The label value for BASELINE variant. Default is `baseline`. | No |

### Stage Config

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
...
  pipeline:
    stages:
      - name: K8S_SYNC
        with:
          ...
      - name: K8S_PRIMARY_ROLLOUT
        with:
          ...
      - name: K8S_CANARY_ROLLOUT
        with:
          ...
      - name: K8S_CANARY_CLEAN
        with:
          ...
      - name: K8S_BASELINE_ROLLOUT
        with:
          ...
      - name: K8S_BASELINE_CLEAN
        with:
          ...        
```

#### `K8S_SYNC`

| Field | Type | Description | Required |
|-|-|-|-|
| addVariantLabelToSelector | bool | Whether the PRIMARY variant label should be added to manifests if they were missing. | No |
| prune | string | Whether the resources that are no longer defined in Git should be removed or not. | No |

#### `K8S_PRIMARY_ROLLOUT`

| Field | Type | Description | Required |
|-|-|-|-|
| suffix | string | Suffix that should be used when naming the PRIMARY variant's resources. Default is "primary". | No |
| createService | bool | Whether the PRIMARY service should be created. | No |
| addVariantLabelToSelector | bool | Whether the PRIMARY variant label should be added to manifests if they were missing. | No |
| prune | string | Whether the resources that are no longer defined in Git should be removed or not. | No |

#### `K8S_CANARY_ROLLOUT`

| Field | Type | Description | Required |
|-|-|-|-|
| replicas | int | How many pods for CANARY workloads. Default is `1` pod. Alternatively, can be specified a string suffixed by "%" to indicate a percentage value compared to the pod number of PRIMARY. | No |
| suffix | string | Suffix that should be used when naming the CANARY variant's resources. Default is `canary`. | No |
| createService | bool | Whether the CANARY service should be created. Default is `false`. | No |
| Patches | [][K8sResourcePatch](#K8sResourcePatch) | List of patches used to customize manifests for CANARY variant. | No |

#### K8sResourcePatch

| Field | Type | Description | Required |
|-|-|-|-|
| target | [K8sResourcePatchTarget](#K8sResourcePatchTarget) | Which manifest, which field will be the target of patch operations. | Yes |
| ops | [][K8sResourcePatchOp](#K8sResourcePatchOp) | List of operations should be applied to the above target. | No |

##### K8sResourcePatchTarget

| Field | Type | Description | Required |
|-|-|-|-|
| kind | string | The resource kind. e.g. `ConfigMap` | Yes |
| name | string | The resource name. e.g. `config-map-name` | Yes |
| documentRoot | string | In case you want to manipulate the YAML or JSON data specified in a field of the manfiest, specify that field's path. The string value of that field will be used as input for the patch operations. Otherwise, the whole manifest will be the target of patch operations. e.g. `$.data.envoy-config` | No |

##### K8sResourcePatchOp

| Field | Type | Description | Required |
|-|-|-|-|
| op | string | The operation type. This must be one of `yaml-replace`, `yaml-add`, `yaml-remove`, `json-replace`, `text-regex`. Default is `yaml-replace`. | No |
| path | string | The path string pointing to the manipulated field. For yaml operations it looks like `$.foo.array[0].bar`. | No |
| value | string | The value string whose content will be used as new value for the field. | No |


#### `K8S_CANARY_CLEAN`

| Field | Type | Description | Required |
|-|-|-|-|

#### `K8S_BASELINE_ROLLOUT`

| Field | Type | Description | Required |
|-|-|-|-|
| replicas | int | How many pods for BASELINE workloads. Default is `1` pod. Alternatively, can be specified a string suffixed by "%" to indicate a percentage value compared to the pod number of PRIMARY | No |
| suffix | string | Suffix that should be used when naming the BASELINE variant's resources. Default is `baseline`. | No |
| createService | bool | Whether the BASELINE service should be created. Default is `false`. | No |

#### `K8S_BASELINE_CLEAN`

| Field | Type | Description | Required |
|-|-|-|-|
