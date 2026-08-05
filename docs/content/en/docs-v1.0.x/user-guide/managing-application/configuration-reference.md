---
title: "Configuration reference"
linkTitle: "Configuration reference"
weight: 9
description: >
  Learn about all the configurable fields in the application configuration file.
---

This page describes all configurable fields for the application configuration file (`app.pipecd.yaml`) in PipeCD v1.

## Example `app.pipecd.yaml`

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-app
  labels:
    env: example
  description: "My example application"
  planner:
    autoRollback: true
  trigger:
    onCommit:
      paths:
        - deployment.yaml
    onOutOfSync:
      disabled: false
      minWindow: 5m
  plugins:
    kubernetes:
      input:
        manifests:
          - deployment.yaml
          - service.yaml
```

The fields under `spec.plugins.<plugin_name>` (such as `input` above) are plugin-specific. See the documentation for each plugin for the fields it accepts.

## Application Configuration

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `apiVersion` | string | `pipecd.dev/v1beta1` | Yes |
| `kind` | string | `Application` | Yes |
| `spec.name` | string | The application name. Required when the application is defined through this configuration file. | Yes* |
| `spec.labels` | map[string]string | Additional attributes to identify the application. | No |
| `spec.description` | string | Notes on the application. | No |
| `spec.planner` | [DeploymentPlanner](#deploymentplanner) | Configuration used while planning deployment. | No |
| `spec.commitMatcher` | [DeploymentCommitMatcher](#deploymentcommitmatcher) | Forcibly use QuickSync or Pipeline when a commit message matches a pattern. | No |
| `spec.pipeline` | [DeploymentPipeline](#deploymentpipeline) | Pipeline for deploying progressively. | No |
| `spec.trigger` | [Trigger](#trigger) | Trigger configuration to determine when a deployment is triggered. | No |
| `spec.postSync` | [PostSync](#postsync) | Configuration to be used once the deployment is triggered successfully. | No |
| `spec.timeout` | duration | Maximum time to execute the deployment before giving up. Default is `6h`. | No |
| `spec.encryption` | [SecretEncryption](#secretencryption) | List of encrypted secrets and targets to decrypt before using. | No |
| `spec.attachment` | [Attachment](#attachment) | List of files to attach to application manifests before using. | No |
| `spec.notification` | [DeploymentNotification](#deploymentnotification) | Additional configuration used while sending notifications to external services. | No |
| `spec.eventWatcher` | [][EventWatcherConfig](#eventwatcherconfig) | List of event watcher configurations. | No |
| `spec.driftDetection` | [DriftDetection](#driftdetection) | Configuration for drift detection. | No |
| `spec.plugins` | map[string]object | Plugin-specific configuration, keyed by plugin name (e.g., `kubernetes`, `terraform`). The value is decoded by each plugin, so see the per-plugin documentation for the fields under each. | No |

## DeploymentPlanner

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `alwaysUsePipeline` | bool | Disable auto-detecting whether to use QuickSync or Pipeline. Always uses the defined pipeline. | No |
| `autoRollback` | bool | Automatically reverts all deployment changes on failure. Default is `true`. | No |

## DeploymentCommitMatcher

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `quickSync` | string | Regular expression. Forces QuickSync when the commit message matches. | No |
| `pipeline` | string | Regular expression. Forces Pipeline when the commit message matches. | No |

## DeploymentPipeline

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `stages` | [][PipelineStage](#pipelinestage) | List of stages to run in sequence. | No |

## PipelineStage

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `name` | string | The stage name (e.g., `K8S_SYNC`, `WAIT_APPROVAL`, `ANALYSIS`). | Yes |
| `desc` | string | Human-readable description of this stage. | No |
| `timeout` | duration | Maximum time for this stage before it is cancelled. Default is `6h`. | No |
| `with` | object | Stage-specific configuration. See per-plugin documentation for available fields. | No |
| `skipOn` | [SkipOptions](#skipoptions) | Conditions under which this stage is skipped. | No |

## SkipOptions

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `commitMessagePrefixes` | []string | Skip this stage if the commit message starts with any of these prefixes. | No |
| `paths` | []string | Skip this stage if only files matching these paths changed. | No |

## Trigger

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `onCommit` | [OnCommit](#oncommit) | Trigger settings based on commit changes. | No |
| `onCommand` | [OnCommand](#oncommand) | Trigger settings based on a received SYNC command. | No |
| `onOutOfSync` | [OnOutOfSync](#onoutofsync) | Trigger settings based on OUT_OF_SYNC state. | No |
| `onChain` | [OnChain](#onchain) | Trigger settings based on a received CHAIN_SYNC command. | No |

## OnCommit

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `disabled` | bool | Exclude this application from being triggered when a new commit touches it. Default is `false`. | No |
| `paths` | []string | List of directories or files whose changes trigger a deployment. Supports regular expressions. | No |
| `ignores` | []string | List of directories or files whose changes are ignored. Supports regular expressions. | No |

## OnCommand

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `disabled` | bool | Exclude this application from being triggered when a SYNC command is received. Default is `false`. | No |

## OnOutOfSync

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `disabled` | bool | Exclude this application from being triggered when it is in OUT_OF_SYNC state. Default is `true`. | No |
| `minWindow` | duration | Minimum time that must elapse since the last deployment. Default is `5m`. | No |

## OnChain

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `disabled` | bool | Exclude this application from being triggered when a CHAIN_SYNC command is received. Default is `true`. | No |

## PostSync

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `chain` | [DeploymentChain](#deploymentchain) | Configuration for triggering a chain of deployments after this one succeeds. | No |

## DeploymentChain

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `applications` | [][ChainApplicationMatcher](#chainapplicationmatcher) | List of application matchers defining which applications to trigger as chain nodes. | Yes |

## ChainApplicationMatcher

At least one of `name`, `kind`, or `labels` must be set.

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `name` | string | Name of the application to match. | No |
| `kind` | string | Kind of the application to match. | No |
| `labels` | map[string]string | Labels of the application to match. | No |

## SecretEncryption

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `encryptedSecrets` | map[string]string | Map of secret name to encrypted value. | No |
| `decryptionTargets` | []string | List of files where the decrypted secrets will be injected. | Yes |

## Attachment

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `sources` | map[string]string | Map of name to file path containing the source data to embed. | No |
| `targets` | []string | List of files where the source data will be embedded. | Yes |

## DeploymentNotification

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `mentions` | [][NotificationMention](#notificationmention) | List of users or groups to notify for specific events. | No |

## NotificationMention

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `event` | string | The event name to trigger this mention. Use `*` for all events. | Yes |
| `slackusers` | []string | List of Slack user IDs to mention. See [Slack formatting docs](https://api.slack.com/reference/surfaces/formatting#mentioning-users). | No |
| `slackgroups` | []string | List of Slack group IDs to mention. See [Slack formatting docs](https://api.slack.com/reference/surfaces/formatting#mentioning-groups). | No |
| `slack` | []string | Deprecated. Use `slackusers` instead. | No |

## EventWatcherConfig

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `matcher` | [EventWatcherMatcher](#eventwatchermatcher) | Defines which event this watcher handles. | Yes |
| `handler` | [EventWatcherHandler](#eventwatcherhandler) | Defines how the matched event is handled. | Yes |

## EventWatcherMatcher

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `name` | string | The event name to match. | Yes |
| `labels` | map[string]string | Additional attributes to uniquely identify the event. | No |

## EventWatcherHandler

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `type` | string | Handler type. Currently only `GIT_UPDATE` is supported. | No |
| `config` | [EventWatcherHandlerConfig](#eventwatcherhandlerconfig) | Configuration for the handler. | Yes |

## EventWatcherHandlerConfig

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `commitMessage` | string | Commit message used when pushing changes. Uses a default message if not set. | No |
| `makePullRequest` | bool | Create a pull request instead of committing directly. | No |
| `replacements` | [][EventWatcherReplacement](#eventwatcherreplacement) | List of replacement targets to update when the event matches. | No |

## EventWatcherReplacement

Only one of `yamlField`, `jsonField`, `HCLField`, or `regex` may be set alongside `file`.

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `file` | string | Path to the file to update. | Yes |
| `yamlField` | string | YAML path to the field to update. Must start with `$`. e.g. `$.foo.bar[0].baz`. | No |
| `jsonField` | string | JSON path to the field to update. | No |
| `HCLField` | string | HCL path to the field to update. | No |
| `regex` | string | Regular expression specifying what to replace. Only the first capturing group `()` is replaced. e.g. `host.xz/foo/bar:(v[0-9].[0-9].[0-9])`. | No |

## DriftDetection

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `ignoreFields` | []string | List of fields to ignore when detecting drift. Format: `apiVersion:kind:namespace:name#fieldPath`. | No |
