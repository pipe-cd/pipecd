---
title: "Installing on ECS Fargate"
linkTitle: "Installing on ECS Fargate"
weight: 4
description: >
  This page describes how to install `piped` as a task on an ECS cluster backed by AWS Fargate.
---

## Prerequisites

### A registered `piped`

- Make sure your `piped` is registered in the Control Plane and that you have its **PIPED_ID** and **PIPED_KEY**.  
- If not, follow the guide to [register a new `piped`](../../../user-guide/managing-controlplane/registering-a-piped/).

### SSH key for Git repositories

- If your Git repositories are private, `piped` requires a private SSH key to access those repositories.
- Please check out [this documentation](https://help.github.com/en/github/authenticating-to-github/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent) for how to generate a new SSH key pair. Then add the public key to your repositories. (If you are using GitHub, you can add it to Deploy Keys at the repository's **Settings** page.)

## Installation

### Preparing the `piped` configuration file

Prepare a `piped` configuration file as follows:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: {PROJECT_ID}
  pipedID: {PIPED_ID}
  pipedKeyData: {BASE64_ENCODED_PIPED_KEY}
  # Write in a format like "host:443" because the communication is done via gRPC.
  apiAddress: {CONTROL_PLANE_API_ADDRESS}

  git:
    sshKeyData: {BASE64_ENCODED_PRIVATE_SSH_KEY}

  repositories:
    - repoId: {REPO_ID_OR_NAME}
      remote: git@github.com:{GIT_ORG}/{GIT_REPO}.git
      branch: {GIT_BRANCH}

  # Optional
  # Enable this piped to handle ECS applications.
  platformProviders:
    - name: ecs-dev
      type: ECS
      config:
        region: {ECS_RUNNING_REGION}

  # Optional
  # Uncomment this if you want to enable this piped to handle Terraform applications.
  #  - name: terraform-dev
  #    type: TERRAFORM

  # Optional
  # Uncomment this if you want to enable Secret Management.
  # See: https://pipecd.dev/docs/user-guide/managing-application/secret-management/
  # secretManagement:
  #   type: KEY_PAIR
  #   config:
  #     privateKeyData: {BASE64_ENCODED_PRIVATE_KEY}
  #     publicKeyData: {BASE64_ENCODED_PUBLIC_KEY}
```

See the [configuration reference](../../../user-guide/managing-piped/configuration-reference/) for the full configuration.

### Storing the configuration in AWS

Store the above configuration data in AWS so it can be referenced from your Fargate task. You can use either [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/) or [AWS Systems Manager Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html).

{{< tabpane >}}
{{< tab lang="bash" header="Store in AWS Secrets Manager" >}}
aws secretsmanager create-secret --name PipedConfig \
  --description "Configuration of piped running as ECS Fargate task" \
  --secret-string `base64 piped-config.yaml`
{{< /tab >}}
{{< tab lang="bash" header="Store in AWS Systems Manager Parameter Store" >}}
aws ssm put-parameter \
  --name PipedConfig \
  --value `base64 piped-config.yaml` \
  --type SecureString
{{< /tab >}}
{{< /tabpane >}}

If you use AWS Secrets Manager, make sure your task role or execution role has permission to read the secret. See [Required IAM permissions for Amazon ECS secrets](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/specifying-sensitive-data-secrets.html) for more details.

### Defining the task definition

Prepare a task definition for your `piped` task. The following examples show how to configure `piped` to read its configuration from AWS Secrets Manager.

{{< tabpane >}}
{{< tab lang="json" header="Piped with Remote-upgrade" >}}
# Enable remote-upgrade feature of piped.
# https://pipecd.dev/docs/user-guide/managing-piped/remote-upgrade-remote-config/#remote-upgrade
# This allows upgrading piped to a new version from the web console.

{
  "family": "piped",
  "executionRoleArn": "{PIPED_TASK_EXECUTION_ROLE_ARN}",
  "containerDefinitions": [
    {
      "name": "piped",
      "essential": true,
      "image": "ghcr.io/pipe-cd/launcher:{{< blocks/latest_version >}}",
      "entryPoint": [
        "sh",
        "-c"
      ],
      "command": [
        "/bin/sh -c \"launcher launcher --config-data=$(echo $CONFIG_DATA)\""
      ],
      "secrets": [
        {
          "valueFrom": "{PIPED_SECRET_MANAGER_ARN}",
          "name": "CONFIG_DATA"
        }
      ]
    }
  ],
  "requiresCompatibilities": [
    "FARGATE"
  ],
  "networkMode": "awsvpc",
  "memory": "512",
  "cpu": "256"
}
{{< /tab >}}
{{< tab lang="json" header="Piped" >}}
# This just installs a piped with the specified version.
# Whenever you want to upgrade that piped to a new version or update its config data you have to restart it.

{
  "family": "piped",
  "executionRoleArn": "{PIPED_TASK_EXECUTION_ROLE_ARN}",
  "containerDefinitions": [
    {
      "name": "piped",
      "essential": true,
      "image": "ghcr.io/pipe-cd/piped:{{< blocks/latest_version >}}",
      "entryPoint": [
        "sh",
        "-c"
      ],
      "command": [
        "/bin/sh -c \"piped piped --config-data=$(echo $CONFIG_DATA)\""
      ],
      "secrets": [
        {
          "valueFrom": "{PIPED_SECRET_MANAGER_ARN}",
          "name": "CONFIG_DATA"
        }
      ]
    }
  ],
  "requiresCompatibilities": [
    "FARGATE"
  ],
  "networkMode": "awsvpc",
  "memory": "512",
  "cpu": "256"
}
{{< /tab >}}
{{< /tabpane >}}

Register the task definition and run a `piped` task:

```console
aws ecs register-task-definition --cli-input-json file://taskdef.json
aws ecs run-task --cluster {ECS_CLUSTER} --task-definition piped
```

Once the task is created, it will run continuously because of the `piped` execution. Since this task is run without [`startedBy`](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_StartTask.html#API_StartTask_RequestSyntax), if `piped` stops, it will not automatically be restarted. To keep `piped` running, define an ECS service to control the task deployment.

### Defining the ECS service

The following is a sample ECS service definition to control the `piped` task deployment:

```json
{
  "cluster": "{ECS_CLUSTER}",
  "serviceName": "piped",
  "desiredCount": 1, 
  "taskDefinition": "{PIPED_TASK_DEFINITION_ARN}",
  "deploymentConfiguration": {
    "minimumHealthyPercent": 0,
    "maximumPercent": 100
  },
  "schedulingStrategy": "REPLICA",
  "launchType": "FARGATE",
  "networkConfiguration": {
    "awsvpcConfiguration": {
      "assignPublicIp": "ENABLED",
      "...": "..."
    }
  }
}
```

Then create the ECS service:

```console
aws ecs create-service \
  --cluster {ECS_CLUSTER} \
  --service-name piped \
  --cli-input-json file://service.json
```

When the service is running, ECS will ensure that exactly one `piped` task is running (because `desiredCount` is 1), keeping your `piped` agent available on Fargate.

