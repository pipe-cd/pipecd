---
title: "Script Run stage"
linkTitle: "Script Run stage"
weight: 4
description: >
   Specific guide for configuring Script Run stage
---

`SCRIPT_RUN` stage is one stage in the pipeline and you can execute any commands.

> Note: This feature is at the alpha status. Currently you can use it on all application kinds, but the rollback feature is only for the application kind of KubernetesApp.

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
  labels:
    env: example
    team: product
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

## Builtin command

Currently, you can use the commands which are installed in the environment for the piped.

For example, If you are using the container platform and the offcial piped container image, you can use the command below.
- git
- ssh
- jq
- curl
- and the builtin commands installed in the base image.

The public piped image available in PipeCD main repo (ref: [Dockerfile](https://github.com/pipe-cd/pipecd/blob/master/cmd/piped/Dockerfile)) is based on [alpine](https://hub.docker.com/_/alpine/) and only has a few UNIX command available (ref: [piped-base Dockerfile](https://github.com/pipe-cd/pipecd/blob/master/tool/piped-base/Dockerfile)). 

If you want to use your commands, you can realize it with either step below.
- Prepare your own environment container image then add [piped binary](https://github.com/pipe-cd/pipecd/releases) to it.
- Build your own container image based on `ghcr.io/pipe-cd/piped` image.

## Default environment values

You can use the envrionment values related to the deployment.

| Name | Description | Example |
|-|-|-|
|SR_DEPLOYMENT_ID| The deployment id | 877625fc-196a-40f9-b6a9-99decd5494a0 |
|SR_APPLICATION_ID| The application id | 8d7609e0-9ff6-4dc7-a5ac-39660768606a |
|SR_APPLICATION_NAME| The application name | example |
|SR_TRIGGERED_AT| The timestamp when the deployment is triggered  | 1719571113 |
|SR_TRIGGERED_COMMIT_HASH| The commit hash that triggered the deployment | 2bf969a3dad043aaf8ae6419943255e49377da0d |
|SR_REPOSITORY_URL| The repository url configured in the piped config  | git@github.com:org/repo.git, https://github.com/org/repo |
|SR_SUMMARY| The summary of the deployment | Sync with the specified pipeline because piped received a command from user via web console or pipectl|
|SR_CONTEXT_RAW| The json encoded string of above values | {"deploymentID":"877625fc-196a-40f9-b6a9-99decd5494a0","applicationID":"8d7609e0-9ff6-4dc7-a5ac-39660768606a","applicationName":"example","triggeredAt":1719571113,"triggeredCommitHash":"2bf969a3dad043aaf8ae6419943255e49377da0d","repositoryURL":"git@github.com:org/repo.git","labels":{"env":"example","team":"product"}} |
|SR_LABELS_XXX| The label attached to the deployment. The env name depends on the label name. For example, if a deployment has the labels `env:prd` and `team:server`, `SR_LABELS_ENV` and `SR_LABELS_TEAM` are registered.  | prd, server |

### Use `SR_CONTEXT_RAW` with jq

You can use jq command to refer to the values from `SR_CONTEXT_RAW`.

```
      - name: SCRIPT_RUN
        with:
          run: |
            echo "Get deploymentID from SR_CONTEXT_RAW"
            echo $SR_CONTEXT_RAW | jq -r '.deploymentID'
            sleep 10
          onRollback: |
            echo "rollback script-run"
```

## Rollback

> Note: Currently, this feature is only for the application kind of KubernetesApp.

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
