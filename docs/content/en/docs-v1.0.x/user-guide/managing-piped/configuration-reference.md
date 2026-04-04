---
title: "Configuration reference"
linkTitle: "Configuration reference"
weight: 9
description: >
  Learn about all the configurable fields in the `piped` configuration file.
---

This page describes all configurable fields for the Piped (`piped.config`) configuration file in PipeCD v1.

In v1, the architecture has shifted to a plugin-based model. The old `platformProviders` have been replaced by `plugins` (which specify the tool binaries to load) and `deployTargets` (where to deploy, nested under plugins). `chartRepositories` and `analysisProviders` have also been moved or removed from the top level.

### Example `piped.config`

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: my-project
  pipedID: my-piped-id
  pipedKeyFile: /etc/piped-secret/piped-key
  apiAddress: grpc.pipecd.dev:443
  plugins:
    - name: k8s_plugin
      url: file:///path/to/k8s_plugin
      port: 8081
      deployTargets:
        - name: dev-cluster
          labels:
            env: dev
          config:
            masterURL: http://cluster-dev
            kubeConfigPath: ./kubeconfig-dev
```

## Piped Configuration

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `apiVersion` | string | `pipecd.dev/v1beta1` | Yes |
| `kind` | string | `Piped` | Yes |
| `spec.projectID` | string | The identifier of the PipeCD project where this piped belongs to. | Yes |
| `spec.pipedID` | string | The generated ID for this piped. | Yes |
| `spec.pipedKeyFile` | string | The path to the file containing the generated key string for this piped. | Yes* |
| `spec.pipedKeyData` | string | Base64 encoded string of Piped key. Either `pipedKeyFile` or `pipedKeyData` must be set. | Yes* |
| `spec.name` | string | The name of this piped. | No |
| `spec.apiAddress` | string | The address used to connect to the Control Plane's API. | Yes |
| `spec.webAddress` | string | The address to the Control Plane's Web interface. | No |
| `spec.syncInterval` | duration | How often to check whether an application should be synced. Default is `1m`. | No |
| `spec.appConfigSyncInterval` | duration | How often to check whether an application configuration file should be synced. Default is `1m`. | No |
| `spec.git` | [PipedGit](#pipedgit) | Configuration for Git executable needed for Git commands. | No |
| `spec.repositories` | [][PipedRepository](#pipedrepository) | List of Git repositories this Piped should watch. | No |
| `spec.plugins` | [][PipedPlugin](#pipedplugin) | List of architectural plugins (e.g., `k8s_plugin`, `terraform_plugin`) the Piped will run. | Yes |
| `spec.notifications` | [Notifications](#notifications) | Configurations for sending deployment notifications. | No |
| `spec.secretManagement` | [SecretManagement](#secretmanagement) | Configuration for decrypting secrets in manifests. | No |
| `spec.eventWatcher` | [PipedEventWatcher](#pipedeventwatcher) | Optional settings for event watcher. | No |
| `spec.appSelector` | map[string]string | List of labels to filter all applications this piped will handle. | No |

## PipedGit

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `username` | string | The username that will be configured for `git` user. Default is `piped`. | No |
| `email` | string | The email that will be configured for `git` user. Default is `pipecd.dev@gmail.com`. | No |
| `sshConfigFilePath` | string | Where to write ssh config file. Default is `$HOME/.ssh/config`. | No |
| `host` | string | The host name. Default is `github.com`. | No |
| `hostName` | string | The hostname or IP address of the remote git server. Default is the same value with Host. | No |
| `sshKeyFile` | string | The path to the private ssh key file. This will be used to clone the source code of the specified git repositories. | No |
| `sshKeyData` | string | Base64 encoded string of SSH key. | No |
| `password` | string | The base64 encoded password for git used while cloning above Git repository via HTTPS. | No |

## PipedRepository

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `repoID` | string | Unique identifier to the repository. This must be unique in the piped scope. | Yes |
| `remote` | string | Remote address of the repository used to clone the source code. e.g. `git@github.com:org/repo.git` | Yes |
| `branch` | string | The branch will be handled. | Yes |

## PipedPlugin

Defines the external plugin binaries that this Piped agent should load to handle specific platforms.

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `name` | string | The name of the plugin (e.g., `k8s_plugin`). | Yes |
| `url` | string | Source to download the plugin binary (schemes: `file`, `https`, `oci`). | Yes |
| `port` | int | The port which the plugin listens to. | Yes |
| `config` | object | Configuration for the plugin. | No |
| `deployTargets` | [][PipedDeployTarget](#pipeddeploytarget) | The destination environments/clusters where the Piped is allowed to deploy applications. | No |

## PipedDeployTarget

Defines the target environments where applications can be deployed.

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `name` | string | The unique name of the deploy target. | Yes |
| `labels` | map[string]string | Attributes to identify the target (e.g., `env: production`). | No |
| `config` | object | The platform-specific connection configuration. | Yes |

## PipedEventWatcher

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `checkInterval` | duration | Interval to fetch the latest event and compare it. | No |
| `gitRepos` | [][PipedEventWatcherGitRepo](#pipedeventwatchergitrepo) | The configuration list of git repositories to be observed. | No |

## PipedEventWatcherGitRepo

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `repoId` | string | Id of the git repository. Must be unique. | Yes |
| `commitMessage` | string | The commit message used to push after replacing values. | No |
| `includes` | []string | The paths to EventWatcher files to be included. e.g. `foo/*.yaml`. | No |
| `excludes` | []string | The paths to EventWatcher files to be excluded. Prioritized over `includes`. | No |

## SecretManagement

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `type` | string | Which management service should be used (`KEY_PAIR`, `GCP_KMS`). | Yes |
| `config` | object | Configuration for the specified secret management type. | Yes |

## Notifications

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `routes` | []NotificationRoute | List of notification routes. | No |
| `receivers` | []NotificationReceiver | List of notification receivers. | No |

### NotificationRoute

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `name` | string | The name of the route. | Yes |
| `receiver` | string | The name of receiver who will receive all matched events. | Yes |
| `events` | []string | List of events that should be routed to the receiver. | No |
| `ignoreEvents` | []string | List of events that should be ignored. | No |
| `apps` | []string | List of applications where their events should be routed. | No |
| `ignoreApps` | []string | List of applications where their events should be ignored. | No |
| `labels` | map[string]string | List of labels where their events should be routed. | No |
| `ignoreLabels` | map[string]string | List of labels where their events should be ignored. | No |

### NotificationReceiver

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `name` | string | The name of the receiver. | Yes |
| `slack` | NotificationReceiverSlack | Configuration for slack receiver. | No |
| `webhook` | NotificationReceiverWebhook | Configuration for webhook receiver. | No |

#### NotificationReceiverSlack

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `hookURL` | string | The hookURL of a slack channel. | Yes |
| `oauthToken` | string | [The token for Slack API use.](https://api.slack.com/authentication/basics) (deprecated)| No |
| `oauthTokenData` | string | Base64 encoded string of [The token for Slack API use.](https://api.slack.com/authentication/basics) | No |
| `oauthTokenFile` | string | The path to the oauthToken file | No |
| `channelID` | string | The channel id which slack api sends to. | No |
| `mentionedAccounts` | []string | The accounts to which slack api refers. This field supports both `@username` and `username` writing styles.| No |
| `mentionedGroups` | []string | The groups to which slack api refers. This field supports both `<!subteam^groupname>` and `groupname` writing styles.| No |

#### NotificationReceiverWebhook

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `url` | string | The URL where notification event will be sent to. | Yes |
| `signatureKey` | string | The HTTP header key used to store the configured signature in each event. Default is "PipeCD-Signature". | No |
| `signatureValue` | string | The value of signature included in header of each event request. It can be used to verify the received events. | No |
| `signatureValueFile` | string | The path to the signature value file. | No |

