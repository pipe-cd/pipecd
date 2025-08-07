---
title: "PipeCD v1 Quick Start Guide"
linkTitle: "Quick Start"
weight: 5
description: >
  Get up and running with PipeCD v1 in under 30 minutes
---

# PipeCD v1 Quick Start Guide

This guide will get you up and running with PipeCD v1's plugin architecture in under 30 minutes. You'll deploy a sample Kubernetes application using the new plugin-based system.

## Prerequisites

Before starting, ensure you have:

- A running Kubernetes cluster (local or remote)
- `kubectl` configured to access your cluster
- Git repository for storing manifests
- PipeCD Control Plane v0.52.0 or later

## Step 1: Set Up Control Plane

If you don't have a PipeCD Control Plane running, deploy one:

### Using Helm

```bash
# Add PipeCD Helm repository
helm repo add pipecd https://charts.pipecd.dev
helm repo update

# Install PipeCD Control Plane
helm install pipecd pipecd/pipecd \
  --namespace pipecd \
  --create-namespace \
  --set config.projects[0].id=quickstart \
  --set config.projects[0].staticAdmin.username=admin \
  --set config.projects[0].staticAdmin.passwordHash='$2a$10$ye96mUqUqTnjUqgwQJbJzel/LJh3jqfKNbqXSUU13EaMLpSVnCXmu'  # "admin"
```

### Access Control Plane

```bash
# Port forward to access the console
kubectl port-forward -n pipecd svc/pipecd 8080:8080

# Open browser to http://localhost:8080
# Login with: admin / admin
```

## Step 2: Create Piped Configuration

Create a v1 Piped configuration file:

```yaml
# piped-v1-config.yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: quickstart
  pipedID: quickstart-piped  # Will be generated in next step
  pipedKeyData: ""          # Will be generated in next step
  apiAddress: localhost:8080
  syncInterval: 1m
  
  repositories:
    - repoID: quickstart-repo
      remote: https://github.com/your-username/pipecd-quickstart.git  # Update this
      branch: main
  
  plugins:
    - name: kubernetes
      port: 7001
      url: https://github.com/pipe-cd/pipecd/releases/download/v1.0.0/plugin_kubernetes_linux_amd64
      deployTargets:
        - name: default-cluster
          labels:
            env: quickstart
          config:
            # For local clusters (kind, minikube, etc.)
            kubeConfigPath: ~/.kube/config
            # For remote clusters, set masterURL and kubeConfigPath
            # masterURL: https://your-cluster.example.com
            # kubeConfigPath: /path/to/kubeconfig
    
    - name: wait
      port: 7002
      url: https://github.com/pipe-cd/pipecd/releases/download/v1.0.0/plugin_wait_linux_amd64
```

## Step 3: Register Piped Agent

### Generate Piped ID and Key

1. Open PipeCD Console at http://localhost:8080
2. Login with `admin` / `admin`
3. Navigate to Settings → Piped
4. Click **"+ ADD"** to create a new Piped
5. Copy the generated **Piped ID** and **Key**

### Update Configuration

Update your `piped-v1-config.yaml` with the generated values:

```yaml
spec:
  pipedID: "<GENERATED_PIPED_ID>"
  pipedKeyData: "<GENERATED_BASE64_KEY>"
```

## Step 4: Download and Start Piped v1

### Download Piped v1 Binary

```bash
# Linux
curl -Lo piped-v1 https://github.com/pipe-cd/pipecd/releases/download/v1.0.0/piped_linux_amd64
chmod +x piped-v1

# macOS  
curl -Lo piped-v1 https://github.com/pipe-cd/pipecd/releases/download/v1.0.0/piped_darwin_amd64
chmod +x piped-v1

# Windows
curl -Lo piped-v1.exe https://github.com/pipe-cd/pipecd/releases/download/v1.0.0/piped_windows_amd64.exe
```

### Start Piped Agent

```bash
# Create tools directory
mkdir -p /tmp/piped-tools

# Start Piped v1
./piped-v1 piped \
  --config-file=piped-v1-config.yaml \
  --tools-dir=/tmp/piped-tools \
  --insecure=true  # Only for local development
```

You should see logs indicating successful plugin startup:

