#!/usr/bin/env bash

# Copyright 2025 The PipeCD Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Purpose: Harmonise '$GITHUB_WORKFLOWS_PATH' and the reversemap '$REVERSEMAP_FILE'
# so that every workflow references GitHub Actions by commit hash (not tag).
#
# Usage: see help function
# Working directory must be the project root.

set -e

GITHUB_WORKFLOWS_PATH="./.github/workflows"
REVERSEMAP_FILE=".gha-reversemap.yml"
YQ_MIN_VERSION="yq (https://github.com/mikefarah/yq/) version v4.45"
GIT_COMMITSHA_LENGTH=40
TMP_OUTPUT="/tmp/$(date -u -Iseconds | cut -d '+' -f1).json"

ERR_LIMITED=47
ERR_YQ_NOT_INSTALLED=60
ERR_BAD_VERSION=72
ERR_FETCH_TAG_FAIL=74
ERR_FETCH_BRANCH_FAIL=76
ERR_NO_SHA=79
ERR_NO_ACTION=84
ERR_NO_LATEST=86
ERR_NO_TAG=87
ERR_VERIFY_FAIL=90

if [[ -n "${GITHUB_TOKEN:-}" ]]; then
    github_auth=("-H" "Authorization: Bearer $GITHUB_TOKEN")
else
    github_auth=()
fi

help() {
    cat <<EOF
Harmonise '$GITHUB_WORKFLOWS_PATH' and the reversemap '$REVERSEMAP_FILE'

Usage:
    $0  <operation> [ARGUMENTS]

Operations:
    verify-mapusage       [WORKFLOW_FILE...]   Check that every workflow references
                                              GitHub Actions by commit hash and that
                                              the hash is the one in the reversemap.
                                              If no files given, all workflows are checked.

    apply-reversemap      [WORKFLOW_FILE...]  Update workflow files so every action
                                              uses the commit hash from the reversemap.
                                              If no files given, all workflows are updated.

    update-action-version ACTION_REF...        Update the reversemap entry for the given
                                              action(s) to their latest release tag.

    update-reversemap      [WORKFLOW_FILE...]  Update the reversemap with sha/tag/urls
                                              from the actions used in the given workflows.

To fix failing CI: run '$0 apply-reversemap' then commit the changes.

EOF
}

_loginfo() {
    echo -e "$(date -Iseconds);INFO;$1"
}

_exit_with_error() {
    echo -e "$(date -Iseconds);ERROR;$2" >&2
    exit "$1"
}

_return() {
    echo "$1"
}

_check_yq_version() {
    if command -v yq >/dev/null 2>&1; then
        INSTALLED_VERSION=$(yq --version 2>/dev/null || true)
        if ! [[ "$INSTALLED_VERSION" > "$YQ_MIN_VERSION" ]]; then
            _exit_with_error $ERR_YQ_NOT_INSTALLED "yq is required (at least $YQ_MIN_VERSION). Install from https://github.com/mikefarah/yq."
        fi
    else
        _exit_with_error $ERR_YQ_NOT_INSTALLED "yq is required. Install from https://github.com/mikefarah/yq."
    fi
}

_fetch_sha_from_upstream_ref() {
    action_ref=$1
    tag_or_branch=$2
    action_ref_safe=$(echo "$action_ref" | cut -d '/' -f 1,2)
    API_GITHUB_TAG="https://api.github.com/repos/${action_ref_safe}/git/refs/tags/${tag_or_branch}"
    API_GITHUB_BRANCH="https://api.github.com/repos/${action_ref_safe}/git/refs/heads/${tag_or_branch}"
    HTTP_STATUS=$(curl -o "$TMP_OUTPUT" -s -w "%{http_code}" "${github_auth[@]}" "$API_GITHUB_TAG")
    case "$HTTP_STATUS" in
        2??) commit_sha=$(jq -r '.object.sha' "$TMP_OUTPUT") ;;
        403|429) _exit_with_error $ERR_LIMITED "GitHub API rate limit. Set GITHUB_TOKEN and retry." ;;
        404) HTTP_STATUS=$(curl -o "$TMP_OUTPUT" -s -w "%{http_code}" "${github_auth[@]}" "$API_GITHUB_BRANCH")
            case "$HTTP_STATUS" in
                2??) commit_sha=$(jq -r '.object.sha' "$TMP_OUTPUT") ;;
                403|429) _exit_with_error $ERR_LIMITED "GitHub API rate limit. Set GITHUB_TOKEN and retry." ;;
                404) _exit_with_error $ERR_BAD_VERSION "${tag_or_branch} is neither a tag nor a branch of ${action_ref_safe}" ;;
                *) _exit_with_error $ERR_FETCH_BRANCH_FAIL "GitHub rejected GET $API_GITHUB_BRANCH with status $HTTP_STATUS" ;;
            esac ;;
        *) _exit_with_error $ERR_FETCH_TAG_FAIL "GitHub rejected GET $API_GITHUB_TAG with status $HTTP_STATUS" ;;
    esac
    _return "$commit_sha"
}

