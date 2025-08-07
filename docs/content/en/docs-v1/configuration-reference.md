---
title: "PipeCD v1 Configuration Reference"
linkTitle: "Configuration Reference"
weight: 40
description: >
  Complete reference for PipeCD v1 Piped and Application configurations
---

# PipeCD v1 Configuration Reference

This document provides a comprehensive reference for configuring PipeCD v1 Piped agents and applications. The v1 architecture introduces significant changes to configuration structure to support the new plugin-based system.

## Piped Configuration

The Piped configuration defines how the Piped agent connects to the control plane and configures plugins for deployment.

### Basic Structure

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  # Core connectivity configuration
  projectID: string
  pipedID: string
  pipedKeyData: string
  pipedKeyFile: string
  apiAddress: string
  webAddress: string
  
  # Plugin configuration
  plugins: []Plugin
  
  # Git repository configuration
  repositories: []GitRepository
  
  # Optional configurations
  syncInterval: duration
  secretManagement: SecretManagement
  notifications: []NotificationConfig
  insights: InsightConfig
```

### Core Configuration Fields

#### Connection Settings

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `projectID` | string | Yes | PipeCD project identifier |
| `pipedID` | string | Yes | Unique identifier for this Piped agent |
| `pipedKeyData` | string | Yes* | Base64 encoded Piped private key |
| `pipedKeyFile` | string | Yes* | Path to Piped private key file |
| `apiAddress` | string | Yes | Control plane gRPC API address (host:port) |
| `webAddress` | string | No | Control plane web address for console links |

*Either `pipedKeyData` or `pipedKeyFile` must be specified.

#### Plugin Configuration

```yaml
plugins:
  - name: string              # Plugin identifier
    port: int                 # Port for plugin gRPC server  
    url: string               # Plugin binary URL or local path
    config: object            # Plugin-specific configuration
    deployTargets: []DeployTarget  # Deployment targets for this plugin
```

#### Plugin Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Unique plugin name (used in application configs) |
| `port` | int | Yes | Unused port for plugin gRPC server |
| `url` | string | Yes | Plugin binary URL or file:// path |
| `config` | object | No | Plugin-specific global configuration |
| `deployTargets` | array | No | Available deployment targets for this plugin |

#### Deploy Target Configuration

```yaml
deployTargets:
  - name: string              # Target identifier
    labels: map[string]string # Target labels for selection
    config: object            # Target-specific configuration
```

### Complete Piped Configuration Example

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: my-project
  pipedID: piped-01
  pipedKeyFile: /etc/piped/key
  apiAddress: control-plane.example.com:8080
  webAddress: https://console.example.com
  syncInterval: 1m
  
  plugins:
    # Kubernetes plugin for container deployments
    - name: kubernetes
      port: 7001
      url: https://github.com/pipe-cd/pipecd/releases/download/v1.0.0/plugin_kubernetes_linux_amd64
      deployTargets:
        - name: dev-cluster
          labels:
            env: development
            region: us-west-1
          config:
            masterURL: https://dev-k8s.example.com
            kubeConfigPath: /etc/kubeconfig/dev
            kubectlVersion: "1.28.0"
        
        - name: prod-cluster
          labels:
            env: production
            region: us-west-1
          config:
            masterURL: https://prod-k8s.example.com
            kubeConfigPath: /etc/kubeconfig/prod
            kubectlVersion: "1.28.0"
    
    # Terraform plugin for infrastructure
    - name: terraform
      port: 7002
      url: https://github.com/pipe-cd/pipecd/releases/download/v1.0.0/plugin_terraform_linux_amd64
      deployTargets:
        - name: aws-infra
          labels:
            provider: aws
            env: production
          config:
            vars:
              - key: region
                value: us-west-1
              - key: environment
                value: prod
    
    # Wait plugin for deployment pauses
    - name: wait
      port: 7003
      url: https://github.com/pipe-cd/pipecd/releases/download/v1.0.0/plugin_wait_linux_amd64
    
    # Script execution plugin
    - name: scriptrun
      port: 7004
      url: https://github.com/pipe-cd/pipecd/releases/download/v1.0.0/plugin_scriptrun_linux_amd64
      config:
        allowedCommands:
          - "kubectl"
          - "helm"
          - "curl"
  
  repositories:
    - repoID: main-repo
      remote: git@github.com:myorg/k8s-manifests.git
      branch: main
    - repoID: infra-repo
      remote: git@github.com:myorg/terraform-infra.git
      branch: main
  
  secretManagement:
    keyProviders:
      - name: gcp-kms
        config:
          projectID: my-gcp-project
          keyRing: pipecd
          cryptoKey: piped-key
  
  notifications:
    routes:
      - name: slack-dev
        receiver: slack-dev
        events:
          - DEPLOYMENT_STARTED
          - DEPLOYMENT_FAILED
        ignoreFields:
          - "time"
    receivers:
      - name: slack-dev
        slack:
          hookURL: https://hooks.slack.com/services/xxx
  
  insights:
    applicationCountInterval: 24h
    insightCollectionInterval: 6h
```

