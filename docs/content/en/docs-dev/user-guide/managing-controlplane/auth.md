---
title: "Authentication and authorization"
linkTitle: "Authentication and authorization"
weight: 3
description: >
  This page describes about PipeCD Authentication and Authorization.
---

![](/images/settings-project-v0.38.x.png)

### Static Admin

When the PipeCD owner [adds a new project](/docs/user-guide/managing-controlplane/adding-a-project/), an admin account will be automatically generated for the project. After that, PipeCD owner sends that static admin information including username, password strings to the project admin, who can use that information to log in to PipeCD web with the admin role.

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

PipeCD provides three built-in roles:

- `Viewer`: has only permissions to view existing resources or data.
- `Editor`: has all viewer permissions, plus permissions for actions that modify state, such as manually syncing application, canceling deployment...
- `Admin`: has all editor permissions, plus permissions for updating project configurations.

The policies of the role are represent like this format.
```
resources=RESOURCE_NAMES;actions=ACTION_NAMES
```

#### Actions for all resources
Currently, PipeCD has provided these resources and actions, and this table shows the being used actions per a resource.

| resource | get | list | create | update | delete |
|:--------------------|:------:|:-------:|:-------:|:-------:|:-------:|
| application     | ○ | ○ | ○ | ○ | ○ |
| deployment      | ○ | ○ |   | ○ |   |
| event           |   | ○ |   |   |   |
| piped           | ○ | ○ | ○ | ○ |   |
| project         | ○ |   |   | ○ |   |
| apiKey          |   | ○ | ○ | ○ |   |
| insight         | ○ |   |   |   |   |

The `*` represents all resources and all action for a resource.

User Group represents a relation with a specific team (GitHub)/group (Google) and an arbitrary role. All users belong to a team/group will have all permissions of that team/group.

You cannot assign multiple roles to a team/group.

![](/images/settings-add-user-group.png)
