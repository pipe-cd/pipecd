---
title: "Concepts"
linkTitle: "Concepts"
weight: 2
description: >
  This page describes several core concepts in PipeCD.
---

![Architecture Overview](/images/architecture-overview.png)
<p style="text-align: center;"
>
Component Architecture
</p>

### Control Plane

The Control Plane is the centralized management service of PipeCD. It coordinates all activities between users, projects, and `piped` instances.

The Control Plane remains the backbone of the system but is now fully plugin-aware. Instead of directly handling deployment logic for specific platforms, it interacts with `piped` agents that run plugin binaries, allowing the Control Plane to manage deployments across any platform supported by plugins.

For more detailed information about Control Plane architecture and components, see [Architecture overview of Control Plane](../user-guide/managing-controlplane/architecture-overview/).

### piped

`piped` is a binary, agent component responsible for executing deployments in PipeCD. `piped` now adopts **plugin-based** **architecture**, transforming from a single-purpose executor into a lightweight runtime capable of running any deployment logic defined by plugins. The `piped` component is designed to be stateless.

### Plugins

![PipeCD v1 Plugin Architecture](/images/pipecdv1-architecture.png)
<p style="text-align: center;"
>
PipeCD v1 Plugin Architecture
</p>

Plugins replace the concept of Platform Providers (also called Cloud Providers) from earlier versions of PipeCD.

A Plugin defines how and where deployments are executed. Each plugin encapsulates the logic required to interact with a specific platform or tool - for example, Kubernetes, Terraform, or ECS.

In this architecture, plugins are the actors who execute deployments on behalf of `piped`. Rather than containing built-in logic for every platform, `piped` now loads plugin binaries at startup and communicates with them over gRPC. This design allows support for new or custom platforms.

>**Note:**  
>Check out the [PipeCD Community Plugins repository](https://github.com/pipe-cd/community-plugins) to browse available plugins and learn how to create your own.

### Project

A Project is a logical group of applications managed together by a team of users.
Each project can connect to multiple `piped` instances running across different environments or clouds.

Projects use role-based access control (RBAC) to manage permissions:

- Viewer – can view applications and deployments within the project.
- Editor – includes Viewer permissions and can perform actions that modify state, such as triggering or cancelling deployments.
- Admin – includes Editor permissions and can manage project settings, members, and associated `piped` instances.

### Application

An Application represents a collection of declarative resources and configurations that represent a deployable unit managed by PipeCD.

Applications are no longer tied to a specific platform or technology such as Kubernetes, Terraform, or ECS.
Instead, they are platform-agnostic objects whose deployment behavior is determined dynamically through **[Plugins](#plugins).**

### Application Configuration

A declarative YAML file that defines how an application is managed by PipeCD.
It specifies metadata, labels, and the deployment logic that `piped` executes through plugins.
Each application has one configuration file stored in its Git directory, typically named app.pipecd.yaml.

### Application Directory

A directory in the Git repository that contains the application’s configuration file and related manifests.
This directory represents the source of truth for the application’s desired state.

### Deployment

A Deployment is the process of bringing an application’s live state in line with its desired state defined in Git.
When a deployment completes successfully, the running environment matches the configuration from the target Git commit.

### Sync Strategy

PipeCD provides 3 different ways to keep your application’s live state consistent with its desired state stored in Git.
Depending on your deployment workflow, you can choose from one of the following sync strategies:

- Quick Sync: A fast, single-step method to sync your deployment with the desired state. PipeCD automatically generates a pipeline composed of predefined Quick Sync stages provided by plugins linked to the application. Quick Sync is generally used when you need to rapidly apply configuration changes without a gradual rollout.

- Pipeline Sync: A customizable, step-by-step sync process that follows the pipeline you define in your application configuration file. Use Pipeline Sync when you need more control over how updates are rolled out.

- Auto Sync: When you trigger a sync without specifying a strategy, `piped` automatically selects the most appropriate method based on your application configuration.

### Analysis Provider
An external product that provides metrics/logs to evaluate deployments, such as `Prometheus`, `Datadog`, `Stackdriver`, `CloudWatch` and so on.

### Stage
A single step within a deployment pipeline. Stages are defined in the application configuration file (`app.pipecd.yaml`) and are executed sequentially by `piped`.
Common built-in stages include `K8S_CANARY_ROLLOUT`, `K8S_PRIMARY_ROLLOUT`, `K8S_CANARY_CLEAN`, `TERRAFORM_PLAN`, `TERRAFORM_APPLY`, `WAIT`, `WAIT_APPROVAL`, `ANALYSIS`, and `SCRIPT_RUN`.

### pipectl
The official command-line tool for interacting with the PipeCD Control Plane API.
`pipectl` allows users to add and sync applications, get deployment status, encrypt secrets, and more from the terminal.
See [Command-line tool: pipectl](../user-guide/command-line-tool/) for usage information.

### Plan Preview
A feature that previews the expected changes of a deployment before it is actually applied.
When integrated with a CI system such as GitHub Actions, Plan Preview posts a comment on pull requests showing which resources would be added, modified, or deleted.
See [Plan preview](../user-guide/managing-application/plan-preview/) for details.

### Event Watcher
A `piped` feature that watches for external events, such as a new container image pushed to a registry or a new Helm chart version published, and automatically updates files in the Git repository to trigger a new deployment.
See [Event watcher](../user-guide/managing-piped/configuring-event-watcher/) for configuration details.

### Deployment Trace
A feature that links a deployment in PipeCD to the CI build that produced its artifacts (e.g., a GitHub Actions run or a Jenkins build).
Deployment Trace provides end-to-end traceability from a commit to the resulting deployment.
See [Deployment Trace](../user-guide/managing-application/deployment-trace/) for details.

### Insights
The analytics dashboard of PipeCD that displays delivery performance metrics such as deployment frequency, lead time, mean time to restore (MTTR), and change failure rate.
See [Insights](../user-guide/observability-and-metrics/insights/) for details.

### Configuration Drift Detection
A `piped` feature that periodically compares the live state of an application's resources with the desired state declared in the Git repository.
When a drift is detected, for example when someone manually edits a resource in a cluster, PipeCD visualizes the difference and can optionally trigger a sync to reconcile the drift.
See [Configuration drift detection](../user-guide/managing-application/configuration-drift-detection/) for details.

### Application Live State
The real-time state of an application's resources as they exist on the target platform.
PipeCD continuously monitors live state and displays it on the web console, including resource status, health, and the relationship between resources (e.g., Deployment → ReplicaSet → Pod).
See [Application live state](../user-guide/managing-application/app-live-state/) for details.

### Secret Management
PipeCD provides a mechanism for safely managing secrets used in application manifests.
Secrets are encrypted using `pipectl` and stored in the Git repository in encrypted form. `piped` decrypts them at deployment time, so plain-text secrets never leave the deployment environment.
See [Secret management](../user-guide/managing-application/secret-management/) for details.

### Launcher
A component (`cmd/launcher`) that manages the lifecycle of the `piped` agent process.
See [Remote upgrade and remote config](../user-guide/managing-piped/remote-upgrade-remote-config/) for more information.