- Start Date: 2023-11-16
- Target Version: 0.46.0

# Summary

To support prgoressive delivery for ECS services accessed via ECS Service Discovery.

# Motivation

<!-- Why are we doing this? What do we expect? -->

- Currently, PipeCD requires ELB & Target Groups for ECS.
- However, some ECS services are not accessed via ELB, instead directly from other ECS services using ECS Service Discovery.
  - e.g. in a service mesh or gRPC backend services(gRPC is not supported by ALB)
- We'll provide progressive delivery for them too.

# Detailed design

<!-- Explain the design in detail including what components should be added or changed, examples of how the feature is used. -->

## ECS access types

There are 4 types of ECS deployment targets. We focus on (3) here.

| No. | type                                 | supported by PipeCD        | use case example                                       |
| --- | ------------------------------------ | -------------------------- | ----------------------------------------------- |
| 1   | a standalone task                    | Yes (only QuickSync)       | jobs                                            |
| 2   | a service under ELB                  | Yes (called `application`) | frontend services                               |
| 3   | a service with ECS Service Discovery | Not yet                    | internal services in a simple service mesh      |
| 4   | a service in App Mesh                | Not yet                    | internal services in a complicated service mesh |

- PipeCD needs to handle them in different ways because they have different ways of access and deployments.
- We focus on (3) in this CFP and Issue #4616 because there are some users facing such cases now.
  - However, we consider extensibility for other types like (4)App Mesh and [ECS Service Connect](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/service-connect.html).

## Pipeline Stages

- We use existing stages of ECS as below.
  - `ECS_SYNC` (for QuickSync)
  - `ECS_CANARY_ROLLOUT`
  - `ECS_PRIMARY_ROLLOUT`
  - `ECS_TRAFFIC_ROUTING`
  - `ECS_CANARY_CLEAN`
- The reason for that is simplicity.  If we add new stages like `ECS_CANARY_ROLLOUT_SERVICE_DISCOVERY`, there'll be too many stages as deployment targets increase, which will make users & PipeCD developers confuse which stages to use.
- The stages will behave slightly different from the current as follows.

## Deployment Flows

- We'll support 3 types of deployments as current.
  - (A) QuickSync
  - (B) Canary
  - (C) Blue/Green
- NOTE: In Canary and Blue/Green, tasks will start to receive traffic while rollout right after deployments, unlike current ECS deployments.
  - That's because ECS Service Discovery automatically register new tasks to the namespace right after deployments.
      <!-- - We would not deregister new tasks from Service Discovery right after the registration, because that would lead to bugs. -->
  - Therefore, in `ECS_CANARY_ROLLOUT` `ECS_PRIMARY_ROLLOUT`, we'll deregister the tasks from Service Discovery in order to stop receiving traffic.
  - alternatives and why not adopted:
    - *not deregister the tasks and start to receive traffic in rollout stages*
      - That would prevent flexible pipelines because rollout stages will have multiple responsibilities.
    - *turn off the Service Discovery option while rollout and turn it on after the deployment*
      - That will reboot tasks.

### (A) QuickSync Flow

| stage      | what's executed                         |
| ---------- | --------------------------------------- |
| `ECS_SYNC` | simply update the primary service tasks |

- the same as the current process.

### (B) Canary Flow

| stage                 | what's executed                                                                                       |
| --------------------- | ----------------------------------------------------------------------------------------------------- |
| `ECS_CANARY_ROLLOUT` (`scale:xx`)  | create the secondary service (if not exist^1) & tasks, and deregister them from Service Discovery^2 |
| `ECS_TRAFFIC_ROUTING`  (`canary:yy`) | register the secondary to Service Discovery^3 (at least 1 task)                                                        |
| `ECS_PRIMARY_ROLLOUT` | update the primary service tasks, and automatically register them to Service Discovery^4                |
| `ECS_TRAFFIC_ROUTING` (`primary:100`) | deregister the secondary from Service Discovery                                                       |
| `ECS_CANARY_CLEAN`    | delete the secondary service & tasks                                                                  |

- ^1: If there are multi `ECS_CANARY_ROLLOUT` stages, the service will not be recreated after the first `ECS_CANARY_ROLLOUT` stage.
- ^2: The target namespace service is the same as the primary service.
- ^3: We can't route traffic to the secondary strictly based on the `canary` value.
  - We only adjust n of primary/canary tasks under Service Discovery.
  - That's because [ECS Service Discovery does not support weighted routing](https://docs.aws.amazon.com/cloud-map/latest/dg/services-values.html#service-creating-values-routing-policy).
- ^4: We need to keep the primary/canary ratio by registering/deregistering.

### (C) Blue/Green Flow

| stage                                    | what's executed                                                                                                                                                              |
| ---------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `ECS_CANARY_ROLLOUT`<br>(`scale:100`)    | create the secondary service (if not exist) & tasks, and deregister them from Service Discovery |
| `ECS_TRAFFIC_ROUTING`<br>(`canary:100`)  | (1) register the secondary to Service Discovery <br> (2) deregister the primary from Service Discovery                                                                        |
| `ECS_PRIMARY_ROLLOUT`                    | update the primary service tasks, and deregister them from Service Discovery ^1                                                                                                |
| `ECS_TRAFFIC_ROUTING`<br>(`primary:100`) | (1) register the primary to Service Discovery <br> (2) deregister the secondary from Service Discovery                                                                       |
| `ECS_CANARY_CLEAN`                       | delete the secondary service & tasks                                                                                                                                         |

- ^1 : to prevent the primary from receiving traffic before `ECS_TRAFFIC_ROUTING` is started.

## Config

We add one config as below:

| key                        | description                                 | values                                                          | default |
| -------------------------- | ------------------------------------------- | --------------------------------------------------------------- | ------- |
| `spec:input:ecsAccessType` | to determine which 'ECS access type' to use | `ELB`, `ECS_SERVICE_DISCOVERY` (`APP_MESH` in the future) ^1 | `ELB`^2   |

- ^1: We don't need `STANDALONE` option for `ecsAccessType` because it's determined by whether `spec:input:serviceDefinitionFile` is configured.
- ^2: In order to prevent users who use ELB now from beging affected.
- Users don't need to configure `spec:input:targetGroups` when not selecting `ELB` for `ecsAccessType`.
- Users don't need to configure `ecsAccessType` when not selecting `ECSApp` as `kind`.

# Unresolved questions

- Right now we would not support cases that a ECS service is accessed from both ELB and other ECS services . They are not common.
  - <https://docs.aws.amazon.com/AmazonECS/latest/developerguide/service-discovery.html>
    > You can configure service discovery for a service that's behind a load balancer, but service discovery traffic is always routed to the task and not the load balancer.
  
# Further info

- App Mesh would be supported by PipeCD in the similar way. ref:
  [Create a pipeline with canary deployments for Amazon ECS using AWS App Mesh](https://aws.amazon.com/jp/blogs/containers/create-a-pipeline-with-canary-deployments-for-amazon-ecs-using-aws-app-mesh/)
