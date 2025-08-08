---
title: "Migration Guide: PipeCD v0 to v1"
linkTitle: "Migration Guide"
weight: 30
description: >
  Step-by-step guide for migrating from PipeCD v0 to v1, including configuration updates and best practices
---

# Migration Guide: PipeCD v0 to v1

This comprehensive guide walks you through migrating from PipeCD v0 to the new plugin-based v1 architecture. The migration is designed to be gradual and non-disruptive, allowing you to move at your own pace.

## Migration Overview

### Key Principles

- **Zero Downtime**: Control plane supports both v0 and v1 Piped agents simultaneously
- **Gradual Migration**: Move projects and applications incrementally
- **Backward Compatibility**: v0 configurations continue to work during transition
- **Flexible Timeline**: No forced migration deadline until end of 2025

### Timeline and Support

- **Current State**: v0 and v1 coexistence
- **Migration Period**: 2024-2025
- **v0 End of Life**: December 31, 2025
- **v1 Feature Development**: All new features will be v1-only

## Pre-Migration Assessment

### 1. Inventory Your Current Setup

Document your existing PipeCD deployment:

```bash
# List all Piped agents
pipectl piped list

# Export application configurations
pipectl application list --project YOUR_PROJECT > apps_v0.yaml

# Document platform providers
pipectl platform-provider list --project YOUR_PROJECT > providers_v0.yaml
```

### 2. Compatibility Check

Verify your setup is ready for migration:

- **Control Plane Version**: Must be v0.52.0 or later
- **Application Types**: All current types are supported in v1
- **Custom Stages**: May require plugin development
- **Platform Integrations**: Check plugin availability

### 3. Migration Planning

Create a migration plan:

```yaml
# migration-plan.yaml
migration:
  phase1:
    - project: development
      applications: ["dev-app-1", "dev-app-2"]
      priority: high
  phase2:
    - project: staging  
      applications: ["staging-app"]
      priority: medium
  phase3:
    - project: production
      applications: ["prod-app-1", "prod-app-2"]
      priority: critical
```

## Configuration Migration

### Platform Providers → Deploy Targets

The biggest change is how deployment targets are configured.

#### v0 Configuration (Piped)

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: my-project
  pipedID: my-piped
  pipedKeyData: <base64-key>
  apiAddress: control-plane.example.com:8080
  
  platformProviders:
    - name: dev-k8s
      type: KUBERNETES
      config:
        masterURL: https://dev-k8s.example.com
        kubeConfigPath: /path/to/dev-kubeconfig
    
    - name: prod-k8s
      type: KUBERNETES  
      config:
        masterURL: https://prod-k8s.example.com
        kubeConfigPath: /path/to/prod-kubeconfig
    
    - name: aws-lambda
      type: LAMBDA
      config:
        region: us-west-2
        profile: lambda-deployer
```

#### v1 Configuration (Piped)

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: my-project
  pipedID: my-piped
  pipedKeyData: <base64-key>
  apiAddress: control-plane.example.com:8080
  
  plugins:
    - name: kubernetes
      port: 7001
      url: https://github.com/pipe-cd/pipecd/releases/download/v1.0.0/plugin_kubernetes_linux_amd64
      deployTargets:
        - name: dev-k8s
          labels:
            env: dev
            cluster: primary
          config:
            masterURL: https://dev-k8s.example.com
            kubeConfigPath: /path/to/dev-kubeconfig
        
        - name: prod-k8s
          labels:
            env: prod
            cluster: primary
          config:
            masterURL: https://prod-k8s.example.com  
            kubeConfigPath: /path/to/prod-kubeconfig
    
    - name: lambda
      port: 7002
      url: https://github.com/pipe-cd/pipecd/releases/download/v1.0.0/plugin_lambda_linux_amd64
      deployTargets:
        - name: aws-lambda
          labels:
            env: prod
            region: us-west-2
          config:
            region: us-west-2
            profile: lambda-deployer
```

### Application Kind → Labels

Application kind moves from a dedicated field to labels.

#### v0 Configuration (Application)

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-app
  kind: KUBERNETES
  platformProvider: dev-k8s
  
  input:
    manifests:
      - deployment.yaml
      - service.yaml
```

#### v1 Configuration (Application)

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
metadata:
  labels:
    kind: KUBERNETES
spec:
  name: my-app
  deployTargets:
    - dev-k8s
  
  plugins:
    kubernetes:
      input:
        manifests:
          - deployment.yaml
          - service.yaml
```

### Pipeline Configuration

Pipeline configurations remain largely the same, with stage names updated for v1 plugins.

#### v0 Pipeline

