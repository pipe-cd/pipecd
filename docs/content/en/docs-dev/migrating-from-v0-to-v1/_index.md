---
title: "Mirgrating from v0 to v1"
linkTitle: "Mirgrating from v0 to v1"
weight: 90
description: >
  Documentation on migrating PipeCD v0 deployments to PipeCD v1
---

This page explains how to migrate your existing PipeCD system to **piped v1**, the new plugin-based architecture that brings modularity and extensibility to PipeCD.

## Overview

PipeCD v1 (internally referred to as *pipedv1*) introduces a **pluggable architecture** that allows developers to add and maintain custom deployment and operational plugins without modifying the core system of PipeCD.

Migration from v0 is designed to be **safe** and **incremental**, allowing you to switch between piped and pipedv1 during the process with minimal disruption.

## Components

| Component | Description | Compatibility |
|------------|--------------|----------------|
| **Control Plane** | Manages projects, deployments, and applications. | Supports both piped and pipedv1 concurrently. |
| **Piped** | Manages the actual deployment and syncing of applications. | Backward compatible - You can switch between versions safely. |

---

## Prerequisites

Before you start, ensure that:

- You are running PipeCD **v0.54.0-rc1 or later**.
- You have the **latest Control Plane** installed.
- You have **pipectl v0.54.0-rc1 or later**.
- You have access to your Control Plane with **API write permissions**.

> **Note:** If you’re new to the plugin architecture, read:

