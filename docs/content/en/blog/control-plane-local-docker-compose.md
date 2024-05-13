---
date: 2024-03-14
title: "Control Plane on local by docker-compose"
linkTitle: "Control Plane on local by docker-compose"
weight: 989
description: "This blog shows you how to install a PipeCD Control Plane on local machine easily."
author: Tetsuya Kikuchi ([@t-kikuc](https://github.com/t-kikuc))
---

Currently, you can deploy and operate a PipeCD Control Plane on a Kubernetes cluster or [on Amazon ECS](/blog/2023/02/07/pipecd-best-practice-02-control-plane-on-ecs/).
However, some developers would like to build a Control Plane more easily for introduction or development.
This blog shows you how to install a PipeCD Control Plane on local machine easily.

This blog is for those who would like to:
- begin using PipeCD and experiment its features instantly, including Control Plane and Piped

### Architecture

The general architecture of PipeCD Control Plane is as below.
![](/images/control-plane-components.png)

> Note: See [Architecture Overview](/docs/user-guide/managing-controlplane/architecture-overview/) doc for details.
In this blog, you will build a Control Plane by these containers:

- Server: pipecd server
- Ops: pipecd ops
- Cache: Redis
- Data Store: MySQL
- File Store: MinIO

### Installation of Control Plane

1. Get the demo codes from https://github.com/pipe-cd/demo/blob/main/control_plane/docker-compose/

2. Execute the below command for docker-compose.yaml you got in [1.].

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

### To confirm: Access the console of the Control Plane

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

> Note: By following commands instead of above one, you can keep data such as Piped or applications on the Control Plane even after restarting/updating the server component.

```sh
# Restart only the server component.
docker-compose rm -fsv pipecd-server
docker-compose up pipecd-server
```
