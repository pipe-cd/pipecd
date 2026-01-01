---
title: "Terraform"
linkTitle: "Terraform"
weight: 30
description: >
  Manage infrastructure as code with Terraform
---

The Terraform plugin enables PipeCD to manage infrastructure as code using Terraform with automated planning, approval workflows, and drift detection.

## Features

- **Automated planning:** Automatic `terraform plan` execution
- **Manual approval gates:** Review plans before applying
- **Drift detection:** Detect infrastructure configuration drift
- **Workspace support:** Manage multiple Terraform workspaces
- **State management:** Secure remote state handling
- **Module support:** Full support for Terraform modules

## Piped Configuration

Configure the Terraform plugin in your Piped configuration (v1):

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: dev
  pipedID: your-piped-id
  pipedKeyFile: /etc/piped-secret/piped-key
  apiAddress: your-control-plane:443
  git:
    sshKeyFile: /etc/piped-secret/ssh-key
  repositories:
    - repoId: examples
      remote: https://github.com/your-org/examples.git
      branch: master
  plugins:
    - name: terraform
      port: 7002
      url: https://github.com/pipe-cd/pipecd/releases/download/pkg/app/pipedv1/plugin/terraform/v0.2.1/terraform_linux_amd64
      deployTargets:
        - name: production
          config:
            vars:
              - project=pipecd
              - region=us-central1
            terraformVersion: 1.5.0
        - name: staging
          config:
            vars:
              - project=pipecd-staging
              - region=us-west1
            terraformVersion: 1.5.0
```

## Application Configuration

Example `.pipe.yaml` for Terraform applications (v1):

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-terraform-app
  labels:
    env: production
    team: infrastructure
  pipeline:
    stages:
      - name: TERRAFORM_PLAN
      - name: WAIT_APPROVAL
      - name: TERRAFORM_APPLY
  plugins:
    terraform:
      input:
        workspace: production
        terraformVersion: 1.5.0
```

## Available Stages

- **TERRAFORM_PLAN:** Generate and display execution plan
- **TERRAFORM_APPLY:** Apply infrastructure changes
- **TERRAFORM_SYNC:** Automatic plan and apply
- **WAIT:** Wait for a specified duration
- **WAIT_APPROVAL:** Manual approval gate

## Examples

### Simple Apply with Approval

Standard workflow with manual approval before applying:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-terraform-app
  labels:
    env: production
  pipeline:
    stages:
      - name: TERRAFORM_PLAN
      - name: WAIT_APPROVAL
      - name: TERRAFORM_APPLY
  plugins:
    terraform:
      input:
        workspace: production
```

### Auto-sync

Automatic plan and apply without approval:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-terraform-app
  labels:
    env: staging
  pipeline:
    stages:
      - name: TERRAFORM_SYNC
  plugins:
    terraform:
      input:
        workspace: staging
```

### With Custom Terraform Version

Specify a particular Terraform version:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-terraform-app
  labels:
    env: production
    team: platform
  pipeline:
    stages:
      - name: TERRAFORM_PLAN
      - name: WAIT
        with:
          duration: 5m
      - name: WAIT_APPROVAL
      - name: TERRAFORM_APPLY
  plugins:
    terraform:
      input:
        workspace: production
        terraformVersion: 1.5.0
```

## Source Code

- [`pkg/app/pipedv1/plugin/terraform/`](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin/terraform)

## See Also

- [Configuration Reference](/docs-dev/user-guide/configuration-reference/)
- [Managing Applications](/docs-dev/user-guide/managing-application/)
- [Migrating to PipeCD V1](/docs-dev/migrating-from-v0-to-v1/)