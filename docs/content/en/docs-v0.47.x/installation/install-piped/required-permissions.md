---
title: "Required Permissions"
linkTitle: "Required Permissions"
weight: 6
description: >
    This page describes what permissions are required for a Piped to deploy applications.
---

A Piped requires some permissions to deploy applications, depending on the platform.

Note: If you run a piped as an ECS task, you need to attach the permissions on the piped task's `task role`, not `task execution role`.

## For ECSApp

You need IAM actions like the following example. You can restrict `Resource`.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ecs:CreateService",
                "ecs:CreateTaskSet",
                "ecs:DeleteTaskSet",
                "ecs:DeregisterTaskDefinition",
                "ecs:DescribeServices",
                "ecs:DescribeTaskSets",
                "ecs:RegisterTaskDefinition",
                "ecs:RunTask",
                "ecs:TagResource",
                "ecs:UpdateService",
                "ecs:UpdateServicePrimaryTaskSet",
                "elasticloadbalancing:DescribeListeners",
                "elasticloadbalancing:DescribeRules",
                "elasticloadbalancing:DescribeTargetGroups",
                "elasticloadbalancing:ModifyListener",
                "elasticloadbalancing:ModifyRule"
            ],
            "Resource": [
                "*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "iam:PassRole"
            ],
            "Resource": [
                "arn:aws:iam::<account-id>:role/<task-execution-role>",
                "arn:aws:iam::<account-id>:role/<task-role>"
            ]
        }
    ]
}
```

## For LambdaApp

You need IAM actions like the following example. You can restrict `Resource`.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "lambda:CreateAlias",
                "lambda:CreateFunction",
                "lambda:GetAlias",
                "lambda:GetFunction",
                "lambda:PublishVersion",
                "lambda:TagResource",
                "lambda:UntagResource",
                "lambda:UpdateAlias",
                "lambda:UpdateFunctionCode",
                "lambda:UpdateFunctionConfiguration"
            ],
            "Resource": [
                "*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "iam:PassRole"
            ],
            "Resource": [
                "arn:aws:iam::<account-id>:role/<function-role>"
            ]
        }
    ]
}
```