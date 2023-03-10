---
title: "Adding a cloud provider"
linkTitle: "Adding cloud provider"
weight: 3
description: >
  This page describes how to add a cloud provider to enable its applications.
---

> NOTE: Starting from version v0.35.0, the CloudProvider concept is being replaced by PlatformProvider. It's a name change due to the PipeCD vision improvement. __The CloudProvider configuration is marked as deprecated, please migrate your piped agent configuration to use PlatformProvider__.

PipeCD supports multiple clouds and multiple application kinds.
Cloud provider defines which cloud and where the application should be deployed to.
So while registering a new application, the name of a configured cloud provider is required.

Currently, PipeCD is supporting these five kinds of cloud providers: `KUBERNETES`, `ECS`, `TERRAFORM`, `CLOUDRUN`, `LAMBDA`.
A new cloud provider can be enabled by adding a [CloudProvider](../configuration-reference/#cloudprovider) struct to the piped configuration file.
A piped can have one or multiple cloud provider instances from the same or different cloud provider kind.

The next sections show the specific configuration for each kind of cloud provider.

### Configuring Kubernetes cloud provider

By default, piped deploys Kubernetes application to the cluster where the piped is running in. An external cluster can be connected by specifying the `masterURL` and `kubeConfigPath` in the [configuration](../configuration-reference/#cloudproviderkubernetesconfig).

And, the default resources (defined at [here](https://github.com/pipe-cd/pipecd/blob/master/pkg/app/piped/platformprovider/kubernetes/resourcekey.go)) from all namespaces of the Kubernetes cluster will be watched for rendering the application state in realtime and detecting the configuration drift. In case you want to restrict piped to watch only a single namespace, let specify the namespace in the [KubernetesAppStateInformer](../configuration-reference/#kubernetesappstateinformer) field. You can also add other resources or exclude resources to/from the watching targets by that field.

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

See [ConfigurationReference](../configuration-reference/#cloudproviderkubernetesconfig) for the full configuration.

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

See [ConfigurationReference](../configuration-reference/#cloudproviderterraformconfig) for the full configuration.

### Configuring Cloud Run cloud provider

Adding a Cloud Run provider requires the name of the Google Cloud project and the region name where Cloud Run service is running. A service account file for accessing to Cloud Run is also required if the machine running the piped does not have enough permissions to access.

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  cloudProviders:
    - name: cloudrun-dev
      type: CLOUDRUN
      config:
        project: {GCP_PROJECT}
        region: {CLOUDRUN_REGION}
        credentialsFile: {PATH_TO_THE_SERVICE_ACCOUNT_FILE}
```

See [ConfigurationReference](../configuration-reference/#cloudprovidercloudrunconfig) for the full configuration.

### Configuring Lambda cloud provider

Adding a Lambda provider requires the region name where Lambda service is running.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  cloudProviders:
    - name: lambda-dev
      type: LAMBDA
      config:
        region: {LAMBDA_REGION}
        profile: default
        credentialsFile: {PATH_TO_THE_CREDENTIAL_FILE}
```

You will generally need your AWS credentials to authenticate with Lambda. Piped provides multiple methods of loading these credentials.
It attempts to retrieve credentials in the following order:
1. From the environment variables. Available environment variables are `AWS_ACCESS_KEY_ID` or `AWS_ACCESS_KEY` and `AWS_SECRET_ACCESS_KEY` or `AWS_SECRET_KEY`.
2. From the given credentials file. (the `credentialsFile field in above sample`)
3. From the pod running in EKS cluster via STS (SecurityTokenService).
4. From the EC2 Instance Role.

Therefore, you don't have to set credentialsFile if you use the environment variables or the EC2 Instance Role. Keep in mind the IAM role/user that you use with your Piped must possess the IAM policy permission for at least `Lambda.Function` and `Lambda.Alias` resources controll (list/read/write).

See [ConfigurationReference](../configuration-reference/#cloudproviderlambdaconfig) for the full configuration.

### Configuring ECS cloud provider

Adding a ECS provider requires the region name where ECS cluster is running.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  cloudProviders:
    - name: ecs-dev
      type: ECS
      config:
        region: {ECS_CLUSTER_REGION}
        profile: default
        credentialsFile: {PATH_TO_THE_CREDENTIAL_FILE}
```

Just same as Lambda cloud provider, there are several ways to authorize Piped agent to enable it performs deployment jobs.
It attempts to retrieve credentials in the following order:
1. From the environment variables. Available environment variables are `AWS_ACCESS_KEY_ID` or `AWS_ACCESS_KEY` and `AWS_SECRET_ACCESS_KEY` or `AWS_SECRET_KEY`.
2. From the given credentials file. (the `credentialsFile field in above sample`)
3. From the pod running in EKS cluster via STS (SecurityTokenService).
4. From the EC2 Instance Role.

See [ConfigurationReference](../configuration-reference/#cloudproviderecsconfig) for the full configuration.
