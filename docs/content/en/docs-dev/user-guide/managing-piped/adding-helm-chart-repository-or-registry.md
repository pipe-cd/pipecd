---
title: "Adding a Helm chart repository or registry"
linkTitle: "Adding Helm chart repo or registry"
weight: 5
description: >
  This page describes how to add a new Helm chart repository or registry.
---

PipeCD supports Kubernetes applications that are using Helm for templating and packaging. In addition to being able to deploy a Helm chart that is sourced from the same Git repository (`local chart`) or from a different Git repository (`remote git chart`), an application can use a chart sourced from a Helm chart repository.

### Adding Helm chart repository

A Helm [chart repository](https://helm.sh/docs/topics/chart_repository/) is a location backed by an HTTP server where packaged charts can be stored and shared. Before an application can be configured to use a chart from a Helm chart repository, that chart repository must be enabled in the related `piped` by adding the [ChartRepository](../configuration-reference/#chartrepository) struct to the piped configuration file.

``` yaml
# piped configuration file
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  chartRepositories:
    - name: pipecd
      address: https://charts.pipecd.dev
```

For example, the above snippet enables the official chart repository of PipeCD project. After that, you can configure the Kubernetes application to load a chart from that chart repository for executing the deployment.

``` yaml
# Application configuration file.
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  input:
    # Helm chart sourced from a Helm Chart Repository.
    helmChart:
      repository: pipecd
      name: helloworld
      version: v0.5.0
```

In case the chart repository is backed by HTTP basic authentication, the username and password strings are required in [configuration](../configuration-reference/#chartrepository).

### Adding Helm chart registry

A Helm chart [registry](https://helm.sh/docs/topics/registries/) is a mechanism enabled by default in Helm 3.8.0 and later that allows the OCI registry to be used for storage and distribution of Helm charts.

Before an application can be configured to use a chart from a registry, that registry must be enabled in the related `piped` by adding the [ChartRegistry](../configuration-reference/#chartregistry) struct to the piped configuration file if authentication is enabled at the registry.

``` yaml
# piped configuration file
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  chartRegistries:
    # Basic authentication
    - type: OCI
      address: registry.example.com
      authType: BASIC
      username: sample-username
      password: sample-password
      insecure: false  # Set to true to allow connections to TLS registry without certs

    # Token authentication
    - type: OCI
      address: registry.example.com
      authType: TOKEN
      token: your-token
      insecure: false

    # Certificate authentication
    - type: OCI
      address: registry.example.com
      authType: CERT
      certFile: /path/to/cert.pem
      keyFile: /path/to/key.pem
      caFile: /path/to/ca.pem  # Optional
      insecure: false
```

The following authentication methods are supported:

1. **Basic Authentication (BASIC)**: Uses username and password for authentication. This is the default authentication type. Passwords are securely passed via stdin to avoid exposure in process arguments.

2. **Token Authentication (TOKEN)**: Uses a bearer token for authentication. Tokens are securely passed via stdin to avoid exposure in process arguments.

3. **Certificate Authentication (CERT)**: Uses client certificates for mutual TLS authentication. This method requires:
   - `certFile`: Path to the client certificate file
   - `keyFile`: Path to the client key file
   - `caFile`: (Optional) Path to the CA certificate file for verifying the server's certificate

Additional options:
- `insecure`: When set to true, allows connections to TLS registry without certificates. Use with caution as this reduces security.
