---
title: "Configuration reference"
linkTitle: "Configuration reference"
weight: 9
description: >
  This page describes all configurable fields in the piped configuration.
---

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: ...
  pipedID: ...
  ...
```

## Piped Configuration

| Field | Type | Description | Required |
|-|-|-|-|
| projectID | string | The identifier of the PipeCD project where this piped belongs to. | Yes |
| pipedID | string | The generated ID for this piped. | Yes |
| pipedKeyFile | string | The path to the file containing the generated key string for this piped. | Yes |
| pipedKeyData | string | Base64 encoded string of Piped key. Either pipedKeyFile or pipedKeyData must be set. | Yes |
| apiAddress | string | The address used to connect to the Control Plane's API in format `host:port`. | Yes |
| syncInterval | duration | How often to check whether an application should be synced. Default is `1m`. | No |
| appConfigSyncInterval | duration | How often to check whether application configuration files should be synced. Default is `1m`. | No |
| git | [Git](#git) | Git configuration needed for Git commands. | No |
| repositories | [][Repository](#gitrepository) | List of Git repositories this piped will handle. | No |
| chartRepositories | [][ChartRepository](#chartrepository) | List of Helm chart repositories that should be added while starting up. | No |
| chartRegistries | [][ChartRegistry](#chartregistry) | List of helm chart registries that should be logged in while starting up. | No |
| cloudProviders | [][CloudProvider](#cloudprovider) | List of cloud providers can be used by this piped. This field is deprecated, use `platformProviders` instead. | No |
| platformProviders | [][PlatformProvider](#platformprovider) | List of platform providers can be used by this piped. | No |
| analysisProviders | [][AnalysisProvider](#analysisprovider) | List of analysis providers can be used by this piped. | No |
| eventWatcher | [EventWatcher](#eventwatcher) | Optional Event watcher settings. | No |
| secretManagement | [SecretManagement](#secretmanagement) | The using secret management method. | No |
| notifications | [Notifications](#notifications) | Sending notifications to Slack, Webhook... | No |
| appSelector | map[string]string | List of labels to filter all applications this piped will handle. Currently, it is only be used to filter the applications suggested for adding from the control plane. | No |

## Git

| Field | Type | Description | Required |
|-|-|-|-|
| username | string | The username that will be configured for `git` user. Default is `piped`. | No |
| email | string | The email that will be configured for `git` user. Default is `pipecd.dev@gmail.com`. | No |
| sshConfigFilePath | string | Where to write ssh config file. Default is `$HOME/.ssh/config`. | No |
| host | string | The host name. Default is `github.com`. | No |
| hostName | string | The hostname or IP address of the remote git server. Default is the same value with Host. | No |
| sshKeyFile | string | The path to the private ssh key file. This will be used to clone the source code of the specified git repositories. | No |
| sshKeyData | string | Base64 encoded string of SSH key. | No |

## GitRepository

| Field | Type | Description | Required |
|-|-|-|-|
| repoID | string | Unique identifier to the repository. This must be unique in the piped scope. | Yes |
| remote | string | Remote address of the repository used to clone the source code. e.g. `git@github.com:org/repo.git` | Yes |
| branch | string | The branch will be handled. | Yes |

## ChartRepository

| Field | Type | Description | Required |
|-|-|-|-|
| type | string | The repository type. Currently, HTTP and GIT are supported. Default is HTTP. | No |
| name | string | The name of the Helm chart repository. Note that is not a Git repository but a [Helm chart repository](https://helm.sh/docs/topics/chart_repository/). | Yes if type is HTTP |
| address | string | The address to the Helm chart repository. | Yes if type is HTTP |
| username | string | Username used for the repository backed by HTTP basic authentication. | No |
| password | string | Password used for the repository backed by HTTP basic authentication. | No |
| insecure | bool | Whether to skip TLS certificate checks for the repository or not. | No |
| gitRemote | string | Remote address of the Git repository used to clone Helm charts. | Yes if type is GIT |
| sshKeyFile | string | The path to the private ssh key file used while cloning Helm charts from above Git repository. | No |

## ChartRegistry

| Field | Type | Description | Required |
|-|-|-|-|
| type | string | The registry type. Currently, only OCI is supported. Default is OCI. | No |
| address | string | The address to the registry. | Yes |
| username | string | Username used for the registry authentication. | No |
| password | string | Password used for the registry authentication. | No |

## CloudProvider

This field is deprecated, please use [PlatformProvider](#platformprovider) instead.

## PlatformProvider

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The name of the platform provider. | Yes |
| type | string | The platform provider type. Must be one of the following values:<br>`KUBERNETES`, `TERRAFORM`, `CLOUDRUN`, `LAMBDA`. | Yes |
| config | [PlatformProviderConfig](#platformproviderconfig) | Specific configuration for the specified type of platform provider. | No |

## PlatformProviderConfig

Must be one of the following structs:

### PlatformProviderKubernetesConfig

| Field | Type | Description | Required |
|-|-|-|-|
| masterURL | string | The master URL of the kubernetes cluster. Empty means in-cluster. | No |
| kubectlVersion | string | Version of kubectl which will be used to connect to your cluster. Empty means the version set on [piped config](../user-guide/managing-piped/configuration-reference/#platformproviderkubernetesconfig) or [default version](https://github.com/pipe-cd/pipecd/blob/master/tool/piped-base/install-kubectl.sh#L24) will be used. | No |
| kubeConfigPath | string | The path to the kubeconfig file. Empty means in-cluster. | No |
| appStateInformer | [KubernetesAppStateInformer](#kubernetesappstateinformer) | Configuration for application resource informer. | No |

### PlatformProviderTerraformConfig

| Field | Type | Description | Required |
|-|-|-|-|
| vars | []string | List of variables that will be set directly on terraform commands with `-var` flag. The variable must be formatted by `key=value`. | No |

### PlatformProviderCloudRunConfig

| Field | Type | Description | Required |
|-|-|-|-|
| project | string | The GCP project hosting the Cloud Run service. | Yes |
| region | string | The region of running Cloud Run service. | Yes |
| credentialsFile | string | The path to the service account file for accessing Cloud Run service. | No |

### PlatformProviderLambdaConfig

| Field | Type | Description | Required |
|-|-|-|-|
| region | string | The region of running Lambda service. | Yes |
| credentialsFile | string | The path to the credential file for logging into AWS cluster. If this value is not provided, piped will read credential info from environment variables. It expects the format [~/.aws/credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html). | No |
| roleARN | string | The IAM role arn to use when assuming an role. Required if you want to use the AWS SecurityTokenService. | No |
| tokenFile | string | The path to the WebIdentity token the SDK should use to assume a role with. Required if you want to use the AWS SecurityTokenService. | No |
| profile | string | The profile to use for logging into AWS cluster. The default value is `default`. | No |

### PlatformProviderECSConfig

| Field | Type | Description | Required |
|-|-|-|-|
| region | string | The region of running ECS cluster. | Yes |
| credentialsFile | string | The path to the credential file for logging into AWS cluster. If this value is not provided, piped will read credential info from environment variables. It expects the format [~/.aws/credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html) | No |
| roleARN | string | The IAM role arn to use when assuming an role. Required if you want to use the AWS SecurityTokenService. | No |
| tokenFile | string | The path to the WebIdentity token the SDK should use to assume a role with. Required if you want to use the AWS SecurityTokenService. | No |
| profile | string | The profile to use for logging into AWS cluster. The default value is `default`. | No |

## KubernetesAppStateInformer

| Field | Type | Description | Required |
|-|-|-|-|
| namespace | string | Only watches the specified namespace. Empty means watching all namespaces. | No |
| includeResources | [][KubernetesResourcematcher](#kubernetesresourcematcher) | List of resources that should be added to the watching targets. | No |
| excludeResources | [][KubernetesResourcematcher](#kubernetesresourcematcher) | List of resources that should be ignored from the watching targets. | No |

### KubernetesResourceMatcher

| Field | Type | Description | Required |
|-|-|-|-|
| apiVersion | string | The APIVersion of the kubernetes resource. | Yes |
| kind | string | The kind name of the kubernetes resource. Empty means all kinds are matching. | No |

## AnalysisProvider

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The unique name of the analysis provider. | Yes |
| type | string | The provider type. Currently, only PROMETHEUS, DATADOG are available. | Yes |
| config | [AnalysisProviderConfig](#analysisproviderconfig) | Specific configuration for the specified type of analysis provider. | Yes |

## AnalysisProviderConfig

Must be one of the following structs:

### AnalysisProviderPrometheusConfig
| Field | Type | Description | Required |
|-|-|-|-|
| address | string | The Prometheus server address. | Yes |
| usernameFile | string | The path to the username file. | No |
| passwordFile | string | The path to the password file. | No |

### AnalysisProviderDatadogConfig
| Field | Type | Description | Required |
|-|-|-|-|
| address | string | The address of Datadog API server. Only "datadoghq.com", "us3.datadoghq.com", "datadoghq.eu", "ddog-gov.com" are available. Defaults to "datadoghq.com" | No |
| apiKeyFile | string | The path to the api key file. | Yes |
| applicationKeyFile | string | The path to the application key file. | Yes |
| apiKey | string | Optional: API Key for Datadog API server. | No |
| applicationKey | string | Optional: Application Key for Datadog API server | No |

## EventWatcher

| Field | Type | Description | Required |
|-|-|-|-|
| checkInterval | duration | Interval to fetch the latest event and compare it with one defined in EventWatcher config files. Defaults to `1m`. | No |
| gitRepos | [][EventWatcherGitRepo](#eventwatchergitrepo) | The configuration list of git repositories to be observed. Only the repositories in this list will be observed by Piped. | No |

### EventWatcherGitRepo

| Field | Type | Description | Required |
|-|-|-|-|
| repoId | string | Id of the git repository. This must be unique within the repos' elements. | Yes |
| commitMessage | string | The commit message used to push after replacing values. Default message is used if not given. | No |
| includes | []string | The paths to EventWatcher files to be included. Patterns can be used like `foo/*.yaml`. | No |
| excludes | []string | The paths to EventWatcher files to be excluded. Patterns can be used like `foo/*.yaml`. This is prioritized if both includes and this are given. | No |

## SecretManagement

| Field | Type | Description | Required |
|-|-|-|-|
| type | string | Which management method should be used. Default is `KEY_PAIR`. | Yes |
| config | [SecretManagementConfig](#secretmanagementconfig) | Configration for using secret management method. | Yes |

## SecretManagementConfig

Must be one of the following structs:

### SecretManagementKeyPair

| Field | Type | Description | Required |
|-|-|-|-|
| privateKeyFile | string | Path to the private RSA key file. | Yes |
| privateKeyData | string | Base64 encoded string of private RSA key. Either privateKeyFile or privateKeyData must be set. | No |
| publicKeyFile | string | Path to the public RSA key file. | Yes |
| publicKeyData | string | Base64 encoded string of public RSA key. Either publicKeyFile or publicKeyData must be set. | No |

### SecretManagementGCPKMS

> WIP

## Notifications

| Field | Type | Description | Required |
|-|-|-|-|
| routes | [][NotificationRoute](#notificationroute) | List of notification routes. | No |
| receivers | [][NotificationReceiver](#notificationreceiver) | List of notification receivers. | No |

### NotificationRoute

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The name of the route. | Yes |
| receiver | string | The name of receiver who will receive all matched events. | Yes |
| events | []string | List of events that should be routed to the receiver. | No |
| ignoreEvents | []string | List of events that should be ignored. | No |
| groups | []string | List of event groups should be routed to the receiver. | No |
| ignoreGroups | []string | List of event groups should be ignored. | No |
| apps | []string | List of applications where their events should be routed to the receiver. | No |
| ignoreApps | []string | List of applications where their events should be ignored. | No |
| labels | map[string]string | List of labels where their events should be routed to the receiver. | No |
| ignoreLabels | map[string]string | List of labels where their events should be ignored. | No |


### NotificationReceiver

| Field | Type | Description | Required |
|-|-|-|-|
| name | string | The name of the receiver. | Yes |
| slack | [NotificationReciverSlack](#notificationreceiverslack) | Configuration for slack receiver. | No |
| webhook | [NotificationReceiverWebhook](#notificationreceiverwebhook) | Configuration for webhook receiver. | No |

#### NotificationReceiverSlack

| Field | Type | Description | Required |
|-|-|-|-|
| hookURL | string | The hookURL of a slack channel. | Yes |
| oauthToken | string | [The token for slack api use.](https://api.slack.com/authentication/basics) | No |
| channelID | string | The channel id which slack api send to. | No |
| mentionedAccounts | []string | The accounts to which slack api referes. This field supports both `@username` and `username` writing styles.| No |

#### NotificationReceiverWebhook

| Field | Type | Description | Required |
|-|-|-|-|
| url | string | The URL where notification event will be sent to. | Yes |
| signatureKey | string | The HTTP header key used to store the configured signature in each event. Default is "PipeCD-Signature". | No |
| signatureValue | string | The value of signature included in header of each event request. It can be used to verify the received events. | No |
| signatureValueFile | string | The path to the signature value file. | No |
