---
title: "Migrating applications across projects"
linkTitle: "Migrating across projects"
weight: 6
description: >
  How to safely migrate PipeCD managed applications from one projects or control plane to another using pipectl transfer.
---

This guide explains how to migrate pipeds and applications from one PipeCD project or control plane to another using the `pipectl transfer` command. This is useful when reorganizing projects, migrating to a new control plane instance, or consolidating multiple control planes.

## Overview

The `pipectl transfer` command provides a safe, reliable way to migrate your PipeCD deployment infrastructure between control planes.

### What gets transferred

- **Pipeds**: All piped (new IDs and API keys will be generated)
- **Applications**: All application configurations (both enabled and disabled)
- **Application state**: Enabled/disabled status is preserved

### What doesn't get transferred

- **Deployment history**: Past deployment records remain in the source control plane
- **Insights data**: Historical metrics and analytics data
- **User accounts**: User permissions and accounts are project-specific

### When to use this feature

- Migrating to a new control plane instance
- Reorganizing projects or consolidating control planes
- Moving applications between projects
- Disaster recovery scenarios

## Prerequisites

Before starting the migration, ensure you have:

1. **pipectl installed**: Version that includes the `transfer` command (v0.56.0 or later)
2. **API keys**: Valid API keys for both source and target control planes/projects

NOTE:

- The application workloads will be untouched during the migration process, since this only transfers the PipeCD application configurations.
- After migration, the applications will be triggered to run a new deployment by piped in the target project/control plane. Since it's just a QUICKSYNC deployment, it's expected to NOT having any changes to the application workloads. To ensure that, **you should check to ensure your application can be triggered to run a QUICKSYNC deployment successfully without any changes in the source project/control plane**.


## Migration Workflow Overview

<!-- TODO: Add architecture diagram showing source/target control plane relationship and data flow -->

The migration process consists of three main phases:

1. **Backup**: Export data from source control plane/project to a local JSON file
2. **Restore Pipeds**: Register pipeds on target control plane/project with new IDs and keys
3. **Restore Applications**: Create applications on target control plane/project after pipeds connect

**Expected Timeline**:
- Backup: Minutes (depending on number of applications)
- Piped restore: Seconds to minutes
- Piped configuration update and restart: Manual step, timing varies
- Application restore: Minutes (depending on number of applications)

NOTE:

- In case you worry about the application being triggered to run a new deployment while the migration is in progress, you can stop your piped agent after step (3) and restart it after the applications migration is complete.

## Step-by-Step Migration Guide

### Phase 1: Preparation

#### 1.1 Create API keys with appropriate permissions

<!-- TODO: Add screenshot showing where to create API keys in the web UI -->

API Keys role:
- For source control plane: Create API key with at least `READ ONLY` permission
- For target control plane: Create API key with at least `READ/WRITE` permission

Create API keys in both control planes via the web UI:
1. Navigate to Settings > API Keys
2. Click "+ ADD" button
3. Select role
4. Copy and securely store the generated API key

#### 1.2 Stop and document current piped

Before migration, stop your piped (from the source control plane/project) and save copies of your current piped configuration files. You'll need to update these with new IDs and keys after the piped restore phase.

### Phase 2: Backup Source Data

#### 2.1 Run backup command

Export all pipeds and applications from the source control plane:

``` console
pipectl transfer backup \
    --address=https://source-control-plane.example.com \
    --api-key=SOURCE_API_KEY \
    --output-file=backup.json
```

Alternative using API key file:

``` console
pipectl transfer backup \
    --address=https://source-control-plane.example.com \
    --api-key-file=/path/to/source-api-key.key \
    --output-file=backup.json
```

#### 2.2 Verify backup file

Check that the backup completed successfully:

``` console
# View backup file summary
cat backup.json | jq '{version, created_at, piped_count: (.pipeds | length), app_count: (.applications | length)}'
```

Expected output:

``` json
{
  "version": "1",
  "created_at": "2024-01-15T10:30:00Z",
  "pipeds": 5,
  "app_count": 42
}
```

### Phase 3: Restore Pipeds

#### 3.1 Register pipeds on target control plane

``` console
pipectl transfer restore piped \
    --address=https://target-control-plane.example.com \
    --api-key=TARGET_API_KEY \
    --input-file=backup.json \
    --output-file=piped-mapping.json
```

The command will output progress:

```
Restoring pipeds..., input-file: backup.json
Found 5 piped(s) in backup (created at 2024-01-15T10:30:00Z)
Registered piped, name: production-piped, old-id: abc123, new-id: def456
Registered piped, name: staging-piped, old-id: ghi789, new-id: jkl012
...
Piped restore completed. Update each piped config with the new ID and key...
```

#### 3.2 Understanding the piped mapping file

The `piped-mapping.json` file contains the mapping between old and new piped IDs:

``` json
{
  "piped_mappings": [
    {
      "old_piped_id": "abc123...",
      "new_piped_id": "def456...",
      "new_key": "newly-issued-piped-api-key",
      "piped_name": "production-piped"
    },
    ...
  ]
}
```

**Security Note**: This file contains newly generated piped API keys. Protect it appropriately and distribute keys securely to each piped environment.

#### 3.3 Update piped configurations