> - [Overview of the plan for plugin-enabled PipeCD](https://pipecd.dev/blog/2024/11/28/overview-of-the-plan-for-pluginnable-pipecd/)
> - [What’s new in pipedv1](https://pipecd.dev/blog/2025/09/02/what-is-new-in-pipedv1-plugin-arch-piped/)

---

## Migration Process Overview

The migration flow involves the following steps:

1. Update `piped` and `pipectl` binaries.  
2. Convert application configurations to the v1 format.  
3. Update the application model in the Control Plane database.  
4. Update piped configuration for v1 (plugins).  
5. Deploy and verify pipedv1.  
6. Optionally, switch back to the old piped.

---

## 1. Update piped and pipectl

Install or upgrade your `pipectl` to **v0.54.0-rc1 or newer**:

```bash
# Example for upgrading pipectl
curl -Lo ./pipectl https://github.com/pipe-cd/pipecd/releases/download/v0.54.0-rc1/pipectl_<OS>_<ARCH>
chmod +x ./pipectl
mv ./pipectl /usr/local/bin/
```

Verify the version:

```bash
pipectl version
```

> **Note:** The OS and the ARCH are supposed to be replaced by your own thing. Additonally you may also want to get a differnt release of it. Check the releases here and get your URL. Example: # OS: "darwin" or "linux" | CPU_ARCH: "arm64" or "amd64"

## 2. Convert Application Configurations to v1 Format

Convert your existing `app.pipecd.yaml` configurations to the new **v1 format**:

```bash
pipectl migrate application-config \
  --config-files=path/to/app1.pipecd.yaml,path/to/app2.pipecd.yaml
```

Or specify an entire directory:

```bash
pipectl migrate application-config --dirs=path/to/apps/
```

After conversion, your `app.pipecd.yaml` files will be compatible with both `piped` and `pipedv1`.
> **Recommended:** Trigger a deployment using your current (v0.54.x) piped after converting configurations to verify backward compatibility.

## 3. Update Application Model in Control Plane Database

Use `pipectl` to update the application models stored in your Control Plane database.

You’ll need an API key with write permission. Check:

- [Generating an API key](https://pipecd.dev/docs-v0.53.x/user-guide/command-line-tool/#authentication)

After obtaining your API key, Run:

```bash
pipectl migrate database \
  --address=<control-plane-address> \
  --api-key-file=<path-to-api-key> \
  --applications=<app-1-id>,<app-2-id>
```

This adds the new `deployTargets` field to your application models, replacing the deprecated `platformProvider`.

You can confirm success by checking the **Application Details** page in the Control Plane UI — the **Deploy Targets** field should now appear.

## 4. Update piped Configuration for v1 (Plugin Definition)

In pipedv1, **platform providers** have been replaced by **plugins**. Each plugin defines its deploy targets and configuration.

### Example (Old piped config)

Your exisiting `piped` config may look similar to the following:

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

### Example (New pipedv1 config)

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

**Changes:**

- The contents of `platformProviders[].config` are now defined under
  `plugins[].deployTargets[].config`.
- Each plugin requires a `url` field that specifies where to download the plugin binary.
- Officially released plugins can be found on the [PipeCD releases page](https://github.com/pipe-cd/pipecd/releases).

- Some examples:
  - [Kubernetes Plugin v0.1.0](https://github.com/pipe-cd/pipecd/releases/tag/pkg%2Fapp%2Fpipedv1%2Fplugin%2Fkubernetes%2Fv0.1.0)

## 5. Installing pipedv1

Once your configuration is updated and validated, you're ready to install the pipedv1 binary. You can install the latest stable or any specific release of `pipedv1` depending on your deployment needs.

All official builds are published under the [Pipecd Releases](https://github.com/pipe-cd/pipecd/releases) page.

Example:

- [v1.0.0-rc5](https://github.com/pipe-cd/pipecd/releases/tag/pipedv1%2Fexp%2Fv1.0.0-rc5)
- [v1.0.0-rc3](https://github.com/pipe-cd/pipecd/releases/tag/pipedv1%2Fexp%2Fv1.0.0-rc3)

>**Warning:**
> Before switching to `pipedv1`, stop your existing Piped process to avoid conflicts, as both versions use the same `piped-id`.

Once you successfully install the `pipedv1` binary, there are a few different ways to deploy it.

### Option 1 - Deploy via Helm (Kubernetes)

If you are deploying Piped as a pod in a Kubernetes cluster, use the following Helm command:

```bash
helm upgrade -i pipedv1-exp oci://ghcr.io/pipe-cd/chart/pipedv1-exp \
  --version=v1.0.0-rc3 \
  --namespace=<NAMESPACE_TO_INSTALL_PIPEDV1> \
  --create-namespace \
  --set-file config.data=<PATH_TO_PIPEDV1_CONFIG_FILE> \
  --set-file secret.data.ssh-key=<PATH_TO_GIT_SSH_KEY>
```

### Option 2 - Run the `pipedv1` directly

Once downloaded, you can also run the `pipedv1` binary directly in your local environment:

```bash
# Download pipedv1 binary
# OS: "darwin" or "linux" | CPU_ARCH: "arm64" or "amd64"
curl -Lo ./piped https://github.com/pipe-cd/pipecd/releases/download/pipedv1%2Fexp%2Fv1.0.0-rc3/pipedv1_v1.0.0-rc3_<OS>_<CPU_ARCH>
chmod +x ./piped

# Run piped binary
./piped piped --config-file=<PATH_TO_PIPEDV1_CONFIG_FILE> --tools-dir=/tmp/piped-bin
```

If your Control Plane is running locally, append the `--insecure=true` flag to skip TLS certificate verification:

```bash
./piped piped --config-file=<PATH_TO_PIPEDV1_CONFIG_FILE> --tools-dir=/tmp/piped-bin --insecure=true
```

---

## Want to switch back?

If you need to roll back to the original `piped` process, you can do so safely at any time.

Simply **stop** the running `pipedv1` instance, then **start** the `pipedv0` service as you normally would.

> **Warning:**
> The configuration files for `piped` and `pipedv1` are **not interchangeable**.
> Make sure to restore or reference the correct configuration file before restarting pipedv0.

Once restarted, the old `pipedv0` should resume operation using its previous configuration and state. Happy Pipe-CDing!

---
