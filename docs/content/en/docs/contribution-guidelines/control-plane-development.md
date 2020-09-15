---
title: "Control plane development"
linkTitle: "Control plane development"
weight: 5
description: >
  This page describes where to find control-plane source code and how to run it locally for debugging.
---

## Source code structure

- [pkg/app/api](https://github.com/pipe-cd/pipe/tree/master/pkg/app/api): contains source code for control-plane api.
- [cmd/api](https://github.com/pipe-cd/pipe/tree/master/cmd/api): entrypoint for binary of control-plane api.

- [pkg/app/web](https://github.com/pipe-cd/pipe/tree/master/pkg/app/web): contains source code for control-plane web.
- [cmd/web](https://github.com/pipe-cd/pipe/tree/master/cmd/web): entrypoint for binary of control-plane web.

- [pkg](https://github.com/pipe-cd/pipe/tree/master/pkg): contains shared source code for all components of both `piped` and `control-plane`.

## How to run API locally

Prepare a ControlPlane configuration file as described at [Installation](https://pipecd.dev/docs/operator-manual/control-plane/installation/) and start running by the following command:

``` console
bazelisk run //cmd/api:api -- server \
  --config-file=/your-path-to-path/control-plane.yaml
```

## How to run Web locally

Run the following command to start http server for web.

```
bazelisk run //cmd/web:web server
```

Then access `http://localhost:9082` on your browser.
