# PipeCD

Continuous Delivery for Declarative Kubernetes, Serverless Application and Infrastructure

## Proposal

https://docs.google.com/document/d/1Z3NqnsxgraD9f55F0TK6e4oLV4296Hd7xUb5GxwdaJQ

## Status

This project is under **PROTOTYPE** development phase.

## Overview

This project aims to explore and develop a unified delivery infrastructure for CA projects.

![](https://github.com/pipe-cd/pipe/blob/master/docs/static/images/architecture-overview.png)

**Powerful**
- Unified Deployment System: kubernetes (plain-yaml, helm, kustomize), terraform, lambda, cloudrun...
- Progressive Deployment Strategies: canary, bluegreen, rolling update
- Automated Analysis: by metrics, log, smoke test...
- Automated Rollback
- Automated Configuration Drift Detection
- Insights shows Delivery Perfomance
- Support Webhook and Slack notifications

**Easy to Use**
- Operations by Pull Request: scale, rollout, rollback by PR
- Realtime Visualization of application state
- Deployment Pipeline to see what is happenning
- Intuitive UI

**Easy to Operate**
- Two seperate components: single binary `piped` and `control-plane`
- `piped` can be run in a Kubernetes cluster, a single VM or even a local machine
- Easy to operate multi-tenancy, multi-cluster

**Safety and Security**
- Support single sign-on (SSO) and role-based access control (RBAC)
- Your credentials are not exposed outside your cluster and not saved in control-plane

## License

Apache License 2.0, see [LICENSE](https://github.com/pipe-cd/pipe/blob/master/README.md).
