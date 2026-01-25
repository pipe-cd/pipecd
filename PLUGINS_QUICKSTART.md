# Quick Start Guide - PipeCD Plugins Registry

## For End Users

### Finding Plugin Versions

**Option 1: Human-Readable Documentation**
Visit `docs/plugins.md` for a formatted table showing all official plugins and their latest versions.

**Option 2: Machine-Readable JSON**
Access `docs/plugins.json` for structured data:
```bash
curl https://raw.githubusercontent.com/pipe-cd/pipecd/master/docs/plugins.json | jq '.plugins[] | {id, name, latestVersion}'
```

**Option 3: GitHub Releases**
As before, visit https://github.com/pipe-cd/pipecd/releases (unchanged)

---

## For Developers

### Running Registry Updates Locally

```bash
# Generate the registry
make gen/plugins-registry

# Validate the registry
make check/plugins-registry

# Or directly with Python
python3 scripts/update-plugins-registry.py
python3 scripts/validate-plugins-registry.py
```

### Adding a New Plugin

1. **Edit** `scripts/update-plugins-registry.py`
2. **Find** the `PLUGINS_CONFIG` list
3. **Add** a new entry:
   ```python
   {
       "id": "my-plugin",
       "name": "My Plugin",
       "description": "What it does",
       "sourcePath": "pkg/app/pipedv1/plugin/myplugin",
       "repository": "https://github.com/pipe-cd/pipecd",
       "repositoryType": "inline",
       "tagPattern": "pkg/app/pipedv1/plugin/myplugin/*",
       "status": "stable",
   }
   ```
4. **Save** and run `make gen/plugins-registry`

### Testing Changes

```bash
# Validate your registry changes
python3 scripts/validate-plugins-registry.py

# Run unit tests
python3 -m pytest scripts/test_registry_scripts.py -v

# Quick validation
python3 scripts/test_registry_scripts.py --quick
```

---

## For Maintainers

### Verifying Automation Works

The GitHub Actions workflow automatically:
- Updates on every release
- Updates every 6 hours
- Creates commits only when versions change

**To manually trigger:**
1. Go to `.github/workflows/update-plugins-registry.yaml`
2. Click "Run workflow" button in GitHub Actions
3. Wait for workflow to complete
4. Check generated files

### Monitoring

Check the workflow execution:
1. GitHub repo → Actions tab
2. Select "update-plugins-registry" workflow
3. Review recent runs
4. Check logs for any errors

### If Something Goes Wrong

**Scenario: Registry shows "unknown" version**
- Plugin may not have releases yet
- Check tag pattern matches actual GitHub tags
- Verify GitHub API is accessible

**Scenario: Workflow fails**
- Check workflow logs in GitHub Actions
- Verify Python dependencies are installed
- Check GitHub API rate limits

**Scenario: Manual fix needed**
```bash
# Manually regenerate and commit
make gen/plugins-registry
git add docs/plugins.json docs/plugins.md
git commit -m "fix: update plugins registry"
git push
```

---

## Files Reference

| File | Purpose | Audience |
|------|---------|----------|
| `docs/plugins.md` | Human-readable registry | Everyone |
| `docs/plugins.json` | Machine-readable API | Tools, websites |
| `docs/plugins.schema.json` | JSON schema validation | Developers |
| `docs/PLUGINS_REGISTRY.md` | Complete documentation | Developers, maintainers |
| `scripts/update-plugins-registry.py` | Update script | Developers, CI/CD |
| `scripts/validate-plugins-registry.py` | Validation script | Developers, CI/CD |
| `scripts/README.md` | Scripts documentation | Developers |
| `.github/workflows/update-plugins-registry.yaml` | GitHub Actions | Maintainers |

---

## Common Commands

```bash
# Generate registry
make gen/plugins-registry

# Validate registry
make check/plugins-registry

# Both together
make gen/plugins-registry && make check/plugins-registry

# View current plugins
cat docs/plugins.json | jq '.plugins[] | {id, latestVersion}'

# Check with token (better rate limits)
GITHUB_TOKEN=<token> python3 scripts/update-plugins-registry.py
```

---

## Integration Examples

### For Website

```javascript
// Fetch plugin versions
fetch('https://raw.githubusercontent.com/pipe-cd/pipecd/master/docs/plugins.json')
  .then(r => r.json())
  .then(registry => {
    registry.plugins.forEach(plugin => {
      console.log(`${plugin.name}: ${plugin.latestVersion}`);
    });
  });
```

### For Documentation

Link to `docs/plugins.md` in your documentation:
```markdown
[Official Plugins](docs/plugins.md)
```

### For CI/CD

```yaml
- name: Validate plugins registry
  run: python3 scripts/validate-plugins-registry.py
```

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| "Module not found" | `pip install requests jsonschema` |
| "Rate limit exceeded" | Use GitHub token: `GITHUB_TOKEN=<token>` |
| Registry shows "unknown" | Check if plugin has releases on GitHub |
| Validation fails | Run `make check/plugins-registry` to see details |
| Workflow not running | Check GitHub Actions settings and workflow status |

---

## More Information

- **Full Documentation:** [`docs/PLUGINS_REGISTRY.md`](docs/PLUGINS_REGISTRY.md)
- **Implementation Details:** [`IMPLEMENTATION_SUMMARY.md`](IMPLEMENTATION_SUMMARY.md)
- **Scripts Documentation:** [`scripts/README.md`](scripts/README.md)
- **Current Plugins:** [`docs/plugins.md`](docs/plugins.md)
- **Registry API:** [`docs/plugins.json`](docs/plugins.json)

---

**Last Updated:** 2026-01-24  
**Status:** ✅ Production Ready
