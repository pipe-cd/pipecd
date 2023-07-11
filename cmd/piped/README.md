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

1. Prepare the piped configuration file. This is an example configuration;
    ```yaml:piped-config.yaml
    apiVersion: pipecd.dev/v1beta1
    kind: Piped
    spec:
    projectID: quickstart
    # FIXME: Replace here with your piped ID.
    pipedID: 7accd470-1786-49ee-ac09-3c4d4e31dc12
    # Base64 encoded string of the piped private key. You can generate it by the following command.
    # echo "your-piped-key" | base64
    # FIXME: Replace here with your piped key file path.
    pipedKeyData: OTl4c2RqdjUxNTF2OW1sOGw5ampndXUyZjB2aGJ4dGw0bHVkamF4Mmc3a3l1enFqY20K
    # Write in a format like "host:443" because the communication is done via gRPC.
    # FIXME: Replace here with your piped address if you connect Piped to a control plane that does not run locally.
    apiAddress: localhost:8080
    # git:
    #   sshKeyFile: /etc/piped-secret/ssh-key
    repositories:
        - repoId: example
        remote: git@github.com:pipe-cd/examples.git
        branch: master
    syncInterval: 1m
    platformProviders:
        - name: example-kubernetes
        type: KUBERNETES
        config:
            # FIXME: Replace here with your kubeconfig absolute file path.
            kubeConfigPath: /path/to/.kube/config
    ```

2. Ensure that your `kube-context` is connecting to the right kubernetes cluster

3. Run the following command to start running `piped` (if you want to connect Piped to a locally running Control Plane, add `INSECURE=true` option)

    ``` console
    make run/piped CONFIG_FILE=piped-config.yaml
    ```
