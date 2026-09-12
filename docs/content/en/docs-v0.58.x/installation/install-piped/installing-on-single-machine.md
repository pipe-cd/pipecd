---
title: "Installing on a single machine"
linkTitle: "Installing on a single machine"
weight: 5
description: >
  This page describes how to install a Piped on a single machine.
---

## Prerequisites

##### Having piped's ID and Key strings
- Ensure that the `piped` has been registered and you are having its PIPED_ID and PIPED_KEY strings.
- If you are not having them, this [page](../../../user-guide/managing-controlplane/registering-a-piped/) guides you how to register a new one.

##### Preparing SSH key
- If your Git repositories are private, `piped` requires a private SSH key to access those repositories.
- Please checkout [this documentation](https://help.github.com/en/github/authenticating-to-github/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent) for how to generate a new SSH key pair. Then add the public key to your repositories. (If you are using GitHub, you can add it to Deploy Keys at the repository's Settings page.)

## Installation

- Downloading the latest `piped` binary for your machine

  https://github.com/pipe-cd/pipecd/releases

- Preparing a piped configuration file as the following:

  ``` yaml
  apiVersion: pipecd.dev/v1beta1
  kind: Piped
  spec:
    projectID: {PROJECT_ID}
    pipedID: {PIPED_ID}
    pipedKeyFile: {PATH_TO_PIPED_KEY_FILE}
    # Write in a format like "host:443" because the communication is done via gRPC.
    apiAddress: {CONTROL_PLANE_API_ADDRESS}
    git:
      sshKeyFile: {PATH_TO_SSH_KEY_FILE}
    repositories:
      - repoId: {REPO_ID_OR_NAME}
        remote: git@github.com:{GIT_ORG}/{GIT_REPO}.git
        branch: {GIT_BRANCH}
    syncInterval: 1m
  ```

See [ConfigurationReference](../../../user-guide/managing-piped/configuration-reference/) for the full configuration.

- Start running the `piped`

  ``` console
  ./piped piped --config-file={PATH_TO_PIPED_CONFIG_FILE}
  ```

