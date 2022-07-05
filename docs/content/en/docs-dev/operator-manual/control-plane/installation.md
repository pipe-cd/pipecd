---
title: "Installation"
linkTitle: "Installation"
weight: 2
description: >
  This page describes how to install control plane on a Kubernetes cluster.
---

## Prerequisites

- Having a running Kubernetes cluster
- Installed [Helm](https://helm.sh/docs/intro/install/) (3.8.0 or later)

## Installation

### 1. Preparing an encryption key

PipeCD requires a key for encrypting sensitive data or signing JWT token while authenticating. You can use one of the following commands to generate an encryption key.

``` console
openssl rand 64 | base64 > encryption-key

# or
cat /dev/urandom | head -c64 | base64 > encryption-key
```

### 2. Preparing Control Plane configuration file and installing

As described at the [architecture overview](/docs/operator-manual/control-plane/architecture-overview/) page, the Control Plane's data can be stored in one of the provided fully-managed or self-managed services. So you have to decide which kind of [data store](/docs/operator-manual/control-plane/architecture-overview/#data-store) and [file store](/docs/operator-manual/control-plane/architecture-overview/#file-store) you want to use and prepare a Control Plane configuration file suitable for that choice.

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
      # Must be a service account with "Cloud Datastore User" and "Cloud Datastore Index Admin" roles
      # since PipeCD needs them to creates the needed Firestore composite indexes in the background.
      credentialsFile: /etc/pipecd-secret/firestore-service-account
  filestore:
    type: GCS
    config:
      bucket: {YOUR_BUCKET_NAME}
      # Must be a service account with "Storage Object Admin (roles/storage.objectAdmin)" role on the given bucket
      # since PipeCD need to write file object such as deployment log file to that bucket.
      credentialsFile: /etc/pipecd-secret/gcs-service-account
```

See [ConfigurationReference](/docs/operator-manual/control-plane/configuration-reference/) for the full configuration.

After all, install the Control Plane as bellow:

``` console
helm install pipecd oci://ghcr.io/pipe-cd/chart/pipecd --version {{< blocks/latest_version >}} --namespace={NAMESPACE} \
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

__Caution__: In case of using `MySQL` as Control Plane's datastore, please note that the implementation of PipeCD requires some features that only available on [MySQL v8](https://dev.mysql.com/doc/refman/8.0/en/), make sure your MySQL service is satisfied the requirement.

### 4. Accessing the PipeCD web

If your installation was including an [ingress](https://github.com/pipe-cd/pipecd/blob/master/manifests/pipecd/values.yaml#L7), the PipeCD web can be accessed by the ingress's IP address or domain.
Otherwise, private PipeCD web can be accessed by using `kubectl port-forward` to expose the installed Control Plane on your localhost:

``` console
kubectl port-forward svc/pipecd 8080 --namespace={NAMESPACE}
```

Now go to [http://localhost:8080](http://localhost:8080) on your browser, you will see a page to login to your project. But before logging in, you need to initialize a new project by following the [next section](/docs/operator-manual/control-plane/adding-a-project/).

## Production Hardening

This part provides guidance for a production hardened deployment of the control plane.

- Publishing the control plane

    You can allow external access to the control plane by enabling the [ingress](https://github.com/pipe-cd/pipecd/blob/master/manifests/pipecd/values.yaml#L7) configuration.

- End-to-End TLS

    After switching to HTTPs, do not forget to set the `api.args.secureCookie` parameter to be `true` to disallow using cookie on unsecured HTTP connection.

    Alternatively in the case of GKE Ingress, PipeCD also requires a TLS certificate for internal use. This can be a self-signed one and generated by this command:

    ``` console
    openssl req -x509 -nodes -days 3650 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN={YOUR_DOMAIN}"
    ```
    Those key and cert can be configured via [`secret.internalTLSKey.data`](https://github.com/pipe-cd/pipecd/blob/master/manifests/pipecd/values.yaml#L118) and [`secret.internalTLSCert.data`](https://github.com/pipe-cd/pipecd/blob/master/manifests/pipecd/values.yaml#L121).

    To enable internal tls connection, please set the `gateway.internalTLS.enabled` parameter to be `true`.

    Otherwise, the `cloud.google.com/app-protocols` annotation is also should be configured as the following:

    ``` yaml
    service:
      port: 443
      annotations:
        cloud.google.com/app-protocols: '{"service":"HTTP2"}'
    ```