```
INFO  Starting Piped agent
INFO  Loading plugin: kubernetes
INFO  Plugin kubernetes started on port 7001
INFO  Loading plugin: wait  
INFO  Plugin wait started on port 7002
INFO  All plugins loaded successfully
INFO  Piped agent started successfully
```

## Step 5: Prepare Sample Application

### Create Git Repository

Create a new Git repository (or use an existing one) for your manifests:

```bash
mkdir pipecd-quickstart
cd pipecd-quickstart
git init
```

### Create Application Manifests

Create the following files:

#### `app.pipecd.yaml` (Application Configuration)

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
metadata:
  labels:
    kind: KUBERNETES
spec:
  name: hello-pipecd
  description: "Hello PipeCD v1 sample application"
  deployTargets:
    - default-cluster
  
  plugins:
    kubernetes:
      input:
        manifests:
          - deployment.yaml
          - service.yaml
        namespace: hello-pipecd
        autoCreateNamespace: true
      
      service:
        kind: Service
        name: hello-pipecd
      
      workloads:
        - kind: Deployment
          name: hello-pipecd
  
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
        desc: "Deploy canary with 50% traffic"
        with:
          replicas: 50%
      
      - name: WAIT
        desc: "Wait 30 seconds for validation"
        with:
          duration: 30s
      
      - name: K8S_PRIMARY_ROLLOUT
        desc: "Promote canary to primary"
      
      - name: K8S_CANARY_CLEAN
        desc: "Clean up canary resources"
```

#### `deployment.yaml` (Kubernetes Deployment)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-pipecd
  labels:
    app: hello-pipecd
spec:
  replicas: 2
  selector:
    matchLabels:
      app: hello-pipecd
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: hello-pipecd
        pipecd.dev/variant: primary
    spec:
      containers:
      - name: hello-pipecd
        image: ghcr.io/pipe-cd/helloworld:v0.45.0
        ports:
        - containerPort: 9085
        env:
        - name: PORT
          value: "9085"
        - name: MESSAGE
          value: "Hello from PipeCD v1!"
```

#### `service.yaml` (Kubernetes Service)

```yaml
apiVersion: v1
kind: Service
metadata:
  name: hello-pipecd
spec:
  selector:
    app: hello-pipecd
  ports:
    - protocol: TCP
      port: 9085
      targetPort: 9085
  type: ClusterIP
```

### Commit and Push

```bash
git add .
git commit -m "Add PipeCD v1 sample application"
git remote add origin https://github.com/your-username/pipecd-quickstart.git
git push -u origin main
```

## Step 6: Register Application

### Register in PipeCD Console

1. Go to PipeCD Console → Applications
2. Click **"+ ADD FROM SUGGESTIONS"** (Piped V1 tab)
3. Select your Piped agent (`quickstart-piped`)
4. Select the discovered application (`hello-pipecd`)
5. Click **"ADD"**

### Verify Registration

You should see the application appear in the Applications list with status "Healthy".

## Step 7: Deploy Application

### Trigger Deployment

Make a change to trigger deployment:

```bash
# Update the image version
sed -i 's/v0.45.0/v0.46.0/g' deployment.yaml

# Update the message
sed -i 's/Hello from PipeCD v1!/Hello from PipeCD v1 - Updated!/g' deployment.yaml

git add deployment.yaml
git commit -m "Update application version"
git push origin main
```

### Monitor Deployment

1. Go to PipeCD Console → Applications → hello-pipecd
2. Click on the new deployment that appears
3. Watch the pipeline progress through stages:
   - **K8S_CANARY_ROLLOUT**: Deploys canary with 50% replicas
   - **WAIT**: Waits 30 seconds
   - **K8S_PRIMARY_ROLLOUT**: Promotes canary to primary
   - **K8S_CANARY_CLEAN**: Removes canary resources

### Verify Deployment

```bash
# Check deployment status
kubectl get deployments -n hello-pipecd

# Check pods
kubectl get pods -n hello-pipecd

# Test the service
kubectl port-forward -n hello-pipecd service/hello-pipecd 9085:9085

# In another terminal
curl http://localhost:9085
# Should return: Hello from PipeCD v1 - Updated!
```

## Step 8: Explore Advanced Features

