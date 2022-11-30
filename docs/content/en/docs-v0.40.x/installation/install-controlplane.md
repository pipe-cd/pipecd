---
title: "Install Control Plane"
linkTitle: "Install Control Plane"
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

![](/images/control-plane-components.png)
<p style="text-align: center;">
Control Plane Architecture
</p>

The Control Plane of PipeCD is constructed by several components, as shown in the above graph (for more in detail please read [Control Plane architecture overview docs](../../user-guide/managing-controlplane/architecture-overview/)). As mentioned in the graph, the PipeCD's data can be stored in one of the provided fully-managed or self-managed services. So you have to decide which kind of [data store](../../user-guide/managing-controlplane/architecture-overview/#data-store) and [file store](../../user-guide/managing-controlplane/architecture-overview/#file-store) you want to use and prepare a Control Plane configuration file suitable for that choice.

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

See [ConfigurationReference](../../user-guide/managing-controlplane/configuration-reference/) for the full configuration.

After all, install the Control Plane as bellow:

``` console
helm upgrade -i pipecd oci://ghcr.io/pipe-cd/chart/pipecd --version {{< blocks/latest_version >}} --namespace={NAMESPACE} \
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

You can find required configurations to use other datastores and filestores from [ConfigurationReference](../../user-guide/managing-controlplane/configuration-reference/).

__Caution__: In case of using `MySQL` as Control Plane's datastore, please note that the implementation of PipeCD requires some features that only available on [MySQL v8](https://dev.mysql.com/doc/refman/8.0/en/), make sure your MySQL service is satisfied the requirement.

### 3. Accessing the PipeCD web

If your installation was including an [ingress](https://github.com/pipe-cd/pipecd/blob/master/manifests/pipecd/values.yaml#L7), the PipeCD web can be accessed by the ingress's IP address or domain.
Otherwise, private PipeCD web can be accessed by using `kubectl port-forward` to expose the installed Control Plane on your localhost:

``` console
kubectl port-forward svc/pipecd 8080 --namespace={NAMESPACE}
```

Now go to [http://localhost:8080](http://localhost:8080) on your browser, you will see a page to login to your project.

Up to here, you have a installed PipeCD's Control Plane. To logging in, you need to initialize a new project.

### 4. Initialize a new project

To create a new project, you need to access to the `ops` pod in your installed PipeCD control plane, using `kubectl port-forward` command:

```console
kubectl port-forward service/pipecd-ops 9082 --namespace={NAMESPACE}
```

Then, access to [http://localhost:9082](http://localhost:9082).

On that page, you will see the list of registered projects and a link to register new projects. Registering a new project requires only a unique ID string and an optional description text.

Once a new project has been registered, a static admin (username, password) will be automatically generated for the project admin, you can use that to login via the login form in the above section.

For more about adding a new project in detail, please read the following [docs](../../user-guide/managing-controlplane/adding-a-project/).

### 4'. Upgrade Control Plane version

To upgrade the PipeCD Control Plane, preparations and commands remain as you do when installing PipeCD Control Plane. Only need to change the version flag in command to the specified version you want to upgrade your PipeCD Control Plane to.

``` console
helm upgrade -i pipecd oci://ghcr.io/pipe-cd/chart/pipecd --version {NEW_VERSION} --namespace={NAMESPACE} \
  --set-file config.data=path-to-control-plane-configuration-file \
  --set-file secret.encryptionKey.data=path-to-encryption-key-file \
  --set-file secret.firestoreServiceAccount.data=path-to-service-account-file \
  --set-file secret.gcsServiceAccount.data=path-to-service-account-file
```

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
