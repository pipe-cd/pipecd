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

## Installation on Kubernetes cluster
### In the cluster-wide mode
This way requires installing cluster-level resources. Piped installed with this way can perform deployment workloads against any other namespaces than the where Piped runs on.

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

### In the namespaced mode
The previous way requires installing cluster-level resources. If you want to restrict Piped's permission within the namespace where Piped runs on, this way is for you.
Most parts are identical to the previous way, but some are slightly different.

- Adding a new cloud provider like below to the previous piped configuration file:

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
    # This is needed to restrict to limit the access range to within a namespace.
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

#### In case on OpenShift less than 4.2
Containers on OpenShift run using arbitrary Linux users. What we don't know the UID in advance causes a couple of problems.
To prevent you from such problems, you should set additional options to the previous `helm upgrade` command.


``` console
  --set args.addLoginUserToPasswd=true \
  --set securityContext.runAsNonRoot=true \
  --set securityContext.runAsUser=YOUR_UID \
  --set securityContext.fsGroup=YOUR_FS_GROUP \
  --set securityContext.runAsGroup=0 \
  --set image.repository="gcr.io/pipecd/piped-okd"
```

Keep in mind OpenShift 4.2 onward doesn't have kind of like this concern.

## Installing on single machine

- Downloading the latest `piped` binary for your machine

  https://github.com/pipe-cd/pipe/releases

- Preparing a piped configuration file as the following:

  ``` yaml
  apiVersion: pipecd.dev/v1beta1
  kind: Piped
  spec:
    projectID: {YOUR_PROJECT_ID}
    pipedID: {YOUR_PIPED_ID}
    pipedKeyFile: {PATH_TO_PIPED_KEY_FILE}
    # Write in a format like "host:443" because the communication is done via gRPC.
    apiAddress: {YOUR_CONTROL_PLANE_ADDRESS}
    webAddress: http://{YOUR_CONTROL_PLANE_ADDRESS}
    git:
      sshKeyFile: {PATH_TO_SSH_KEY_FILE}
    repositories:
      - repoId: {REPO_ID_OR_NAME}
        remote: git@github.com:{YOUR_GIT_ORG}/{YOUR_GIT_REPO}.git
        branch: {YOUR_GIT_BRANCH}
    syncInterval: 1m
  ```

- Start running the `piped`

  ``` console
  ./piped piped --config-file=PATH_TO_PIPED_CONFIG_FILE
  ```

