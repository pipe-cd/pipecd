---
title: "Installation"
linkTitle: "Installation"
weight: 1
description: >
  This page describes how to install a piped.
---

## Prerequisites

##### Having piped's ID and Key strings
- Ensure that the `piped` has been registered and you are having its PIPED_ID and PIPED_KEY strings.
- If you are not having them, this [page](/docs/operator-manual/control-plane/registering-a-piped/) guides you how to register a new one.

##### Preparing a SSH key
- If your Git repositories are private, `piped` requires a private SSH key to access those repositories.
- Please checkout [this documentation](https://help.github.com/en/github/authenticating-to-github/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent) for how to generate a new SSH key. Then add the public key to your repositories. (If you are using GitHub, you can add it to Deploy Keys at the repository's Settings page.)

## Installation

### Installing on a Kubernetes cluster

- Adding `pipecd` helm chart repository

  ```
  helm repo add pipecd https://charts.pipecd.dev
  ```

- Preparing a piped configuration file as the following:

  ``` yaml
  apiVersion: pipecd.dev/v1beta1
  kind: Piped
  spec:
    projectID: YOUR_PROJECT_ID
    pipedID: YOUR_PIPED_ID
    pipedKeyFile: /etc/piped-secret/piped-key
    # Write in a format like "host:443" because the communication is done via gRCP.
    apiAddress: YOUR_CONTROL_PLANE_ADDRESS
    webAddress: http://YOUR_CONTROL_PLANE_ADDRESS
    git:
      sshKeyFile: /etc/piped-secret/ssh-key
    repositories:
      - repoId: REPO_ID_OR_NAME
        remote: git@github.com:YOUR_GIT_ORG/YOUR_GIT_REPO.git
        branch: master
    syncInterval: 1m
  ```

- Installing by using `Helm 3`

  ``` console
  helm repo update

  helm upgrade -i dev-piped pipecd/piped --version=VERSION --namespace=NAMESPACE \
    --set args.insecure=true \
    --set-file config.data=PATH_TO_PIPED_CONFIG_FILE \
    --set-file secret.pipedKey.data=PATH_TO_PIPED_KEY_FILE \
    --set-file secret.sshKey.data=PATH_TO_PRIVATE_SSH_KEY_FILE
  ```
Note: Be sure to set `args.insecure=false` if your control-plane is already TLS-enabled.

### Installing on a single machine

- Downloading the latest `piped` binary for your machine

  https://github.com/pipe-cd/pipe/releases

- Preparing a piped configuration file as the following:

  ``` yaml
  apiVersion: pipecd.dev/v1beta1
  kind: Piped
  spec:
    projectID: YOUR_PROJECT_ID
    pipedID: YOUR_PIPED_ID
    pipedKeyFile: PATH_TO_PIPED_KEY_FILE
    # Write in a format like "host:443" because the communication is done via gRCP.
    apiAddress: YOUR_CONTROL_PLANE_ADDRESS
    webAddress: http://YOUR_CONTROL_PLANE_ADDRESS
    git:
      sshKeyFile: PATH_TO_SSH_KEY_FILE
      sshConfigFilePath: PATH_TO_SSH_CONFIG_FILE
    repositories:
      - repoId: REPO_ID_OR_NAME
        remote: git@github.com:YOUR_GIT_ORG/YOUR_GIT_REPO.git
        branch: master
    syncInterval: 1m
  ```

- Start running the `piped`

  ``` console
  ./piped piped \
  --config-file=PATH_TO_PIPED_CONFIG_FILE \
  --bin-dir=PATH_TO_BIN_DIR
  ```

Note that the `BIN_DIR` must be the directory to which `piped` has write permission because `piped` installs the tools to deploy applications underneath it.
