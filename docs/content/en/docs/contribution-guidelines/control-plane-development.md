---
title: "Control plane development"
linkTitle: "Control plane development"
weight: 5
description: >
  This page describes where to find control-plane source code and how to run it locally for debugging.
---

## Source code structure

- [cmd/server](https://github.com/pipe-cd/pipe/tree/master/cmd/server): entrypoint for binary of control-plane server.

- [pkg/app/api](https://github.com/pipe-cd/pipe/tree/master/pkg/app/api): contains source code for control-plane api.
- [pkg/app/web](https://github.com/pipe-cd/pipe/tree/master/pkg/app/web): contains source code for control-plane web.
- [pkg](https://github.com/pipe-cd/pipe/tree/master/pkg): contains shared source code for all components of both `piped` and `control-plane`.

## How to run server locally

Prepare a ControlPlane configuration file as described at [Installation](https://pipecd.dev/docs/operator-manual/control-plane/installation/) and start running by the following command:

``` console
bazelisk run //cmd/server:server -- server \
--config-file=absolute-path-to-control-plane-config.yaml \
--encryption-key-file=absolute-path-to-a-random-key-file
```

Then go to `http://localhost:9082` on your browser to access PipeCD's web.
