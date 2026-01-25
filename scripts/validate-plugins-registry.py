#!/usr/bin/env python3
"""
PipeCD Plugins Registry Validator

This script validates the plugins registry against the JSON schema and performs
semantic checks on the registry data.

Usage:
    python3 scripts/validate-plugins-registry.py [--registry docs/plugins.json]
"""

import json
import sys
from pathlib import Path
from typing import List, Tuple
from urllib.parse import urlparse

try:
    import jsonschema
except ImportError:
    print("Error: jsonschema module not found. Install with: pip install jsonschema", file=sys.stderr)
    sys.exit(1)


class PluginRegistryValidator:
    """Validate plugin registry"""

    def __init__(self, registry_path: str = "docs/plugins.json", schema_path: str = "docs/plugins.schema.json"):
        """Initialize validator
        
        Args:
            registry_path: Path to plugins.json
            schema_path: Path to plugins.schema.json
        """
        self.registry_path = Path(registry_path)
        self.schema_path = Path(schema_path)
        self.errors: List[str] = []
        self.warnings: List[str] = []

    def load_json(self, path: Path) -> dict:
        """Load and parse JSON file
        
        Args:
            path: Path to JSON file
        
        Returns:
            Parsed JSON data
        
        Raises:
            FileNotFoundError: If file doesn't exist
            json.JSONDecodeError: If JSON is invalid
        """
        with open(path, "r") as f:
            return json.load(f)

    def validate_schema(self) -> bool:
        """Validate registry against JSON schema
        
        Returns:
            True if valid, False otherwise
        """
        try:
            registry = self.load_json(self.registry_path)
            schema = self.load_json(self.schema_path)
            
            validator = jsonschema.Draft7Validator(schema)
            
            for error in validator.iter_errors(registry):
                path = ".".join(str(p) for p in error.absolute_path)
                self.errors.append(f"Schema validation error at {path}: {error.message}")
            
            return len(self.errors) == 0
        except (FileNotFoundError, json.JSONDecodeError) as e:
            self.errors.append(f"Failed to load files: {e}")
            return False

    def validate_semantic(self) -> bool:
        """Perform semantic checks on registry
        
        Returns:
            True if valid, False otherwise
        """
        try:
            registry = self.load_json(self.registry_path)
        except (FileNotFoundError, json.JSONDecodeError) as e:
            self.errors.append(f"Failed to load registry: {e}")
            return False
        
        valid = True
        plugins = registry.get("plugins", [])
        seen_ids = set()
        
        for i, plugin in enumerate(plugins):
            plugin_id = plugin.get("id")
            
            # Check for duplicate IDs
            if plugin_id in seen_ids:
                self.errors.append(f"Plugin {i}: Duplicate plugin ID '{plugin_id}'")
                valid = False
            seen_ids.add(plugin_id)
            
            # Validate plugin ID format
            if not plugin_id or not plugin_id.replace("-", "").replace("_", "").isalnum():
                self.errors.append(f"Plugin {i}: Invalid plugin ID format '{plugin_id}'")
                valid = False
            
            # Validate repository URL format
            repo_url = plugin.get("repository", "")
            try:
                result = urlparse(repo_url)
                if not all([result.scheme, result.netloc]):
                    self.errors.append(f"Plugin {plugin_id}: Invalid repository URL '{repo_url}'")
                    valid = False
            except Exception as e:
                self.errors.append(f"Plugin {plugin_id}: Failed to parse repository URL: {e}")
                valid = False
            
            # Validate release URL format
            release_url = plugin.get("releaseUrl", "")
            try:
                result = urlparse(release_url)
                if not all([result.scheme, result.netloc]):
                    self.errors.append(f"Plugin {plugin_id}: Invalid release URL '{release_url}'")
                    valid = False
            except Exception as e:
                self.errors.append(f"Plugin {plugin_id}: Failed to parse release URL: {e}")
                valid = False
            
            # Check if version looks valid (basic check)
            version = plugin.get("latestVersion", "")
            if version == "unknown":
                self.warnings.append(f"Plugin {plugin_id}: Version is unknown (not yet released?)")
            elif not version.startswith(("v", "pkg/")):
                # Allow both v-prefixed and path-prefixed tags
                if not any(c.isdigit() for c in version):
                    self.warnings.append(f"Plugin {plugin_id}: Version '{version}' doesn't look like a semantic version")
            
            # Check source path format
            source_path = plugin.get("sourcePath", "")
            if not source_path:
                self.errors.append(f"Plugin {plugin_id}: Missing sourcePath")
                valid = False
            elif not source_path.startswith(("pkg/", ".")):
                self.warnings.append(f"Plugin {plugin_id}: sourcePath '{source_path}' has unexpected format")
            
            # Validate repository type
            repo_type = plugin.get("repositoryType", "")
            if repo_type not in ("inline", "external"):
                self.errors.append(f"Plugin {plugin_id}: Invalid repositoryType '{repo_type}'")
                valid = False
            
            # Validate status
            status = plugin.get("status", "")
            if status not in ("stable", "beta", "alpha", "deprecated"):
                self.errors.append(f"Plugin {plugin_id}: Invalid status '{status}'")
                valid = False
        
        # Check for minimum number of plugins
        if len(plugins) < 3:
            self.warnings.append(f"Registry contains only {len(plugins)} plugins (expected at least 3)")
        
        return valid

    def run(self) -> bool:
        """Run all validations
        
        Returns:
            True if all validations pass, False otherwise
        """
        print("Validating PipeCD Plugins Registry...", file=sys.stderr)
        
        # Schema validation
        print("  Checking JSON schema...", file=sys.stderr)
        schema_valid = self.validate_schema()
        
        # Semantic validation
        print("  Checking semantic validity...", file=sys.stderr)
        semantic_valid = self.validate_semantic()
        
        # Report results
        if self.errors:
            print("\n❌ ERRORS:", file=sys.stderr)
            for error in self.errors:
                print(f"  - {error}", file=sys.stderr)
        
        if self.warnings:
            print("\n⚠️  WARNINGS:", file=sys.stderr)
            for warning in self.warnings:
                print(f"  - {warning}", file=sys.stderr)
        
        if not self.errors:
            print("\n✓ Registry is valid", file=sys.stderr)
        
        return schema_valid and semantic_valid


def main():
    """Main entry point"""
    registry_path = "docs/plugins.json"
    schema_path = "docs/plugins.schema.json"
    
    # Parse command line arguments
    if "--registry" in sys.argv:
        idx = sys.argv.index("--registry")
        if idx + 1 < len(sys.argv):
            registry_path = sys.argv[idx + 1]
    
    if "--schema" in sys.argv:
        idx = sys.argv.index("--schema")
        if idx + 1 < len(sys.argv):
            schema_path = sys.argv[idx + 1]
    
    validator = PluginRegistryValidator(registry_path=registry_path, schema_path=schema_path)
    valid = validator.run()
    
    sys.exit(0 if valid else 1)


if __name__ == "__main__":
    main()
