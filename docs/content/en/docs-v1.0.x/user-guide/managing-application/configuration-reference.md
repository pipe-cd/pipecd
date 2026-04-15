---
title: "Configuration reference"
linkTitle: "Configuration reference"
weight: 9
description: >
  Learn about all the configurable fields in the application configuration file.
---

This page describes all configurable fields in the application configuration for PipeCD v1.

Unlike previous versions, PipeCD v1 unifies all application types under a single `Application` kind. The specific platform (Kubernetes, Terraform, Cloud Run, etc.) is now defined using a platform label.

### Example `app.config`

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
labels:
  pipecd.dev/platform: KUBERNETES # Or TERRAFORM, CLOUDRUN, LAMBDA, ECS
spec:
  name: my-app
  description: "My unified v1 application"
  plugins: {}
  pipeline: {}
```

## Application Configuration

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `apiVersion` | string | `pipecd.dev/v1beta1` | Yes |
| `kind` | string | `Application` | Yes |
| `labels` | map[string]string | Additional attributes to identify applications. Must include the `pipecd.dev/platform` label to specify the platform type. | Yes |
| `spec.name` | string | The application name. | Yes |
| `spec.description` | string | Notes on the Application. | No |
| `spec.planner` | DeploymentPlanner | Configuration used while planning the deployment. | No |
| `spec.commitMatcher` | DeploymentCommitMatcher | Forcibly use QuickSync or Pipeline when commit message matched the specified pattern. | No |
| `spec.pipeline` | Pipeline | Pipeline definition for progressive delivery. | No |
| `spec.trigger` | Trigger | Configuration used to determine if a new deployment should be triggered. | No |
| `spec.postSync` | PostSync | Extra actions to execute once the deployment is triggered. | No |
| `spec.timeout` | duration | The maximum length of time to execute deployment before giving up. Default is `6h`. | No |
| `spec.encryption` | SecretEncryption | List of encrypted secrets and targets that should be decoded before using. | No |
| `spec.attachment` | Attachment | List of files that should be attached to application manifests before using. | No |
| `spec.notification` | DeploymentNotification | Additional configuration for sending notifications to external services. | No |
| `spec.eventWatcher` | []EventWatcherConfig | List of the configurations for the event watcher. | No |
| `spec.driftDetection` | [DriftDetection](#driftdetection) | Configuration for drift detection. | No |
| `spec.plugins` | map[string]any | List of the configurations for plugins. This field is plugin-specific. | No |

*(Note: The `spec.plugins` structures depend on the value of the `pipecd.dev/platform` label. See the specific Plugin documentation for deeper fields).*

## DriftDetection

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `ignoreFields` | []string | List of `apiVersion:kind:namespace:name#fieldPath` to ignore in diffs. | No |

## Planner

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `alwaysUsePipeline` | bool | Whether to always use the pipeline for deployment instead of a QuickSync. | No |
| `autoRollback` | bool | Whether to automatically rollback to the previous state when the deployment fails. | No |

## Pipeline

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `stages` | []PipelineStage | List of stages to be executed sequentially during the deployment pipeline. | Yes |

## Trigger

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `onCommit` | OnCommit | Configuration for triggering deployments upon new commits. | No |
| `onCommand` | OnCommand | Configuration for triggering deployments via manual command. | No |
| `onOutOfSync` | OnOutOfSync | Configuration for triggering deployments when drift is detected. | No |

## PostSync

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| `chain` | []PostSyncPlugin | List of plugins or tasks to execute after a successful synchronization. | Yes |
