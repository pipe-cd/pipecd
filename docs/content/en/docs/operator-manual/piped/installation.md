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

##### Preparing SSH key
- If your Git repositories are private, `piped` requires a private SSH key to access those repositories.
- Please checkout [this documentation](https://help.github.com/en/github/authenticating-to-github/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent) for how to generate a new SSH key. Then add the public key to your repositories. (If you are using GitHub, you can add it to Deploy Keys at the repository's Settings page.)

## Installation

### Installing on Kubernetes cluster

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
    # Write in a format like "host:443" because the communication is done via gRPC.
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
    --set-file config.data=PATH_TO_PIPED_CONFIG_FILE \
    --set-file secret.pipedKey.data=PATH_TO_PIPED_KEY_FILE \
    --set-file secret.sshKey.data=PATH_TO_PRIVATE_SSH_KEY_FILE
  ```

Note: Be sure to set `--set args.insecure=true` if your control-plane has not TLS-enabled yet.

See [values.yaml](https://github.com/pipe-cd/manifests/blob/master/manifests/piped/values.yaml) for the full values.

### Installing on Kubernetes cluster in the namespaced mode
The previous way requires installing cluster-level resources. If you want to restrict Piped's permission within the namespace as the same as Piped, this way is for you.
Most part is identical to the previous way, but some part is slightly different.

- Adding a new cloud provider like below to the previous piped configuration file:

  ``` yaml
  apiVersion: pipecd.dev/v1beta1
  kind: Piped
  spec:
    cloudProviders:
      - name: my-kubernetes
        type: KUBERNETES
        config:
          appStateInformer:
            namespace: {YOUR_NAMESPACE}
  ```

- Then installing it with the following options:

  ``` console
  helm repo update

  helm upgrade -i dev-piped pipecd/piped --version=VERSION --namespace=NAMESPACE \
    --set-file config.data=PATH_TO_PIPED_CONFIG_FILE \
    --set-file secret.pipedKey.data=PATH_TO_PIPED_KEY_FILE \
    --set-file secret.sshKey.data=PATH_TO_PRIVATE_SSH_KEY_FILE \
    --set args.enableDefaultKubernetesCloudProvider=false \
    --set rbac.scope=namespace
  ```

### Installing on single machine

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
    # Write in a format like "host:443" because the communication is done via gRPC.
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
  --tools-dir=PATH_TO_TOOLS_DIR
  ```

Note that the `TOOLS_DIR` must be the directory to which `piped` has write permission because `piped` installs the tools to deploy applications underneath it.