```yaml
spec:
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 10%
      - name: WAIT
        with:
          duration: 30s
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_CANARY_CLEAN
```

#### v1 Pipeline

```yaml
spec:
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT  # Same stage names
        with:
          replicas: 10%
      - name: WAIT
        with:
          duration: 30s
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_CANARY_CLEAN
```

## Step-by-Step Migration Process

### Step 1: Upgrade Control Plane

Ensure your control plane supports both v0 and v1:

```bash
# Check current version
pipectl version

# Upgrade control plane (if needed)
# Follow your standard upgrade procedure
helm upgrade pipecd pipecd/pipecd --version v0.52.0
```

### Step 2: Prepare v1 Piped Configuration

Create v1 configuration based on your v0 setup:

```bash
# Generate v1 config from v0 (example script)
./migrate-config.sh piped-v0-config.yaml > piped-v1-config.yaml
```

Example migration script:

```bash
#!/bin/bash
# migrate-config.sh

v0_config=$1
v1_config=${2:-piped-v1-config.yaml}

# Extract basic info
project_id=$(yq '.spec.projectID' $v0_config)
piped_id=$(yq '.spec.pipedID' $v0_config) 
api_address=$(yq '.spec.apiAddress' $v0_config)

# Generate v1 config template
cat > $v1_config << EOF
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: $project_id
  pipedID: $piped_id
  apiAddress: $api_address
  pipedKeyData: <UPDATE_THIS>
  
  plugins:
EOF

# Convert platform providers to plugins
yq '.spec.platformProviders[]' $v0_config | while read -r provider; do
  # Add conversion logic based on provider type
  echo "  # TODO: Convert provider: $provider"
done
```

### Step 3: Deploy v1 Piped Agent

Deploy v1 Piped alongside your existing v0 Piped:

```bash
# Download v1 Piped binary
curl -Lo piped-v1 https://github.com/pipe-cd/pipecd/releases/download/v1.0.0/piped_linux_amd64
chmod +x piped-v1

# Start v1 Piped with new configuration
./piped-v1 piped --config-file=piped-v1-config.yaml --tools-dir=/tmp/piped-bin
```

Or using Docker:

```yaml
# docker-compose.yml
version: '3.8'
services:
  piped-v1:
    image: pipecd/piped:v1.0.0
    volumes:
      - ./piped-v1-config.yaml:/etc/piped/config.yaml
      - ./kubeconfig:/etc/kubeconfig
    command: piped --config-file=/etc/piped/config.yaml
    environment:
      - KUBECONFIG=/etc/kubeconfig
```

### Step 4: Migrate Applications

Migrate applications one by one or in groups:

#### 4.1 Update Application Configuration

```yaml
# Before (v0)
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: my-app
  kind: KUBERNETES
  platformProvider: dev-k8s

# After (v1)  
apiVersion: pipecd.dev/v1beta1
kind: Application
metadata:
  labels:
    kind: KUBERNETES
spec:
  name: my-app
  deployTargets:
    - dev-k8s
  plugins:
    kubernetes:
      input:
        manifests:
          - deployment.yaml
          - service.yaml
```

#### 4.2 Test Migration

```bash
# Validate v1 configuration
pipectl application validate app.pipecd.yaml

# Deploy test application
git add app.pipecd.yaml
git commit -m "Migrate to PipeCD v1"
git push origin main
```

#### 4.3 Monitor Deployment

Watch the deployment in PipeCD console:
- Verify v1 Piped agent picks up the application
- Check deployment logs and status
- Validate functionality matches v0 behavior

### Step 5: Validate and Cleanup

After successful migration:

```bash
# Verify v1 deployment
pipectl application get my-app --project my-project

# Check deployment history
pipectl deployment list --application my-app

# Gradually migrate remaining applications
# Once all apps are migrated, decommission v0 Piped
```

## Common Migration Scenarios

### Kubernetes Applications

#### Basic Kubernetes App

```yaml
# v0 → v1 changes
platformProvider: dev-k8s → deployTargets: [dev-k8s]
kind: KUBERNETES → labels.kind: KUBERNETES

# Plugin configuration
plugins:
  kubernetes:
    input:
      manifests: [...] # Same as v0
      namespace: my-app
```

#### Multi-Cluster Kubernetes

```yaml
# v1 enables easier multi-cluster deployment
spec:
  deployTargets:
    - dev-cluster
    - staging-cluster
    - prod-cluster
  
  plugins:
    kubernetes:
      input:
        manifests:
          - deployment.yaml
          - service.yaml
```

### Terraform Applications

