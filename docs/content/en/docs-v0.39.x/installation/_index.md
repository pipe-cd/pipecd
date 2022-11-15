---
title: "Installation"
linkTitle: "Installation"
weight: 4
description: >
  Complete guideline for installing and configuring PipeCD on your own.
---

Before starting to install PipeCD, let’s have a look at PipeCD’s components, determine your role, and which components you will interact with while installing/using PipeCD. You’re recommended to read about PipeCD’s [Control Plane](../concepts/#control-plane) and [Piped](../concepts/#piped) on the concepts page.

![](/images/architecture-overview-with-roles.png)
<p style="text-align: center;">
PipeCD's components with roles
</p>

Basically, there are two types of users/roles that exist in the PipeCD system, which are:
- Developers/Production team: Users who use PipeCD to manage their applications’ deployments. You will interact with Piped and may or may not need to install Piped by yourself.
- Operators/Platform team: Users who operate the PipeCD for other developers can use it. You will interact with the Control Plane and Piped, you will be the one who installs the Control Plane and keeps it up for other Pipeds to connect while managing their applications’ deployments.

This section contains the guideline for installing PipeCD's Control Plane and Piped step by step. You can choose what to read based on your roles.
