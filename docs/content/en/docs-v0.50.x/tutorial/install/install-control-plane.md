---
title: "2. Install Control Plane"
linkTitle: "2. Install Control Plane"
weight: 2
description: >
  Install a Control Plane on local by docker-compose.
---

# 2. Install Control Plane

In this page, you will install a Control Plane on local by `docker-compose`.

## Installation

1. Execute the following command on your `src/install/control-plane/docker-compose.yaml`.
    ```sh
    docker-compose up
    ```

    After the command, you will see logs as below.
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

## Confirmation

1. Access the console running on [http://localhost:8080](http://localhost:8080)
2. Enter the following value and click `CONTINUE`.
   - `Project Name`: `tutorial`

    ![signin-project](/images/tutorial/install/signin-project.png)

3. Enter the following values and click `LOGIN`.
   - `Username`: `hello-pipecd`
   - `Password`: `hello-pipecd`

    ![signin-user](/images/tutorial/install/signin-user.png)

4. If successful, the following page will appear.

    ![applications-page](/images/tutorial/install/applications.png)


## See Also

- [Architecture Overview](https://pipecd.dev/docs/user-guide/managing-controlplane/architecture-overview/)
- [Managing Control Plane](https://pipecd.dev/docs/user-guide/managing-controlplane/)
- [Installing Control Plane on Kubernetes](https://pipecd.dev/docs/installation/install-control-plane/installing-controlplane-on-k8s/)

---

[Next: 3. Install Piped >](../install-piped/)

[< Previous: 1. Setup Git Repository](../setup-git-repository/)