### Try Different Deployment Strategies

Update your `app.pipecd.yaml` to try different strategies:

#### Blue-Green Deployment

```yaml
pipeline:
  stages:
    - name: K8S_BASELINE_ROLLOUT
      desc: "Deploy baseline version"
      with:
        replicas: 100%
    
    - name: K8S_CANARY_ROLLOUT
      desc: "Deploy new version (blue-green)"
      with:
        replicas: 100%
    
    - name: WAIT
      desc: "Manual validation period"
      with:
        duration: 2m
    
    - name: K8S_PRIMARY_ROLLOUT
      desc: "Switch traffic to new version"
    
    - name: K8S_BASELINE_CLEAN
      desc: "Remove old version"
```

#### Script-Based Deployment

Add the ScriptRun plugin to your Piped config:

```yaml
plugins:
  - name: scriptrun
    port: 7003
    url: https://github.com/pipe-cd/pipecd/releases/download/v1.0.0/plugin_scriptrun_linux_amd64
```

Then use custom scripts in your pipeline:

```yaml
pipeline:
  stages:
    - name: SCRIPT_RUN
      desc: "Pre-deployment validation"
      with:
        script: |
          #!/bin/bash
          echo "Running pre-deployment checks..."
          kubectl cluster-info
          echo "Cluster is ready for deployment"
    
    - name: K8S_SYNC
      desc: "Deploy application"
    
    - name: SCRIPT_RUN
      desc: "Post-deployment testing"
      with:
        script: |
          #!/bin/bash
          echo "Running post-deployment tests..."
          # Add your custom tests here
          echo "All tests passed!"
```

## Troubleshooting

### Piped Agent Issues

```bash
# Check Piped logs
./piped-v1 piped --config-file=piped-v1-config.yaml --log-level=debug

# Verify plugin connectivity
kubectl logs -n hello-pipecd deployment/hello-pipecd

# Check plugin status in PipeCD Console
# Go to Settings → Piped → Your Piped → View Details
```

### Plugin Issues

```bash
# Check plugin downloads
ls -la ~/.piped/plugins/

# Test plugin manually
~/.piped/plugins/kubernetes --help

# Verify plugin ports are not in use
netstat -tulpn | grep :7001
```

### Common Issues

1. **Plugin download fails**: Check internet connectivity and GitHub access
2. **Kubernetes connection fails**: Verify kubeconfig and cluster access
3. **Application not detected**: Ensure `app.pipecd.yaml` is in repository root
4. **Pipeline fails**: Check stage configuration and plugin compatibility

## Next Steps

Congratulations! You've successfully deployed an application using PipeCD v1's plugin architecture. Here are some next steps:

### Explore More Plugins

- **[Terraform Plugin](../plugins/terraform/)**: Deploy infrastructure as code
- **[Lambda Plugin](../plugins/lambda/)**: Deploy AWS Lambda functions
- **[Community Plugins](../community/plugins/)**: Discover community-developed plugins

### Advanced Configuration

- **[Multi-Cluster Deployments](../advanced/multi-cluster/)**: Deploy across multiple environments
- **[Custom Stages](../advanced/custom-stages/)**: Implement custom deployment logic
- **[Security Configuration](../advanced/security/)**: Secure your deployment pipeline

### Plugin Development

- **[Plugin Development Guide](../plugin-development/)**: Create your own plugins
- **[SDK Reference](../sdk-reference/)**: Detailed SDK documentation
- **[Plugin Examples](../plugin-examples/)**: Learn from real-world examples

### Migration

- **[Migration Guide](../migration-guide/)**: Migrate from PipeCD v0 to v1
- **[Configuration Reference](../configuration-reference/)**: Complete configuration documentation

## Getting Help

- **Documentation**: [pipecd.dev/docs-v1](https://pipecd.dev/docs-v1)
- **Slack**: [#pipecd](https://cloud-native.slack.com/archives/C01B27F9T0X)
- **GitHub**: [Issues](https://github.com/pipe-cd/pipecd/issues)
- **Community**: [Discussions](https://github.com/pipe-cd/pipecd/discussions)

You're now ready to harness the full power of PipeCD v1's plugin architecture for your deployment needs!
