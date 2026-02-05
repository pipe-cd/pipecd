# Simple Kubernetes Application Example

This example demonstrates how to deploy a basic Kubernetes application using PipeCD v1 with the plugin architecture.

## Overview

This is the simplest example of using PipeCD v1 to manage a Kubernetes deployment. It shows:
- Basic application configuration with the Kubernetes plugin
- Standard Kubernetes manifests (Deployment and Service)
- Simple sync deployment using the K8S_SYNC pipeline stage

## Files Included

### `app.pipecd.yaml`
The PipeCD v1 application configuration file that defines:
- **Application metadata**: Name and labels for organization
- **Plugin configuration**: Kubernetes plugin with target namespace
- **Deployment pipeline**: Single K8S_SYNC stage for straightforward deployment

This file tells PipeCD how to deploy and manage your application. The `.pipecd.yaml` suffix is required for PipeCD to discover it.

### `deployment.yaml`
A standard Kubernetes Deployment manifest that:
- Creates a Deployment named `simple-app`
- Runs 2 replicas of the `nginx:latest` container
- Exposes port 80 for web traffic

### `service.yaml`
A standard Kubernetes Service manifest that:
- Creates a ClusterIP Service named `simple-app`
- Routes traffic to the nginx pods on port 80
- Enables internal cluster communication

## Prerequisites

Before using this example, you need:
1. A running PipeCD v1 control plane
2. A registered Piped agent with Kubernetes plugin configured
3. Access to a Kubernetes cluster where the application will be deployed
4. This repository configured in your Piped's `repositories` section

## How to Use

### 1. Register the Application

In the PipeCD UI:
1. Navigate to the **Applications** page
2. Click the **"+ ADD"** button
3. Select **"ADD FROM SUGGESTIONS"**
4. Choose your Piped and the `kubernetes-default` deploy target
5. Select this application from the list
6. Click **"SAVE"**

### 2. Deploy

Once registered, PipeCD will:
1. Detect the `app.pipecd.yaml` file in this repository
2. Execute the K8S_SYNC pipeline stage
3. Apply both Kubernetes manifests (`deployment.yaml` and `service.yaml`)
4. Monitor the deployment status and report progress

### 3. Verify Deployment

Check your Kubernetes cluster to confirm the deployment:
kubectl get deployments
kubectl get pods
kubectl get services

You should see:
- A `simple-app` Deployment with 2/2 ready pods
- Two running pods with names like `simple-app-xxxxxxxxx-xxxxx`
- A ClusterIP Service named `simple-app` on port 80

## Expected Outcome

After successful deployment:
- **Deployment**: 2 nginx pods running and healthy
- **Service**: ClusterIP Service exposing the application internally
- **PipeCD Dashboard**: Shows application status as "HEALTHY" with deployment history
- **Continuous Sync**: PipeCD monitors and syncs the application state with Git

Any changes to the manifests in Git will trigger automatic synchronization.