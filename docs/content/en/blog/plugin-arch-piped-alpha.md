---
date: 2025-06-16
title: "Plugin architecture Piped alpha version has been released"
linkTitle: "Plugin architecture Piped alpha"
weight: 980
description: ""
author: Khanh Tran ([@khanhtc1202](https://github.com/khanhtc1202))
categories: ["Announcement"]
tags: ["Plugin", "New Feature"]
---

It has been months from the [Piped plugin architecture](https://pipecd.dev/blog/2024/11/28/overview-of-the-plan-for-pluginnable-pipecd/) blog, and after over a year of development, we happy to announce that the alpha version of Piped is ready 🎉

In this blog, we will show you how to run the plugin-arch piped with some built in plugins developed by the maintainers team.

_If you want to know more about the plugin-arch piped internal, please read the [pipedv1 README](https://github.com/pipe-cd/pipecd/blob/master/cmd/pipedv1/README.md) in pipecd repo._

## Prerequisites

- kubectl and a k8s cluster (They are not required if you won't use the kubernetes plugin)

## 1. Setup Control Plane

1. Run a Control Plane that your piped will connect to. If you want to run a Control Plane locally, see [How to run Control Plane locally](https://github.com/pipe-cd/pipecd/blob/master/cmd/pipecd/README.md#how-to-run-control-plane-locally).
    - The Control Plane version must be v0.52.0 or later.

2. Generate a new piped key/ID.

    2.1. Access the Control Plane console.

    2.2. Go to the piped list page. (https://{console-address}/settings/piped)

    2.4. Add a new piped via the `+ADD` button.

    2.5. Copy the generated piped ID and base64 encoded key.

## 2. Run plugin-arched piped

1. Download piped binary.

```sh
# OS="darwin" or "linux" CPU_ARCH="arm64" or "amd64"
curl -Lo ./piped_kubecon_jp_2025 https://github.com/pipe-cd/pipecd/releases/download/kubecon-jp-2025/piped_kubecon_jp_2025_{OS}_{CPU_ARCH}

chmod +x ./piped_kubecon_jp_2025
```

2. Create a piped config file like the following.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  apiAddress: {CONTROL_PLANE_API_ADDRESS} # like "localhost:8080"
  projectID: {PROJECT_ID}
  pipedID: {PIPED_ID}
  pipedKeyData: {BASE64_ENCODED_PIPED_KEY} # or use pipedKeyFile
  repositories:
    - repoID: repo1
      remote: https://github.com/your-account/your-repo
      branch: xxx
  # See https://pipecd.dev/docs/user-guide/managing-piped/configuration-reference/ for details of above.
  # platformProviders is not necessary.

  plugins:
    - name: kubernetes
      port: 7001 # Any unused port
      url: https://github.com/pipe-cd/pipecd/releases/download/kubecon-jp-2025/plugin_kubernetes_kubecon_jp_2025_darwin_arm64 # choose binary on the release for your own OS and CPU arch
      deployTargets:
        - name: cluster1
          config: 
            masterURL: https://127.0.0.1:61337   # shown by kubectl cluster-info
            kubeConfigPath: /path/to/kubeconfig
            kubectlVersion: 1.33.0
    - name: wait
      port: 7002 # Any unused port
      url: https://github.com/pipe-cd/pipecd/releases/download/kubecon-jp-2025/plugin_wait_kubecon_jp_2025_darwin_arm64 # choose binary on the release for your own OS and CPU arch

    - name: example-stage
      port: 7003 # Any unused port
      url: https://github.com/pipe-cd/community-plugins/releases/download/kubecon-jp-2025/plugin_example-stage_kubecon_jp_2025_darwin_arm64 # choose binary on the release for your own OS and CPU arch
      config:
        - commonMessage: "[common message]"
```

3. Run piped

```sh
./piped_kubecon_jp_2025 piped --config-file=/path/to/piped-config.yaml --tools-dir=/tmp/piped-bin
```

- If your Control Plane runs on local, add `--insecure=true` to the command to skip TLS certificate checks.


## 3. Deploy an application

Let's create application with kubernetes plugin.

1. Create an app.pipecd.yaml like the following.

    ```yaml
    apiVersion: pipecd.dev/v1beta1
    kind: Application
    spec:
      name: canary
      labels:
        env: example
        team: product
      pipeline:
        stages:
          # Deploy the workloads of CANARY variant. In this case, the number of
          # workload replicas of CANARY variant is 10% of the replicas number of PRIMARY variant.
          - name: K8S_CANARY_ROLLOUT
            with:
              replicas: 10%
          # Wait 10 seconds before going to the next stage.
          - name: WAIT
            with:
              duration: 10s
          # Update the workload of PRIMARY variant to the new version.
          - name: K8S_PRIMARY_ROLLOUT
          # Destroy all workloads of CANARY variant.
          - name: K8S_CANARY_CLEAN
    ```

2. Create Resources in the same dir as app.pipecd.yaml.

    deployment.yaml
    ```yaml
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: canary
      labels:
        app: canary
    spec:
      replicas: 2
      revisionHistoryLimit: 2
      selector:
        matchLabels:
          app: canary
          pipecd.dev/variant: primary
      template:
        metadata:
          labels:
            app: canary
            pipecd.dev/variant: primary
        spec:
          containers:
          - name: helloworld
            image: ghcr.io/pipe-cd/helloworld:v0.32.0
            args:
              - server
            ports:
            - containerPort: 9085
    ```

    service.yaml
    ```yaml
    apiVersion: v1
    kind: Service
    metadata:
      name: canary
    spec:
      selector:
        app: canary
      ports:
        - protocol: TCP
          port: 9085
          targetPort: 9085
    ```

3. Push the app.pipecd.yaml and the resources to your remote repository.
4. On the Control Plane console, register the application via `PIPED V1 ADD FROM SUGGESTIONS` tab.

## See also

Want to read more about the built in plugins of the new piped. Checkout the following docs

- kubernetes plugin: [README.md](https://github.com/pipe-cd/pipecd/blob/master/pkg/app/pipedv1/plugin/kubernetes/README.md)
- wait stage plugin: [README.md](https://github.com/pipe-cd/pipecd/blob/master/pkg/app/pipedv1/plugin/wait/README.md)

## And more

Another hot news, we want to share you this time. As part of PipeCD team activity at [KubeCon Japan 2025](https://events.linuxfoundation.org/kubecon-cloudnativecon-japan/program/schedule/), the [piped community plugins repository](https://github.com/pipe-cd/community-plugins) is officially opened 🎉

If you're interested, join us, discuss pipecd at KubeCon Japan and fire some PRs to community-plugins repo. We are happy to help you build your desired deployment pipeline and plugins.

Let's build and connect!
