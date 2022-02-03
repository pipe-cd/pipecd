---
title: "Quickstart"
linkTitle: "Quickstart"
weight: 3
description: >
  This page describes how to quickly get started with PipeCD on Kubernetes.
---

This page is a guideline for installing PipeCD in your Kubernetes cluster and deploying a "hello world" application to that same Kubernetes cluster.

### Prerequisites
- Having a Kubernetes cluster
- Installed [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) and [helm3](https://helm.sh/docs/intro/install/)
- Forked the [Examples](https://github.com/pipe-cd/examples) repository

### 1. Installing control plane

``` console
helm repo add pipecd https://charts.pipecd.dev

helm install pipecd pipecd/pipecd -n pipecd \
  --create-namespace \
  --dependency-update \
  --values https://raw.githubusercontent.com/pipe-cd/manifests/{{< blocks/latest_version >}}/quickstart/control-plane-values.yaml
```

Once installed, use `kubectl port-forward` to expose the web console on your localhost:

``` console
kubectl -n pipecd port-forward svc/pipecd 8080
```

The PipeCD web console will be available at [http://localhost:8080](http://localhost:8080). To login, you can use the configured static admin account as below:
- project name: `quickstart`
- username: `hello-pipecd`
- password: `hello-pipecd`

![](/images/quickstart-login.png)

### 2. Installing a `piped`
Before running a piped, you have to register it on the web and take the generated ID and Key strings.

Navigate to the `Piped` tab on the same page as before, click on the `Add` button. Then you enter as:

![](/images/quickstart-adding-piped.png)

Click on the `Save` button, and then you can see the piped-id and secret-key.
Be sure to keep a copy for later use.

![](/images/quickstart-piped-registered.png)

Then complete the installation by running the following command after replacing `{PIPED_ID}`, `{PIPED_KEY}`, `{FORKED_GITHUB_ORG}` with what you just got:

``` console
helm install piped pipecd/piped -n pipecd \
  --set quickstart.enabled=true \
  --set quickstart.pipedId={PIPED_ID} \
  --set secret.data.piped-key={PIPED_KEY} \
  --set quickstart.gitRepoRemote=https://github.com/{FORKED_GITHUB_ORG}/examples.git
```

### 3. Registering a kubernetes application
Navigate to the `Applications` page, click on the `Add` button on the top left corner.

Go to the `ADD FROM GIT` tab, then select:
- Piped: `dev` (you just registered)
- CloudProvider: `kubernetes-default`

You should see a lot of suggested applications.

Select the `canary` application and click the `ADD` button to register.

![](/images/quickstart-adding-application-from-git.png)

After a bit, the first deployment would be complete automatically to sync the application to the state specified in the current Git commit.

![](/images/quickstart-first-deployment.png)

### 4. Let's deploy!
Let's get started with deployment! All you have to do is to make a PR to update the image tag, scale the replicas, or change the manifests.

For instance, open the `kubernetes/canary/deployment.yaml` under the forked examples' repository, then change the tag from `v0.1.0` to `v0.2.0`.

![](/images/quickstart-update-image-tag.png)

After a short wait, a new deployment will be started to update to `v0.2.0`.

![](/images/quickstart-deploying.png)

### 5. Cleanup
When you’re finished experimenting with PipeCD, you can uninstall with:

``` console
helm -n pipecd uninstall piped
helm -n pipecd uninstall pipecd
kubectl delete deploy canary -n pipecd
kubectl delete svc canary -n pipecd
```

### What's next?
You want to know some details on how to set up for a production environment? Visit [Operating Control Plane](/docs/operator-manual/control-plane/) at first. After reading that, the [Operating Piped](/docs/operator-manual/piped/) page will be for you.
