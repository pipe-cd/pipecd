<p align="center">
  <img src="https://github.com/pipe-cd/pipe/blob/master/docs/static/images/logo.png" width="180"/>
</p>

<p align="center">
  Continuous Delivery for Declarative Kubernetes, Serverless Application and Infrastructure
  <br/>
  <a href="https://pipecd.dev"><strong>Explore PipeCD docs »</strong></a>
</p>

#

## Overview

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
