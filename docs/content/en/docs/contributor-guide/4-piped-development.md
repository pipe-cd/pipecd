---
title: "Piped Development"
linkTitle: "Piped Development"
weight: 4
description: >
  This page describes where to find piped source code and how to run it locally for debugging.
---

## Source code structure

- [pkg/app/piped](https://github.com/pipe-cd/pipe/tree/master/pkg/app/piped): Contains source code for only `piped`. 
- [cmd/piped](https://github.com/pipe-cd/pipe/tree/master/cmd/piped): Entrypoint for `piped` binary.
- [pkg](https://github.com/pipe-cd/pipe/tree/master/pkg): Contains shared source code for all components of both `piped` and `control-plane`.

## How to run it locally

In order to run `piped` at the local environment for debugging while development without connecting to a real control plane,
we prepared a fake control-plane to be used.

1. Prepare a `.dev` directory at the root of repository that contains:
- a `piped.key` file with a fake key (e.g. `"hello-pipecd"`)
- a piped configuration file `piped-config.yaml` (e.g. [pkg/config/testdata/piped/dev-config.yaml](https://github.com/pipe-cd/pipe/blob/master/pkg/config/testdata/piped/dev-config.yaml))

2. Ensure that your `kube-context` is connecting to right kubernetes cluster

2. Run the following command to start running `piped`

``` console
bazelisk run --run_under="cd $PWD && " //cmd/piped:piped -- piped \
--log-encoding=humanize \
--use-fake-api-client=true \
--project-id=local-dev-project \
--piped-id=local-dev-piped \
--piped-key-file=.dev/piped.key \
--bin-dir=/tmp/piped-bin \
--config-file=pkg/config/testdata/piped/dev-config.yaml
```
