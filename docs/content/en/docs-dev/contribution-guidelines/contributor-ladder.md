---
title: "Contributor Ladder"
linkTitle: "Contributor Ladder"
weight: 5
description: >
  How PipeCD recognizes and celebrates community contributors through our contributor ladder.
---

The **PipeCD Contributor Ladder** is designed to welcome new contributors, recognize ongoing contributions, and provide clear paths for community progression.

Every merged contribution makes a difference to PipeCD — whether fixing a bug, improving documentation, designing features, or enhancing test coverage.

---

## Ladder Tiers

The contributor ladder currently tracks contributions based on **merged pull requests** to the [`pipe-cd/pipecd`](https://github.com/pipe-cd/pipecd) repository:

| Tier | Merged PRs | Description & Recognition |
| :--- | :---: | :--- |
| 🌱 **Newcomer** | **1** | Welcome to the community! Listed in the Newcomer section of [`CONTRIBUTORS.md`](https://github.com/pipe-cd/pipecd/blob/master/CONTRIBUTORS.md). |
| 🛠️ **Contributor** | **2–4** | Continued active involvement and consistent contributions across any part of the project. |
| 🚀 **Core Contributor** | **5+** | Established and trusted contributors with a strong track record. Eligible to request membership in the `pipe-cd` GitHub organization. |

---

## How Counting Works

- **Merged PRs Only:** Only pull requests that have been reviewed, approved, and merged are counted toward ladder tiers. Merely opening a PR does not count until it is merged.
- **Automated Daily Refresh:** The contributor ladder at [`CONTRIBUTORS.md`](https://github.com/pipe-cd/pipecd/blob/master/CONTRIBUTORS.md) is regenerated automatically every day via a GitHub Actions workflow.
- **No Manual Edits:** The list is maintained by automation; please do not submit manual pull requests to edit `CONTRIBUTORS.md`.
- **Bot Exclusions:** Automated bots and service accounts (such as `dependabot[bot]` and `github-actions[bot]`) are excluded from the ladder.

---

## Becoming a Member of the PipeCD GitHub Organization

Once you reach the **Core Contributor** tier (5+ merged PRs), you are invited to apply for membership in the `pipe-cd` GitHub organization!

### Requirements:
1. Have at least **5 merged PRs** in repositories under the `pipe-cd` organization.
2. Have attended a [PipeCD Community Meeting](https://zoom-lfx.platform.linuxfoundation.org/meeting/96831504919?password=2f60b8ec-5896-40c8-aa1d-d551ab339d00).
3. Reach out to the maintainers in the `#pipecd` channel on [CNCF Slack](https://cloud-native.slack.com/) or during a community meeting.

---

## Getting Started

Looking for ways to begin or climb the ladder?
- Check out issues labeled [**good first issue**](https://github.com/pipe-cd/pipecd/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22) for beginner-friendly tasks.
- Help improve our docs by reviewing the [Contribute to PipeCD Documentation](../contributing-documentation/) guide.
- Read our [General Contribution Guide](../contributing/) for details on local environment setup, testing, and DCO sign-off.
- Join the discussion in our `#pipecd` channel on [CNCF Slack](https://cloud-native.slack.com/).

---

## Future Roadmap

This contributor ladder is being rolled out in iterative phases as discussed in [GitHub Issue #6548](https://github.com/pipe-cd/pipecd/issues/6548):

- **Phase 1 (Current):** Automated tier ladder based on merged PR count in `CONTRIBUTORS.md`.
- **Phase 2 (Planned):** A comprehensive contribution scoring model taking into account issue triage, code reviews, blog posts, and architectural proposals.
- **Phase 3 (Planned):** Community leaderboard announcements and recognition in the PipeCD Slack channel.
