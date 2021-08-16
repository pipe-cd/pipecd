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
| Beta | Usable in production. Documented. |
| Stable | Production hardened. Backward compatibility. Documented. |

## PipeCD Features

### Kubernetes Deployment

| Feature | Phase |
|-|-|
| Quick Sync Deployment | Beta |
| Deployment with the Specified Pipeline (canary, bluegreen...) | Beta |
| Automated Rollback | Beta |
| [Automated Configuration Drift Detection](/docs/user-guide/configuration-drift-detection/) | Beta |
| [Application Live State](/docs/user-guide/application-live-state/) | Beta |
| Support Helm | Beta |
| Support Kustomize | Beta |
| Support Istio Mesh | Beta |
| Support SMI Mesh | Incubating |
| Support [AWS App Mesh](https://aws.amazon.com/app-mesh/) | Incubating |
| [Plan Preview](/docs/user-guide/plan-preview) | Alpha |

### Terraform Deployment

| Feature | Phase |
|-|-|
| Quick Sync Deployment | Beta |
| Deployment with the Specified Pipeline | Beta |
| Automated Rollback | Beta |
| [Automated Configuration Drift Detection](/docs/user-guide/configuration-drift-detection/) | Incubating |
| [Application Live State](/docs/user-guide/application-live-state/) | Incubating |
| [Plan Preview](/docs/user-guide/plan-preview) | Alpha |

### CloudRun Deployment

| Feature | Phase |
|-|-|
| Quick Sync Deployment | Beta |
| Deployment with the Specified Pipeline | Beta |
| Automated Rollback | Beta |
| [Automated Configuration Drift Detection](/docs/user-guide/configuration-drift-detection/) | Incubating |
| [Application Live State](/docs/user-guide/application-live-state/) | Incubating |
| [Plan Preview](/docs/user-guide/plan-preview) | Alpha |

### Lambda Deployment

| Feature | Phase |
|-|-|
| Quick Sync Deployment | Beta |
| Deployment with the Specified Pipeline | Beta |
| Automated Rollback | Beta |
| [Automated Configuration Drift Detection](/docs/user-guide/configuration-drift-detection/) | Incubating |
| [Application Live State](/docs/user-guide/application-live-state/) | Incubating |
| [Plan Preview](/docs/user-guide/plan-preview) | Alpha |

### Amazon ECS Deployment

| Feature | Phase |
|-|-|
| Quick Sync Deployment | Alpha |
| Deployment with the Specified Pipeline | Alpha |
| Automated Rollback | Alpha |
| [Automated Configuration Drift Detection](/docs/user-guide/configuration-drift-detection/) | Incubating |
| [Application Live State](/docs/user-guide/application-live-state/) | Incubating |
| Support [AWS App Mesh](https://aws.amazon.com/app-mesh/) | Incubating |
| [Plan Preview](/docs/user-guide/plan-preview) | Alpha |

### Piped's Core

| Feature | Phase |
|-|-|
| [Wait Stage](/docs/user-guide/adding-a-wait-stage/) | Beta |
| [Wait Manual Approval Stage](/docs/user-guide/adding-a-manual-approval/) | Beta |
| [Notification](/docs/operator-manual/piped/configuring-notifications/) to Slack | Beta |
| [Notification](/docs/operator-manual/piped/configuring-notifications/) to Webhook | Incubating |
| [Secrets Management](/docs/user-guide/secret-management/) | Beta |
| [Event Watcher](/docs/user-guide/event-watcher/) | Alpha |
| [Command-line tool (pipectl) and API for external services](/docs/user-guide/command-line-tool/) | Beta |
| Support executing custom stage | Incubating |
| [ADA](/docs/user-guide/automated-deployment-analysis/) (Automated Deployment Analysis) by Prometheus metrics | Alpha |
| [ADA](/docs/user-guide/automated-deployment-analysis/) by Datadog metrics | Alpha |
| [ADA](/docs/user-guide/automated-deployment-analysis/) by Stackdriver metrics | Incubating |
| [ADA](/docs/user-guide/automated-deployment-analysis/) by Stackdriver log | Incubating |
| [ADA](/docs/user-guide/automated-deployment-analysis/) by CloudWatch metrics | Incubating |
| [ADA](/docs/user-guide/automated-deployment-analysis/) by CloudWatch log | Incubating |
| [ADA](/docs/user-guide/automated-deployment-analysis/) by HTTP request (smoke test...) | Incubating |

### ControlPlane's Core

| Feature | Phase |
|-|-|
| Project/Environment/Piped/Application/Deployment Management | Beta |
| Rendering Deployment Pipeline in Realtime | Beta |
| Canceling a Deployment from Web | Beta |
| Triggering a Sync/Deployment from Web | Beta |
| Authentication by Username/Password for Static Admin | Beta |
| GitHub & GitHub Enterprise SSO | Beta |
| Google SSO | Incubating |
| Data Store - Support GCP [Firestore](https://cloud.google.com/firestore) | Beta |
| Data Store - Support [MySQL v8.0](https://www.mysql.com/) | Alpha |
| File Store - Support GCP [GCS](https://cloud.google.com/storage) | Beta |
| File Store - Support AWS [S3](https://aws.amazon.com/s3/) | Alpha |
| File Store - Support [Minio](https://github.com/minio/minio) | Alpha |
| [Insights](/docs/user-guide/insights/) shows delivery performance | Incubating |
| Deployment Chain - Allow rolling out to multiple clusters gradually or promoting across environments | Incubating |
| Collecting piped's metrics and enabling their dashboards | Incubating |
