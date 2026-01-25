#!/usr/bin/env python3
"""
PipeCD Plugins Registry Generator

This script automatically detects and updates plugin versions from GitHub releases.
It generates both machine-readable (plugins.json) and human-readable (plugins.md) registry files.

Usage:
    python3 scripts/update-plugins-registry.py [--token TOKEN] [--output-dir docs]

Environment variables:
    GITHUB_TOKEN: GitHub API token (optional, for higher rate limits)
"""

import json
import os
import re
import sys
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Dict, List, Optional
from dataclasses import dataclass, asdict

import requests


@dataclass
class Plugin:
    """Plugin metadata"""
    id: str
    name: str
    description: str
    sourcePath: str
    repository: str
    repositoryType: str  # "inline" or "external"
    latestVersion: str
    releaseUrl: str
    tagPattern: str
    status: str  # "stable", "beta", "alpha"


class PluginRegistryGenerator:
    """Generate plugin registry from GitHub releases"""

    # Base configuration of known plugins
    PLUGINS_CONFIG = [
        {
            "id": "kubernetes",
            "name": "Kubernetes Plugin",
            "description": "Deploy applications to Kubernetes clusters",
            "sourcePath": "pkg/app/pipedv1/plugin/kubernetes",
            "repository": "https://github.com/pipe-cd/pipecd",
            "repositoryType": "inline",
            "tagPattern": "pkg/app/pipedv1/plugin/kubernetes/*",
            "status": "stable",
        },
        {
            "id": "terraform",
            "name": "Terraform Plugin",
            "description": "Deploy infrastructure using Terraform",
            "sourcePath": "pkg/app/pipedv1/plugin/terraform",
            "repository": "https://github.com/pipe-cd/pipecd",
            "repositoryType": "inline",
            "tagPattern": "pkg/app/pipedv1/plugin/terraform/*",
            "status": "stable",
        },
        {
            "id": "cloudrunservice",
            "name": "Cloud Run Plugin",
            "description": "Deploy services to Google Cloud Run",
            "sourcePath": "pkg/app/pipedv1/plugin/cloudrun",
            "repository": "https://github.com/pipe-cd/pipecd",
            "repositoryType": "inline",
            "tagPattern": "pkg/app/pipedv1/plugin/cloudrunservice/*",
            "status": "stable",
        },
        {
            "id": "wait",
            "name": "Wait Stage Plugin",
            "description": "Add delay/wait stages in deployment pipelines",
            "sourcePath": "pkg/app/pipedv1/plugin/wait",
            "repository": "https://github.com/pipe-cd/pipecd",
            "repositoryType": "inline",
            "tagPattern": "pkg/app/pipedv1/plugin/wait/*",
            "status": "stable",
        },
        {
            "id": "waitapproval",
            "name": "Wait Approval Stage Plugin",
            "description": "Add manual approval gates in deployment pipelines",
            "sourcePath": "pkg/app/pipedv1/plugin/waitapproval",
            "repository": "https://github.com/pipe-cd/pipecd",
            "repositoryType": "inline",
            "tagPattern": "pkg/app/pipedv1/plugin/waitapproval/*",
            "status": "stable",
        },
        {
            "id": "scriptrun",
            "name": "Script Run Plugin",
            "description": "Execute custom scripts as deployment stages",
            "sourcePath": "pkg/app/pipedv1/plugin/scriptrun",
            "repository": "https://github.com/pipe-cd/pipecd",
            "repositoryType": "inline",
            "tagPattern": "pkg/app/pipedv1/plugin/scriptrun/*",
            "status": "stable",
        },
        {
            "id": "analysis",
            "name": "Analysis Plugin",
            "description": "Analyze deployment metrics and logs",
            "sourcePath": "pkg/app/pipedv1/plugin/analysis",
            "repository": "https://github.com/pipe-cd/pipecd",
            "repositoryType": "inline",
            "tagPattern": "pkg/app/pipedv1/plugin/analysis/*",
            "status": "stable",
        },
        {
            "id": "kubernetes-multicluster",
            "name": "Kubernetes Multi-cluster Plugin",
            "description": "Deploy applications across multiple Kubernetes clusters",
            "sourcePath": "pkg/app/pipedv1/plugin/kubernetes_multicluster",
            "repository": "https://github.com/pipe-cd/pipecd",
            "repositoryType": "inline",
            "tagPattern": "pkg/app/pipedv1/plugin/kubernetes_multicluster/*",
            "status": "stable",
        },
        {
            "id": "piped-plugin-sdk-go",
            "name": "PipeCD Plugin SDK for Go",
            "description": "Official SDK for developing PipeCD plugins in Go",
            "sourcePath": "pkg/plugin/sdk",
            "repository": "https://github.com/pipe-cd/piped-plugin-sdk-go",
            "repositoryType": "external",
            "tagPattern": "v*",
            "status": "stable",
        },
    ]

    def __init__(self, token: Optional[str] = None, output_dir: str = "docs"):
        """Initialize registry generator
        
        Args:
            token: GitHub API token (optional, for higher rate limits)
            output_dir: Output directory for generated files
        """
        self.token = token or os.environ.get("GITHUB_TOKEN")
        self.output_dir = Path(output_dir)
        self.output_dir.mkdir(parents=True, exist_ok=True)
        self.session = requests.Session()
        if self.token:
            self.session.headers.update({"Authorization": f"token {self.token}"})

    def get_latest_version(self, owner: str, repo: str, tag_pattern: str) -> Optional[str]:
        """Get latest version for a plugin from GitHub releases
        
        Args:
            owner: Repository owner (e.g., "pipe-cd")
            repo: Repository name (e.g., "pipecd")
            tag_pattern: Tag pattern to match (e.g., "pkg/app/pipedv1/plugin/kubernetes/*")
        
        Returns:
            Latest version tag or None if not found
        """
        try:
            # Use GitHub API to list releases and tags
            # First try releases API (official releases)
            url = f"https://api.github.com/repos/{owner}/{repo}/releases"
            params = {"per_page": 100}
            
            while url:
                resp = self.session.get(url, params=params, timeout=10)
                resp.raise_for_status()
                releases = resp.json()
                
                # Convert tag pattern to regex (e.g., "pkg/app/pipedv1/plugin/kubernetes/*" -> regex)
                pattern_regex = tag_pattern.replace("*", ".*")
                pattern_regex = "^" + pattern_regex + "$"
                
                matching_versions = []
                for release in releases:
                    tag = release["tag_name"]
                    if re.match(pattern_regex, tag):
                        # Extract version number from tag
                        version = tag.split("/")[-1] if "/" in tag else tag
                        matching_versions.append((version, tag, release["published_at"]))
                
                if matching_versions:
                    # Sort by semantic version (descending) and return latest
                    matching_versions.sort(
                        key=lambda x: self._parse_version(x[0]), 
                        reverse=True
                    )
                    return matching_versions[0][1]
                
                # Check if there's a next page
                if "link" in resp.headers:
                    links = resp.headers["link"].split(",")
                    next_url = None
                    for link in links:
                        if 'rel="next"' in link:
                            next_url = link.split(";")[0].strip("<>")
                    url = next_url
                else:
                    break
            
            # Fallback: try tags API if releases API didn't find anything
            url = f"https://api.github.com/repos/{owner}/{repo}/tags"
            params = {"per_page": 100}
            
            while url:
                resp = self.session.get(url, params=params, timeout=10)
                resp.raise_for_status()
                tags = resp.json()
                
                pattern_regex = tag_pattern.replace("*", ".*")
                pattern_regex = "^" + pattern_regex + "$"
                
                matching_versions = []
                for tag_obj in tags:
                    tag = tag_obj["name"]
                    if re.match(pattern_regex, tag):
                        version = tag.split("/")[-1] if "/" in tag else tag
                        matching_versions.append((version, tag))
                
                if matching_versions:
                    matching_versions.sort(
                        key=lambda x: self._parse_version(x[0]), 
                        reverse=True
                    )
                    return matching_versions[0][1]
                
                # Check pagination
                if "link" in resp.headers:
                    links = resp.headers["link"].split(",")
                    next_url = None
                    for link in links:
                        if 'rel="next"' in link:
                            next_url = link.split(";")[0].strip("<>")
                    url = next_url
                else:
                    break
            
        except requests.RequestException as e:
            print(f"Error fetching versions for {owner}/{repo}: {e}", file=sys.stderr)
        
        return None

    @staticmethod
    def _parse_version(version: str) -> tuple:
        """Parse semantic version string to tuple for comparison
        
        Args:
            version: Version string (e.g., "v0.1.0", "0.1.0")
        
        Returns:
            Tuple of integers for version comparison
        """
        # Remove 'v' prefix if present
        version = version.lstrip("v")
        
        # Extract numeric parts
        parts = []
        for part in version.split("."):
            # Extract leading digits
            match = re.match(r"(\d+)", part)
            if match:
                parts.append(int(match.group(1)))
            else:
                parts.append(0)
        
        # Pad to 3 parts (major.minor.patch)
        while len(parts) < 3:
            parts.append(0)
        
        return tuple(parts[:3])

    def generate_registry(self) -> Dict[str, Any]:
        """Generate plugin registry with latest versions"""
        plugins = []
        
        for config in self.PLUGINS_CONFIG:
            # Parse repository URL to get owner and repo
            repo_url = config["repository"]
            match = re.match(r"https://github\.com/([^/]+)/([^/]+)", repo_url)
            if not match:
                print(f"Could not parse repository URL: {repo_url}", file=sys.stderr)
                continue
            
            owner, repo = match.groups()
            
            print(f"Fetching latest version for {config['id']}...", file=sys.stderr)
            latest_tag = self.get_latest_version(owner, repo, config["tagPattern"])
            
            if not latest_tag:
                print(f"Warning: Could not find version for {config['id']}", file=sys.stderr)
                latest_version = "unknown"
                release_url = config["repository"]
            else:
                latest_version = latest_tag.split("/")[-1] if "/" in latest_tag else latest_tag
                release_url = f"{repo_url}/releases/tag/{latest_tag}"
            
            plugin = Plugin(
                id=config["id"],
                name=config["name"],
                description=config["description"],
                sourcePath=config["sourcePath"],
                repository=config["repository"],
                repositoryType=config["repositoryType"],
                latestVersion=latest_version,
                releaseUrl=release_url,
                tagPattern=config["tagPattern"],
                status=config["status"],
            )
            plugins.append(plugin)
        
        registry = {
            "version": "1.0",
            "lastUpdated": datetime.now(timezone.utc).isoformat() + "Z",
            "description": "PipeCD Official Plugins Registry - automatically updated on releases",
            "plugins": [asdict(p) for p in plugins],
            "metadata": {
                "updateFrequency": "On every release, every 6 hours via scheduled workflow",
                "dataFormat": "JSON Schema v7",
                "apiVersion": "1.0.0",
                "notes": "This file is auto-generated. See .github/workflows/update-plugins-registry.yaml for automation details.",
            }
        }
        
        return registry

    def save_json_registry(self, registry: Dict[str, Any]) -> bool:
        """Save registry as JSON file
        
        Args:
            registry: Registry data
        
        Returns:
            True if file was written, False if no changes
        """
        output_file = self.output_dir / "plugins.json"
        new_content = json.dumps(registry, indent=2)
        
        # Check if file exists and content is the same
        if output_file.exists():
            with open(output_file, "r") as f:
                existing_content = f.read()
            if existing_content == new_content:
                print(f"No changes to {output_file}", file=sys.stderr)
                return False
        
        with open(output_file, "w") as f:
            f.write(new_content)
        
        print(f"Updated {output_file}", file=sys.stderr)
        return True

    def generate_markdown(self, registry: Dict[str, Any]) -> str:
        """Generate human-readable markdown registry
        
        Args:
            registry: Registry data
        
        Returns:
            Markdown content
        """
        plugins = registry["plugins"]
        last_updated = registry["lastUpdated"]
        
        # Parse timestamp for display
        dt = datetime.fromisoformat(last_updated.replace("Z", "+00:00"))
        formatted_date = dt.strftime("%Y-%m-%d")
        
        # Generate table
        table_rows = ["| Plugin | Latest Version | Repository | Documentation |"]
        table_rows.append("|--------|----------------|------------|---|")
        
        for plugin in plugins:
            doc_link = f"[Repo]({plugin['repository']})"
            if plugin["repositoryType"] == "inline":
                doc_link = f"[Docs](https://pipecd.dev/docs/)"
            
            table_rows.append(
                f"| [{plugin['name']}](#{plugin['id'].replace('-', '')}) | "
                f"{plugin['latestVersion']} | "
                f"[GitHub]({plugin['repository']}) | "
                f"{doc_link} |"
            )
        
        # Generate details section
        details = []
        for plugin in plugins:
            details.append(f"### {plugin['name']}\n")
            details.append(f"{plugin['description']}\n")
            details.append(f"- **Latest Version:** {plugin['latestVersion']}\n")
            details.append(f"- **Release URL:** {plugin['releaseUrl']}\n")
            details.append(f"- **Source:** [{plugin['sourcePath']}](../../{plugin['sourcePath']})\n")
            details.append(f"- **Status:** {plugin['status'].title()}\n")
        
        markdown = f"""# PipeCD Official Plugins

This document lists all official PipeCD plugins with their latest released versions.

**Last updated:** {formatted_date}

---

## Quick Reference

{chr(10).join(table_rows)}

---

## Plugin Details

{chr(10).join(details)}
---

## How Plugin Versions Are Tracked

- **Inline plugins** (in pipecd repo): Released as Git tags with format `pkg/app/pipedv1/plugin/{{name}}/v{{version}}`
  - Example: `pkg/app/pipedv1/plugin/kubernetes/v0.1.0`
  - Release page: https://github.com/pipe-cd/pipecd/releases

- **External plugins** (separate repos): Released in their own repositories
  - Example: `piped-plugin-sdk-go` uses tags `v{{version}}`

## Updating This Registry

This registry is **automatically updated** by a GitHub Actions workflow on:
- Every new plugin release
- Every 6 hours (scheduled)

Changes are committed to the repository when new versions are detected.

For details on the automation, see [.github/workflows/update-plugins-registry.yaml](../../.github/workflows/update-plugins-registry.yaml).

## Related Documentation

- **Plugin Architecture RFC:** [docs/rfcs/0015-pipecd-plugin-arch-meta.md](../rfcs/0015-pipecd-plugin-arch-meta.md)
- **Plugin Development Guide:** [PipeCD Documentation](https://pipecd.dev/docs/developer-guide/plugin-development/)
- **Plugin Release Process:** [.github/workflows/plugin_release.yaml](../../.github/workflows/plugin_release.yaml)
- **All Releases:** https://github.com/pipe-cd/pipecd/releases

---

**Note:** This registry is machine-generated. To update version information or add new plugins, modify the registry generation scripts or submit an issue/PR to the [PipeCD repository](https://github.com/pipe-cd/pipecd).
"""
        return markdown

    def save_markdown_registry(self, markdown: str) -> bool:
        """Save registry as Markdown file
        
        Args:
            markdown: Markdown content
        
        Returns:
            True if file was written, False if no changes
        """
        output_file = self.output_dir / "plugins.md"
        
        # Check if file exists and content is the same
        if output_file.exists():
            with open(output_file, "r") as f:
                existing_content = f.read()
            if existing_content == markdown:
                print(f"No changes to {output_file}", file=sys.stderr)
                return False
        
        with open(output_file, "w") as f:
            f.write(markdown)
        
        print(f"Updated {output_file}", file=sys.stderr)
        return True

    def run(self) -> bool:
        """Generate and save registry files
        
        Returns:
            True if any files were updated, False if no changes
        """
        print("Generating PipeCD Plugins Registry...", file=sys.stderr)
        
        registry = self.generate_registry()
        json_updated = self.save_json_registry(registry)
        markdown = self.generate_markdown(registry)
        md_updated = self.save_markdown_registry(markdown)
        
        return json_updated or md_updated


def main():
    """Main entry point"""
    token = os.environ.get("GITHUB_TOKEN")
    output_dir = os.environ.get("OUTPUT_DIR", "docs")
    
    # Parse command line arguments
    if "--token" in sys.argv:
        idx = sys.argv.index("--token")
        if idx + 1 < len(sys.argv):
            token = sys.argv[idx + 1]
    
    if "--output-dir" in sys.argv:
        idx = sys.argv.index("--output-dir")
        if idx + 1 < len(sys.argv):
            output_dir = sys.argv[idx + 1]
    
    generator = PluginRegistryGenerator(token=token, output_dir=output_dir)
    updated = generator.run()
    
    # Exit with code 0 (success) regardless, as the script always completes successfully
    sys.exit(0)


if __name__ == "__main__":
    main()
