# WAIT stage plugin

## Overview

`WAIT` stage is a stage that waits for the specified duration and then proceeds to the next stage.

cf. The spec is almost the same as pipedv0:
https://pipecd.dev/docs-v0.52.x/user-guide/managing-application/customizing-deployment/adding-a-wait-stage/

## Plugin Configuration

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  plugins:
  - name: wait
    port: 7002 # any unused port
    url: file:///path/to/.piped/plugins/wait # or remoteUrl(TBD)
```

`config` and `deployTargets` are not supported.

## Application Configuration

### WAIT stage options

| Field | Type | Description | Required | Default |
|-|-|-|-|-|
| duration | duration | The duration to wait. e.g. 30s | Yes | |
