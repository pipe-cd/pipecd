---
date: 2022-02-07
title: "PipeCD best practice 02 - control plane on ECS"
linkTitle: "PipeCD best practice 02"
weight: 996
description: "This blog is a guideline for you to operate your own PipeCD on Amazon ECS."
author: Yohei Namba ([@kevin_namba](https://twitter.com/kevin_namba))
---

This blog is a part of PipeCD best practice series, a guideline for you to operate your own PipeCD.
Currently, you can deploy and operate the PipeCD control plane on a Kubernetes cluster easily, but some developers that would like to introduce PipeCD can not prepare Kubernetes environments. If you have the same problem, this blog is for you. We will show you how to deploy the PipeCD control plane on Amazon ECS.

### Architecture
> Note: Please refer to [architecture-overview](docs/user-guide/managing-controlplane/architecture-overview/) docs for definitions of PipeCD components such as server, ops, cache, datastore and filestore.

![](/images/control-plane-on-ecs.png)

Following the above graph for PipeCD control plane runs on Amazon ECS, we have to prepare these next components

### Secrets Manager
You should put config files (control-plane-config.yaml, envoy-config.yaml) on Secrets Manager because config files contain some credentials such as database passwords.
The examples of config files for ECS are [here](https://github.com/pipe-cd/control-plane-aws-ecs-terraform-demo/tree/main/config). Please edit these files according to your environment.
```
aws secretsmanager create-secret --name control-plane-config \
--description "Configuration of control plane" \
--secret-string `base64 control-plane-config.yaml`
aws secretsmanager create-secret --name envoy-config \
--description "Configuration of control plane" \
--secret-string `base64 envoy-config.yaml`
```
You should also put encryption key
```
aws secretsmanager create-secret --name encryption-key \
--description "Encryption key for control plane" \
--secret-string `openssl rand 64 | base64`
```

### RDS(datastore)
It is possible to use RDS as a datastore. Edit your configuration file for the control plane according to your RDS setting.
```yaml
  datastore:
    type: MYSQL
    config: 
        url: root:password@tcp(endpoint_of_rds:3306) 
        database: quickstart
```

### Redis(cache)
It is possible to use Redis as a cache. Note the endpoint of Redis for the task definition.


### S3(filestore)
It is possible to use S3 as a filestore. The filestore contains state files that describe your secure infrastructure, so make sure to make the bucket private. Only allow pipecd-server(ECS) to access this bucket.

### ECS
You need to create two different services for pipecd-server and a pipecd-ops because they have the same ports and different permissions. The pipecd-server can be accessed by external clients such as piped or web clients, so this service includes the pipe-cd gateway and this service must be connected to the application loadbalancer. The pipecd-ops can only be accessed by admin users, so this service must only be accessed via SSM session manager.
ECS agent sets config files as environment variables in container from secrets manager. Create configuration files from environment variables in the container as below.

```bash
echo $ENVOY_CONFIG; echo $ENVOY_CONFIG | base64 -d >> envoy-config.yaml
```

> Note: Attach IAM policy to get secrets from Secrets Manager to task execution role.

> Note: If you manage both of RDS and ECS by terraform, you can rewite datastore endpoint in the configuration file.
`sed -i -e s/pipecd-mysql/${var.db_instance_address}/ control-plane-config.yaml;`

> Note: Attach IAM policy to access filestore to task role.


#### task definitions examples (using terraform variables)
1. gateway and server
```
{
name  = "pipecd-gateway"
image = var.gateway_image_url
portMappings = [
  {
    hostPort      = 9090
    containerPort = 9090
    protocol      = "tcp"
  }
]
essential = false
command = [
  "/bin/sh -c 'echo $ENVOY_CONFIG; echo $ENVOY_CONFIG | base64 -d >> envoy-config.yaml; envoy -c envoy-config.yaml;'"
]
entrypoint = [
  "sh",
  "-c"
]
secrets = [
  {
    "name" : "ENVOY_CONFIG",
    "valueFrom" : "arn:aws:secretsmanager:${data.aws_region.current.id}:${data.aws_caller_identity.self.account_id}:secret:${var.envoy_config_secret}"
  },
]
},
{
name  = "pipecd-server"
image = var.server_image_url
portMappings = [
]
command = [
  "/bin/sh -c 'echo $CONTROL_PLANE_CONFIG; echo $CONTROL_PLANE_CONFIG | base64 -d >> control-plane-config.yaml; sed -i -e s/pipecd-mysql/${var.db_instance_address}/ control-plane-config.yaml; echo $ENCRYPTION_KEY >> encryption-key; pipecd server --insecure-cookie=true --cache-address=${var.redis_host}:6379 --config-file=control-plane-config.yaml --enable-grpc-reflection=false --encryption-key-file=encryption-key --log-encoding=humanize --metrics=true;'"
]
entrypoint = [
  "sh",
  "-c"
]
secrets = [
  {
    "name" : "ENCRYPTION_KEY",
    "valueFrom" : "arn:aws:secretsmanager:${data.aws_region.current.id}:${data.aws_caller_identity.self.account_id}:secret:${var.encryption_key_secret}"
  },
  {
    "name" : "CONTROL_PLANE_CONFIG",
    "valueFrom" : "arn:aws:secretsmanager:${data.aws_region.current.id}:${data.aws_caller_identity.self.account_id}:secret:${var.control_plane_config_secret}"
  },
]
},
```

2. ops
```
{
name  = "pipecd-ops"
image = var.ops_image_url
portMappings = [
  {
    "name" : "http",
    "protocol" : "tcp",
    "containerPort" : 9082,
    "appProtocol" : "http"
  },
  {
    "name" : "http",
    "protocol" : "tcp",
    "containerPort" : 9085,
    "appProtocol" : "http"
  },
]
command = [
  "/bin/sh -c 'echo $CONTROL_PLANE_CONFIG; echo $CONTROL_PLANE_CONFIG | base64 -d >> control-plane-config.yaml; sed -i -e s/pipecd-mysql/${var.db_instance_address}/ control-plane-config.yaml; echo $ENCRYPTION_KEY >> encryption-key; pipecd ops --cache-address=${var.redis_host}:6379 --config-file=control-plane-config.yaml --log-encoding=humanize --metrics=true;'"
]
entrypoint = [
  "sh",
  "-c"
]
secrets = [
  {
    "name" : "ENCRYPTION_KEY",
    "valueFrom" : "arn:aws:secretsmanager:${data.aws_region.current.id}:${data.aws_caller_identity.self.account_id}:secret:${var.encryption_key_secret}"
  },
  {
    "name" : "CONTROL_PLANE_CONFIG",
    "valueFrom" : "arn:aws:secretsmanager:${data.aws_region.current.id}:${data.aws_caller_identity.self.account_id}:secret:${var.control_plane_config_secret}"
  },
]
},
```

### ALB
You must prepare two target groups for both HTTP and gRPC. Make two hosts and listner rules as below. Listner protocol should be HTTPS becuase it uses gRPC.

![](/images/control-plane-alb.png)

### Terraform example
PipeCD gives [control-plane-aws-ecs-terraform-demo](https://github.com/pipe-cd/control-plane-aws-ecs-terraform-demo), which we use Terraform to prepare PipeCD controlplane components and install it on air.

#### Prepare
1. Prepare SSL certificate
Prepare your domain and SSL certificate from AWS certificate manager.

2. Create a s3 bucket for terraform backend
Write bucket name to `00-main.tf`
```
terraform {
  backend "s3" {
    bucket  = "example-pipecd-control-plane-tfstate" #your bucket name for terraform backend
    region  = "ap-northeast-1"
    key     = "tfstate"
    profile = "pipecd-control-planeg-terraform" #your profile
  }
  required_providers {
    aws = {
      version = "~> 3.34.0"
    }
  }
}
```

3. Edit `variables.tf` for your project
```
//export
locals {
  alb = {
    certificate_arn = ""
  }
  
  redis = {
    node_type = "cache.t2.micro"
  }

  rds = {
    node_type = "db.t3.micro"
  }

  ecs = {
    memory = "1024"
    cpu = "512"
  }
}
```

4. Create a S3 bucket for filestore and write the bucket name for it to `control-plane-config.yaml` and `variables.tf`
```
apiVersion: "pipecd.dev/v1beta1"
kind: ControlPlane
spec:
  datastore:
    type: MYSQL
    config:
      url: sample:test@tcp(pipecd-mysql:3306)
      database: quickstart
  filestore:
    type: S3
    config: # edit here
        bucket: example-pipecd-control-plane-filestore 
        region: ap-northeast-1
  projects:
    - id: quickstart
        staticAdmin:
          username: hello-pipecd
          passwordHash: "$2a$10$ye96mUqUqTnjUqgwQJbJzel/LJibRhUnmzyypACkvrTSnQpVFZ7qK" # bcrypt value of "hello-pipecd"
```
```
//export
locals {
  s3 = { # These must be unique in the world.
    filestore_bucket = "${local.project}-filestore" # edit here
  }
}
```
5. Write config of RDS for datastore to `control-plane-config.yaml`
Note: Do not edit hostname (pipecd-mysql) because it will be edited autimaticaly by terraform.
```
apiVersion: "pipecd.dev/v1beta1"
kind: ControlPlane
spec:
  datastore:
    type: MYSQL
    config: # edit here
      url: sample:test@tcp(pipecd-mysql:3306)
      database: quickstart
  filestore:
    type: S3
    config: 
      bucket: example-pipecd-control-plane-filestore 
      region: ap-northeast-1
  projects:
    - id: quickstart
      staticAdmin:
        username: hello-pipecd
        passwordHash: "$2a$10$ye96mUqUqTnjUqgwQJbJzel/LJibRhUnmzyypACkvrTSnQpVFZ7qK" # bcrypt value of "hello-pipecd"
```

6. Put an encryption key and config file in Secrets Manager and write the path to `variables.tf`
```
locals {
  sm = {
    control_plane_config_secret = ""
    envoy_config_secret         = ""
    encryption_key_secret       = ""
  }
}
```

#### Deploy
```
terraform apply
```

#### login admin console
You can login pipecd-ops via ECS exec.
```
aws ssm start-session --target ecs:${CLUSTER}_${TASK_ID}_${CONTAINER_ID} --document-name AWS-StartPortForwardingSession --parameters '{"portNumber":["9082"],"localPortNumber":["19082"]}'
```

That's all! Now you have your own PipeCD controlplane and it's ready to go!!
