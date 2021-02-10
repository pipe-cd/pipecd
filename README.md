<p align="center">
  <img src="https://github.com/pipe-cd/pipe/blob/master/docs/static/images/logo.png" width="180"/>
</p>

<p align="center">
  Continuous Delivery for Declarative Kubernetes, Serverless and Infrastructure Applications
  <br/>
  <a href="https://pipecd.dev"><strong>Explore PipeCD docs »</strong></a>
</p>

#

![](https://github.com/pipe-cd/pipe/blob/master/docs/static/images/deployment-details.png)

### Overview

PipeCD provides a unified continuous delivery solution for multiple application kinds on multi-cloud that empowers engineers to deploy faster with more confidence, a GitOps tool that enables doing deployment operations by pull request on Git.

![](https://github.com/pipe-cd/pipe/blob/master/docs/static/images/architecture-overview.png)

**Visibility**
- Deployment pipeline UI shows clarify what is happening
- Separate logs viewer for each individual deployment
- Realtime visualization of application state
- Deployment notifications to slack, webhook endpoints
- Insights show metrics like lead time, deployment frequency, MTTR and change failure rate to measure delivery performance

**Automation**
- Automated deployment analysis to measure deployment impact based on metrics, logs, emitted requests
- Automatically roll back to the previous state as soon as analysis or a pipeline stage fails
- Automatically detect configuration drift to notify and render the changes
- Automatically trigger a new deployment when a defined event has occurred (e.g. container image pushed, helm chart published, etc)

**Safety and Security**
- Support single sign-on and role-based access control
- Credentials are not exposed outside the cluster and not saved in the control-plane
- Piped makes only outbound requests and can run inside a restricted network
- Built-in secrets management

**Multi-provider & Multi-Tenancy**
- Support multiple application kinds on multi-cloud including Kubernetes, Terraform, Cloud Run, AWS Lambda
- Support multiple analysis providers including Prometheus, Datadog, Stackdriver, and more
- Easy to operate multi-cluster, multi-tenancy by separating control-plane and piped

#

### Contributing

We'd love you to join us! Please see the [Contributor Guide](https://pipecd.dev/docs/contribution-guidelines/).

#

### License

Apache License 2.0, see [LICENSE](https://github.com/pipe-cd/pipe/blob/master/LICENSE).

