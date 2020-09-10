---
title: "Architecture overview"
linkTitle: "Architecture overview"
weight: 1
description: >
  This page describes the architecture of control plane.
---

> TBA

![](/images/control-plane-components.png)
<p style="text-align: center;">
Component Architecture
</p>

### Piped

A single binary component that you run in your cluster, your local network to handle the deployment tasks.
It can be run inside a Kubernetes cluster by simply starting a Pod or a Deployment.
This component is designed to be stateless so it can also be run in a single VM or even your local machine.

### Control Plane

A centralized component that manages deployment data and provides gPRC API for connecting `piped`s as well as all web-functionalities of PipeCD such as
authentication, showing deployment list/details, application list/details, delivery insights...

Control Plane contains the following components:
- `api`: a service to provide api for piped, web and hook requests.
- `web`: a service to serve static files for web.
- `cache`: a redis cache service for caching internal data.
- `datastore`: data storage for storing deployment, application data
  - this can be a fully-managed service such as `Firestore`, `DynamoDB`...
  - or a self-managed such as `MongoDB`
- `filestore`: file storage for storing logs, application states
  - this can a fully-managed service such as `GCS`, `S3`...
  - or a self-managed service such as `Minio`
