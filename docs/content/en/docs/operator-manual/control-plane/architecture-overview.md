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

##### API

`api` is the most important service in the control plane. It handles all incoming requests by providing a gRPC server for handling incoming gRPC requests from `piped`s, web clients, and providing an HTTP server for handling incoming HTTP requests from third party services such as auth callback, webhook requests. This service can be easily scaled by updating the pod number.

##### Web

`web` is a service providing an HTTP server for serving all static assets for web rendering such as HTML, JS, CSS... This service can be easily scaled by updating the pod number.

##### Cache

`cache` is a single pod service for caching internal data used by `api` service. Currently, this `cache` service is using the `redis` docker image. You can configure the control plane to use a fully-managed redis cache service instead of launching a cache pod in your cluster.

##### Ops

`ops` is a single pod service for operating PipeCD owner's tasks.
For example, it provides an internal web page for adding and managing projects; it periodically removes the old data; it collects and saves the deployment insights.

##### Data Store

`Data store` is a storage for storing the application, deployment data. This can be a fully-managed service such as GCP [Firestore](https://cloud.google.com/firestore), AWS [DynamoDB](https://aws.amazon.com/dynamodb/), or a self-managed service such as [MongoDB](https://www.mongodb.com/). When installing the control plane, you have to choose one of the provided data store services.

##### File Store

`File store` is a storage for storing stage logs, application live states. This can be a fully-managed service such as GCP [GCS](https://cloud.google.com/storage), AWS [S3](https://aws.amazon.com/s3/), or a self-managed service such as [Minio](https://github.com/minio/minio). When installing the control plane, you have to choose one of the provided data store services.
