---
title: "Confidently review your changes with Plan Preview"
linkTitle: "Plan preview"
weight: 6
description: >
  Enables the ability to preview the deployment plan against a given commit before merging.
---

In order to help developers review the pull request with a better experience and more confidence to approve it to trigger the actual deployments,
PipeCD provides a way to preview the deployment plan of all updated applications by that pull request.

Here are what will be included currently in the result of plan-preview process:

- which application will be deployed once the pull request got merged
- which deployment strategy (QUICK_SYNC or PIPELINE_SYNC) will be used
- which resources will be added, deleted, or modified

This feature will available for all application kinds: KUBERNETES, TERRAFORM, CLOUD_RUN, LAMBDA and Amazon ECS.

![](/images/plan-preview-comment.png)
<p style="text-align: center;">
PlanPreview with GitHub actions <a href="https://github.com/pipe-cd/actions-plan-preview">pipe-cd/actions-plan-preview</a>
</p>

## Prerequisites

- Ensure the version of your Piped is at least `v0.11.0`.
- Having an API key that has `READ_WRITE` role to authenticate with PipeCD's Control Plane. A new key can be generated from `settings/api-key` page of your PipeCD web.

## Usage

Plan-preview result can be requested by using `pipectl` command-line tool as below:

``` console
pipectl plan-preview \
  --address={ PIPECD_CONTROL_PLANE_ADDRESS } \
  --api-key={ PIPECD_API_KEY } \
  --repo-remote-url={ REPO_REMOTE_GIT_SSH_URL } \
  --head-branch={ HEAD_BRANCH } \
  --head-commit={ HEAD_COMMIT } \
  --base-branch={ BASE_BRANCH } \
  --sort-label-keys={ SORT_LABEL_KEYS }
```

You can run it locally or integrate it to your CI system to run automatically when a new pull request is opened/updated. Use `--help` to see more options.

``` console
pipectl plan-preview --help
```

## GitHub Actions

If you are using GitHub Actions, you can seamlessly integrate our prepared [actions-plan-preview](https://github.com/pipe-cd/actions-plan-preview) to your workflows. This automatically comments the plan-preview result on the pull request when it is opened or updated. You can also trigger to run plan-preview manually by leave a comment `/pipecd plan-preview` on the pull request.
