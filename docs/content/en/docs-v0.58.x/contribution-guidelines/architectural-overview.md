---
title: "Architectural overview"
linkTitle: "Architectural overview"
weight: 3
description: >
  This page describes the architecture of PipeCD.
---

![](/images/architecture-overview.png)
<p style="text-align: center;">
Component Architecture
</p>

### Piped

A single binary component runs in your cluster, your local network to handle the deployment tasks.
It can be run inside a Kubernetes cluster by simply starting a Pod or a Deployment.
This component is designed to be stateless, so it can also be run in a single VM or even your local machine.

### Control Plane

A centralized component manages deployment data and provides gRPC API for connecting `piped`s as well as all web-functionalities of PipeCD such as
authentication, showing deployment list/details, application list/details, delivery insights...

Control Plane contains the following components:
- `server`: a service to provide api for piped, web and serve static assets for web.
- `ops`: a service to provide administrative features for Control Plane owner like adding/managing projects.
- `cache`: a redis cache service for caching internal data.
- `datastore`: data storage for storing deployment, application data
  - this can be a fully-managed service such as `Firestore`, `Cloud SQL`...
  - or a self-managed such as `MySQL`
- `filestore`: file storage for storing logs, application states
  - this can a fully-managed service such as `GCS`, `S3`...
  - or a self-managed service such as `Minio`

For more information, see [Architecture overview of Control Plane](../../user-guide/managing-controlplane/architecture-overview/).
