#!/usr/bin/env python3
"""
PipeCD Contributor Ladder Generator

Queries the GitHub API to tally merged pull requests per author,
categorizes contributors into defined ladder tiers, and generates
the CONTRIBUTORS.md file.

Configurable constants at the top for easy maintenance and tuning.
"""

import argparse
import datetime
import os
import subprocess
import sys
import time
from typing import Any, Dict, List, Optional, Tuple

import requests

# ==============================================================================
# CONFIGURATION CONSTANTS (Easily tunable)
# ==============================================================================

# Default repository to query
DEFAULT_REPO = os.environ.get("GITHUB_REPOSITORY", "pipe-cd/pipecd")

# Tier thresholds based on merged PR count
# Defined in descending order: (Name, Badge, Min PRs, Max PRs, Description)
# Provisional thresholds as agreed in issue #6548:
#   🌱 Newcomer: 1 merged PR
#   🛠️ Contributor: 2–4 merged PRs
#   🚀 Core Contributor: 5+ merged PRs
TIERS = [
    {
        "name": "Core Contributor",
        "badge": "🚀",
        "min_prs": 5,
        "max_prs": None,
        "description": "5+ merged PRs",
    },
    {
        "name": "Contributor",
        "badge": "🛠️",
        "min_prs": 2,
        "max_prs": 4,
        "description": "2–4 merged PRs",
    },
    {
        "name": "Newcomer",
        "badge": "🌱",
        "min_prs": 1,
        "max_prs": 1,
        "description": "1 merged PR",
    },
]

# Bot and system accounts to exclude from the contributor ladder
EXCLUDED_USERS = {
    "dependabot[bot]",
    "github-actions[bot]",
    "pipe-cd-bot",
    "pipe-cd-bot[bot]",
    "codecov[bot]",
    "imgbot[bot]",
    "fossabot",
    "ghost",
}


def is_bot_user(login: str, user_type: Optional[str] = None) -> bool:
    """Check if the user is a known bot or service account."""
    clean_login = login.lower()
    if clean_login in {u.lower() for u in EXCLUDED_USERS}:
        return True
    if clean_login.endswith("[bot]"):
        return True
    if user_type == "Bot":
        return True
    return False


def get_auth_token(cli_token: Optional[str] = None) -> Optional[str]:
    """Retrieve GitHub token from CLI arg, environment, or GitHub CLI."""
    if cli_token:
        return cli_token

    for env_var in ("GITHUB_TOKEN", "GH_TOKEN"):
        token = os.environ.get(env_var)
        if token:
            return token.strip()

    # Fallback to `gh auth token` if available locally
    try:
        output = subprocess.check_output(
            ["gh", "auth", "token"], stderr=subprocess.DEVNULL
        )
        token = output.decode("utf-8").strip()
        if token:
            return token
    except Exception:
        pass

    return None


def fetch_merged_prs_via_pulls(
    repo: str, headers: Dict[str, str], session: requests.Session
) -> Tuple[Dict[str, Dict[str, Any]], int]:
    """
    Fetch all closed pull requests via /repos/{owner}/{repo}/pulls,
    filtering for merged_at != None.
    Handles unlimited pagination without GitHub Search API's 1,000-result cap.
    """
    url = f"https://api.github.com/repos/{repo}/pulls"
    params = {
        "state": "closed",
        "per_page": 100,
        "page": 1,
    }

    contributors: Dict[str, Dict[str, Any]] = {}
    total_merged = 0

    print(f"Fetching pull requests for {repo} via pulls API...")
    while True:
        resp = session.get(url, headers=headers, params=params)

        # Handle rate limiting
        if resp.status_code == 403 and "rate limit" in resp.text.lower():
            reset_time = int(resp.headers.get("X-RateLimit-Reset", time.time() + 60))
            sleep_duration = max(reset_time - int(time.time()), 1)
            print(f"Rate limit exceeded. Waiting {sleep_duration}s for reset...")
            time.sleep(sleep_duration)
            continue

        resp.raise_for_status()
        prs = resp.json()
        if not prs:
            break

        for pr in prs:
            if not pr.get("merged_at"):
                continue

            total_merged += 1
            user = pr.get("user")
            if not user:
                continue

            login = user.get("login", "")
            user_type = user.get("type", "")
            if is_bot_user(login, user_type):
                continue

            if login not in contributors:
                contributors[login] = {
                    "login": login,
                    "html_url": user.get("html_url", f"https://github.com/{login}"),
                    "count": 0,
                }
            contributors[login]["count"] += 1

        link_header = resp.headers.get("Link", "")
        if 'rel="next"' not in link_header:
            break

        current_page = params["page"]
        if current_page % 10 == 0:
            print(f"  Processed {current_page} pages ({total_merged} merged PRs found so far)...")
        params["page"] += 1

    return contributors, total_merged


