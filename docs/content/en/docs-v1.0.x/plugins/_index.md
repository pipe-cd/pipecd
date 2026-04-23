---
title: "Plugins"
linkTitle: "Plugins"
weight: 2
description: >
  Learn more about Plugins in PipeCD v1.
---

PipeCD V1 uses a plugin-based architecture where each deployment target (Kubernetes, ECS, etc.) is handled by a dedicated plugin. Plugins are configured in the piped configuration and loaded automatically at startup.

## Available plugins

| Plugin | Description | Status |
|---|---|---|
| [ECS](./ecs/) | Deploy applications to Amazon ECS using task sets and the `EXTERNAL` deployment controller | Alpha |
