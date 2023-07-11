# Piped Agent
## Development

## Prerequisites

- [Go 1.19 or later](https://go.dev/)

## Repositories
- [pipecd](https://github.com/pipe-cd/pipecd): contains all source code and documentation of PipeCD project.

## Commands

- `make build/go`: builds all go modules including pipecd, piped, pipectl.
- `make test/go`: runs all unit tests of go modules.
- `make run/piped`: runs Piped locally (for more information, see [here](#how-to-run-piped-agent-locally)).
- `make gen/code`: generate Go and Typescript code from protos and mock configs. You need to run it if you modified any proto or mock definition files.

For the full list of available commands, please see the Makefile at the root of the repository.

## Setup Control Plane

1. Prepare Control Plane that piped connects. If you want to run a control plane locally, see [How to run Control Plane locally](https://github.com/pipe-cd/pipecd/tree/master/cmd/pipecd#how-to-run-control-plane-locally).

2. Access to Control Plane console, go to Piped list page and add a new piped. Then, copy generated Piped ID and key for `piped-config.yaml`

## How to run Piped agent locally

1. Prepare the piped configuration file `piped-config.yaml`. You can find an example config [here](https://github.com/pipe-cd/pipecd/tree/master/cmd/piped/piped-config.yaml).

2. Ensure that your `kube-context` is connecting to the right kubernetes cluster

3. Run the following command to start running `piped` (if you want to connect Piped to a locally running Control Plane, add `INSECURE=true` option)

    ``` console
    make run/piped CONFIG_FILE=cmd/piped/piped-config.yaml
    ```
