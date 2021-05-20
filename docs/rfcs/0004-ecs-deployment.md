- Start Date: 2020-05-19
- Target Version: 1.0.1

# Summary

This RFC proposes adding a new service deployment from PipeCD: AWS ECS deployment.

# Motivation

PipeCD aims to support a wide range of deployable services, currently, [Terraform deployment](https://pipecd.dev/docs/feature-status/#terraform-deployment) and [Lambda deployment](https://pipecd.dev/docs/feature-status/#lambda-deployment) are supported. ECS deployment meets a lot of requests we received.

# Detailed design

// TODO

# Unresolved questions

Service auto scaling is not supported when using an external deployment controller.
