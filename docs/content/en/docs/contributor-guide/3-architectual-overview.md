---
title: "Architectual Overview"
linkTitle: "Architectual Overview"
weight: 3
description: >
  This page describes the architecture of PipeCD.
---

> WIP

## Diagram

![](/images/architecture.png)

## Components

## Services

- `piped`: A component that runs inside the target cloud to execute the deployment tasks.
- `api`: A service to provide api for external service like web and hook requests.
- `web`: A service for serving static files for web.
