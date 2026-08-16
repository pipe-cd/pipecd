---
date: 2026-08-16
title: "PipeCD's Newcomers Meetings: July 16 & August 6, 2026"
linkTitle: "Newcomers Meetings, Jul 16 & Aug 6, 2026"
weight: 968
description: "A recap of PipeCD's July 16 and August 6 2026 Newcomers meetings: the v0 to v1 evolution, how to start contributing, and where to find the community."
author: Harshit Ghagre ([@harshitghagre](https://github.com/harshitghagre))
categories: ["Announcement"]
tags: ["Newcomers", "Community", "Contributing", "Onboarding"]
---

Twice this summer, I sat down with a group of people who had never touched PipeCD before, some of them new to CNCF projects altogether, and walked them through the project from scratch. That's what the Newcomers meeting is for: no assumption that you've read the docs, no assumption you know what GitOps even means yet. Just an open call, a screen share, and an hour to answer whatever's on your mind.

We ran two of these sessions this year, on July 16 and August 6, 2026:

- [PipeCD's Newcomers meeting - July 16, 2026](https://youtu.be/hBXkaWD9190?si=bc3UQTsmZYVTvzzq)
- [PipeCD's Newcomers meeting - August 6, 2026](https://youtu.be/JhUGsf0bo4M?si=2YjKUCtcDwROankz)

If you were there, this is a recap. If you weren't, consider this the next best thing, and an open invitation to the next one.

## Why the project looks different than it did a year ago

The first thing I try to get out of the way early is the v0-to-v1 question, because it trips up almost everyone who opens the repo cold. If you clone PipeCD today, you'll find two versions of the piped agent living side by side, and that's confusing until you know why.

The original piped, what we now call v0, has every deployment stage built directly into the agent. Kubernetes sync, Terraform plan and apply, canary analysis, all of it compiled in. That was simple in the beginning, but it meant that adding support for a new platform meant changing piped itself, and that doesn't scale to a community project with contributors who each care about a different corner of the deployment world.

pipedv1 is the answer to that. Stage executors became plugins: separate binaries that piped loads and talks to over gRPC. Kubernetes, ECS, Lambda, Cloud Run, Terraform, wait, wait-approval, analysis, each one its own module now, each one able to move independently. I usually point newcomers to three things to read after the meeting if they want the full picture: [the original design post](https://pipecd.dev/blog/2024/11/28/overview-of-the-plan-for-pluginnable-pipecd/), [what actually changed in practice](https://pipecd.dev/blog/2025/09/02/what-is-new-in-pipedv1-plugin-arch-piped/), and [the migration guide](https://pipecd.dev/docs-dev/migrating-from-v0-to-v1/) if you're moving a real installation over.

What matters most if you're about to open a pull request:

- `platformProviders` and `cloudProviders` in piped config don't exist in v1. Plugins define `deployTargets` instead.
- Every application config uses `kind: Application` now, not the old per-platform kinds.
- Plugin code under `pkg/app/pipedv1/plugin/` is its own Go module and can only import the [piped-plugin-sdk-go](https://github.com/pipe-cd/piped-plugin-sdk-go) SDK, never the main `github.com/pipe-cd/pipecd` module.
- Legacy `piped` (`pkg/app/piped/`) and `pipedv1` (`pkg/app/pipedv1/`) both live in the repo right now, because the migration is gradual. Before you start on a `good first issue`, check which side it's on. If it's not obvious, ask in the issue, that's a completely normal question to ask.

## What actually happens when you send a pull request

The second half of both meetings is usually spent on the part people are most nervous about: what happens after they write code. Open source can feel like a black box from the outside even when the process behind it is ordinary, so we walk through it step by step.

It starts with an issue. [`good first issue`](https://github.com/pipe-cd/pipecd/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22) is the label to search for, and a comment saying you'd like to take it is enough to get assigned, though the expectation is a PR within about a week, so it's worth claiming something you can actually start on soon. For anything that isn't a one-line fix, I encourage people to describe their planned approach in the issue before writing code. It's a five-minute step that saves a lot of rework later.

From there, the mechanics: keep the diff small, somewhere around 300 lines is the rule of thumb, and split anything bigger. Sign off every commit with `git commit -s` for DCO. Write commit messages in the present tense with a capital first letter, the way `Add imports to Terraform plan result` reads. Target `master`. And before opening the PR, run `make check` locally, it's the same build, lint, and test steps CI will run anyway, so you find out about failures on your own machine first.

The [full contributor guide](https://pipecd.dev/docs-dev/contribution-guidelines/contributing/) covers the rest, the license header, the PR template for user-facing changes, how the bug and feature templates work. And code isn't the only door in: [documentation](https://pipecd.dev/docs-dev/contribution-guidelines/contributing-documentation/) and [blog posts](https://pipecd.dev/docs-dev/contribution-guidelines/contributing-blogs/) (this one included) need contributors too, and so does issue triage and conversation in [GitHub Discussions](https://github.com/pipe-cd/pipecd/discussions).

## Where the conversation continues after the meeting ends

Every Newcomers meeting ends the same way: with people asking where to go next. Here's what I tell them.

`#pipecd` on [CNCF Slack](https://cloud-native.slack.com/) is where most of the day-to-day questions get answered. The broader [PipeCD Development and Community Meeting](https://zoom-lfx.platform.linuxfoundation.org/meeting/96831504919?password=2f60b8ec-5896-40c8-aa1d-d551ab339d00) runs every two weeks with project news and issue triage, notes are at [bit.ly/pipecd-mtg-notes](https://bit.ly/pipecd-mtg-notes). Newcomers meetings themselves happen periodically, not on a fixed schedule, so Slack and the meeting notes are the place to watch for the next date. If you want something more structured, [LFX Mentorship](https://github.com/cncf/mentoring/blob/main/programs/lfx-mentorship/README.md#program-guidelines) is worth a look, and [one mentee wrote honestly about her first month](https://pipecd.dev/blog/2026/04/08/my-first-30-days-as-an-lfx-mentee-with-pipecd/) if you want to know what that actually feels like. And once you've merged 5 PRs across the pipe-cd organization and attended a public community meeting, you're eligible for membership in the GitHub org, though plenty of people contribute happily without ever going that route.

## See you at the next one

Both sessions, in the end, come down to the same thing: showing people that the project is smaller and friendlier up close than it looks from the outside. If you missed these two, the recordings are linked above. And if you're reading this because you're weighing your first contribution, the door's open, come find us in [#pipecd on Slack](https://cloud-native.slack.com/) and we'll figure out where to start.
