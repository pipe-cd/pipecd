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

### 2. Preparing an encryption key

PipeCD requires a key for encrypting sensitive data or signing JWT token while authenticating. You can use one of the following commands to generate an encryption key.

``` console
openssl rand 64 | base64 > encryption-key

# or
cat /dev/urandom | head -c64 | base64 > encryption-key
```

### 3. Preparing control-plane configuration file and installing

As described at the [architecture overview](/docs/operator-manual/control-plane/architecture-overview/) page, the control-plane's data can be stored in one of the provided fully-managed or self-managed services. So you have to decide which kind of [data store](/docs/operator-manual/control-plane/architecture-overview/#data-store) and [file store](/docs/operator-manual/control-plane/architecture-overview/#file-store) you want to use and prepare a control-plane configuration file suitable for that choice.

#### Using Firestore and GCS

PipeCD requires a GCS bucket and service account files to access Firestore and GCS service. Here is an example of configuration file:

``` yaml
apiVersion: "pipecd.dev/v1beta1"
kind: ControlPlane
spec:
  stateKey: {RANDOM_STRING}
  datastore:
    type: FIRESTORE
    config:
      namespace: pipecd
      environment: dev
      project: {YOUR_GCP_PROJECT_NAME}
      # Must be a service account with "Cloud Datastore User" and "Cloud Datastore Index Admin" roles.
      credentialsFile: /etc/pipecd-secret/firestore-service-account
  filestore:
    type: GCS
    config:
      bucket: {YOUR_BUCKET_NAME}
      credentialsFile: /etc/pipecd-secret/gcs-service-account
```

PipeCD automatically creates needed Firestore composite indexes in the background. So you have to prepare the service account with not only `Cloud Datastore User` but `Cloud Datastore Index Admin` roles to perform this behavior.

See [ConfigurationReference](/docs/operator-manual/control-plane/configuration-reference/) for the full configuration.

After all, install the control-plane as bellow:

``` console
helm install pipecd pipecd/pipecd --version={VERSION} --namespace={NAMESPACE} \
  --set-file config.data=path-to-control-plane-configuration-file \
  --set-file secret.encryptionKey.data=path-to-encryption-key-file \
  --set-file secret.firestoreServiceAccount.data=path-to-service-account-file \
  --set-file secret.gcsServiceAccount.data=path-to-service-account-file
```

Currently, besides `Firestore` PipeCD supports other databases as its datastore such as `MySQL`. Also as for filestore, PipeCD supports `AWS S3` and `MINIO` either.

For example, in case of using `MySQL` as datastore and `MINIO` as filestore, the ControlPlane configuration will be as follow:

```yaml
apiVersion: "pipecd.dev/v1beta1"
kind: ControlPlane
spec:
  stateKey: {RANDOM_STRING}
  datastore:
    type: MYSQL
    config:
      url: {YOUR_MYSQL_ADDRESS}
      database: {YOUR_DATABASE_NAME}
  filestore:
    type: MINIO
    config:
      endpoint: {YOUR_MINIO_ADDRESS}
      bucket: {YOUR_BUCKET_NAME}
      accessKeyFile: /etc/pipecd-secret/minio-access-key
      secretKeyFile: /etc/pipecd-secret/minio-secret-key
      autoCreateBucket: true
```

You can find required configurations to use other datastores and filestores from [ConfigurationReference](/docs/operator-manual/control-plane/configuration-reference/).

__Caution__: In case of using `MySQL` as control-plane's datastore, please note that the implementation of PipeCD requires some features that only available on [MySQL v8](https://dev.mysql.com/doc/refman/8.0/en/), make sure your MySQL service is satisfied the requirement.

### 4. Accessing the PipeCD web

If your installation was including an [ingress](https://github.com/pipe-cd/manifests/blob/master/manifests/pipecd/values.yaml#L6), the PipeCD web can be accessed by the ingress's IP address or domain.
Otherwise, private PipeCD web can be accessed by using `kubectl port-forward` to expose the installed control-plane on your localhost:

``` console
kubectl port-forward svc/pipecd 8080 --namespace={NAMESPACE}
```

Now go to [http://localhost:8080](http://localhost:8080) on your browser, you will see a page to login to your project. But before logging in, you need to initialize a new project by following the [next section](/docs/operator-manual/control-plane/adding-a-project/).

## Production Hardening

This part provides guidance for a production hardened deployment of the control plane.

- Publishing the control plane

    You can allow external access to the control plane by enabling the [ingress](https://github.com/pipe-cd/manifests/blob/master/manifests/pipecd/values.yaml#L6) configuration.

- End-to-End TLS

    After switching to HTTPs, do not forget to set the `api.args.secureCookie` parameter to be `true` to disallow using cookie on unsecured HTTP connection.

    Alternatively in the case of GKE Ingress, PipeCD also requires a TLS certificate for internal use. This can be a self-signed one and generated by this command:

    ``` console
    openssl req -x509 -nodes -days 3650 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN={YOUR_DOMAIN}"
    ```
    Those key and cert can be configured via [`secret.internalTLSKey.data`](https://github.com/pipe-cd/manifests/blob/master/manifests/pipecd/values.yaml#L83) and [`secret.internalTLSCert.data`](https://github.com/pipe-cd/manifests/blob/master/manifests/pipecd/values.yaml#L86).

    To enable internal tls connection, please set the `gateway.internalTLS.enabled` parameter to be `true`.

    Otherwise, the `cloud.google.com/app-protocols` annotation is also should be configured as the following:

    ``` yaml
    service:
      port: 443
      annotations:
        cloud.google.com/app-protocols: '{"service":"HTTP2"}'
    ```
