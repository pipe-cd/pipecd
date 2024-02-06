---
title: "Installing on ECS Fargate"
linkTitle: "Installing on ECS Fargate"
weight: 4
description: >
  This page describes how to install Piped as a task on ECS cluster backed by AWS Fargate.
---

## Prerequisites

##### Having piped's ID and Key strings
- Ensure that the `piped` has been registered and you are having its PIPED_ID and PIPED_KEY strings.
- If you are not having them, this [page](../../../user-guide/managing-controlplane/registering-a-piped/) guides you how to register a new one.

##### Preparing SSH key
- If your Git repositories are private, `piped` requires a private SSH key to access those repositories.
- Please checkout [this documentation](https://help.github.com/en/github/authenticating-to-github/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent) for how to generate a new SSH key pair. Then add the public key to your repositories. (If you are using GitHub, you can add it to Deploy Keys at the repository's Settings page.)

## Installation

- Preparing a piped configuration file as follows:

  ``` yaml
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
    # Enable this Piped to handle ECS application.
    platformProviders:
      - name: ecs-dev
        type: ECS
        config:
          region: {ECS_RUNNING_REGION}
  
    # Optional
    # Uncomment this if you want to enable this Piped to handle Terraform application.
    #  - name: terraform-dev
    #    type: TERRAFORM
  
    # Optional
    # Uncomment this if you want to enable SecretManagement feature.
    # https://pipecd.dev//docs/user-guide/managing-application/secret-management/
    # secretManagement:
    #   type: KEY_PAIR
    #   config:
    #     privateKeyData: {BASE64_ENCODED_PRIVATE_KEY}
    #     publicKeyData: {BASE64_ENCODED_PUBLIC_KEY}
  ```

- Store the above configuration data to AWS to enable using it while creating piped task. Both [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/) and [AWS Systems Manager Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html) are available to address this task.

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

- Prepare task definition for your piped task. Basically, you can just define your piped TaskDefinition as normal TaskDefinition, the only thing that needs to be beware is, in case you used [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/) to store piped configuration, to enable your piped accesses it's configuration we created as a secret on above, you need to add `secretsmanager:GetSecretValue` policy to your piped task `executionRole`. Read more in [Required IAM permissions for Amazon ECS secrets](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/specifying-sensitive-data-secrets.html).

  A sample TaskDefinition for Piped as follows

  {{< tabpane >}}
  {{< tab lang="json" header="Piped with Remote-upgrade" >}}
# Enable remote-upgrade feature of Piped.
# https://pipecd.dev/docs/user-guide/managing-piped/remote-upgrade-remote-config/#remote-upgrade
# This allows upgrading Piped to a new version from the web console.

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
      ],
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
# This just installs a Piped with the specified version.
# Whenever you want to upgrade that Piped to a new version or update its config data you have to restart it.

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
      ],
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

  Register this piped task definition and start piped task:

  ```console
  aws ecs register-task-definition --cli-input-json file://taskdef.json
  aws ecs run-task --cluster {ECS_CLUSTER} --task-definition piped
  ```

  Once the task is created, it will run continuously because of the piped execution. Since this task is run without [startedBy](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_StartTask.html#API_StartTask_RequestSyntax) setting, in case the piped is stopped, it will not automatically be restarted. To do so, you must define an ECS service to control piped task deployment.

  A sample Service definition to control piped task deployment.

  ```json
  {
    "cluster": "{ECS_CLUSTER}",
    "serviceName": "piped",
    "desiredCount": 1, # This must be 1.
    "taskDefinition": "{PIPED_TASK_DEFINITION_ARN}",
    "deploymentConfiguration": {
      "minimumHealthyPercent": 0,
      "maximumPercent": 100
    },
    "schedulingStrategy": "REPLICA",
    "launchType": "FARGATE",
    "networkConfiguration": {
      "awsvpcConfiguration": {
        "assignPublicIp": "ENABLED", # This is need to enable ECS deployment to pull piped container images.
        ...
      }
    }
  }
  ```

  Then start your piped task controller service.

  ```console
  aws ecs create-service \
    --cluster {ECS_CLUSTER} \
    --service-name piped \
    --cli-input-json file://service.json
  ```