def fetch_merged_prs_via_search(
    repo: str, headers: Dict[str, str], session: requests.Session
) -> Tuple[Dict[str, Dict[str, Any]], int]:
    """
    Fetch merged pull requests via GitHub Search API (/search/issues).
    Note: GitHub Search API limits total results to 1,000 items per query.
    If total_count > 1000, prompts switching to the pulls API.
    """
    url = "https://api.github.com/search/issues"
    query = f"repo:{repo} is:pr is:merged"
    params = {
        "q": query,
        "per_page": 100,
        "page": 1,
    }

    contributors: Dict[str, Dict[str, Any]] = {}
    total_merged = 0

    print(f"Fetching merged pull requests for {repo} via Search API...")
    while True:
        resp = session.get(url, headers=headers, params=params)

        if resp.status_code == 403 and "rate limit" in resp.text.lower():
            reset_time = int(resp.headers.get("X-RateLimit-Reset", time.time() + 60))
            sleep_duration = max(reset_time - int(time.time()), 1)
            print(f"Search API rate limit exceeded. Waiting {sleep_duration}s...")
            time.sleep(sleep_duration)
            continue

        if resp.status_code == 422:
            print("Notice: GitHub Search API limit (1,000 results) reached.")
            break

        resp.raise_for_status()
        data = resp.json()
        total_count = data.get("total_count", 0)

        items = data.get("items", [])
        if not items:
            break

        for item in items:
            total_merged += 1
            user = item.get("user")
            if not user:
                continue

            login = user.get("login", "")
            user_type = user.get("type", "")
            if is_bot_user(login, user_type):
                continue

            if login not in contributors:
                contributors[login] = {
                    "login": login,
                    "html_url": user.get("html_url", f"https://github.com/{login}"),
                    "count": 0,
                }
            contributors[login]["count"] += 1

        if params["page"] * params["per_page"] >= total_count:
            break

        params["page"] += 1
        if params["page"] > 10:  # Search API max 1000 items
            print(f"Notice: Reached 1,000 items cap in Search API (total_count in repo: {total_count}).")
            break

    return contributors, total_merged


def assign_tiers(
    contributors: Dict[str, Dict[str, Any]]
) -> Dict[str, List[Dict[str, Any]]]:
    """
    Assign contributors into tiers and sort within each tier by:
      1. PR count descending
      2. Username alphabetically (case-insensitive) ascending as tiebreak.
    """
    tiered: Dict[str, List[Dict[str, Any]]] = {tier["name"]: [] for tier in TIERS}

    for user_info in contributors.values():
        count = user_info["count"]
        for tier in TIERS:
            min_prs = tier["min_prs"]
            max_prs = tier["max_prs"]
            if min_prs is not None and count < min_prs:
                continue
            if max_prs is not None and count > max_prs:
                continue
            tiered[tier["name"]].append(user_info)
            break

    # Sort each tier: count DESC, login ASC
    for tier_name in tiered:
        tiered[tier_name].sort(key=lambda u: (-u["count"], u["login"].lower()))

    return tiered


