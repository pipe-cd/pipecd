---
title: "Configuring a Plugin"
linkTitle: "Configuring a Plugin"
weight: 2
description: >
  This page describes how to configure a plugin in PipeCD.
---

Starting PipeCD v1, you can deploy your application to multiple platforms using plugins.

A plugin represents a deployment capability (like Kubernetes, Terraform, etc.). Each plugin can have one or more `deployTargets`, where a deploy target represents the environment where your application will be deployed.

Currently, the official plugins maintained by the PipeCD Maintainers are:

- Kubernetes
- Terraform
- Analysis
- ScriptRun
- Wait
- Wait Approval

We are working towards releasing more plugins in the future.

>**Note:**
> We also have the [PipeCD Community Plugins repository](https://github.com/pipe-cd/community-plugins) for plugins made by the PipeCD Community.

A plugin is added to the piped configuration inside the `spec.plugins` array and providing the plugin’s executable URL, the port it should run on, and any deploy targets that belong to it. For more details, see the [configuration reference for plugins](../configuration-reference/#plugins).

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  repositories:
  plugins:
    - name: plugin_name
      port: 7001
      url: url_to_plugin_binary
      deployTargets:
        - name:
          config: {}
```

Check out the latest plugin releases on [GitHub](https://github.com/pipe-cd/pipecd/releases).

---

> **Note:** Detailed configuration guides for specific plugins have been moved to the [Plugins](../../plugins/) section.
