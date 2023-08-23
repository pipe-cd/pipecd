---
title: "Custom Sync"
linkTitle: "Custom Sync"
weight: 4
description: >
   Specific guide for configuring Custom Sync
---

`CUSTOM_SYNC` is one stage in the pipeline and you can define scripts to deploy run in this stage.

> Note: This feature is at the alpha status and only available for application of kind AWS Lambda.

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
    # Must add this configuration to force use CUSTOM_SYNC stage.
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
1. You can use `CUSTOM_SYNC` with any current supporting application kind, but keep `alwaysUsePipeline` true to not run the application kind's default `QUICK_SYNC`.
2. Only one `CUSTOM_SYNC` stage should be used in an application pipeline.
3. The commands run with the enviroment variable `PATH` that refers `~/.piped/tools` at first.

The public piped image available in PipeCD main repo (ref: [Dockerfile](https://github.com/pipe-cd/pipecd/blob/master/cmd/piped/Dockerfile)) is based on [alpine](https://hub.docker.com/_/alpine/) and only has a few UNIX command available (ref: [piped-base Dockerfile](https://github.com/pipe-cd/pipecd/blob/master/tool/piped-base/Dockerfile)). If you want to use your commands (`sam` in the above example), you can:

- Prepare your own environment container image then add [piped binary](https://github.com/pipe-cd/pipecd/releases) to it.
- Build your own container image based on `ghcr.io/pipe-cd/piped` image.
- Manually update your running piped container (not recommended).

## Auto Rollback

When `autoRollback` is enabled, the deployment will be rolled back in the same way as [Rolling Back](../../rolling-back-a-deployment).

When the rolling back process is triggered in the pipeline including `CUSTOM_SYNC`, `CUSTOM_SYNC_ROLLBACK` stage will be added to the deployment pipeline.
`CUSTOM_SYNC_ROLLBACK` is different from `ROLLBACK` that applications set defaultly, it runs the same commands as `CUSTOM_SYNC` in the runnning commit to reverts all the applied changes.

![](/images/custom-sync-rollback.png)
