# 🔐 Testing OIDC with Keycloak

This setup provides a ready-to-use local Keycloak instance for testing OIDC authentication flows. It includes realm import automation, default test users, client roles, realm roles, group-based mappings, and preconfigured mappers.

---

## ✅ How to Run

```bash
make setup-local-oidc
```

You may override the following environment variables:

```bash
PIPECD_CONTROL_PLANE_ADDRESS=http://localhost:8080 \
LOCAL_KEYCLOAK_ADDRESS=http://192.168.1.100:8081 \
make setup-local-oidc
```

📦 Dependencies

- Docker
- docker-compose (v2+ recommended)
- `envsubst` (usually part of GNU `gettext`)

## ⚙️ Default Configuration

| Item                       | Value                                                   |
| -------------------------- | ------------------------------------------------------- |
| Realm Name                 | `pipecd-test-realm`                                     |
| Admin Console              | `http://localhost:8081/` or `http://<private-ip>:8081/` |
| Default Admin User         | `admin` / `password`                                    |
| Imported Users             | `Admin`, `Editor`, `Viewer` (all password: `password`)  |
| Default PipeCD OIDC Client | `pipecd-keycloak`                                       |
| Issuer URL                 | `http://<private-ip>:8081/realms/pipecd-test-realm`     |
| Redirect URI(PipeCD)       | `http://localhost:8080/auth/callback`                   |
| Imported Groups            | `Admin`, `Editor`, `Viewer`                             |
| Imported Realm Roles       | `Admin`, `Editor`, `Viewer`                             |
| Imported Client Roles      | `Admin`, `Editor`, `Viewer`                             |

### ⚠️ Note

- This configuration is for **local testing only**.
- The default `issuer` uses the host's **private IP** and port 8081 (injected via shell).
- If you plan to use **HTTPS**, be sure to configure HTTPS endpoints for both Keycloak and PipeCD.
  - You can read more about [PipeCD Keycloak Connecting Issues](https://github.com/pipe-cd/pipecd/issues/5770#issuecomment-2888943178)

---

## 👥 Roles and Groups

This test realm demonstrates **three types of role-based access control**:

| Type        | Description                         | Example Claim  |
| ----------- | ----------------------------------- | -------------- |
| Realm Role  | Assigned at realm level             | `realm_roles`  |
| Client Role | Specific to OIDC client             | `client_roles` |
| Group Role  | Inferred from user group membership | `groups`       |

### 🛠️ Switching `rolesClaimKey` in PipeCD config

You can control which claim is used by PipeCD via:

```yaml
rolesClaimKey: realm_roles | client_roles | groups
```

For example:

```yaml
rolesClaimKey: client_roles
```

### 🔧 `realm.json` Customization

There are two main ways to customize the `realm.json` used in this setup: via the Keycloak admin console or by directly editing the JSON file.

---

#### 🖥️ Using the Keycloak Console

- You can modify users, roles, groups, clients, and token mappers via the Keycloak admin UI (`http://localhost:8081/`).
- After making changes in the console, **don't forget to export your updated realm** if you want to persist or version-control it.

---

#### 📝 Editing `realm.json` Directly

You can also edit the `realm.json` file manually if you prefer Git-based configuration.

Here are the key sections:

| Purpose      | Path                              | Description                                              |
| ------------ | --------------------------------- | -------------------------------------------------------- |
| Users        | `users[*]`                        | Define test users with username, password, roles, groups |
| Realm Roles  | `roles.realm[*]`                  | Define realm-level roles (e.g. `Admin`, `Viewer`)        |
| Client Roles | `roles.client.pipecd-keycloak[*]` | Define roles specific to the OIDC client                 |
| Groups       | `groups[*]`                       | Define user groups with assigned roles                   |

#### 🔄 Role Mappings

| Type                | Path                            | Description                                        |
| ------------------- | ------------------------------- | -------------------------------------------------- |
| Realm Role Mapping  | `scopeMappings[*]`              | Maps realm roles to the client                     |
| Token Claim Mapping | `clients[*].protocolMappers[*]` | Controls how user info appears in ID/access tokens |

Example claim mappers included:

- `client_roles`
- `realm_roles`
- `groups`

Each claim is already mapped to the corresponding token via `protocolMappers`.

To change role usage in your app (e.g. PipeCD), update:

```yaml
rolesClaimKey: realm_roles | client_roles | groups
```
