---
title: "CI Overview"
linkTitle: "CI Overview"
weight: 30
description: >
  An overview of PipeCD’s GitHub Actions CI workflows, what each workflow does,
  and when it runs.
---

## Who Is This Document For?

This document is intended for contributors who:
- submit pull requests to PipeCD
- want to understand which CI workflows run and why
- need basic context to debug CI failures
---

## CI Philosophy

PipeCD’s CI is designed to:
- provide fast feedback on pull requests
- prevent regressions on the default and release branches
- keep workflows explicit and deterministic
- automate release and publishing tasks safely

---

## Workflow Overview

Below is a brief description of each GitHub Actions workflow under
`.github/workflows/`.

### `test.yaml`
**What it does:** Runs Go unit tests to validate core functionality.  
**When it runs:** On pull requests and pushes to the default branch.

---

### `lint.yaml`
**What it does:** Runs linters and static analysis to enforce code quality.  
**When it runs:** On pull requests and pushes to the default branch.

---

### `build.yaml`
**What it does:** Builds PipeCD binaries to ensure the project compiles correctly.  
**When it runs:** On pull requests and pushes to the default branch.

---

### `gen.yaml`
**What it does:** Verifies that generated files are up to date.  
**When it runs:** On pull requests to prevent uncommitted generated changes.

---

### `web.yaml`
**What it does:** Builds and validates the PipeCD Web UI.  
**When it runs:** On pull requests that include web-related changes.

---

### `chart.yaml`
**What it does:** Validates Helm charts used to deploy PipeCD.  
**When it runs:** On pull requests affecting Helm chart files.

---

### `codeql-analysis.yaml`
**What it does:** Runs CodeQL security analysis to detect potential vulnerabilities.  
**When it runs:** On a scheduled basis and on selected pushes.

---

### `labeler.yaml`
**What it does:** Automatically applies labels to pull requests based on changed files.  
**When it runs:** When a pull request is opened or updated.

---

### `stale.yaml`
**What it does:** Marks and closes inactive issues and pull requests.  
**When it runs:** On a scheduled basis.

---

### `cherry_pick.yaml`
**What it does:** Automates cherry-picking changes into release branches.  
**When it runs:** When manually triggered by maintainers.

---

### `plugin_release.yaml`
**What it does:** Builds and publishes PipeCD plugins.  
**When it runs:** During plugin release workflows.

---

### `publish_binary.yaml`
**What it does:** Builds and publishes PipeCD binaries.  
**When it runs:** During release workflows.

---

### `publish_pipedv1_exp.yaml`
**What it does:** Builds and publishes experimental `pipedv1` container images.  
**When it runs:** On prerelease or experimental release triggers.

---

### `publish_image_chart.yaml`
**What it does:** Publishes Helm charts and related container images.  
**When it runs:** During release workflows.

---

### `publish_tool.yaml`
**What it does:** Builds and publishes PipeCD CLI and tooling images.  
**When it runs:** On tool release events.

---

### `publish_site.yaml`
**What it does:** Builds and publishes the PipeCD documentation website.  
**When it runs:** On documentation updates and release events.

---

### `prerelease.yaml`
**What it does:** Prepares prerelease artifacts for testing and validation.  
**When it runs:** When a prerelease is triggered.

---

### `release.yaml`
**What it does:** Orchestrates the full PipeCD release process.  
**When it runs:** When a release is created.

---

### `build_tool.yaml`
**What it does:** Builds internal tooling images used by CI and release workflows.  
**When it runs:** On pushes affecting tool-related code.

---

### `code-butler.yaml`
**What it does:** Runs automated code maintenance tasks.  
**When it runs:** On a scheduled basis or when triggered by maintainers.

---
## Debugging CI Failures

When a CI check fails:
1. Identify which workflow failed.
2. Read the workflow name to understand its responsibility.
3. Check the job logs for the specific failure.

Most CI failures are scoped to a single concern (tests, lint, build, or web),
which usually makes the root cause clear.




## References

- `.github/workflows/`
- Contributing Guide
- GitHub Actions documentation
