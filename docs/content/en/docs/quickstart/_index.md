---
title: "Quickstart"
linkTitle: "Quickstart"
weight: 3
description: >
  This page describes how to quickly get started with PipeCD on Kubernetes.
---

This guides you to install PipeCD in your kubernetes and deploy a `helloworld` application to that Kubernetes cluster.

### Prerequisites
- Have a cluster running a compatible version of Kubernetes
- Installed [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- Installed [helm3](https://helm.sh/docs/intro/install/)
- Forked the [Examples](https://github.com/pipe-cd/examples) repository

### 1. Cloning pipe repository

Navigate to the root of the repository once cloned.

```bash
git clone https://github.com/pipe-cd/manifests.git
cd pipe
```

### 2. Installing control plane

```bash
helm install pipecd ./manifests/pipecd --values ./quickstart/values.yaml
```

### 3. Accessing the PipeCD web
PipeCD comes with an embedded web based UI.
First up, you connect to the control-plane you just installed with kubectl port-forwarding:

```bash
kubectl port-forward svc/pipecd 8080:443
```

Point your web browser to [http://localhost:8080](http://localhost:8080) to login with static admin user.

![](/images/quickstart-login.png)

Enter the project name, username and password. Be sure to give the following:
- Project Name: `quickstart`
- Username: `hello-pipecd`
- Password: `hello-pipecd`

### 4. Adding an environment
To add a new [Environment](http://localhost:1313/docs/concepts/#environment), go to the `Environment` tab at `Settings` page and click on the `Add` button to add a new environment to the project.

Then you give the environment name and its description as shown below:

![](/images/quickstart-adding-environment.png)


### 5. Installing a `piped`
First up, you need to take a piped-id and a secret-key on the Web.

Navigate to the `Piped` tab on the same page as before, click on the `Add` button. Then you enter as:

![](/images/quickstart-adding-piped.png)

Click on the `Save` button, and then you can see the piped-id and secret-key.
Be sure to keep a copy for later use.

![](/images/quickstart-piped-registered.png)



Open [`./quickstart/piped-config.yaml`](https://github.com/pipe-cd/pipe/blob/master/quickstart/piped-config.yaml) with your editor and:
- replace `YOUR_USERNAME` with your git username
- replace `YOUR_PIPED_ID` with the piped-id you have copied before

You can complete the installation by running the following after replacing `YOUR_PIPED_SECRET_KEY` with what you just got:

```bash
helm install dev-piped ./manifests/piped \
  --set args.insecure=true \
  --set secret.pipedKey.data=YOUR_PIPED_SECRET_KEY \
  --set-file config.data=./quickstart/piped-config.yaml
```

### 6. Configuring a kubernetes application
Navigate to the `Application` page, click on the `Add` button. Then give as:

![](/images/quickstart-adding-application.png)

While you can see the select box for the deployment config template, skip it at this point.

After a bit, the first deployment would be complete automatically.

![](/images/quickstart-first-deployment.png)

### 7. Let's deploy!
Let's get started with deployment! All we have to do is to update the image tag.

Open the `canary/deployment.yaml` under the forked examples repository, then change the tag from `v0.1.0` to `v0.2.0`.

![](/images/quickstart-update-image-tag.png)

After a short wait, a new deployment will be started to update to `v0.2.0`.

![](/images/quickstart-deploying.png)

### 8. Cleanup
When you’re finished experimenting with PipeCD, you can uninstall with:

```bash
helm uninstall piped
helm uninstall pipecd
```
