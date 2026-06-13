---
title: "1. Deploy Simply"
linkTitle: "1. Deploy Simply"
weight: 1
description: >
  Deploy an application to your platform in a simple way.
---

# 1. Deploy Simply

In this page, you will deploy an application to your platform in a simple way.

## 1. Prepare config files

1-1. Edit configuration files of one platform under `/src/deploy/` in your cloned repository as bellow.

> **Note:** Each platform has a deployment config file (`app.pipecd.yaml`) + platform specific files.

- For **Kubernetes**:
  - You will run a [helloworld](https://github.com/pipe-cd/pipecd/pkgs/container/helloworld) service.
  - Use `kubernetes/simple/`. You do not need to edit.
- For **Google Cloud Run**:
  - You will run a [helloworld](https://github.com/pipe-cd/pipecd/pkgs/container/helloworld) service.
  - Use `cloudrun/simple/`. You do not need to edit.
- For **Amazon ECS**:
  - You will run an nginx service.
  - Edit `ecs/simple/` as below.
    - `app.pipecd.yaml`: Edit `targetGroupArn`.
    - `servicedef.yaml`: Edit `cluster`, `securityGroups`, and `subnets`.
    - `taskdef.yaml`: Edit `executionRoleArn`.
- For **AWS Lambda**:
  - You will create a function of your own image.
  - Edit `lambda/simple/` as below.
    - `function.yaml`: Edit `role` and `image`.
- For **Terraform**:
  - You will generate a file on local.
  - Edit `terraform/simple/` as below.
    - `main.tf`: Edit `path` and `filename`.

1-2. Commit and push the changes to remote.


## 2. Register the application

2-1. Go to the applications page. [http://localhost:8080](http://localhost:8080)

2-2. Click `+ ADD`.

2-3. Enter values and click `SAVE`->`SAVE`.
   - `Piped`: Your Piped
   - `Platform Provider`: The platform
   - `Application`: The application you configured in [1.](#1-prepare-config-files)

![add-application-input](/images/tutorial/deploy/application-add-input.png)

2-4. If successful, you will see the dialog like the following image.

![application-is-added](/images/tutorial/deploy/application-is-added.png)


## 3. Watch the first deployment

3-1. Go to the deployments page. [http://localhost:8080/deployments](http://localhost:8080/deployments)

3-2. Wait until a new deployment automatically appears. Then click it to see details.

> **Note:** You do NOT need to invoke a deployment by yourself since your Piped automatically starts it. **This is GitOps.**

![deployment-before-appear](/images/tutorial/deploy/deployment-before-appear.png)
![deployment-appear](/images/tutorial/deploy/deployment-appear.png)


3-3. Wait until the status becomes `SUCCESS`/`FAILURE`. If it shows `FAILURE`, check the log in the page.

![deployment-deploying](/images/tutorial/deploy/deployment-deploying.png)


3-4. When the status becomes `SUCCESS`, the deployment is successfully finished.

![deployment-success](/images/tutorial/deploy/deployment-success.png)

3-5. Check your platform (Kubernetes cluster, cloud console, etc.) to confirm the result.
     

## 4. Edit the config and watch a new deployment

4-1. Edit the config file you deployed. (e.g. Change the image tag, sizing, etc.)

4-2. Commit and push the change to remote.

4-3. Go to the deployments page again. [http://localhost:8080/deployments](http://localhost:8080/deployments)

4-4. A new deployment will start in a few minutes.

---

[Next: Deploy with a Customized Pipeline >](../deploy-with-pipeline/)

[< Previous: Deploy](../)
