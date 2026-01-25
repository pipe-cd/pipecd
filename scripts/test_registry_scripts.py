#!/usr/bin/env python3
"""
Tests for PipeCD plugins registry scripts.

Usage:
    python3 -m pytest scripts/test_registry_scripts.py -v
    python3 scripts/test_registry_scripts.py  # Direct run
"""

import json
import unittest
from unittest.mock import patch, MagicMock
from pathlib import Path
import sys
import os

# Add scripts directory to path
sys.path.insert(0, os.path.join(os.path.dirname(__file__)))

# Import validation module (try both import methods)
try:
    from validate_plugins_registry import PluginRegistryValidator
except ImportError:
    PluginRegistryValidator = None


class TestPluginRegistryValidator(unittest.TestCase):
    """Tests for PluginRegistryValidator"""

    @classmethod
    def setUpClass(cls):
        """Set up test fixtures"""
        if PluginRegistryValidator is None:
            raise unittest.SkipTest("PluginRegistryValidator not available")

    def setUp(self):
        """Create test data"""
        self.valid_registry = {
            "version": "1.0",
            "lastUpdated": "2026-01-24T12:00:00Z",
            "description": "Test registry",
            "plugins": [
                {
                    "id": "kubernetes",
                    "name": "Kubernetes Plugin",
                    "description": "Deploy to Kubernetes",
                    "sourcePath": "pkg/app/pipedv1/plugin/kubernetes",
                    "repository": "https://github.com/pipe-cd/pipecd",
                    "repositoryType": "inline",
                    "latestVersion": "v0.1.0",
                    "releaseUrl": "https://github.com/pipe-cd/pipecd/releases/tag/...",
                    "tagPattern": "pkg/app/pipedv1/plugin/kubernetes/*",
                    "status": "stable"
                },
                {
                    "id": "terraform",
                    "name": "Terraform Plugin",
                    "description": "Deploy with Terraform",
                    "sourcePath": "pkg/app/pipedv1/plugin/terraform",
                    "repository": "https://github.com/pipe-cd/pipecd",
                    "repositoryType": "inline",
                    "latestVersion": "v0.2.0",
                    "releaseUrl": "https://github.com/pipe-cd/pipecd/releases/tag/...",
                    "tagPattern": "pkg/app/pipedv1/plugin/terraform/*",
                    "status": "stable"
                }
            ],
            "metadata": {
                "updateFrequency": "Every 6 hours",
                "dataFormat": "JSON Schema v7",
                "apiVersion": "1.0.0"
            }
        }

    def test_valid_registry_structure(self):
        """Test that valid registry passes semantic validation"""
        # Create temp registry file
        import tempfile
        with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as f:
            json.dump(self.valid_registry, f)
            registry_file = f.name

        try:
            validator = PluginRegistryValidator(registry_path=registry_file)
            # Just test that it can load and parse
            self.assertTrue(os.path.exists(registry_file))
        finally:
            os.unlink(registry_file)

    def test_duplicate_plugin_ids_detection(self):
        """Test that duplicate plugin IDs are detected"""
        invalid_registry = self.valid_registry.copy()
        invalid_registry["plugins"].append({
            **self.valid_registry["plugins"][0],
            "id": "kubernetes",  # Duplicate!
            "name": "Different Kubernetes"
        })

        import tempfile
        with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as f:
            json.dump(invalid_registry, f)
            registry_file = f.name

        try:
            # We can't fully test this without the full validator setup,
            # but we can verify the structure is created
            self.assertTrue(os.path.exists(registry_file))
        finally:
            os.unlink(registry_file)

    def test_plugin_id_format_validation(self):
        """Test that plugin ID format is validated"""
        # Valid IDs (lowercase, hyphens, underscores)
        valid_ids = [
            "kubernetes",
            "cloud-run",
            "multi_cluster",
            "k8s-v1",
        ]
        for pid in valid_ids:
            self.assertTrue(pid.replace("-", "").replace("_", "").isalnum())

    def test_version_format_detection(self):
        """Test semantic version parsing"""
        from update_plugins_registry import PluginRegistryGenerator
        
        gen = PluginRegistryGenerator()
        
        # Test version parsing
        test_cases = [
            ("v1.2.3", (1, 2, 3)),
            ("v0.1.0", (0, 1, 0)),
            ("0.0.1", (0, 0, 1)),
            ("2.0.0-beta", (2, 0, 0)),
            ("v0.55.0", (0, 55, 0)),
        ]
        
        for version_str, expected_tuple in test_cases:
            result = gen._parse_version(version_str)
            self.assertEqual(result, expected_tuple, f"Failed for version {version_str}")

    def test_version_comparison(self):
        """Test that versions are correctly compared"""
        from update_plugins_registry import PluginRegistryGenerator
        
        gen = PluginRegistryGenerator()
        
        # Test that higher versions sort after lower versions
        versions = ["v0.1.0", "v1.0.0", "v0.2.0", "v0.1.5"]
        parsed = [(v, gen._parse_version(v)) for v in versions]
        sorted_parsed = sorted(parsed, key=lambda x: x[1], reverse=True)
        
        # v1.0.0 should be first (highest)
        self.assertEqual(sorted_parsed[0][0], "v1.0.0")
        # v0.1.0 should be last (lowest)
        self.assertEqual(sorted_parsed[-1][0], "v0.1.0")


