---
title: "Adding a cloud provider"
linkTitle: "Adding cloud provider"
weight: 3
description: >
  This page describes how to add a cloud provider to enable its applications.
---

PipeCD supports multiple clouds and multiple application kinds.
Cloud provider defines which cloud and where the application should be deployed to.
So while registering a new application, the name of a configured cloud provider is required.

Currently, PipeCD is supporting these four kinds of cloud providers: `KUBERNETES`, `TERRAFORM`, `CLOUDRUN`, `LAMBDA`.
A new cloud provider can be enabled by adding a [CloudProvider](/docs/operator-manual/piped/configuration-reference/#cloudprovider) struct to the piped configuration file.
A piped can have one or multiple cloud provider instances from the same or different cloud provider kind.

The next sections show the specific configuration for each kind of cloud provider.

### Configuring Kubernetes cloud provider

By default, piped deploys Kubernetes application to the cluster where the piped is running in. An external cluster can be connected by specifying the `masterURL` and `kubeConfigPath` in the [configuration](/docs/operator-manual/piped/configuration-reference/#cloudproviderkubernetesconfig).

And, the default resources (defined at [here](https://github.com/pipe-cd/pipe/blob/master/pkg/app/piped/cloudprovider/kubernetes/resourcekey.go#L24-L74)) from all namespaces of the Kubernetes cluster will be watched for rendering the application state in realtime and detecting the configuration drift. In case you want to restrict piped to watch only a single namespace, let specify the namespace in the [KubernetesAppStateInformer](/docs/operator-manual/piped/configuration-reference/#kubernetesappstateinformer) field. You can also add other resources or exclude resources to/from the watching targets by that field.

Below configuration snippet just specifies a name and type of cloud provider. It means the cloud provider `kubernetes-dev` will connect to the Kubernetes cluster where the piped is running in, and this cloud provider watches all of the predefined resources from all namespaces inside that cluster.

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  cloudProviders:
    - name: kubernetes-dev
      type: KUBERNETES
```

See [ConfigurationReference](/docs/operator-manual/piped/configuration-reference/#cloudproviderkubernetesconfig) for the full configuration.

### Configuring Terraform cloud provider

A terraform cloud provider contains a list of shared terraform variables that will be applied while running the deployment of its applications.

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  cloudProviders:
    - name: terraform-dev
      type: TERRAFORM
      config:
        vars:
          - "project=pipecd"
```

See [ConfigurationReference](/docs/operator-manual/piped/configuration-reference/#cloudproviderterraformconfig) for the full configuration.

### Configuring CloudRun cloud provider

Adding a CloudRun provider requires the name of the Google Cloud project and the region name where CloudRun service is running. A service account file for accessing to CloudRun is also required if the machine running the piped does not have enough permissions to access.

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  cloudProviders:
    - name: cloudrun-dev
      type: CLOUDRUN
      config:
        project: gcp-project
        region: cloudrun-region
        credentialsFile: path-to-the-service-account-file
```

See [ConfigurationReference](/docs/operator-manual/piped/configuration-reference/#cloudprovidercloudrunconfig) for the full configuration.

### Configuring Lambda cloud provider

> TBA
