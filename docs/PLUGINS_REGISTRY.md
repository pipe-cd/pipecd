# PipeCD Plugins Registry System Documentation

## Overview

The PipeCD Plugins Registry System provides a centralized, automatically-updated source of truth for all official PipeCD plugin versions. This eliminates the need to manually search through GitHub releases to find the latest plugin versions.

## Components

### 1. Registry Files

#### `docs/plugins.json` - Machine-Readable Registry
Structured JSON registry containing all plugin metadata and versions. Designed for programmatic consumption by tools, websites, and automation.

**Key fields:**
- `plugins[]`: Array of plugin objects
- `plugins[].latestVersion`: Latest stable release version
- `plugins[].releaseUrl`: Direct link to the release
- `plugins[].tagPattern`: Git tag pattern used to detect releases
- `lastUpdated`: ISO 8601 timestamp of last update

**Usage example:**
```python
import json

with open('docs/plugins.json') as f:
    registry = json.load(f)
    
for plugin in registry['plugins']:
    print(f"{plugin['name']}: {plugin['latestVersion']}")
```

#### `docs/plugins.md` - Human-Readable Registry
Formatted Markdown document providing an easy-to-read overview of all plugins with their current versions. Suitable for documentation, websites, and manual reference.

**Includes:**
- Quick reference table with all plugin versions
- Detailed description of each plugin
- Links to documentation and release pages
- Explanation of the version tracking system

#### `docs/plugins.schema.json` - JSON Schema
JSON Schema v7 definition for validating registry files. Ensures data integrity and enables IDE/tool validation.

### 2. Update Scripts

#### `scripts/update-plugins-registry.py`
Main Python script that:
1. Connects to GitHub API
2. Queries each plugin's repository for latest releases/tags
3. Matches release tags against plugin-specific patterns
4. Generates both JSON and Markdown registry files
5. Implements semantic version comparison for accurate version detection

**Features:**
- Support for both inline (in pipecd repo) and external plugins
- Customizable tag patterns for flexible release tracking
- Rate limit handling with GitHub API
- Token support for authenticated requests (higher limits)

**Usage:**
```bash
# With GitHub token for higher rate limits
GITHUB_TOKEN=<token> python3 scripts/update-plugins-registry.py

# Custom output directory
python3 scripts/update-plugins-registry.py --output-dir /path/to/output

# Without token
python3 scripts/update-plugins-registry.py
```

#### `scripts/validate-plugins-registry.py`
Validation script that:
1. Validates registry against JSON schema
2. Performs semantic checks on data integrity
3. Reports errors and warnings
4. Suitable for CI/CD pipelines

**Usage:**
```bash
python3 scripts/validate-plugins-registry.py

# Custom paths
python3 scripts/validate-plugins-registry.py --registry docs/plugins.json --schema docs/plugins.schema.json
```

### 3. GitHub Actions Workflow

**File:** `.github/workflows/update-plugins-registry.yaml`

**Triggers:**
1. **On Release** - When new plugins or components are released
2. **On Scheduled Trigger** - Every 6 hours for consistency
3. **On Workflow Success** - After main pipecd release workflow completes
4. **Manual Trigger** - Via `workflow_dispatch` for testing

**Process:**
1. Checks out repository
2. Sets up Python environment
3. Installs dependencies
4. Runs `update-plugins-registry.py`
5. Detects changes in registry files
6. **If changes detected:**
   - Commits changes with descriptive message
   - Pushes to master branch
   - Creates PR for visibility (on workflow_run trigger)

**Safety features:**
- Only commits if actual version changes detected
- Uses `pipecd-bot` account for commits
- Signed commits for security
- Minimal diff detection prevents churn

## Plugin Configuration

### Adding New Plugins

To add a new official plugin to the registry:

