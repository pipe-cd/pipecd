---
title: "Installation"
linkTitle: "Installation"
weight: 4
description: >
  Complete guideline for installing and configuring PipeCD on your own.
---

Before starting to install PipeCD, let's have a look at PipeCD's components, determine your role, and which components you will interact with while installing/using PipeCD. You're recommended to read about PipeCD's [Control Plane](../concepts/#control-plane) and [`piped`](../concepts/#piped) on the concepts page.

![](/images/architecture-overview-with-roles.png)
<p style="text-align: center;">
PipeCD's components with roles
</p>

Basically, there are two types of users/roles that exist in the PipeCD system:

- Developers/Production team: Users who use PipeCD to manage their applications' deployments. You will interact with `piped` and may or may not need to install `piped` by yourself.

- Operators/Platform team: Users who operate PipeCD for other developers. You will interact with the Control Plane and `piped`, you will be the one who installs the Control Plane and keeps it up for other `piped` instances to connect while managing their applications' deployments.

This section contains guidelines for installing PipeCD's Control Plane and `piped` step by step. You can choose what to read based on your roles.
