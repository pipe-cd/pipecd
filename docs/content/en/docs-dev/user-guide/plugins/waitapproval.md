---
title: "Wait Approval"
linkTitle: "Wait Approval"
weight: 70
description: >
  Manual approval gates in pipelines
---

The Wait Approval plugin enables PipeCD to add manual approval gates in deployment pipelines, requiring specified users to approve deployment progression.

## Features

- **Manual Gates:** Require manual approval to proceed with deployment
- **Multiple Approvers:** Require approval from specific team members
- **Minimum Approvers:** Specify minimum number of approvals required
- **Flexible Authorization:** Define who can approve specific deployments
- **Audit Trail:** Track who approved deployments and when
- **Timeout Support:** Optional timeout for approval requests

## Piped Configuration

Configure the Wait Approval plugin in your Piped configuration:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: dev
  pipedID: xxx
  plugins:
    - name: waitapproval
      port: 7006
      url: https://github.com/pipe-cd/pipecd/releases/download/...
```

## Application Configuration

Add approval gates in `.pipe.yaml`:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-app
  labels:
    env: production
  pipeline:
    stages:
      - name: WAIT_APPROVAL
        timeout: 1h
        with:
          approvers:
            - devops-team
          minApproverNum: 1
  plugins: {}
```

## Available Stages

- **WAIT_APPROVAL:** Wait for manual approval before proceeding

## Stage Configuration

### WAIT_APPROVAL

Add a manual approval gate requiring specified approvers.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-app
  labels:
    env: production
  pipeline:
    stages:
      - name: WAIT_APPROVAL
        timeout: 2h
        with:
          approvers:
            - john@example.com
            - alice@example.com
          minApproverNum: 1
  plugins: {}
```

## Configuration Fields

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| approvers | []string | List of users/groups who can approve | Yes |
| minApproverNum | int | Minimum number of approvals needed | No (default: 1) |

## Examples

### Single Approver Gate

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: api-service
  labels:
    env: production
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 20%
      - name: WAIT_APPROVAL
        with:
          approvers:
            - devops-lead
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_CANARY_CLEAN
  plugins:
    kubernetes:
      input:
        manifests:
          - deployment.yaml
```

### Multiple Approvers with Minimum Count

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: database-migration
  labels:
    env: production
    criticality: high
  pipeline:
    stages:
      - name: WAIT_APPROVAL
        timeout: 24h
        with:
          approvers:
            - database-admin
            - platform-lead
            - devops-lead
          minApproverNum: 2
  plugins:
    terraform:
      input:
        terraformVersion: 1.5.0
```

## Source Code

- [Wait Approval Plugin](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin/waitapproval)

## See Also

- [Configuration Reference](/docs-dev/user-guide/configuration-reference/)
