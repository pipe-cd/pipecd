---
title: "Adding a git repository"
linkTitle: "Adding git repository"
weight: 2
description: >
  Learn how to add a git repository in the `piped` configuration file.

---
In the `piped` configuration file, all the git repositories that you want to be tracked by the `piped` agent are specified.

A Git repository contains one or more deployable applications where each application is put inside a directory called as [application directory](../../../concepts/#application-directory).

An application directory contains an **application configuration file** as well as application manifests.
The `piped` periodically checks the new commits and fetches the needed manifests from those repositories for executing the deployment.

A single `piped` can be configured to handle one or more Git repositories.
In order to enable a new Git repository, add a new [GitRepository](../configuration-reference/#gitrepository) block to the `repositories` field in the `piped` configuration file.

For example, in the following snippet, `piped` will take the `master` branch of [pipe-cd/examples](https://github.com/pipe-cd/examples) repository as a target Git repository for deployments.

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

In most cases, you would want to deal with private git repositories. For accessing private repositories, `piped` needs a private SSH key, which can be configured while [installing](../../../installation/install-piped/installing-on-kubernetes/) with `secret.sshKey` in the Helm chart.

``` console
helm install dev-piped pipecd/piped --version={VERSION} \
  --set-file config.data={PATH_TO_PIPED_CONFIG_FILE} \
  --set-file secret.data.piped-key={PATH_TO_PIPED_KEY_FILE} \
  --set-file secret.data.ssh-key={PATH_TO_PRIVATE_SSH_KEY_FILE}
```

You can see the [git configuration reference](../configuration-reference/#git) to learn more about the configurable fields.

Currently, `piped` allows configuring only one private SSH key for all specified git repositories. For working with multiple repositories, you either have to configure the same SSH key for all of those private repositories, or use separate `piped`s for each repository. We are working on adding support for multiple SSH keys for a single `piped`.
