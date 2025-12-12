---
date: 2026-01-05
title: "How to set up a local development environment"
linkTitle: "How to set up a local development environment"
weight: 990 
author: Shivansh Sahu ([@shivansh-gohem](https://github.com/shivansh-gohem))
categories: ["Tutorial", "Contributing"]
tags: ["Local Development", "Go", "Kubernetes", "Dev Setup"]
---


## Why a Local Development Environment is Essential

If you plan to contribute code, fix bugs, or add new plugins to PipeCD, you should run the control plane and Piped agent directly from the source code. This local development setup lets you test your Go changes live before submitting a pull request.

This guide walks through building and running the components from source so you can iterate quickly and confidently.

**Estimated setup time: 20–30 minutes**

## Prerequisites

Make sure you have the following installed and configured:

- **Go** (v1.21 or higher)
- **Docker**
- **kubectl**
- **Kind** (required to create the local cluster)
- **Node.js and Yarn** (required to build the web UI)
- **A fork of the PipeCD repository**

To contribute to PipeCD, we recommend setting up a fork:

1. Fork `pipe-cd/pipecd` on GitHub.  
2. Clone your fork locally:
```bash
git clone https://github.com/<YOUR_USERNAME>/pipecd.git
cd pipecd
```
## 1. Prepare and Start the Local Cluster

This step prepares a local Kubernetes cluster and container registry required for running the PipeCD control plane locally.

### Update dependencies

Run these commands to ensure your local Go modules and web dependencies are up to date. Starting the environment may fail if these are outdated.
```bash
make update/go-deps
make update/web-deps
```

### Start local registry and cluster

A helper command starts a local kind cluster and a container registry. This command also automatically creates the `pipecd` namespace where the components will run.
```bash
make up/local-cluster
```
After the cluster starts, export the kubeconfig to ensure kubectl can communicate with it:

```bash
kind export kubeconfig --name pipecd
```
## 2. Run the PipeCD Control Plane (from source)

The control plane provides the web UI, API, and metadata storage. Running it from source ensures you are testing your latest changes.

### Start the control plane

This command compiles the Go code, builds the web assets, and starts the control plane server locally.
```bash
make run/pipecd
```
### Access the UI

Once the control plane is running, forward the port to access the UI from your browser. Open a new terminal and run:
```bash
kubectl port-forward -n pipecd svc/pipecd 8080
```
Then open your browser:

- URL: <http://localhost:8080?project=quickstart>  
- Username: `hello-pipecd`  
- Password: `hello-pipecd`

## 3. Configure and Run the Piped Agent (from source)

The Piped agent connects the control plane to your local Kubernetes cluster. You will run this agent from source as well.

### Register Piped in the UI

1. Go to **Settings → Piped** (or open <http://localhost:8080/settings/piped?project=quickstart>).  
2. Click **+ ADD**, give it a name (for example, `dev`), and save.  
3. **Crucial step:** Copy the generated **Piped ID** and **Base64 encoded key** immediately.

### Create the Piped configuration

Create a file named `piped-config.yaml` in your repo root:
```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: quickstart
  # Replace here with your piped ID.
  pipedID: <COPIED_PIPED_ID>
  # Base64 encoded string of the piped private key.
  # Replace here with your piped base64 key.
  pipedKeyData: <COPIED_ENCODED_PIPED_KEY>
  apiAddress: localhost:8080
  repositories:
  - repoId: example
    remote: git@github.com:pipe-cd/examples.git
    branch: master
  syncInterval: 1m
  platformProviders:
  - name: example-kubernetes
    type: KUBERNETES
    config:
      # Replace here with your kubeconfig absolute file path.
      kubeConfigPath: /path/to/.kube/config
```

### Run Piped from source

Use your local code and the config file you just created:
```bash
make run/piped CONFIG_FILE=piped-config.yaml INSECURE=true
```

## Cleaning Up

To stop and remove the local cluster and registry when you are done, run:

```bash
make down/local-cluster
```
## Next Steps

Once Piped starts and shows as connected in the UI, your local development environment is ready. You can now build, test, and iterate on your PipeCD changes.
