# Piped Agent
## Development

## Prerequisites

- [Go 1.19](https://go.dev/)

## Repositories
- [pipecd](https://github.com/pipe-cd/pipecd): contains all source code and documentation of PipeCD project.

## Commands

- `make build/go`: builds all go modules including pipecd, piped, pipectl.
- `make test/go`: runs all unit tests of go modules.
- `make run/piped`: runs Piped locally (for more information, see [here](#how-to-run-piped-agent-locally)).
- `make gen/code`: generate Go and Typescript code from protos and mock configs. You need to run it if you modified any proto or mock definition files.

For the full list of available commands, please see the Makefile at the root of repository.

## How to run Piped agent locally

1. Prepare the piped configuration file `piped-config.yaml`

2. Ensure that your `kube-context` is connecting to the right kubernetes cluster

3. Run the following command to start running `piped` (if you want to connect Piped to a locally running Control Plane, add `INSECURE=true` option)

    ``` console
    make run/piped CONFIG_FILE=piped-config.yaml
    ```
