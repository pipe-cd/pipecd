- Start Date: (2021-03-02)
- Target Version: (1.x / 2.x)

# Summary

This RFC proposes adding a new service deployment from PipeCD: AWS ECS deployment.

# Motivation

PipeCD aims to support wide range of deployable services, currently [Terraform deployment](https://pipecd.dev/docs/feature-status/#terraform-deployment) and [Lambda deployment](https://pipecd.dev/docs/feature-status/#lambda-deployment) are supported. ECS deployment meets a lot of requests we received.

# Detailed design

The deployment configuration is used to customize the way to do the deployment. In the case of AWS ECS deployment, current common stages (`WAIT`, `WAIT_APPROVAL`, `ANALYSIS`) are all inherited, besides with the stages for ECS deployment `ECS_SYNC`, `ECS_CANARY_ROLLOUT` and `ECS_TRAFFIC_ROUTING`.

In case of `ECS_SYNC`, PipeCD simply applys `serviceDefinition` and `taskDefinition`.
In case you use pipeline, you need to use `External` as a deployment controller in your `serviceDefinition`.
In case of `ECS_TRAFFIC_ROUTING`, PipeCD changes the configuration of the load balancer you specified in .pipe.yaml in order to change the traffic routing state.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: ECSApp
spec:
  input:
    name: Sample
    serviceDefinition: path/to/servicedef.json # optional
    taskDefinition: path/to/taskdef.json       # required
    loadBalancerInfo:                          # optional
        containerName: sample-app
        containerPort: 80
  pipeline:
    stages:
      # Deploy workloads of the new version.
      # But this is still receiving no traffic.
      - name: ECS_CANARY_ROLLOUT
      # Change the traffic routing state where
      # the new version will receive the specified percentage of traffic.
      # This is known as multi-phase canary strategy.
      - name: ECS_TRAFFIC_ROUTING
        with:
          newVersion: 10
      # Optional: We can also add an ANALYSIS stage to verify the new version.
      # If this stage finds any not good metrics of the new version,
      # a rollback process to the previous version will be executed.
      - name: ANALYSIS
      # Change the traffic routing state where
      # thre new version will receive 100% of the traffic.
      - name: ECS_TRAFFIC_ROUTING
        with:
          newVersion: 100
```

## Architecture

Just as current Lambda but under `pkg/cloudprovider/ecs` package.

# Alternatives

The deployment configuration sample as bellow:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: ECSApp
spec:
  input:
    name: Sample
    serviceDefinition: path/to/servicedef.json # optional
    taskDefinition: path/to/taskdef.json       # required
    loadBalancerInfo:                          # optional
        containerName: sample-app
        containerPort: 80
```

# Unresolved questions

Service auto scaling is not supported when using an external deployment controller.