## Application Configuration

Application configurations define how applications are deployed using v1 plugins.

### Basic Structure

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
metadata:
  labels:
    kind: string              # Application type (KUBERNETES, TERRAFORM, etc.)
spec:
  name: string                # Application name
  deployTargets: []string     # Target deployment environments
  
  plugins: map[string]object  # Plugin-specific configurations
  
  pipeline: Pipeline          # Deployment pipeline definition
  
  # Optional configurations
  description: string
  labels: map[string]string
  timeout: duration
  trigger: TriggerConfig
```

### Application Fields

#### Metadata

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `metadata.labels.kind` | string | Yes | Application type (KUBERNETES, TERRAFORM, LAMBDA, etc.) |

#### Core Specification

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Application name |
| `deployTargets` | array | Yes | List of deploy target names from Piped config |
| `plugins` | object | Yes | Plugin-specific configuration |
| `pipeline` | object | No | Deployment pipeline stages |

### Plugin-Specific Configurations

#### Kubernetes Plugin Configuration

```yaml
plugins:
  kubernetes:
    input:
      manifests: []string           # List of manifest files
      kubectlVersion: string        # kubectl version override
      kustomizeVersion: string      # kustomize version override
      kustomizeOptions: map[string]string
      helmVersion: string           # helm version override
      helmChart: HelmChart          # Helm chart configuration
      helmOptions: HelmOptions      # Helm deployment options
      namespace: string             # Target namespace
      autoCreateNamespace: bool     # Auto-create namespace
    
    service: K8sResourceReference   # Service resource reference
    workloads: []K8sResourceReference  # Workload resource references
    variantLabel: KubernetesVariantLabel  # Variant labeling configuration
```

##### Helm Chart Configuration

```yaml
helmChart:
  path: string        # Local chart path
  name: string        # Chart name
  version: string     # Chart version
  repository: string  # Chart repository URL
```

##### Helm Options

```yaml
helmOptions:
  releaseName: string              # Helm release name
  setValues: map[string]string     # --set values
  valueFiles: []string             # -f value files
  setFiles: map[string]string      # --set-file values
  apiVersions: []string            # Kubernetes API versions
  kubeVersion: string              # Kubernetes version
```

##### Resource References

```yaml
# Service configuration
service:
  kind: string    # Resource kind (e.g., "Service")
  name: string    # Resource name

# Workload configuration
workloads:
  - kind: string  # Resource kind (e.g., "Deployment")
    name: string  # Resource name
```

##### Variant Labeling

```yaml
variantLabel:
  key: string             # Label key (default: "pipecd.dev/variant")
  primaryValue: string    # Primary variant value (default: "primary")
  canaryValue: string     # Canary variant value (default: "canary")
  baselineValue: string   # Baseline variant value (default: "baseline")
```

#### Terraform Plugin Configuration

```yaml
plugins:
  terraform:
    input:
      workspace: string             # Terraform workspace
      terraformVersion: string      # Terraform version
      vars: []TerraformVar         # Variable definitions
      varFiles: []string           # Variable files
      autoApprove: bool            # Skip interactive approval
```

##### Terraform Variables

```yaml
vars:
  - key: string      # Variable name
    value: string    # Variable value
  - key: string
    value: string
```

#### Lambda Plugin Configuration

```yaml
plugins:
  lambda:
    input:
      functionManifests: []string   # Lambda function manifests
      ignoreSSLErrors: bool         # Ignore SSL errors
      timeout: duration             # Deployment timeout
```

#### ScriptRun Plugin Configuration

```yaml
plugins:
  scriptrun:
    input:
      script: string              # Script to execute
      env: map[string]string      # Environment variables
      onRollback: string         # Rollback script
```

### Pipeline Configuration

Pipelines define the deployment workflow using stages provided by plugins.

```yaml
pipeline:
  stages:
    - name: string          # Stage name
      desc: string          # Stage description  
      timeout: duration     # Stage timeout
      with: object          # Stage-specific configuration
```

#### Common Kubernetes Stages

##### K8S_CANARY_ROLLOUT

```yaml
- name: K8S_CANARY_ROLLOUT
  desc: "Deploy canary variant"
  timeout: 10m
  with:
    replicas: string        # Replica count or percentage (e.g., "50%")
    suffix: string          # Resource name suffix (default: "canary")
    createService: bool     # Create canary service (default: false)
    patches: []K8sResourcePatch  # Resource patches
