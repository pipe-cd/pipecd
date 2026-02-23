---
date: 2026-02-24
weight: 977
title: "Install PipeCD locally using Kind"
linkTitle: "Install PipeCD locally using Kind"
author: Cornelius Emase ([@lochipi](https://github.com/lochipi))
categories: ["Tutorial"]
tags: ["PipeCD", "Kind", "Kubernetes", "Calico", "Local Development"]
---
## Video walkthrough

Prefer a video? Watch the step-by-step walkthrough here:

- https://youtu.be/rXRxLYCcquQ

## Big picture

This guide walks you through setting up a local PipeCD development environment using Kind (Kubernetes in Docker). The goal is to get PipeCD running on your machine in a simple way so you can catch up quickly with PipeCD.

By the end of this guide, you will have:

- A working Kubernetes cluster running in Docker.
- Calico installed and providing cluster networking.
- The PipeCD control plane running.
- A connected Piped agent ready to deploy applications.

This setup typically takes 20–30 minutes, depending on your system and network speed.

## When to use this guide

Use this guide if you want to:

- Try PipeCD locally without cloud infrastructure.
- Understand how PipeCD components fit together.
- Experiment safely before moving to production.

This is not a production setup. For production environments, refer to the official [PipeCD installation guide](https://pipecd.dev/docs-v1.0.x/installation/).

## Before you begin

### Required

You’ll need a Linux machine with:

- Docker installed and running.
- At least 2 CPUs (4+ recommended).
- 8GB RAM recommended.
- ~10GB free disk space.

You can also run Linux in a virtual machine (VM).

### Tools

Make sure the following tools are available:

```bash
docker --version
kubectl version --client
kind version
```

If any of these commands fail, install the missing tool before continuing.

Resources:

- https://docs.docker.com/engine/install/
- https://kind.sigs.k8s.io/
- https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/

## Concepts

### Why Kind?

Kind runs Kubernetes inside Docker containers. It’s lightweight, fast to reset, and ideal for local development.

### Why Calico?

PipeCD expects a properly configured Kubernetes network. Calico is a widely used CNI that works well 
with Kind and closely resembles real-world setups.

In this guide, we explicitly install Calico instead of relying on Kind’s default networking.

## How to

### 1. Create a Kubernetes cluster with Kind

We’ll start by creating a multi-node cluster: one control plane and two workers.

#### Create a Kind configuration file

You can create this file in the home directory or whichever directory you would like 
to(current development working directory);

```yaml
# kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
- role: worker
networking:
  disableDefaultCNI: true
```

We disable the default CNI so we can install Calico ourselves.

#### Create the cluster

```bash
kind create cluster --name mycluster --config kind-config.yaml
```

Verify the cluster is reachable:

```bash
kubectl cluster-info
kubectl get nodes
```

At this point, nodes may show `NotReady`. This is expected until networking is installed.

### 2. Install Calico

Next, install Calico to provide pod networking.

```bash
kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.27.0/manifests/calico.yaml
```

Watch the rollout:

```bash
kubectl get pods -n kube-system -w
```

You’ll see Calico pods initializing. This can take a few minutes.

#### Confirm Calico is running

```bash
kubectl get pods -n kube-system
```

Wait until:

- All `calico-node` pods are `Running`.
- `calico-kube-controllers` is `Running`.
- `coredns` pods are `Running`.

Then confirm node readiness:

```bash
kubectl get nodes
```

All nodes should now be `Ready`.

> Note:
> If your Kubernetes nodes are not in the `Ready` state, stop here.
> Fix networking (most likely Calico or other CNI) before installing PipeCD.
> PipeCD depends on a healthy cluster and will not work correctly otherwise.

### 3. Create the PipeCD namespace

```bash
kubectl create namespace pipecd
```

### 4. Deploy the PipeCD control plane

PipeCD provides a quickstart manifest that deploys all required components.

```bash
kubectl apply -n pipecd -f https://raw.githubusercontent.com/pipe-cd/pipecd/master/quickstart/manifests/control-plane.yaml
```

Watch the pods:

```bash
kubectl get pods -n pipecd -w
```

Wait until all PipeCD pods are in the `Running` state.

### 5. Access the PipeCD web UI

Forward the service locally:

```bash
kubectl port-forward -n pipecd svc/pipecd 8080:80
```

Open your browser and visit:

- http://localhost:8080?project=quickstart

You should now see the PipeCD UI login page.

> **Note**:
> To log in, you can use the configured static admin account as below:
>
> **username**: hello-pipecd
> **password**: hello-pipecd

![PipeCD UI login page](/images/install-kind-pipecd-ui-login.png)

### 6. Create and register a Piped agent

PipeCD uses `piped` as the agent that runs inside your cluster and executes deployments.

#### Create a Piped in the UI

In the PipeCD UI:

1. Go to **Settings → Piped**.
2. Create a new Piped.
3. Copy the generated:
   - Piped ID
   - Piped Key (base64 encoded)

Example values (for illustration only):

```bash
PIPED_ID=piped-sample-123
PIPED_KEY=LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0t...
```

#### Deploy Piped to the cluster

Export your values:

```bash
export PIPED_ID="piped-sample-123"
export PIPED_KEY="LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0t..."
```

Apply the Piped manifest, replacing placeholders using `sed`:

```bash
curl -s https://raw.githubusercontent.com/pipe-cd/pipecd/master/quickstart/manifests/piped.yaml | \
sed -e "s/<YOUR_PIPED_ID>/${PIPED_ID}/g" \
    -e "s/<YOUR_PIPED_KEY_DATA>/${PIPED_KEY}/g" | \
kubectl apply -n pipecd -f -
```
This command:

- Downloads the manifest.
- Injects your Piped credentials.
- Applies it directly to the cluster.

#### Verify Piped is running

```bash
kubectl get pods -n pipecd
```

You should see a `piped` pod in the `Running` state.

In the PipeCD UI, the Piped should now appear as Connected.
Connected Piped instances are marked with a blue dot next to the Piped name.

![Piped connected status in PipeCD UI](/images/install-kind-pipecd-piped-connected.png)

## Final checks

Before considering the setup complete, verify:

```bash
kubectl get nodes
kubectl get pods -A
kubectl get pods -n pipecd
```

All components should be healthy (`Running`).

![Cluster and PipeCD status checks](/images/install-kind-pipecd-final-checks-status.png)

## Cleanup

When you're done, you can safely tear down the cluster to free system resources and avoid 
state-related issues.

To delete the Kind cluster run:

```bash
kind delete cluster --name mycluster
```

## What’s next?

To prepare your PipeCD for a production environment, please [visit the Installation guideline](https://pipecd.dev/docs-v1.0.x/installation/). 
For guidelines to use PipeCD to deploy your application in daily usage, please visit the 
[User guide docs](https://pipecd.dev/docs-v1.0.x/user-guide/).
