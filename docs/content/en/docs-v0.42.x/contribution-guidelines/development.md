---
title: "Development"
linkTitle: "Development"
weight: 2
description: >
  This page describes how to build, test PipeCD source code at your local environment.
---

## Prerequisites

- [Go 1.19](https://go.dev/)
- [Docker](https://www.docker.com/)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) (If you want to run Control Plane locally)
- [helm 3.8](https://helm.sh/docs/intro/install/) (If you want to run Control Plane locally)

## Repositories
- [pipecd](https://github.com/pipe-cd/pipecd): contains all source code and documentation of PipeCD project.
- [examples](https://github.com/pipe-cd/examples): contains various generated examples to demonstrate how to use PipeCD.

## Commands

- `make build/go`: builds all go modules including pipecd, piped, pipectl.
- `make build/web`: builds the static files for web.

- `make test/go`: runs all unit tests of go modules.
- `make test/web`: runs all unit tests of web.
- `make test/integration`: runs integration tests.

- `make run/piped`: runs Piped locally (for more information, see [here](#how-to-run-piped-agent-locally)).
- `make run/site`: runs PipeCD site locally (requires [hugo](https://github.com/gohugoio/hugo) with `_extended` version `0.92.1` or later to be installed).

- `make gen/code`: generate Go and Typescript code from protos and mock configs. You need to run it if you modified any proto or mock definition files.

For the full list of available commands, please see the Makefile at the root of repository.

## How to run Control Plane locally

1. Start running a Kubernetes cluster

    ``` console
    make kind-up
    ```

    Once it is no longer used, run `make kind-down` to delete it.

2. Install Control Plane into the local cluster

    ``` console
    make run/pipecd
    ```

    Once all components are running up, use `kubectl port-forward` to expose the installed Control Plane on your localhost:

    ``` console
    kubectl -n pipecd port-forward svc/pipecd 8080
    ```

3. Access to the local Control Plane web console

    Point your web browser to [http://localhost:8080](http://localhost:8080) to login with the configured static admin account: project = `quickstart`, username = `hello-pipecd`, password = `hello-pipecd`.

## How to run Piped agent locally

1. Prepare the piped configuration file `piped-config.yaml`

2. Ensure that your `kube-context` is connecting to the right kubernetes cluster

3. Run the following command to start running `piped` (if you want to connect Piped to a locally running Control Plane, add `INSECURE=true` option)

    ``` console
    make run/piped CONFIG_FILE=piped-config.yaml
    ```

## Docs and workaround with docs

PipeCD official site contains multiple versions of documentation, all placed under the `/docs/content/en` directory, which are:
- `/docs`: stable version docs, usually synced with the latest released version docs.
- `/docs-dev`: experimental version docs, contains docs for not yet released features or changes.
- `/docs-v0.x.x`: contains docs for specified version family (a version family is all versions which in the same major release).

Basically, we have two simple rules:
- Do not touch to the `/docs` content directly.
- Keep stable docs version synced with the latest released docs version.

Here are the flow of docs contribution regard some known scenarios:
1. Update docs that are related to a specified version (which is not the latest released version):
In such case, update the docs under `/docs-v0.x.x` is enough.
2. Update docs for not yet released features or changes:
In such case, update the docs under `/docs-dev` is enough.
3. Update docs that are related to the latest released docs version:
- Change the docs' content that fixes the issue under `/docs-dev` and `/docs-v0.x.x`, they share the same file structure so you should find the right files in both directories.
- Use `make gen/stable-docs` command to sync the latest released version docs under `/docs-v0.x.x` to `/docs`

If you find any issues related to the docs, we're happy to accept your help.
