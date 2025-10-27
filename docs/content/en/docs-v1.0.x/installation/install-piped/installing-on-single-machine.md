---
title: "Installing on a single machine"
linkTitle: "Installing on a single machine"
weight: 5
description: >
  This page describes how to install a Piped on a single machine.
---


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

### Install the `piped` V1 binary

Download the latest `piped` V1. See the [latest releases](https://github.com/pipe-cd/pipecd/releases) and find out the right binary for your machine.

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

See [ConfigurationReference](../../../user-guide/managing-piped/configuration-reference/) for the full configuration.

>**Note:**
>`Piped`'s plugins are versioned independently from PipeCD. Official plugins are maintained and monitored by the PipeCD Maintainers. See the [latest releases](https://github.com/pipe-cd/pipecd/releases) for more information.
>
>With PipeCD V1, we have also added support for community plugins. See the [Community plugins repository on GitHub](https://github.com/pipe-cd/community-plugins)

## Run the `piped`

After you have configured your Piped configuration file, execute the `piped` binary and specify the path to the Piped configuration file.

  ``` console
  #Replace `<PATH_TO_PIPED_CONFIG_FILE>` with the path to your Piped configuration file.
  ./piped pipedv1 --config-file={PATH_TO_PIPED_CONFIG_FILE}
  ```

If you followed all steps correctly, you should have a running `piped` process on your system.
