---
title: "Configuring a Plugin"
linkTitle: "Configuring a Plugin"
weight: 2
description: >
  This page describes how to configure a plugin in PipeCD.
---

Starting PipeCD v1, you can deploy your application to multiple platforms using plugins.

A plugin represents a deployment capability (like Kubernetes, Terraform, etc.). Each plugin can have one or more `deployTargets`, where a deploy target represents the environment where your application will be deployed.

Currently, the official plugins maintained by the PipeCD Maintainers are:

- Kubernetes
- Terraform
- Analysis
- ScriptRun
- Wait
- Wait Approval

We are working towards releasing more plugins in the future.

>**Note:**  
> We also have the [PipeCD Community Plugins repository](https://github.com/pipe-cd/community-plugins) for plugins made by the PipeCD Community.

A plugin is added to the piped configuration inside the `spec.plugins` array and providing the plugin’s executable URL, the port it should run on, and any deploy targets that belong to it. For more details, see the [configuration reference for plugins](../configuration-reference/#plugins).

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  repositories:
  plugins:
    - name: plugin_name
      port: 7001
      url: url_to_plugin_binary
      deployTargets:
        - name:
          config: {}
```

Check out the latest plugin releases on [GitHub](https://github.com/pipe-cd/pipecd/releases).

---

Now, we will see how you can configure plugins for different application types.

## Configuring Kubernetes plugin

The Kubernetes plugin enables PipeCD to manage Kubernetes application deployments.

By default, piped deploys Kubernetes application to the cluster where the piped is running in. An external cluster can be connected by specifying the `masterURL` and `kubeConfigPath` in [`deployTargets`](../configuration-reference/#kubernetesplugin).

<!-- And, the default resources (defined [here](https://github.com/pipe-cd/pipecd/blob/master/pkg/app/piped/platformprovider/kubernetes/resourcekey.go)) from all namespaces of the Kubernetes cluster will be watched for rendering the application state in realtime and detecting the configuration drift. In case you want to restrict piped to watch only a single namespace, let specify the namespace in the [KubernetesAppStateInformer](../configuration-reference/#kubernetesappstateinformer) field. You can also add other resources or exclude resources to/from the watching targets by that field. -->

Below is an example configuration for using the Kubernetes plugin:

``` yaml
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

See [Configuration Reference for Kubernetes plugin](../configuration-reference/#kubernetesplugin) for complete configuration details.

### Configuring Terraform plugin

The Terraform plugin enables Piped to run Terraform-based deployments.
A deploy target represents a Terraform execution environment (e.g., dev/prod workspace, shared variables, drift detection settings).

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  plugins:
    - name: terraform
      port: 7002
      url: https://github.com/.../terraform_v0.3.0_linux_amd64
      deployTargets:
        - name: tf-dev
          config:
            vars:
              - "project=pipecd"
```

See [Configuration Reference for Terraform plugin](../configuration-reference/#terraformplugin) for complete configuration details.
