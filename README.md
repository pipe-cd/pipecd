# PipeCD

Continuous Delivery for Declarative Kubernetes Application and Infrastructure

## Proposal

https://docs.google.com/document/d/1Z3NqnsxgraD9f55F0TK6e4oLV4296Hd7xUb5GxwdaJQ

## Status

This project is under **PROTOTYPE** development phase.

## Overview

This project aims to explore and develop a unified delivery infrastructure for CA projects.

![](https://github.com/pipe-cd/pipe/blob/master/docs/static/images/architecture-overview.png)

**Powerful**
- Unifed Deployment System: kubernetes (plain-yaml, helm, kustomize), terraform, lambda, cloudrun...
- Progressive Deployment Strategies: canary, bluegreen, rolling update
- Automated Analysis: by metrics, log, smoke test...
- Automated Rollback
- Automated Configuration Drift Detection
- Insights shows Delivery Perfomance

**Easy to Use**
- Operations by Pull Request: scale, rollout, rollback by PR
- Realtime Visualization of application state
- Deployment Pipeline to see what is happenning
- Intuitive UI

**Easy to Operate**
- Two seperate components: single binary `piped` and `control-plane`
- `piped` can be run in a Kubernetes cluster, a single VM or even a local machine
- Easy to operate multi-tenancy, multi-cluster
- Security: your credentials are not exposed outside of your cluster

## License

Apache License 2.0, see [LICENSE](https://github.com/pipe-cd/pipe/blob/master/README.md).
