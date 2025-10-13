#!/bin/bash

# Script to fetch the latest plugin versions from GitHub API
# This script helps track the latest released version of official plugins

set -e

REPO_OWNER="pipe-cd"
REPO_NAME="pipecd"
API_BASE="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to get latest version for a plugin
get_plugin_version() {
    local plugin_name="$1"
    local plugin_path="pkg/app/pipedv1/plugin/${plugin_name}"
    
    # Get releases that match the plugin path pattern
    local releases=$(curl -s "${API_BASE}/releases?per_page=50" | jq -r --arg path "$plugin_path" '.[] | select(.tag_name | startswith($path)) | .tag_name')
    
    if [ -z "$releases" ]; then
        echo "none"
        return
    fi
    
    # Extract version from tag_name (format: pkg/app/pipedv1/plugin/pluginname/vX.Y.Z)
    local latest_version=$(echo "$releases" | head -1 | sed "s|${plugin_path}/||")
    echo "$latest_version"
}

# Function to get release date
get_release_date() {
    local tag_name="$1"
    curl -s "${API_BASE}/releases/tags/${tag_name}" | jq -r '.published_at' | cut -d'T' -f1
}

echo -e "${BLUE}PipeCD Official Plugins - Latest Versions${NC}"
echo "=============================================="
echo

# List of official plugins
plugins=(
    "kubernetes"
    "kubernetes_multicluster"
    "terraform"
    "cloudrun"
    "analysis"
    "scriptrun"
    "wait"
    "waitapproval"
)

# Header
printf "%-25s %-15s %-15s %s\n" "Plugin" "Version" "Release Date" "GitHub Release"
echo "-----------------------------------------------------------------------------------------"

for plugin in "${plugins[@]}"; do
    version=$(get_plugin_version "$plugin")
    
    if [ "$version" != "none" ]; then
        tag_name="pkg/app/pipedv1/plugin/${plugin}/${version}"
        release_date=$(get_release_date "$tag_name")
        github_link="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/tag/${tag_name}"
        
        printf "%-25s ${GREEN}%-15s${NC} %-15s %s\n" \
            "$plugin" \
            "$version" \
            "$release_date" \
            "$github_link"
    else
        printf "%-25s ${RED}%-15s${NC} %-15s %s\n" \
            "$plugin" \
            "No releases" \
            "-" \
            "-"
    fi
done

echo
echo -e "${YELLOW}Note:${NC} All plugins are currently released together with PipeCD core."
echo -e "For the complete list of releases, visit: ${BLUE}https://github.com/${REPO_OWNER}/${REPO_NAME}/releases${NC}"