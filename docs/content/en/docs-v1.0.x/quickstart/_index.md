---
title: "Quickstart"
linkTitle: "Quickstart"
weight: 3
description: >
  This page describes how to quickly get started with PipeCD on Kubernetes.
---

PipeCD consists of two core components: The Control Plane and Piped (see [PipeCD Concepts](../concepts/)).

- **The Control Plane** can be thought of as a web service application that can be installed anywhere. It provides the web UI, API endpoints, and metadata storage.

- **Piped** is a lightweight agent that connects your infrastructure to the Control Plane. In PipeCD v1, each plugin (an external component) implements the deployment and synchronization logic for a specific application kind, such as Kubernetes or Terraform.

In this quickstart, you’ll install both components on a Kubernetes cluster and deploy a sample “hello world” application.

>Note:
>
>- It's not required to install the PipeCD control plane to the cluster where your applications are running (See [PipeCD best practices](/blog/2021/12/29/pipecd-best-practice-01-operate-your-own-pipecd-cluster/) to understand more about PipeCD in real life use cases).
>- If you want to experiment with PipeCD freely or don't have a Kubernetes cluster, we recommend using [this Tutorial](https://github.com/pipe-cd/tutorial).

---

### Prerequisites

- A running Kubernetes cluster.  
- [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/) installed and configured to connect to your cluster.

---

### 1. Installation

#### 1.1. Installing PipeCD Control Plane

```bash
kubectl create namespace pipecd
kubectl apply -n pipecd -f https://raw.githubusercontent.com/pipe-cd/pipecd/master/quickstart/manifests/control-plane.yaml
```

The PipeCD control plane will be installed with a default project named `quickstart`. To access the PipeCD Control Plane UI, run the following command

```bash
kubectl port-forward -n pipecd svc/pipecd 8080
```

You can access the PipeCD console at [http://localhost:8080?project=quickstart](http://localhost:8080?project=quickstart)

To login, you can use the configured static admin account as below:

- username: `hello-pipecd`
- password: `hello-pipecd`

And you will access the main page of PipeCD Control Plane console, which looks like this

![An image of the Control Plane](/images/pipecd-control-plane-mainpage.png)

For more about PipeCD control plane management, please check [Managing ControlPlane](/docs/user-guide/managing-controlplane/).

#### 1.2. Installing Piped

Next, in order to perform CD tasks, you need to install a Piped agent to the cluster.

From your logged in tab, navigate to the PipeCD setting page at [http://localhost:8080/settings/piped?project=quickstart](http://localhost:8080/settings/piped?project=quickstart).

You will find the `+ADD` button around the top left of this page, click there and insert information to register the Piped agent (for example, `dev`).

![Interface to add piped](/images/quickstart-adding-piped.png)

Click on the `Save` button, and then you can see the piped-id and secret-key.

![A successfully registered piped](/images/quickstart-piped-registered.png)

You need to copy two values, `Piped Id` and `Base64 Encoded Piped Key`, and fill in `<COPIED_PIPED_ID>` and `<COPIED_ENCODED_PIPED_KEY>` respectively this command:

```bash
$ curl -s https://raw.githubusercontent.com/pipe-cd/pipecd/refs/heads/master/quickstart/manifests/pipedv1-exp.yaml | \
  sed -e 's/<YOUR_PIPED_ID>/<COPIED_PIPED_ID>/g' \
      -e 's/<YOUR_PIPED_KEY_DATA>/<COPIED_ENCODED_PIPED_KEY>/g' | \
  kubectl apply -n pipecd -f -
```

For more information about Piped management, please check [Managing Piped](/docs/user-guide/managing-piped/).

That's all! You are ready to use PipeCD to manage your application's deployment.

You can check the readiness of all PipeCD components via command

```bash
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

---

### 2. Deploy a Kubernetes application with PipeCD

Above is all that is necessary to set up your own PipeCD (both control plane and agent), let's use the installed one to deploy your first Kubernetes application with PipeCD.

Navigate to the `Applications` page, click on the `+ADD` button on the top left corner.

Go to the `PIPED V1 ADD FROM SUGGESTIONS` tab, then select:

- Piped that you have just registered (e.g. `dev`)
- The deployment target (e.g. 'kubernetes')

You should see a lot of suggested applications. Select one of listed applications and click the `SAVE` button to register.

![Adding the application](/images/quickstart-adding-application-from-suggestions-v1.png)

After a bit, the first deployment is complete and will automatically sync the application to the state specified in the current Git commit.

![Preview your deployment](/images/quickstart-first-deployment.png)

For more information on manging application deployment with PipeCD, see [Managing application](/docs/user-guide/managing-application/)

---

### 3. Cleanup

When you’re finished experimenting with PipeCD quickstart mode, you can uninstall it using:

``` bash
kubectl delete ns pipecd
```

---

### What next?

To prepare your PipeCD for a production environment, please visit the [Installation](../installation/) guideline. For guidelines to use PipeCD to deploy your application in daily usage, please visit the [User guide](../user-guide/) docs.

To set up the development environment and start contributing to PipeCD, please visit the [Contributor guide](../contribution-guidelines/) docs.
