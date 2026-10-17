---
title: "Configuring Terraform application"
linkTitle: "Terraform"
weight: 2
description: >
  Specific guide to configuring deployment for Terraform application.
---

## Quick Sync

By default, when the [pipeline](../../../configuration-reference/#terraform-application) was not specified, PipeCD triggers a quick sync deployment for the merged pull request.
Quick sync for a Terraform deployment does `terraform plan` and if there are any changes detected it applies those changes automatically.

## Sync with the specified pipeline

The [pipeline](../../../configuration-reference/#terraform-application) field in the application configuration is used to customize the way to do the deployment.
You can add a manual approval before doing `terraform apply` or add an analysis stage after applying the changes to determine the impact of those changes.

These are the provided stages for Terraform application you can use to build your pipeline:

- `TERRAFORM_PLAN`
  - do the terraform plan and show the changes will be applied
- `TERRAFORM_APPLY`
  - apply all the infrastructure changes

and other common stages:
- `WAIT`
- `WAIT_APPROVAL`
- `ANALYSIS`

See the description of each stage at [Customize application deployment](../../customizing-deployment/).

## Module location

Terraform module can be loaded from:

- the same git repository with the application directory, we call as a `local module`
- a different git repository, we call as a `remote module`

## Reference

See [Configuration Reference](../../../configuration-reference/#terraform-application) for the full configuration.
