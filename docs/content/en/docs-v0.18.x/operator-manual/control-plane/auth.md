---
title: "Authentication and authorization"
linkTitle: "Authentication and authorization"
weight: 6
description: >
  This page describes about PipeCD Authentication and Authorization.
---

![](/images/settings-project.png)

### Static Admin

When the PipeCD owner [adds a new project](/docs/operator-manual/control-plane/adding-a-project/), an admin account will be automatically generated for the project. After that, PipeCD owner sends that static admin information including username, password strings to the project admin, who can use that information to log in to PipeCD web with the admin role.

After logging, the project admin should change the provided username and password. Or disable the static admin account after configuring the single sign-on for the project.

### Single Sign-On (SSO)

Single sign-on (SSO) allows users to log in to PipeCD by relying on a trusted third-party service such as GitHub, GitHub Enterprise, Google Gmail, Bitbucket...

Before configuring the SSO, you need an OAuth application of the using service. For example, GitHub SSO requires creating a GitHub OAuth application as described in this page:

https://docs.github.com/en/developers/apps/creating-an-oauth-app

The authorization callback URL should be `https://YOUR_PIPECD_ADDRESS/auth/callback`.

![](/images/settings-update-sso.png)

The project can be configured to use a shared SSO configuration (shared OAuth application) instead of needing a new one. In that case, while creating the project, the PipeCD owner specifies the name of the shared SSO configuration should be used, and then the project admin can skip configuring SSO at the settings page.

### Role-Based Access Control (RBAC)

Role-based access control (RBAC) allows restricting access on the PipeCD web-based on the roles of user groups within the project. Before using this feature, the SSO must be configured.

PipeCD provides three roles:

- `viewer`: has only permissions to view application, deployment list, and details.
- `editor`: has all viewer permissions, plus permissions for actions that modify state, such as manually syncing application, canceling deployment...
- `admin`: has all editor permissions, plus permissions for updating project configurations.

Configuring RBAC means setting up 3 teams (GitHub) /groups (Google) corresponding to 3 above roles. All users belong to a team/group will have all permissions of that team/group.

![](/images/settings-update-rbac.png)
