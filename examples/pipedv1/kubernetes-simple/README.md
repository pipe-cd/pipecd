# Simple Kubernetes App - PipeCD v1

This example demonstrates PipeCD v1 (pipedv1) with plugin architecture.

## Key Differences from piped (v0)

| Feature | v0 | v1 |
|---------|----|----|
| File name | `app.pipecd.yaml` | `.app.yaml` |
| Kind | `KubernetesApp` | `Application` |
| Config | `input:` section | `plugins:` section |
| Deployment | Implicit quick sync | Explicit `pipeline:` stages |

## How to use

1. Configure piped with kubernetes plugin enabled
2. Point PipeCD to this directory
3. Application will sync using K8S_SYNC stage
