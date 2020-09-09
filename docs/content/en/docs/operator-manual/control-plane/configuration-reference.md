---
title: "Configuration reference"
linkTitle: "Configuration reference"
weight: 6
description: >
  This page describes all configurable fields in the control-plane configuration.
---

> TBA

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: ControlPlane
spec:
```

### Control Plane Configuration

| Field | Type | Description | Required |
|-|-|-|-|
| address | string | The address of control-plane api | Yes |
| stateKey | string | A randomly generated string used to sign oauth state | Yes |
| dataStore | [DataStore](/docs/operator-manual/control-plane/configuration-reference/#datastore) | Storage for storing application, deployment data | Yes |
| fileStore | [FileStore](/docs/operator-manual/control-plane/configuration-reference/#filestore) | File storage for storing deployment logs and application states | Yes |
| cache | [Cache](/docs/operator-manual/control-plane/configuration-reference/#cache) | List of cloud providers can be used by this piped | No |
| projects | [][Project](/docs/operator-manual/control-plane/configuration-reference/#project) | List of debugging/quickstart projects | No |

### DataStore

| Field | Type | Description | Required |
|-|-|-|-|
| type | string | Which type of data store should be used. Can be one of the following values<br>`FIRESTORE`, `DYNAMODB`, `MONGODB` | Yes |

### FileStore

| Field | Type | Description | Required |
|-|-|-|-|
| type | string | Which type of file store should be used. Can be one of the following values<br>`GCS`, `S3`, `MINIO` | Yes |

### Cache

| Field | Type | Description | Required |
|-|-|-|-|
| ttl | duration | The name of the chart repository | Yes |

### Project

| Field | Type | Description | Required |
|-|-|-|-|
| id | string | The unique identifier of the project | Yes |
