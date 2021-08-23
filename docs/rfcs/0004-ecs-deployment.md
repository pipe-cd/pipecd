- Start Date: 2021-05-19
- Target Version: 1.0.1

# Summary

This RFC proposes adding a new service deployment from PipeCD: AWS ECS deployment.

# Motivation

PipeCD aims to support a wide range of deployable services, currently, [Terraform deployment](https://pipecd.dev/docs/feature-status/#terraform-deployment) and [Lambda deployment](https://pipecd.dev/docs/feature-status/#lambda-deployment) are supported. ECS deployment meets a lot of requests we received.

# Detailed design

Before start, suppose we have a map of ECS's keywords which corresponds to the related keywords in other platform
| ECS  | Common use |
|---|---|
| Container | Generic Container |
| Task Definition | Define the way how ECS should handle containers |
| Task | The instantiation of a task definition within a ECS cluster (related to Pod in k8s) |
| Service | Contains information that ECS to run and maintain a specified number of instances of a task definition simultaneously in an ECS cluster (related to Deployment in k8s) |
| Task Set | Contains information that allows an ECS service to run multiple versions of an application via variations in the task definition associated with it (related to ReplicaSet in k8s) |

In the simplest case, to enable deploy containers on to an AWS ECS cluster, we need:
- Container's image
- Task Definition which has to be registered with the AWS ECS service (it will be versioned)
- Service which will be "applied" to AWS ECS cluster, the service has versioned taskDefinitionArn information so that ECS cluster could deploy as defined

Note:
- No TaskSet is required because the "deployments" are handled by ECS itself.
- To enable customize the deployment process, we have to take a look at the [ECS Deployment Controller](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_DeploymentController.html). There are 3 types of controller are supported currently: `ECS`, `CODE_DEPLOY` and `EXTERNAL`. In case of PipeCD, we requires users to defined there ECS Applications' deployments configuration to use `EXTERNAL`.

With `EXTERNAL` deployment controller is used and required by PipeCD workflow, the above simplest case will be changed to:
- Container's image
- Task Definition which registered with AWS ECS service
- Service which has to define 3 most important spec: `DesiredCount`, `SchedulingStrategy` and `DeploymentController` which always be set to `EXTERNAL`. The `service.TaskDefinitionArn` should be set to null and it's unable to changed in case deployment controller of type `EXTERNAL` is used.
- TaskSet required to create Task with specified TaskDefinition, those Task will be controlled by Service after created (as it's replica set function). 

Insprised by the [AWS Blue/Green deployments with ECS External Deployment Controller](https://aws.amazon.com/blogs/containers/blue-green-deployments-with-the-ecs-external-deployment-controller/), in case of PipeCD's AWS ECS deployments, all current common stages (`WAIT`, `WAIT_APPROVAL`, `ANALYSIS`) are inherited, besides with stages for ECS deployment such as `ECS_SYNC`, `ECS_CANARY_ROLLOUT`, `ECS_TRAFFIC_ROUTING` and `ECS_CLEAN`.

- `ECS_SYNC` simply applies a service definition and a task definition.
- `ECS_CANARY_ROLLOUT` deploys workloads of the new version.
- `ECS_TRAFFIC_ROUTING` changes the configuration of the load balancer (you have to prepare and configure it in `.pipe.yaml` to enable using this stage) in order to change the traffic routing state.
- `ECS_CLEAN` remove old version workloads.

A simplest case (`ECS_SYNC` strategy), `.pipe.yaml` for PipeCD ECS Application would look like:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: ECSApp
spec:
  input:
    serviceDefinition: path/to/servicedef.yaml
    taskDefinition: path/to/taskdef.yaml
    targetGroups:
      primary:
        targetGroupArn: {PRIMARY_TARGET_GROUP_ARN}
        containerName: service
        containerPort: 80
```

In case of canary release strategy

```yaml
apiVersion: pipecd.dev/v1beta1
kind: ECSApp
spec:
  input:
    serviceDefinition: path/to/servicedef.json
    taskDefinition: path/to/taskdef.json
    targetGroups:
      primary:
        targetGroupArn: {PRIMARY_TARGET_GROUP_ARN}
        containerName: service
        containerPort: 80
      canary:
        targetGroupArn: {CANARY_TARGET_GROUP_ARN}
        containerName: service
        containerPort: 80
  pipeline:
    stages:
      # Deploy the workloads of CANARY variant, the number of workload
      # for CANARY variant is equal to 10% of PRIMARY's workload.
      # But this is still receiving no traffic.
      - name: ECS_CANARY_ROLLOUT
        with:
          scale: 10%
      # Change the traffic routing state where
      # the CANARY workloads will receive the specified percentage of traffic.
      # This is known as multi-phase canary strategy.
      - name: ECS_TRAFFIC_ROUTING
        with:
          canary: 20
      # Optional: We can also add an ANALYSIS stage to verify the new version.
      # If this stage finds any not good metrics of the new version,
      # a rollback process to the previous version will be executed.
      - name: ANALYSIS
      # Update the workload of PRIMARY variant to the new version.
      - name: ECS_PRIMARY_ROLLOUT
      # Change the traffic routing state where
      # the PRIMARY workloads will receive 100% of the traffic.
      - name: ECS_TRAFFIC_ROUTING
        with:
          primary: 100
      # Destroy all workloads of CANARY variant.
      - name: ECS_CANARY_CLEAN
```

In case of blue/green release strategy

```yaml
apiVersion: pipecd.dev/v1beta1
kind: ECSApp
spec:
  input:
    serviceDefinition: path/to/servicedef.json
    taskDefinition: path/to/taskdef.json
    targetGroups:
      primary:
        targetGroupArn: {PRIMARY_TARGET_GROUP_ARN}
        containerName: service
        containerPort: 80
      canary:
        targetGroupArn: {CANARY_TARGET_GROUP_ARN}
        containerName: service
        containerPort: 80
  pipeline:
    stages:
      # Deploy the workloads of CANARY variant with the number of
      # workload for CANARY variant is the same as PRIMARY variant.
      # But this is still receiving no traffic.
      - name: ECS_CANARY_ROLLOUT
        with:
          scale: 100%
      # Change the traffic routing state where
      # the CANARY workloads will receive 100% of the traffic.
      - name: ECS_TRAFFIC_ROUTING
        with:
          canary: 100
      # Optional: We can also add an ANALYSIS stage to verify the new version.
      # If this stage finds any not good metrics of the new version,
      # a rollback process to the previous version will be executed.
      - name: ANALYSIS
      # Update the workload of PRIMARY variant to the new version.
      - name: ECS_PRIMARY_ROLLOUT
      # Change the traffic routing state where
      # the PRIMARY workloads will receive 100% of the traffic.
      - name: ECS_TRAFFIC_ROUTING
        with:
          primary: 100
      # Destroy all workloads of CANARY variant.
      - name: ECS_CANARY_CLEAN
```

# Unresolved questions

Service auto scaling is not supported when using an external deployment controller.