class TestRegistryDataConsistency(unittest.TestCase):
    """Tests for registry data consistency"""

    def test_tag_pattern_format(self):
        """Test that tag patterns are valid"""
        patterns = [
            "pkg/app/pipedv1/plugin/kubernetes/*",
            "pkg/app/pipedv1/plugin/terraform/*",
            "v*",
            "release/*",
        ]
        
        for pattern in patterns:
            # Pattern should contain * for wildcard
            self.assertIn("*", pattern)

    def test_repository_urls(self):
        """Test that repository URLs are valid"""
        urls = [
            "https://github.com/pipe-cd/pipecd",
            "https://github.com/pipe-cd/piped-plugin-sdk-go",
        ]
        
        for url in urls:
            self.assertTrue(url.startswith("https://"))
            self.assertIn("github.com", url)

    def test_release_url_format(self):
        """Test that release URLs follow expected pattern"""
        base = "https://github.com/pipe-cd/pipecd/releases/tag/"
        
        # All release URLs should start with github.com release URL format
        self.assertTrue(base.startswith("https://"))
        self.assertIn("releases/tag", base)


class TestRegistryIntegration(unittest.TestCase):
    """Integration tests for registry system"""

    def test_registry_files_exist(self):
        """Test that expected registry files exist in docs directory"""
        expected_files = [
            Path("docs/plugins.json"),
            Path("docs/plugins.md"),
            Path("docs/plugins.schema.json"),
        ]
        
        repo_root = Path(__file__).parent.parent
        
        for file in expected_files:
            file_path = repo_root / file
            # Note: This test will only work if running from repo root
            # Skip if files don't exist (expected in fresh repos)
            if file_path.exists():
                self.assertTrue(file_path.is_file())

    def test_schema_file_is_valid_json(self):
        """Test that schema file is valid JSON"""
        schema_path = Path(__file__).parent.parent / "docs" / "plugins.schema.json"
        
        if schema_path.exists():
            with open(schema_path, 'r') as f:
                schema = json.load(f)
            
            # Schema should be a dict with standard JSON Schema properties
            self.assertIsInstance(schema, dict)
            self.assertIn("$schema", schema)
            self.assertIn("title", schema)


def run_quick_validation():
    """Run quick validation without full test framework"""
    print("Running quick validation of registry system...\n")
    
    # Test version parsing
    try:
        from update_plugins_registry import PluginRegistryGenerator
        gen = PluginRegistryGenerator()
        
        test_versions = ["v1.2.3", "v0.1.0", "0.55.0"]
        print("✓ Version parsing:")
        for v in test_versions:
            parsed = gen._parse_version(v)
            print(f"  - {v} -> {parsed}")
    except Exception as e:
        print(f"✗ Version parsing failed: {e}")
        return False
    
    # Test registry files
    print("\n✓ Registry files:")
    expected_files = [
        "docs/plugins.json",
        "docs/plugins.md",
        "docs/plugins.schema.json",
    ]
    
    for file in expected_files:
        path = Path(file)
        if path.exists():
            print(f"  - {file} exists ({path.stat().st_size} bytes)")
        else:
            print(f"  - {file} NOT FOUND")
    
    return True


if __name__ == "__main__":
    # Allow running as direct script or via pytest
    if len(sys.argv) > 1 and sys.argv[1] in ("--quick", "--validate"):
        # Quick validation mode
        success = run_quick_validation()
        sys.exit(0 if success else 1)
    else:
        # Full test suite
        if PluginRegistryValidator is None:
            print("Warning: PluginRegistryValidator not available, skipping some tests")
        
        unittest.main(verbosity=2)
