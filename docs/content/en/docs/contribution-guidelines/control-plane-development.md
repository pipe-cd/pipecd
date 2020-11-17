---
title: "Control plane development"
linkTitle: "Control plane development"
weight: 5
description: >
  This page describes where to find control-plane source code and how to run it locally for debugging.
---

## Source code structure

- [cmd/server](https://github.com/pipe-cd/pipe/tree/master/cmd/server): entrypoint for binary of control-plane server.

- [pkg/app/api](https://github.com/pipe-cd/pipe/tree/master/pkg/app/api): contains source code for control-plane api.
- [pkg/app/web](https://github.com/pipe-cd/pipe/tree/master/pkg/app/web): contains source code for control-plane web.
- [pkg](https://github.com/pipe-cd/pipe/tree/master/pkg): contains shared source code for all components of both `piped` and `control-plane`.

## How to run server locally

### Running up server

Prepare a ControlPlane configuration file as described at [Installation](https://pipecd.dev/docs/operator-manual/control-plane/installation/) and start running by the following command:

``` console
bazelisk run //cmd/server:server -- server \
--config-file=absolute-path-to-control-plane-config.yaml \
--encryption-key-file=absolute-path-to-a-random-key-file
```

Because we are using grpc-web for communicating between web-client and server, so we may need a **local Envoy instance**.

### Integrating with Envoy

You can install [Envoy](https://www.envoyproxy.io/docs/envoy/latest/start/install) locally or running it on [Docker](https://docs.docker.com/get-docker/).

Prepare for Envoy configuration file as below and save as `pipe-envoy-config.yaml`:

**This config file syntax is running on _Envoy_ version _1.10.0_**

```yaml
admin:
  access_log_path: /dev/stdout
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 9095

static_resources:
  listeners:
  - name: ingress
    address:
      socket_address:
        address: 0.0.0.0
        port_value: 9090
    filter_chains:
    - filters:
      - name: envoy.http_connection_manager
        config:
          access_log:
            name: envoy.file_access_log
            config:
              path: /dev/stdout
            filter:
              not_health_check_filter: {}
          codec_type: auto
          idle_timeout: 600s
          stat_prefix: ingress_http
          http_filters:
          - name: envoy.grpc_web
          - name: envoy.cors # If cors is enable
          - name: envoy.router
          route_config:
            virtual_hosts:
            - name: envoy
              domains:
                - '*'
              cors:
                allow_origin:
                - http://localhost:9090
                allow_methods: GET, PUT, DELETE, POST, OPTIONS
                allow_headers: keep-alive,user-agent,cache-control,content-type,content-transfer-encoding,custom-header-1,x-accept-content-transfer-encoding,x-accept-response-streaming,x-user-agent,x-grpc-web,grpc-timeout,authorization
                allow_credentials: true
                max_age: "1728000"
                expose_headers: custom-header-1,grpc-status,grpc-message
              routes:
                - match:
                    prefix: /pipe.api.service.pipedservice.PipedService/
                    grpc:
                  route:
                    cluster: server-piped-api
                - match:
                    prefix: /pipe.api.service.webservice.WebService/
                    grpc:
                  route:
                    cluster: server-web-api
                - match:
                    prefix: /
                  route:
                    cluster: server-http
  clusters:
    - name: server-piped-api
      http2_protocol_options: {}
      connect_timeout: 0.25s
      type: strict_dns
      lb_policy: round_robin
      load_assignment:
        cluster_name: server-piped-api
        endpoints:
        - lb_endpoints:
          - endpoint:
              address:
                socket_address:
                  address: localhost
                  port_value: 9080
    - name: server-web-api
      http2_protocol_options: {}
      connect_timeout: 0.25s
      type: strict_dns
      lb_policy: round_robin
      load_assignment:
        cluster_name: server-web-api
        endpoints:
        - lb_endpoints:
          - endpoint:
              address:
                socket_address:
                  address: localhost
                  port_value: 9081
    - name: server-http
      connect_timeout: 0.25s
      type: strict_dns
      lb_policy: round_robin
      load_assignment:
        cluster_name: server-http
        endpoints:
        - lb_endpoints:
          - endpoint:
              address:
                socket_address:
                  address: localhost
                  port_value: 9082
```

#### Using local Envoy

```
envoy -c pipe-envoy-config.yaml 
```

#### Using Envoy running on Docker

1. Create simple Dockerfile.

```Dockerfile
FROM envoyproxy/envoy:v1.10.0
COPY absolute-path-to-pipe-envoy-config.yaml /etc/envoy/envoy.yaml
CMD /usr/local/bin/envoy -c /etc/envoy/envoy.yaml
```

2. Run Docker command build for this `Dockerfile`.

```
docker build -t envoy:v1 .
```

3. And now you can execute it with:

```
docker run -d --net="host" --name web-envoy -p 9095:9095 -p 9090:9090 web-envoy:1.0
```

Then go to `http://localhost:9090` on your browser to access PipeCD's web.