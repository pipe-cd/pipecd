---
title: "Script Run stage"
linkTitle: "Script Run stage"
weight: 4
description: >
   Specific guide for configuring Script Run stage
---

`SCRIPT_RUN` stage is one stage in the pipeline and you can execute any commands.

> Note: This feature is at the alpha status and only available for application of kind AWS Lambda.

## How to configure SCRIPT_RUN stage

Add a `SCRIPT_RUN` to your pipeline and write commands. 
The commands run in the directory where this application configuration file exists.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: canary-with-script-run
  labels:
    env: example
    team: product
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 10%
      - name: WAIT
        with:
          duration: 10s
      - name: SCRIPT_RUN
        with:
          env:
            MSG: "execute script1"
          run: |
            echo $MSG
            sleep 10
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_CANARY_CLEAN
      - name: SCRIPT_RUN
        with:
          env:
            MSG: "execute script2"
          run: |
            echo $MSG
            sleep 10
```

![](/images/custom-sync.png)

Note:
1. You can use `SCRIPT_RUN` stage with any current supporting application kind.


2. The public piped image available in PipeCD main repo (ref: [Dockerfile](https://github.com/pipe-cd/pipecd/blob/master/cmd/piped/Dockerfile)) is based on [alpine](https://hub.docker.com/_/alpine/) and only has a few UNIX command available (ref: [piped-base Dockerfile](https://github.com/pipe-cd/pipecd/blob/master/tool/piped-base/Dockerfile)). If you want to use your commands, you can:

- Prepare your own environment container image then add [piped binary](https://github.com/pipe-cd/pipecd/releases) to it.
- Build your own container image based on `ghcr.io/pipe-cd/piped` image.
- Manually update your running piped container (not recommended).


