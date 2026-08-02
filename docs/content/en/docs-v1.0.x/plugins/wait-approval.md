---
title: "Wait approval plugin"
linkTitle: "Wait approval"
weight: 50
description: >
  Pause the pipeline until a user approves.
---

The `waitapproval` plugin provides the `WAIT_APPROVAL` stage, which pauses the pipeline until the required number of users approve the deployment from the PipeCD console. Because it is a stage plugin, `WAIT_APPROVAL` can be added to any deployment pipeline.

## Prerequisites

Register the plugin in the piped configuration:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  # ...
  plugins:
    - name: waitapproval
      port: 7005
      url: file:///path/to/plugin/binary  # or an https:// release URL
```

## The WAIT_APPROVAL stage

Add a `WAIT_APPROVAL` stage and set the approvers and the number of approvals required under `with`. For example, require one approval before applying an infrastructure change:

```yaml
pipeline:
  stages:
    - name: TERRAFORM_PLAN
    - name: WAIT_APPROVAL
      with:
        approvers:
          - user-a
          - user-b
        minApproverNum: 1
    - name: TERRAFORM_APPLY
```

The stage waits until at least `minApproverNum` users have approved it from the console, then continues.

## Configuration reference

### WAIT_APPROVAL stage options

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| approvers | []string | Users designated as approvers of the deployment. | Yes |
| minApproverNum | int | Number of approvals required before the pipeline continues. | No (default `1`) |