For each piped, update its configuration file with the new ID and key:

1. Locate the piped's configuration file (typically `piped-config.yaml`)
2. Find the mapping for this piped in `piped-mapping.json`
3. Update the configuration:

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: target-project # Update to target project
  pipedID: def456...  # Update with new_piped_id
  pipedKeyData: newly-issued-piped-api-key  # Update with new_key
  apiAddress: target-control-plane.example.com:443  # Update to target control plane
  # ... rest of configuration remains the same
```

#### 3.4 Restart piped and verify availability

After updating configurations, restart each piped agent and check it availability.

1. Open the target control plane web UI
2. Navigate to Settings > Pipeds
3. Confirm all pipeds show as "Online" with a green indicator
4. Check that repository sync has completed (may take a few minutes)

**Important**: Do not proceed to application restore until all pipeds are online and have completed their initial repository sync. This ensures Git repositories are registered before applications are created.

### Phase 4: Restore Applications

#### 4.1 Wait for repository synchronization

After pipeds connect, they automatically synchronize their configured Git repositories. This process registers repositories with the control plane, which is required before applications can be created.

Once your pipeds are online (green indicator in the web UI), configured Git repositories should be synced.

#### 4.2 Restore applications

Once all pipeds are connected and repositories are synced:

``` console
pipectl transfer restore application \
    --address=https://target-control-plane.example.com \
    --api-key=TARGET_API_KEY \
    --input-file=backup.json \
    --piped-id-mapping-file=piped-mapping.json
```

The command will output progress:

```
Restoring applications...
Loaded mapping for 5 piped(s)
Restored application, name: frontend-app, old-id: app1, new-id: app101
Restored application, name: backend-api, old-id: app2, new-id: app102 (disabled)
...
Application restore completed, restored: 40, failed: 2
```

#### 4.3 Understanding partial restore behavior

The restore command continues even if individual applications fail. Failed applications are logged with details:

```
Warning: failed to restore application frontend-app: repository not found
Warning: failed to restore application backend-api: missing piped mapping for piped xyz
```

Common failure reasons:
- Repository not yet synced by piped
- Piped mapping not found (piped wasn't migrated)
- Application configuration validation errors
- Insufficient permissions

Failed applications are skipped and logged. You can manually create them later via the web UI or `pipectl application add` command.

#### 4.4 Verify application creation

<!-- TODO: Add screenshot showing application list in web UI -->

Check that applications were created successfully:

1. Open the target control plane web UI
2. Navigate to Applications
3. Verify all expected applications appear in the list
4. Check that enabled/disabled status matches the source

Via command line:

``` console
pipectl application list \
    --address=https://target-control-plane.example.com \
    --api-key=TARGET_API_KEY
```

## Understanding Migration Behavior

### Zero Downtime Migration

The migration process does **not** touch your actual workloads (Kubernetes pods, ECS tasks, Lambda functions, etc.). Your applications continue running normally on their existing infrastructure during the entire migration process.

### QuickSync Deployment After Migration

When applications are created on the target control plane, the piped treats them as "new" applications and automatically triggers a **QuickSync** deployment strategy. This means:

1. **Initial sync occurs**: The piped compares the desired state (in Git) with the live state
2. **Expected outcome**: If your live workloads match the Git configuration, QuickSync will detect no changes and complete successfully without modifying anything
3. **Zero downtime**: Workloads remain unchanged during this process

### First Deployment Behavior

**Important**: The initial QuickSync deployment after migration is treated as a first-time deployment. If this deployment fails:

- **No rollback occurs**: There's no previous "good state" to roll back to on the target control plane
- **Manual intervention**: You'll need to investigate and fix the issue manually
- **Source unaffected**: Your applications continue running normally, managed by the source control plane until you decommission it

### State Transition

Applications transition from "managed by source control plane" to "managed by target control plane" at the application restore phase. During this transition:

- Deployment history does not carry over
- Drift detection starts fresh from the current state
- Application appears as "new" with no previous deployment records

## Validation and Testing (optional)

After completing the migration, perform these validation steps:

### 5.1 Verify all pipeds are connected

``` console
pipectl piped list \
    --address=https://target-control-plane.example.com \
    --api-key=TARGET_API_KEY
```

All pipeds should show as online.

### 5.2 Verify all applications are created

``` console
pipectl application list \
    --address=https://target-control-plane.example.com \
    --api-key=TARGET_API_KEY | jq '. | length'
```

Compare the count with your backup file's application count.

### 5.3 Check application sync status

<!-- TODO: Add screenshot showing application sync status and QuickSync deployments -->

Navigate to the Applications page in the web UI. Each application should show:
- Sync status: Synced (green)
- Recent deployment: QuickSync completed successfully


## Rollback Considerations (optional)

### When to roll back

Consider rolling back if:
- Many applications fail to restore
- Critical applications show sync failures (after long waiting period)

### How to roll back

Since workloads are not touched during migration, rollback is straightforward:

1. **Keep source control plane running**: Don't decommission the source control plane until you're confident in the migration
2. **Revert piped configurations**: Update piped configs back to source control plane settings
3. **Restart piped agents**: Reconnect pipeds to source control plane
4. **Clean up target control plane**: Optionally delete migrated applications from target (they were never deployed)

