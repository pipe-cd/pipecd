---
title: "Kubernetes multi-cluster plugin"
linkTitle: "Kubernetes multi-cluster"
weight: 10
description: >
  Deploy a single Kubernetes application to multiple clusters with one pipeline.
---

> **Note:** The `kubernetes_multicluster` plugin is in **Alpha**. Its configuration schema may change in future releases.

The `kubernetes_multicluster` plugin deploys one application to **multiple Kubernetes clusters** from a single `app.pipecd.yaml`. Every pipeline stage runs against all of the application's deploy targets in parallel, unless a stage is restricted to a subset of targets. It supports the same deployment strategies as the single-cluster Kubernetes plugin (quick sync, canary, baseline, primary rollout, and traffic routing) and extends them with per-cluster manifests, per-cluster tool versions, and per-stage target filtering.

## Prerequisites

1. **Register each cluster as a deploy target** in the piped configuration. Add a `kubernetes_multicluster` plugin block with one `deployTargets` entry per cluster:

   ```yaml
   apiVersion: pipecd.dev/v1beta1
   kind: Piped
   spec:
     # ...
     plugins:
       - name: kubernetes_multicluster
         port: 7002
         url: file:///path/to/plugin/binary  # or an https:// release URL
         deployTargets:
           - name: cluster-eu
             config:
               kubeConfigPath: /etc/piped/kube/cluster-eu
           - name: cluster-us
             config:
               kubeConfigPath: /etc/piped/kube/cluster-us
   ```

2. **The plugin downloads `kubectl` automatically.** Piped fetches `kubectl` (and `kustomize`/`helm` when used) via the tool registry so the binary does not need to be pre-installed on the piped host. Pin a version with `kubectlVersion` if you need a specific one.

3. When registering the application in the control plane, **select every cluster** the app should deploy to as its deploy targets.

## Quick sync

With no `pipeline` defined, the plugin performs a **quick sync** (`K8S_MULTI_SYNC`): it applies all manifests to every selected cluster. This minimal `app.pipecd.yaml` deploys the same manifests to two clusters simultaneously:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: my-app
  plugins:
    kubernetes_multicluster:
      input:
        manifests:
          - deployment.yaml
          - service.yaml
        kubectlVersion: 1.32.2
```

## Sync with the specified pipeline

Define a `pipeline` to control the rollout. The stages run in order, and each stage executes on all deploy targets in parallel:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: my-app
  pipeline:
    stages:
      - name: K8S_MULTI_CANARY_ROLLOUT
        with:
          replicas: 1
      - name: K8S_MULTI_PRIMARY_ROLLOUT
      - name: K8S_MULTI_CANARY_CLEAN
  plugins:
    kubernetes_multicluster:
      input:
        manifests:
          - deployment.yaml
        kubectlVersion: 1.32.2
```

