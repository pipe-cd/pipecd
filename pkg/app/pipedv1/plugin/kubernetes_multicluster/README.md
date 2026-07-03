# Kubernetes Multicluster Plugin

## Overview

The Kubernetes Multicluster plugin supports deploying applications across multiple Kubernetes clusters from a single pipeline definition. It handles sync state, traffic routing, and rollback across all target clusters.

> [!CAUTION]
> The configuration format is unstable and may change in the future.

### Quick sync

Quick sync applies manifests to all target clusters immediately, without a pipeline.

It will be planned in one of the following cases:
- no pipeline was specified in the application configuration file
- the sync was triggered manually with the quick sync option

### Pipeline sync

You can configure a pipeline to deploy progressively across clusters with stages like canary rollout, traffic routing, primary rollout, and baseline comparison.

The plugin defines three variants for each application: primary, canary, and baseline.
- `primary` runs the current stable version.
- `canary` runs the proposed change and receives a configurable percentage of traffic.
- `baseline` runs the same version as primary, used as a comparison target for canary analysis.

Once deployment completes, only the `primary` variant remains across all clusters.

These are the stages provided by the Kubernetes Multicluster plugin:

- `K8S_MULTI_SYNC`: sync all target clusters with the current manifests (used for quick sync)
- `K8S_MULTI_PRIMARY_ROLLOUT`: roll out the new version as the PRIMARY variant to all target clusters
- `K8S_MULTI_CANARY_ROLLOUT`: roll out the new version as the CANARY variant to all target clusters
- `K8S_MULTI_CANARY_CLEAN`: remove CANARY variant resources from all target clusters
- `K8S_MULTI_BASELINE_ROLLOUT`: roll out the current version as the BASELINE variant to all target clusters
- `K8S_MULTI_BASELINE_CLEAN`: remove BASELINE variant resources from all target clusters
- `K8S_MULTI_TRAFFIC_ROUTING`: route traffic between PRIMARY, CANARY, and BASELINE variants across all target clusters
- `K8S_MULTI_ROLLBACK`: rollback all target clusters to the previous state (automatically added when `autoRollback` is enabled)

## Directory Structure

```
kubernetes_multicluster/
├── main.go         # Plugin entrypoint
├── config/         # Configuration types for piped, deploy targets, application, and stages
├── deployment/     # Stage execution and pipeline planning
├── provider/       # Kubernetes API client, manifest loading, kubectl/kustomize/helm wrappers
├── livestate/      # Live state fetcher (reports current resource state across clusters)
├── planpreview/    # Plan preview (shows diff before deployment)
├── toolregistry/   # Tool registry for kubectl, kustomize, and helm binaries
└── example/        # Example application configurations
```

## How to Build & Run

### Build

Build the plugin binary using the `build/plugin` Makefile target from the repository root:

```bash
# Build all plugins (including kubernetes_multicluster)
make build/plugin

# Build only the kubernetes_multicluster plugin
make build/plugin PLUGINS=kubernetes_multicluster
```

The binary will be placed at `~/.piped/plugins/kubernetes_multicluster` by default.

### Run

Configure your piped to load the plugin by adding it to the piped config:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  plugins:
  - name: kubernetes_multicluster
    port: 7002 # any unused port
    url: file:///home/<user>/.piped/plugins/kubernetes_multicluster
    deployTargets:
      - name: cluster1
        config:
          masterURL: https://127.0.0.1:61337
          kubeConfigPath: /path/to/kubeconfig/for/cluster1
      - name: cluster2
        config:
          masterURL: https://127.0.0.1:62082
          kubeConfigPath: /path/to/kubeconfig/for/cluster2
```

Then start piped as usual. The plugin process will be launched automatically on the specified port.

> [!NOTE]
> If you build the plugin manually with `go build` and place the binary at a custom path, piped may use a previously cached binary instead. Pass `--force-plugin-redownload` when starting piped to ensure it picks up your locally built binary:
> ```bash
> piped --force-plugin-redownload ...
> ```

### Prepare clusters locally

To test locally with [kind](https://kind.sigs.k8s.io/):

```sh
kind create cluster --name cluster1
kind export kubeconfig --name cluster1 --kubeconfig /path/to/kubeconfig/for/cluster1

