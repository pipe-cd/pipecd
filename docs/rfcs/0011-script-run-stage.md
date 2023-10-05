- Start Date: (fill me in with today's date, YYYY-MM-DD)
- Target Version: (1.x / 2.x)

# Summary

This RFC introduces a new way to enable users to use "script run stages” that users can execute any commands in their pipelines.

# Motivation

Currently, users can use only stages that PipeCD has already defined. However some users want to define new stages by their use-cases as below. 

- Deploying infrastructure by tools other than that PipeCD supports (terraform and kubernetes) such as SAM, cloud formation….
- Running End to End tests
- Interacting with external systems
- Performing database migrations
- notifying the deployed result

`CUSTOM_SYNC` is implemented for the above use-cases, but it is for sync. 
So more simply, some users want to execute commands.

# Detailed design

## feature

1. execute any commands in their pipeline.

```
apiVersion: pipecd.dev/v1beta1
kind: LambdaApp
spec:
  encryptedSecrets:
    password: encrypted-secrets
  pipeline:
    stages:
      - name: SCRIPT_RUN
        env:
          AWS_PROFILE: default
        runs:
          - "echo {{ .encryptedSecrets.password }} | sudo -S su"
          - "sam build"
          - "sam deploy -g --profile $AWS_PROFILE"
```

2. combine with other stage

Compared to CUSTOM_SYNC, this stage can be combined with other stage.
For example, 

```
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 10%
      - name: WAIT_APPROVAL
        with:
          timeout: 30m
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_CANARY_CLEAN
      - name: SCRIPT_RUN
        env:
          SLACK_WEBHOOK_URL: ""
        runs:
          - "curl -X POST -H 'Content-type: application/json' --data '{"text":"successfully deployed!!"}' $SLACK_WEBHOOK_URL"
```

3. execute any commands when rollback

Users can define commands to execute when rolling back with `runsOnRollback`.
If `runsOnRollback` is not set, nothing to execute when rolling back.

```
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 10%
      - name: WAIT_APPROVAL
        with:
          timeout: 30m
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_CANARY_CLEAN
      - name: SCRIPT_RUN
        env:
          SLACK_WEBHOOK_URL: ""
        runs:
          - "curl -X POST -H 'Content-type: application/json' --data '{"text":"successfully deployed!!"}' $SLACK_WEBHOOK_URL"
        runsOnRollback:
          - "curl -X POST -H 'Content-type: application/json' --data '{"text":"failed to deploy: rollback"}' $SLACK_WEBHOOK_URL"
```

# Alternatives

**What's the difference between "CUSTOM_SYNC" and "SCRIPT_SYNC"?**

"CUSTOM_SYNC" is one of the stages to **sync**, but "SCRIPT_RUN" is the stage to **execute commands**.

# Unresolved questions

It might be better to add `runsOnRollback` on each sync stage if users want to control when rollback.
This is a just idea.
