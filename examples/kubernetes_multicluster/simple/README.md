# Kubernetes Multi-Cluster Simple Example

This example shows the first pipedv1 Kubernetes multicluster app in PipeCD. It deploys the same nginx manifests to two clusters through the kubernetes-multicluster plugin, using `cluster-a` and `cluster-b` as deploy targets.

`K8S_MULTI_SYNC` performs a quick sync across every selected cluster at the same time. It applies the manifests in this directory to each target and, with `prune: true`, removes resources that are no longer tracked in Git.

## Prerequisites

- A pipedv1 setup with the `kubernetes-multicluster` plugin configured.
- Matching deploy targets named `cluster-a` and `cluster-b`.
- See [.piped/config-example.yaml](.piped/config-example.yaml) for the expected plugin config shape.

## What's different from v0 examples

v1 runs the kubernetes-multicluster plugin as a separate gRPC process registered in Piped config, instead of embedding the logic in the main binary like the older examples.