```yaml
# v0 configuration
spec:
  kind: TERRAFORM
  platformProvider: terraform-provider

# v1 configuration  
metadata:
  labels:
    kind: TERRAFORM
spec:
  deployTargets:
    - terraform-target
  
  plugins:
    terraform:
      input:
        workspace: production
        vars:
          - key: environment
            value: prod
```

### Lambda Applications

```yaml
# v0 configuration
spec:
  kind: LAMBDA
  platformProvider: aws-lambda

# v1 configuration
metadata:
  labels:
    kind: LAMBDA
spec:
  deployTargets:
    - aws-lambda
  
  plugins:
    lambda:
      input:
        functionManifests:
          - function.yaml
```

## Advanced Migration Topics

### Custom Deployment Strategies

If you have custom deployment logic in v0, you may need to:

1. **Develop Custom Plugin**: Create plugin for your specific needs
2. **Use ScriptRun Plugin**: Execute custom scripts as deployment stages
3. **Combine Plugins**: Use multiple plugins for complex workflows

### Secret Management Migration

```yaml
# v0 secret configuration
spec:
  secretManagement:
    keyProviders:
      - name: gcp-kms
        config:
          projectID: my-project
          keyRing: pipecd
          cryptoKey: piped-key

# v1 secret configuration (unchanged)
spec:
  secretManagement:
    keyProviders:
      - name: gcp-kms
        config:
          projectID: my-project
          keyRing: pipecd
          cryptoKey: piped-key
```

### Git Repository Migration

Git repository configuration remains the same:

```yaml
# Both v0 and v1
spec:
  repositories:
    - repoID: main-repo
      remote: git@github.com:myorg/manifests.git
      branch: main
```

## Migration Tools and Automation

### Configuration Converter

Create automated tools to convert configurations:

```python
#!/usr/bin/env python3
# migrate_config.py

import yaml
import sys

def convert_piped_config(v0_config):
    """Convert v0 Piped config to v1 format"""
    v1_config = {
        'apiVersion': 'pipecd.dev/v1beta1',
        'kind': 'Piped',
        'spec': {
            'projectID': v0_config['spec']['projectID'],
            'pipedID': v0_config['spec']['pipedID'],
            'apiAddress': v0_config['spec']['apiAddress'],
            'plugins': []
        }
    }
    
    # Convert platform providers to plugins
    plugins_map = {}
    
    for provider in v0_config['spec'].get('platformProviders', []):
        plugin_type = provider['type'].lower()
        
        if plugin_type not in plugins_map:
            plugins_map[plugin_type] = {
                'name': plugin_type,
                'port': 7000 + len(plugins_map),
                'url': f'https://github.com/pipe-cd/pipecd/releases/download/v1.0.0/plugin_{plugin_type}_linux_amd64',
                'deployTargets': []
            }
        
        deploy_target = {
            'name': provider['name'],
            'config': provider['config']
        }
        
        plugins_map[plugin_type]['deployTargets'].append(deploy_target)
    
    v1_config['spec']['plugins'] = list(plugins_map.values())
    return v1_config

def convert_app_config(v0_config):
    """Convert v0 Application config to v1 format"""
    v1_config = {
        'apiVersion': 'pipecd.dev/v1beta1',
        'kind': 'Application',
        'metadata': {
            'labels': {
                'kind': v0_config['spec']['kind']
            }
        },
        'spec': {
            'name': v0_config['spec']['name'],
            'deployTargets': [v0_config['spec']['platformProvider']],
        }
    }
    
    # Convert plugin-specific configuration
    plugin_name = v0_config['spec']['kind'].lower()
    
    if 'input' in v0_config['spec']:
        v1_config['spec']['plugins'] = {
            plugin_name: {
                'input': v0_config['spec']['input']
            }
        }
    
    # Copy pipeline if present
    if 'pipeline' in v0_config['spec']:
        v1_config['spec']['pipeline'] = v0_config['spec']['pipeline']
    
    return v1_config

if __name__ == '__main__':
    if len(sys.argv) != 3:
        print("Usage: migrate_config.py <input.yaml> <piped|app>")
        sys.exit(1)
    
    with open(sys.argv[1], 'r') as f:
        config = yaml.safe_load(f)
    
    if sys.argv[2] == 'piped':
        result = convert_piped_config(config)
    elif sys.argv[2] == 'app':
        result = convert_app_config(config)
    else:
        print("Type must be 'piped' or 'app'")
        sys.exit(1)
    
    print(yaml.dump(result, default_flow_style=False))
```

### Batch Migration Script

