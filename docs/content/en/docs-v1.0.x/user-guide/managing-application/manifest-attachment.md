---
title: "Manifest attachment"
linkTitle: "Manifest attachment"
weight: 10
description: >
  Attach configuration cross manifest files while deployment.
---

For insensitive data which needs to be attached/mounted as a configuration of other resources, Kubernetes ConfigMaps is a simple and bright idea. How about the other application kinds, which need something as simple as k8s ConfigMaps? PipeCD has attachment feature for your usecase.

## Configuration

Suppose you have `config.yaml` file which contains

```yaml
mysql:
  rootPassword: "test"
  database: "pipecd"
```

Then your application configuration will be configured like this

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: secret-management
  labels:
    env: example
    team: xyz
  attachment:
    sources:
      config: config.yaml
    targets:
      - taskdef.yaml
```

The `config.yaml` file is used as an attachment that can be referenced by other files. The content in `config.yaml` will be reffered to as `config`, and the target files that are configured to use the `config` can be defined under targets. In this case, it is the `taskdef.yaml` file.

And the "target" file, which uses `config.yaml` file content, can be configured as:

```yaml
...
containerDefinitions:
  - command: "echo {{ .attachment.config }}"
    image: nginx:1
    cpu: 100
    memory: 100
    name: web
...
```

In all cases, `Piped` will perform attaching the attachment file content at last, right before using it to handle any deployment tasks.

>**Tip:**
>
>This feature can be used in combo with PipeCD [SecretManagement feature](../secret-management). You can encrypt your secret data using PipeCD secret encryption function, it will be decrypted and placed in your configuration files; then the PipeCD attachment feature will attach that decrypted configuration to the manifest of resource, which requires that configuration.

<!-- See examples for detail. -->

<!-- ## Examples

- [examples/ecs/attachment](https://github.com/pipe-cd/examples/tree/master/ecs/attachment) -->