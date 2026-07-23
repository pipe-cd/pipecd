---
title: "Kubernetes plugin"
linkTitle: "Kubernetes"
weight: 10
description: >
  Deploy an application to a Kubernetes cluster.
---

The `kubernetes` plugin deploys an application to a single Kubernetes cluster. It supports quick sync and pipeline-based rollouts with canary, baseline, and primary variants, and can shift traffic between variants using the PodSelector or Istio method.

By default, `piped` deploys the application to the cluster it runs in. To deploy to an external cluster, set `masterURL` and `kubeConfigPath` in the plugin's [`deployTargets`](#kubernetesdeploytargetconfig).

## Prerequisites

1. **Register the plugin in the piped configuration.** Add a `kubernetes` plugin block with a `deployTargets` entry for the cluster:

   ```yaml
   apiVersion: pipecd.dev/v1beta1
   kind: Piped
   spec:
     # ...
     plugins:
       - name: kubernetes
         port: 7001
         url: file:///path/to/plugin/binary  # or an https:// release URL
         deployTargets:
           - name: local
             config:
               # Empty masterURL/kubeConfigPath means the cluster piped runs in.
               kubectlVersion: 1.32.2
   ```

2. **The plugin downloads `kubectl` automatically.** `piped` fetches `kubectl` (and `kustomize`/`helm` when used) via the tool registry, so the binary does not need to be pre-installed on the `piped` host. Pin a version with `kubectlVersion` if you need a specific one.

3. When registering the application in the control plane, **select the deploy target** the application should deploy to.

## Quick sync

With no `pipeline` defined, the plugin performs a **quick sync** (`K8S_SYNC`): it applies all manifests to the cluster. This minimal `app.pipecd.yaml` deploys the manifests in the application directory:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: my-app
  plugins:
    kubernetes:
      input:
        manifests:
          - deployment.yaml
          - service.yaml
        kubectlVersion: 1.32.2
```

## Sync with the specified pipeline

Define a `pipeline` to control the rollout. The stages run in order:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: my-app
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 1
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_CANARY_CLEAN
  plugins:
    kubernetes:
      input:
        manifests:
          - deployment.yaml
        kubectlVersion: 1.32.2
```

