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
- Please check out [this documentation](https://help.github.com/en/github/authenticating-to-github/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent) for how to generate a new SSH key pair. Then add the public key to your repositories. (If you are using GitHub, you can add it to Deploy Keys at the repository's Settings page.)

If your Git repositories are private, `piped` needs an SSH key to access them.

- Generate a new SSH key pair by following [GitHub’s guide to generating an SSH Key](https://help.github.com/en/github/authenticating-to-github/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent).  

>Note:
>If you are using GitHub, add the **public key** as a deploy key to your repositories.

### Install the `piped` V1 binary

Download the latest `piped` V1. See the [latest releases](https://github.com/pipe-cd/pipecd/releases) and find out the right binary for your machine.

## Installation

### Preparing the `Piped` configuration file

Plugins are external binaries that have to be referenced in the piped configuration file. There are no plugins set by default.

An example of a piped V1 configuration file using the [Example-stage plugin](https://github.com/pipe-cd/community-plugins/tree/main/plugins/example-stage):

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
  plugins: {}
```

See [ConfigurationReference](../../../user-guide/managing-piped/configuration-reference/) for the full configuration.

>**Note:**
>`Piped`'s plugins are versioned independently from PipeCD. Official plugins are maintained and monitored by the PipeCD Maintainers. See the [latest releases](https://github.com/pipe-cd/pipecd/releases) for more information.
>
>We now also have a repository for community-built plugins. See the [Community plugins repository on GitHub](https://github.com/pipe-cd/community-plugins) to know more.

## Run the `piped`

After you have configured your Piped configuration file, execute the `piped` binary and specify the path to the Piped configuration file.

  ``` console
  #Replace `<PATH_TO_PIPED_CONFIG_FILE>` with the path to your Piped configuration file.
  ./piped pipedv1 --config-file={PATH_TO_PIPED_CONFIG_FILE}
  ```

If you followed all steps correctly, you should have a running `piped` process on your system.
