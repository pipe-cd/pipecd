---
title: "Development"
linkTitle: "Development"
weight: 2
description: >
  This page describes how to build, test PipeCD source code at your local environment.
---

## Prerequisites

Only `bazel` is required to build and test `PipeCD` project, but you don't need to install `bazel` directly.
Instead of that we are using [`bazelisk`](https://github.com/bazelbuild/bazelisk) for automatically picking a good version of Bazel for building `PipeCD` project.

You can install `bazelisk` by `go get` as the following:
```
go get github.com/bazelbuild/bazelisk
```

or directly install its [binary](https://github.com/bazelbuild/bazelisk/releases) from the release page.
For more information, you might want to read the [installation notes of `bazelisk`](https://github.com/bazelbuild/bazelisk#requirements).

## Repositories
- [pipe](https://github.com/pipe-cd/pipe): contains all source code and documentation of PipeCD project.
- [manifests](https://github.com/pipe-cd/manifests): contains all automatically generated release manifests for both `piped` and `control-plane` components.
- [examples](https://github.com/pipe-cd/examples): contains various generated examples to demonstrate how to use PipeCD.

## Build and test with Bazel

- `make build`: builds all binaries in tree.
- `make test`: runs all unit tests.
- `make dep`: updates `go.mod` and bazel `WORKSPACE`. Run this command after adding a new go dependency or update the version of a dependency.
- `make gazelle`: generates `BUILD.bazel` files for go code. Run this command after adding a new `import` in go code or adding a new go file.
- `make buildifier`: formats bazel BUILD and .bzl files with a standard convention.
- `make clean`: cleans all bazel cache.
- `make expose-generated-go`: exposes generated Go files (`.pb.go`, `.mock.go`...) to editors and IDEs.
- `make site`: runs PipeCD site (https://pipecd.dev) locally (requires [hugo](https://github.com/gohugoio/hugo) with `_extended` version `0.88.1` or later to be installed).

**NOTE**: The first time of running a bazel command will take some minutes because bazel needs to download all required dependencies. From the second time it will be very fast.

## Docs and workaround with docs

PipeCD official site contains multiple versions documentation, all placed under `/docs/content/en` directory, which are:
- `/docs`: stable version docs, usually synced with the latest released version docs.
- `/docs-dev`: experimental version docs, contains docs for not yet released features or changes.
- `/docs-v0.x.x`: contains docs for specified version family (a version family is all versions which in the same major release).

Basically, we have two simple rules:
- Do not touch to the `/docs` content directly.
- Keep stable docs version synced with the latest released docs version.

Here is the flow of docs contribution regard some known scenarios:
1. Update docs that are related to a specified version (which is not the latest released version):
In such case, update the docs under `/docs-v0.x.x` is enough.
2. Update docs for not yet released features or changes:
In such case, update the docs under `/docs-dev` is enough.
3. Update docs that are related to the latest released docs version:
- Change the docs' content that fixes the issue under `/docs-dev` and `/docs-v0.x.x`, they share the same file structure so you should find the right files in both directories.
- Use `make sync-stable-docs` command to sync the latest released version docs under `/docs-v0.x.x` to `/docs`

If you find any issues related to the docs, we're happy to accept your help.
