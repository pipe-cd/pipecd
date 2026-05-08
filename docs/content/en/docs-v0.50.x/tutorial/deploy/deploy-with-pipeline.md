---
title: "2. Deploy with a Customized Pipeline"
linkTitle: "2. Deploy with Pipeline"
weight: 2
description: >
  Customize the deployment pipeline with strategies like canary or blue/green.
---

# 2. Deploy with a Customized Pipeline

In this page, you will customize the deployment pipeline by editing the `app.pipecd.yaml` deployment configuration.

## Overview

In the previous step, you deployed an application using the default pipeline (Quick Sync).
PipeCD also supports customized pipelines with stages such as:

- **K8S_CANARY_ROLLOUT** / **ECS_CANARY_ROLLOUT**: Deploy a canary version alongside the stable version.
- **K8S_TRAFFIC_ROUTING** / **ECS_TRAFFIC_ROUTING**: Shift traffic between the stable and canary versions.
- **WAIT_APPROVAL**: Pause the pipeline until a manual approval is given on the console.
- **ANALYSIS**: Automatically analyze metrics to determine the deployment's health.
- **K8S_PRIMARY_ROLLOUT** / **ECS_PRIMARY_ROLLOUT**: Promote the canary version to primary.
- **K8S_CANARY_CLEAN** / **ECS_CANARY_CLEAN**: Clean up the canary resources.

## Steps

### 1. Edit `app.pipecd.yaml`

Edit the `app.pipecd.yaml` in the directory of the platform you deployed in the previous step.

For example, for **Kubernetes**, edit `src/deploy/kubernetes/simple/app.pipecd.yaml`:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: tutorial-kubernetes-simple
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 50%
      - name: WAIT_APPROVAL
        with:
          approvers:
            - hello-pipecd
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_CANARY_CLEAN
```

For **Amazon ECS**, edit `src/deploy/ecs/simple/app.pipecd.yaml`:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: ECSApp
spec:
  name: tutorial-ecs-simple
  pipeline:
    stages:
      - name: ECS_CANARY_ROLLOUT
        with:
          scale: 50%
      - name: WAIT_APPROVAL
        with:
          approvers:
            - hello-pipecd
      - name: ECS_PRIMARY_ROLLOUT
      - name: ECS_CANARY_CLEAN
```

### 2. Edit a manifest

Make some change to the manifest file as well (e.g., update the image tag or replica count) so that Piped detects a new change to deploy.

### 3. Push and watch

1. Commit and push the changes to remote.

   ```sh
   git add .
   git commit -m "Add pipeline stages"
   git push origin main
   ```

2. Go to the deployments page. [http://localhost:8080/deployments](http://localhost:8080/deployments)

3. A new deployment will start in a few minutes. You will see the pipeline stages in the deployment details.

4. When the pipeline reaches the `WAIT_APPROVAL` stage, click `APPROVE` on the console to continue.

5. After approval, the remaining stages will execute automatically.

## See Also

- [Configuring Deployment Pipeline](https://pipecd.dev/docs/user-guide/managing-application/customizing-deployment/)
- [Pipeline Stages Reference](https://pipecd.dev/docs/user-guide/managing-application/customizing-deployment/)

---

[Next: Congratulations! >](../../tutorial/next-step/)

[< Previous: Deploy Simply](../deploy-simply/)
