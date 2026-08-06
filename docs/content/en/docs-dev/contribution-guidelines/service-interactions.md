---
title: "APIService, WebService, and PipedService"
linkTitle: "APIService, WebService, and PipedService"
weight: 4
description: >
  This page explains what APIService, WebService, and PipedService are and how they interact, for contributors working on the control plane server.
---

The control plane's `server` component (`cmd/controlplane`) exposes three separate gRPC services. Each one is defined by its own `.proto` file, implemented by its own handler type, and run as its own gRPC server with its own authentication mechanism. This page describes each service, how a request reaches it, and how the three services end up affecting each other even though none of them calls the others directly.

### The three services

| Service | Proto file | Handler | Who calls it |
|---|---|---|---|
| `PipedService` | `pkg/app/server/service/pipedservice/service.proto` | `grpcapi.PipedAPI` (`pkg/app/server/grpcapi/piped_api.go`) | The `piped` agent (legacy `pkg/app/piped` and `pipedv1`) and the `launcher` |
| `APIService` | `pkg/app/server/service/apiservice/service.proto` | `grpcapi.API` (`pkg/app/server/grpcapi/api.go`) | `pipectl` and other external integrations |
| `WebService` | `pkg/app/server/service/webservice/service.proto` | `grpcapi.WebAPI` (`pkg/app/server/grpcapi/web_api.go`) | The web frontend (`web/`), via grpc-web |

The intent of each service is documented directly above its `service` definition in the proto file:

- `PipedService`: "contains all RPC definitions for piped. All of these RPCs are only called by piped and authenticated by using PIPED_TOKEN."
- `APIService`: "contains all RPC definitions for external service, pipectl. All of these RPCs are authenticated by using API key."
- `WebService`: "contains all RPC definitions for web client. All of these RPCs are only called by web client and authenticated by using ID_TOKEN."

### Each service runs as its own gRPC server

`cmd/controlplane/server.go` starts three independent `*grpc.Server` instances, each on its own port and with its own auth interceptor:

| Service | Default port (flag) | Auth interceptor | Credential verified against |
|---|---|---|---|
| `PipedService` | `9080` (`--piped-api-port`) | `rpc.WithPipedTokenAuthUnaryInterceptor` | `pipedverifier`, checking `ProjectStore` and `PipedStore` |
| `APIService` | `9083` (`--api-port`) | `rpc.WithAPIKeyAuthUnaryInterceptor` | `apikeyverifier`, checking `APIKeyStore` |
| `WebService` | `9081` (`--web-api-port`) | `rpc.WithJWTAuthUnaryInterceptor` | JWT verifier + `webservice.NewRBACAuthorizer` |

A fourth process, an HTTP server on port `9082`, serves auth callbacks, webhooks, and the static web assets; it is not one of the three gRPC services but is started alongside them.

Each RPC in `webservice/service.proto` declares the RBAC resource/action it requires, e.g.:

```proto
rpc RegisterPiped(RegisterPipedRequest) returns (RegisterPipedResponse) {
    option (model.rbac).resource = PIPED;
    option (model.rbac).action = CREATE;
}
```

`JWTUnaryServerInterceptor` (`pkg/rpc/rpcauth/interceptor.go`) reads the signed token from the `token` cookie, verifies it, and calls the `RBACAuthorizer` with the RPC's full method name and the caller's role before letting the request through. `PipedService` and `APIService` RPCs have no such per-method RBAC options; they instead call helpers like `requireAPIKey` inline in the handler.

### How requests reach the right server

In front of the three gRPC servers, an Envoy proxy (configured in `manifests/controlplane/templates/envoy-configmap.yaml`) listens on port `9090` and routes each request to the matching upstream cluster based on the gRPC service name in the request path:

- `/grpc.service.pipedservice.PipedService/*` → the `PipedService` server on `9080`
- `/grpc.service.webservice.WebService/*` → the `WebService` server on `9081`
- `/grpc.service.apiservice.APIService/*` → the `APIService` server on `9083`

The `envoy.filters.http.grpc_web` filter on this listener also translates grpc-web requests (used by the browser-based web frontend, see `web/src/api/client.ts`) into plain gRPC before forwarding them to `WebService`. `ext_authz` is explicitly disabled per-route for all three services in this config, since each gRPC server already authenticates the request itself via its own interceptor.

### The services don't call each other directly

`PipedAPI`, `API`, and `WebAPI` are constructed independently in `cmd/controlplane/server.go`, but they are wired up with references to the *same* underlying datastore, cache, and filestore instances (for example, the same `applicationlivestatestore.Store` value `alss` is passed into both `grpcapi.NewPipedAPI` and `grpcapi.NewWebAPI`). None of the three services holds a client to either of the others. Instead, interaction between them happens because one service writes state that another later reads from that shared storage. Two examples grounded in the code:

#### 1. Triggering work: Web/API write a Command, Piped picks it up

When a user clicks "Sync" in the web UI, `WebAPI.SyncApplication` (`pkg/app/server/grpcapi/web_api.go`) builds a `model.Command` and saves it through the shared `commandStore`. `APIService.SyncApplication` (`pkg/app/server/grpcapi/api.go`), used by `pipectl`, does the same thing.

The `piped` agent never receives that command directly. Instead, its local `commandstore` (`pkg/app/piped/apistore/commandstore/store.go`, and the equivalent in `pipedv1`) polls `PipedService.ListUnhandledCommands` every 5 seconds (`defaultSyncInterval`), picks up the pending command, executes it, and reports the outcome back via `PipedService.ReportCommandHandled`, which updates the same command record.

The web UI, meanwhile, polls `WebService.GetCommand` every 3 seconds (`FETCH_COMMANDS_INTERVAL` in `web/src/queries/commands/use-fetch-command.ts`) until the command's status is no longer `COMMAND_NOT_HANDLED_YET`. `APIService.GetCommand` serves the equivalent purpose for `pipectl`. This request/poll/report/poll cycle — documented in the comment above `ListUnhandledCommands` in `pipedservice/service.proto` — is how a web or API action ends up being executed by a piped that the control plane cannot call into directly (piped only ever makes outbound connections to the control plane).

#### 2. Reporting state: Piped writes, Web reads

`PipedAPI.ReportApplicationLiveState` and `ReportApplicationLiveStateEvents` write the application's live state into the shared `applicationLiveStateStore`. `WebAPI.GetApplicationLiveState` reads the snapshot back out of that same store. The two handlers never call each other; they only agree on the shape and location of the stored data.

The same pattern applies to deployment status: piped reports lifecycle changes through `ReportDeploymentPlanned`, `ReportDeploymentStatusChanged`, and `ReportDeploymentCompleted` (all `PipedService` RPCs), which persist to the shared `DeploymentStore`. `WebService.ListDeployments`/`GetDeployment` and `APIService.ListDeployments`/`GetDeployment` read from that same store to show deployment status to users and external tools.

### Where to look when extending one of these services

- RPC and message definitions live in each service's `service.proto` file, under `pkg/app/server/service/{pipedservice,webservice,apiservice}/`.
- After editing a `.proto` file, run `make gen/code` to regenerate the `.pb.go` files; never edit generated files by hand.
- Handler implementations live in `pkg/app/server/grpcapi/{piped_api,web_api,api}.go`, with matching `_test.go` files in the same package.
- New Go files need the standard PipeCD license header (see `.github/copilot-instructions.md`).
