---
title: "Control Plane development"
linkTitle: "Control Plane development"
weight: 5
description: >
  This page describes where to find Control Plane source code and how to run it locally for debugging.
---

## Source code structure

- [cmd/pipecd](https://github.com/pipe-cd/pipecd/tree/master/cmd/pipecd): entrypoint for binary of Control Plane server.
- [pkg/app/server](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/server): contains source code for Control Plane api.
- [web](https://github.com/pipe-cd/pipecd/tree/master/web): contains source code for the web console.
- [pkg](https://github.com/pipe-cd/pipecd/tree/master/pkg): contains shared source code for all components of both `Piped` and `Control Plane`.

## How to run Control Plane locally

### Prerequisites
- Installing [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)

### Start running a Kubernetes cluster

``` console
make kind-up
```

Once it is no longer used, run `make kind-down` to delete it.

### Installing Control Plane into the local cluster

``` console
make run/pipecd
```

Once all components are running up, use `kubectl port-forward` to expose the installed Control Plane on your localhost:

``` console
kubectl -n pipecd port-forward svc/pipecd 8080
```

Point your web browser to [http://localhost:8080](http://localhost:8080) to login with the configured static admin account: project = `quickstart`,
username = `hello-pipecd`, password = `hello-pipecd`.
