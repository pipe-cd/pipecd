---
date: 2025-09-03
title: "Migrate to pipedv1 (plugin-arch piped)"
linkTitle: "Migrate to pipedv1"
weight: 978
description: ""
author: Khanh Tran ([@khanhtc1202](https://github.com/khanhtc1202))
categories: ["Announcement"]
tags: ["Plugin", "New Feature"]
---

Dear PipeCD users! After a long wait and a lot of [works](https://github.com/pipe-cd/pipecd/issues/5259), we finally got the RC version of plugin-arch piped (aka pipedv1) with a supporting Control Plane and prepared migration tools to help you safely change your application from piped to pipedv1.

This post is a quick guide to help you easily migrate your PipeCD system to the RC version of PipeCD v1. The migration task is managed by this [issue](https://github.com/pipe-cd/pipecd/issues/5542), you can check it out for up to date information.

NOTE: If you have not yet familiar with plugin-arch piped, it's recommend to read [this overview pluginable pipecd](https://pipecd.dev/blog/2024/11/28/overview-of-the-plan-for-pluginnable-pipecd/) and [what is new in pipedv1](https://pipecd.dev/blog/2025/09/02/what-is-new-in-pipedv1-plugin-arch-piped/) blog posts.

## Overview of migration flows

While developing pipedv1, we always considered ensuring backward compatibility, as well as minimizing the cost of migrating to pipedv1.

There are 2 components in a PipeCD system:

- Control Plane: Guaranteed to support piped and pipedv1 simultaneously, changes have been introduced regularly in PipeCD releases.

- Piped: Backward compatibility guaranteed, can switch back to older piped versions with minimal impact to your current workflow.

To switch to pipedv1, you need to update both ControlPlane and pipedv1.

> NOTE: Please refer [v0.54.0-rc1](https://github.com/pipe-cd/pipecd/releases/tag/v0.54.0-rc1) or later as the latest version of the old PipeCD in this docs. While [v1.0.0-rc3](https://github.com/pipe-cd/pipecd/releases/tag/pipedv1%2Fexp%2Fv1.0.0-rc3) or later will be considered as the latest version of the pipedv1.

## Migrate PipeCD ControlPlane

As mentioned above, ControlPlane is guaranteed to support piped and pipedv1 simultaneously; changes have been introduced periodically in PipeCD releases.

__You need to install the latest version of PipeCD to ensure ControlPlane supports pipedv1 properly__.

## Migrate Piped and its managed applications

After ensuring your PipeCD ControlPlane is up to date, here is a tested flow by the PipeCD dev team to migrate your piped and the applications managed by it to pipedv1.

### 0. Update piped and pipectl version

Please update your piped version to [v0.54.0-rc1](https://github.com/pipe-cd/pipecd/releases/tag/v0.54.0-rc1) or later to ensure piped support v1-formated application configuration. You can update your version of piped in any way you usually would.

The [pipectl](https://pipecd.dev/docs-v0.53.x/user-guide/command-line-tool/) migrate command is prepared tool that helps you while migrating to pipedv1.

```bash
$ pipectl migrate -h                   
Do migration tasks.

Usage:
  pipectl migrate [command]

Available Commands:
  application-config Do migration tasks for application config.
  database           Do migration tasks for database.
```

Please use `pipectl` version [v0.54.0-rc1](https://github.com/pipe-cd/pipecd/releases/tag/v0.54.0-rc1) or later for the next following steps.

### 1. Convert application configuration to v1 format

To convert application configuration file (aka. `app.pipecd.yaml`) to v1-formated application configuration, you can run the following command

```bash
$ pipectl migrate application-config --config-files=path-to-app1-config-file,path-to-app2-config-file
```

You can also use `--dirs` flag to specify directories that contain multiple applications at once.

The latest version of piped (aka [v0.54.0-rc1](https://github.com/pipe-cd/pipecd/releases/tag/v0.54.0-rc1)) contains support for application configuration in v1 format. This means you can safely use piped with a v1-formatted application configuration.

NOTE: After migrating your application configuration file, it's recommended that you trigger deployment for that application to ensure the current piped version you're using is backward compatible. If the application deployment finishes successfully, for later steps of the migration, you can switch your piped back to using the old version instead of pipedv1 at any time.

### 2. Update application model in ControlPlane database

For this step, you need to obtain an pipectl [API key](https://pipecd.dev/docs-v0.53.x/user-guide/command-line-tool/#authentication) with the write permisson on PipeCD ControlPlane.

To update application model in ControlPlane database, you can run the following command

```bash
$ pipectl migrate database --address=control-plane-address --api-key-file=api-key-file --applications=app-1-id,app-2-id
```

This command will add `deployTargets` to the application model object in the database, based on its current `platformProvider` value.

The expected change after this step can be check via application detail page. The UI should show your application `Deploy Targets` instead of `Platform Provider`.

### 3. Prepare pipedv1 settings and update to pipedv1

The step requires some manual update on your current piped config file.

Suppose you have piped configuration file like below

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: dev
  pipedID: piped-id
  pipedKeyData: piped-key-data
  apiAddress: controlplane.pipecd.dev:443
  git:
    sshKeyFile: /etc/piped-secret/ssh-key
  repositories:
    - repoId: examples
      remote: https://github.com/pipe-cd/examples.git
      branch: master
  syncInterval: 1m
  platformProviders:
    - name: kubernetes-dev
      type: KUBERNETES
      config:
        kubectlVersion: 1.32.4
        kubeConfigPath: /users/home/.kube/config
```

You need to add configuration for all plugins that you want your pipedv1 to use while perform application deployment

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: dev
  pipedID: piped-id
  pipedKeyData: piped-key-data
  apiAddress: controlplane.pipecd.dev:443
  git:
    sshKeyFile: /etc/piped-secret/ssh-key
  repositories:
    - repoId: examples
      remote: https://github.com/pipe-cd/examples.git
      branch: master
  syncInterval: 1m
  plugins:
    - name: kubernetes
      port: 7001
      url: https://github.com/pipe-cd/pipecd/releases/download/k8s-plugin
      deployTargets:
        - name: local
          config:
            kubectlVersion: 1.32.4
            kubeConfigPath: /users/home/.kube/config
    - name: wait
      port: 7002
      url: https://github.com/pipe-cd/pipecd/releases/download/wait-plugin
```

All `platformProviders[].config` is moved to under `plugins[].deployTargets[].config`, so basically you can just move it around.

The `plugins[].url` configuration is the path to where you can download the plugin binary. Currently published official plugins can be found at the [pipecd repo release tab](https://github.com/pipe-cd/pipecd/releases).

For example, [kubernetes plugin v0.1.0](https://github.com/pipe-cd/pipecd/releases/tag/pkg%2Fapp%2Fpipedv1%2Fplugin%2Fkubernetes%2Fv0.1.0).

Next, you can install the pipedv1 version [v1.0.0-rc3](https://github.com/pipe-cd/pipecd/releases/tag/pipedv1%2Fexp%2Fv1.0.0-rc3) with your prepared piped config above.

> NOTE: Be sure to stop your current piped before switching to pipedv1, since we use the same piped-id.

If you're familiar with piped installation as pod in Kubernetes cluster, use the following command to install pipedv1

```bash
$ helm upgrade -i pipedv1-exp oci://ghcr.io/pipe-cd/chart/pipedv1-exp \
       --version=v1.0.0-rc3 \
       --namespace=<NAMESPACE_TO_INSTALL_PIPEDV1> \
       --create-namespace \
       --set-file config.data=<PATH_TO_PIPEDV1_CONFIG_FILE> \
       --set-file secret.data.ssh-key=<PATH_TO_GIT_SSH_KEY>
```

If you prefer to run pipedv1 binary on your environment directly, check the following commands

```bash
# Download pipedv1 binary
# OS="darwin" or "linux" CPU_ARCH="arm64" or "amd64"
$ curl -Lo ./piped https://github.com/pipe-cd/pipecd/releases/download/pipedv1%2Fexp%2Fv1.0.0-rc3/pipedv1_v1.0.0-rc3_<OS>_<CPU_ARCH>
$ chmod +x ./piped

# Run piped binary
$ piped piped --config-file=<PATH_TO_PIPEDV1_CONFIG_FILE> --tools-dir=/tmp/piped-bin
```

If your Control Plane runs on local, add `--insecure=true` to the command to skip TLS certificate checks.

### 4. Trigger deployment for your v1-formated application

When no errors occur in the above steps, from this step, you have a pipedv1 ready to work with features preserved from the old piped version. You can use your PipeCD as usual. Let's trigger some deployments or register a new application!

NOTE: `Add from suggestion` feature works on pipedv1 as well, but be sure your repository contains `Kind: Application` application configuration.

### 4'. Switch back to the old piped

If you want to switch back to using the old piped, you can do it at any time by stopping the pipedv1 and starting the piped as you used to. Please note that the piped configuration is different from the pipedv1 configuration.

NOTE: You can keep your application configuration in v1 format since piped versions later than [v0.54.0-rc1](https://github.com/pipe-cd/pipecd/releases/tag/v0.54.0-rc1) support v1-formatted application configuration.

## Conclusion

Above is a basic guide to help you migrate from piped to pipedv1. Let's try pipedv1, and let us know if there is anything you are interested in via the GitHub [issue tracker](https://github.com/pipe-cd/pipecd/issues).

Happy PipeCD-ing!