_yq_update_reversemap() {
    action_ref=$1
    action_tag=$2
    action_sha=$3
    action_ref_safe=$(echo "$action_ref" | cut -d '/' -f 1,2)
    yq ".\"${action_ref}\".sha = \"${action_sha}\"" -i "$REVERSEMAP_FILE"
    yq ".\"${action_ref}\".sha-url = \"https://github.com/${action_ref_safe}/commit/${action_sha}\"" -i "$REVERSEMAP_FILE"
    yq ".\"${action_ref}\".tag = \"${action_tag}\"" -i "$REVERSEMAP_FILE"
    yq ".\"${action_ref}\".tag-url = \"https://github.com/${action_ref_safe}/tree/${action_tag}\"" -i "$REVERSEMAP_FILE"
}

_update_reversemap_with() {
    local filename=$1
    local action_fullref action_ref action_tag length
    for action_fullref in $(yq '.jobs[].steps[] | select(has("uses")) | .uses' "$filename" 2>/dev/null); do
        [[ "$action_fullref" == docker://* ]] && continue
        [[ "$action_fullref" != *@* ]] && continue
        action_ref=$(echo "$action_fullref" | cut -d '@' -f 1)
        action_tag=$(echo "$action_fullref" | cut -d '@' -f 2 | sed 's/[[:space:]]*#.*//' | awk '{print $1}')
        length=${#action_tag}
        _loginfo "ref=$action_ref version=$action_tag len=$length"
        if [[ $length -ne $GIT_COMMITSHA_LENGTH ]]; then
            action_sha=$(_fetch_sha_from_upstream_ref "$action_ref" "$action_tag")
            _loginfo "action=$action_ref tag=$action_tag sha=$action_sha"
            _yq_update_reversemap "$action_ref" "$action_tag" "$action_sha"
        fi
    done
}

_fetch_latest_tag() {
    local action_ref=$1
    local action_ref_safe latest latest_json
    action_ref_safe=$(echo "$action_ref" | cut -d '/' -f 1,2)
    HTTP_STATUS=$(curl -o "$TMP_OUTPUT" -s -w "%{http_code}" "${github_auth[@]}" "https://api.github.com/repos/${action_ref_safe}/releases/latest")
    if [[ "$HTTP_STATUS" -ge 200 && "$HTTP_STATUS" -lt 300 ]]; then
        latest_json=$(cat "$TMP_OUTPUT")
    elif [[ "$HTTP_STATUS" == 403 || "$HTTP_STATUS" == 429 ]]; then
        _exit_with_error $ERR_LIMITED "GitHub API rate limit. Set GITHUB_TOKEN and retry."
    elif [[ "$HTTP_STATUS" == 404 ]]; then
        _exit_with_error $ERR_NO_ACTION "No action named '$action_ref'"
    else
        _exit_with_error $ERR_NO_LATEST "GitHub API returned $HTTP_STATUS"
    fi
    latest=$(jq -r .tag_name <<<"$latest_json")
    [[ -z "$latest" || "$latest" = null ]] && _exit_with_error $ERR_NO_TAG "No tag_name in response for $action_ref"
    _return "$latest"
}

_update_action_version_infile() {
    local file=$1 action_ref=$2 action_sha=$3
    local sed_ref
    sed_ref=$(echo "$action_ref" | sed 's/[\/&]/\\&/g')
    if [[ "$(uname -s)" = "Darwin" ]]; then
        sed -E -i '' "s;(uses:) ${sed_ref}@[^[:space:]]+[^[:space:]]*;\1 ${action_ref}@${action_sha};g" "$file"
    else
        sed -E -i "s;(uses:) ${sed_ref}@[^[:space:]]+[^[:space:]]*;\1 ${action_ref}@${action_sha};g" "$file"
    fi
}

_get_sha_from_reversemap() {
    local action_ref=$1 query
    query=$(yq ".\"${action_ref}\".sha" "$REVERSEMAP_FILE" 2>/dev/null)
    [[ -z "$query" || "$query" = null ]] && _exit_with_error $ERR_NO_SHA "No sha for $action_ref in $REVERSEMAP_FILE"
    _return "$query"
}

run_verify_mapusage() {
    local files file ref action version goodsha
    files=("$@")
    if [[ ${#files[@]} -eq 0 ]]; then
        files=("${GITHUB_WORKFLOWS_PATH}"/*.yaml "${GITHUB_WORKFLOWS_PATH}"/*.yml)
    fi
    local failed=false
    local shadict
    shadict=$(yq -o json 'map_values(.sha)' "$REVERSEMAP_FILE")
    for file in "${files[@]}"; do
        [[ -f "$file" ]] || continue
        for ref in $(yq '.jobs[].steps[].uses?' "$file" 2>/dev/null); do
            [[ "$ref" == null || -z "$ref" ]] && continue
            [[ "$ref" == docker://* ]] && continue
            [[ "$ref" != *@* ]] && continue
            action=$(echo "$ref" | cut -d'@' -f1)
            version=$(echo "$ref" | cut -d'@' -f2 | sed 's/[[:space:]]*#.*//' | awk '{print $1}')
            if ! [[ "$version" =~ ^[0-9a-f]{40}$ ]]; then
                _loginfo "$file uses $ref (version '$version' is not a 40-char commit hash)"
                failed=true
                continue
            fi
            goodsha=$(jq -r --arg action "$action" '.[$action] // empty' <<<"$shadict")
            if [[ -z "$goodsha" ]]; then
                _loginfo "$file uses $ref but reversemap has no entry for $action"
                failed=true
            elif [[ "$version" != "$goodsha" ]]; then
                _loginfo "$file uses $ref (hash $version does not match reversemap $goodsha)"
                failed=true
            fi
        done
    done
    if [[ "$failed" = true ]]; then
        _exit_with_error $ERR_VERIFY_FAIL "Workflows must reference actions by the commit hash in $REVERSEMAP_FILE. Run: $0 apply-reversemap"
    fi
}

run_apply_reversemap() {
    local files file action_fullref action_ref action_sha
    files=("$@")
    if [[ ${#files[@]} -eq 0 ]]; then
        files=("${GITHUB_WORKFLOWS_PATH}"/*.yaml "${GITHUB_WORKFLOWS_PATH}"/*.yml)
    fi
    for file in "${files[@]}"; do
        [[ -f "$file" ]] || continue
        _loginfo "applying $REVERSEMAP_FILE to $file"
        for action_fullref in $(yq '.jobs[].steps[] | select(.uses) | .uses' "$file" 2>/dev/null); do
            [[ "$action_fullref" == docker://* ]] && continue
            [[ "$action_fullref" != *@* ]] && continue
            action_ref=$(echo "$action_fullref" | cut -d'@' -f1)
            action_sha=$(_get_sha_from_reversemap "$action_ref")
            _loginfo "$action_ref -> $action_sha"
            _update_action_version_infile "$file" "$action_ref" "$action_sha"
        done
    done
}

run_update_action_version() {
    local action_refs=("$@") action_ref latest_tag action_sha
    [[ ${#action_refs[@]} -eq 0 ]] && _exit_with_error 1 "Usage: $0 update-action-version OWNER/REPO [OWNER/REPO...]"
    for action_ref in "${action_refs[@]}"; do
        _loginfo "updating $action_ref to latest release in $REVERSEMAP_FILE"
        latest_tag=$(_fetch_latest_tag "$action_ref")
        action_sha=$(_fetch_sha_from_upstream_ref "$action_ref" "$latest_tag")
        _yq_update_reversemap "$action_ref" "$latest_tag" "$action_sha"
    done
}

run_update_reversemap() {
    local files=("$@") file
    if [[ ${#files[@]} -eq 0 ]]; then
        files=("${GITHUB_WORKFLOWS_PATH}"/*.yaml "${GITHUB_WORKFLOWS_PATH}"/*.yml)
    fi
    for file in "${files[@]}"; do
        [[ -f "$file" ]] || continue
        _loginfo "updating $REVERSEMAP_FILE from $file"
        _update_reversemap_with "$file"
    done
}

run_cli() {
    local op=${1:-}
    shift || true
    case "$op" in
        help|--help|-h)
            help
            ;;
        verify-mapusage)
            _check_yq_version
            _loginfo "verifying all workflows use commit hashes from $REVERSEMAP_FILE"
            run_verify_mapusage "$@"
            ;;
        apply-reversemap)
            _check_yq_version
            run_apply_reversemap "$@"
            ;;
        update-action-version)
            _check_yq_version
            run_update_action_version "$@"
            ;;
        update-reversemap)
            _check_yq_version
            run_update_reversemap "$@"
            ;;
        *)
            help
            exit 1
            ;;
    esac
    rm -f "$TMP_OUTPUT"
}

run_cli "$@"
