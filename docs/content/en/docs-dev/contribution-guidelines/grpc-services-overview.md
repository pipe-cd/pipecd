---
title: "gRPC services overview"
linkTitle: "gRPC services overview"
weight: 4
description: >
  This page describes the gRPC services exposed by the Control Plane `server` component, how each of them authenticates its clients, and the request patterns shared across them.
---

The Control Plane `server` component exposes three separate gRPC services. Each one is served by its own `grpc.Server` instance, listening on its own port, and protected by exactly one authentication mechanism suited to its caller. The services are defined in:

- `pkg/app/server/service/apiservice/service.proto` ([APIService](#apiservice))
- `pkg/app/server/service/webservice/service.proto` ([WebService](#webservice))
- `pkg/app/server/service/pipedservice/service.proto` ([PipedService](#pipedservice))

### Request flow

| Service        | Typical caller                 | Credential           | Transport                |
| :------------- | :----------------------------- | :------------------- | :----------------------- |
| `APIService`   | `pipectl`, external automation | API key              | `authorization` metadata |
| `WebService`   | PipeCD web browser client      | Signed JWT (session) | `token` cookie           |
| `PipedService` | `piped` / `pipedv1` agent      | Piped token          | `authorization` metadata |

Each service is registered and started independently in `cmd/pipecd/server.go`, with its own port option (`apiPort`, `webAPIPort`, `pipedAPIPort`) and its own auth interceptor (`rpc.WithAPIKeyAuthUnaryInterceptor`, `rpc.WithJWTAuthUnaryInterceptor`, `rpc.WithPipedTokenAuthUnaryInterceptor`). A request that reaches the wrong port, or presents the wrong credential type for that port, is rejected before it reaches any handler code.

Every unary request also passes through `rpc.WithRequestValidationUnaryInterceptor()`: if the request message implements a `Validate() error` method (generated from the `(validate.rules)` field options declared in the `.proto` files), that validation runs first, and a failure short-circuits the call with `InvalidArgument`.

### Authentication

Credentials are carried in the gRPC `authorization` metadata as `"<TYPE> <data>"` (see `pkg/rpc/rpcauth/auth.go`):

- `API-KEY <key>` — verified against a stored, hashed API key. On success the resolved `model.APIKey` (which carries the owning `ProjectId` and a `READ_ONLY`/`READ_WRITE` role) is attached to the request context.
- `PIPED-TOKEN <projectID>,<pipedID>,<pipedKey>` — a comma-joined triple built with `rpcauth.MakePipedToken`. On success `projectID`, `pipedID`, and `pipedKey` are attached to the request context and can be read back with `rpcauth.ExtractPipedToken`.

`WebService` does **not** use this header. Its `JWTUnaryServerInterceptor` reads the JWT from a `token` cookie instead (`pkg/rpc/rpcauth/interceptor.go`, `pkg/jwt`). That cookie is set by the HTTP login handler (`pkg/app/server/httpapi/auth_handler.go`) after a successful static-admin login or SSO callback, and it carries the user's subject and RBAC role claims. After verifying the JWT signature, the interceptor also enforces per-RPC RBAC (see [WebService](#webservice) below) before the handler runs.

### APIService

`APIService` is the surface used by `pipectl` and other external integrations. It authenticates with an API key generated from the web console (`WebService.GenerateAPIKey`) and does not use the RBAC role model — instead, each handler calls `requireAPIKey(ctx, requiredRole, logger)`, which:

- Rejects the call if no verified key is on the context (`Unauthenticated`).
- Allows the call if the key's role is `READ_WRITE`.
- Allows the call only if the key's role is `READ_ONLY` **and** the RPC only requires `READ_ONLY` (`PermissionDenied` otherwise).

Its RPCs cover application management (`AddApplication`, `GetApplication`, `ListApplications`, `UpdateApplication`, `DeleteApplication`, `Enable/DisableApplication`, `RenameApplicationConfigFile`, `UpdateApplicationDeployTargets`), deployment reads (`GetDeployment`, `ListDeployments`, `ListStageLogs`), piped registration/management (`RegisterPiped`, `GetPiped`, `UpdatePiped`, `Enable/DisablePiped`), triggering actions (`SyncApplication`, `RequestPlanPreview`, `GetPlanPreviewResults`), event ingestion (`RegisterEvent`), command lookup (`GetCommand`), and secret encryption (`Encrypt`).

Every handler additionally checks that the resource being accessed belongs to the caller's own project (`key.ProjectId == resource.ProjectId`). Unlike `WebService`, most of these ownership checks in `APIService` return `InvalidArgument` rather than `PermissionDenied` — callers should not rely on the status code alone to distinguish "not found in your project" from "malformed input".

### WebService

`WebService` backs the PipeCD web console. Beyond the JWT cookie, most RPCs declare their access requirement directly in the `.proto` file:

```protobuf
rpc ListApplications(ListApplicationsRequest) returns (ListApplicationsResponse) {
    option (model.rbac).resource = APPLICATION;
    option (model.rbac).action = LIST;
}
```

These `(model.rbac)` options are compiled by `protoc-gen-auth` into a generated `authorizer` (`pkg/app/server/service/webservice/service.pb.auth.go`, built via `webservice.NewRBACAuthorizer`). For every call, `JWTUnaryServerInterceptor` looks up the caller's project RBAC roles/policies (cached for 10 minutes) and confirms the claimed role grants the declared `resource`/`action` pair before invoking the handler; otherwise it returns `PermissionDenied`.

A few RPCs opt out of this check with `option (model.rbac).ignored = true;` (currently only `GetCommand`) — the caller still needs a valid session cookie, but no resource/action policy is evaluated.

RBAC only proves the caller's role is allowed to perform an action on a resource _kind_; it does not know which project a specific application, deployment, or piped belongs to. Handlers therefore separately compare the resource's `ProjectId` against `claims.Role.ProjectId` (directly, or through a cached lookup — see [Common request patterns](#common-request-patterns)) and return `PermissionDenied` on a mismatch.

`WebService` RPCs are grouped by domain: Piped lifecycle (`RegisterPiped`, `UpdatePiped`, `RecreatePipedKey`, `DeleteOldPipedKeys`, `Enable/DisablePiped`, `ListPipeds`, `GetPiped`, `UpdatePipedDesiredVersion`, `RestartPiped`, `ListReleasedVersions`, `ListDeprecatedNotes`), Application (`Add/Update/Enable/Disable/DeleteApplication`, `ListApplications`, `SyncApplication`, `GetApplication`, `GenerateApplicationSealedSecret`, `ListUnregisteredApplications`), Deployment (`ListDeployments`, `GetDeployment`, `GetStageLog`, `CancelDeployment`, `SkipStage`, `ApproveStage`), deployment tracing and chains, live state, Account/Project settings (static admin, SSO, RBAC roles and user groups, `GetMe`), Command lookup, API key management, Insights, and Events.

### PipedService

`PipedService` is used exclusively by the `piped`/`pipedv1` agent. There are no `(model.rbac)` options on this service — any piped that presents a valid piped token may call any RPC — but almost every handler still validates that the specific application or deployment being touched belongs to _that_ piped (`validateAppBelongsToPiped` / `validateDeploymentBelongsToPiped`), returning `PermissionDenied` otherwise.

Its RPCs fall into a few groups:

- **Piped lifecycle**: `ReportStat` (periodic metrics), `ReportPipedMeta` (startup metadata: platform providers, plugins, repositories), `GetDesiredVersion`.
- **Application sync**: `ListApplications`, `ReportApplicationSyncState`, `ReportApplicationDeployingStatus`, `Report`/`GetApplicationMostRecentDeployment`, `UpdateApplicationConfigurations`, `ReportUnregisteredApplicationConfigurations`.
- **Deployment lifecycle**: `CreateDeployment`, `GetDeployment`, `ListNotCompletedDeployments`, `ReportDeploymentPlanned`, `ReportDeploymentStatusChanged`, `ReportDeploymentCompleted`, metadata persistence (`SaveDeploymentSharedMetadata`, `SaveDeploymentPluginMetadata`, `SaveStageMetadata` — `SaveDeploymentMetadata` is deprecated in favor of the shared/plugin variants), and stage logs (`ReportStageLogs`, `ReportStageLogsFromLastCheckpoint`, `ReportStageStatusChanged`).
- **Commands**: `ListUnhandledCommands`, `ReportCommandHandled` — see [Command polling](#command-polling).
- **Application live state**: `ReportApplicationLiveState`, `ReportApplicationLiveStateEvents`.
- **Events**: `GetLatestEvent`, `ListEvents`, `ReportEventStatuses`.
- **Deployment chains**: `CreateDeploymentChain`, `InChainDeploymentPlannable`.
- **Plugin shared state**: `GetApplicationSharedObject`/`PutApplicationSharedObject`, the pipedv1 replacement for the deprecated `GetLatestAnalysisResult`/`PutLatestAnalysisResult` pair.

### Common request patterns

A few conventions repeat across all three services' handler implementations in `pkg/app/server/grpcapi/`:

- **Ownership validation before mutation.** After extracting the caller's identity, handlers fetch the target resource and compare its `ProjectId` (WebService/APIService) or `PipedId` (PipedService) against the caller's own identity, rejecting mismatches.
- **In-memory ownership caches.** Because the same ownership check would otherwise hit the datastore on every request, resource-to-owner lookups are cached with a TTL (24h in `WebAPI`/`PipedAPI`, via `memorycache.NewTTLCache`) — e.g. `appProjectCache`, `pipedProjectCache`, `deploymentProjectCache` in `WebAPI`; `appPipedCache`, `deploymentPipedCache` in `PipedAPI`.
- **Sensitive-data redaction.** `model.Piped` and `model.Project` objects call `RedactSensitiveData()` before being returned to `WebService`/`APIService` clients, stripping fields such as key hashes and SSO secrets.
- **Fire-and-poll actions.** RPCs that trigger asynchronous work on a piped (`SyncApplication`, `CancelDeployment`, `SkipStage`, `ApproveStage`, `RestartPiped`, `RequestPlanPreview`) don't perform the action inline — they create a `model.Command` and return its `command_id` immediately. See [Command polling](#command-polling) for how it's actually carried out.

### Pagination

`ListApplications`, `ListDeployments`, `ListDeploymentTraces`, `ListDeploymentChains`, and `ListEvents` (on `WebService` and/or `APIService`) share the same cursor-based pagination built on `pkg/datastore.ListOptions`/`Iterator`:

- The request carries a page size (`limit` on `APIService`, `page_size` on `WebService`) and an opaque `cursor` string.
- The response carries the next `cursor`; an empty cursor means there is no further page.
- Passing the returned cursor back as the request's `cursor` continues from where the previous page left off.

Label filtering (the `labels` field on these requests) is applied **after** the datastore query, because composite indexes aren't created for every label combination. When a `labels` filter is set, a handler may transparently issue several backing `List` calls — each time advancing the cursor — until it collects a full page of label-matching results or the backing store runs out of data, so a single logical "page" can cost more than one datastore round trip.

`PipedService.ListApplications` and `ListNotCompletedDeployments` are intentionally unpaginated (`ListNotCompletedDeploymentsResponse` still has a `cursor` field, but the handler does not use it — see the `// TODO: Support pagination` comment on `PipedAPI.ListApplications`): a piped's own managed application/deployment set is expected to be small enough to fetch in full each time.

### Command polling

Server-issued actions (deployment sync, cancel, stage skip/approve, piped restart, plan preview) don't call the piped directly — the control plane has no open connection to push to. Instead:

1. A `WebService`/`APIService`/`PipedService` handler persists a `model.Command` (via `commandstore.Store`) targeted at a specific `piped_id`.
2. Each `piped`/`pipedv1` agent runs a `commandstore` that polls `PipedService.ListUnhandledCommands` on a fixed interval (5 seconds, `defaultSyncInterval` in `pkg/app/piped/apistore/commandstore/store.go` and `pkg/app/pipedv1/apistore/commandstore/store.go`), and sorts the returned commands into per-domain lists (application, deployment, stage, plan-preview, piped) for the relevant component to pick up.
3. Once a component finishes handling its command, the piped calls `PipedService.ReportCommandHandled` with the resulting `model.CommandStatus`, any `metadata`, and optional `output` bytes (persisted through a `commandOutputPutter`).
4. Callers that need the result poll back through the same service that created the command — e.g. `APIService.GetPlanPreviewResults` re-reads the command via `GetCommand`, and once it is marked handled, fetches the stored `output` through a `commandOutputGetter` and decodes it. Until then it returns `NotFound` ("waiting for result") if still within `command_handle_timeout` (default 5 minutes), or a result carrying an `Error` field once the piped is confirmed offline or the timeout has elapsed.

Because of this poll/report cycle, an interaction triggered from the web console or `pipectl` is only guaranteed to reach the piped within roughly one polling interval, not immediately.

### Errors

All three services return standard `google.golang.org/grpc/status` errors.

- **Datastore errors** are normalized by `gRPCStoreError()` (`pkg/app/server/grpcapi/grpcapi.go`):
  - `datastore.ErrNotFound`, `filestore.ErrNotFound`, `stagelogstore.ErrNotFound` → `NotFound`
  - `datastore.ErrInvalidArgument` → `InvalidArgument`
  - `datastore.ErrAlreadyExists` → `AlreadyExists`
  - `datastore.ErrUserDefined` → `FailedPrecondition` (the underlying error message is passed through as-is)
  - anything else → `Internal`, with a generic `"Failed to <action>"` message; the real error is logged server-side but not returned to the client.
- **Authentication failures** (missing/malformed credentials, failed API key/piped token/JWT verification) return `Unauthenticated`.
- **Authorization failures** (insufficient RBAC role on `WebService`, wrong API key role, or a resource that doesn't belong to the caller's project/piped) return `PermissionDenied` — except that several `APIService` ownership checks return `InvalidArgument` instead (see [APIService](#apiservice)), so error codes alone don't consistently separate "not allowed" from "bad input" on that service.
- **Payload validation failures** (from the `(validate.rules)` field options in the `.proto` files, enforced by `RequestValidationUnaryServerInterceptor` before any handler runs) return `InvalidArgument` with a message naming the offending field.
