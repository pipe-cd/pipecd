---
title: "Piped development"
linkTitle: "Piped development"
weight: 4
description: >
  This page describes where to find piped source code and how to run it locally for debugging.
---

## Source code structure

- [pkg/app/piped](https://github.com/pipe-cd/pipe/tree/master/pkg/app/piped): contains source code for only `piped`.
- [cmd/piped](https://github.com/pipe-cd/pipe/tree/master/cmd/piped): entrypoint for `piped` binary.
- [pkg](https://github.com/pipe-cd/pipe/tree/master/pkg): contains shared source code for all components of both `piped` and `control-plane`.

## How to run it locally

In order to run `piped` at the local environment for debugging while development without connecting to a real control plane,
we prepared a fake control-plane to be used.

1. Prepare a `.dev` directory at the root of repository that contains:
- A `piped.key` file with a fake key as the following
  ```
  hello-pipecd
  ```

- A piped configuration file `piped-config.yaml`

2. Ensure that your `kube-context` is connecting to the right kubernetes cluster

2. Run the following command to start running `piped`

``` console
bazelisk run --run_under="cd $PWD && " //cmd/piped:piped -- piped \
--use-fake-api-client=true \
--bin-dir=/tmp/piped-bin \
--config-file=.dev/dev-config.yaml
```
