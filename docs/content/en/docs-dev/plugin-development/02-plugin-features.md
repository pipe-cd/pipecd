---
title: "Plugin Features to Implement"
weight: 2
description: >
  Determining the features and stages our custom plugin will provide.
---

First, let's decide what kind of features our custom plugin will provide.

In this book, we will implement a plugin that manages files on the local machine where `piped` is running. While `piped` is often run in containerized environments (such as Kubernetes), it can also be run directly on VMs or bare-metal machines. Practically, a local file-system plugin like this could be used to deploy configuration files onto a VM—taking over some of the tasks typically handled by tools like Chef or Ansible.

### Specific Features to Implement

Each PipeCD Application will have a designated target deployment directory. The plugin's job is to copy files located under the same directory as the `app.pipecd.yaml` configuration file from the Git repository directly into that target directory.

We will define two deployment stages for this plugin:

1. **`FILE_DIFF`**: This stage compares the files in the Git repository with the files actually deployed in the target directory, printing the differences to the deployment log. This is similar to running `terraform plan` or `kubectl diff`.
2. **`FILE_SYNC`**: This stage executes the actual sync operation by copying all files from the Git repository to the target directory and deleting any orphaned files in the target directory that do not exist in the Git repository.
