---
title: "Overview"
linkTitle: "Overview"
weight: 1
description: >
  This page describes how to add a new application.
---

> TBA

Before deploying an application, the application must be registered from web UI to configure what piped should handle the application or where to deploy it. An application must belong to exactly one environment and can be handled by one registered piped. Currently, PipeCD supports the following application kinds:

- Kubernetes applicaiton
- Terraform application
- CloudRun application
- Lambda application

1. Registering a new application from Web UI.
2. Adding a deployment configuration file (`.piped.yaml`) to application directory in Git.
