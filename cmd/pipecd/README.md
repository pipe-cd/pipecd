# PipeCD Controle Plane
## Development

## Prerequisites

- [Go 1.24 or later](https://go.dev/)
- [NodeJS v20 or later](https://nodejs.org/en/)
- [Docker](https://www.docker.com/)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) (If you want to run Control Plane locally)
- [helm 3.8](https://helm.sh/docs/intro/install/) (If you want to run Control Plane locally)

## Repositories
- [pipecd](https://github.com/pipe-cd/pipecd): contains all source code and documentation of PipeCD project.

## Commands

- `make build/go`: builds all go modules including pipecd, piped, pipectl.
- `make test/go`: runs all unit tests of go modules.
- `make test/integration`: runs integration tests.
- `make gen/code`: generate Go and Typescript code from protos and mock configs. You need to run it if you modified any proto or mock definition files.

For the full list of available commands, please see the Makefile at the root of repository.

## How to run Control Plane locally

1. Start running a Kubernetes cluster (If you don't have any Kubernetes cluster to use)

    ``` console
    make kind-up
    ```

    Once it is no longer used, run `make kind-down` to delete it.

2. Install the web dependencies module

    ``` console
    make update/web-deps
    ```

3. Install Control Plane into the local cluster

    ``` console
    make run/pipecd
    ```

    Once all components are running up, use `kubectl port-forward` to expose the installed Control Plane on your localhost:

    ``` console
    kubectl -n pipecd port-forward svc/pipecd 8080
    ```

4. Access to the local Control Plane web console

    Point your web browser to [http://localhost:8080](http://localhost:8080) to login with the configured static admin account: project = `quickstart`, username = `hello-pipecd`, password = `hello-pipecd`.

## How to run Piped locally and add an application to your cluster
See [How to run Piped agent locally](https://github.com/pipe-cd/pipecd/tree/master/cmd/piped#how-to-run-piped-agent-locally).

## How to set up OIDC provider(Keycloak) locally

Run `make setup-local-oidc` to set up local OIDC provider(keycloak). This will create a new Keycloak realm and an OIDC client for PipeCD. See [Local Keycloak](../../hack/oidc/README.md) for more details.
