---
title: "Architecture overview"
linkTitle: "Architecture overview"
weight: 1
description: >
  This page describes the architecture of control plane.
---

![](/images/control-plane-components.png)
<p style="text-align: center;">
Component Architecture
</p>

The control plane is a centralized part of PipeCD. It contains several services as below to manage the application, deployment data and handle all requests from `piped`s and web clients:

##### Server

`server` handles all incoming gRPC requests from `piped`s, web clients, incoming HTTP requests such as auth callback from third party services.
It also serves all web assets including HTML, JS, CSS...
This service can be easily scaled by updating the pod number.

##### Cache

`cache` is a single pod service for caching internal data used by `server` service. Currently, this `cache` service is powered by `redis`.
You can configure the control plane to use a fully-managed redis cache service instead of launching a cache pod in your cluster.

##### Ops

`ops` is a single pod service for operating PipeCD owner's tasks.
For example, it provides an internal web page for adding and managing projects; it periodically removes the old data; it collects and saves the deployment insights.

##### Data Store

`Data store` is a storage for storing model data such as applications and deployments. This can be a fully-managed service such as GCP [Firestore](https://cloud.google.com/firestore), GCP [Cloud SQL](https://cloud.google.com/sql) or AWS [RDS](https://aws.amazon.com/rds/) (currently we choose [MySQL v8](https://www.mysql.com/) as supported relational data store). You can also configure the control plane to use a self-managed MySQL server.
When installing the control plane, you have to choose one of the provided data store services.

##### File Store

`File store` is a storage for storing stage logs, application live states. This can be a fully-managed service such as GCP [GCS](https://cloud.google.com/storage), AWS [S3](https://aws.amazon.com/s3/), or a self-managed service such as [Minio](https://github.com/minio/minio).
When installing the control plane, you have to choose one of the provided file store services.
