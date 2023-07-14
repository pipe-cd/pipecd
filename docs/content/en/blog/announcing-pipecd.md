---
date: 2020-10-06
title: "Announcing PipeCD"
linkTitle: "Announcing PipeCD"
weight: 1000
description: "Continuous Delivery for declarative Kubernetes, Serverless and Infrastructure applications"
author: Le Van Nghia ([@nghialv](https://twitter.com/nghialv2607))
images: 
 - images/deployment-details.png
---

Today we are excited to announce the open-source availability of PipeCD: a continuous delivery system for declarative Kubernetes, Serverless, and Infrastructure applications.
PipeCD aims to provide a unified CD solution for multiple application kinds on multi-cloud that empowers engineers to deploy faster with more confidence.
It is also available as a GitOps tool that enables doing deployment operations by pull request on Git.

<br>

![](/images/deployment-details.png)
<p style="text-align: center;">
Deployment Details Page
</p>
<br>

### Background

As one of our Developer Productivity team's missions, we aim to empower engineers to deploy their services faster, more frequently with reliability.
Martin Fowler, in his book _Continuous Delivery_, points out that "The biggest risk to any software effort is that you end up building something that isn't useful. The earlier and more frequently you get working software in front of real users, the quicker you get feedback to find out how valuable it really is."

Recently, with the popularity of cloud services and the container technology, engineers have even more options in choosing the infrastructure model, the cloud services, which are most suitable for their team's requirements.
At [CyberAgent](https://www.cyberagent.co.jp/en/), we have a large number of services from many teams where each team can have a different infrastructure model and a different cloud service. Some big projects also deploy their services on multi-cloud.
This diversification leads to a problem that we are facing, lacking a robust CD system for all teams.

So we decided to create a new CD system that provides a unified interface for many application kinds to improve the developer experience.

### Key Features

While designing PipeCD, we focused on the following [4 key features](https://pipecd.dev) with the aim of creating a CD system that provides a good experience for both developers and operators.

**Visibility**

Visibility is one of the most requested factors we received when surveying engineers in our company.
Visibility is the ability to see what's going on in the cluster, the ability to see how each component of the application had been deployed, the ability to know quickly why the deployment was failed.
Visibility for a team leader is the ability to know the delivery performance of the team and what metrics should be improved.
With PipeCD, we always strive to maximize the visibility for engineers, operators as well as team leaders. Currently, it includes:

- Deployment pipeline UI shows clarify what is happening
- Separate logs viewer for each individual deployment
- Realtime visualization of application component and state
- Deployment notifications to slack, webhook endpoints
- Insights show the delivery performance

In addition, the entire state of the service is managed through Git, so you can view the whole state of the cluster and all audit logs provided by Git.

**Automation**

Automation reduces or removes repetitive overhead of frequent releases. So maximizing automation helps to minimize human error during the deployment process, as well as reduce the amount of work engineers need to do.
PipeCD has the following automated functionalities:

- Automated deployment analysis based on metrics, logs, emitted requests
- Automatically roll back to the previous state as soon as analysis or a pipeline stage fails
- Automatically detect configuration drift to notify and render the changes
- Automatically watch and detect the new container images to deploy

**Secure**

The CD system often carries a lot of credentials needed to access the cluster and the necessary services of the teams.
Ensuring the safety of the teams is always on our top priority.
So while designing the PipeCD, we decided not to store those credentials in a central place. Instead of that, all user's credentials always stay inside their clusters.

- Support single sign-on and role-based access control
- Credentials are not exposed outside the cluster and not saved in the control-plane
- Piped makes only outbound requests and can run inside a restricted network

**Multi-provider & Multi-Tenancy**

Multi-provider means supporting multiple cloud services, multiple container registries, multiple monitoring services for doing deployment analysis.
You can use PipeCD to deploy your Kubernetes applications, CloudRun, AWS Lambda application and Terraform application.
It also supports doing progressive delivery with canary and blue-green strategy.

- Easy to operate multi-cluster, multi-tenancy by separating control-plane and piped
- Support multiple analysis providers including Prometheus, Datadog, Stackdriver, and more
- Support multiple application kinds on multi-cloud including Kubernetes, Terraform, Cloud Run, AWS Lambda

<br>

While designing PipeCD, we simplified its architecture by minimizing the number of components, so you do not have to install many things to enable all features.
In addition, PipeCD also supports storing data in several fully-managed services to minimize the operating cost.

Currently, we have completed the basic features and many of the features are in the alpha status. We are working hard to increase the stability and planning to release a stable version in the next months.

### Community

PipeCD team hopes to receive a warm welcome and feedback from the open-source community.
We value every contribution and invite you to join us on GitHub, Slack and Twitter.

- Visit our website and documentation at [https://pipecd.dev](https://pipecd.dev)
- Check out the code at [https://github.com/pipe-cd/pipecd](https://github.com/pipe-cd/pipecd) or explore the [examples](https://pipecd.dev/docs/examples/)
- Join us on Slack [@cloud-native/pipecd](https://slack.cncf.io) to chat with other developers
- Follow us on Twitter [@pipecd_dev](https://twitter.com/pipecd_dev) to get the latest news

PipeCD team is hiring engineers/interns to work on PipeCD. Please contact us if you are interested.

### Thanks

Finally, we would like to thank the existing open-source CD projects like Spinnaker, FluxCD, ArgoCD... PipeCD has been built on many great ideas from those great projects.
PipeCD team would also like to thank CyberAgent's engineers and collaborators from other companies, who have sent us so much valuable feedback throughout the development process to this day.