See [Pipeline stages](#pipeline-stages) for every available stage and its options.

## Pipeline stages

Stages are listed under `spec.pipeline.stages`, with options under `with`. By default a stage runs on **all** of the application's deploy targets; add `multiTargets` to restrict it to named clusters (see [Per-stage target filtering](#per-stage-target-filtering)).

### K8S_MULTI_SYNC

Applies all manifests to every selected deploy target. This is the stage that runs during a quick sync when no `pipeline` block is defined in your `app.pipecd.yaml`. The available options are the same as those in the `quickSync` block of the application config.

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| addVariantLabelToSelector | bool | Add the PRIMARY variant label to workload selectors if missing. | false |
| prune | bool | Remove resources that are no longer defined in Git. | false |
| multiTargets | []string | Restrict the stage to these deploy targets. | all |

### K8S_MULTI_PRIMARY_ROLLOUT

Rolls out the PRIMARY (stable) variant to all selected targets using the manifests defined in Git. You can optionally create a dedicated Service for the PRIMARY variant and enable pruning so that any resources removed from Git are also deleted from the cluster.

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| suffix | string | Suffix for the PRIMARY variant's resource names. | primary |
| createService | bool | Create a Service for the PRIMARY variant. | false |
| addVariantLabelToSelector | bool | Add the PRIMARY variant label to workload selectors if missing. | false |
| prune | bool | Remove resources no longer defined in Git. | false |
| multiTargets | []string | Restrict the stage to these deploy targets. | all |

### K8S_MULTI_CANARY_ROLLOUT

Creates CANARY variant workloads alongside the currently running version. This lets you send a portion of traffic to the new version and compare its behaviour against the stable version before deciding to promote or roll back. You can optionally create a dedicated Service for the canary and apply manifest patches to customise the variant before it is deployed.

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| replicas | int or string | Number of CANARY pods. An integer, or a percentage of PRIMARY (e.g. `"50%"`). | 1 |
| suffix | string | Suffix for the CANARY variant's resource names. | canary |
| createService | bool | Create a Service for the CANARY variant. | false |
| patches | [][K8sResourcePatch](#k8sresourcepatch) | Patches applied to manifests before generating the CANARY variant. | - |
| multiTargets | []string | Restrict the stage to these deploy targets. | all |

### K8S_MULTI_CANARY_CLEAN

Removes the CANARY variant resources that were created by `K8S_MULTI_CANARY_ROLLOUT`. This stage is typically placed at the end of a canary pipeline to clean up after promotion or after a rollback decision.

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| multiTargets | []string | Restrict the stage to these deploy targets. | all |

### K8S_MULTI_BASELINE_ROLLOUT

Creates BASELINE variant workloads from the **running** (currently live) manifests, not the target manifests. This means the baseline is an exact copy of what is already in production, giving you a stable reference point to compare against the canary during analysis.

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| replicas | int or string | Number of BASELINE pods. An integer, or a percentage of PRIMARY. | 1 |
| suffix | string | Suffix for the BASELINE variant's resource names. | baseline |
| createService | bool | Create a Service for the BASELINE variant. | false |
| multiTargets | []string | Restrict the stage to these deploy targets. | all |

### K8S_MULTI_BASELINE_CLEAN

Removes the BASELINE variant resources that were created by `K8S_MULTI_BASELINE_ROLLOUT`. Place this stage at the end of a canary/baseline pipeline to clean up after the analysis is complete, whether you promoted or rolled back.

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| multiTargets | []string | Restrict the stage to these deploy targets. | all |

### K8S_MULTI_TRAFFIC_ROUTING

Shifts traffic between the PRIMARY, CANARY, and BASELINE variants. The routing method is set by the application-level [`trafficRouting`](#kubernetestrafficrouting) config and determines how the split is applied.

**PodSelector** (default) works by updating the Service's `spec.selector` to point entirely at one variant. This means one variant must receive 100% of traffic at a time and baseline is not supported with this method.

**Istio** works by updating the `VirtualService` route weights, which allows traffic to be split by percentage across all three variants simultaneously.

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| all | string | Send all traffic to one variant: `primary`, `canary`, or `baseline`. | - |
| primary | int | Percentage of traffic to the PRIMARY variant. | 0 |
| canary | int | Percentage of traffic to the CANARY variant. | 0 |
| baseline | int | Percentage of traffic to the BASELINE variant. | 0 |
| multiTargets | []string | Restrict the stage to these deploy targets. | all |

PodSelector: switch all traffic to the canary variant:

```yaml
- name: K8S_MULTI_TRAFFIC_ROUTING
  with:
    all: canary
```

Istio: split traffic 80/20 between primary and canary (requires `trafficRouting.method: istio` on the app)

```yaml
- name: K8S_MULTI_TRAFFIC_ROUTING
  with:
    primary: 80
    canary: 20
```

### K8S_MULTI_ROLLBACK

Restores the application to the previously running manifests on every target. It prunes resources no longer defined in those manifests and cleans up any leftover CANARY or BASELINE variant resources. This stage is triggered automatically when a deployment fails or is cancelled with rollback. You do not add it to your pipeline; PipeCD inserts it automatically.

## Multi-cluster features

### One application, many clusters

A single `app.pipecd.yaml` deploys to every cluster selected as a deploy target. Every stage runs against all targets in parallel, and the deployment result reports the outcome per cluster, so one cluster failing does not hide the status of the others.

By default all targets share the same `input.manifests`. To give each cluster its own manifests or tool versions, declare `multiTargets`.

### Per-cluster manifests with multiTargets

Each `multiTargets` entry binds a deploy target name to its own manifest paths (and optional tool overrides). Use this when clusters need different manifests - for example region-specific config:

```yaml
spec:
  plugins:
    kubernetes_multicluster:
      input:
        multiTargets:
          - target:
              name: cluster-eu
            manifests:
              - ./cluster-eu/deployment.yaml
          - target:
              name: cluster-us
            manifests:
              - ./cluster-us/deployment.yaml
        kubectlVersion: 1.32.2
```

Per-target overrides available on each entry:

- `kubectlVersion` : use a different `kubectl` version for this cluster (takes precedence over `input.kubectlVersion`).
- `kustomizeDir`: point this cluster at its own kustomize overlay directory.
- `kustomizeVersion` / `kustomizeOptions` : per-cluster kustomize control.

```yaml
multiTargets:
  - target:
      name: cluster-eu
    kustomizeDir: ./overlays/eu
  - target:
      name: cluster-us
    kustomizeDir: ./overlays/us
    kubectlVersion: 1.31.0
```

### Per-stage target filtering

Any stage can be limited to a subset of clusters with the `multiTargets` option, which takes a list of deploy target **names**. When `multiTargets` is empty (the default) the stage runs on all targets. When it contains names, the stage runs only on those clusters. Any name that does not match a configured deploy target is silently ignored.

This lets you roll out cautiously. For example, canary on `cluster-eu` only, then promote to primary everywhere:

```yaml
pipeline:
  stages:
    - name: K8S_MULTI_CANARY_ROLLOUT
      with:
        replicas: 1
        multiTargets: [cluster-eu]
    - name: K8S_MULTI_PRIMARY_ROLLOUT          # all targets
    - name: K8S_MULTI_CANARY_CLEAN
      with:
        multiTargets: [cluster-eu]
```

## Livestate and drift detection

The plugin reports the live state of the application per cluster. Each deploy target gets its own sync status (`SYNCED`, `OUT_OF_SYNC`, or `UNKNOWN`) and resource tree, aggregated under a single application entry in the UI.

Drift is detected by comparing the live cluster state against the manifests in Git. Two controls scope what counts as drift:

- **Ignore annotation** : annotate a resource with `pipecd.dev/ignore-drift-detection: "true"` and changes to it no longer flip the application to `OUT_OF_SYNC`.
- **Include/exclude filtering** : narrow which resource kinds the informer watches via `appStateInformer` in the deploy target config. For example, to stop Secret changes from being reported as drift:

  ```yaml
  deployTargets:
    - name: cluster-us
      config:
        kubeConfigPath: /etc/piped/kube/cluster-us
        appStateInformer:
          excludeResources:
            - apiVersion: v1
              kind: Secret
  ```

## Plan preview

Before a pipeline runs, plan preview shows the manifest diff that the deployment will apply. For multi-cluster applications, the diff is rendered **per deploy target** so you can see exactly what will change on each cluster (including "No changes were detected" for clusters that are unaffected). Secret values are masked as `***` in the preview.

## Helm authentication for private repositories and registries

Credentials for private Helm chart repositories and OCI registries are configured in the **piped** configuration under the plugin's `config` block and never in `app.pipecd.yaml`. When the plugin starts, it runs `helm repo add` / `helm repo update` for each configured repository and `helm registry login` for each OCI registry, once, before any deployment. Passwords are passed to `helm` via stdin so they do not appear in the process list.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  plugins:
    - name: kubernetes_multicluster
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

The `spec` of a `KubernetesApp` shares the common application fields (`name`, `labels`, `pipeline`, ...) and adds the following under `plugins.kubernetes_multicluster`:

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| input | [KubernetesDeploymentInput](#kubernetesdeploymentinput) | Input for the deployment such as manifests, tool versions, and per-cluster targets. | Yes |
| quickSync | [K8S_MULTI_SYNC options](#k8s_multi_sync) | Options applied when the application is deployed via quick sync (no pipeline). | No |
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
| kustomizeOptions | map[string]string | Options passed to `kustomize` commands (e.g. `load-restrictor: LoadRestrictionsNone`). Mutually exclusive with `helmChart`. | No |
| helmVersion | string | Version of `helm` to use. | No |
| helmChart | [InputHelmChart](#inputhelmchart) | Where to fetch the Helm chart. | No |
| helmOptions | [InputHelmOptions](#inputhelmoptions) | Configurable parameters for `helm` commands. | No |
| namespace | string | The namespace where manifests are applied. | No |
| autoCreateNamespace | bool | Automatically create the namespace if it does not exist. Default is `false`. | No |
| multiTargets | [][KubernetesMultiTarget](#kubernetesmultitarget) | Per-cluster manifest paths and tool overrides. See [Multi-cluster features](#multi-cluster-features). | No |

### KubernetesMultiTarget

Declares per-cluster input. Each entry targets one deploy target by name and may override manifests and tool settings for that cluster.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| target.name | string | Name of the deploy target this entry applies to. | Yes |
| target.labels | map[string]string | Labels for the target. | No |
| manifests | []string | Manifest files for this target. Overrides the top-level `input.manifests` for this cluster. | No |
| kubectlVersion | string | `kubectl` version for this target. Takes precedence over `input.kubectlVersion`. | No |
| kustomizeDir | string | Path to a kustomize overlay directory for this target. | No |
| kustomizeVersion | string | `kustomize` version for this target. | No |
| kustomizeOptions | map[string]string | `kustomize` options for this target. | No |

### KubernetesDeployTargetConfig

Configured under `plugins[].deployTargets[].config` in the **piped** configuration, one per cluster.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| masterURL | string | Master URL of the cluster. Empty means in-cluster. | No |
| kubeConfigPath | string | Path to the kubeconfig file. Empty means in-cluster. | No |
| kubectlVersion | string | Default `kubectl` version for this cluster. | No |
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
| gitRemote | string | Git remote address where the chart is located. Empty means the same repository. | No |
| ref | string | Commit SHA or tag for the remote git. | No |
| path | string | Relative path to the chart directory (for a local chart). | No |
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

Used by `K8S_MULTI_CANARY_ROLLOUT.patches` to customize manifests before generating the CANARY variant.

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
