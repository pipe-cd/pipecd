---
title: "Canary deployment for Kubernetes app with Istio"
linkTitle: "Canary k8s app with Istio"
weight: 1
description: >
  How to enable canary deployment for Kubernetes application with Istio.
---

> Canary release is a technique to reduce the risk of introducing a new software version in production by slowly rolling out the change to a small subset of users before rolling it out to the entire infrastructure and making it available to everybody.
> -- <cite>[martinfowler.com/canaryrelease](https://martinfowler.com/bliki/CanaryRelease.html)</cite>

With Istio, we can accomplish this goal by configuring a sequence of rules that route a percentage of traffic to each [variant](../../managing-application/defining-app-configuration/kubernetes/#sync-with-the-specified-pipeline) of the application.
And with PipeCD, you can enable and automate the canary strategy for your Kubernetes application even easier.

In this guide, we will show you how to configure the application configuration file to send 10% of traffic to the new version and keep 90% to the primary variant. Then after waiting for manual approval, you will complete the migration by sending 100% of traffic to the new version.

Complete source code for this example is hosted in [pipe-cd/examples](https://github.com/pipe-cd/examples/tree/master/kubernetes/mesh-istio-canary) repository.

## Before you begin

- Add a new Kubernetes application by following the instructions in [this guide](../../managing-application/adding-an-application/)
- Ensure having `pipecd.dev/variant: primary` [label](https://github.com/pipe-cd/examples/blob/master/kubernetes/mesh-istio-canary/deployment.yaml#L17) and [selector](https://github.com/pipe-cd/examples/blob/master/kubernetes/mesh-istio-canary/deployment.yaml#L12) in the workload template
- Ensure having at least one Istio's `DestinationRule` and defining the needed subsets (`primary` and `canary`) with `pipecd.dev/variant` label

``` yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: mesh-istio-canary
spec:
  host: mesh-istio-canary.default.svc.cluster.local
  subsets:
  - name: primary
    labels:
      pipecd.dev/variant: primary
  - name: canary
    labels:
      pipecd.dev/variant: canary
```

- Ensure having at least one Istio's `VirtualService` manifest and all traffic is routed to the `primary`

``` yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: mesh-istio-canary
spec:
  hosts:
    - mesh-istio-canary.pipecd.dev
  gateways:
    - mesh-istio-canary
  http:
    - route:
      - destination:
          host: mesh-istio-canary.default.svc.cluster.local
          subset: primary
        weight: 100
```

## Enabling canary strategy

- Add the following application configuration file into the application directory in Git.

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 50%
      - name: K8S_TRAFFIC_ROUTING
        with:
          canary: 10
          primary: 90
      - name: WAIT_APPROVAL
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_TRAFFIC_ROUTING
        with:
          primary: 100
      - name: K8S_CANARY_CLEAN
  trafficRouting:
    method: istio
    istio:
      host: mesh-istio-canary.default.svc.cluster.local
```

- Send a PR to update the container image version in the Deployment manifest and merge it to trigger a new deployment. PipeCD will plan the deployment with the specified canary strategy.

![](/images/example-canary-kubernetes-istio.png)
<p style="text-align: center;">
Deployment Details Page
</p>

- Now you have an automated canary deployment for your application. 🎉

## Understanding what happened

In this example, you configured the application configuration file to migrate traffic from an old to a new version of the application using Istio's weighted routing feature.

- Stage 1: `K8S_CANARY_ROLLOUT` ensures that the workloads of canary variant (new version) should be deployed. But at this time, they still handle nothing, all traffic are handled by workloads of primary variant.
The number of workloads (e.g. pod) for canary variant is configured to be 50% of the replicas number of primary varant.

![](/images/example-canary-kubernetes-istio-stage-1.png)

- Stage 2: `K8S_TRAFFIC_ROUTING` ensures that 10% of traffic should be routed to canary variant and 90% to primary variant. Because the `trafficRouting` is configured to use Istio, PipeCD will find Istio's VirtualService resource of this application to control the traffic percentage.

![](/images/example-canary-kubernetes-istio-stage-2.png)

- Stage 3: `WAIT_APPROVAL` waits for a manual approval from someone in your team.

- Stage 4: `K8S_PRIMARY_ROLLOUT` ensures that all resources of primary variant will be updated to the new version.

![](/images/example-canary-kubernetes-istio-stage-4.png)

- Stage 5: `K8S_TRAFFIC_ROUTING` ensures that all traffic should be routed to primary variant. Now primary variant is running the new version so it means all traffic is handled by the new version.

![](/images/example-canary-kubernetes-istio-stage-5.png)

- Stage 6: `K8S_CANARY_CLEAN` ensures all created resources for canary variant should be destroyed.

![](/images/example-canary-kubernetes-istio-stage-6.png)
