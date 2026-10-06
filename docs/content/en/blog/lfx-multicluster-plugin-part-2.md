---
date: 2026-10-06
title: "Kubernetes Multi-Cluster Plugin for PipeCD, Part 2 — Per-Cluster Filtering, SDK Extensions and Production Hardening"
linkTitle: "Kubernetes Multi-Cluster Plugin for PipeCD, Part 2"
weight: 969
author: Mohammed Firdous ([@mohammedfirdouss](https://github.com/mohammedfirdouss))
categories: ["Contribution"]
tags: ["Kubernetes", "Plugin", "LFX Mentorship"]
---

*This post follows [Building the Kubernetes Multi-Cluster Plugin for PipeCD](https://pipecd.dev/blog/2026/04/10/building-the-kubernetes-multi-cluster-plugin-for-pipecd-lfx-mentorship/). Part 1 covered the six progressive delivery stages. This post covers the work that came after, which made those stages more flexible, easier to debug, and more reliable.*

---

## Where Part 1 Left Off

At the end of Part 1, the plugin could run Canary, Baseline, Primary, Traffic Routing, and their cleanup stages on all clusters at the same time. However, every stage always ran on every cluster, and when something failed, the plugin did not tell you *which* cluster failed.

Part 1 ended by listing the next steps: DetermineStrategy and drift detection. This post covers both. It also covers plan preview, Helm/OCI authentication, config validation, health checks for all workload types, better rollback, and settings you can set per cluster.

All screenshots in this post come from a local setup. It has a PipeCD control plane in one kind cluster, and one piped that deploys to two more kind clusters. These two clusters are registered as the deploy targets `cluster-eu` and `cluster-us`.

![PipeCD applications list showing the demo applications deployed to cluster-eu and cluster-us](/images/multicluster-part2-applications-overview.png)

---

## Per-Stage Cluster Filtering

Every stage ran on every cluster. In a real pipeline, you may want the canary to run only on `cluster-eu`, the baseline only on `cluster-us`, and the primary rollout on all clusters.

PR [#6757](https://github.com/pipe-cd/pipecd/pull/6757) added a `multiTargets` option to each stage:

```yaml
stages:
  - name: K8S_MULTI_CANARY_ROLLOUT
    with:
      multiTargets: ["cluster-eu"]
  - name: K8S_MULTI_BASELINE_ROLLOUT
    with:
      multiTargets: ["cluster-us"]
  - name: K8S_MULTI_PRIMARY_ROLLOUT
    # no multiTargets = all clusters
```

Each stage keeps only the clusters that are both in its `multiTargets` list and in the application's deploy targets. If `multiTargets` is empty, the stage runs on all clusters. This is the same as the old behaviour, so existing pipelines keep working.

The filtering happens before the stage starts applying changes. Clusters that are not in the list are never touched. This means no partial changes and no rollout to a cluster by mistake. In the config code, every stage options type has an optional `MultiTargets []string` field:

```go
type K8sCanaryRolloutStageOptions struct {
    Replicas      unit.Replicas `json:"replicas"`
    Suffix        string        `json:"suffix" default:"canary"`
    CreateService bool          `json:"createService"`
    Patches       []K8sResourcePatch `json:"patches"`
    // Empty means all deploy targets.
    MultiTargets  []string      `json:"multiTargets,omitempty"`
}
```

The canary stage log only mentions `cluster-eu`:

![Per-stage cluster filtering: K8S_MULTI_CANARY_ROLLOUT stage log running on cluster-eu only](/images/multicluster-part2-stage-filtering-canary-log.png)

While the pipeline is paused, you can check the clusters directly. The canary exists only on `cluster-eu`, and the baseline exists only on `cluster-us`.

![kubectl output: canary-demo-canary on cluster-eu only, canary-demo-baseline on cluster-us only](/images/multicluster-part2-stage-filtering-kubectl.png)

---

## Per-Cluster Visibility: DeployTargetStatus in the SDK

When a stage ran on several clusters and one of them failed, you got a single combined error and had to work out which cluster caused it.

PR [#6861](https://github.com/pipe-cd/pipecd/pull/6861) added a `DeployTargetStatuses` field to the plugin SDK's `ExecuteStageResponse`:

```go
type DeployTargetStatus struct {
    // Name is the deploy target name.
    Name    string
    // Status is the outcome for this specific target.
    Status  StageStatus
    // Message is an optional human-readable detail.
    Message string
}
```

PR [#6812](https://github.com/pipe-cd/pipecd/pull/6812) used this field in every stage of the plugin. The stage dispatch in `plugin.go` now looks like this:

```go
case StageK8sMultiCanaryRollout:
    status, dtStatuses := p.executeK8sMultiCanaryRolloutStage(ctx, input, dts)
    return &sdk.ExecuteStageResponse{
        Status:               status,
        DeployTargetStatuses: dtStatuses,
    }, nil
```

Each stage now returns a status for each cluster, together with the overall status. The stage log also names the cluster that failed. In the run below, `cluster-us` got a manifest that the API server rejects. The stage fails, and the log shows which cluster caused the failure. The apply that was still running on `cluster-eu` is cancelled instead of being left to finish:

![Failed K8S_MULTI_CANARY_ROLLOUT stage: the log names cluster-us as the failed target](/images/multicluster-part2-per-cluster-failure-canary-log.png)

Right now, the plugin returns these per-cluster statuses through the SDK. Showing them in a separate per-cluster view in the web UI is planned as future work.

---

## Plan Preview

Before running a pipeline, operators often want to see which manifests will be added, updated, or removed, without applying anything. PipeCD supports this with its plan preview API.

PR [#6685](https://github.com/pipe-cd/pipecd/pull/6685) added plan preview to the multi-cluster plugin. It returns one diff *for each deploy target* instead of one combined diff. `cluster-eu` and `cluster-us` can load different manifests, so their diffs can be different, and one combined diff would hide that.

```go
// Multi-target: produce one PlanPreviewResult per deploy target.
results := make([]sdk.PlanPreviewResult, 0, len(dts))
for _, dt := range dts {
    newManifests, err := loadManifests(ctx, loader, input, &targetDS, targetSpec, mt)
    // ...
    result, err := provider.DiffList(oldManifests, newManifests, ...)
    results = append(results, toResult(result, dt.Name))
}
return &sdk.GetPlanPreviewResponse{Results: results}, nil
```

Each result includes the deploy target name, a summary, and the full diff, with secrets hidden. If no `multiTargets` are set, the plugin loads the manifests once and returns a single result.

In this example, a branch changes the replica count for `cluster-eu` only:

![pipectl plan-preview output showing 1 changed manifest for cluster-eu and no changes for cluster-us](/images/multicluster-part2-plan-preview-per-cluster.png)

One gap turned up while taking that screenshot: the preview applied the per-cluster `manifests` override but not the per-cluster `kustomizeDir`, `kustomizeVersion` and `kustomizeOptions` that the deployment path already used, so applications with per-cluster Kustomize overlays always previewed as "No changes". [#7437](https://github.com/pipe-cd/pipecd/pull/7437) makes the preview use the same per-cluster settings as the deployment.

---

## Helm/OCI Initializer Hook

The plugin's `Initialize` hook runs once when piped starts, before any deployment. PR [#6723](https://github.com/pipe-cd/pipecd/pull/6723) added Helm repository setup and OCI registry login to this hook.

```go
func (i *initializer) Initialize(ctx context.Context,
    input *sdk.InitializeInput[...]) error {

    helm := provider.NewHelm("", helmPath, input.Logger)

    // Add and update HTTP Helm chart repositories.
    for _, repo := range input.Config.HTTPHelmChartRepositories() {
        if err := helm.AddRepository(ctx, repo); err != nil {
            return err
        }
    }
    helm.UpdateRepositories(ctx)

    // Log in to OCI registries.
    for _, registry := range input.Config.ChartRegistries {
        if !registry.IsOCI() {
            continue
        }
        if err := helm.LoginToOCIRegistry(ctx,
            registry.Address, registry.Username, registry.Password); err != nil {
            return err
        }
    }
    return nil
}
```

You set both of these once, at the plugin level in the piped config:

```yaml
plugins:
  - name: kubernetes_multicluster
    config:
      chartRepositories:
        - name: demo-charts
          address: http://127.0.0.1:8879
      chartRegistries:
        - type: OCI
          address: localhost:5002
          username: piped
          password: <registry-password>
```

Because login happens in `Initialize`, it runs once at startup, not on every deployment. In the screenshot below, the `helm` binary that piped uses is first rejected by a password-protected registry. After piped starts and the hook runs, the same binary can pull from it:

![piped startup log adding the chart repository and logging in to the OCI registry, with helm pull failing before and succeeding after](/images/multicluster-part2-helm-oci-initialize.png)

This prepares the plugin for remote charts. Today the plugin only renders local charts (`helmChart.path`). Rendering a chart directly from a repository or registry is the next step, and the login it needs is already done.

---

## Config Validation

PR [#6747](https://github.com/pipe-cd/pipecd/pull/6747) added a `Validate()` method to `KubernetesApplicationSpec`. It catches two kinds of mistakes when the config is read, instead of later during the deployment:

```go
func (s *KubernetesApplicationSpec) Validate() error {
    // helmChart and kustomizeOptions are mutually exclusive.
    if s.Input.HelmChart != nil && len(s.Input.KustomizeOptions) > 0 {
        return errors.New("helmChart and kustomizeOptions are mutually exclusive")
    }
    // Every multiTarget entry must have a name.
    for i, mt := range s.Input.MultiTargets {
        if mt.Target.Name == "" {
            return fmt.Errorf("multiTargets[%d].target.name must not be empty", i)
        }
    }
    return nil
}
```

The first check stops you from setting both a Helm chart and Kustomize options in the same spec, because the two cannot be used together. The second check makes sure every `multiTargets` entry has a name, so the filtering has something to match.

If you commit a config that breaks either rule, the deployment fails while the config is being read. No stage reaches any cluster:

![PipeCD deployment marked FAILURE after an invalid application config was committed](/images/multicluster-part2-config-validation-failed-deployment.png)

The piped log shows the reason. The workload that was already running stays the same on both clusters:

![Terminal showing the two invalid configs, the validation errors from the piped log, and the unchanged DaemonSet on both clusters](/images/multicluster-part2-config-validation-terminal.png)

---

## Better Rollback

In Part 1, rollback restored the previous primary. Two things were missing.

**Removing leftover resources** ([#6748](https://github.com/pipe-cd/pipecd/pull/6748)): a resource might exist in the previous deployment but be removed from Git before the rollback. In that case, it stayed in the cluster. Now rollback compares what is running with what should be running after the rollback, and deletes anything extra.

**Cleaning up Canary and Baseline** ([#6660](https://github.com/pipe-cd/pipecd/pull/6660)): if a pipeline fails halfway, the canary and baseline pods keep running on the clusters that got them. Before, rollback only restored the primary. Now, after it re-applies the previous primary manifests, rollback also finds and deletes all `variant=canary` and `variant=baseline` resources on every cluster.

The screenshot below shows the rollback after the failed canary stage from earlier. The baseline had been rolled out to both clusters, and the canary had reached `cluster-eu`. Rollback removes all of them:

![K8S_MULTI_ROLLBACK stage log: primary restored, then canary and baseline variants removed on both clusters](/images/multicluster-part2-rollback-variant-cleanup-log.png)

![kubectl confirming no canary or baseline resources remain after rollback, primary back on the previous image](/images/multicluster-part2-post-rollback-clean-clusters.png)

---

## Health Checks for All Workload Types

The plugin checked the health of Deployments and StatefulSets, but it skipped DaemonSets, ReplicaSets, and Pods without any warning. PR [#6807](https://github.com/pipe-cd/pipecd/pull/6807) added the missing checks:

- **DaemonSet**: healthy when `NumberReady == DesiredNumberScheduled`
- **ReplicaSet**: healthy when `ReadyReplicas >= Replicas`
- **Pod**: healthy when all containers are running and none are in a crash loop

Without these checks, a DaemonSet that could not start on half of the nodes would still show as healthy.

![PipeCD application view: daemonset-demo reported Synced and Healthy, with a tab per cluster](/images/multicluster-part2-daemonset-health-livestate.png)

---

## Drift Detection

When the plugin read the live state of a cluster, it loaded all resources but did not correctly keep only the ones that belong to the application. PR [#6673](https://github.com/pipe-cd/pipecd/pull/6673) fixed this. The plugin now selects resources using the labels PipeCD adds at deploy time: `pipecd.dev/managed-by`, `pipecd.dev/piped`, and `pipecd.dev/application`.

PR [#6672](https://github.com/pipe-cd/pipecd/pull/6672) updated `DetermineStrategy` to handle multiple clusters. Before, it looked only at the first cluster to decide between QuickSync and the full pipeline. Now it uses the full pipeline if *any* cluster has a workload change.

Drift is shown for each cluster. In this example, someone ran `kubectl scale` on `cluster-us` only. The diff for `cluster-us` shows `replicas: 4`, while Git says `2`:

![PipeCD UI showing Out of Sync, with a separate diff for cluster-eu and cluster-us](/images/multicluster-part2-drift-detected-cluster-us.png)

The screenshot also shows the two things that make this report noisier than it should be. The auto-generated `Endpoints` and `EndpointSlice` objects of the Service are listed as deletions, and until [#7438](https://github.com/pipe-cd/pipecd/pull/7438) the plugin stamped the head commit of the repository onto the manifests it compared against, so any commit to the repository, even one for another application, reported every application as out of sync. The comparison now leaves the commit hash out.

---

## Per-Cluster Settings

Three PRs added settings that you can now set for each cluster instead of only once for all clusters. In each case, the plugin starts with the top-level `spec.input` value and uses the per-cluster value instead if one is set:

```go
kustomizeVersion := spec.Input.KustomizeVersion  // global default
if multiTarget != nil && multiTarget.KustomizeVersion != "" {
    kustomizeVersion = multiTarget.KustomizeVersion  // per-target override
}
```

The three additions:

- **`kustomizeDir`** ([#6718](https://github.com/pipe-cd/pipecd/pull/6718)): the folder to use for Kustomize. This helps when each cluster uses a different overlay from the same repo.
- **`kustomizeVersion` and `kustomizeOptions`** ([#6749](https://github.com/pipe-cd/pipecd/pull/6749)): use a different Kustomize version or options for each cluster.
- **Config hash for StatefulSet and DaemonSet** ([#6697](https://github.com/pipe-cd/pipecd/pull/6697)): StatefulSets and DaemonSets now roll out again when a ConfigMap or Secret they use changes, as Deployments already did. The plugin does this by adding a hash of the config data as a pod annotation.

A full multi-cluster config with per-cluster settings looks like this:

```yaml
input:
  multiTargets:
    - target:
        name: cluster-eu
      kustomizeDir: overlays/eu
      kustomizeVersion: "5.3.0"
    - target:
        name: cluster-us
      kustomizeDir: overlays/us
      kustomizeVersion: "5.4.0"
    - target:
        name: cluster-asia
      # no overrides, so it uses the top-level input values
```

Here `overlays/eu` sets 2 replicas and `overlays/us` sets 3. One sync gives a different result on each cluster:

![kubectl output: kustomize-demo has 2 replicas and region=eu on cluster-eu, 3 replicas and region=us on cluster-us](/images/multicluster-part2-per-target-kustomize-kubectl.png)

---

## Stage Rename

PR [#6880](https://github.com/pipe-cd/pipecd/pull/6880) added a `K8S_MULTI_` prefix to all stage names. This was needed because the single-cluster Kubernetes plugin uses the same stage names (`K8S_CANARY_ROLLOUT`, `K8S_PRIMARY_ROLLOUT`, and so on). When both plugins were loaded together, the names clashed, and PipeCD could not tell which plugin should run the stage.

To migrate, find and replace the stage names in your pipeline YAML:

| Old name | New name |
|---|---|
| `K8S_SYNC` | `K8S_MULTI_SYNC` |
| `K8S_ROLLBACK` | `K8S_MULTI_ROLLBACK` |
| `K8S_CANARY_ROLLOUT` | `K8S_MULTI_CANARY_ROLLOUT` |
| `K8S_CANARY_CLEAN` | `K8S_MULTI_CANARY_CLEAN` |
| `K8S_BASELINE_ROLLOUT` | `K8S_MULTI_BASELINE_ROLLOUT` |
| `K8S_BASELINE_CLEAN` | `K8S_MULTI_BASELINE_CLEAN` |
| `K8S_PRIMARY_ROLLOUT` | `K8S_MULTI_PRIMARY_ROLLOUT` |
| `K8S_TRAFFIC_ROUTING` | `K8S_MULTI_TRAFFIC_ROUTING` |

---

## Where Things Stand

The plugin now covers most of what the single-cluster plugin does. It has six progressive delivery stages, per-stage cluster filtering, per-cluster status, plan preview, Helm/OCI login, config validation, health checks for all five workload types, better rollback, drift detection, and per-cluster settings. The main piece still missing is templating Helm charts straight from a chart repository or OCI registry; the plugin templates local charts today, and the login added in the initializer is the groundwork for that. In total, 39 PRs were merged during the mentorship term.

The plugin is part of the open-source PipeCD project. Contributions are welcome.

## Links

- [PipeCD repository](https://github.com/pipe-cd/pipecd)
- [Part 1: Building the six stages](https://pipecd.dev/blog/2026/04/10/building-the-kubernetes-multi-cluster-plugin-for-pipecd-lfx-mentorship/)
- [Plugin docs](https://pipecd.dev/docs-v1.0.x/plugins/official/kubernetes-multicluster/)
- [LFX Mentorship Program](https://mentorship.lfx.linuxfoundation.org)