kind create cluster --name cluster2
kind export kubeconfig --name cluster2 --kubeconfig /path/to/kubeconfig/for/cluster2
```

Refer to [cmd/pipecd/README.md](../../../../../cmd/pipecd/README.md) to set up the PipeCD control plane, and [cmd/pipedv1/README.md](../../../../../cmd/pipedv1/README.md) to run piped locally.

## How to Test

### Unit tests

Run unit tests from the repository root:

```bash
go test ./pkg/app/pipedv1/plugin/kubernetes_multicluster/...
```

Unit tests use [envtest](https://book.kubebuilder.io/reference/envtest) to spin up a real Kubernetes API server in memory. No external cluster is needed.

### Integration tests

To test against real clusters, build and run piped with the plugin pointing to real clusters. The plugin resolves kubeconfig in this order:

1. `kubeConfigPath` in deploy target config: explicit kubeconfig file path
2. `masterURL` in deploy target config: explicit API server URL
3. In-cluster config: when running inside a pod

The simplest approach for local development is to use kind clusters as shown above.

## Examples

There are example application configurations under `./example`:

| Name | Description |
|------|-------------|
| [simple](./example/simple/) | Deploy the same resources to multiple clusters. |
| [multi-sources-template-none](./example/multi-sources-template-none/) | Deploy different resources to each cluster. |

## Plugin Configuration

### Piped Config

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  plugins:
  - name: kubernetes_multicluster
    port: 7002
    url: file:///path/to/.piped/plugins/kubernetes_multicluster
    deployTargets:
      - name: cluster1
        config:
          masterURL: https://127.0.0.1:61337
          kubeConfigPath: /path/to/kubeconfig/for/cluster1
          kubectlVersion: 1.32.0
      - name: cluster2
        config:
          masterURL: https://127.0.0.1:62082
          kubeConfigPath: /path/to/kubeconfig/for/cluster2
```

