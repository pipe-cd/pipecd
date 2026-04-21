---
title: "Terraform Plugin"
linkTitle: "Terraform"
weight: 20
description: >
  How to configure the Terraform plugin.
---

The Terraform plugin enables Piped to run Terraform-based deployments.
A deploy target represents a Terraform execution environment (e.g., dev/prod workspace, shared variables, drift detection settings).

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  ...
  plugins:
    - name: terraform
      port: 7002
      url: https://github.com/.../terraform_v0.3.0_linux_amd64
      deployTargets:
        - name: tf-dev
          config:
            vars:
              - "project=pipecd"
```

See [Configuration Reference for Terraform plugin](../user-guide/managing-piped/configuration-reference/#terraformplugin) for complete configuration details.