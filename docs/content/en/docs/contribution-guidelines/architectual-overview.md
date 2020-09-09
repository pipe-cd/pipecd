---
title: "Architectual overview"
linkTitle: "Architectual overview"
weight: 3
description: >
  This page describes the architecture of PipeCD.
---

> TBA

![](/images/architecture-overview.png)
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
- `api`: A service to provide api for external service like web and hook requests.
- `web`: A service for serving static files for web.
- `cache`: A redis cache service for caching internal data.
- `datastore`: Data storage for storing deployment, application data. This can be a fully-managed service such as `Firestore`, `DynamoDB`... or a self-managed such as `MongoDB`.
- `filestore`: File storage for storing logs, application states. This can a fully-managed service such as `GCS`, `S3`... or a self-managed service such as `Minio`.
