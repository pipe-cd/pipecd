---
title: "Installing on Kubernetes cluster"
linkTitle: "Installing on Kubernetes cluster"
weight: 1
description: >
  This page describes how to install Piped on Kubernetes cluster.
---

## Prerequisites

##### Having piped's ID and Key strings
- Ensure that the `piped` has been registered and you are having its PIPED_ID and PIPED_KEY strings.
- If you are not having them, this [page](/docs/operator-manual/control-plane/registering-a-piped/) guides you how to register a new one.

##### Preparing SSH key
- If your Git repositories are private, `piped` requires a private SSH key to access those repositories.
- Please checkout [this documentation](https://help.github.com/en/github/authenticating-to-github/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent) for how to generate a new SSH key pair. Then add the public key to your repositories. (If you are using GitHub, you can add it to Deploy Keys at the repository's Settings page.)

## In the cluster-wide mode
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
    projectID: {PROJECT_ID}
    pipedID: {PIPED_ID}
    pipedKeyFile: /etc/piped-secret/piped-key
    webAddress: {CONTROL_PLANE_WEB_ADDRESS}
    # Write in a format like "host:443" because the communication is done via gRPC.
    apiAddress: {CONTROL_PLANE_API_ADDRESS}
    git:
      sshKeyFile: /etc/piped-secret/ssh-key
    repositories:
      - repoId: {REPO_ID_OR_NAME}
        remote: git@github.com:{GIT_ORG}/{GIT_REPO}.git
        branch: {GIT_BRANCH}
    syncInterval: 1m
  ```

- Installing by using `Helm 3`

  ``` console
  helm repo update

  helm upgrade -i dev-piped pipecd/piped --version={{< blocks/latest_version >}} --namespace={NAMESPACE} \
    --set-file config.data={PATH_TO_PIPED_CONFIG_FILE} \
    --set-file secret.pipedKey.data={PATH_TO_PIPED_KEY_FILE} \
    --set-file secret.sshKey.data={PATH_TO_PRIVATE_SSH_KEY_FILE}
  ```

Note: Be sure to set `--set args.insecure=true` if your control-plane has not TLS-enabled yet.

See [values.yaml](https://github.com/pipe-cd/manifests/blob/master/manifests/piped/values.yaml) for the full values.

## In the namespaced mode
The previous way requires installing cluster-level resources. If you want to restrict Piped's permission within the namespace where Piped runs on, this way is for you.
Most parts are identical to the previous way, but some are slightly different.

- Adding a new cloud provider like below to the previous piped configuration file:

  ``` yaml
  apiVersion: pipecd.dev/v1beta1
  kind: Piped
  spec:
    projectID: {PROJECT_ID}
    pipedID: {PIPED_ID}
    pipedKeyFile: /etc/piped-secret/piped-key
    webAddress: {CONTROL_PLANE_WEB_ADDRESS}
    # Write in a format like "host:443" because the communication is done via gRPC.
    apiAddress: {CONTROL_PLANE_API_ADDRESS}
    git:
      sshKeyFile: /etc/piped-secret/ssh-key
    repositories:
      - repoId: REPO_ID_OR_NAME
        remote: git@github.com:{GIT_ORG}/{GIT_REPO}.git
        branch: {GIT_BRANCH}
    syncInterval: 1m
    # This is needed to restrict to limit the access range to within a namespace.
    cloudProviders:
      - name: my-kubernetes
        type: KUBERNETES
        config:
          appStateInformer:
            namespace: {NAMESPACE}
  ```

- Then installing it with the following options:

  ``` console
  helm repo update

  helm upgrade -i dev-piped pipecd/piped --version={{< blocks/latest_version >}} --namespace={NAMESPACE} \
    --set-file config.data={PATH_TO_PIPED_CONFIG_FILE} \
    --set-file secret.pipedKey.data={PATH_TO_PIPED_KEY_FILE} \
    --set-file secret.sshKey.data={PATH_TO_PRIVATE_SSH_KEY_FILE} \
    --set args.enableDefaultKubernetesCloudProvider=false \
    --set rbac.scope=namespace
  ```

#### In case on OpenShift less than 4.2
OpenShift uses an arbitrarily assigned user ID when it starts a container.
Starting from OpenShift 4.2, it also inserts that user into /etc/passwd for using by the application inside the container,
but before that version, the assigned user is missing in that file. That blocks workloads of gcr.io/pipecd/piped image.
Therefore if you are running on OpenShift with a version before 4.2, please use gcr.io/pipecd/piped-okd image with the following command:

``` console
  helm upgrade -i dev-piped pipecd/piped --version={{< blocks/latest_version >}} --namespace={NAMESPACE} \
    --set-file config.data={PATH_TO_PIPED_CONFIG_FILE} \
    --set-file secret.pipedKey.data={PATH_TO_PIPED_KEY_FILE} \
    --set-file secret.sshKey.data={PATH_TO_PRIVATE_SSH_KEY_FILE} \
    --set args.enableDefaultKubernetesCloudProvider=false \
    --set rbac.scope=namespace
    --set args.addLoginUserToPasswd=true \
    --set securityContext.runAsNonRoot=true \
    --set securityContext.runAsUser={UID} \
    --set securityContext.fsGroup={FS_GROUP} \
    --set securityContext.runAsGroup=0 \
    --set image.repository="gcr.io/pipecd/piped-okd"
```
