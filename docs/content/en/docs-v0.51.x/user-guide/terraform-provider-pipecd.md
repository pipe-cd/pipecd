---
title: "PipeCD Terraform provider"
linkTitle: "PipeCD Terraform provider"
weight: 10
description: >
  This page describes how to manage PipeCD resources with Terraform using terraform-provider-pipecd.
---

Besides using web UI and command line tool, PipeCD community also provides Terraform module, [terraform-provider-pipecd](https://registry.terraform.io/providers/pipe-cd/pipecd/latest), which allows you to manage PipeCD resources.
This provider enables us to add, update, and delete PipeCD resources as  Infrastructure as Code (IaC). Storing resources as code in a version control system like Git repository ensures more reliability, security, and makes it more friendly for engineers to manage PipeCD resources with the power of Git leverage.

## Usage

### Setup Terraform provider
Add terraform block to declare that you use PipeCD Terraform provider. You need to input a controle plane's host name and API key via provider settings or environment variables. API key is available on the web UI.

```hcl
terraform {
  required_providers {
    pipecd = {
      source  = "pipe-cd/pipecd"
      version = "0.1.0"
    }
  }
  required_version = ">= 1.4"
}

provider "pipecd" {
  # pipecd_host    = "" // optional, if not set, read from environments as PIPECD_HOST
  # pipecd_api_key = "" // optional, if not set, read from environments as PIPECD_API_KEY
}
```

### Manage Piped agent
Add `pipecd_piped` resource to manage a Piped agent.

```hcl
resource "pipecd_piped" "mypiped" {
  name        = "mypiped"
  description = "This is my piped"
  id          = "my-piped-id"
}
```

### Adding a new application
Add `pipecd_application` resource to manage an application.

```hcl
// CloudRun Application
resource "pipecd_application" "main" {
  kind              = "CLOUDRUN"
  name              = "example-application"
  description       = "This is the simple application"
  platform_provider = "cloudrun-inproject"
  piped_id          = "your-piped-id"
  git = {
    repository_id = "examples"
    remote        = "git@github.com:pipe-cd/examples.git"
    branch        = "master"
    path          = "cloudrun/simple"
    filename      = "app.pipecd.yaml"
  }
}
```

### You want more?

We always want to add more needed resources into the Terraform provider. Please let the maintainers know what resources you want to add by creating issues in the [pipe-cd/terraform-provider-pipecd](https://github.com/pipe-cd/terraform-provider-pipecd/) repository. We also welcome your pull request to contribute!
