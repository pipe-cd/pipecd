---
title: "Authentication and authorization"
linkTitle: "Authentication and authorization"
weight: 4
description: >
  This page describes about PipeCD Authentication and Authorization.
---

### Single Sign-On (SSO)

PipeCD authn and authz using GitHub Oauth. You will need to create GitHub OAuth App and get your Oauth Client ID and Client Secret.

To create GitHub Oauth App, please see the [documentation](https://developer.github.com/apps/building-oauth-apps/creating-an-oauth-app/).

After getting your Oauth Client ID and Client Secret, please edit the pipecd configuration and deploy it. A configuration will be like this:

```
  github:
    baseUrl: https://github.com
    clientId: <client-id>
    clientSecret: <client-secret>
```

> TBA

### Role-Based Access Control (RBAC)

> TBA
