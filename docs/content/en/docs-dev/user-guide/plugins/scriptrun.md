---
title: "Script Run"
linkTitle: "Script Run"
weight: 50
description: >
  Execute custom scripts during deployment
---

The Script Run plugin enables PipeCD to execute custom scripts as part of deployment pipelines, allowing for flexible deployment logic and integration with external systems.

## Features

- **Custom Scripts:** Execute arbitrary scripts during deployment
- **Environment Variables:** Access deployment context via environment variables
- **Rollback Scripts:** Define custom rollback logic
- **Shell Support:** Execute bash, sh, and other shell scripts
- **Context Information:** Access deployment metadata in scripts
- **Error Handling:** Proper error handling and logging

## Piped Configuration

Configure the Script Run plugin in your Piped configuration:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: dev
  pipedID: xxx
  plugins:
    - name: scriptrun
      port: 7005
      url: https://github.com/pipe-cd/pipecd/releases/download/...
```

## Application Configuration

Define script execution in `.pipe.yaml`:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-app
  labels:
    env: production
  pipeline:
    stages:
      - name: SCRIPT_RUN
        with:
          run: ./deploy.sh
  plugins:
    scriptrun: {}
```

## Available Stages

- **SCRIPT_RUN:** Execute custom deployment script
- **SCRIPT_RUN_ROLLBACK:** Execute custom rollback script

## Stage Configuration

### SCRIPT_RUN

Execute a deployment script.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: deployment-app
  labels:
    env: production
  pipeline:
    stages:
      - name: SCRIPT_RUN
        with:
          run: ./scripts/deploy.sh
          env:
            ENVIRONMENT: production
            REGION: us-east-1
  plugins:
    scriptrun: {}
```

### SCRIPT_RUN_ROLLBACK

Execute a rollback script on deployment failure.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: deployment-app
  labels:
    env: production
  pipeline:
    stages:
      - name: SCRIPT_RUN
        with:
          run: ./deploy.sh
      - name: SCRIPT_RUN_ROLLBACK
        with:
          onRollback: ./scripts/rollback.sh
  plugins:
    scriptrun: {}
```

## Examples

### Simple Script Execution

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: deployment-app
  labels:
    env: production
  pipeline:
    stages:
      - name: SCRIPT_RUN
        with:
          run: ./deploy.sh
  plugins:
    scriptrun: {}
```

### Script with Environment Variables

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: complex-deployment
  labels:
    env: production
  pipeline:
    stages:
      - name: SCRIPT_RUN
        with:
          run: ./scripts/deploy.sh
          env:
            ENVIRONMENT: production
            REGION: us-central1
            VERSION: "1.0.0"
          onRollback: ./scripts/rollback.sh
  plugins:
    scriptrun: {}
```

## Source Code

- [Script Run Plugin](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin/scriptrun)

## See Also

- [Configuration Reference](/docs-dev/user-guide/configuration-reference/)
