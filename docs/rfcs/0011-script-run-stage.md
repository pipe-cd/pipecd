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

```yaml
apiVersion: pipecd.dev/v1beta1
kind: LambdaApp
spec:
  encryptedSecrets:
    password: encrypted-secrets
  pipeline:
    stages:
      - name: SCRIPT_RUN
        with:
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

```yaml
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
        with:
            env:
              SLACK_WEBHOOK_URL: ""
            runs:
              - "curl -X POST -H 'Content-type: application/json' --data '{"text":"successfully deployed!!"}' $SLACK_WEBHOOK_URL"
```

## when to rollback

Users can define commands to execute with `onRollback` when rolling back.
If `onRollback` is not set, nothing to execute when rolling back.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  pipeline:
    stages:
      - name: SCRIPT_RUN
        with:
          env:
            SLACK_WEBHOOK_URL: ""
          runs:
            - "curl -X POST -H 'Content-type: application/json' --data '{"text":"successfully deployed!!"}' $SLACK_WEBHOOK_URL"
          onRollback:
            - "curl -X POST -H 'Content-type: application/json' --data '{"text":"failed to deploy: rollback"}' $SLACK_WEBHOOK_URL"
```

**SCRIPT_SYNC stage also rollbacks** when the deployment status is `DeploymentStatus_DEPLOYMENT_CANCELLED` or `DeploymentStatus_DEPLOYMENT_FAILURE` even though other rollback stage is also executed.

For example, here is a deploy pipeline combined with other k8s stages.
The result status of the pipeline is FAIL or CANCELED, piped rollbacks the stages `K8S_CANARY_ROLLOUT`, `K8S_PRIMARY_ROLLOUT`, and `SCRIPT_RUN`.

```yaml
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
        with:
          env:
            SLACK_WEBHOOK_URL: ""
          runs:
            - "curl -X POST -H 'Content-type: application/json' --data '{"text":"successfully deployed!!"}' $SLACK_WEBHOOK_URL"
          onRollback:
            - "curl -X POST -H 'Content-type: application/json' --data '{"text":"failed to deploy: rollback"}' $SLACK_WEBHOOK_URL"
```

## prepare environment for execution

Commands are executed on the container of piped or on the host OS(standalone).

## CUSTOM_SYNC in the future

"CUSTOM_SYNC" stage will be deprecated because the "SCRIPT_RUN" has also similar features.
I expect users to use "SCRIPT_RUN" for executing any command before or after other stages.


# Alternatives

## What's the difference between "CUSTOM_SYNC" and "SCRIPT_SYNC"?

"CUSTOM_SYNC" is one of the stages to **sync**, but "SCRIPT_RUN" is the stage to **execute commands**.

## How about other CD tools?

### Argo

#### strategy

- **execute command with another k8s resources**
#### details
Resource Hooks 
https://argo-cd.readthedocs.io/en/stable/user-guide/resource_hooks/#resource-hooks
> Hooks are ways to run scripts before, during, and after a Sync operation. Hooks can also be run if a Sync operation fails at any point.

- There are four points to execute command.
	- PreSync: before sync
	- Sync: during sync
	- PostSync: after sync
	- SyncFaill: failed to sync
- **To execute command, ArgoCD applys k8s resources such as Job or Pod, [[Argo Workflows]] ...**
- users set some annotations and ArgoCD detects them to control the order to execute command
	- e.g. https://argo-cd.readthedocs.io/en/stable/user-guide/resource_hooks/#using-a-hook-to-send-a-slack-message

#### pros/cons

**pros**
- can separate respolibility for delivery and executing any command

**cons**
- users need to prepare and manage the resource to execute any command

### FluxCD
- There is no functions to realize that

### Flagger

#### strategy
- Call api set as webhooks on each points and execute command on the api
- Flagger just call api registerd as webhooks.

#### details

**Webhooks**
https://fluxcd.io/flagger/usage/webhooks/#load-testing
>Flagger will call each webhook URL and determine from the response status code (HTTP 2xx) if the canary is failing or not.

- There are some webhook points.
	- confirm-rollout
	- pre-rollout
	- rollout
	- confirm-traffic-increase
	- confirm-promotion
	- post-rollout
	- rollback
	- event
- Flagger calls webhooks on each points.
	- e.g. load testing with Flagger: https://fluxcd.io/flagger/usage/webhooks/#load-testing
- if users want to execute command with webhook, set the command text to `metadata` section and prepare webhook handler to execute cmd parsed from metadata.

#### pros/cons

**pros**
- can separate respolibility for delivery and executing any command

**cons**
- users need to prepare api for the webhooks.


# Unresolved questions

- It might be better to change from alpine to debian and so on to provide users popular commands (e.g. curl) by default.
- It might be better to add `runsOnRollback` on each sync stage if users want to control when rollback.
This is a just idea.
