- Start Date: 2023-02-20
- Target Version: 0.41.4

# Summary

This RFC introduces a new way to enable users to use “custom stages” that users defined in their pipelines.

# Motivation

Currently, users can use only stages that PipeCD have already defined. But some users want to define new stages by their use-cases as bellow. 

- Deploying infrastructure by tools other than that PipeCD supports (terraform and kubernetes) such as SAM, cloud formation….
- Running integration tests
- Interacting with external systems
- Performing database migrations

# Detailed design

## Application Configuration

Users can include custom stages that users defined in a pipeline.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Lambda
spec:
  name: wait-approval
  labels:
    env: example
    team: product
  pipeline:
    stages:
      - name: SAM_DEPLOY
        runs:
         - "sam build"
         - "sam deploy -g"
```

Users can also use custom stages with default stages.
```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: wait-approval
  labels:
    env: example
    team: product
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 10%
      - name: WAIT_APPROVAL
        with:
          approvers:
            - nghialv
      - name: CUSTOM_WEB_HOOK
        runs:
          - "curl https://hooks.slack.com"
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_CANARY_CLEAN
```

## Variable/Secret Management

Users can use [encrypted-secret]([https://pipecd.dev/docs/user-guide/managing-application/secret-management/](https://pipecd.dev/docs/user-guide/managing-application/secret-management/)) and environment variables in scripts as bellow.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: CustomApp
spec:
  encryptedSecrets:
    password: encrypted-secrets
  variables:
    AWS_PROFILE: default
  runs:
    - "echo {{ .encryptedSecrets.password }} | sudo -S su"
    - "sam build"
    - "sam deploy -g --profile {{ .AWS_PROFILE }}"
```

## Binary Management

When PipeCD runs script, it find commands in the specified directory (~/.piped/tools). The field externalBinary can manage these command binaries. If command binaries are not in the directory or command version is different from specified version, PipeCD downloads commands by installScript. The install script is run in a temporary directory that PipeCD creates.

Enumerate external binaries that users want to use in a piped configuration file. They can use {{ .BinDir }} that is replacement of the directory (~/.piped/tools) where binary script should be installed and {{ .Version }} that is replacement of the value of the field `version` .

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  externalBinaries:
    - command: "sam"
      version: 1.7.3
      installScript: |
        wget https://github.com/aws/aws-sam-cli/releases/download/v{{ .Version }}/aws-sam-cli-macos-arm64.pkg
        echo {{ .encryptedSecrets.password }} | sudo -S installer -pkg aws-sam-cli-macos-arm64.pkg -target {{ .BinDir }}
        mv sam sam-{{ .Version }}
```

# Alternatives
## How to define custom stages?
1. Add bin executer to exisisting XXX_SYNC
- Must implement for every app kind.
- What app kind is SAM app?

1. Add QUICK_SYNC stage
- It's a stage with a lot of freedom, but it's inconsistent with the use cases being narrowed down.

1. Add OtherApp to manage custom stages
Users can define quick sync jobs by themselves. After PipeCD detect a new commit, PipeCD runs scripts users defined. This application is not related with platform providers, so the kind of application is “CustomApp”

```yaml
apiVersion: pipecd.dev/v1beta1
kind: OtherApp
spec:
  runs:
    - "sam build"
    - "sam deploy -g"
```

Users can make a pipeline that is composed of custom stages that users defined. They can also use stages (WAIT, WAIT_APPROVAL, ANALYSIS, ROLLBACK)that are not related with platform providers.

- Users can define quick sync without a pipeline.
- But it is difficult to understand of the meaning of "OtherApp". Every applications' kind shouold be either of platform providers.

## Where users define custom stage
1. Define custom stages only in a piped config file and refer defined stage name in an application config file.

piped config file

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  platformProviders:
    - name: sam-app
      type: CUSTOM
  customStages:
    - name: SAM_DEPLOY     
      runs:
        - "sam build"
        - "sam deploy -g --profile {{ .AWS_PROFILE }}"
    externalBinary:
      - command: "sam"
        version: 1.7.3
        installScript: |
          wget https://github.com/aws/aws-sam-cli/releases/download/v{{ .Version }}/aws-sam-cli-macos-arm64.pkg
          echo {{ .encryptedSecrets.password }} | sudo -S installer -pkg aws-sam-cli-macos-arm64.pkg -target {{ .BinDir }}
          mv sam sam-{{ .Version }}
```

application config file

```yaml
apiVersion: pipecd.dev/v1beta1
kind: CustomApp
spec:
  pipelines:
    - name: SAM_DEPLOY
     encryptedSecrets:
        password: encrypted-secrets
      variables:
        AWS_PROFILE: default
    - name: WAIT_APPROVAL
```

- It is difficult to edit custom stage setting because piped config should be placed where developers can not access easily.

1. Define custom stages only in an application config file

application config file

```yaml
apiVersion: pipecd.dev/v1beta1
kind: CustomApp
spec:
  pipelines:
    - name: CUSTOM_STAGE
      id: sum-deploy   
      encryptedSecrets:
        password: encrypted-secrets
      variables:
        AWS_PROFILE: default
      runs:
        - "sam build"
        - "sam deploy -g --profile {{ .AWS_PROFILE }}"
      externalBinary:
        - command: "sam"
          version: 1.7.3
          installScript: |
            wget https://github.com/aws/aws-sam-cli/releases/download/v{{ .Version }}/aws-sam-cli-macos-arm64.pkg
            echo {{ .encryptedSecrets.password }} | sudo -S installer -pkg aws-sam-cli-macos-arm64.pkg -target {{ .BinDir }}
            mv sam sam-{{ .Version }}
    - name: WAIT_APPROVAL
      with:
        approvers:
          - nghialv
```

- The application config file will be large and complicated as the number of custom stages increase.
- Users must write custom stage configurations in every application config file.