```

##### K8S_PRIMARY_ROLLOUT

```yaml
- name: K8S_PRIMARY_ROLLOUT
  desc: "Update primary workloads"
  timeout: 10m
  with:
    suffix: string          # Resource name suffix (default: "primary")
    createService: bool     # Create primary service (default: true)
    prune: bool            # Prune unused resources (default: false)
```

##### K8S_BASELINE_ROLLOUT

```yaml
- name: K8S_BASELINE_ROLLOUT
  desc: "Deploy baseline variant for analysis"
  timeout: 10m
  with:
    replicas: string        # Replica count or percentage
    suffix: string          # Resource name suffix (default: "baseline")
    createService: bool     # Create baseline service (default: false)
```

##### K8S_CANARY_CLEAN / K8S_BASELINE_CLEAN

```yaml
- name: K8S_CANARY_CLEAN
  desc: "Remove canary resources"
  timeout: 5m

- name: K8S_BASELINE_CLEAN
  desc: "Remove baseline resources"
  timeout: 5m
```

#### Resource Patches

```yaml
patches:
  - target:
      kind: string          # Resource kind
      name: string          # Resource name
      documentRoot: string  # YAML path to patch (optional)
    ops:
      - op: string          # Operation: yaml-replace, yaml-add, yaml-remove, json-replace, text-regex
        path: string        # Target path (e.g., "$.spec.replicas")
        value: string       # New value
```

#### Wait Stages

```yaml
- name: WAIT
  desc: "Wait before proceeding"
  timeout: 5m
  with:
    duration: duration      # Wait duration (e.g., "30s", "5m")
```

#### Script Execution Stages

```yaml
- name: SCRIPT_RUN
  desc: "Execute custom script"
  timeout: 10m
  with:
    script: string          # Script content or file path
    env:                    # Environment variables
      KEY: value
```

### Complete Application Examples

#### Simple Kubernetes Application

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
metadata:
  labels:
    kind: KUBERNETES
spec:
  name: hello-world
  deployTargets:
    - dev-cluster
  
  plugins:
    kubernetes:
      input:
        manifests:
          - deployment.yaml
          - service.yaml
        namespace: hello-world
        autoCreateNamespace: true
      
      service:
        kind: Service
        name: hello-world
      
      workloads:
        - kind: Deployment
          name: hello-world
```

#### Canary Deployment with Analysis

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
metadata:
  labels:
    kind: KUBERNETES
spec:
  name: api-service
  deployTargets:
    - prod-cluster
  
  plugins:
    kubernetes:
      input:
        manifests:
          - deployment.yaml
          - service.yaml
          - configmap.yaml
        namespace: api
  
  pipeline:
    stages:
      # Deploy canary with 10% traffic
      - name: K8S_CANARY_ROLLOUT
        desc: "Deploy canary version"
        timeout: 10m
        with:
          replicas: 10%
          createService: true
          patches:
            - target:
                kind: Service
                name: api-service
              ops:
                - op: yaml-add
                  path: "$.metadata.annotations"
                  value: |
                    traffic.pipecd.dev/canary: "10"
      
      # Wait for metrics
      - name: WAIT
        desc: "Wait for metrics collection"
        with:
          duration: 5m
      
      # Promote to primary
      - name: K8S_PRIMARY_ROLLOUT
        desc: "Promote to primary"
        timeout: 10m
      
      # Clean up canary
      - name: K8S_CANARY_CLEAN
        desc: "Clean up canary resources"
        timeout: 5m
```

#### Helm Application

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
metadata:
  labels:
    kind: KUBERNETES
spec:
  name: nginx-helm
  deployTargets:
    - dev-cluster
  
  plugins:
    kubernetes:
      input:
        helmChart:
          name: nginx
          version: "1.15.0"
          repository: https://charts.bitnami.com/bitnami
        
        helmOptions:
          releaseName: my-nginx
          setValues:
            replicaCount: "3"
            service.type: LoadBalancer
          valueFiles:
            - values-dev.yaml
        
        namespace: nginx
        autoCreateNamespace: true
```

#### Terraform Infrastructure

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
metadata:
  labels:
    kind: TERRAFORM
spec:
  name: vpc-infrastructure
  deployTargets:
    - aws-infra
  
  plugins:
    terraform:
      input:
        workspace: production
        terraformVersion: "1.5.0"
        vars:
          - key: environment
            value: prod
          - key: vpc_cidr
            value: "10.0.0.0/16"
        varFiles:
          - terraform.tfvars
        autoApprove: false
  
  pipeline:
    stages:
      - name: TERRAFORM_PLAN
        desc: "Generate execution plan"
        timeout: 10m
      
      - name: WAIT_APPROVAL
        desc: "Manual approval required"
        timeout: 24h
      
      - name: TERRAFORM_APPLY
        desc: "Apply infrastructure changes"
        timeout: 30m
