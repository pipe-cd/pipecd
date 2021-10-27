---
title: "Feature Status"
linkTitle: "Feature Status"
weight: 8
description: >
  This page lists the relative maturity of every PipeCD features.
---

Please note that the phases (Incubating, Alpha, Beta, and Stable) are applied to individual features within the project, not to the project as a whole.

## Feature Phase Definitions

| Phase | Definition |
|-|-|
| Incubating | Under planning/developing the prototype and still not ready to be used. |
| Alpha | Demo-able, works end-to-end but has limitations. No guarantees on backward compatibility. |
| Beta | **Usable in production**. Documented. |
| Stable | Production hardened. Backward compatibility. Documented. |

## PipeCD Features

### Kubernetes Deployment

| Feature | Phase |
|-|-|
| Quick sync deployment | Beta |
| Deployment with a defined pipeline (e.g. canary, analysis) | Beta |
| [Automated rollback](/docs/user-guide/rolling-back-a-deployment/) | Beta |
| [Automated configuration drift detection](/docs/user-guide/configuration-drift-detection/) | Beta |
| [Application live state](/docs/user-guide/application-live-state/) | Beta |
| Support Helm | Beta |
| Support Kustomize | Beta |
| Support Istio service mesh | Beta |
| Support SMI service mesh | Incubating |
| Support [AWS App Mesh](https://aws.amazon.com/app-mesh/) | Incubating |
| [Plan preview](/docs/user-guide/plan-preview) | Beta |

### Terraform Deployment

| Feature | Phase |
|-|-|
| Quick sync deployment | Beta |
| Deployment with a defined pipeline (e.g. manual-approval) | Beta |
| [Automated rollback](/docs/user-guide/rolling-back-a-deployment/) | Beta |
| [Automated configuration drift detection](/docs/user-guide/configuration-drift-detection/) | Incubating |
| [Application live state](/docs/user-guide/application-live-state/) | Incubating |
| [Plan preview](/docs/user-guide/plan-preview) | Beta |

### CloudRun Deployment

| Feature | Phase |
|-|-|
| Quick sync deployment | Beta |
| Deployment with a defined pipeline (e.g. canary, analysis) | Beta |
| [Automated rollback](/docs/user-guide/rolling-back-a-deployment/) | Beta |
| [Automated configuration drift detection](/docs/user-guide/configuration-drift-detection/) | Incubating |
| [Application live state](/docs/user-guide/application-live-state/) | Incubating |
| [Plan preview](/docs/user-guide/plan-preview) | Alpha |

### Lambda Deployment

| Feature | Phase |
|-|-|
| Quick sync deployment | Beta |
| Deployment with a defined pipeline (e.g. canary, analysis) | Beta |
| [Automated rollback](/docs/user-guide/rolling-back-a-deployment/) | Beta |
| [Automated configuration drift detection](/docs/user-guide/configuration-drift-detection/) | Incubating |
| [Application live state](/docs/user-guide/application-live-state/) | Incubating |
| [Plan preview](/docs/user-guide/plan-preview) | Alpha |

### Amazon ECS Deployment

| Feature | Phase |
|-|-|
| Quick sync deployment | Alpha |
| Deployment with a defined pipeline (e.g. canary, analysis) | Alpha |
| [Automated rollback](/docs/user-guide/rolling-back-a-deployment/) | Beta |
| [Automated configuration drift detection](/docs/user-guide/configuration-drift-detection/) | Incubating |
| [Application live state](/docs/user-guide/application-live-state/) | Incubating |
| Support [AWS App Mesh](https://aws.amazon.com/app-mesh/) | Incubating |
| [Plan preview](/docs/user-guide/plan-preview) | Alpha |

### Piped's Core

| Feature | Phase |
|-|-|
| [Deployment wait stage](/docs/user-guide/adding-a-wait-stage/) | Beta |
| [Deployment manual approval stage](/docs/user-guide/adding-a-manual-approval/) | Beta |
| [Notification](/docs/operator-manual/piped/configuring-notifications/) to Slack | Beta |
| [Notification](/docs/operator-manual/piped/configuring-notifications/) to external service via webhook | Alpha |
| [Secrets management](/docs/user-guide/secret-management/) - Storing secrets safely in the Git repository | Beta |
| [Event watcher](/docs/user-guide/event-watcher/) - Updating files in Git automatically for given events | Alpha |
| [Pipectl](/docs/user-guide/command-line-tool/) - Command-line tool for interacting with control-plane | Beta |
| Deployment plugin - Allow executing user-created deployment plugin | Incubating |
| [ADA](/docs/user-guide/automated-deployment-analysis/) (Automated Deployment Analysis) by Prometheus metrics | Alpha |
| [ADA](/docs/user-guide/automated-deployment-analysis/) by Datadog metrics | Alpha |
| [ADA](/docs/user-guide/automated-deployment-analysis/) by Stackdriver metrics | Incubating |
| [ADA](/docs/user-guide/automated-deployment-analysis/) by Stackdriver log | Incubating |
| [ADA](/docs/user-guide/automated-deployment-analysis/) by CloudWatch metrics | Incubating |
| [ADA](/docs/user-guide/automated-deployment-analysis/) by CloudWatch log | Incubating |
| [ADA](/docs/user-guide/automated-deployment-analysis/) by HTTP request (smoke test...) | Incubating |
| [Remote upgrade](/docs/operator-manual/piped/remote-upgrade-remote-config/#remote-upgrade) - Ability to upgrade Piped from the web console | Alpha |
| [Remote config](/docs/operator-manual/piped/remote-upgrade-remote-config/#remote-config) - Watch and reload configuration from a remote location such as Git | Alpha |

### ControlPlane's Core

| Feature | Phase |
|-|-|
| Project/Environment/Piped/Application/Deployment management | Beta |
| Rendering deployment pipeline in realtime | Beta |
| Canceling a deployment from console | Beta |
| Triggering a deployment manually from console | Beta |
| Authentication by username/password for static admin | Beta |
| GitHub & GitHub Enterprise SSO | Beta |
| Google SSO | Incubating |
| Support GCP [Firestore](https://cloud.google.com/firestore) as data store | Beta |
| Support [MySQL v8.0](https://www.mysql.com/) as data store | Beta |
| Support GCP [GCS](https://cloud.google.com/storage) as file store | Beta |
| Support AWS [S3](https://aws.amazon.com/s3/) as file store | Beta |
| Support [Minio](https://github.com/minio/minio) as file store | Beta |
| [Insights](/docs/user-guide/insights/) - Show the delivery performance of a team or an application | Incubating |
| Deployment Chain - Allow rolling out to multiple clusters gradually or promoting across environments | Incubating |
| [Metrics](/docs/operator-manual/control-plane/metrics/) - Dashboards for PipeCD and Piped metrics | Beta |
