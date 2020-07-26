---
title: "Authentication and Authorization"
linkTitle: "Authentication and Authorization"
weight: 4
description: >
  This page describes about PipeCD Authentication and Authorization.
---

> WIP

How to setup PipeCD Authentication and Authorization
========================================================

PipeCD authn and authz using GitHub Oauth. you will need to create GitHub OAuth App and get your Oauth Client ID and Client Secret.

To create GitHub Oauth App, please see the [documentation](https://developer.github.com/apps/building-oauth-apps/creating-an-oauth-app/).

After you get your Oauth Client ID and Client Secret, please edit the pipecd configuration and deploy it. A configuration will be like this:

```
  github:
    baseUrl: https://github.com
    clientId: <client-id>
    clientSecret: <client-secret>
```

Role-Based Access Control
========================
