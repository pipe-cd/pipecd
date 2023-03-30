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

## External Tool Management
You can manage external tools `CUSTOM_SYNC` stages use other than PipeCD supports. You have to install `asdf` because this feature manage tools by `asdf` internally. If you want to know more about asdf, please see [official documents](https://asdf-vm.com/).
You can use specified version of external tools globally according to your piped configuration file. With following setting, piped runs `asdf global aws-sam-cli 1.7.3` when piped starts running.
```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  externalTools:
    - command: "aws-sam-cli"
      version: 1.7.3
```

You can also set specified version locally to your application folder where your application configuration file exists. With following setting, piped runs `asdf local aws-sam-cli 1.77.3` in the folder when piped starts running.
```
apiVersion: pipecd.dev/v1beta1
kind: LambdaApp
spec:
...
  pipeline:
    stages:
      - name: CUSTOM_SYNC
        with:
          run: |
            sam --version
          externalTools:
            - package: "aws-sam-cli"
              version: "1.77.0"
```
