---
title: "Deployment chain"
linkTitle: "Deployment chain"
weight: 11
description: >
  Specific guide for configuring chain of deployments.
---

For users who want to use PipeCD to build a complex deployment flow, which contains multiple applications across multiple application kinds and roll out them to multiple clusters gradually or promoting across environments, this guideline will show you how to use PipeCD to achieve that requirement.

## Configuration

The idea of this feature is to trigger the whole deployment chain when a specified deployment is triggered. To enable trigger the deployment chain, we need to add a configuration section named `postSync` which contains all configurations that be used when the deployment is triggered. For this `Deployment Chain` feature, configuration for it is under `postSync.chain` section.

A canonical configuration looks as below:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: TerraformApp
spec:
  input:
    ...
  pipeline:
    ...
  postSync:
    chain:
      applications:
        # Find all applications with name `application-2` and trigger them.
        - name: application-2
        # Fill all applications with name `application-3` of kind `KUBERNETES`
        # and trigger them.
        - name: application-3
          kind: KUBERNETES
```

As a result, the above configuration will be used to create a deployment chain like the below figure

![](/images/deployment-chain-figure.png)

In the context of the deployment chain in PipeCD, a chain is made up of many `blocks`, and each block contains multiple `nodes` which is the reference to a deployment. The first block in the chain always contains only one node, which is the deployment that triggers the whole chain. Other blocks of the chain are built using filters which are configurable via `postSync.chain.applications` section. As for the above example, the second block `Block 2` contains 2 different nodes, which are 2 different PipeCD applications with the same name `application-2`.

__Tip__:

1. If you followed all the configuration references and built your deployment chain configuration, but some deployments in your defined chain are not triggered as you want, please re-check those deployments [`trigger configuration`](../triggering-a-deployment/#trigger-configuration). The `onChain` trigger is __disabled by default__; you need to enable that configuration to enable your deployment to be triggered as a node in the deployment chain.
2. Values configured under `postSync.chain.applications` - we call it __Application matcher__'s values are merged using `AND` operator. Currently, only `name` and `kind` are supported, but `labels` will also be supported soon.

See [Examples](../../examples/#deployment-chain) for more specific.

## Deployment chain characteristic

Something you need to care about while creating your deployment chain with PipeCD

1. The deployment chain blocks are run in sequence, one by one. But all nodes in the same block are run in parallel, you should ensure that all nodes(deployments) in the same block do not depend on each other.
2. Once a node in a block has finished with `FAILURE` or `CANCELLED` status, the containing block will be set to fail, and all other nodes which have not yet finished will be set to `CANCELLED` status (those nodes will be rolled back if they're in the middle of its deploying process). Consequently, all blocks after that failed block will be set to `CANCELLED` status and be stopped.

## Console view

![](/images/deployment-chain-console.png)

The UI for this deployment chain feature currently is under deployment, we can only __view deployments in chain one by one__ on the deployments page and deployment detail page as usual.

## Reference

See [Configuration Reference](../../configuration-reference/#postsync) for the full configuration.
