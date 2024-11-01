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
- If you are not having them, this [page](../../../user-guide/managing-controlplane/registering-a-piped/) guides you how to register a new one.

##### Preparing SSH key
- If your Git repositories are private, `piped` requires a private SSH key to access those repositories.
- Please checkout [this documentation](https://help.github.com/en/github/authenticating-to-github/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent) for how to generate a new SSH key pair. Then add the public key to your repositories. (If you are using GitHub, you can add it to Deploy Keys at the repository's Settings page.)

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
  {{< tab lang="bash" header="Piped with Remote-upgrade" >}}
# Enable remote-upgrade feature of Piped.
# https://pipecd.dev/docs/user-guide/managing-piped/remote-upgrade-remote-config/#remote-upgrade
# This allows upgrading Piped to a new version from the web console.
# But we still need to restart Piped when we want to update its config data.

helm upgrade -i dev-piped oci://ghcr.io/pipe-cd/chart/piped --version={{< blocks/latest_version >}} --namespace={NAMESPACE} \
  --set launcher.enabled=true \
  --set-file config.data={PATH_TO_PIPED_CONFIG_FILE} \
  --set-file secret.data.piped-key={PATH_TO_PIPED_KEY_FILE} \
  --set-file secret.data.ssh-key={PATH_TO_PRIVATE_SSH_KEY_FILE}
  {{< /tab >}}
  {{< tab lang="bash" header="Piped with Remote-upgrade and Remote-config" >}}
# Enable both remote-upgrade and remote-config features of Piped.
# https://pipecd.dev/docs/user-guide/managing-piped/remote-upgrade-remote-config/#remote-config
# Beside of the ability to upgrade Piped to a new version from the web console,
# remote-config allows loading the Piped config stored in a remote location such as a Git repository.
# Whenever the config data is changed, it loads the new config and restarts Piped to use that new config.

helm upgrade -i dev-piped oci://ghcr.io/pipe-cd/chart/piped --version={{< blocks/latest_version >}} --namespace={NAMESPACE} \
  --set launcher.enabled=true \
  --set launcher.configFromGitRepo.enabled=true \
  --set launcher.configFromGitRepo.repoUrl=git@github.com:{GIT_ORG}/{GIT_REPO}.git \
  --set launcher.configFromGitRepo.branch={GIT_BRANCH} \
  --set launcher.configFromGitRepo.configFile={RELATIVE_PATH_TO_PIPED_CONFIG_FILE_IN_GIT_REPO} \
  --set launcher.configFromGitRepo.sshKeyFile=/etc/piped-secret/ssh-key \
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
  {{< tab lang="bash" header="Piped with Remote-upgrade" >}}
# Enable remote-upgrade feature of Piped.
# https://pipecd.dev/docs/user-guide/managing-piped/remote-upgrade-remote-config/#remote-upgrade
# This allows upgrading Piped to a new version from the web console.
# But we still need to restart Piped when we want to update its config data.

helm upgrade -i dev-piped oci://ghcr.io/pipe-cd/chart/piped --version={{< blocks/latest_version >}} --namespace={NAMESPACE} \
  --set launcher.enabled=true \
  --set-file config.data={PATH_TO_PIPED_CONFIG_FILE} \
  --set-file secret.data.piped-key={PATH_TO_PIPED_KEY_FILE} \
  --set-file secret.data.ssh-key={PATH_TO_PRIVATE_SSH_KEY_FILE} \
  --set args.enableDefaultKubernetesCloudProvider=false \
  --set rbac.scope=namespace
  {{< /tab >}}
  {{< tab lang="bash" header="Piped with Remote-upgrade and Remote-config" >}}
# Enable both remote-upgrade and remote-config features of Piped.
# https://pipecd.dev/docs/user-guide/managing-piped/remote-upgrade-remote-config/#remote-config
# Beside of the ability to upgrade Piped to a new version from the web console,
# remote-config allows loading the Piped config stored in a remote location such as a Git repository.
# Whenever the config data is changed, it loads the new config and restarts Piped to use that new config.

helm upgrade -i dev-piped oci://ghcr.io/pipe-cd/chart/piped --version={{< blocks/latest_version >}} --namespace={NAMESPACE} \
  --set launcher.enabled=true \
  --set launcher.configFromGitRepo.enabled=true \
  --set launcher.configFromGitRepo.repoUrl=git@github.com:{GIT_ORG}/{GIT_REPO}.git \
  --set launcher.configFromGitRepo.branch={GIT_BRANCH} \
  --set launcher.configFromGitRepo.configFile={RELATIVE_PATH_TO_PIPED_CONFIG_FILE_IN_GIT_REPO} \
  --set launcher.configFromGitRepo.sshKeyFile=/etc/piped-secret/ssh-key \
  --set-file secret.data.piped-key={PATH_TO_PIPED_KEY_FILE} \
  --set-file secret.data.ssh-key={PATH_TO_PRIVATE_SSH_KEY_FILE} \
  --set args.enableDefaultKubernetesCloudProvider=false \
  --set rbac.scope=namespace
  {{< /tab >}}
  {{< /tabpane >}}

#### In case on OpenShift less than 4.2

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
  {{< tab lang="bash" header="Piped with Remote-upgrade" >}}
# Enable remote-upgrade feature of Piped.
# https://pipecd.dev/docs/user-guide/managing-piped/remote-upgrade-remote-config/#remote-upgrade
# This allows upgrading Piped to a new version from the web console.
# But we still need to restart Piped when we want to update its config data.

helm upgrade -i dev-piped oci://ghcr.io/pipe-cd/chart/piped --version={{< blocks/latest_version >}} --namespace={NAMESPACE} \
  --set launcher.enabled=true \
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
  --set launcher.image.repository="ghcr.io/pipe-cd/launcher-okd"
  {{< /tab >}}
  {{< tab lang="bash" header="Piped with Remote-upgrade and Remote-config" >}}
# Enable both remote-upgrade and remote-config features of Piped.
# https://pipecd.dev/docs/user-guide/managing-piped/remote-upgrade-remote-config/#remote-config
# Beside of the ability to upgrade Piped to a new version from the web console,
# remote-config allows loading the Piped config stored in a remote location such as a Git repository.
# Whenever the config data is changed, it loads the new config and restarts Piped to use that new config.

helm upgrade -i dev-piped oci://ghcr.io/pipe-cd/chart/piped --version={{< blocks/latest_version >}} --namespace={NAMESPACE} \
  --set launcher.enabled=true \
  --set launcher.configFromGitRepo.enabled=true \
  --set launcher.configFromGitRepo.repoUrl=git@github.com:{GIT_ORG}/{GIT_REPO}.git \
  --set launcher.configFromGitRepo.branch={GIT_BRANCH} \
  --set launcher.configFromGitRepo.configFile={RELATIVE_PATH_TO_PIPED_CONFIG_FILE_IN_GIT_REPO} \
  --set launcher.configFromGitRepo.sshKeyFile=/etc/piped-secret/ssh-key \
  --set-file secret.data.piped-key={PATH_TO_PIPED_KEY_FILE} \
  --set-file secret.data.ssh-key={PATH_TO_PRIVATE_SSH_KEY_FILE} \
  --set args.enableDefaultKubernetesCloudProvider=false \
  --set rbac.scope=namespace
  --set args.addLoginUserToPasswd=true \
  --set securityContext.runAsNonRoot=true \
  --set securityContext.runAsUser={UID} \
  --set securityContext.fsGroup={FS_GROUP} \
  --set securityContext.runAsGroup=0 \
  --set launcher.image.repository="ghcr.io/pipe-cd/launcher-okd"
  {{< /tab >}}
  {{< /tabpane >}}
