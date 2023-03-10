---
title: "Custom Sync"
linkTitle: "Custom Sync"
weight: 11
description: >
   Specific guide for configuring Custom Sync
---

When you deploy to your infrastructure, you can use platform provider (`Terraform`, `Kuberetes`, `ECS`, `Cloud Run` and `Lambda`).
If you would like to deloy your infrastructure other than those platform providers such as AWS CloudFormation, AWS SAM or Google Cloud Deployment Manager, you can use Custom Sync.
`CUSTOM_SYNC` is one of the stage in the pipeline and you can define scripts to deploy run in this stage.

## How to configure Custom Sync
Add a `CUSTOM_SYNC` to your pipeline and write commands to deploy your infrastructure. 
The commands run in the directory where this application configuration file exists.
```yaml
apiVersion: pipecd.dev/v1beta1
kind: LambdaApp
spec:
  name: sam-simple
  labels:
    env: example
    team: abc
  planner:
    alwaysUsePipeline: true
  pipeline:
    stages:
      - name: CUSTOM_SYNC
        with:
          envs:
            AWS_PROFILE: "sample"
          run: |
            cd sam-app
            sam build
            echo y | sam deploy --profile $AWS_PROFILE
```

![](/images/custom-sync.png)

Note:
1. You can chose any application kind, but keep `alwaysUsePipeline` true not to run the application kind's default `QUICK_SYNC`.
2. PipeCD has not ever supported `CUSTOM_SYNC` stages more than one stage. Use one `CUSTOM_SYNC` stage in your pipeline.
3. If you chose `Kuberetes`, `ECS`, `Cloud Run` or `Lambda` as an application kind, you can not use rollout stages or promote stages because these stages are not related with `CUSTOM_SYNC` stages.
4. The commands run with the enviroment variable `PATH` that refers `~/.piped/tools` at first.

## Auto Rollback
When `autoRollback` is enabled, the deployment will be rolled back in the same way as [Rolling Back](../rolling-back-a-deployment).
When the rolling back process is triggered in the pipeline including `CUSTOM_SYNC`, `CUSTOM_SYNC_ROLLBACK` stage will be added to the deployment pipeline.
`CUSTOM_SYNC_ROLLBACK` is different from `ROLLBACK` that applications set defaultly, it runs the same commands as `CUSTOM_SYNC` in the runnning commit to reverts all the applied changes.

![](/images/custom-sync-rollback.png)

## External Binary Management
`CUSTOM_SYNC` stage runs script with the enviroment variable `PATH` that refers `~/.piped/tools` at first. The field externalBinary can manage binaries to install in the directory. If command binaries are not in the directory or command version is different from specified version, PipeCD downloads commands by installScript. The install script is run in a temporary directory that PipeCD creates.

Enumerate external binaries that users want to use in a piped configuration file. They can use {{ .BinDir }} that is replacement of the directory (~/.piped/tools) where binary script should be installed and {{ .Version }} that is replacement of the value of the field version.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  externalBinaries:
    - command: "sam"
      version: 1.7.3
      installScript: |
        wget https://github.com/aws/aws-sam-cli/releases/v{{ .Version }}download/aws-sam-cli-linux-x86_64.zip
        sha256sum aws-sam-cli-linux-x86_64.zip
        unzip aws-sam-cli-linux-x86_64.zip -d sam-installation
        sudo ./sam-installation/install
        mv /usr/local/bin/sam {{ .binDir }}/sam-{{ .Version }}
```
Note:
There are some cases that the binary depending on other libraries may not be installed in the `~/.piped/tools` directory because libraries' path depends on the directory where binary is installed.