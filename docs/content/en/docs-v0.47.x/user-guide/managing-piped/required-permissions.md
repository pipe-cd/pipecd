---
title: "Required Permissions"
linkTitle: "Required Permissions"
weight: 12
description: >
    This page describes what permissions are required for a Piped to deploy applications.
---

A Piped requires some permissions to deploy applications, depending on the platform.

## For Amazon ECS

You need the following IAM actions.

- `ecs:CreateService`
- `ecs:CreateTaskSet`
- `ecs:DeleteTaskSet`
- `ecs:DeregisterTaskDefinition`
- `ecs:DescribeServices`
- `ecs:DescribeTaskSets`
- `ecs:RegisterTaskDefinition`
- `ecs:RunTask`
- `ecs:TagResource`
- `ecs:UpdateService`
- `ecs:UpdateServicePrimaryTaskSet`

- `elasticloadbalancing:DescribeListeners`
- `elasticloadbalancing:DescribeRules`
- `elasticloadbalancing:DescribeTargetGroups`
- `elasticloadbalancing:ModifyListener`
- `elasticloadbalancing:ModifyRule`

## For AWS Lambda

You need the following IAM actions.

- `lambda:CreateAlias`
- `lambda:CreateFunction`
- `lambda:GetAlias`
- `lambda:GetFunction`
- `lambda:PublishVersion`
- `lambda:TagResource`
- `lambda:UntagResource`
- `lambda:UpdateAlias`
- `lambda:UpdateFunctionCode`
- `lambda:UpdateFunctionConfiguration`

