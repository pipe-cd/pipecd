- Start Date: (2021-03-02)
- Target Version: 1.1.0

# Summary

This RFC proposes adding a new service deployment from PipeCD: AWS ECS deployment.

# Motivation

PipeCD aims to support a wide range of deployable services, currently, [Terraform deployment](https://pipecd.dev/docs/feature-status/#terraform-deployment) and [Lambda deployment](https://pipecd.dev/docs/feature-status/#lambda-deployment) are supported. ECS deployment meets a lot of requests we received.

# Detailed design

The deployment configuration is used to customize the way to do the deployment. In the case of AWS ECS deployment, current common stages (`WAIT`, `WAIT_APPROVAL`, `ANALYSIS`) are all inherited, besides with the stages for ECS deployment `ECS_SYNC`, `ECS_CANARY_ROLLOUT` and `ECS_TRAFFIC_ROUTING`.

`ECS_SYNC` simply applies a service definition and a task definition. It uses API such as [`CreateService`](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_CreateService.html), [`UpdateService`](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_UpdateService.html), [`RegisterTaskDefinition`](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_RegisterTaskDefinition.html) and [`CreateTaskSet`](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_CreateTaskSet.html).

`ECS_CANARY_ROLLOUT` deploys workloads of the new version. it uses API such as [`DescribeServices`](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_DescribeServices.html), [`CreateTaskSet`](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_CreateTaskSet.html), [`UpdateServicePrimaryTaskSet`](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_UpdateServicePrimaryTaskSet.html) and [`DeleteTaskSet`](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_DeleteTaskSet.html).

`ECS_TRAFFIC_ROUTING` changes the configuration of the load balancer you specified in .pipe.yaml in order to change the traffic routing state. It uses API such as [`DescribeTargetGroups`](https://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_DescribeTargetGroups.html), [`DescribeListeners`](https://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_DescribeListeners.html), [`ModifyListener`](https://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_ModifyListener.html), [`DescribeRules`](https://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_DescribeRules.html) and [`ModifyRule`](https://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_ModifyRule.html).

In case of using pipeline, you need to use `External` as a deployment controller in your `serviceDefinition`.

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

## Example

### Canary Deployment

```yaml
apiVersion: pipecd.dev/v1beta1
kind: ECSApp
spec:
  input:
    name: Sample
    serviceDefinition: path/to/servicedef.json
    taskDefinition: path/to/taskdef.json
    loadBalancerInfo:
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
          canary: 10
      # Optional: We can also add an ANALYSIS stage to verify the new version.
      # If this stage finds any not good metrics of the new version,
      # a rollback process to the previous version will be executed.
      - name: ANALYSIS
      # Change the traffic routing state where
      # thre new version will receive 100% of the traffic.
      - name: ECS_TRAFFIC_ROUTING
        with:
          canary: 100
```

### Blue-Green Deployment

```yaml
apiVersion: pipecd.dev/v1beta1
kind: ECSApp
spec:
  input:
    name: Sample
    serviceDefinition: path/to/servicedef.json
    taskDefinition: path/to/taskdef.json
    loadBalancerInfo:
        containerName: sample-app
        containerPort: 80
  pipeline:
    stages:
      # Deploy workloads of the new version.
      # But this is still receiving no traffic.
      - name: ECS_CANARY_ROLLOUT
      # Change the traffic routing state where
      # the new version will receive all traffic.
      # This is known as multi-phase canary strategy.
      - name: ECS_TRAFFIC_ROUTING
        with:
          all: canary
      - name: ANALYSIS
```

## Architecture

Just as current Lambda but under `pkg/cloudprovider/ecs` package.

# Unresolved questions

Service auto scaling is not supported when using an external deployment controller.
