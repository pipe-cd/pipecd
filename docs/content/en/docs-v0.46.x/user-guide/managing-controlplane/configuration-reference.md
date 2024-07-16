---
title: "Configuration reference"
linkTitle: "Configuration reference"
weight: 6
description: >
  This page describes all configurable fields in the Control Plane configuration.
---

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: ControlPlane
spec:
  address: https://your-pipecd-address
  ...
```

## Control Plane Configuration

| Field | Type | Description | Required |
|-|-|-|-|
| stateKey | string | A randomly generated string used to sign oauth state. | Yes |
| datastore | [DataStore](#datastore) | Storage for storing application, deployment data. | Yes |
| filestore | [FileStore](#filestore) | File storage for storing deployment logs and application states. | Yes |
| cache | [Cache](#cache) | Internal cache configuration. | No |
| address | string | The address to the control plane. This is required if SSO is enabled. | No |
| insightCollector | [InsightCollector](#insightcollector) | Option to run collector of Insights feature. | No |
| sharedSSOConfigs | [][SharedSSOConfig](#sharedssoconfig) | List of shared SSO configurations that can be used by any projects. | No |
| projects | [][Project](#project) | List of debugging/quickstart projects. Please note that do not use this to configure the projects running in the production. | No |

## DataStore

| Field | Type | Description | Required |
|-|-|-|-|
| type | string | Which type of data store should be used. Can be one of the following values<br>`FIRESTORE`, `MYSQL`. | Yes |
| config | [DataStoreConfig](#datastoreconfig) | Specific configuration for the datastore type. This must be one of these DataStoreConfig. | Yes |

## DataStoreConfig

Must be one of the following objects:

### DataStoreFireStoreConfig

| Field | Type | Description | Required |
|-|-|-|-|
| namespace | string | The root path element considered as a logical namespace, e.g. `pipecd`. | Yes |
| environment | string | The second path element considered as a logical environment, e.g. `dev`. All pipecd collections will have path formatted according to `{namespace}/{environment}/{collection-name}`. | Yes |
| collectionNamePrefix | string | The prefix for collection name. This can be used to avoid conflicts with existing collections in your Firestore database. | No |
| project | string | The name of GCP project hosting the Firestore. | Yes |
| credentialsFile | string | The path to the service account file for accessing Firestores. | No |


### DataStoreMySQLConfig

| Field | Type | Description | Required |
|-|-|-|-|
| url | string | The address to MySQL server. Should attach with the database port info as `127.0.0.1:3307` in case you want to use another port than the default value. | Yes |
| database | string | The name of database. | No (If you set it via URL) |
| usernameFile | string | Path to the file containing the username. | No |
| passwordFile | string | Path to the file containing the password. | No |


## FileStore

| Field | Type | Description | Required |
|-|-|-|-|
| type | string | Which type of file store should be used. Can be one of the following values<br>`GCS`, `S3`, `MINIO` | Yes |
| config | [FileStoreConfig](#filestoreconfig) | Specific configuration for the filestore type. This must be one of these FileStoreConfig. | Yes |

## FileStoreConfig

Must be one of the following objects:

### FileStoreGCSConfig

| Field | Type | Description | Required |
|-|-|-|-|
| bucket | string | The bucket name. | Yes |
| credentialsFile | string | The path to the service account file for accessing GCS. | No |

### FileStoreS3Config

| Field | Type | Description | Required |
|-|-|-|-|
| bucket | string | The AWS S3 bucket name. | Yes |
| region | string | The AWS region name. | Yes |
| profile | string | The AWS profile name. Default value is `default`. | No |
| credentialsFile | string | The path to AWS credential file. Requires only if you want to auth by specified credential file, by default PipeCD will use `$HOME/.aws/credentials` file. | No |
| roleARN | string | The IAM role arn to use when assuming an role. Requires only if you want to auth by `WebIdentity` pattern. | No |
| tokenFile | string | The path to the WebIdentity token PipeCD should use to assume a role with. Requires only if you want to auth by `WebIdentity` pattern. | No |

### FileStoreMinioConfig

| Field | Type | Description | Required |
|-|-|-|-|
| endpoint | string | The address of Minio. | Yes |
| bucket | string | The bucket name. | Yes |
| accessKeyFile | string | The path to the access key file. | No |
| secretKeyFile | string | The path to the secret key file. | No |
| autoCreateBucket | bool | Whether the given bucket should be made automatically if not exists. | No |

## Cache

| Field | Type | Description | Required |
|-|-|-|-|
| ttl | duration | The time that in-memory cache items are stored before they are considered as stale. | Yes |

## Project

| Field | Type | Description | Required |
|-|-|-|-|
| id | string | The unique identifier of the project. | Yes |
| desc | string | The description about the project. | No |
| staticAdmin | [ProjectStaticUser](#projectstaticuser) | Static admin account of the project. | Yes |

## ProjectStaticUser

| Field | Type | Description | Required |
|-|-|-|-|
| username | string | The username string. | Yes |
| passwordHash | string | The bcrypt hashed value of the password string. | Yes |

## InsightCollector

| Field | Type | Description | Required |
|-|-|-|-|
| application | [InsightCollectorApplication](#insightcollectorapplication) | Application metrics collector. | No |
| deployment | [InsightCollectorDeployment](#insightcollectordeployment) | Deployment metrics collector. | No |

## InsightCollectorApplication

| Field | Type | Description | Required |
|-|-|-|-|
| enabled | bool | Whether to enable. Default is `true` | No |
| schedule | string | When collector will be executed. Default is `0 * * * *` | No |

## InsightCollectorDeployment

| Field | Type | Description | Required |
|-|-|-|-|
| enabled | bool | Whether to enable. Default is `true` | No |
| schedule | string | When collector will be executed. Default is `30 * * * *` | No |
| chunkMaxCount | int | The maximum number of deployment items could be stored in a chunk. Default is `1000` | No |

## SharedSSOConfig

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The unique name of the configuration. | Yes |
| provider | string | The SSO service provider. Can be one of the following values<br>`GITHUB`, `GOOGLE`... | Yes |
| sessionTtl | int | The time to live of session for SSO login. Unit is `hour`. Default is 7 * 24 hours. | No |
| github | [SSOConfigGitHub](#ssoconfiggithub) | GitHub sso configuration. | No |

## SSOConfigGitHub

| Field | Type | Description | Required |
|-|-|-|-|
| clientId | string | The client id string of GitHub oauth app. | Yes |
| clientSecret | string | The client secret string of GitHub oauth app. | Yes |
| baseUrl | string | The address of GitHub service. Required if enterprise. | No |
| uploadUrl | string | The upload url of GitHub service. | No |
| proxyUrl | string | The address of the proxy used while communicating with the GitHub service. | No |
