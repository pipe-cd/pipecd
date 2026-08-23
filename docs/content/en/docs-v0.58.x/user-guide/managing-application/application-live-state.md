---
title: "Application live state"
linkTitle: "Application live state"
weight: 7
description: >
  The live states of application components as well as their health status.
---

By default, `piped` continuously monitors the running resources/components of all deployed applications to determine their state and then sends those results to the control plane. The application state will be visualized and rendered on the application details page in real time. This helps developers see what is running in the cluster as well as their health status. The application state includes:
- visual graph of application resources/components. Each resource/component node includes its metadata and health status.
- health status of the whole application. Application health status is `HEALTHY` if and only if the health statuses of all of its resources/components are `HEALTHY`.

![](/images/application-details.png)
<p style="text-align: center;">
Application Details Page
</p>

By clicking on the resource/component node, a popup will be revealed from the right side to show more details about that resource/component.