def generate_markdown(
    repo: str,
    total_merged: int,
    contributors: Dict[str, Dict[str, Any]],
    tiered: Dict[str, List[Dict[str, Any]]],
    timestamp: datetime.datetime,
) -> str:
    """Generate clean, GitHub-flavored Markdown for CONTRIBUTORS.md."""
    utc_str = timestamp.strftime("%Y-%m-%d %H:%M UTC")
    total_contributors = len(contributors)

    lines = [
        "<!-- AUTO-GENERATED — DO NOT EDIT DIRECTLY -->",
        "<!-- Generated by scripts/generate_contributors.py -->",
        "",
        "# PipeCD Contributors",
        "",
        "This contributor ladder recognizes and celebrates everyone who has contributed merged pull requests to [PipeCD](https://github.com/{repo}).".format(
            repo=repo
        ),
        "",
        f"> **Last updated:** {utc_str}  ",
        f"> **Total Merged PRs:** {total_merged:,} | **Total Contributors:** {total_contributors:,}  ",
        "> *Auto-refreshed daily via GitHub Actions.*",
        "",
        "---",
        "",
    ]

    for tier in TIERS:
        name = tier["name"]
        badge = tier["badge"]
        desc = tier["description"]
        users = tiered.get(name, [])

        lines.append(f"## {badge} {name}s ({desc})")
        lines.append("")

        if not users:
            lines.append("*No contributors in this tier yet.*")
            lines.append("")
            continue

        lines.append(f"*{len(users)} contributor{'s' if len(users) != 1 else ''} in this tier*")
        lines.append("")
        lines.append("| Contributor | Merged PRs |")
        lines.append("|:---|:---:|")

        for user in users:
            login = user["login"]
            html_url = user["html_url"]
            count = user["count"]
            lines.append(f"| [@{login}]({html_url}) | {count} |")

        lines.append("")

    lines.extend(
        [
            "---",
            "",
            "### How the Contributor Ladder Works",
            "",
            "- **Merged PRs Only:** Only merged pull requests are counted (opening a PR does not count until merged).",
            "- **Sorting:** Contributors within each tier are ordered by merged PR count (descending), with username alphabetically as a tiebreaker.",
            "- **Daily Updates:** This file is updated daily via automated workflow. If your pull request was recently merged, your profile will be added in the next scheduled run.",
            "- **Full Guidelines:** Learn more about contributor roles and progression in our [Contributor Ladder Guide](https://pipecd.dev/docs-dev/contribution-guidelines/contributor-ladder/).",
            "",
        ]
    )

    return "\n".join(lines)


def main():
    parser = argparse.ArgumentParser(description="Generate PipeCD CONTRIBUTORS.md")
    parser.add_argument(
        "--repo",
        default=DEFAULT_REPO,
        help=f"Target repository (owner/repo). Default: {DEFAULT_REPO}",
    )
    parser.add_argument(
        "--output",
        default="CONTRIBUTORS.md",
        help="Path to output markdown file. Default: CONTRIBUTORS.md",
    )
    parser.add_argument(
        "--token",
        default=None,
        help="GitHub token (defaults to GITHUB_TOKEN / GH_TOKEN env or gh auth)",
    )
    parser.add_argument(
        "--method",
        choices=["pulls", "search"],
        default="pulls",
        help="API method to query: 'pulls' (default, retrieves full history) or 'search' (/search/issues).",
    )

    args = parser.parse_args()

    token = get_auth_token(args.token)
    headers = {
        "Accept": "application/vnd.github.v3+json",
        "User-Agent": "pipecd-contributor-ladder-generator",
    }
    if token:
        headers["Authorization"] = f"Bearer {token}"
        print("Authenticated with GitHub token.")
    else:
        print("Warning: Running unauthenticated. GitHub API rate limits will be restricted.")

    session = requests.Session()

    try:
        if args.method == "search":
            contributors, total_merged = fetch_merged_prs_via_search(
                args.repo, headers, session
            )
        else:
            contributors, total_merged = fetch_merged_prs_via_pulls(
                args.repo, headers, session
            )
    except requests.RequestException as e:
        print(f"Error fetching data from GitHub API: {e}", file=sys.stderr)
        sys.exit(1)

    print(f"Total merged PRs counted: {total_merged}")
    print(f"Total contributors found: {len(contributors)}")

    tiered = assign_tiers(contributors)
    now = datetime.datetime.now(datetime.timezone.utc)
    markdown_content = generate_markdown(args.repo, total_merged, contributors, tiered, now)

    # Ensure target directory exists
    output_dir = os.path.dirname(args.output)
    if output_dir:
        os.makedirs(output_dir, exist_ok=True)

    with open(args.output, "w", encoding="utf-8") as f:
        f.write(markdown_content)

    print(f"Successfully generated {args.output}")


if __name__ == "__main__":
    main()
