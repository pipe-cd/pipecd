---
title: "Piped development"
linkTitle: "Piped development"
weight: 4
description: >
  This page describes where to find piped source code and how to run it locally for debugging.
---

## Source code structure

- [pkg/app/piped](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/piped): contains source code for only `piped`.
- [cmd/piped](https://github.com/pipe-cd/pipecd/tree/master/cmd/piped): entrypoint for `piped` binary.
- [pkg](https://github.com/pipe-cd/pipecd/tree/master/pkg): contains shared source code for all components of both `piped` and `control-plane`.

## How to run it locally

1. Prepare the piped configuration file `piped-config.yaml`

2. Ensure that your `kube-context` is connecting to the right kubernetes cluster

2. Run the following command to start running `piped`

``` console
bazelisk run --run_under="cd $PWD && " //cmd/piped:piped -- piped \
--tools-dir=/tmp/piped-bin \
--config-file=piped-config.yaml
```

## How to run it locally as docker container

``` bash
# Compile the current source code to build a new Docker image
# and then load it into the local docker client as bazel/cmd/piped:image.
make load-piped-image

docker run bazel/cmd/piped:image --help
```
