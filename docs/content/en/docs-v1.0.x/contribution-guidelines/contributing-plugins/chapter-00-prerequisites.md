---
title: "Prerequisites and Setup"
linkTitle: "Prerequisites"
weight: 2
description: >
  What you need and how to set up a local PipeCD development environment.
---

← Back to Book Index: ./_index.md

# Prerequisites and Setup

This chapter covers what you should know and the tools you need before starting to build a PipeCD plugin, plus quick steps to set up a local development environment.

## 1 — What You Need to Know Before Starting

- Go programming: be comfortable reading and writing Go structs, interfaces, and functions. If you need a refresher, see https://go.dev/tour
- Kubernetes basics: understand core objects such as Deployments and Services and how manifests describe them
- Basic understanding of what a CD (continuous delivery) platform does and the role of a control plane vs. agents/plugins
- Familiarity with gRPC is helpful but not required — the PipeCD plugin SDK abstracts most of the gRPC boilerplate

## 2 — Required Tools

- Go 1.22 or later — download from https://go.dev/dl/
- kubectl 1.32 or later — https://kubernetes.io/docs/tasks/tools/
- Docker — for running a local PipeCD control plane when you want to test end-to-end
- git — for source control and contributing
- A text editor or IDE with Go support (VS Code with the Go extension is recommended)

## 3 — Setting Up a Local PipeCD Environment (high level)

- Clone the PipeCD repository:

```bash
git clone https://github.com/pipe-cd/pipecd.git
cd pipecd
```

- The `pipedv1`-related binaries live under `cmd/` (for example `cmd/pipedv1/`)
- Official plugins live under `pkg/app/pipedv1/plugin/`
- To build a plugin locally, change to the plugin directory and run `go build`:

```bash
cd pkg/app/pipedv1/plugin/kubernetes
go build ./...
```

- Refer to the piped configuration schema at https://github.com/pipe-cd/pipecd/blob/master/pkg/configv1/piped.go to understand how plugins are registered in `piped` config
- For full installation and local runtime instructions, see the official docs: https://pipecd.dev/docs

## 4 — Verifying Your Setup

Checklist you can run through locally:

- [ ] `go version` shows Go 1.22 or later
- [ ] `kubectl version --client` shows kubectl 1.32 or later
- [ ] You can clone the pipecd repository and `cd cmd/pipedv1/`
- [ ] You can build the kubernetes plugin:

```bash
cd pkg/app/pipedv1/plugin/kubernetes && go build ./...
```
- [ ] You can open a proto file and recognize a gRPC service definition, for example: https://github.com/pipe-cd/pipecd/blob/master/pkg/plugin/api/v1alpha1/deployment/api.proto

## 5 — A Note on the Plugin SDK

- The plugin SDK is located in `pkg/plugin/sdk/` inside the pipecd repository and provides helpers for building a gRPC server and wiring your implementation into `piped`.
- Plugins are standalone Go binaries that implement one or more interfaces (Deployment, Livestate, etc.). The SDK handles the gRPC server setup; you implement the business logic.
- Key SDK types and entry points live in `pkg/plugin/sdk/plugin.go` — consult the file in the repository for API details.

---

Next Chapter → ./chapter-01-introduction.md
