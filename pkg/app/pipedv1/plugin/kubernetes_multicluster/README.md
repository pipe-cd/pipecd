# Kubernetes Multicluster Plugin

## Specification

The current specification is described in the [RFC](../../../../../docs/rfcs/0014-multi-cluster-deployment-for-k8s.md).  
The configuration format is unstable and may change in the future.

Note:
- Currently, only QuickSync is supported.

## Try k8s multicluster plugin locally

**Switch to the upstream commit of the master branch**

The whole fixes for plugins are stored in the master branch and some of the versions don't have them.

```sh
git switch master
git pull
```

**Prepare the PipeCD Control Plane**

Please refer to [pipe-cd/pipecd/cmd/pipecd/README.md](../../../../../cmd/pipecd/README.md) to set up the Control Plane in your local environment.

**Prepare two k8s clusters**

```sh
kind create cluster --name cluster1
kind export kubeconfig --name cluster1 --kubeconfig /path/to/kubeconfig/for/cluster1

kind create cluster --name cluster2
kind export kubeconfig --name cluster2 --kubeconfig /path/to/kubeconfig/for/cluster2
```

**Start pipedv1 locally**

Please refer to [pipe-cd/pipecd/cmd/pipedv1/README.md](../../../../../cmd/pipedv1/README.md) to set up the Control Plane in your local environment.  
At this time, please modify the `spec.plugins` section of the piped config as shown below.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  plugins:
  - name: kubernetes_multicluster
    port: 7002 # any unused port
    url: file:///path/to/.piped/plugins/kubernetes_multicluster # It's OK using any value for now because it's a dummy. We will implement it later.
    deployTargets:
    - name: cluster1
      config:
        masterURL: https://127.0.0.1:61337   # shown by kubectl cluster-info
        kubeConfigPath: /path/to/kubeconfig/for/cluster1
    - name: cluster2
      config:
        masterURL: https://127.0.0.1:62082   # shown by kubectl cluster-info
        kubeConfigPath: /path/to/kubeconfig/for/cluster2
```

**Prepare the manifest**

- Please create a new repository for the manifest because an error will occur if there is an app.pipecd.yaml in the format before supporting the plugin mechanism for now.
For example usage, see the [Examples](#examples) section below.

**Register the application**

![adding-application](./docs/static/adding-application.png)

At this time, please select multiple DeployTargets.

## Examples
There are examples under `./example`.

| Name | Description |
|------|-------------|
| [simple](./example/simple/) | Deploy the same resources to the multiple clusters. |
| [multi-sources-template-none](./example/multi-sources-template-none/) | Deploy the different resources to the multiple clusters. |

## Config Reference

### Piped Config

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  plugins:
    # Plugin name.
  - name: kubernetes_multicluster
    # Port number used to listen plugin server.
    port: 7002
    # The URL where the plugin binary located.
    # It's OK using any value for now because it's a dummy. We will implement it later.
    url: file:///path/to/.piped/plugins/kubernetes_multicluster
    # List of the information for each target platform.
    # This is alternative for the platform providers.
    deployTargets:
      # Then name of deploy target.
    - name: cluster1
      # The plugin-specific config.
      config:
        # The master URL of the kubernetes cluster.
        # Empty means in-cluster.
        masterURL: https://127.0.0.1:61337
        # The path to the kubeconfig file.
        # Empty means in-cluster.
        kubeConfigPath: /path/to/kubeconfig/for/cluster1
        # Version of kubectl will be used.
        kubectlVersion: 1.32.0
    - name: cluster2
      config:
        masterURL: https://127.0.0.1:62082
        kubeConfigPath: /path/to/kubeconfig/for/cluster2
```

### app.pipecd.yaml

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  ...
  plugins:
    kubernetes_multicluster:
      input:
        # List of manifest files in the application directory used to deploy.
        # Empty means all manifest files in the directory will be used.
        manifests:
          - deployment.yaml
          - service.yaml
        # Version of kubectl will be used.
        kubectlVersion: 1.32.0
        # The namespace where manifests will be applied.
        namespace: example
        # Automatically create a new namespace if it does not exist.
        # Default is false.
        autoCreateNamespace: false
        # List of the setting for each deploy targets.
        # You can also set config for them separately.
        multiTargets:
            # The identity of the deploy target.
            # You can specify deploy target by some of them.
          - target:
              # The name of deploy target
              name: cluster1
            # List of manifest files in the application directory used to deploy.
            # Empty means all manifest files in the directory will be used.
            manifests:
              - ./cluster1/deployment.yaml
            # Version of kubectl will be used.
            kubectlVersion: 1.32.2
          - target:
              name: cluster2
            manifests:
              - ./cluster2/deployment.yaml
            kubectlVersion: 1.32.2
```
