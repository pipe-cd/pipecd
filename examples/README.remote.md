# Examples

A repository contains some examples for PipeCD.

**NOTE**: This repository is automatically synced from the examples directory of [pipe-cd/pipecd](https://github.com/pipe-cd/pipecd/tree/master/examples) repository. If you want to make a pull request, please send it to [pipe-cd/pipecd](https://github.com/pipe-cd/pipecd) repository.

</br>

### Kubernetes Applications

| Running on Play | Name                                                                        | Description |
|-----------------|-----------------------------------------------------------------------------|-------------|
| [link](https://play.pipecd.dev/applications/558401f0-8a35-494a-a9ba-dd0afe79824e?project=play) | [simple](https://github.com/pipe-cd/examples/tree/master/kubernetes/simple) | Deploy plain-yaml manifests in application directory without using pipeline. |
| -- | [helm-local-chart](https://github.com/pipe-cd/examples/tree/master/kubernetes/helm-local-chart) | Deploy a helm chart sourced from the same Git repository. |
| [link](https://play.pipecd.dev/applications/36347720-8f03-417d-8465-094f7d4eb4b1?project=play) | [helm-remote-chart](https://github.com/pipe-cd/examples/tree/master/kubernetes/helm-remote-chart) | Deploy a helm chart sourced from a [Helm Chart Repository](https://helm.sh/docs/topics/chart_repository/). |
| [link](https://play.pipecd.dev/applications/f7fc49cf-71e1-4932-8ba4-8863eeace077?project=play) | [helm-remote-git-chart](https://github.com/pipe-cd/examples/tree/master/kubernetes/helm-remote-git-chart) | Deploy a helm chart sourced from another Git repository. |
| [link](https://play.pipecd.dev/applications/a01c3ebb-89d2-4569-bef7-d659412daa11?project=play) | [kustomize-local-base](https://github.com/pipe-cd/examples/tree/master/kubernetes/kustomize-local-base) | Deploy a kustomize package that just uses the local bases from the same Git repository. |
| -- | [kustomize-remote-base](https://github.com/pipe-cd/examples/tree/master/kubernetes/kustomize-remote-base) | Deploy a kustomize package that uses remote bases from other Git repositories. |
| [link](https://play.pipecd.dev/applications/374119cd-f3a8-47f2-93db-99f58855e5a4?project=play) | [canary](https://github.com/pipe-cd/examples/tree/master/kubernetes/canary) | Deloyment pipeline with canary strategy. |
| -- | [canary-by-config-change](https://github.com/pipe-cd/examples/tree/master/kubernetes/canary-by-config-change) | Deployment pipeline with canary strategy when ConfigMap was changed. |
| -- | [canary-patch](https://github.com/pipe-cd/examples/tree/master/kubernetes/canary-patch) | Demonstrate how to customize manifests for Canary variant using [patches](https://pipecd.dev/docs/user-guide/configuration-reference/#kubernetescanaryrolloutstageoptions) option. |
| [link](https://play.pipecd.dev/applications/b8575010-9619-4141-bb0e-6d58ee5d09c9?project=play) | [bluegreen](https://github.com/pipe-cd/examples/tree/master/kubernetes/bluegreen) | Deployment pipeline with bluegreen strategy. This also contains a manual approval stage. |
| -- | [mesh-istio-canary](https://github.com/pipe-cd/examples/tree/master/kubernetes/mesh-istio-canary) | Deployment pipeline with canary strategy by using Istio for traffic routing.  |
| -- | [mesh-istio-bluegreen](https://github.com/pipe-cd/examples/tree/master/kubernetes/mesh-istio-bluegreen) | Deployment pipeline with bluegreen strategy by using Istio for traffic routing. |
| -- | [mesh-smi-canary](https://github.com/pipe-cd/examples/tree/master/kubernetes/mesh-smi-canary) | Deployment pipeline with canary strategy by using SMI for traffic routing. |
| -- | [mesh-smi-bluegreen](https://github.com/pipe-cd/examples/tree/master/kubernetes/mesh-smi-bluegreen) | Deployment pipeline with bluegreen strategy by using SMI for traffic routing. |
| [link](https://play.pipecd.dev/applications/72dbd53e-a90a-41b3-8503-44af2edeb507?project=play) | [wait-approval](https://github.com/pipe-cd/examples/tree/master/kubernetes/wait-approval) | Deployment pipeline that contains a manual approval stage. |
| -- | [multi-steps-canary](https://github.com/pipe-cd/examples/tree/master/kubernetes/multi-steps-canary) | Deployment pipeline with multiple canary steps. |
| [link](https://play.pipecd.dev/applications/913a0bde-1f38-41e3-9f56-75910b8988a9?project=play) | [analysis-by-metrics](https://github.com/pipe-cd/examples/tree/master/kubernetes/analysis-by-metrics) | Deployment pipeline with analysis stage by metrics. |
| -- | [analysis-by-http](https://github.com/pipe-cd/examples/tree/master/kubernetes/analysis-by-http) | Deployment pipeline with analysis stage by running http requests. |
| -- | [analysis-by-log](https://github.com/pipe-cd/examples/tree/master/kubernetes/analysis-by-log) | Deployment pipeline with analysis stage by checking logs. |
| -- | [analysis-with-baseline](https://github.com/pipe-cd/examples/tree/master/kubernetes/analysis-with-baseline) | Deployment pipeline with analysis stage by comparing baseline and canary. |
| -- | [secret-management](https://github.com/pipe-cd/examples/tree/master/kubernetes/secret-management) | Demonstrate how to manage sensitive data by using [Secret Management](https://pipecd.dev/docs/user-guide/secret-management/) feature. |

### Terraform Applications

| Running on Play | Name                                                                        | Description |
|-----------------|-----------------------------------------------------------------------------|-------------|
| [link](https://play.pipecd.dev/applications/ece10473-0cdb-4fec-96a1-a3df8f2e3c6e?project=play) | [simple](https://github.com/pipe-cd/examples/tree/master/terraform/simple) |  Automatically applies when any changes were detected.  |
| -- | [local-module](https://github.com/pipe-cd/examples/tree/master/terraform/local-module) | Deploy application that using local terraform modules from the same Git repository. |
| -- | [remote-module](https://github.com/pipe-cd/examples/tree/master/terraform/remote-module) | Deploy application that using remote terraform modules from other Git repositories. |
| [link](https://play.pipecd.dev/applications/4726503e-68e0-40a0-b9cb-9761567f4745?project=play) | [wait-approval](https://github.com/pipe-cd/examples/tree/master/terraform/wait-approval) | Deployment pipeline that contains a manual approval stage. |
| -- | [autorollback](https://github.com/pipe-cd/examples/tree/master/terraform/auto-rollback) |  Automatically rollback the changes when deployment was failed.  |
| [link](https://play.pipecd.dev/applications/33b9b73b-acf2-4cd4-9e0c-ab2e9fad86d1?project=play) | [secret-management](https://github.com/pipe-cd/examples/tree/master/terraform/secret-management) | Demonstrate how to manage sensitive data by using [Secret Management](https://pipecd.dev/docs/user-guide/secret-management/) feature. |

### Cloud Run Applications

| Running on Play | Name                                                                        | Description |
|-----------------|-----------------------------------------------------------------------------|-------------|
| [link](https://play.pipecd.dev/applications/64eee87f-7fae-4760-81cc-c6e66f1b48c9?project=play) | [simple](https://github.com/pipe-cd/examples/tree/master/cloudrun/simple) | Quick sync by rolling out the new version and switching all traffic to it. |
| [link](https://play.pipecd.dev/applications/845613b4-f997-4682-9529-98f089480394?project=play) | [canary](https://github.com/pipe-cd/examples/tree/master/cloudrun/canary) | Deployment pipeline with canary strategy. |
| [link](https://play.pipecd.dev/applications/c1fcbca1-c3ed-41f6-b8d9-0a1ee28df5c3?project=play) | [wait-approval](https://github.com/pipe-cd/examples/tree/master/cloudrun/wait-approval) | Deployment pipeline that contains a manual approval stage. |
| -- | [analysis](https://github.com/pipe-cd/examples/tree/master/cloudrun/analysis) | Deployment pipeline that contains an analysis stage. |
| -- | [secret-management](https://github.com/pipe-cd/examples/tree/master/cloudrun/secret-management) | Demonstrate how to manage sensitive data by using [Secret Management](https://pipecd.dev/docs/user-guide/secret-management/) feature. |

### Lambda Applications

| Running on Play | Name                                                                        | Description |
|-----------------|-----------------------------------------------------------------------------|-------------|
| -- | [simple](https://github.com/pipe-cd/examples/tree/master/lambda/simple) | Quick sync by rolling out the new version and switching all traffic to it. |
| -- | [canary](https://github.com/pipe-cd/examples/tree/master/lambda/canary) | Deployment pipeline with canary strategy. |
| -- | [analysis](https://github.com/pipe-cd/examples/tree/master/lambda/analysis) | Deployment pipeline that contains an analysis stage. |
| -- | [secret-management](https://github.com/pipe-cd/examples/tree/master/lambda/secret-management) | Demonstrate how to manage sensitive data by using [Secret Management](https://pipecd.dev/docs/user-guide/secret-management/) feature. |

### ECS Applications

| Running on Play | Name                                                                        | Description |
|-----------------|-----------------------------------------------------------------------------|-------------|
| -- | [simple](https://github.com/pipe-cd/examples/tree/master/ecs/simple) | Quick sync by rolling out the new version and switching all traffic to it. |
| -- | [canary](https://github.com/pipe-cd/examples/tree/master/ecs/canary) | Deployment pipeline with canary strategy. |
| -- | [bluegreen](https://github.com/pipe-cd/examples/tree/master/ecs/bluegreen) | Deployment pipeline with blue-green strategy. |
| -- | [standalone-task](https://github.com/pipe-cd/examples/tree/master/ecs/standalone-task) | Deployment pipeline for an ECS standalone task (no service required). |
| -- | [secret-management](https://github.com/pipe-cd/examples/tree/master/ecs/secret-management) | Demonstrate how to manage sensitive data by using [Secret Management](https://pipecd.dev/docs/user-guide/managing-application/secret-management/) feature. |
| -- | [attachment](https://github.com/pipe-cd/examples/tree/master/ecs/attachment) | Demonstrate how to manage insensitive data and import it into application manifests while deployment using [Attachment](https://pipecd.dev/docs/user-guide/managing-application/manifest-attachment/) feature. |
