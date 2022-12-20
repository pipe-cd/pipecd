---
title: "Adding a git repository"
linkTitle: "Adding git repository"
weight: 2
description: >
  This page describes how to add a new Git repository.
---

In the `piped` configuration file, we specify the list of Git repositories should be handled by the `piped`.
A Git repository contains one or more deployable applications where each application is put inside a directory called as [application directory](../../../concepts/#application-directory).
That directory contains an application configuration file as well as application manifests.
The `piped` periodically checks the new commits and fetches the needed manifests from those repositories for executing the deployment.

A single `piped` can be configured to handle one or more Git repositories.
In order to enable a new Git repository, let's add a new [GitRepository](../configuration-reference/#gitrepository) block to the `repositories` field in the `piped` configuration file.

For example, with the following snippet, `piped` will take the `master` branch of [pipe-cd/examples](https://github.com/pipe-cd/examples) repository as a target Git repository for doing deployments.

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  repositories:
    - repoId: examples
      remote: git@github.com:pipe-cd/examples.git
      branch: master
```

In most of the cases, we want to deal with private Git repositories. For accessing those private repositories, `piped` needs a private SSH key, which can be configured while [installing](../../../installation/install-piped/installing-on-kubernetes/) with `secret.sshKey` in the Helm chart.

``` console
helm install dev-piped pipecd/piped --version={VERSION} \
  --set-file config.data={PATH_TO_PIPED_CONFIG_FILE} \
  --set-file secret.data.piped-key={PATH_TO_PIPED_KEY_FILE} \
  --set-file secret.data.ssh-key={PATH_TO_PRIVATE_SSH_KEY_FILE}
```

You can see this [configuration reference](../configuration-reference/#git) for more configurable fields about Git commands.

Currently, `piped` allows configuring only one private SSH key for all specified Git repositories. So you can configure the same SSH key for all of those private repositories, or break them into separate `piped`s. In the near future, we also want to update `piped` to support loading multiple SSH keys.
