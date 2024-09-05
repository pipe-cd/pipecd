---
title: "Script Run stage"
linkTitle: "Script Run stage"
weight: 4
description: >
   Specific guide for configuring Script Run stage
---

`SCRIPT_RUN` stage is one stage in the pipeline and you can execute any commands.

> Note: This feature is at the alpha status and currently only for the application kind of KubernetesApp.

## How to configure SCRIPT_RUN stage

Add a `SCRIPT_RUN` to your pipeline and write commands. 

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

You can define the command as `run`.
Also, if you want to some values as variables, you can define them as `env`.

The commands run in the directory where this application configuration file exists.

![](/images/script-run.png)

### Execute the script file

If your script is so long, you can separate the script as a file.
You can put the file with the app.pipecd.yaml in the same dir and then you can execute the script like this.

```
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: script-run
  pipeline:
    stages:
      - name: SCRIPT_RUN
        with:
          run: |
            sh script.sh
```

```
.
├── app.pipecd.yaml
└── script.sh
```

## Builtin commands

Currently, you can use the commands which are installed in the environment for the piped.

For example, If you are using the container platform and the offcial piped container image, you can use the command below.
- git
- ssh
- curl
- commands installed by piped in $PIPED_TOOL_DIR (check at runtime)
- built-in commands installed in the base image

The public piped image available in PipeCD main repo (ref: [Dockerfile](https://github.com/pipe-cd/pipecd/blob/master/cmd/piped/Dockerfile)) is based on [alpine](https://hub.docker.com/_/alpine/) and only has a few UNIX commands available (ref: [piped-base Dockerfile](https://github.com/pipe-cd/pipecd/blob/master/tool/piped-base/Dockerfile)). 

If you want to use your commands, you can realize it with either step below.
- Prepare your own environment container image then add [piped binary](https://github.com/pipe-cd/pipecd/releases) to it.
- Build your own container image based on `ghcr.io/pipe-cd/piped` image.

# When to rollback

You can define the command as `onRollback` to execute when to rollback similar to `run`.
Execute the command to rollback SCRIPT_RUN to the point where the deployment was canceled or failed.

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
            R_MSG: "rollback script1"
          run: |
            echo $MSG
            sleep 10
          onRollback: |
            echo $R_MSG
            sleep 10
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_CANARY_CLEAN
```

![](/images/script-run-onRollback.png)

The command defined as `onRollback` is executed as `SCRIPT_RUN_ROLLBACK` stage after each `ROLLBACK` stage.

When there are multiple SCRIPT_RUN stages, they are executed in the same order as SCRIPT_RUN on the pipeline.
Also, only for the executed SCRIPT_RUNs are rollbacked.

For example, consider when deployment proceeds in the following order from 1 to 7.
```
1. K8S_CANARY_ROLLOUT
2. WAIT
3. SCRIPT_RUN
4. K8S_PRIMARY_ROLLOUT
5. SCRIPT_RUN
6. K8S_CANARY_CLEAN
7. SCRIPT_RUN
```

Then
- If 3 is canceled or fails while running, only SCRIPT_RUN of 3 will be rollbacked.
- If 4 is canceled or fails while running, only SCRIPT_RUN of 3 will be rollbacked.
- If 6 is canceled or fails while running, only SCRIPT_RUNs 3 and 5 will be rollbacked. The order of executing is 3 -> 5.
