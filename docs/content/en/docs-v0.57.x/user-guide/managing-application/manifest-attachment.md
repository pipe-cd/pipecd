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
kind: ECSApp
spec:
  name: secret-management
  labels:
    env: example
    team: xyz
  input:
    ...
  attachment:
    sources:
      config: config.yaml
    targets:
      - taskdef.yaml
```

The configuration says that: The file `config.yaml` will be used as an attachment for others, its content will be referred as `config`. The target files, that can use the `config.yaml` file as an attachment, are currently configured to `taskdef.yaml` file.

And in the "target" file, which uses `config.yaml` file content

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

__Tip__:

This feature can be used in combo with PipeCD [SecretManagement feature](../secret-management). You can encrypt your secret data using PipeCD secret encryption function, it will be decrypted and placed in your configuration files; then the PipeCD attachment feature will attach that decrypted configuration to the manifest of resource, which requires that configuration.

See examples for detail.

## Examples

- [examples/ecs/attachment](https://github.com/pipe-cd/examples/tree/master/ecs/attachment)
