# PipeCD Scripts

Utility scripts for PipeCD development and maintenance.

## plugins-registry

Scripts for managing and updating the PipeCD plugins registry.

### update-plugins-registry.py

Automatically detects the latest versions of all official PipeCD plugins from GitHub releases.

**Features:**
- Fetches release information from GitHub API
- Matches releases against plugin-specific tag patterns
- Generates both JSON (machine-readable) and Markdown (human-readable) registries
- Supports both inline and external plugins
- Semantic version comparison for accurate latest version detection

**Usage:**
```bash
# Generate with default settings
python3 scripts/update-plugins-registry.py

# With GitHub token (higher rate limits)
GITHUB_TOKEN=<token> python3 scripts/update-plugins-registry.py

# Custom output directory
python3 scripts/update-plugins-registry.py --output-dir /path/to/output
```

**Output files:**
- `docs/plugins.json` - Structured JSON registry
- `docs/plugins.md` - Human-readable documentation

### validate-plugins-registry.py

Validates the plugins registry for correctness and consistency.

**Checks:**
- JSON schema validation against `docs/plugins.schema.json`
- Semantic checks (duplicate IDs, valid URLs, version formats)
- URL format validation
- Required field validation

**Usage:**
```bash
# Validate with default paths
python3 scripts/validate-plugins-registry.py

# Custom paths
python3 scripts/validate-plugins-registry.py \
  --registry docs/plugins.json \
  --schema docs/plugins.schema.json
```

**Exit codes:**
- `0` - All validations passed
- `1` - Validation failed

## Automation

These scripts are automatically run by GitHub Actions:

- **Update Workflow**: `.github/workflows/update-plugins-registry.yaml`
  - Runs on plugin releases
  - Runs on schedule (every 6 hours)
  - Commits changes when versions update

- **Validation**: Part of CI pipeline
  - Validates registry after updates
  - Prevents invalid data from being committed

## Dependencies

### update-plugins-registry.py
- Python 3.8+
- `requests` library

Install: `pip install requests`

### validate-plugins-registry.py
- Python 3.8+
- `jsonschema` library

Install: `pip install jsonschema`

## Development

### Running locally

1. **Set up environment**
   ```bash
   python3 -m venv venv
   source venv/bin/activate  # or `venv\Scripts\activate` on Windows
   pip install requests jsonschema
   ```

2. **Generate registry**
   ```bash
   python3 scripts/update-plugins-registry.py
   ```

3. **Validate registry**
   ```bash
   python3 scripts/validate-plugins-registry.py
   ```

4. **Check output**
   ```bash
   cat docs/plugins.json | jq .
   cat docs/plugins.md
   ```

### Adding new plugins

Edit `scripts/update-plugins-registry.py` and add entry to `PLUGINS_CONFIG`:

```python
{
    "id": "my-plugin",
    "name": "My Plugin",
    "description": "Plugin description",
    "sourcePath": "pkg/app/pipedv1/plugin/myplugin",
    "repository": "https://github.com/pipe-cd/pipecd",
    "repositoryType": "inline",
    "tagPattern": "pkg/app/pipedv1/plugin/myplugin/*",
    "status": "stable",
}
```

Then run the update script to fetch the latest version.

### Testing

```bash
# Test with mock GitHub API responses
python3 -m pytest scripts/test_registry_scripts.py

# Manual testing
python3 scripts/validate-plugins-registry.py
echo $?  # Check exit code
```

## Documentation

- **Full Documentation**: [`docs/PLUGINS_REGISTRY.md`](../PLUGINS_REGISTRY.md)
- **Registry Files**:
  - Machine-readable: [`docs/plugins.json`](../plugins.json)
  - Human-readable: [`docs/plugins.md`](../plugins.md)
  - Schema: [`docs/plugins.schema.json`](../plugins.schema.json)
- **Automation**: [`.github/workflows/update-plugins-registry.yaml`](../.github/workflows/update-plugins-registry.yaml)

## Troubleshooting

### Rate limit exceeded

**Problem**: `403: API rate limit exceeded`

**Solution**: Use GitHub token for higher rate limits
```bash
export GITHUB_TOKEN=<your_token>
python3 scripts/update-plugins-registry.py
```

### Module not found errors

**Problem**: `ModuleNotFoundError: No module named 'requests'`

**Solution**: Install dependencies
```bash
pip install requests jsonschema
```

### Invalid JSON output

**Problem**: Generated `plugins.json` fails validation

**Solution**: Check validation output
```bash
python3 scripts/validate-plugins-registry.py
```

Fix any errors and regenerate.

## Contributing

To improve these scripts:

1. Update script files directly
2. Add test cases to `scripts/test_registry_scripts.py` (if created)
3. Test locally before committing
4. Submit PR with changes

## License

Apache License 2.0 - See [LICENSE](../LICENSE) for details.
