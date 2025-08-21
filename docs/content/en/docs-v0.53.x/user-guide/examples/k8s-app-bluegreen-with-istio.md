---
title: "BlueGreen deployment for Kubernetes app with Istio"
linkTitle: "BlueGreen k8s app with Istio"
weight: 2
description: >
  How to enable blue-green deployment for Kubernetes application with Istio.
---

Similar to [canary deployment](../k8s-app-canary-with-istio/), PipeCD allows you to enable and automate the blue-green deployment strategy for your application based on Istio's weighted routing feature.

In both canary and blue-green strategies, the old version and the new version of the application get deployed at the same time.
But while the canary strategy slowly routes the traffic to the new version, the blue-green strategy quickly routes all traffic to one of the versions.

In this guide, we will show you how to configure the application configuration file to apply the blue-green strategy.

Complete source code for this example is hosted in [pipe-cd/examples](https://github.com/pipe-cd/examples/tree/master/kubernetes/mesh-istio-bluegreen) repository.

## Before you begin

- Add a new Kubernetes application by following the instructions in [this guide](../../managing-application/adding-an-application/)
- Ensure having `pipecd.dev/variant: primary` [label](https://github.com/pipe-cd/examples/blob/master/kubernetes/mesh-istio-bluegreen/deployment.yaml#L17) and [selector](https://github.com/pipe-cd/examples/blob/master/kubernetes/mesh-istio-bluegreen/deployment.yaml#L12) in the workload template
- Ensure having at least one Istio's `DestinationRule` and defining the needed subsets (`primary` and `canary`) with `pipecd.dev/variant` label

``` yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: mesh-istio-bluegreen
spec:
  host: mesh-istio-bluegreen
  subsets:
  - name: primary
    labels:
      pipecd.dev/variant: primary
  - name: canary
    labels:
      pipecd.dev/variant: canary
  trafficPolicy:
    tls:
      mode: ISTIO_MUTUAL
```

- Ensure having at least one Istio's `VirtualService` manifest and all traffic is routed to the `primary`

``` yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: mesh-istio-bluegreen
spec:
  hosts:
    - mesh-istio-bluegreen.pipecd.dev
  gateways:
    - mesh-istio-bluegreen
  http:
    - route:
      - destination:
          host: mesh-istio-bluegreen
          subset: primary
        weight: 100
```

## Enabling blue-green strategy

- Add the following application configuration file into the application directory in the Git repository.

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 100%
      - name: K8S_TRAFFIC_ROUTING
        with:
          all: canary
      - name: WAIT_APPROVAL
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_TRAFFIC_ROUTING
        with:
          all: primary
      - name: K8S_CANARY_CLEAN
  trafficRouting:
    method: istio
    istio:
      host: mesh-istio-bluegreen
```

- Send a PR to update the container image version in the Deployment manifest and merge it to trigger a new deployment. PipeCD will plan the deployment with the specified blue-green strategy.

![](/images/example-bluegreen-kubernetes-istio.png)
<p style="text-align: center;">
Deployment Details Page
</p>

- Now you have an automated blue-green deployment for your application. 🎉

## Understanding what happened

In this example, you configured the application configuration file to switch all traffic from an old to a new version of the application using Istio's weighted routing feature.

- Stage 1: `K8S_CANARY_ROLLOUT` ensures that the workloads of canary variant (new version) should be deployed. But at this time, they still handle nothing, all traffic is handled by workloads of primary variant.
The number of workloads (e.g. pod) for canary variant is configured to be 100% of the replicas number of primary varant.

![](/images/example-bluegreen-kubernetes-istio-stage-1.png)

- Stage 2: `K8S_TRAFFIC_ROUTING` ensures that all traffic should be routed to canary variant. Because the `trafficRouting` is configured to use Istio, PipeCD will find Istio's VirtualService resource of this application to control the traffic percentage.
(You can add an [ANALYSIS](../../managing-application/customizing-deployment/automated-deployment-analysis/) stage after this to validate the new version. When any negative impacts are detected, an auto-rollback stage will be executed to switch all traffic back to the primary variant.)

![](/images/example-bluegreen-kubernetes-istio-stage-2.png)

- Stage 3: `WAIT_APPROVAL` waits for a manual approval from someone in your team.

- Stage 4: `K8S_PRIMARY_ROLLOUT` ensures that all resources of primary variant will be updated to the new version.

![](/images/example-bluegreen-kubernetes-istio-stage-4.png)

- Stage 5: `K8S_TRAFFIC_ROUTING` ensures that all traffic should be routed to primary variant. Now primary variant is running the new version so it means all traffic is handled by the new version.

![](/images/example-bluegreen-kubernetes-istio-stage-5.png)

- Stage 6: `K8S_CANARY_CLEAN` ensures all created resources for canary variant should be destroyed.

![](/images/example-bluegreen-kubernetes-istio-stage-6.png)
