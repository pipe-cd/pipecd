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
|Incubating| Early in design or prototype stage. Not ready for general use.|
|Alpha|Works end to end but may be incomplete or unstable. Backward compatibility is not guaranteed.|
|Beta|Ready for production use and well documented, but may still evolve.|
|Stable|Production-proven, fully supported, and backward compatible.|

## Plugins

<!-- ### Kubernetes Plugin

| Feature | Phase |
|-|-|
| Quick sync deployment | Beta |
| Deployment with a defined pipeline (e.g. canary, analysis) | Beta |
| [Automated rollback](../user-guide/managing-application/rolling-back-a-deployment/) | Beta |
| [Automated configuration drift detection](../user-guide/managing-application/configuration-drift-detection/) | Beta |
| [Application live state](../user-guide/managing-application/application-live-state/) | Beta |
| Prune resources | Alpha |
| Support Helm | Beta |
| Support Kustomize | Beta |
| Support Istio service mesh | Beta |
| Support SMI service mesh | Incubating |
| [Plan preview](../user-guide/plan-preview) | Beta |
| [Manifest attachment](../user-guide/managing-application/manifest-attachment) | Alpha |

### Terraform Plugin

| Feature | Phase |
|-|-|
| Quick sync deployment | Beta |
| Deployment with a defined pipeline (e.g. manual-approval) | Beta |
| [Automated rollback](../user-guide/managing-application/rolling-back-a-deployment/) | Beta |
| [Automated configuration drift detection](../user-guide/managing-application/configuration-drift-detection/) | Alpha |
| [Application live state](../user-guide/managing-application/application-live-state/) | Incubating |
| Prune resources | Incubating |
| [Plan preview](../user-guide/plan-preview) | Beta |
| [Manifest attachment](../user-guide/managing-application/manifest-attachment) | Alpha |

### Cloud Run Plugin

| Feature | Phase |
|-|-|
| Quick sync deployment | Beta |
| Deployment with a defined pipeline (e.g. canary, analysis) | Beta |
| [Automated rollback](../user-guide/managing-application/rolling-back-a-deployment/) | Beta |
| [Automated configuration drift detection](../user-guide/managing-application/configuration-drift-detection/) | Beta |
| [Application live state](../user-guide/managing-application/application-live-state/) | Beta |
| Prune resources | Incubating |
| [Plan preview](../user-guide/plan-preview) | Beta |
| [Manifest attachment](../user-guide/managing-application/manifest-attachment) | Alpha |

Note: These are statuses for Cloud Run service. Cloud Run job has not been supported yet.

### Lambda Plugin

| Feature | Phase |
|-|-|
| Quick sync deployment | Beta |
| Deployment with a defined pipeline (e.g. canary, analysis) | Beta |
| [Automated rollback](../user-guide/managing-application/rolling-back-a-deployment/) | Beta |
| [Automated configuration drift detection](../user-guide/managing-application/configuration-drift-detection/) | Alpha |
| [Application live state](../user-guide/managing-application/application-live-state/) | Alpha |
| Prune resources | Incubating |
| [Plan preview](../user-guide/plan-preview) | Alpha |
| [Manifest attachment](../user-guide/managing-application/manifest-attachment) | Alpha |

### Amazon ECS Plugin

| Feature | Phase |
|-|-|
| Quick sync deployment | Beta |
| Deployment with a defined pipeline (e.g. canary, analysis) | Beta |
| [Automated rollback](../user-guide/managing-application/rolling-back-a-deployment/) | Beta |
| [Automated configuration drift detection](../user-guide/managing-application/configuration-drift-detection/) | Alpha *1 |
| [Application live state](../user-guide/managing-application/application-live-state/) | Alpha *1 |
| Quick sync deployment for [ECS Service Discovery](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/service-discovery.html) | Alpha |
| Deployment with a defined pipeline for [ECS Service Discovery](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/service-discovery.html) | Alpha |
| Prune resources | Incubating |
| [Plan preview](../user-guide/plan-preview) | Alpha |
| [Manifest attachment](../user-guide/managing-application/manifest-attachment) | Alpha |

