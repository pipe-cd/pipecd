# examples

A repository contains some examples for PipeCD.

**NOTE**: This repository is automatically synced from the examples directory of [pipe-cd/pipe](https://github.com/pipe-cd/pipe/tree/master/examples) repository. If you want to make a pull request, please send it to [pipe-cd/pipe](https://github.com/pipe-cd/pipe) repository.

</br>

### Kubernetes Applications

| Name                                                                        | Description |
|-----------------------------------------------------------------------------|-------------|
| [simple](https://github.com/pipe-cd/examples/tree/master/kubernetes/simple) | Deploy plain-yaml manifests in application directory without using pipeline. |
| [helm-local-chart](https://github.com/pipe-cd/examples/tree/master/kubernetes/helm-local-chart) | Deploy a helm chart sourced from the same Git repository. |
| [helm-remote-chart](https://github.com/pipe-cd/examples/tree/master/kubernetes/helm-remote-chart) | Deploy a helm chart sourced from a [Helm Chart Repository](https://helm.sh/docs/topics/chart_repository/). |
| [helm-remote-git-chart](https://github.com/pipe-cd/examples/tree/master/kubernetes/helm-remote-git-chart) | Deploy a helm chart sourced from another Git repository. |
| [kustomize-local-base](https://github.com/pipe-cd/examples/tree/master/kubernetes/kustomize-local-base) | Deploy a kustomize package that just uses the local bases from the same Git repository. |
| [kustomize-remote-base](https://github.com/pipe-cd/examples/tree/master/kubernetes/kustomize-remote-base) | Deploy a kustomize package that uses remote bases from other Git repositories. |
| [canary](https://github.com/pipe-cd/examples/tree/master/kubernetes/canary) | Deloyment pipeline with canary strategy. |
| [canary-by-config-change](https://github.com/pipe-cd/examples/tree/master/kubernetes/canary-by-config-change) | Deployment pipeline with canary strategy when ConfigMap was changed. |
| [canary-patch](https://github.com/pipe-cd/examples/tree/master/kubernetes/canary-patch) | Demonstrate how to customize manifests for Canary variant using [patches](https://pipecd.dev/docs/user-guide/configuration-reference/#kubernetescanaryrolloutstageoptions) option. |
| [bluegreen](https://github.com/pipe-cd/examples/tree/master/kubernetes/bluegreen) | Deployment pipeline with bluegreen strategy. This also contains a manual approval stage. |
| [mesh-istio-canary](https://github.com/pipe-cd/examples/tree/master/kubernetes/mesh-istio-canary) | Deployment pipeline with canary strategy by using Istio for traffic routing.  |
| [mesh-istio-bluegreen](https://github.com/pipe-cd/examples/tree/master/kubernetes/mesh-istio-bluegreen) | Deployment pipeline with bluegreen strategy by using Istio for traffic routing. |
| [mesh-smi-canary](https://github.com/pipe-cd/examples/tree/master/kubernetes/mesh-smi-canary) | Deployment pipeline with canary strategy by using SMI for traffic routing. |
| [mesh-smi-bluegreen](https://github.com/pipe-cd/examples/tree/master/kubernetes/mesh-smi-bluegreen) | Deployment pipeline with bluegreen strategy by using SMI for traffic routing. |
| [wait-approval](https://github.com/pipe-cd/examples/tree/master/kubernetes/wait-approval) | Deployment pipeline that contains a manual approval stage. |
| [multi-steps-canary](https://github.com/pipe-cd/examples/tree/master/kubernetes/multi-steps-canary) | Deployment pipeline with multiple canary steps. |
| [analysis-by-metrics](https://github.com/pipe-cd/examples/tree/master/kubernetes/analysis-by-metrics) | Deployment pipeline with analysis stage by metrics. |
| [analysis-by-http](https://github.com/pipe-cd/examples/tree/master/kubernetes/analysis-by-http) | Deployment pipeline with analysis stage by running http requests. |
| [analysis-by-log](https://github.com/pipe-cd/examples/tree/master/kubernetes/analysis-by-log) | Deployment pipeline with analysis stage by checking logs. |
| [analysis-with-baseline](https://github.com/pipe-cd/examples/tree/master/kubernetes/analysis-with-baseline) | Deployment pipeline with analysis stage by comparing baseline and canary. |
| [secret-management](https://github.com/pipe-cd/examples/tree/master/kubernetes/secret-management) | Demonstrate how to manage sensitive data by using [Secret Management](https://pipecd.dev/docs/user-guide/secret-management/) feature. |

### Terraform Applications

| Name                                                                        | Description |
|-----------------------------------------------------------------------------|-------------|
| [simple](https://github.com/pipe-cd/examples/tree/master/terraform/simple) |  Automatically applies when any changes were detected.  |
| [local-module](https://github.com/pipe-cd/examples/tree/master/terraform/local-module) | Deploy application that using local terraform modules from the same Git repository. |
| [remote-module](https://github.com/pipe-cd/examples/tree/master/terraform/remote-module) | Deploy application that using remote terraform modules from other Git repositories. |
| [wait-approval](https://github.com/pipe-cd/examples/tree/master/terraform/wait-approval) | Deployment pipeline that contains a manual approval stage. |
| [autorollback](https://github.com/pipe-cd/examples/tree/master/terraform/auto-rollback) |  Automatically rollback the changes when deployment was failed.  |
| [secret-management](https://github.com/pipe-cd/examples/tree/master/terraform/secret-management) | Demonstrate how to manage sensitive data by using [Secret Management](https://pipecd.dev/docs/user-guide/secret-management/) feature. |

### CloudRun Applications

| Name                                                                        | Description |
|-----------------------------------------------------------------------------|-------------|
| [simple](https://github.com/pipe-cd/examples/tree/master/cloudrun/simple) | Quick sync by rolling out the new version and switching all traffic to it. |
| [canary](https://github.com/pipe-cd/examples/tree/master/cloudrun/canary) | Deployment pipeline with canary strategy. |
| [analysis](https://github.com/pipe-cd/examples/tree/master/cloudrun/analysis) | Deployment pipeline that contains an analysis stage. |
| [secret-management](https://github.com/pipe-cd/examples/tree/master/cloudrun/secret-management) | Demonstrate how to manage sensitive data by using [Secret Management](https://pipecd.dev/docs/user-guide/secret-management/) feature. |

### Lambda Applications

| Name                                                                        | Description |
|-----------------------------------------------------------------------------|-------------|
| [simple](https://github.com/pipe-cd/examples/tree/master/lambda/simple) | Quick sync by rolling out the new version and switching all traffic to it. |
| [canary](https://github.com/pipe-cd/examples/tree/master/lambda/canary) | Deployment pipeline with canary strategy. |
| [analysis](https://github.com/pipe-cd/examples/tree/master/lambda/analysis) | Deployment pipeline that contains an analysis stage. |
| [secret-management](https://github.com/pipe-cd/examples/tree/master/lambda/secret-management) | Demonstrate how to manage sensitive data by using [Secret Management](https://pipecd.dev/docs/user-guide/secret-management/) feature. |

### ECS Applications

| Name                                                                        | Description |
|-----------------------------------------------------------------------------|-------------|
| [simple](https://github.com/pipe-cd/examples/tree/master/ecs/simple) | Quick sync by rolling out the new version and switching all traffic to it. |
| [canary](https://github.com/pipe-cd/examples/tree/master/ecs/canary) | Deployment pipeline with canary strategy. |
| [bluegreen](https://github.com/pipe-cd/examples/tree/master/ecs/bluegreen) | Deployment pipeline with blue-green strategy. |
| [secret-management](https://github.com/pipe-cd/examples/tree/master/ecs/secret-management) | Demonstrate how to manage sensitive data by using [Secret Management](https://pipecd.dev/docs/user-guide/secret-management/) feature. |

**Note** that the `.kapetanios` directory is for our CI configurations. It has nothing to do with PipeCD.
