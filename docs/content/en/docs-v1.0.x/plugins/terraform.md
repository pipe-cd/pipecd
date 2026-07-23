---
title: "Terraform plugin"
linkTitle: "Terraform"
weight: 20
description: >
  Apply infrastructure changes with Terraform.
---

The `terraform` plugin runs Terraform-based deployments. It reads the Terraform files in the application directory and applies them through a pipeline, running `terraform plan` and `terraform apply` as pipeline stages. A deploy target represents a Terraform execution environment: its variables and drift-detection setting.

## Prerequisites

1. **Register the plugin in the piped configuration.** Add a `terraform` plugin block with one or more `deployTargets`:

   ```yaml
   apiVersion: pipecd.dev/v1beta1
   kind: Piped
   spec:
     # ...
     plugins:
       - name: terraform
         port: 7002
         url: file:///path/to/plugin/binary  # or an https:// release URL
         deployTargets:
           - name: dev
             config:
               vars:
                 - "project=pipecd"
               driftDetectionEnabled: true
   ```

2. **The plugin downloads `terraform` automatically.** `piped` fetches the `terraform` binary via the tool registry, so it does not need to be pre-installed on the `piped` host. The version defaults to `0.13.0`; pin a different one with `terraformVersion` in the application configuration.

3. Put the application's Terraform files (`*.tf`) in the application directory alongside `app.pipecd.yaml`.

## Quick sync

With no `pipeline` defined, the plugin performs a **quick sync** (`TERRAFORM_APPLY`): it applies any detected changes to reach the state described by the Terraform files. This minimal `app.pipecd.yaml` deploys the Terraform files in the application directory:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: TerraformApp
spec:
  name: my-infra
  plugins:
    terraform:
      workspace: dev
      vars:
        - "project=pipecd"
```

## Sync with the specified pipeline

Define a `pipeline` to run `terraform plan` before applying. A common pattern is to plan, pause for a manual approval, then apply:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: TerraformApp
spec:
  name: my-infra
  pipeline:
    stages:
      - name: TERRAFORM_PLAN
      - name: WAIT_APPROVAL
        with:
          approvers:
            - user-a
      - name: TERRAFORM_APPLY
  plugins:
    terraform:
      workspace: dev
```

See [Pipeline stages](#pipeline-stages) for every available stage and its options.

## Pipeline stages

Stages are listed under `spec.pipeline.stages`, with options under `with`.

### TERRAFORM_PLAN

Runs `terraform plan` and shows the changes the deployment would apply. Use it before `TERRAFORM_APPLY` to review the plan (optionally behind a `WAIT_APPROVAL` stage).

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| exitOnNoChanges | bool | Exit the pipeline with a success status when the plan detects no changes. | false |

### TERRAFORM_APPLY

Runs `terraform apply` to apply the changes described by the Terraform files. This is also the stage that runs during a quick sync when no `pipeline` block is defined. It takes no options.

### TERRAFORM_ROLLBACK

Restores the infrastructure to the previous commit's Terraform files by running `terraform apply` against them. This stage is triggered automatically when a deployment fails or is cancelled with rollback. You do not add it to your pipeline; PipeCD inserts it automatically.

## Livestate and drift detection

The plugin detects drift by running `terraform plan` and comparing the actual infrastructure against the Terraform files in Git. If the plan reports changes, the application is marked `OUT_OF_SYNC`; otherwise it is `SYNCED`. Drift detection is controlled per deploy target by `driftDetectionEnabled`, which is enabled by default. It is an evolving feature and its behaviour may change in future releases.

## Plan preview

Before a pipeline runs, plan preview runs `terraform plan` and shows the changes the deployment would apply, so you can review them on the pull request.

## Configuration reference

### TerraformApplicationSpec

The `spec` of a `TerraformApp` shares the [common application fields](../user-guide/managing-application/configuration-reference/) (`name`, `labels`, `pipeline`, ...) and adds the following under `plugins.terraform`:

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| workspace | string | The Terraform workspace name. Empty means the `default` workspace. | No |
| terraformVersion | string | Version of `terraform` to use. Empty means the default version (`0.13.0`). | No |
| vars | []string | Variables passed to `terraform` commands with `-var`. Each entry is `key=value` (e.g. `image_id=ami-abc123`). | No |
| varFiles | []string | Variable files passed to `terraform` commands with `-var-file`. | No |
| commandFlags | [TerraformCommandFlags](#terraformcommandflags) | Additional flags passed to `terraform` commands. | No |
| commandEnvs | [TerraformCommandEnvs](#terraformcommandenvs) | Additional environment variables set while running `terraform` commands. | No |

### TerraformCommandFlags

Extra flags appended to the underlying `terraform` commands.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| shared | []string | Flags applied to every `terraform` command. | No |
| init | []string | Flags applied to `terraform init`. | No |
| plan | []string | Flags applied to `terraform plan`. | No |
| apply | []string | Flags applied to `terraform apply`. | No |

### TerraformCommandEnvs

Extra environment variables set when running the underlying `terraform` commands.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| shared | []string | Environment variables applied to every `terraform` command. | No |
| init | []string | Environment variables applied to `terraform init`. | No |
| plan | []string | Environment variables applied to `terraform plan`. | No |
| apply | []string | Environment variables applied to `terraform apply`. | No |

### DeployTargetConfig

Configured under `plugins[].deployTargets[].config` in the **piped** configuration, one per Terraform execution environment.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| vars | []string | Variables passed to `terraform` commands with `-var` for this target. Each entry is `key=value` (e.g. `image_id=ami-abc123`). | No |
| driftDetectionEnabled | bool | Enable drift detection for this target. Default is `true`. | No |
