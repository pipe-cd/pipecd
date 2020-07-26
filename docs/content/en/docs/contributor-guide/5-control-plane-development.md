---
title: "Control Plane Development"
linkTitle: "Control Plane Development"
weight: 5
description: >
  This page describes where to find control-plane source code and how to run it locally for debugging.
---

## Source code structure

- [pkg/app/api](https://github.com/pipe-cd/pipe/tree/master/pkg/app/api): Contains source code for control-plane api. 
- [cmd/api](https://github.com/pipe-cd/pipe/tree/master/cmd/api): Entrypoint for binary of control-plane api.

- [pkg/app/web](https://github.com/pipe-cd/pipe/tree/master/pkg/app/web): Contains source code for control-plane web. 
- [cmd/web](https://github.com/pipe-cd/pipe/tree/master/cmd/web): Entrypoint for binary of control-plane web.

- [pkg](https://github.com/pipe-cd/pipe/tree/master/pkg): Contains shared source code for all components of both `piped` and `control-plane`.

## How to run API locally

### Using fake response data

If you use web mock response, please write the following config.

``` yaml
apiVersion: "pipecd.dev/v1beta1"
kind: ControlPlane
spec:
  datastore: {}
  filestore: {}
  cache: {}
```

You can run mock control plane in local machine as follows:

``` console
bazelisk run //cmd/api:api -- server \
  --config-file=/your-path-to-path/control-plane-mock.yaml \
  --use-fake-response=true \
  --enable-grpc-reflection=true
```

### Using a database

Prepare a configuration file in anywhere. The following is a sample configuration file.

``` yaml
apiVersion: "pipecd.dev/v1beta1"
kind: ControlPlane
spec:
  datastore:
    type: FIRESTORE
    config:
      namespace: sandbox
      project: pipecd
      credentialsFile: "/your-path-to-path/firestore-service-account-credential.json"
  filestore:
    type: GCS
    config:
      bucket: stage-logs-sandbox 
      credentialsFile: "/your-path-to-path/gcs-service-account-credential.json"
  cache:
    redisAddress: "localhost:6379"
    ttl: 5m
```

You can run control plane in local machine as follows:

``` console
bazelisk run //cmd/api:api -- server \
  --config-file=/your-path-to-path/control-plane.yaml
```

## How to run Web locally

Run the following command to start http server for web.

```
bazelisk run //cmd/web:web server
```

Then access `http://localhost:9082` on your browser.
