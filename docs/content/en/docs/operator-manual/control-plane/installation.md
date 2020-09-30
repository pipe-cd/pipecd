---
title: "Installation"
linkTitle: "Installation"
weight: 2
description: >
  This page describes how to install control plane on a Kubernetes cluster.
---

## Prerequisites

- Having a running Kubernetes cluster
- Installed [helm3](https://helm.sh/docs/intro/install/)

## Installation

### 1. Adding helm chart repository

Installing the control-plane will be done via the helm chart sourced in [pipe-cd/manifests](https://github.com/pipe-cd/manifests/tree/master/manifests/pipecd) GitHub repository. That chart is published in the Helm chart repository at `https://charts.pipecd.dev`.

So before installing PipeCD, let's add the above Helm chart repository to your Helm client by the following command:

``` console
helm repo add pipecd https://charts.pipecd.dev
```

### 2. Preparing a signing key

PipeCD requires a key for signing JWT token while authenticating. You can use the following command to generate a signing key.

``` console
openssl rand 64 -out token-signing-key
```

### 3. Preparing control-plane configuration file and installing

As described at the [architecture overview](/docs/operator-manual/control-plane/architecture-overview/) page, the control-plane's data can be stored in one of the provided fully-managed or self-managed services. So you have to decide which kind of [data store](/docs/operator-manual/control-plane/architecture-overview/#data-store) and [file store](/docs/operator-manual/control-plane/architecture-overview/#file-store) you want to use and prepare a control-plane configuration file suitable for that choice.

#### Using Firestore and GCS

PipeCD requires a GCS bucket and service account files to access Firestore and GCS service. Here is an example of configuration file:

``` yaml
apiVersion: "pipecd.dev/v1beta1"
kind: ControlPlane
spec:
  stateKey: random-string
  datastore:
    type: FIRESTORE
    config:
      namespace: pipecd
      environment: dev
      project: gcp-project-name
      credentialsFile: /etc/pipecd-secret/firestore-service-account
  filestore:
    type: GCS
    config:
      bucket: bucket-name
      credentialsFile: /etc/pipecd-secret/gcs-service-account
```

See [ConfigurationReference](/docs/operator-manual/control-plane/configuration-reference/) for the full configuration.

After all, install the control-plane as bellow:

``` console
helm install pipecd pipecd/pipecd --version=VERSION --namespace=NAMESPACE \
  --set-file config.data=path-to-control-plane-configuration-file \
  --set-file secret.signingKey.data=path-to-signing-key-file \
  --set-file secret.firestoreServiceAccount.data=path-to-service-account-file \
  --set-file secret.gcsServiceAccount.data=path-to-service-account-file
```

#### Using DynamoDB and S3

> TBA

#### Using MongoDB and Minio

``` yaml
apiVersion: "pipecd.dev/v1beta1"
kind: ControlPlane
spec:
  stateKey: random-string
  datastore:
    type: MONGODB
    config:
      url: mongodb-address
      database: database-name
  filestore:
    type: MINIO
    config:
      endpoint: minio-address
      bucket: bucket-name
      accessKeyFile: /etc/pipecd-secret/minio-access-key
      secretKeyFile: /etc/pipecd-secret/minio-secret-key
      autoCreateBucket: true
```

See [ConfigurationReference](/docs/operator-manual/control-plane/configuration-reference/) for the full configuration.

After all, install the control-plane as bellow:

``` console
helm install pipecd pipecd/pipecd --version=VERSION --namespace=NAMESPACE \
  --set-file config.data=path-to-control-plane-configuration-file \
  --set-file secret.signingKey.data=path-to-signing-key-file \
  --set-file secret.minioAccessKey.data=path-to-minio-access-key-file
  --set-file secret.minioSecretKey.data=path-to-minio-secret-key-file
```

### 4. Accessing the PipeCD web

If your installation was including an [ingress](https://github.com/pipe-cd/manifests/blob/master/manifests/pipecd/values.yaml#L6), the PipeCD web can be accessed by the ingress's IP address or domain.
Otherwise, private PipeCD web can be accessed by using `kubectl port-forward` to expose the installed control-plane on your localhost:

``` console
kubectl port-forward svc/pipecd 8080:443
```

Point your web browser to [http://localhost:8080](http://localhost:8080), then you will see a field where you can give your [Project](https://pipecd.dev/docs/concepts/#project).
Before moving forward, you need to create a project, the [next section](/docs/operator-manual/control-plane/adding-a-project/) will help you with that.
