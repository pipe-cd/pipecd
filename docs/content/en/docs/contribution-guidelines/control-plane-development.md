---
title: "Control plane development"
linkTitle: "Control plane development"
weight: 5
description: >
  This page describes where to find control-plane source code and how to run it locally for debugging.
---

## Source code structure

- [cmd/pipecd](https://github.com/pipe-cd/pipe/tree/master/cmd/pipecd): entrypoint for binary of control-plane server.
- [pkg/app/api](https://github.com/pipe-cd/pipe/tree/master/pkg/app/api): contains source code for control-plane api.
- [pkg/app/web](https://github.com/pipe-cd/pipe/tree/master/pkg/app/web): contains source code for control-plane web.
- [pkg](https://github.com/pipe-cd/pipe/tree/master/pkg): contains shared source code for all components of both `piped` and `control-plane`.

## How to run server locally

### Running up server

Prepare a ControlPlane configuration file as described at [Installation](https://pipecd.dev/docs/operator-manual/control-plane/installation/) and start running by the following command:

``` console
bazelisk run //cmd/pipecd -- server \
--config-file=absolute-path-to-control-plane-config.yaml \
--encryption-key-file=absolute-path-to-a-random-key-file
```

Because we are using **grpc-web** for communicating between web-client and server, so we may need a **local Envoy instance**.

### Integrating with Envoy

You can install [Envoy](https://www.envoyproxy.io/docs/envoy/latest/start/install) locally or running it on [Docker](https://docs.docker.com/get-docker/).

We already prepared `local/envoy-config.yaml` **but syntax is running on _Envoy_ version _1.10.0_**

_At root of repository:_

#### Running local Envoy

```
envoy -c local/envoy-config.yaml 
```

#### Running Envoy on Docker

1. Run Docker command build for `local/Dockerfile`

```
docker build -t web-envoy:v1 ./local
```

3. And now you can execute it with:

```
docker run -d --net="host" --name web-envoy -p 9095:9095 -p 9090:9090 web-envoy:1.0
```

Then go to `http://localhost:9090` on your browser to access PipeCD's web.