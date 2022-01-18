---
title: "Control plane development"
linkTitle: "Control plane development"
weight: 5
description: >
  This page describes where to find control-plane source code and how to run it locally for debugging.
---

## Source code structure

- [cmd/pipecd](https://github.com/pipe-cd/pipecd/tree/master/cmd/pipecd): entrypoint for binary of control-plane server.
- [pkg/app/server](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/server): contains source code for control-plane api.
- [pkg/app/web](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/web): contains source code for control-plane web.
- [pkg](https://github.com/pipe-cd/pipecd/tree/master/pkg): contains shared source code for all components of both `piped` and `control-plane`.

## How to run control-plane locally

### Prerequisites
- Installing [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
- Running a local Kubernetes cluster by `make kind-up`. (Cluster that is no longer used can be deleted by `make kind-down`.)

### Pushing the images to local container registry

``` console
make push
```

**NOTE: Since it uses the commit hash as an image tag, you need to make a new commit whenever you change the source code.**

This command compiles the local source code to build the docker images for `pipecd`, `piped` and then pushes them to the local container registery which was enabled by `make kind-up`.

### Rendering the local manifests

Because the `manifests` directory at [pipe-cd/pipe](https://github.com/pipe-cd/pipe) are just containing the manifest templates, they cannot be used to install directly. The following command helps rendering those templates locally. The installable manifests will be stored at `.rendered-manifests` directory.

``` console
make render-manifests
```

#### Installing control-plane into the local cluster

Now, you can use the rendered manifests at `.rendered-manifests` to install control-plane to the local cluster.

Here is the command to install [quickstart](/docs/quickstart/)'s control-plane:

``` console
helm -n pipecd install pipecd .rendered-manifests/pipecd --dependency-update --create-namespace --values ./quickstart/control-plane-values.yaml
```

Once all components are running up, use `kubectl port-forward` to expose the installed control-plane on your localhost:

``` console
kubectl -n pipecd port-forward svc/pipecd 8080
```

Point your web browser to [http://localhost:8080](http://localhost:8080) to login with the configured static admin account: project = `quickstart`,
username = `hello-pipecd`, password = `hello-pipecd`.