1. **Edit** `scripts/update-plugins-registry.py`
2. **Add entry** to `PLUGINS_CONFIG` list with:
   - `id`: Unique identifier (slug format)
   - `name`: Human-readable name
   - `description`: Short description
   - `sourcePath`: Path in repository
   - `repository`: GitHub repository URL
   - `repositoryType`: "inline" or "external"
   - `tagPattern`: Git tag pattern (supports wildcards)
   - `status`: "stable", "beta", "alpha", or "deprecated"

3. **Example:**
```python
{
    "id": "my-plugin",
    "name": "My Custom Plugin",
    "description": "Does something cool",
    "sourcePath": "pkg/app/pipedv1/plugin/myplugin",
    "repository": "https://github.com/pipe-cd/pipecd",
    "repositoryType": "inline",
    "tagPattern": "pkg/app/pipedv1/plugin/myplugin/*",
    "status": "stable",
}
```

### Tag Pattern Syntax

Tag patterns use glob-like syntax:
- `pkg/app/pipedv1/plugin/kubernetes/*` - Matches plugin tags in pipecd repo
- `v*` - Matches semantic version tags in external repos
- `release/*` - Matches release/* prefixed tags

Patterns are converted to regex internally with `*` as wildcard.

## Version Detection Algorithm

1. **Fetch releases** from GitHub API with pagination
2. **Match tags** against plugin's `tagPattern`
3. **Extract version** from matched tag
4. **Compare versions** using semantic versioning
5. **Return latest** stable version

### Supported Version Formats

- Semantic versions: `v1.2.3`, `1.2.3`, `0.1.0`
- Path-prefixed tags: `pkg/app/pipedv1/plugin/kubernetes/v1.0.0`
- Pre-release versions: `v1.0.0-beta.1`

## Current Official Plugins

The registry currently tracks the following 9 official plugins:

1. **Kubernetes Plugin** - Deploy to Kubernetes clusters
2. **Terraform Plugin** - Infrastructure as Code deployments
3. **Cloud Run Plugin** - Google Cloud Run deployments
4. **Wait Stage Plugin** - Add delays to pipelines
5. **Wait Approval Plugin** - Add manual approval gates
6. **Script Run Plugin** - Execute custom scripts
7. **Analysis Plugin** - Deployment metrics analysis
8. **Kubernetes Multi-cluster Plugin** - Multi-cluster deployments
9. **PipeCD Plugin SDK for Go** - Plugin development SDK

## Integration Points

### pipecd.dev Website

The `plugins.json` file can be consumed by the pipecd.dev website to:
- Display plugin version badges
- Show plugin compatibility matrix
- Link to release notes
- Display download links

**Endpoint:** `https://raw.githubusercontent.com/pipe-cd/pipecd/master/docs/plugins.json`

### Documentation

The `plugins.md` file provides:
- Quick reference table for users
- Integration into PipeCD documentation
- Links to individual plugin documentation
- Version history context

### Continuous Integration

Use validation in CI pipelines:
```yaml
- name: Validate plugins registry
  run: python3 scripts/validate-plugins-registry.py
```

### Package Managers

The registry can integrate with package managers to:
- Verify latest plugin versions
- Automate plugin updates
- Check compatibility

## Data Format Specifications

### JSON Registry Format (v1.0)

```json
{
  "version": "1.0",
  "lastUpdated": "2026-01-24T12:00:00Z",
  "description": "...",
  "plugins": [
    {
      "id": "kubernetes",
      "name": "Kubernetes Plugin",
      "description": "...",
      "sourcePath": "pkg/app/pipedv1/plugin/kubernetes",
      "repository": "https://github.com/pipe-cd/pipecd",
      "repositoryType": "inline",
      "latestVersion": "v0.1.0",
      "releaseUrl": "https://github.com/pipe-cd/pipecd/releases/tag/...",
      "tagPattern": "pkg/app/pipedv1/plugin/kubernetes/*",
      "status": "stable"
    }
  ],
  "metadata": {...}
}
```

## Backward Compatibility

- **Version format**: Semver with optional 'v' prefix
- **Schema**: Can be extended with new optional fields
- **Tag patterns**: Flexible wildcard matching
- **Repository types**: Extensible for future plugin hosting locations

## Security Considerations

### Token Handling
- GitHub token passed only via environment variable `GITHUB_TOKEN`
- Never logged or displayed
- Workflow uses `secrets.GITHUB_TOKEN` (GitHub-provided)
- Higher rate limits with token (5000 req/hr vs 60 req/hr)

### Commits
- Signed commits recommended (if pipecd-bot has key configured)
- Minimal commits (only when versions change)
- Clear commit message for audit trail

### API Security
- Uses GitHub's official API
- HTTPS only for all requests
- Handles rate limits gracefully
- No sensitive data in JSON registry

## Troubleshooting

### Registry shows "unknown" for a plugin version

**Causes:**
1. Plugin hasn't been released yet
2. Tag pattern doesn't match any GitHub tags
3. GitHub API rate limit exceeded

**Solutions:**
1. Verify plugin has releases on GitHub
2. Check tag pattern matches actual tags
3. Run script with valid GitHub token

### Workflow not updating registry

**Causes:**
1. No changes to versions detected (working as designed)
2. GitHub Actions disabled on repository
3. Workflow has syntax errors

**Solutions:**
1. Check recent releases against registry
2. Verify workflow is enabled in repo settings
3. Review workflow logs in GitHub Actions tab

### Validation fails

**Causes:**
1. Manual edits corrupted JSON syntax
2. Missing required fields
3. Invalid URL formats

**Solutions:**
1. Run generator script to regenerate
2. Check schema validation errors
3. Fix any malformed entries

## Testing

### Local Testing

```bash
# Generate registry
python3 scripts/update-plugins-registry.py

# Validate generated registry
python3 scripts/validate-plugins-registry.py

# Manual inspection
cat docs/plugins.json | jq '.plugins[] | {id, latestVersion}'
```

### CI/CD Testing

The workflow includes automatic validation:
1. Syntax check
2. Schema validation
3. Semantic validation
4. No-change detection (prevents unnecessary commits)

## Future Enhancements

Potential improvements to the system:

1. **Plugin Compatibility Matrix**
   - Track which plugin versions work with which PipeCD versions
   - Add compatibility field to registry

2. **Release Notes Integration**
   - Parse release notes from GitHub
   - Include changelog snippets in registry

3. **Plugin Dependencies**
   - Track SDK version requirements
   - Highlight compatibility issues

4. **Plugin Search/Discovery**
   - Web UI for browsing plugins
   - Filter by status, category, platform

5. **Community Plugins**
   - Track community-developed plugins
   - Separate registry or combined with official plugins

6. **Mirrors/CDN**
   - Cache registry in CDN
   - Reduce GitHub API dependency

## Maintenance

### Regular Tasks

- **Weekly**: Monitor workflow execution for errors
- **Monthly**: Review plugin status classifications
- **Quarterly**: Update documentation as needed
- **Per Release**: Verify registry updates automatically

### Support

For issues or questions:
- Create GitHub issues in [pipe-cd/pipecd](https://github.com/pipe-cd/pipecd)
- Join PipeCD Slack: [#pipecd](https://app.slack.com/client/T08PSQ7BQ/C01B27F9T0X)
- Check workflow logs for specific errors

## References

- **PipeCD Repository**: https://github.com/pipe-cd/pipecd
- **Plugin Release Workflow**: `.github/workflows/plugin_release.yaml`
- **Plugin Architecture RFC**: `docs/rfcs/0015-pipecd-plugin-arch-meta.md`
- **PipeCD Documentation**: https://pipecd.dev/docs/
- **GitHub API Reference**: https://docs.github.com/en/rest

---

**Last updated:** 2026-01-24  
**Maintained by:** PipeCD Project  
**License:** Apache 2.0
