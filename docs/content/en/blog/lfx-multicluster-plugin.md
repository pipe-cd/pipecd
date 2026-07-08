---
date: 2026-04-10
title: "Building the Kubernetes Multi-Cluster Plugin for PipeCD — LFX Mentorship"
linkTitle: "Building the Kubernetes Multi-Cluster Plugin for PipeCD"
weight: 970
author: Mohammed Firdous ([@mohammedfirdouss](https://github.com/mohammedfirdouss))
categories: ["Contribution"]
tags: ["Kubernetes", "Plugin", "LFX Mentorship"]
---

If you had told me last year that I would be working with Kubernetes and all things clusters, deployments and service meshes, I would have brushed it off. I am truly grateful for the journey thus far.

Earlier last month, I got accepted as an LFX Mentee for Term 1 of this calendar year. For me it is such a big deal, given my background, and how much effort has been put in behind the scenes to get to this stage.

I'm currently a mentee in the LFX Mentorship program working on [PipeCD](https://pipecd.dev), an open-source GitOps continuous delivery platform. For the past four weeks, I've been building out the `kubernetes_multicluster` plugin specifically implementing the deployment pipeline stages that handle canary, primary and baseline deployments across multiple clusters.

---

## What is PipeCD and what is this plugin?

PipeCD is an open-source GitOps CD platform that manages deployments across different infrastructure targets like Kubernetes, ECS, Terraform, Lambda and more. Each target type has a plugin that knows how to deploy to it.

The `kubernetes_multicluster` plugin is for teams running the same application across multiple Kubernetes clusters say US, EU and Asia and needing all of them to stay in sync through a single pipeline. Rolling out a new version across clusters one at a time, manually, with no coordination, is error-prone and slow. The plugin lets you define one pipeline that runs across every cluster at the same time, with canary and baseline checks before anything hits production.

## Progressive Delivery and Why These Stages Exist

Before a new version reaches all users, it goes through stages. A canary sends a small slice of traffic to the new version first. A baseline runs the *current* version at the same scale so you have a fair comparison. Primary is the actual promotion. Clean stages remove the temporary resources when you're done.

This pattern is called progressive delivery, because you roll out gradually, check things look good, then commit. If something looks wrong at the canary stage, you stop there. Nothing has touched production yet.

The `kubernetes_multicluster` plugin runs all of this across every cluster at the same time. One pipeline, every cluster, same stages.

A full pipeline looks like this:

```yaml
stages:
  - name: K8S_CANARY_ROLLOUT
  - name: K8S_BASELINE_ROLLOUT
  - name: K8S_TRAFFIC_ROUTING
  - name: K8S_PRIMARY_ROLLOUT
  - name: K8S_CANARY_CLEAN
  - name: K8S_BASELINE_CLEAN
```

Each of these is a stage I built. The sections below go through what each one does.

## What I Built

### K8S_CANARY_ROLLOUT

The canary stage deploys the new version of your app as a small slice alongside the existing production deployment. If your app normally runs 3 pods, canary might spin up 1 pod (or 20%) of the new version enough to catch problems without affecting most users.

It loads manifests from Git, creates copies of all workloads with a `-canary` suffix, scales them down to the configured replica count, adds a `pipecd.dev/variant=canary` label, and applies them to every target cluster in parallel. The original deployment is never touched this stage only ever adds resources.

![Canary rollout stage log applying manifests to cluster-eu and cluster-us](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/swcu1ppt38ltw87wwbol.png)

![Canary rollout success — deploy targets: cluster-eu + cluster-us](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/0lfhrddbu6r3mt01tlrs.png)

---

### K8S_CANARY_CLEAN

Once the canary window is over, whether you promoted or rolled back, the canary pods are just sitting in every cluster doing nothing. `K8S_CANARY_CLEAN` removes them.

It finds all resources with the label `pipecd.dev/variant=canary` for the application and deletes them in order: Services first, then Deployments, then everything else. The order matters as you don't want to remove the Deployment while the Service is still sending traffic to it.

One thing worth noting: the query is scoped strictly to canary-labelled resources. Even if something goes wrong in the deletion logic, it cannot touch primary resources.

![K8S_CANARY_CLEAN stage log deleting simple-canary resources from both clusters](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/6vzb1fikt47ax3z53bjy.png)

![K8S_CANARY_ROLLOUT → K8S_CANARY_CLEAN pipeline — both stages green on cluster-eu and cluster-us](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/5980a6b46dtiqv9oetz7.png)

---

### K8S_PRIMARY_ROLLOUT

After the canary looks good, you promote the new version to primary, the workload actually serving all your users. This stage takes the manifests from Git, adds the `pipecd.dev/variant=primary` label, and applies them across all clusters in parallel.

It also has a `prune` option: after applying, it checks what's currently running in the cluster against what was just applied, and deletes anything that's no longer in Git. Useful when you remove a resource from your manifests and want the cluster to reflect that.

![K8S_PRIMARY_ROLLOUT success deploy targets: cluster-eu + cluster-us](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/5i4pfqi10ef38i7ltn55.png)

![kubectl confirming simple 2/2 updated in both cluster-eu and cluster-us](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/kbi6iplj1l8g5wbkwgql.png)

---

### K8S_BASELINE_ROLLOUT

This one took me a while to understand and it is the stage I find most interesting to explain as well.

When you're running a canary, the natural thing is to compare it against primary. The issue is that's not a fair comparison primary is handling far more traffic than canary, under different conditions.

Baseline gives you a fairer comparison. You take the *current* version (not the new one) and run it at the same scale as canary. Now your cluster has:

```plaintext
simple             2/2   ← production, current version
simple-canary      1/1   ← new version, being tested
simple-baseline    1/1   ← current version at canary scale
```

You compare canary vs baseline, same number of pods, same traffic conditions. If canary is worse, it's obvious.

The key difference from every other rollout stage is one line of code. Canary and primary load manifests from the new Git commit (`TargetDeploymentSource`). Baseline loads from what's currently running (`RunningDeploymentSource`):

```go
// canary.go — new version
manifests, err := p.loadManifests(ctx, ..., &input.Request.TargetDeploymentSource, ...)

// baseline.go — current version
manifests, err := p.loadManifests(ctx, ..., &input.Request.RunningDeploymentSource, ...)
```

![K8S_BASELINE_ROLLOUT stage log loading manifests from running deployment source](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/a61aeapcwnqquh3v3vdh.png)

![K8S_BASELINE_ROLLOUT](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/0r26i800uc46kfmauo5p.png)

![K8S_BASELINE_ROLLOUT](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/y6lowk38efysvvy0vbmk.png)

![kubectl showing simple, simple-baseline, simple-canary all running in both clusters](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/8n8ailmbo0gae31sf58l.png)

---

### K8S_BASELINE_CLEAN

Once the analysis is done, baseline resources get cleaned up the same way as canary find everything labelled `pipecd.dev/variant=baseline` and delete it in order. No configuration needed. It doesn't matter whether `createService: true` was set during rollout, it finds whatever is there and removes it.

![K8S_BASELINE_CLEAN](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/ddfhlxgi1i40r0x8dgh2.png)

![K8S_BASELINE_CLEAN stage log deleting baseline resources from both clusters](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/fsxadrvt48rbr4hi113m.png)

![K8S_BASELINE_CLEAN stage log deleting baseline resources from both clusters](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/apa1tdqa751bgl2t9g6c.png)

![K8S_BASELINE_CLEAN](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/rder7h0ylhxn0wkaykep.png)

![kubectl confirming no baseline resources remain in cluster-eu or cluster-us](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/8epcyha332ir2xz3jhua.png)

---

### K8S_TRAFFIC_ROUTING

Canary and baseline pods exist in the cluster but get no traffic until this stage runs. Without it, you're analysing pods that nobody is actually hitting. This stage is what sends real user traffic to them.

Two methods are supported:

**PodSelector** (no service mesh needed): changes the Kubernetes Service selector to point at one variant. All-or-nothing 100% to canary or 100% back to primary.

![PodSelector traffic routing full pipeline success across cluster-eu and cluster-us](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/7uqb18dwo6jyhr0aekpw.png)

![PodSelector](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/qddkmg9qs6unwk0bklu3.png)

![PodSelector](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/piwbmpywrxmf4ykbcomp.png)

![PodSelector](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/2u1kqy9jasqhqwtwebbd.png)

**Istio**: updates VirtualService route weights to split traffic across all three variants at once for example, primary 80%, canary 10%, baseline 10%. Also supports `editableRoutes` to limit which named routes the stage is allowed to modify.

One small thing I added on top of the traffic routing stage: per-route logging. When the stage runs, it now logs each route it processes whether it was skipped (because it's not in `editableRoutes`) or updated with new weights. Before this, the log just said "Successfully updated traffic routing" with no detail. Now you can see exactly which routes changed and to what percentages, which is useful when debugging a misconfigured VirtualService.

![Istio traffic routing stage log per-route logging showing which routes were updated in both clusters](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/498uvcytrlppjpxqcq05.png)

![Istio](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/6ql3nmfb19psi5gfu8a2.png)

![Istio](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/oxy1gv5cojotd6ks253m.png)

![Istio](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/tva9uwcmd7prb6qf72dm.png)

![Full Istio pipeline, all 7 stages green on cluster-eu and cluster-us](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/6oj4elm5asfgergkp2se.png)

---

## Something I Found Interesting

The thing that surprised me was how `errgroup` handles running across multiple clusters without much extra code.

Every stage needs to run against N clusters, not one. A simple for-loop would run them one at a time slow, and if cluster 2 fails you don't find out until cluster 1 is already done.

`errgroup` runs all clusters at the same time and returns the first error:

```go
eg, ctx := errgroup.WithContext(ctx)
for _, tc := range targetClusters {
    tc := tc
    eg.Go(func() error {
        return canaryRollout(ctx, tc.deployTarget, ...)
    })
}
return eg.Wait()
```

All clusters run in parallel. If any one fails, the stage fails immediately. The same pattern is used across every stage, so adding a new stage is mostly just writing the per-cluster logic the concurrency part is already solved.

## What's Next

The next piece is `DetermineStrategy`, that is the logic that decides what kind of deployment to trigger based on what changed in Git. After that, livestate drift detection so PipeCD can flag when a cluster has drifted from what Git says it should be.

To get involved, check out the PipeCD project and come join us on Slack.

## Links

- [PipeCD repository](https://github.com/pipe-cd/pipecd)
- [LFX Mentorship Program](https://mentorship.lfx.linuxfoundation.org)
- [Issue #6446, kubernetes_multicluster plugin](https://github.com/pipe-cd/pipecd/issues/6446)
- [PR #6629 K8S_TRAFFIC_ROUTING](https://github.com/pipe-cd/pipecd/pull/6629)
- [PR #6648 Per-route logging in K8S_TRAFFIC_ROUTING](https://github.com/pipe-cd/pipecd/pull/6648)
- [Slack #PipeCD](https://app.slack.com/client/T08PSQ7BQ/C01B27F9T0X)
