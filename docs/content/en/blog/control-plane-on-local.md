---
date: 2024-03-12
title: "Control Plane on local by docker-compose"
linkTitle: "Control Plane on local by docker-compose"
weight: 989
description: ""
author: Tetsuya Kikuchi ([@t-kikuc](https://github.com/t-kikuc))
---

Currently, you can deploy and operate a PipeCD Control Plane on a Kubernetes cluster or [on Amazon ECS](./control-plane-on-ecs.md).
However, some developers would like to build a Control Plane more easily for introduction or development.
This blog shows you how to install a PipeCD Control Plane on local machine easily.

Intended readers are those who would like to:
- begin using PipeCD and experiment its features instantly, including Control Plane and Piped
- develop PipeCD and debug your codes easily

### Architecture

The general architecture of PipeCD Control Plane is as below.
![](/images/control-plane-components.png)

NOTE: See [Architecture Overview](docs/user-guide/managing-controlplane/architecture-overview/) doc for details.

In this blog, you will build the Control Plane by these containers:

- Server: pipecd server
- Ops: pipecd ops
- Cache: Redis
- Data Store: MySQL
- File Store: MinIO

### Prerequisites

- You have [Docker Engine](https://docs.docker.com/engine/)
- You have [docker-compose](https://docs.docker.jp/compose/install.html)

### Installation of Control Plane

You can install a Control Plane by just executing this command:

```sh
docker-compose up
```

After executing the above command, you will see logs like below if success.

```log
pipecd-server-1   | successfully loaded control-plane configuration
pipecd-server-1   | successfully connected to file store
pipecd-server-1   | successfully connected to data store
pipecd-server-1   | grpc server will be run without tls
pipecd-server-1   | grpc server will be run without tls
pipecd-server-1   | grpc server is running on [::]:9080
pipecd-server-1   | grpc server is running on [::]:9083
pipecd-server-1   | grpc server will be run without tls
pipecd-server-1   | admin server is running on 9085
pipecd-server-1   | grpc server is running on [::]:9081
pipecd-server-1   | start running http server on :9082
```

### Access the console of the Control Plane to confirm

1. Access http://localhost:8080 on your web browser.

2. Enter a value as below and `CONTINUE`.
   - `Project Name`: `control-plane-local`

3. Enter values as below and `LOGIN`.
   - `Username`: `hello-pipecd`
   - `Password`: `hello-pipecd`

4. You will see the applications page. Success!

![](/images/control-plane-local-console.png)

### Clean up

To clean up the Control Plane, execute the below command.

```sh
docker-compose down
```

NOTE: By following commands instead of avobe one, you can keep data such as Piped or applications even after restarting/updating the server component.

```sh
# Restart only the server component.
docker-compose rm -fsv pipecd-server
docker-compose up pipecd-server
```
