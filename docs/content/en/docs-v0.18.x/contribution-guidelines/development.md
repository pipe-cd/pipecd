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
- [pipecd](https://github.com/pipe-cd/pipecd): contains all source code and documentation of PipeCD project.
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
