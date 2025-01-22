---
title: "Authentication and authorization"
linkTitle: "Authentication and authorization"
weight: 3
description: >
  This page describes about PipeCD Authentication and Authorization.
---

![](/images/settings-project-v0.38.x.png)

### Static Admin

When the PipeCD owner [adds a new project](../adding-a-project/), an admin account will be automatically generated for the project. After that, PipeCD owner sends that static admin information including username, password strings to the project admin, who can use that information to log in to PipeCD web with the admin role.

After logging, the project admin should change the provided username and password. Or disable the static admin account after configuring the single sign-on for the project.

### Single Sign-On (SSO)

Single sign-on (SSO) allows users to log in to PipeCD by relying on a trusted third-party service.

**Supported service**

- GitHub
- Generic OIDC

> Note: In the future, we want to support such as Google Gmail, Bitbucket...

#### Github

Before configuring the SSO, you need an OAuth application of the using service. For example, GitHub SSO requires creating a GitHub OAuth application as described in this page:

https://docs.github.com/en/developers/apps/creating-an-oauth-app

The authorization callback URL should be `https://YOUR_PIPECD_ADDRESS/auth/callback`.

![](/images/settings-update-sso.png)

#### Generic OIDC

PipeCD supports any OIDC provider, with tested providers including Keycloak, Auth0, and AWS Cognito. The only supported authentication flow currently is the Authorization Code Grant.

Requirements:

- The IdToken will be used to decide the user's role and username.
- The IdToken must contain information about the Username and Role.
  - Supported Claims Key for Username (in order of priority): `username`, `preferred_username`,`name`, `cognito:username`
  - Supported Claims Key for Role (in order of priority): `groups`, `roles`, `cognito:groups`, `custom:roles`, `custom:groups`
  - Supported Claims Key for Avatar (in order of priority): `picture`, `avatar_url`

Provider Configuration Examples:

##### Keycloak

- **Client authentication**: On
- **Valid redirect URIs**: `https://YOUR_PIPECD_ADDRESS/auth/callback`
- **Client scopes**: Add a new mapper to the `<client-id>-dedicated` scope. For instance, map Group Membership information to the groups claim (Full group path should be off).

- **Control Plane configuration**:

  ```yaml
  apiVersion: "pipecd.dev/v1beta1"
  kind: ControlPlane
  spec:
    sharedSSOConfigs:
      - name: oidc
        provider: OIDC
        oidc:
          clientId: <CLIENT_ID>
          clientSecret: <CLIENT_SECRET>
          issuer: https://<KEYCLOAK_ADDRESS>/realms/<REALM>
          redirect_uri: https://<PIPECD_ADDRESS>/auth/callback
          scopes:
            - openid
            - profile
  ```

##### Auth0

- **Allowed Callback URLs**: `https://YOUR_PIPECD_ADDRESS/auth/callback`
- **Control Plane configuration**:

  ```yaml
  apiVersion: "pipecd.dev/v1beta1"
  kind: ControlPlane
  spec:
    sharedSSOConfigs:
      - name: oidc
        provider: OIDC
        oidc:
          clientId: <CLIENT_ID>
          clientSecret: <CLIENT_SECRET>
          issuer: https://<AUTH0_ADDRESS>
          redirect_uri: https://<PIPECD_ADDRESS>/auth/callback
          scopes:
            - openid
            - profile
  ```

- **Roles/Groups Claims**
  For Role or Groups information mapping using Auth0 Actions, here is an example for setting `custom:roles`:

  ```javascript
  exports.onExecutePostLogin = async (event, api) => {
    let namespace = "custom";
    if (namespace && !namespace.endsWith("/")) {
      namespace += ":";
    }
    api.idToken.setCustomClaim(namespace + "roles", event.authorization.roles);
  };
  ```

##### AWS Cognito

- **Allowed Callback URLs**: `https://YOUR_PIPECD_ADDRESS/auth/callback`

- **Control Plane configuration**:

  ```yaml
  apiVersion: "pipecd.dev/v1beta1"
  kind: ControlPlane
  spec:
    sharedSSOConfigs:
      - name: oidc
        provider: OIDC
        oidc:
          clientId: <CLIENT_ID>
          clientSecret: <CLIENT_SECRET>
          issuer: https://cognito-idp.<AWS_REGION>.amazonaws.com/<USER_POOL_ID>
          redirect_uri: https://<PIPECD_ADDRESS>/auth/callback
          scopes:
            - openid
            - profile
  ```

The project can be configured to use a shared SSO configuration (shared OAuth application) instead of needing a new one. In that case, while creating the project, the PipeCD owner specifies the name of the shared SSO configuration should be used, and then the project admin can skip configuring SSO at the settings page.

### Role-Based Access Control (RBAC)

Role-based access control (RBAC) allows restricting access on the PipeCD web-based on the roles of user groups within the project. Before using this feature, the SSO must be configured.

PipeCD provides three built-in roles:

- `Viewer`: has only permissions to view existing resources or data.
- `Editor`: has all viewer permissions, plus permissions for actions that modify state, such as manually syncing application, canceling deployment...
- `Admin`: has all editor permissions, plus permissions for updating project configurations.

#### Configuring the PipeCD's roles

The below table represents PipeCD's resources with actions on those resources.

| resource    | get | list | create | update | delete |
| :---------- | :-: | :--: | :----: | :----: | :----: |
| application |  ○  |  ○   |   ○    |   ○    |   ○    |
| deployment  |  ○  |  ○   |        |   ○    |        |
| event       |     |  ○   |        |        |        |
| piped       |  ○  |  ○   |   ○    |   ○    |        |
| project     |  ○  |      |        |   ○    |        |
| apiKey      |     |  ○   |   ○    |   ○    |        |
| insight     |  ○  |      |        |        |        |

Each role is defined as a combination of multiple policies under this format.

```
resources=RESOURCE_NAMES;actions=ACTION_NAMES
```

The `*` represents all resources and all actions for a resource.

```
resources=*;actions=ACTION_NAMES
resources=RESOURCE_NAMES;actions=*
resources=*;actions=*
```

#### Configuring the PipeCD's user groups

User Group represents a relation with a specific team (GitHub)/group (Google) and an arbitrary role. All users belong to a team/group will have all permissions of that team/group.

In case of using the GitHub team as a PipeCD user group, the PipeCD user group must be set in lowercase. For example, if your GitHub team is named `ORG/ABC-TEAM`, the PipeCD user group would be set as `ORG/abc-team`. (It's follow the GitHub team URL as github.com/orgs/{organization-name}/teams/{TEAM-NAME})

Note: You CANNOT assign multiple roles to a team/group, should create a new role with suitable permissions instead.

![](/images/settings-add-user-group.png)