See [Pipeline stages](#pipeline-stages) for every available stage and its options.

## Pipeline stages

Stages are listed under `spec.pipeline.stages`, with options under `with`.

### K8S_SYNC

Applies all manifests to the cluster. This is the stage that runs during a quick sync when no `pipeline` block is defined in your `app.pipecd.yaml`.

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| addVariantLabelToSelector | bool | Add the PRIMARY variant label to workload selectors if missing. | false |
| prune | bool | Remove resources that are no longer defined in Git. | false |

### K8S_PRIMARY_ROLLOUT

Rolls out the PRIMARY (stable) variant using the manifests defined in Git. You can optionally create a dedicated Service for the PRIMARY variant and enable pruning so that any resources removed from Git are also deleted from the cluster.

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| suffix | string | Suffix for the PRIMARY variant's resource names. | primary |
| createService | bool | Create a Service for the PRIMARY variant. | false |
| addVariantLabelToSelector | bool | Add the PRIMARY variant label to workload selectors if missing. | false |
| prune | bool | Remove resources no longer defined in Git. | false |

### K8S_CANARY_ROLLOUT

Creates CANARY variant workloads alongside the currently running version. This lets you send a portion of traffic to the new version and compare its behaviour against the stable version before deciding to promote or roll back. You can optionally create a dedicated Service for the canary and apply manifest patches to customise the variant before it is deployed.

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| replicas | int or string | Number of CANARY pods. An integer, or a percentage of PRIMARY (e.g. `"50%"`). | 1 |
| suffix | string | Suffix for the CANARY variant's resource names. | canary |
| createService | bool | Create a Service for the CANARY variant. | false |
| patches | [][K8sResourcePatch](#k8sresourcepatch) | Patches applied to manifests before generating the CANARY variant. | - |

### K8S_CANARY_CLEAN

Removes the CANARY variant resources that were created by `K8S_CANARY_ROLLOUT`. This stage is typically placed at the end of a canary pipeline to clean up after promotion or after a rollback decision. It takes no options.

### K8S_BASELINE_ROLLOUT

Creates BASELINE variant workloads from the **running** (currently live) manifests, not the target manifests. This means the baseline is an exact copy of what is already in production, giving you a stable reference point to compare against the canary during analysis.

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| replicas | int or string | Number of BASELINE pods. An integer, or a percentage of PRIMARY. | 1 |
| suffix | string | Suffix for the BASELINE variant's resource names. | baseline |
| createService | bool | Create a Service for the BASELINE variant. | false |

### K8S_BASELINE_CLEAN

Removes the BASELINE variant resources that were created by `K8S_BASELINE_ROLLOUT`. Place this stage at the end of a canary/baseline pipeline to clean up after the analysis is complete, whether you promoted or rolled back. It takes no options.

### K8S_TRAFFIC_ROUTING

Shifts traffic between the PRIMARY, CANARY, and BASELINE variants. The routing method is set by the application-level [`trafficRouting`](#kubernetestrafficrouting) config and determines how the split is applied.

**PodSelector** (default) works by updating the Service's `spec.selector` to point entirely at one variant. This means one variant must receive 100% of traffic at a time, and baseline is not supported with this method.

**Istio** works by updating the `VirtualService` route weights, which allows traffic to be split by percentage across all three variants simultaneously.

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| all | string | Send all traffic to one variant: `primary`, `canary`, or `baseline`. | - |
| primary | int | Percentage of traffic to the PRIMARY variant. | 0 |
| canary | int | Percentage of traffic to the CANARY variant. | 0 |
| baseline | int | Percentage of traffic to the BASELINE variant. | 0 |

PodSelector: switch all traffic to the canary variant:

```yaml
- name: K8S_TRAFFIC_ROUTING
  with:
    all: canary
```

Istio: split traffic 80/20 between primary and canary (requires `trafficRouting.method: istio` on the app):

```yaml
- name: K8S_TRAFFIC_ROUTING
  with:
    primary: 80
    canary: 20
```

### K8S_ROLLBACK

Restores the application to the previously running manifests. This stage is triggered automatically when a deployment fails or is cancelled with rollback. You do not add it to your pipeline; PipeCD inserts it automatically.

## Livestate and drift detection

The plugin reports the live state of the application, and detects drift by comparing the live cluster state against the manifests in Git. Two controls scope what counts as drift:

- **Ignore annotation**: annotate a resource with `pipecd.dev/ignore-drift-detection: "true"` and changes to it no longer flip the application to `OUT_OF_SYNC`.
- **Include/exclude filtering**: narrow which resource kinds the informer watches via `appStateInformer` in the deploy target config. For example, to stop Secret changes from being reported as drift:

  ```yaml
  deployTargets:
    - name: local
      config:
        appStateInformer:
          excludeResources:
            - apiVersion: v1
              kind: Secret
  ```

## Plan preview

Before a pipeline runs, plan preview shows the manifest diff that the deployment will apply. Secret values are masked as `***` in the preview.

## Helm authentication for private repositories and registries

Credentials for private Helm chart repositories and OCI registries are configured in the **piped** configuration under the plugin's `config` block, never in `app.pipecd.yaml`. When the plugin starts, it runs `helm repo add` / `helm repo update` for each configured repository and `helm registry login` for each OCI registry, once, before any deployment. Passwords are passed to `helm` via stdin so they do not appear in the process list.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  plugins:
    - name: kubernetes
      # port / url / deployTargets ...
      config:
        chartRepositories:
          - type: HTTP                       # ChartMuseum, Nexus, Artifactory, ...
            name: my-private-repo
            address: https://charts.example.com
            username: my-user
            password: my-password
            insecure: false
        chartRegistries:
          - type: OCI                        # GHCR, ECR, Google Artifact Registry, ...
            address: ghcr.io
            username: my-github-user
            password: my-pat-token
```

`chartRepositories` fields: `type` (only `HTTP`), `name`, `address`, `username`, `password`, `insecure`.
`chartRegistries` fields: `type` (only `OCI`), `address`, `username`, `password`.

## Configuration reference

### KubernetesApplicationSpec

The `spec` of a `KubernetesApp` shares the [common application fields](../user-guide/managing-application/configuration-reference/) (`name`, `labels`, `pipeline`, ...) and adds the following under `plugins.kubernetes`:

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| input | [KubernetesDeploymentInput](#kubernetesdeploymentinput) | Input for the deployment such as manifests and tool versions. | Yes |
| quickSync | [K8S_SYNC options](#k8s_sync) | Options applied when the application is deployed via quick sync (no pipeline). | No |
| workloads | [][K8sResourceReference](#k8sresourcereference) | Which resources are treated as the application's workloads. Empty means all `Deployment`s. | No |
| service | [K8sResourceReference](#k8sresourcereference) | Which resource is treated as the application's Service. Empty means the first `Service`. | No |
| variantLabel | [KubernetesVariantLabel](#kubernetesvariantlabel) | The label used to distinguish variant (primary/canary/baseline) manifests. | No |
| trafficRouting | [KubernetesTrafficRouting](#kubernetestrafficrouting) | Traffic routing configuration. Defaults to the PodSelector method. | No |

### KubernetesDeploymentInput

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| manifests | []string | List of manifest files in the application directory used to deploy. Empty means all manifest files in the directory. | No |
| kubectlVersion | string | Version of `kubectl` to use. | No |
| kustomizeVersion | string | Version of `kustomize` to use. | No |
| kustomizeOptions | map[string]string | Options passed to `kustomize` commands (e.g. `load-restrictor: LoadRestrictionsNone`). | No |
| helmVersion | string | Version of `helm` to use. | No |
| helmChart | [InputHelmChart](#inputhelmchart) | Where to fetch the Helm chart. | No |
| helmOptions | [InputHelmOptions](#inputhelmoptions) | Configurable parameters for `helm` commands. | No |
| namespace | string | The namespace where manifests are applied. | No |
| autoCreateNamespace | bool | Automatically create the namespace if it does not exist. Default is `false`. | No |

### KubernetesDeployTargetConfig

Configured under `plugins[].deployTargets[].config` in the **piped** configuration.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| masterURL | string | Master URL of the cluster. Empty means the cluster `piped` runs in. | No |
| kubeConfigPath | string | Path to the kubeconfig file. Empty means the cluster `piped` runs in. | No |
| kubectlVersion | string | Default `kubectl` version for this deploy target. | No |
| appStateInformer | [KubernetesAppStateInformer](#kubernetesappstateinformer) | Scopes which resources the livestate informer watches. | No |

### KubernetesAppStateInformer

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| namespace | string | Only watch this namespace. Empty means all namespaces. | No |
| includeResources | [][KubernetesResourceMatcher](#kubernetesresourcematcher) | Resources to add to the watch targets. | No |
| excludeResources | [][KubernetesResourceMatcher](#kubernetesresourcematcher) | Resources to ignore from the watch targets. | No |

### KubernetesResourceMatcher

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| apiVersion | string | API version of the resource. | Yes |
| kind | string | Kind of the resource. Empty means all kinds match. | No |

### KubernetesVariantLabel

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| key | string | The label key. Default `pipecd.dev/variant`. | No |
| primaryValue | string | Label value for the PRIMARY variant. Default `primary`. | No |
| canaryValue | string | Label value for the CANARY variant. Default `canary`. | No |
| baselineValue | string | Label value for the BASELINE variant. Default `baseline`. | No |

### InputHelmChart

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| path | string | Relative path from the repository root to the chart directory (for a local chart). | No |
| repository | string | Name of an added Helm chart repository. | No |
| name | string | Chart name. | No |
| version | string | Chart version. | No |

### InputHelmOptions

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| releaseName | string | Release name of the Helm deployment. Defaults to the application name. | No |
| setValues | map[string]string | Values passed via `--set`. | No |
| valueFiles | []string | Value files to load (must be within the application directory). | No |
| setFiles | map[string]string | File paths whose contents are passed via `--set-file`. | No |
| apiVersions | []string | Supported Kubernetes API versions passed via `--api-versions`. | No |
| kubeVersion | string | Kubernetes version for `Capabilities.KubeVersion`. | No |

### K8sResourceReference

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| kind | string | Kind of the referenced resource. | Yes |
| name | string | Name of the referenced resource. | Yes |

### KubernetesTrafficRouting

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| method | string | Routing method: `podselector` or `istio`. Default `podselector`. | No |
| istio | [IstioTrafficRouting](#istiotrafficrouting) | Istio-specific configuration (used when `method: istio`). | No |

### IstioTrafficRouting

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| editableRoutes | []string | Routes in the VirtualService that may be updated. Empty means all routes. | No |
| host | string | The service host. | No |
| virtualService | [K8sResourceReference](#k8sresourcereference) | Reference to the VirtualService manifest. Empty means the first one. | No |

### K8sResourcePatch

Used by `K8S_CANARY_ROLLOUT.patches` to customize manifests before generating the CANARY variant.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| target.kind | string | Kind of the manifest to patch. | Yes |
| target.name | string | Name of the manifest to patch. | Yes |
| target.documentRoot | string | Path to a field whose string value is patched, e.g. `$.data.values\.yaml`. Empty means the whole manifest. | No |
| ops | [][K8sResourcePatchOp](#k8sresourcepatchop) | List of patch operations to apply. | Yes |

### K8sResourcePatchOp

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| op | string | Operation type. Only `yaml-replace` is supported. Default is `yaml-replace`. | No |
| path | string | JSONPath expression pointing to the field to replace, e.g. `$.spec.replicas`. | Yes |
| value | string | New value to set at `path`. | Yes |