```bash
#!/bin/bash
# batch_migrate.sh

APPS_DIR="./applications"
BACKUP_DIR="./backup_v0"

# Create backup
mkdir -p $BACKUP_DIR
cp -r $APPS_DIR/* $BACKUP_DIR/

# Migrate all applications
for app_file in $APPS_DIR/*.yaml; do
    echo "Migrating $app_file..."
    
    # Convert configuration
    python3 migrate_config.py "$app_file" app > "${app_file}.v1"
    
    # Validate new configuration
    if pipectl application validate "${app_file}.v1"; then
        mv "${app_file}.v1" "$app_file"
        echo "✓ $app_file migrated successfully"
    else
        echo "✗ $app_file migration failed"
        rm "${app_file}.v1"
    fi
done
```

## Testing and Validation

### Pre-Migration Testing

```bash
# Test v1 Piped with dummy application
cat > test-app.yaml << EOF
apiVersion: pipecd.dev/v1beta1
kind: Application
metadata:
  labels:
    kind: KUBERNETES
spec:
  name: test-app
  deployTargets:
    - dev-k8s
  plugins:
    kubernetes:
      input:
        manifests:
          - test-deployment.yaml
EOF

# Deploy and verify
git add test-app.yaml
git commit -m "Test v1 migration"
git push origin test-branch
```

### Post-Migration Validation

```bash
# Compare v0 vs v1 deployment behavior
./compare_deployments.sh app-name v0-deployment-id v1-deployment-id

# Verify all features work
pipectl application sync my-app --project my-project
pipectl deployment list --application my-app
```

## Troubleshooting Migration Issues

### Common Problems

#### 1. Plugin Download Failures

```yaml
# Problem: Plugin URL unreachable
plugins:
  - name: kubernetes
    url: https://invalid-url.com/plugin

# Solution: Use local path during migration
plugins:
  - name: kubernetes
    url: file:///opt/pipecd/plugins/kubernetes
```

#### 2. Configuration Validation Errors

```bash
# Validate configuration before applying
pipectl application validate app.pipecd.yaml --strict

# Common fixes:
# - Update stage names for v1 plugins
# - Fix deploy target references
# - Update plugin configuration format
```

#### 3. Deploy Target Mismatches

```yaml
# Problem: Application references non-existent deploy target
spec:
  deployTargets:
    - non-existent-target

# Solution: Ensure deploy target exists in Piped config
plugins:
  - name: kubernetes
    deployTargets:
      - name: non-existent-target  # Add this
        config: {...}
```

### Rollback Procedures

If migration fails, you can rollback:

```bash
# Stop v1 Piped
systemctl stop piped-v1

# Restore v0 application configurations
git checkout v0-backup-branch
git push origin main

# v0 Piped will resume handling deployments
```

## Migration Checklist

### Pre-Migration

- [ ] Control plane is v0.52.0 or later
- [ ] All current applications are inventoried
- [ ] Migration plan is documented
- [ ] Backup of all configurations created
- [ ] Team is trained on v1 concepts

### During Migration

- [ ] v1 Piped configuration created and validated
- [ ] v1 Piped deployed and running
- [ ] Test application migrated successfully
- [ ] Application configurations converted
- [ ] Pipeline functionality verified
- [ ] Monitoring and logging working

### Post-Migration

- [ ] All applications migrated to v1
- [ ] v0 Piped decommissioned
- [ ] Team updated on v1 workflows
- [ ] Documentation updated
- [ ] Migration artifacts cleaned up

## Best Practices

### 1. Gradual Migration

- Start with development environments
- Migrate non-critical applications first
- Test thoroughly at each stage
- Have rollback plans ready

### 2. Configuration Management

- Use version control for all configurations
- Validate configurations before applying
- Document migration decisions
- Keep backups of working v0 configs

### 3. Team Coordination

- Communicate migration timeline clearly
- Train team on v1 concepts
- Establish new workflows for v1
- Update runbooks and documentation

### 4. Monitoring

- Monitor both v0 and v1 deployments during transition
- Set up alerts for migration issues
- Track migration progress
- Measure performance differences

## Support and Resources

### Getting Help

- **Documentation**: [pipecd.dev/docs-v1](https://pipecd.dev/docs-v1)
- **Slack**: [#pipecd-migration](https://cloud-native.slack.com/)
- **GitHub**: [Migration Issues](https://github.com/pipe-cd/pipecd/issues)
- **Community Calls**: Monthly migration office hours

### Migration Tools

- **Config Converter**: Automated v0 to v1 conversion
- **Validation Tools**: Configuration validation utilities
- **Migration Scripts**: Batch migration automation
- **Rollback Tools**: Quick rollback utilities

The migration to PipeCD v1 represents a significant improvement in flexibility and extensibility. While it requires careful planning and execution, the benefits of the plugin architecture make it a worthwhile investment for your deployment infrastructure.
