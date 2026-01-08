---
title: "Contribute to PipeCD"
linkTitle: "Contribute to PipeCD"
weight: 1
description: >
  This page describes how to contribute to PipeCD.
---

PipeCD is an open source project that anyone in the community can use, improve, and enjoy. We'd love you to join us!

## Code of Conduct

PipeCD follows the [CNCF Code of Conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md). Please read it to understand which actions are acceptable and which are not.

## Ways to Contribute

There are many ways to contribute, and many don't involve writing code:

- **Use PipeCD** — Follow the [Quickstart](/docs/quickstart/) and report issues you encounter
- **Triage issues** — Browse [open issues](https://github.com/pipe-cd/pipecd/issues), provide workarounds, ask for clarification, or suggest labels
- **Fix bugs** — Issues labeled [good first issue](https://github.com/pipe-cd/pipecd/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22) are great starting points
- **Improve docs** — See [Contribute to PipeCD Documentation](../contributing-documentation/)
- **Participate in discussions** — Share ideas in [GitHub Discussions](https://github.com/pipe-cd/pipecd/discussions)

## Join the Community

### Slack

We have a `#pipecd` channel on [CNCF Slack](https://cloud-native.slack.com/) for discussions and help. You can also assist other users in the channel.

### Community Meetings

We host [PipeCD Development and Community Meetings](https://zoom-lfx.platform.linuxfoundation.org/meeting/96831504919?password=2f60b8ec-5896-40c8-aa1d-d551ab339d00) every 2 weeks where we share project news, demos, answer questions, and triage issues. View the [Meeting Notes and Agenda](https://bit.ly/pipecd-mtg-notes).

### Become a Member

If you'd like to join the pipe-cd GitHub organization:

- Have at least 5 PRs merged to repositories in the pipe-cd organization
- Attend a PipeCD public community meeting

Membership isn't required to contribute—it's for those who want to contribute long-term or take ownership of features.

## Development Setup

### Project Structure

PipeCD consists of several components:

| Component | Description |
|-----------|-------------|
| [cmd/pipecd](https://github.com/pipe-cd/pipecd/tree/master/cmd/pipecd) | Control Plane — manages deployment data and provides gRPC API |
| [cmd/piped](https://github.com/pipe-cd/pipecd/tree/master/cmd/piped) | Piped agent — runs in your cluster |
| [cmd/pipectl](https://github.com/pipe-cd/pipecd/tree/master/cmd/pipectl) | Command-line tool |
| [cmd/launcher](https://github.com/pipe-cd/pipecd/tree/master/cmd/launcher) | Command executor for remote upgrade |
| [web](https://github.com/pipe-cd/pipecd/tree/master/web) | Web UI |
| [docs](https://github.com/pipe-cd/pipecd/tree/master/docs) | Documentation |

### Prerequisites

- Go (check `go.mod` for version)
- Node.js and Yarn (for web development)
- Docker
- kubectl
- kind (Kubernetes in Docker)

### Starting a Local Environment

#### 1. Update Dependencies

```bash
make update/go-deps
make update/web-deps
```

> Starting a local environment may fail if dependencies are not up to date.

#### 2. Start Local Cluster and Registry

```bash
make up/local-cluster
```

This creates the `pipecd` Kubernetes namespace and starts a local registry.

To clean up later:

```bash
make down/local-cluster
```

#### 3. Run Control Plane

```bash
make run/pipecd
```

To stop:

```bash
make stop/pipecd
```

#### 4. Port Forward

In a separate terminal:

```bash
kubectl port-forward -n pipecd svc/pipecd 8080
```

#### 5. Access the UI

Open [http://localhost:8080?project=quickstart](http://localhost:8080?project=quickstart)

Login credentials:
- **Username:** `hello-pipecd`
- **Password:** `hello-pipecd`

### Running Piped Agent

1. Ensure Control Plane is running and accessible.

2. Register a new Piped:
   - Go to Settings → Piped (or [http://localhost:8080/settings/piped?project=quickstart](http://localhost:8080/settings/piped?project=quickstart))
   - Add a new piped and copy the generated Piped ID and base64 key

3. Create `piped-config.yaml` in the repository root:

   > **Tip:** Create a `.dev` folder (gitignored) for multiple config files.

   ```yaml
   apiVersion: pipecd.dev/v1beta1
   kind: Piped
   spec:
     projectID: quickstart
     pipedID: <YOUR_PIPED_ID>
     pipedKeyData: <YOUR_PIPED_BASE64_KEY>
     apiAddress: localhost:8080
     repositories:
       - repoId: example
         remote: git@github.com:pipe-cd/examples.git
         branch: master
     syncInterval: 1m
     plugins:
       - name: kubernetes
         port: 7001
         url: <PLUGIN_DOWNLOAD_URL>  # Get from https://github.com/pipe-cd/pipecd/releases
         deployTargets:
           - name: local
             config:
               kubeConfigPath: /path/to/.kube/config
   ```

   > **Note:** Plugins are versioned independently from PipeCD. Download URLs for official plugins can be found on the [PipeCD releases page](https://github.com/pipe-cd/pipecd/releases). Look for releases tagged with `pkg/app/pipedv1/plugin/kubernetes/`.

4. Start the Piped agent:

   ```bash
   make run/piped CONFIG_FILE=path/to/piped-config.yaml INSECURE=true
   ```

### Useful Commands

Check the [Makefile](https://github.com/pipe-cd/pipecd/blob/master/Makefile) for available commands. Common ones:

```bash
make check          # Run all checks (lint, test, etc.)
make test/go        # Run Go tests
make test/web       # Run web tests
make build/go       # Build Go binaries
```

## Pull Request Process

### Before You Start

1. **Find or create an issue** — Check if an issue exists, or open one for new features/bugs
2. **Claim the issue** — Comment "I'd like to work on this" to get assigned (expected to submit PR within 7 days)
3. **Investigate first** — Discuss your approach in the issue before coding to reduce back-and-forth on the PR
4. **Focus on one issue** — Especially for newcomers, work on one issue at a time

### Submitting a Pull Request

1. **Keep PRs small** — Aim for ~300 lines of diff. Split larger changes into multiple PRs.

2. **Use descriptive titles** — Follow commit message style: present tense, capitalize first letter

   ```
   Add imports to Terraform plan result
   ```

3. **Sign off commits** — Required for DCO compliance:

   ```bash
   git commit -s -m "Your commit message"
   ```

4. **Run checks before submitting:**

   ```bash
   make check
   ```

5. **Target the `master` branch** — All PRs should be opened against `master`

### User-Facing Changes

If your change affects users, update the PR description:

```md
**Does this PR introduce a user-facing change?**:
- **How are users affected by this change**:
- **Is this breaking change**:
- **How to migrate (if breaking change)**:
```

### Licensing

By contributing, you agree to license your contributions under the Apache License Version 2.0. Add this header to new files:

```go
// Copyright 2025 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
```

## Reporting Issues

### Bugs

File bugs at [GitHub Issues](https://github.com/pipe-cd/pipecd/issues/new?assignees=&labels=kind%2Fbug&projects=&template=bug-report.md):

- **One issue, one bug**
- **Provide reproduction steps** — List all steps to reproduce the issue

### Security Issues

Report security vulnerabilities privately via Slack or Twitter DM to maintainers listed in [MAINTAINERS.md](https://github.com/pipe-cd/pipecd/blob/master/MAINTAINERS.md).

### Feature Requests

- [Enhancement request](https://github.com/pipe-cd/pipecd/issues/new?assignees=&labels=kind%2Fenhancement&projects=&template=enhancement.md) — Improvements to existing features
- [Feature request](https://github.com/pipe-cd/pipecd/issues/new?assignees=&labels=kind%2Ffeature&projects=&template=new-feature.md) — Entirely new features

## What Happens Next?

The maintainers will review your PR. We'll help with obvious issues and work with you to get it merged. Thank you for contributing!
