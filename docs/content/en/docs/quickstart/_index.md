---
title: "Quickstart"
linkTitle: "Quickstart"
weight: 3
description: >
  This page describes how to quickly get started with PipeCD on Kubernetes.
---

This guides you to install PipeCD in your kubernetes and deploy a `helloworld` application to that Kubernetes cluster.

### Requirements
- Have a cluster running a compatible version of Kubernetes
- Installed [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- Installed [helm3](https://helm.sh/docs/intro/install/)


### 1. Cloning pipe repository

Navigate to the root of the repository once cloned.

```bash
git clone https://github.com/pipe-cd/pipe.git
cd pipe
```

### 2. Installing control plane
First up, add pipecd helm chart repository.
```bash
helm repo add pipecd https://pipecd-charts.storage.googleapis.com
```

Then you can install with:

```bash
helm install pipecd pipecd/pipecd --values ./quickstart/values.yaml
```

### 3. Accessing the API server
Connect to the API server you just installed with kubectl port-forwarding:

```bash
kubectl port-forward svc/pipecd 8080:443
```

### 4. Logging in with static admin user
PipeCD comes with an embedded web based UI.
Point your web browser to [http://localhost:8080](http://localhost:8080) ensure the API server has started successfully.

![](/images/login.png)

Enter the project name, username and password. Be sure to give the following:
- Project Name: `quickstart`
- Username: `hello-pipecd`
- Password: `hello-pipecd`

### 5. Adding an environment
To add a new [Environment](http://localhost:1313/docs/concepts/#environment), go to the `Environment` tab at `Settings` page and click on the `Add` button to add a new environment to the project.

Then you give the environment name and its description as shown below:

![](/images/adding-environment.png)


### 6. Registering a `piped`
Before installing the `piped`, you have to register it on the Web.
Navigate to the `Piped` tab on the same page as before, and then click on the `Add` button.

![](/images/adding-piped.png)

Click on the `Save` button, and then you can see the piped-id and secret key.
Be sure to keep a copy for later use.

![](/images/piped-registered.png)

### 7. Installing a `piped`

Fork the [Examples](https://github.com/pipe-cd/examples) repository.

Once a repository named `YOUR_USERNAME/examples` has been created,
open [`./quickstart/piped-config.yaml`](https://github.com/pipe-cd/pipe/blob/master/quickstart/piped-config.yaml) and:
- replace `YOUR_USERNAME` with your git username
- replace `YOUR_PIPED_ID` with the piped-id you have copied before

You can complete the installation by running the following after replacing `YOUR_PIPED_SECRET_KEY` with yours:

```bash
helm install dev-piped pipecd/piped \
  --set args.address=pipecd:443 \
  --set secret.pipedKey.data=YOUR_PIPED_SECRET_KEY \
  --set-file config.data=./quickstart/piped-config.yaml
```

### 8. Configuring a kubernetes application

> TBA

### 9. Let's deploy!

> TBA


### 10. Cleanup
When you’re finished experimenting with PipeCD, you can uninstall with:

```bash
helm uninstall piped
helm uninstall pipecd
```
