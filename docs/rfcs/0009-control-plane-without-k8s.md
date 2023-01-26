- Start Date: 2022-01-18
- Target Version: 0.41.3

# Summary

Support install PipeCD control plane on other platform which is not k8s

# Motivation

Now we can deploy the control plane to kubernetes cluster, but some developers that would like to introduce PipeCD can not prepare kubernetes environments. We want to support  install PipeCD control plane on platforms other than kubernetes.

# Detailed design

1. Control Plane on docker-compose
    - Developers can deploy control plane on a single machine.
    - Developers do not have to prepare datastore and filestore by themselves or they can easily use database or filestore on local machine.
    
2. Control Plane on managed container services (ex. ECS)
    - Pipecd can give the way to deploy control plane as Terraform template.
    - They can easily use managed database or storage system on cloud as datastore and filestore.
    ![image](assets/control-plane-on-aws.jpg)
    Note:
    - Devide pipecd-server and pipecd-ops to different services because they have the same port and have different authorization.
    - Pay attention to brocking public access to s3.
        - Add IAM role to ECS to access S3.

# Alternatives

1. Control Plane without container image
    - This alternative supports deploy control plane as a binary such as Piped.
    - It is stressful for developers to set up networking by themselves.

2. Abolish envoy and use ALB instead of envoy for access distribution
    - We must move configuration for envoy to ALB.
    - We must make the same number of target groups as pipecd-server has.
# Unresolved questions

None