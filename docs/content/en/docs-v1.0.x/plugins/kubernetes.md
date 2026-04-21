---
title: "Kubernetes Plugin"
linkTitle: "Kubernetes"
weight: 10
description: >
  How to configure the Kubernetes plugin.
---

The Kubernetes plugin enables PipeCD to manage Kubernetes application deployments.

By default, piped deploys Kubernetes application to the cluster where the piped is running in. An external cluster can be connected by specifying the `masterURL` and `kubeConfigPath` in [`deployTargets`](../user-guide/managing-piped/configuration-reference/#kubernetesplugin).

Below is an example configuration for using the Kubernetes plugin:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  plugins:
    - name: kubernetes
      port: 7001
      url: https://github.com/pipe-cd/pipecd/releases/download/pkg%2Fapp%2Fpipedv1%2Fplugin%2Fkubernetes%2Fv0.1.0/kubernetes_v0.1.0_darwin_arm64
      deployTargets:
        - name: local
          config:
            kubectlVersion: 1.32.4
            kubeConfigPath: /path/to/.kube/config
```

See [Configuration Reference for Kubernetes plugin](../user-guide/managing-piped/configuration-reference/#kubernetesplugin) for complete configuration details.