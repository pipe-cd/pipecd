---
title: "Configuration reference"
linkTitle: "Configuration reference"
weight: 11
description: >
  This page describes all configurable fields in the piped configuration.
---

> TBA

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: ...
  pipedID: ...
  ...
```

### Piped Configuration

| Field | Type | Description | Required |
|-|-|-|-|
| projectID | string | The project ID this piped belongs to | Yes |
| pipedID | string | The generated ID for this piped | Yes |
| pipedKeyFile | string | The path to the file containing generated Key for this piped | Yes |
| webURL | string | The URL address of PipeCD web | Yes |
| syncInterval | duration | How offten to check whether an application should be synced. Default is `1m` | No |
| git | [Git](/docs/operator-manual/piped/configuration-reference/#git) | Git configuration needed for Git commands  | No |
| repositories | [][Repository](/docs/operator-manual/piped/configuration-reference/#repository) | List of Git repositories this piped will handle | No |
| chartRepositories | [][ChartRepository](/docs/operator-manual/piped/configuration-reference/#chartrepository) | List of Helm chart repositories that should be added while starting up | No |
| cloudProviders | [][CloudProvider](/docs/operator-manual/piped/configuration-reference/#cloudprovider) | List of cloud providers can be used by this piped | No |
| analysisProviders | [][AnalysisProvider](/docs/operator-manual/piped/configuration-reference/#analysisprovider) | List of analysis providers can be used by this piped | No |
| notifications | [Notifications](/docs/operator-manual/piped/configuration-reference/#notifications) | Notification to Slack, Webhook... | No |

### Git

| Field | Type | Description | Required |
|-|-|-|-|
| Username | string | The username that will be configured to `git` | No |

### Repository

| Field | Type | Description | Required |
|-|-|-|-|
| repoID | string | Unique identifier to the repository. This must be unique in the piped scope | Yes |

### ChartRepository

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The name of the chart repository | Yes |

### CLoudProvider

| Field | Type | Description | Required |
|-|-|-|-|
| Name | string | The name of the cloud provider | Yes |

### AnalysisProvider

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The unique name of the analysis provider. | Yes |
| kind | string | The provider type. Currently only PROMETHEUS is available. | Yes |
| prometheus | [AnalysisProviderPrometheus](/docs/operator-manual/piped/configuration-reference/#analysisproviderprometheus) | Configuration needed to connect to Prometheus. | No |

### AnalysisProviderPrometheus
| Field | Type | Description | Required |
|-|-|-|-|
| address | string | The Prometheus server address. | Yes |
| usernameFile | string | The path to the username file. | No |
| passwordFile | string | The path to the password file. | No |

### Notifications

| Field | Type | Description | Required |
|-|-|-|-|
| routes | [][Notification.Route](/docs/operator-manual/piped/configuration-reference/#notificationroute) | List of notification routes | No |

### Notification.Route

| Field | Type | Description | Required |
|-|-|-|-|
| Name | string | The name of the analysis provider | Yes |
