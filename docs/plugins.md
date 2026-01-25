# PipeCD Official Plugins

This document lists all official PipeCD plugins with their latest released versions.

**Last updated:** 2026-01-24

---

## Quick Reference

| Plugin | Latest Version | Repository | Documentation |
|--------|----------------|------------|---|
| [Kubernetes](#kubernetes-plugin) | v0.1.0 | [GitHub](https://github.com/pipe-cd/pipecd) | [Docs](https://pipecd.dev/docs/user-guide/managing-deployments/deployment-pipeline/#kubernetes-deployment) |
| [Terraform](#terraform-plugin) | v0.1.0 | [GitHub](https://github.com/pipe-cd/pipecd) | [Docs](https://pipecd.dev/docs/user-guide/managing-deployments/deployment-pipeline/#terraform-deployment) |
| [Cloud Run](#cloud-run-plugin) | v0.1.0 | [GitHub](https://github.com/pipe-cd/pipecd) | [Docs](https://pipecd.dev/docs/user-guide/managing-deployments/deployment-pipeline/#cloud-run-deployment) |
| [Wait Stage](#wait-stage-plugin) | v0.1.0 | [GitHub](https://github.com/pipe-cd/pipecd) | [Docs](https://pipecd.dev/docs/user-guide/managing-deployments/deployment-pipeline/#wait-stage) |
| [Wait Approval](#wait-approval-stage-plugin) | v0.1.0 | [GitHub](https://github.com/pipe-cd/pipecd) | [Docs](https://pipecd.dev/docs/user-guide/managing-deployments/deployment-pipeline/#wait-approval-stage) |
| [Script Run](#script-run-plugin) | v0.1.0 | [GitHub](https://github.com/pipe-cd/pipecd) | [Docs](https://pipecd.dev/docs/user-guide/managing-deployments/deployment-pipeline/#script-run-stage) |
| [Analysis](#analysis-plugin) | v0.1.0 | [GitHub](https://github.com/pipe-cd/pipecd) | [Docs](https://pipecd.dev/docs/user-guide/managing-deployments/deployment-pipeline/#analysis-stage) |
| [Kubernetes Multi-cluster](#kubernetes-multi-cluster-plugin) | v0.1.0 | [GitHub](https://github.com/pipe-cd/pipecd) | [Docs](https://pipecd.dev/docs/user-guide/managing-deployments/deployment-pipeline/#kubernetes-multicluster-deployment) |
| [Plugin SDK for Go](#pipecd-plugin-sdk-for-go) | v0.3.0 | [GitHub](https://github.com/pipe-cd/piped-plugin-sdk-go) | [Repo](https://github.com/pipe-cd/piped-plugin-sdk-go) |

---

## Plugin Details

### Kubernetes Plugin

Deploy and manage applications on Kubernetes clusters with declarative, GitOps-driven deployments.

- **Latest Version:** v0.1.0
- **Release URL:** https://github.com/pipe-cd/pipecd/releases/tag/pkg/app/pipedv1/plugin/kubernetes/v0.1.0
- **Source:** [pkg/app/pipedv1/plugin/kubernetes](../../pkg/app/pipedv1/plugin/kubernetes)
- **Status:** Stable

### Terraform Plugin

Deploy and manage infrastructure as code using Terraform with automated drift detection and remediation.

- **Latest Version:** v0.1.0
- **Release URL:** https://github.com/pipe-cd/pipecd/releases/tag/pkg/app/pipedv1/plugin/terraform/v0.1.0
- **Source:** [pkg/app/pipedv1/plugin/terraform](../../pkg/app/pipedv1/plugin/terraform)
- **Status:** Stable

### Cloud Run Plugin

Deploy services to Google Cloud Run serverless platform.

- **Latest Version:** v0.1.0
- **Release URL:** https://github.com/pipe-cd/pipecd/releases/tag/pkg/app/pipedv1/plugin/cloudrun/v0.1.0
- **Source:** [pkg/app/pipedv1/plugin/cloudrun](../../pkg/app/pipedv1/plugin/cloudrun)
- **Status:** Stable

### Wait Stage Plugin

Add timed delay stages to deployment pipelines for staged rollouts and coordinated deployments.

- **Latest Version:** v0.1.0
- **Release URL:** https://github.com/pipe-cd/pipecd/releases/tag/pkg/app/pipedv1/plugin/wait/v0.1.0
- **Source:** [pkg/app/pipedv1/plugin/wait](../../pkg/app/pipedv1/plugin/wait)
- **Status:** Stable

### Wait Approval Stage Plugin

Add manual approval gates to deployment pipelines for controlled rollouts.

- **Latest Version:** v0.1.0
- **Release URL:** https://github.com/pipe-cd/pipecd/releases/tag/pkg/app/pipedv1/plugin/waitapproval/v0.1.0
- **Source:** [pkg/app/pipedv1/plugin/waitapproval](../../pkg/app/pipedv1/plugin/waitapproval)
- **Status:** Stable

### Script Run Plugin

Execute custom shell scripts as part of deployment pipelines for flexible stage definitions.

- **Latest Version:** v0.1.0
- **Release URL:** https://github.com/pipe-cd/pipecd/releases/tag/pkg/app/pipedv1/plugin/scriptrun/v0.1.0
- **Source:** [pkg/app/pipedv1/plugin/scriptrun](../../pkg/app/pipedv1/plugin/scriptrun)
- **Status:** Stable

### Analysis Plugin

Analyze deployment metrics, logs, and performance to determine deployment success and health.

- **Latest Version:** v0.1.0
- **Release URL:** https://github.com/pipe-cd/pipecd/releases/tag/pkg/app/pipedv1/plugin/analysis/v0.1.0
- **Source:** [pkg/app/pipedv1/plugin/analysis](../../pkg/app/pipedv1/plugin/analysis)
- **Status:** Stable

### Kubernetes Multi-cluster Plugin

Deploy applications across multiple Kubernetes clusters with coordinated rollouts.

- **Latest Version:** v0.1.0
- **Release URL:** https://github.com/pipe-cd/pipecd/releases/tag/pkg/app/pipedv1/plugin/kubernetes_multicluster/v0.1.0
- **Source:** [pkg/app/pipedv1/plugin/kubernetes_multicluster](../../pkg/app/pipedv1/plugin/kubernetes_multicluster)
- **Status:** Stable

### PipeCD Plugin SDK for Go

Official SDK for developing custom PipeCD plugins using Go. Build plugins that extend PipeCD with new deployment platforms and stages.

- **Latest Version:** v0.3.0
- **Release URL:** https://github.com/pipe-cd/piped-plugin-sdk-go/releases/tag/v0.3.0
- **Repository:** [pipe-cd/piped-plugin-sdk-go](https://github.com/pipe-cd/piped-plugin-sdk-go) (External)
- **Source Mirror:** [pkg/plugin/sdk](../../pkg/plugin/sdk)
- **Status:** Stable

---

## How Plugin Versions Are Tracked

- **Inline plugins** (in pipecd repo): Released as Git tags with format `pkg/app/pipedv1/plugin/{name}/v{version}`
  - Example: `pkg/app/pipedv1/plugin/kubernetes/v0.1.0`
  - Release page: https://github.com/pipe-cd/pipecd/releases

- **External plugins** (separate repos): Released in their own repositories
  - Example: `piped-plugin-sdk-go` uses tags `v{version}`
  - Release page: https://github.com/pipe-cd/piped-plugin-sdk-go/releases

## Updating This Registry

This registry is **automatically updated** by a GitHub Actions workflow on:
- Every new plugin release
- Every 6 hours (scheduled)

Changes are committed to the repository when new versions are detected.

For details on the automation, see [.github/workflows/update-plugins-registry.yaml](../../.github/workflows/update-plugins-registry.yaml).

## Related Documentation

- **Plugin Architecture RFC:** [docs/rfcs/0015-pipecd-plugin-arch-meta.md](../rfcs/0015-pipecd-plugin-arch-meta.md)
- **Plugin Development Guide:** [PipeCD Documentation](https://pipecd.dev/docs/developer-guide/plugin-development/)
- **Plugin Release Process:** [.github/workflows/plugin_release.yaml](../../.github/workflows/plugin_release.yaml)
- **All Releases:** https://github.com/pipe-cd/pipecd/releases

---

**Note:** This registry is machine-generated. To update version information or add new plugins, modify the registry generation scripts or submit an issue/PR to the [PipeCD repository](https://github.com/pipe-cd/pipecd).
