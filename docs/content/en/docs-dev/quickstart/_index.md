---
title: "Quickstart"
linkTitle: "Quickstart"
weight: 3
description: >
  This page describes how to quickly get started with PipeCD on Kubernetes.
---

This page is a guideline for installing PipeCD into your Kubernetes cluster and deploying a "hello world" application to that same Kubernetes cluster.

Note:

- It's not required to install the PipeCD control plane to the cluster where your applications are running. Please read this [blog post](/blog/2021/12/29/pipecd-best-practice-01-operate-your-own-pipecd-cluster/) to understand more about PipeCD in real life use cases.
- If you want to experiment with PipeCD freely or don't have a Kubernetes cluster, we recommend [this Tutorial](https://github.com/pipe-cd/tutorial).

### Prerequisites
- Having a Kubernetes cluster and connect to it via [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/).

### 1. Installing PipeCD in quickstart mode

Across the [PipeCD concepts](/docs/concepts/), PipeCD platform is constructed by 2 components: Control Plane and Piped (the agent).

#### 1.1. Installing PipeCD Control Plane

```console
$ kubectl create namespace pipecd
$ kubectl apply -n pipecd -f https://raw.githubusercontent.com/pipe-cd/pipecd/master/quickstart/manifests/control-plane.yaml
```

The PipeCD control plane will be installed with a default project named `quickstart`. To access the PipeCD Control Plane UI, run the following command

```console
$ kubectl port-forward -n pipecd svc/pipecd 8080
```

You can access to the PipeCD console at [http://localhost:8080?project=quickstart](http://localhost:8080?project=quickstart)

To login, you can use the configured static admin account as below:
- username: `hello-pipecd`
- password: `hello-pipecd`

And you will access the main page of PipeCD Control Plane console, which looks like this

![](/images/pipecd-control-plane-mainpage.png)

For more about PipeCD control plane management, please check [Managing ControlPlane](/docs/user-guide/managing-controlplane/).

#### 1.2. Installing Piped

Next, in order to perform CD tasks, you need to install a Piped agent to the cluster.

From your logged in tab, navigate to the PipeCD setting page at [http://localhost:8080/settings/piped?project=quickstart](http://localhost:8080/settings/piped?project=quickstart).

You will find the `+ADD` button around the top left of this page, click there and insert information to register the `piped`.

![](/images/quickstart-adding-piped.png)

Click on the `Save` button, and then you can see the piped-id and secret-key.

![](/images/quickstart-piped-registered.png)

You need to copy two values, `Piped Id` and `Base64 Encoded Piped Key`, and fill in `<COPIED_PIPED_ID>` and `<COPIED_ENCODED_PIPED_KEY>` respectively this below command

```console
$ curl -s https://raw.githubusercontent.com/pipe-cd/pipecd/master/quickstart/manifests/piped.yaml | \
  sed -e 's/<YOUR_PIPED_ID>/<COPIED_PIPED_ID>/g' \
      -e 's/<YOUR_PIPED_KEY_DATA>/<COPIED_ENCODED_PIPED_KEY>/g' | \
  kubectl apply -n pipecd -f -
```

For more about Piped management, please check [Managing Piped](/docs/user-guide/managing-piped/).

That's all! You are ready to use PipeCD to manage your application's deployment.

You can check the readiness of all PipeCD components via command

```console
$ kubectl get pod -n pipecd
NAME                              READY   STATUS    RESTARTS      AGE
pipecd-cache-56c7c65ddc-xqcst     1/1     Running   0             38m
pipecd-gateway-58589b55f9-9nbrv   1/1     Running   0             38m
pipecd-minio-677999d5bb-xnb78     1/1     Running   0             38m
pipecd-mysql-6fff49fbc7-hkvt4     1/1     Running   0             38m
pipecd-ops-779d6844db-nvbwn       1/1     Running   0             38m
pipecd-server-5769df7fcb-9hc45    1/1     Running   1 (38m ago)   38m
piped-8477b5d55d-74s5v            1/1     Running   0             97s
```

### 2. Deploy a Kubernetes application with PipeCD

Above is all that is necessary to set up your own PipeCD (both control plane and agent), let's use the installed one to deploy your first Kubernetes application with PipeCD.

Navigate to the `Applications` page, click on the `+ADD` button on the top left corner.

Go to the `ADD FROM SUGGESTIONS` tab, then select:
- Piped: `dev` (you just registered)
- PlatformProvider: `kubernetes-default`

You should see a lot of suggested applications. Select one of listed applications and click the `SAVE` button to register.

![](/images/quickstart-adding-application-from-suggestions.png)

After a bit, the first deployment is complete and will automatically sync the application to the state specified in the current Git commit.

![](/images/quickstart-first-deployment.png)

For more about manage applications' deployment with PipeCD, referrence to [Managing application](/docs/user-guide/managing-application/)

### 3. Cleanup
When you’re finished experimenting with PipeCD quickstart mode, you can uninstall it using:

``` console
$ kubectl delete ns pipecd
```

### What's next?

To prepare your PipeCD for a production environment, please visit the [Installation](../installation/) guideline. For guidelines to use PipeCD to deploy your application in daily usage, please visit the [User guide](../user-guide/) docs.
