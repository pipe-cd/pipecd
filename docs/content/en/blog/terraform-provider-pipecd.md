---
date: 2023-07-27
title: "Terraform Provider for PipeCD is Out!"
linkTitle: "Terraform Provider for PipeCD"
weight: 990
description: "In this post, we announce the release of Terraform provider for PipeCD."
author: Kenta Kozuka ([@kentakozuka](https://twitter.com/kenta_kozuka))
categories: ["Announcement"]
tags: ["Ecosystem", "Terraform"]
---

Now PipeCD resources can be managed by Terraform🎉

[terraform-provider-pipecd](https://github.com/pipe-cd/terraform-provider-pipecd), which is a Terraform provider for PipeCD, has been released!
Thanks to [@arabian9ts](https://github.com/arabian9ts), [@kurochan](https://twitter.com/kuro_m88), and [@sivchari](https://twitter.com/sivchari) for creating this provider!

Until now, there were two methods for adding or updating piped and applications in PipeCD: manually accessing the Control Plane through a browser or using the provided CLI, pipectl. While these methods were intuitive and easy to understand, they had the drawback of being prone to operational errors as the number of managed piped and applications increased. Even incorporating pipectl into shell scripts resulted in a procedural approach. It's not very cool to have a tool that provides GitOps capabilities but can't be managed using GitOps principles, right?

The newly released terraform-provider-pipecd addresses these issues. With this provider, you can easily deploy PipeCD resources using the declarative syntax of Terraform.

Using terraform-provider-pipecd, you can enjoy the following benefits:
1. **Easy deployment**: It simplifies the process of deploying PipeCD.
2. **Declarative code management**: By managing it as code, you can track changes and reviews through pull requests.
3. **Integration with continuous deployment**: When combined with CI/CD tools, it enables automation of configuration testing and the addition or modification of PipeCD resources.

Currently, the PipeCD resources that can be managed with terraform-provider-pipecd are:
- Piped Agent
- Application

In the future, based on community demand, more resources that can be managed might be added. Please try out the provider and provide feedback on Github.

## How to Use
While this blog post doesn't delve into detailed usage, it provides an example of registering a CloudRun Application with PipeCD using HCL code:

```hcl
// Setup the provider
terraform {
  required_providers {
    pipecd = {
      source  = "pipe-cd/pipecd"
      version = "0.1.0"
    }
  }
  required_version = ">= 1.4"
}

// Declaration of CloudRun Application
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

With this simple code, you can define the configuration information for an application. Then, just use the `terraform` command to perform an `apply`.

For more details, check out the [documentation](/docs/user-guide/terraform-provider-pipecd/) and the [example](https://github.com/pipe-cd/terraform-provider-pipecd/tree/main/example) in the repository.

Please watch our Community Meeting #2. [@arabian9ts](https://github.com/arabian9ts), a maintainer of terraform-provider-pipecd explaned how to add an application step by step.

{{< youtube B8NWgzLNe_o>}}

## Born in the Community
This project was originally developed by [@arabian9ts](https://github.com/arabian9ts), [@kurochan](https://github.com/kurochan), and [@sivchari](https://github.com/sivchari) to make the PipeCD's management better, and later donated to the pipe-cd organization. Many thanks for their contribution! (And they stay to be the maintainer of the projects to improve it more. How amazing is it?)

PipeCD is supported by the power of such community efforts. We encourage everyone to join the PipeCD community and participate.