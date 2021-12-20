---
title: "Quickstart"
linkTitle: "Quickstart"
weight: 3
description: >
  This page describes how to quickly get started with PipeCD on Kubernetes.
---

This guides you to install PipeCD in your kubernetes and deploy a `helloworld` application to that Kubernetes cluster. For further reading about PipeCD's ideas, please visit [overview](/docs/overview/) and [concepts](/docs/concepts/).

### Prerequisites
- Having a Kubernetes cluster
- Installed [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- Installed [helm3](https://helm.sh/docs/intro/install/)
- Forked the [Examples](https://github.com/pipe-cd/examples) repository

### 1. Installing control plane

``` console
helm repo add pipecd https://charts.pipecd.dev

helm install pipecd pipecd/pipecd -n pipecd --dependency-update --create-namespace \
  --values https://raw.githubusercontent.com/pipe-cd/manifests/{{< blocks/latest_version >}}/quickstart/control-plane-values.yaml
```

### 2. Accessing the PipeCD web
PipeCD comes with an embedded web-based UI.
First up, using kubectl port-forward to expose the installed control-plane on your localhost:

``` console
kubectl -n pipecd port-forward svc/pipecd 8080
```

Point your web browser to [http://localhost:8080](http://localhost:8080) to login with the configured static admin account.

![](/images/quickstart-login.png)

Enter the project name, username and password. Be sure to give the following:
- Project Name: `quickstart`
- Username: `hello-pipecd`
- Password: `hello-pipecd`

### 3. Installing a `piped`
Before running a piped, you have to register it on the web and take the generated ID and Key strings.

Navigate to the `Piped` tab on the same page as before, click on the `Add` button. Then you enter as:

![](/images/quickstart-adding-piped.png)

Click on the `Save` button, and then you can see the piped-id and secret-key.
Be sure to keep a copy for later use.

![](/images/quickstart-piped-registered.png)


Open [`./quickstart/piped-values.yaml`](https://github.com/pipe-cd/manifests/blob/master/quickstart/piped-values.yaml) with your editor and:
- replace `FORKED_REPO_URL` with forked repository of [Examples](https://github.com/pipe-cd/examples), such as `https://github.com/YOUR_ORG/examples.git`
- replace `YOUR_PIPED_ID` with the piped-id you have copied before

You can complete the installation by running the following after replacing `{YOUR_PIPED_SECRET_KEY}` with what you just got:

``` console
helm install piped pipecd/piped -n pipecd \
  --values https://raw.githubusercontent.com/pipe-cd/manifests/{{< blocks/latest_version >}}/quickstart/piped-values.yaml \
  --set secret.pipedKey.data={YOUR_PIPED_SECRET_KEY}
```

### 4. Registering a kubernetes application
Navigate to the `Application` page, click on the `Add` button on the top left corner.

Go to the `ADD FROM GIT` tab, then select:
- Piped: `dev` (you just registered)
- CloudProvider: `kubernetes-default`

You should see a lot of application suggestions.

Click the `canary` row. Make sure Kind is set to KUBENETES, and then click the ADD button.

![](/images/quickstart-adding-application.png)

After a bit, the first deployment would be complete automatically to sync the application to the state specified in the current Git commit.

![](/images/quickstart-first-deployment.png)

### 5. Let's deploy!
Let's get started with deployment! All you have to do is to make a PR to update the image tag, scale the replicas, or change the manifests.

For instance, open the `kubernetes/canary/deployment.yaml` under the forked examples' repository, then change the tag from `v0.1.0` to `v0.2.0`.

![](/images/quickstart-update-image-tag.png)

After a short wait, a new deployment will be started to update to `v0.2.0`.

![](/images/quickstart-deploying.png)

### 6. Cleanup
When you’re finished experimenting with PipeCD, you can uninstall with:

``` console
helm -n pipecd uninstall piped
helm -n pipecd uninstall pipecd
kubectl delete deploy canary -n pipecd
kubectl delete svc canary -n pipecd
```

### What's next?
You want to know some details on how to set up for a production environment? Visit [Operating Control Plane](/docs/operator-manual/control-plane/) at first. After reading that, the [Operating Piped](/docs/operator-manual/piped/) page will be for you.
