---
title: "Script run plugin"
linkTitle: "Script run"
weight: 70
description: >
  Run arbitrary commands as a pipeline stage.
---

The `scriptrun` plugin provides the `SCRIPT_RUN` stage, which runs arbitrary shell commands as a step in the pipeline. Because it is a stage plugin, `SCRIPT_RUN` can be added to any deployment pipeline. It is useful for tasks such as smoke tests, notifications, or custom checks between deployment stages.

## Prerequisites

Register the plugin in the piped configuration:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  # ...
  plugins:
    - name: scriptrun
      port: 7006
      url: file:///path/to/plugin/binary  # or an https:// release URL
```

## The SCRIPT_RUN stage

Add a `SCRIPT_RUN` stage and set the command to run under `with.run`. You can pass environment variables with `env`, and provide an `onRollback` command that runs if the deployment is rolled back:

```yaml
pipeline:
  stages:
    - name: K8S_CANARY_ROLLOUT
    - name: SCRIPT_RUN
      with:
        env:
          APP_URL: https://example.com
        run: |
          curl -sSf "$APP_URL/healthz"
        onRollback: |
          echo "rolling back"
    - name: K8S_PRIMARY_ROLLOUT
```

If `onRollback` is set, PipeCD runs it during rollback through an automatically inserted `SCRIPT_RUN_ROLLBACK` stage. You do not add `SCRIPT_RUN_ROLLBACK` to your pipeline.

## Configuration reference

### SCRIPT_RUN stage options

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| run | string | The command(s) to run. | Yes |
| env | map[string]string | Environment variables to set when running the command. | No |
| onRollback | string | The command(s) to run if the deployment is rolled back. | No |
