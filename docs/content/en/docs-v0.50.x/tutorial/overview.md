---
title: "Mechanism Overview"
linkTitle: "Mechanism Overview"
weight: 1
description: >
  Quickly learn about the basic mechanism of PipeCD.
---

## Architecture

PipeCD follows a client-server architecture with two main components:

- **Control Plane**: The centralized component that manages the web console, API, and data storage. It is responsible for managing projects, applications, pipeds, and deployments.
- **Piped**: The agent that runs in your cluster or environment. It watches for changes in your Git repository and executes deployments according to the pipeline definition.

![Architecture overview](/images/architecture-overview.png)

### How It Works

1. **Developer** pushes a change to the Git repository (e.g., updating a manifest or configuration).
2. **Piped** detects the change by periodically polling the Git repository.
3. **Piped** compares the desired state (from Git) with the current state (in the target platform).
4. If there is a difference, **Piped** automatically triggers a deployment.
5. **Piped** executes the deployment pipeline stages defined in `app.pipecd.yaml`.
6. The deployment progress and results are reported to the **Control Plane** and displayed on the web console.

### Scalability in an Organization

PipeCD is designed to scale across large organizations:

- **Multiple Pipeds**: You can run multiple Piped agents, each managing different environments or teams.
- **Multi-cluster / Multi-cloud**: Each Piped can connect to different platforms (Kubernetes, Cloud Run, ECS, Lambda, Terraform, etc.).
- **Single Control Plane**: All Pipeds connect to one Control Plane, providing a unified view of all deployments across the organization.
- **Project-based isolation**: Different teams can have their own projects within the same Control Plane.

## See Also

- [Concepts](../../docs-v0.50.x/concepts/)
- [Architecture Overview](../../docs-v0.50.x/user-guide/managing-controlplane/architecture-overview/)

---

[Next: Prerequisites >](../prerequisites/)
