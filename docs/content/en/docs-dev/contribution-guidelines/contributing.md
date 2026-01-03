---
title: "Contribute to PipeCD"
linkTitle: "Contribute to PipeCD"
weight: 1
description: >
  This page describes how to contribute to PipeCD.
---

PipeCD is an open source project that anyone in the community can use, improve, and enjoy. We'd love you to join us!

## Ways to contribute

There are many ways to contribute, and many don't involve writing code:

- **Use PipeCD** — Follow the [Quickstart](/docs/quickstart/) and report issues you encounter
- **Triage issues** — Browse [open issues](https://github.com/pipe-cd/pipecd/issues), provide workarounds, ask for clarification, or suggest labels
- **Fix bugs** — Issues labeled [good first issue](https://github.com/pipe-cd/pipecd/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22) are great starting points
- **Improve docs** — See [Contribute to PipeCD Documentation](../contributing-documentation/)
- **Write blog posts** — See [Contribute to PipeCD Blogs](../contributing-blogs/)
- **Build plugins** — See [Contribute to piped Plugins](../contributing-plugins/)

## Join the community

### Slack

We have a `#pipecd` channel on [CNCF Slack](https://cloud-native.slack.com/) for discussions and help.

### Community meetings

We host [PipeCD Development and Community Meetings](https://zoom-lfx.platform.linuxfoundation.org/meeting/96831504919?password=2f60b8ec-5896-40c8-aa1d-d551ab339d00) every 2 weeks. View the [Meeting Notes and Agenda](https://bit.ly/pipecd-mtg-notes).

## Development setup

PipeCD consists of several components:

| Component | Description |
|-----------|-------------|
| [cmd/pipecd](https://github.com/pipe-cd/pipecd/tree/master/cmd/pipecd) | Control Plane — manages deployment data and provides gRPC API |
| [cmd/piped](https://github.com/pipe-cd/pipecd/tree/master/cmd/piped) | Piped agent — runs in your cluster |
| [cmd/pipectl](https://github.com/pipe-cd/pipecd/tree/master/cmd/pipectl) | Command-line tool |
| [web](https://github.com/pipe-cd/pipecd/tree/master/web) | Web UI |
| [docs](https://github.com/pipe-cd/pipecd/tree/master/docs) | Documentation |

### Quick start

```bash
# Update dependencies
make update/go-deps
make update/web-deps

# Start local cluster and registry
make up/local-cluster

# Run Control Plane
make run/pipecd

# Port forward (in another terminal)
kubectl port-forward -n pipecd svc/pipecd 8080

# Access UI at http://localhost:8080?project=quickstart
# Login: hello-pipecd / hello-pipecd
```

For full setup instructions including Piped agent, see the [CONTRIBUTING.md](https://github.com/pipe-cd/pipecd/blob/master/CONTRIBUTING.md#development).

## Pull request process

1. **Claim the issue** — Comment "I'd like to work on this" to get assigned
2. **Investigate first** — Discuss your approach in the issue before coding
3. **Keep PRs small** — Aim for ~300 lines of diff
4. **Sign off commits** — Use `git commit -s` for DCO compliance
5. **Run checks** — Execute `make check` before submitting

All PRs should target the `master` branch.

## Code of Conduct

PipeCD follows the [CNCF Code of Conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md).

## Full contributing guide

For complete details on development, commit messages, licensing, and more, see [CONTRIBUTING.md](https://github.com/pipe-cd/pipecd/blob/master/CONTRIBUTING.md).