```

#### Multi-Plugin Application

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
metadata:
  labels:
    kind: KUBERNETES
spec:
  name: complex-app
  deployTargets:
    - prod-cluster
  
  plugins:
    kubernetes:
      input:
        manifests:
          - k8s/
        namespace: complex-app
  
  pipeline:
    stages:
      # Pre-deployment script
      - name: SCRIPT_RUN
        desc: "Pre-deployment checks"
        timeout: 5m
        with:
          script: |
            #!/bin/bash
            echo "Running pre-deployment validation..."
            kubectl get nodes
            helm version
      
      # Deploy canary
      - name: K8S_CANARY_ROLLOUT
        desc: "Deploy canary"
        timeout: 10m
        with:
          replicas: 20%
      
      # Wait and analyze
      - name: WAIT
        desc: "Wait for analysis"
        with:
          duration: 10m
      
      # Custom analysis script
      - name: SCRIPT_RUN
        desc: "Custom analysis"
        timeout: 5m
        with:
          script: |
            #!/bin/bash
            echo "Running custom analysis..."
            # Custom analysis logic here
            exit 0  # 0 = success, non-zero = failure
      
      # Promote to primary
      - name: K8S_PRIMARY_ROLLOUT
        desc: "Promote to primary"
        timeout: 10m
      
      # Cleanup
      - name: K8S_CANARY_CLEAN
        desc: "Clean up canary"
        timeout: 5m
      
      # Post-deployment script
      - name: SCRIPT_RUN
        desc: "Post-deployment tasks"
        timeout: 5m
        with:
          script: |
            #!/bin/bash
            echo "Running post-deployment tasks..."
            # Notification, cleanup, etc.
```

## Configuration Validation

### Piped Configuration Validation

```bash
# Validate Piped configuration
pipectl piped validate piped-config.yaml

# Check plugin availability
pipectl plugin list --piped-config piped-config.yaml

# Test plugin connectivity
pipectl plugin test kubernetes --piped-config piped-config.yaml
```

### Application Configuration Validation

```bash
# Validate application configuration
pipectl application validate app.pipecd.yaml

# Validate against specific Piped
pipectl application validate app.pipecd.yaml --piped-config piped-config.yaml

# Strict validation (catch warnings as errors)
pipectl application validate app.pipecd.yaml --strict
```

## Environment Variable Substitution

Both Piped and Application configurations support environment variable substitution:

```yaml
# Piped configuration
spec:
  pipedKeyData: $PIPED_KEY_DATA
  apiAddress: $CONTROL_PLANE_ADDRESS

# Application configuration  
plugins:
  kubernetes:
    input:
      namespace: $ENVIRONMENT
      
# Set environment variables
export PIPED_KEY_DATA=$(echo -n "secret-key" | base64)
export CONTROL_PLANE_ADDRESS="control-plane.example.com:8080"
export ENVIRONMENT="production"
```

## Configuration Best Practices

### 1. Security

- Store sensitive data in external secret management systems
- Use `pipedKeyFile` instead of `pipedKeyData` when possible
- Limit plugin permissions through configuration
- Regularly rotate Piped keys

### 2. Organization

- Use consistent naming conventions for deploy targets
- Group related applications in the same repository
- Use labels for deploy target selection
- Document configuration decisions

### 3. Maintainability

- Version your configurations alongside code
- Use configuration validation in CI/CD
- Keep plugin versions pinned for stability
- Document custom stage configurations

### 4. Performance

- Configure appropriate timeouts for stages
- Use resource limits where supported
- Monitor plugin performance
- Cache plugin binaries locally when possible

## Troubleshooting Configuration Issues

### Common Validation Errors

```bash
# Plugin not found
Error: Plugin 'kubernetes' not configured in Piped

# Fix: Add plugin to Piped configuration
plugins:
  - name: kubernetes
    ...

# Deploy target not found  
Error: Deploy target 'prod-cluster' not found

# Fix: Add deploy target to plugin configuration
deployTargets:
  - name: prod-cluster
    ...

# Invalid stage configuration
Error: Stage 'K8S_CANARY_ROLLOUT' configuration invalid

# Fix: Check stage-specific configuration
with:
  replicas: "10%"  # Must be string, not number
```

### Configuration Debugging

```bash
# Enable debug logging
export PIPECD_LOG_LEVEL=debug

# Validate with verbose output
pipectl application validate app.pipecd.yaml --verbose

# Check plugin status
pipectl plugin status --piped-id my-piped
```

This configuration reference provides the foundation for setting up and managing PipeCD v1 deployments. For specific plugin configurations and advanced use cases, refer to the individual plugin documentation.