| Field | Type | Description | Required |
|-|-|-|-|
| deployTargets | [][DeployTargetConfig](#deploytargetconfig) | List of Kubernetes clusters to deploy to | Yes |

#### DeployTargetConfig

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | Name of the deploy target. | Yes |
| labels | map[string]string | Labels for the deploy target. | No |
| config | [KubernetesDeployTargetConfig](#kubernetesdeploytargetconfig) | Cluster-specific configuration. | No |

##### KubernetesDeployTargetConfig

| Field | Type | Description | Required |
|-|-|-|-|
| masterURL | string | The API server URL of the Kubernetes cluster. Empty means in-cluster. | No |
| kubeConfigPath | string | Path to the kubeconfig file. Empty means in-cluster. | No |
| kubectlVersion | string | Version of kubectl to use. Piped downloads it automatically via the tool registry. | No |
| appStateInformer | [KubernetesAppStateInformer](#kubernetesappstateinformer) | Configuration for the resource state watcher. | No |

##### KubernetesAppStateInformer

| Field | Type | Description | Required |
|-|-|-|-|
| namespace | string | Only watch this namespace. Empty means watch all namespaces. | No |
| includeResources | [][KubernetesResourceMatcher](#kubernetesresourcematcher) | Additional resource types to watch. | No |
| excludeResources | [][KubernetesResourceMatcher](#kubernetesresourcematcher) | Resource types to exclude from watching. | No |

##### KubernetesResourceMatcher

| Field | Type | Description | Required |
|-|-|-|-|
| apiVersion | string | The API version of the resource (e.g. `apps/v1`). | No |
| kind | string | The kind of the resource. Empty means all kinds. | No |

### Application Config

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  plugins:
    kubernetes_multicluster:
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
| input | [KubernetesDeploymentInput](#kubernetesdeploymentinput) | Input for deployment such as manifests, kubectl version, helm/kustomize options. | No |
| quickSync | [K8sSyncStageOptions](#k8ssyncstageoptions) | Options for the quick sync stage. | No |
| workloads | [][K8sResourceReference](#k8sresourcereference) | Which resources are treated as the application workload. Empty means all Deployments. | No |
| service | [K8sResourceReference](#k8sresourcereference) | Which resource is treated as the Service. Empty means the first Service found. | No |
| variantLabel | [KubernetesVariantLabel](#kubernetesvariantlabel) | The label used to distinguish variant resources. | No |
| trafficRouting | [KubernetesTrafficRouting](#kubernetestrafficrouting) | Traffic routing method. Default is PodSelector. | No |

#### KubernetesDeploymentInput

| Field | Type | Description | Required |
|-|-|-|-|
| manifests | []string | Manifest files to deploy. Empty means all files in the app directory. | No |
| kubectlVersion | string | Version of kubectl to use. | No |
| kustomizeVersion | string | Version of kustomize to use. | No |
| kustomizeOptions | map[string]string | Options to pass to kustomize commands. | No |
| helmVersion | string | Version of helm to use. | No |
| helmChart | [InputHelmChart](#inputhelmchart) | Where to fetch the Helm chart from. | No |
| helmOptions | [InputHelmOptions](#inputhelmoptions) | Options to pass to helm commands. | No |
| namespace | string | Namespace to apply manifests to. | No |
| autoCreateNamespace | bool | Automatically create the namespace if it does not exist. Default is `false`. | No |
| multiTargets | [][KubernetesMultiTarget](#kubernetesmultitarget) | Per-cluster overrides for manifests and tool versions. | No |

#### KubernetesMultiTarget

Allows setting per-cluster manifest overrides within a single deployment input.

| Field | Type | Description | Required |
|-|-|-|-|
| target | [KubernetesMultiTargetDeployTarget](#kubernetesmultitargetdeploytarget) | Identifies the deploy target this override applies to. | Yes |
| manifests | []string | Manifest files to use for this cluster. Empty means use the top-level manifests. | No |
| kubectlVersion | string | kubectl version override for this cluster. | No |
| kustomizeDir | string | Kustomize directory override for this cluster. | No |
| kustomizeVersion | string | Kustomize version override for this cluster. | No |
| kustomizeOptions | map[string]string | Kustomize options override for this cluster. | No |

#### KubernetesMultiTargetDeployTarget

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | Name of the deploy target. | Yes |
| labels | map[string]string | Labels to match against deploy target labels. | No |

#### InputHelmChart

| Field | Type | Description | Required |
|-|-|-|-|
| path | string | Path to the local chart directory (relative to the app directory). | No |
| gitRemote | string | Git remote URL of the chart repository. | No |
| ref | string | Commit SHA or tag for the remote git chart. | No |
| repository | string | Name of an added Helm chart repository. | No |
| name | string | Name of the chart in the repository. | No |
| version | string | Version of the chart. | No |

#### InputHelmOptions

| Field | Type | Description | Required |
|-|-|-|-|
| releaseName | string | Helm release name. Default is the application name. | No |
| setValues | map[string]string | Values to set with `--set`. | No |
| valueFiles | []string | Value files to load with `-f`. | No |
| setFiles | map[string]string | File paths for values with `--set-file`. | No |
| apiVersions | []string | Kubernetes API versions for `helm template --api-versions`. | No |
| kubeVersion | string | Kubernetes version for `Capabilities.KubeVersion`. | No |

#### KubernetesVariantLabel

| Field | Type | Description | Required |
|-|-|-|-|
| key | string | Label key. Default is `pipecd.dev/variant`. | No |
| primaryValue | string | Label value for the PRIMARY variant. Default is `primary`. | No |
| canaryValue | string | Label value for the CANARY variant. Default is `canary`. | No |
| baselineValue | string | Label value for the BASELINE variant. Default is `baseline`. | No |

#### KubernetesTrafficRouting

| Field | Type | Description | Required |
|-|-|-|-|
| method | string | Traffic routing method. `podselector` or `istio`. Default is `podselector`. | No |
| istio | [IstioTrafficRouting](#istiotrafficrouting) | Istio-specific configuration. Required when method is `istio`. | No |

#### IstioTrafficRouting

| Field | Type | Description | Required |
|-|-|-|-|
| editableRoutes | []string | Routes in the VirtualService that can be modified. Empty means all routes. | No |
| host | string | The service host. | No |
| virtualService | [K8sResourceReference](#k8sresourcereference) | Reference to the VirtualService resource. Empty means the first VirtualService found. | No |

#### K8sResourceReference

| Field | Type | Description | Required |
|-|-|-|-|
| kind | string | Kind of the Kubernetes resource. | No |
| name | string | Name of the Kubernetes resource. | No |

### Stage Config

```yaml
pipeline:
  stages:
    - name: K8S_MULTI_CANARY_ROLLOUT
      with:
        replicas: 10%
    - name: K8S_MULTI_TRAFFIC_ROUTING
      with:
        canary: 10
    - name: K8S_MULTI_PRIMARY_ROLLOUT
    - name: K8S_MULTI_CANARY_CLEAN
```

#### `K8S_MULTI_SYNC`

| Field | Type | Description | Required |
|-|-|-|-|
| addVariantLabelToSelector | bool | Add the PRIMARY variant label to manifest selectors if missing. Default is `false`. | No |
| prune | bool | Remove resources no longer defined in Git. Default is `false`. | No |
| multiTargets | []string | Limit this stage to a subset of deploy targets by name. Empty means all targets. | No |

#### `K8S_MULTI_PRIMARY_ROLLOUT`

| Field | Type | Description | Required |
|-|-|-|-|
| suffix | string | Suffix for PRIMARY variant resource names. Default is `primary`. | No |
| createService | bool | Create a Service for the PRIMARY variant. Default is `false`. | No |
| addVariantLabelToSelector | bool | Add the PRIMARY variant label to selectors if missing. Default is `false`. | No |
| prune | bool | Remove resources no longer defined in Git. Default is `false`. | No |
| multiTargets | []string | Limit this stage to a subset of deploy targets by name. Empty means all targets. | No |

#### `K8S_MULTI_CANARY_ROLLOUT`

| Field | Type | Description | Required |
|-|-|-|-|
| replicas | int or string | Number of canary pods. Integer for absolute count, string with `%` for percentage of primary. Default is `1`. | No |
| suffix | string | Suffix for CANARY variant resource names. Default is `canary`. | No |
| createService | bool | Create a Service for the CANARY variant. Default is `false`. | No |
| patches | [][K8sResourcePatch](#k8sresourcepatch) | Patches to customize manifests for the CANARY variant. | No |
| multiTargets | []string | Limit this stage to a subset of deploy targets by name. Empty means all targets. | No |

#### `K8S_MULTI_CANARY_CLEAN`

| Field | Type | Description | Required |
|-|-|-|-|
| multiTargets | []string | Limit this stage to a subset of deploy targets by name. Empty means all targets. | No |

#### `K8S_MULTI_BASELINE_ROLLOUT`

| Field | Type | Description | Required |
|-|-|-|-|
| replicas | int or string | Number of baseline pods. Integer for absolute count, string with `%` for percentage of primary. Default is `1`. | No |
| suffix | string | Suffix for BASELINE variant resource names. Default is `baseline`. | No |
| createService | bool | Create a Service for the BASELINE variant. Default is `false`. | No |
| multiTargets | []string | Limit this stage to a subset of deploy targets by name. Empty means all targets. | No |

#### `K8S_MULTI_BASELINE_CLEAN`

| Field | Type | Description | Required |
|-|-|-|-|
| multiTargets | []string | Limit this stage to a subset of deploy targets by name. Empty means all targets. | No |

#### `K8S_MULTI_TRAFFIC_ROUTING`

| Field | Type | Description | Required |
|-|-|-|-|
| all | string | Route all traffic to one variant. `primary`, `canary`, or `baseline`. | No |
| primary | int | Percentage of traffic to route to the PRIMARY variant. | No |
| canary | int | Percentage of traffic to route to the CANARY variant. | No |
| baseline | int | Percentage of traffic to route to the BASELINE variant. | No |
| multiTargets | []string | Limit this stage to a subset of deploy targets by name. Empty means all targets. | No |

#### `K8S_MULTI_ROLLBACK`

No configuration options. This stage is automatically appended to the pipeline when `autoRollback` is enabled. It restores all target clusters to the state before the deployment started.

#### K8sResourcePatch

| Field | Type | Description | Required |
|-|-|-|-|
| target | [K8sResourcePatchTarget](#k8sresourcepatchtarget) | The resource to patch. | Yes |
| ops | [][K8sResourcePatchOp](#k8sresourcepatchop) | The patch operations to apply. | Yes |

#### K8sResourcePatchTarget

| Field | Type | Description | Required |
|-|-|-|-|
| kind | string | Kind of the resource to patch. | No |
| name | string | Name of the resource to patch. | No |
| documentRoot | string | Field path within the resource to use as the patch target. Empty means the whole manifest. | No |

#### K8sResourcePatchOp

| Field | Type | Description | Required |
|-|-|-|-|
| op | string | Operation type. Currently only `yaml-replace` is supported. Default is `yaml-replace`. | No |
| path | string | JSONPath to the field to replace (e.g. `$.spec.replicas`). | Yes |
| value | string | New value for the field. | Yes |
