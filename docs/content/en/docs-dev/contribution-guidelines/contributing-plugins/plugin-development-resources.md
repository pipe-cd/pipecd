---
title: "Plugin development resources"
linkTitle: "Plugin development resources"
weight: 5
description: >
  Links and short notes for developing PipeCD plugins.
---

> **Note:**
> This section is still a work in progress. A full tutorial and an in-docs translation of the Zenn book are planned over time.
>
> For a hands-on walkthrough, read [**Build and learn PipeCD plugins**](https://zenn.dev/warashi/books/try-and-learn-pipecd-plugin) (Zenn; Japanese title *作って学ぶ PipeCD プラグイン*). The same book is linked again in [Links](#links), together with other references. Use your browser's translate feature to read this in English. Verify commands and field names against this documentation and the [`pipecd`](https://github.com/pipe-cd/pipecd) repository.

Use this page together with [Contribute to PipeCD plugins](../), which covers layout, `make` targets, and how to open a pull request.

## How the pieces fit together

- **Control plane**: registers projects and `piped` agents; the web UI shows deployment status.
- **`piped`**: runs your plugins as separate binaries and talks to them over gRPC.
- **Plugins**: carry out deployments (and optionally live state and drift) for a platform or tool.

If you are new to PipeCD v1, read [Migrating from v0 to v1](/docs-dev/migrating-from-v0-to-v1/).

### Terms

| Term | Meaning |
|------|---------|
| **Application** | Git content for one deployable unit: manifests plus a `*.pipecd.yaml` file. |
| **Deployment** | One run of the deployment pipeline for an app (from Git, a trigger, or the UI). |
| **Deploy target** | Where the plugin deploys, set under `spec.plugins` in the `piped` config. |
| **Pipeline** | Ordered **stages** (for example sync, canary, wait) from the application config. |
| **Stage** | One step in the pipeline; your plugin implements the stages it supports. |

For plugin interfaces (**Deployment**, **LiveState**, **Drift**), see [Plugin types](../#plugin-types). A first plugin usually implements **Deployment** only.

## Local development

Use a `piped` v1 build that matches your work. From the repo you can run `make run/piped` as in the [contributing guide](../#build-and-test). To install a release binary, see [Installing on a single machine](/docs-dev/installation/install-piped/installing-on-single-machine/).

Example when running the v1 `piped` CLI:

```console
./piped run --config-file=PATH_TO_PIPED_CONFIG --insecure=true
```

Use `--insecure=true` only when it matches your control plane (for example plain HTTP or your dev TLS settings).

The install guide linked above uses the same `run` subcommand. If another page or tutorial shows different syntax, run `./piped --help` on your binary to match your build.

Older blog posts or books may target an older `piped`. Prefer this documentation and the default branch of [`pipecd`](https://github.com/pipe-cd/pipecd).

## `piped` config and application config

You need both:

1. **`piped` configuration**: control plane address, Git `repositories`, and `spec.plugins` (URL, port, `deployTargets`, plugin-specific fields). See [Piped configuration reference](/docs-dev/user-guide/managing-piped/configuration-reference/).
2. **Application configuration**: the `*.pipecd.yaml` in Git (plugin, pipeline, deploy options). See [Adding an application](/docs-dev/user-guide/managing-application/adding-an-application/).

Code in your plugin reads the application config via its own types (often under `config/`). `deployTargets` come from the `piped` config.

## Example plugins

| Plugin | Notes |
|--------|-------|
| [kubernetes](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin/kubernetes) | Full official plugin |
| [wait](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin/wait) | Small official example |
| [example-stage](https://github.com/pipe-cd/community-plugins/tree/main/plugins/example-stage) | Community sample (see [Installing on a single machine](/docs-dev/installation/install-piped/installing-on-single-machine/)) |

## Cache under ~/.piped

After rebuilding a plugin, `piped` may still use files under **`~/.piped`** (including **`~/.piped/plugins`**). If your changes do not show up, remove those directories or clear the cache your team uses, then restart `piped`.

## Debugging

- **Web UI**: deployment and stage status.
- **Stdout**: logs from `piped` and the plugin (verbose but quick for local work).
- **Stage logs**: through the SDK; see [`StageLogPersister`](https://github.com/pipe-cd/pipecd/blob/master/pkg/plugin/sdk/logpersister/persister.go) in the repo.

## Links

| Resource | Notes |
|----------|-------|
| [**Build and learn PipeCD plugins** (Zenn)](https://zenn.dev/warashi/books/try-and-learn-pipecd-plugin) | Japanese tutorial book (*作って学ぶ PipeCD プラグイン*) |
| [DeepWiki (pipecd)](https://deepwiki.com/pipe-cd/pipecd) | Unofficial wiki-style overview of the repository. |
| [Plugin Architecture RFC](https://github.com/pipe-cd/pipecd/blob/master/docs/rfcs/0015-pipecd-plugin-arch-meta.md) | Design (RFC) |
| [Plugin Architecture overview (blog)](https://pipecd.dev/blog/2024/11/28/overview-of-the-plan-for-pluginnable-pipecd/) | Architecture overview |
| [Plugin alpha release (blog)](https://pipecd.dev/blog/2025/06/16/plugin-architecture-piped-alpha-version-has-been-released/) | Piped plugin alpha |
| [#pipecd (CNCF Slack)](https://cloud-native.slack.com/) | Community chat |

See also [Contributing to PipeCD](/docs-dev/contribution-guidelines/contributing/) for local control plane setup and pull requests.
