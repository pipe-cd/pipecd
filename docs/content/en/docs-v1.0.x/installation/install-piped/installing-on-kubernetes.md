---
title: "Running the Kubernetes plugin"
linkTitle: "Running the Kubernetes plugin"
weight: 1
description: >
  This page describes how to run the official Kubernetes plugin on your `piped`.
---

The Kubernetes plugin for PipeCD V1 adds support for deploying and synchronizing applications on Kubernetes containers.

## Prerequisites

### A registered `piped`

- Make sure your `piped` is registered in the Control Plane and that you have its **PIPED_ID** and **PIPED_KEY**.  
- If not, follow the guide to [register a new `Piped`](../../../user-guide/managing-controlplane/registering-a-piped/).

### SSH Key for Git Repositories

- If your Git repositories are private, `piped` requires a private SSH key to access those repositories.
- Please checkout [this documentation](https://help.github.com/en/github/authenticating-to-github/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent) for how to generate a new SSH key pair. Then add the public key to your repositories. (If you are using GitHub, you can add it to Deploy Keys at the repository's Settings page.)


 If your Git repositories are private, `piped` needs an SSH key to access them.

- Generate a new SSH key pair by following [GitHub’s guide to generating an SSH Key](https://help.github.com/en/github/authenticating-to-github/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent).  

>Note:
>If you are using GitHub, add the **public key** as a deploy key to your repositories.

## Installation

### Preparing the `Piped` configuaration file

In PipeCD V0, the default platform provider was Kubernetes. In PipeCD V1, since platforms have been replaced by plugins, there are not defaults set. Plugins are external binaries that have to be referenced in the piped configuration file.

An example of an old piped configuration file:

```yaml
 apiVersion: pipecd.dev/v1beta1
  kind: Piped
  spec:
    projectID: {PROJECT_ID}
    pipedID: {PIPED_ID}
    pipedKeyFile: /etc/piped-secret/piped-key
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

An example of the a piped V1 configuration file using the official Kubernetes plugin:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: dev
  pipedID: {PIPED_ID}
  pipedKeyData: {PIPED_KEY}
  apiAddress: {CONTROL_PLANE_API_ADDRESS} 
  repositories:
    - repoId: {REPO_ID_OR_NAME}
      remote: git@github.com:{GIT_ORG}/{GIT_REPO}.git
      branch: {GIT_BRANCH}
  syncInterval: 1m
  plugins:
    - name: kubernetes
      port: 7001
      url: https://github.com/pipe-cd/pipecd/releases/download/pkg%2Fapp%2Fpipedv1%2Fplugin%2Fkubernetes%2Fv0.1.0/kubernetes_v0.1.0_darwin_arm64
      deployTargets:
        - name: local
          config:
            kubectlVersion: 1.32.4
            kubeConfigPath: /Users/sawanteeshaan/.kube/config
    - name: wait
      port: 7002
      url: https://github.com/pipe-cd/pipecd/releases/download/pkg%2Fapp%2Fpipedv1%2Fplugin%2Fwait%2Fv0.1.0/wait_v0.1.0_darwin_arm64
      # url: file:///Users/s12228/.piped/plugins/wait
    - name: example-stage
      port: 7003
      url: https://github.com/pipe-cd/community-plugins/releases/download/plugins%2Fexample-stage%2Fv0.1.0/example-stage_v0.1.0_darwin_arm64
      config:
        commonMessage: "Hello Middle Earth! This is the wave from UFO!"
```

>**Note:**
>`Piped`'s plugins are versioned independently from PipeCD. Official plugins are maintained and monitored by the PipeCD Maintainers. See the [latest releases](https://github.com/pipe-cd/pipecd/releases) for more information.
>
>With PipeCD V1, we have also added support for community plugins. See the [Community plugins repository on GitHub](https://github.com/pipe-cd/community-plugins)

## In the cluster-wide mode

This way requires installing cluster-level resources. Piped installed with this way can perform deployment workloads against any other namespaces than the where Piped runs on.

- Preparing a piped configuration file as the following

  ``` yaml
  apiVersion: pipecd.dev/v1beta1
  kind: Piped
  spec:
    projectID: {PROJECT_ID}
    pipedID: {PIPED_ID}
    pipedKeyFile: /etc/piped-secret/piped-key
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

See [ConfigurationReference](../../../user-guide/managing-piped/configuration-reference/) for the full configuration.

- Installing by using [Helm](https://helm.sh/docs/intro/install/) (3.8.0 or later)

  {{< tabpane >}}
  {{< tab lang="bash" header="Piped" >}}
# This command just installs a Piped with the specified version.
# Whenever you want to upgrade that Piped to a new version or update its config data
# you have to restart it by re-running this command.

helm upgrade -i dev-piped oci://ghcr.io/pipe-cd/chart/piped --version={{< blocks/latest_version >}} --namespace={NAMESPACE} \
  --set-file config.data={PATH_TO_PIPED_CONFIG_FILE} \
  --set-file secret.data.piped-key={PATH_TO_PIPED_KEY_FILE} \
  --set-file secret.data.ssh-key={PATH_TO_PRIVATE_SSH_KEY_FILE}
  {{< /tab >}}
  
  {{< /tabpane >}}

  Note: Be sure to set `--set args.insecure=true` if your Control Plane has not TLS-enabled yet.

  See [values.yaml](https://github.com/pipe-cd/pipecd/blob/master/manifests/piped/values.yaml) for the full values.

## In the namespaced mode
The previous way requires installing cluster-level resources. If you want to restrict Piped's permission within the namespace where Piped runs on, this way is for you.
Most parts are identical to the previous way, but some are slightly different.

- Adding a new cloud provider like below to the previous piped configuration file

  ``` yaml
  apiVersion: pipecd.dev/v1beta1
  kind: Piped
  spec:
    projectID: {PROJECT_ID}
    pipedID: {PIPED_ID}
    pipedKeyFile: /etc/piped-secret/piped-key
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
    platformProviders:
      - name: my-kubernetes
        type: KUBERNETES
        config:
          appStateInformer:
            namespace: {NAMESPACE}
  ```

- Installing by using [Helm](https://helm.sh/docs/intro/install/) (3.8.0 or later)

  {{< tabpane >}}
  {{< tab lang="bash" header="Piped" >}}
# This command just installs a Piped with the specified version.
# Whenever you want to upgrade that Piped to a new version or update its config data
# you have to restart it by re-running this command.

helm upgrade -i dev-piped oci://ghcr.io/pipe-cd/chart/piped --version={{< blocks/latest_version >}} --namespace={NAMESPACE} \
  --set-file config.data={PATH_TO_PIPED_CONFIG_FILE} \
  --set-file secret.data.piped-key={PATH_TO_PIPED_KEY_FILE} \
  --set-file secret.data.ssh-key={PATH_TO_PRIVATE_SSH_KEY_FILE} \
  --set args.enableDefaultKubernetesCloudProvider=false \
  --set rbac.scope=namespace
  {{< /tab >}}
  {{< /tabpane >}}

## In case on OpenShift less than 4.2

OpenShift uses an arbitrarily assigned user ID when it starts a container.
Starting from OpenShift 4.2, it also inserts that user into `/etc/passwd` for using by the application inside the container,
but before that version, the assigned user is missing in that file. That blocks workloads of `ghcr.io/pipe-cd/piped` image.
Therefore if you are running on OpenShift with a version before 4.2, please use `ghcr.io/pipe-cd/piped-okd` image with the following command:

- Installing by using [Helm](https://helm.sh/docs/intro/install/) (3.8.0 or later)

  {{< tabpane >}}
  {{< tab lang="bash" header="Piped" >}}

# This command just installs a Piped with the specified version.

# Whenever you want to upgrade that Piped to a new version or update its config data

# you have to restart it by re-running this command.


helm upgrade -i dev-piped oci://ghcr.io/pipe-cd/chart/piped --version={{< blocks/latest_version >}} --namespace={NAMESPACE} \
  --set-file config.data={PATH_TO_PIPED_CONFIG_FILE} \
  --set-file secret.data.piped-key={PATH_TO_PIPED_KEY_FILE} \
  --set-file secret.data.ssh-key={PATH_TO_PRIVATE_SSH_KEY_FILE} \
  --set args.enableDefaultKubernetesCloudProvider=false \
  --set rbac.scope=namespace
  --set args.addLoginUserToPasswd=true \
  --set securityContext.runAsNonRoot=true \
  --set securityContext.runAsUser={UID} \
  --set securityContext.fsGroup={FS_GROUP} \
  --set securityContext.runAsGroup=0 \
  --set image.repository="ghcr.io/pipe-cd/piped-okd"
  {{< /tab >}}
  {{< /tabpane >}}
