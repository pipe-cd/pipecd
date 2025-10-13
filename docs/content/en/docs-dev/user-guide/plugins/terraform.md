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

Configure the Terraform plugin in your Piped configuration:


```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  platforms:
    - name: terraform-default
      type: TERRAFORM
      config:
        vars:
          - "project=pipecd"
          - "region=us-central1"
```

## Application Configuration

Example `.pipe.yaml` for Terraform applications:


```yaml
apiVersion: pipecd.dev/v1beta1
kind: TerraformApp
spec:
  input:
    workspace: production
    terraformVersion: 1.5.0
  pipeline:
    stages:
      - name: TERRAFORM_PLAN
      - name: WAIT_APPROVAL
        with:
          approvers:
            - user1@example.com
      - name: TERRAFORM_APPLY
```


## Available Stages

- **TERRAFORM_PLAN:** Generate and display execution plan
- **TERRAFORM_APPLY:** Apply infrastructure changes
- **TERRAFORM_SYNC:** Automatic plan and apply

## Examples

### Simple Apply with Approval

```yaml
apiVersion: pipecd.dev/v1beta1
kind: TerraformApp
spec:
  pipeline:
    stages:
      - name: TERRAFORM_PLAN
      - name: WAIT_APPROVAL
      - name: TERRAFORM_APPLY
```

### Auto-sync

```yaml
apiVersion: pipecd.dev/v1beta1
kind: TerraformApp
spec:
  pipeline:
    stages:
      - name: TERRAFORM_SYNC
```

## Source Code

- [`pkg/app/pipedv1/plugin/terraform/`](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin/terraform)

## See Also

- [Configuration Reference](/docs-dev/user-guide/configuration-reference/#terraform-application)
- [Managing Applications](/docs-dev/user-guide/managing-application/)