*1.  Not supported yet for standalone tasks. -->

## Piped

| Feature | Phase |
|-|-|
| [Deployment wait stage](../user-guide/managing-application/customizing-deployment/adding-a-wait-stage/) | Beta |
| [Deployment manual approval stage](../user-guide/managing-application/customizing-deployment/adding-a-manual-approval/) | Beta |
| [Notification](../user-guide/managing-piped/configuring-notifications/) to Slack | Beta |
| [Notification](../user-guide/managing-piped/configuring-notifications/) to external service via webhook | Beta |
| [Secrets management](../user-guide/managing-application/secret-management/) - Storing secrets safely in the Git repository | Beta |
| [Event watcher](../user-guide/event-watcher/) - Updating files in Git automatically for given events | Beta |
| [Pipectl](../user-guide/command-line-tool/) - Command-line tool for interacting with Control Plane | Beta |
| Deployment plugin - Allow executing user-created deployment plugin | Incubating |
| [ADA](../user-guide/managing-application/customizing-deployment/automated-deployment-analysis/) (Automated Deployment Analysis) by Prometheus metrics | Beta |
| [ADA](../user-guide/managing-application/customizing-deployment/automated-deployment-analysis/) by Datadog metrics | Beta |
| [ADA](../user-guide/managing-application/customizing-deployment/automated-deployment-analysis/) by Stackdriver metrics | Incubating |
| [ADA](../user-guide/managing-application/customizing-deployment/automated-deployment-analysis/) by Stackdriver log | Incubating |
| [ADA](../user-guide/managing-application/customizing-deployment/automated-deployment-analysis/) by CloudWatch metrics | Incubating |
| [ADA](../user-guide/managing-application/customizing-deployment/automated-deployment-analysis/) by CloudWatch log | Incubating |
| [ADA](../user-guide/managing-application/customizing-deployment/automated-deployment-analysis/) by HTTP request (smoke test...) | Incubating |
| [Remote upgrade](../user-guide/managing-piped/remote-upgrade-remote-config/#remote-upgrade) - Ability to upgrade Piped from the web console | Beta |
| [Remote config](../user-guide/managing-piped/remote-upgrade-remote-config/#remote-config) - Watch and reload configuration from a remote location such as Git | Beta |

## Control Plane

| Feature | Phase |
|-|-|
| Project/Piped/Application/Deployment management | Beta |
| Rendering deployment pipeline in realtime | Beta |
| Canceling a deployment from console | Beta |
| Triggering a deployment manually from console | Beta |
| RBAC on PipeCD resources such as Application, Piped... | Beta |
| Authentication by username/password for static admin | Beta |
| GitHub & GitHub Enterprise Server SSO | Beta |
| Support GCP [Firestore](https://cloud.google.com/firestore) as data store | Beta |
| Support [MySQL v8.0](https://www.mysql.com/) as data store | Beta |
| Support file store as data store | Alpha - Deprecated (remove soon) |
| Support GCP [GCS](https://cloud.google.com/storage) as file store | Beta |
| Support AWS [S3](https://aws.amazon.com/s3/) as file store | Beta |
| Support [Minio](https://github.com/minio/minio) as file store | Beta |
| [Insights](../user-guide/insights/) - Show the delivery performance of a team or an application | Beta |
| [Deployment Chain](../user-guide/managing-application/deployment-chain/) - Allow rolling out to multiple clusters gradually or promoting across environments | Alpha |
| [Metrics](../user-guide/managing-controlplane/metrics/) - Dashboards for PipeCD and Piped metrics | Beta |

## Pipectl

Check [pipectl](../user-guide/command-line-tool/) docs for available commands.
