---
title: "Deployment Chain"
linkTitle: "Deployment Chain"
weight: 99
description: >
  Specific guide for configuring chain of deployments.
---

For users who want to use PipeCD to build a complex deployment flow, which contains multiple applications across multiple application kinds, this guideline will show you how to use PipeCD to archive that requirement.

## Configuration

The idea of this feature is to trigger the whole deployment chain when a specified deployment is triggered. To enable trigger the deployment chain, we need to add a configuration section named `postSync` which contains all configurations that be used when the deployment is triggered. For this `Deployment Chain` feature, configuration for it is under `postSync.chain` section.

A canonical configuration look as below:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Kubernetes
spec:
  input:
    ...
  pipeline:
    ...
  postSync:
    chain:
      applications:
        - name: Application 2
        - name: Application 3
```

As a result, the above configuration will be used to create a deployment chain like the below figure

![](/images/deployment-chain-figure.png)

In the context of the deployment chain in PipeCD, a chain is made up of many `blocks`, and each block contains multiple `nodes` which is the reference to a deployment. The first block in the chain always contains only one node, which is the deployment that triggers the whole chain. Other blocks of the chain are built using filters which are configurable via `postSync.chain.applications` field. As for the above example, the second block `Block 2` contains 3 different nodes, which are 3 different PipeCD applications with the same name `Application 2`.

See [Examples](/docs/examples/#deployment-chain) for more specific.

## Deployment chain characteristic

Something you need to care about while creating your deployment chain with PipeCD

1. The deployment chain block is run in sequence, one by one. But all nodes in the same block are run in parallel, you should ensure that all deployments in the same block do not depend on each other.
2. Once a block in the chain is finished with `FAILURE` or `CANCELLED` status, the chain will be set to fail, and all blocks after that finished block will be set to `CANCELLED` status.
3. Once a node in a block is finished with `FAILURE` or `CANCELLED` status, the block will be set to fail, and all other nodes that are not yet finished will be set to `CANCELLED` status (those nodes will be rolled back if they're in the middle of deploying process).

## Console view

TBD

## Reference

See [Configuration Reference](/docs/user-guide/configuration-reference/) for the full configuration.
