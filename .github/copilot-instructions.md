# GitHub Copilot Instructions for PipeCD

This file provides guidance to GitHub Copilot when working with code in this repository.

## Overview

PipeCD is a GitOps-style continuous delivery platform that provides unified deployment operations across multiple platforms (Kubernetes, Terraform, AWS ECS/Lambda, GCP Cloud Run, etc.). The system consists of a Control Plane server (`pipecd`), deployment agents (`piped`), a CLI tool (`pipectl`), and a React web frontend.

## Build and Test Commands

### Pre-commit Check
```bash
make check   # Runs build + lint + test + generated code check + DCO check
```

### Go
```bash
make build/go                          # Build all Go binaries to .artifacts/
make build/go MOD=controlplane         # Build single binary (controlplane, piped, pipectl, launcher)
make test/go                           # Run all Go tests
make test/go MODULES=pkg/plugin/sdk    # Test specific module
go test -run TestFooBar ./pkg/foo/...  # Run single test
make lint/go                           # Lint via Docker (golangci-lint)
make lint/go FIX=true                  # Auto-fix linting issues
```

### Plugin Development
Each plugin is a separate Go module under `pkg/app/pipedv1/plugin/`:
```bash
go -C pkg/app/pipedv1/plugin/kubernetes test -race ./...
go -C pkg/app/pipedv1/plugin/ecs test -race ./...
make build/plugin PLUGINS=kubernetes,ecs    # Build specific plugins
```

### Web (React/TypeScript)
```bash
make run/web            # Dev server at localhost:9090 with MSW mocks
make build/web          # Production build
make test/web           # Tests with coverage
make lint/web           # Lint frontend
make lint/web FIX=true  # Auto-fix linting issues
yarn --cwd web typecheck
```

### Code Generation
```bash
make gen/code           # Regenerate .pb.go and .pb.validate.go from .proto files
```

### Local Development Environment
```bash
make up/local-cluster   # Start local kind cluster + registry
make run/controlplane   # Build and deploy control plane to local cluster
# In another terminal:
kubectl port-forward -n pipecd svc/pipecd 8080
# Access UI at http://localhost:8080?project=quickstart
# Login: hello-pipecd / hello-pipecd

make run/piped CONFIG_FILE=/path/to/piped-config.yaml INSECURE=true
make down/local-cluster # Teardown everything
```

## Architecture

### Components

**Control Plane (cmd/controlplane)**
- Central gRPC server managing deployments, applications, pipeds, and users
- Stores all state in datastore (Firestore or MySQL)
- Serves web dashboard and APIs
- Located: `cmd/controlplane/`, `pkg/app/server/`

**Piped Agent (cmd/piped, cmd/pipedv1)**
- Lightweight agent running in target environments
- Polls control plane for deployments
- Executes platform-specific deployment stages
- Reports status and live state back to control plane
- Located: `cmd/piped/`, `pkg/app/piped/` (legacy), `pkg/app/pipedv1/` (plugin-based)

**CLI Tool (cmd/pipectl)**
- Manages applications, deployments, pipeds
- Handles encryption, plan previews, data migration
- Located: `cmd/pipectl/`, `pkg/app/pipectl/`

**Launcher (cmd/launcher)**
- Manages piped lifecycle and updates
- Located: `cmd/launcher/`

**Web Frontend (web/)**
- React + TypeScript dashboard
- Real-time updates via WebSocket
- Auto-generated API clients from protobuf

### Directory Structure

```
cmd/                    # Entry points for each component
pkg/
  app/                  # Application logic
    server/             # Control plane implementation
      grpcapi/          # gRPC API handlers
      service/          # Proto service definitions (pipedservice, webservice, apiservice)
      *store/           # Data stores (application, deployment, piped, etc.)
    piped/              # Legacy piped agent
    pipedv1/            # Next-gen plugin-based piped
      plugin/           # Platform plugins (each is own Go module)
    pipectl/            # CLI implementation
  model/                # Protobuf domain models (.proto → .pb.go)
  config/               # Config parsing (legacy piped v0)
  configv1/             # Config parsing (pipedv1)
  datastore/            # Database abstraction (Firestore, MySQL)
  plugin/
    sdk/                # Plugin SDK (published separately as piped-plugin-sdk-go)
    api/                # Plugin gRPC API definitions
  rpc/                  # gRPC utilities, auth interceptors
  git/                  # Git operations
  cache/                # Caching layer
  crypto/               # Encryption/decryption
web/                    # React frontend
manifests/              # Helm charts for deployment
```

### Plugin Architecture

**Plugin SDK** (`pkg/plugin/sdk/`):
- Published separately: https://github.com/pipe-cd/piped-plugin-sdk-go
- Defines interfaces: `Deployment`, `LiveState`, `PlanPreview`
- Plugins communicate with piped via gRPC
- Each legacy plugin under `pkg/app/pipedv1/plugin/` is a separate Go module

**Platform Plugins**:
- kubernetes, ecs, lambda, cloudrun, terraform, scriptrun, analysis, waitapproval, wait
- Each implements deployment stages, live state reporting, plan preview generation

### Data Flow

```
Git Push → Piped EventWatcher → CreateDeployment RPC → Control Plane
                                                     ↓
Piped polls ListNotCompletedDeployments ← Control Plane stores Deployment
       ↓
DeploymentController schedules stages
       ↓
Planner generates pipeline → Executor runs stages → Reporter updates status
                                                           ↓
                                              Control Plane updates deployment
                                                           ↓
                                              Web UI shows real-time updates
```

### Communication Protocols

- **Piped ↔ Control Plane**: gRPC (`pkg/app/server/service/pipedservice/`)
- **Web ↔ Control Plane**: gRPC-Web + WebSocket (`pkg/app/server/service/webservice/`)
- **CLI ↔ Control Plane**: gRPC (`pkg/app/server/service/apiservice/`)
- **Plugin ↔ Piped**: gRPC (`pkg/plugin/api/`)

## Multi-Module Monorepo

This repository contains multiple Go modules:

**Main module**: `/go.mod` (pipecd, piped, pipectl, launcher, helloworld)

**Plugin SDK**: `/pkg/plugin/sdk/go.mod` (published independently)

**Legacy plugins**: Each plugin has own `go.mod`:
- `pkg/app/pipedv1/plugin/kubernetes/go.mod`
- `pkg/app/pipedv1/plugin/ecs/go.mod`
- etc.

**Tools**: `tool/actions-gh-release/go.mod`, `tool/codegen/protoc-gen-auth/go.mod`

When testing/building modules:
```bash
make test/go MODULES=.,pkg/plugin/sdk,pkg/app/pipedv1/plugin/kubernetes
go -C pkg/app/pipedv1/plugin/kubernetes test ./...
```

## Key Conventions

### License Header
Every new Go file must start with this header (the year should be the year first published, not necessarily the current year):
```go
// Copyright 2024 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
```

### Banned Imports (enforced by depguard)
- `sync/atomic` → use `go.uber.org/atomic` instead
- `io/ioutil` → use `os` or `io` functions instead
- pipedv1 code must import `github.com/pipe-cd/pipecd/pkg/configv1`, NOT `pkg/config`
- Plugin code under `pkg/app/pipedv1/plugin/` must NOT import from `github.com/pipe-cd/pipecd` (main module). Only `github.com/pipe-cd/piped-plugin-sdk-go` SDK is permitted.

### Protobuf / Generated Files
- Models in `pkg/model/` defined in `.proto` files
- Compiled to `.pb.go` and `.pb.validate.go`
- Never manually edit generated files — run `make gen/code` instead

### Commits
- Sign off every commit: `git commit -s` (DCO required)
- Commit message: single sentence, present tense, capital first letter
- Example: "Add imports to Terraform plan result"
- PRs target `master` branch

## Development Workflow

### Working on Control Plane
1. Make code changes in `cmd/controlplane/` or `pkg/app/server/`
2. Run `make run/controlplane` (builds, pushes to local registry, deploys to kind cluster)
3. Port-forward: `kubectl port-forward -n pipecd svc/pipecd 8080`
4. Access UI at `http://localhost:8080?project=quickstart`

### Working on Piped
1. Make code changes in `cmd/piped/` or `pkg/app/piped/`
2. Create piped config file (see CONTRIBUTING.md for template)
3. Run `make run/piped CONFIG_FILE=/path/to/config.yaml INSECURE=true`

### Working on Plugins
1. Navigate to plugin directory: `cd pkg/app/pipedv1/plugin/kubernetes`
2. Make changes
3. Test: `go test -race ./...`
4. Build: `make build/plugin PLUGINS=kubernetes` (from root)

### Working on Web Frontend
1. Make changes in `web/src/`
2. Run `make run/web` for dev server with MSW mocks
3. Test: `make test/web`
4. Typecheck: `yarn --cwd web typecheck`

### Adding New Models
1. Edit `.proto` file in `pkg/model/`
2. Run `make gen/code` to regenerate Go code
3. Commit both `.proto` and generated files

### Testing Before PR
```bash
make check   # This runs all checks that CI will run
```

## Important Technical Details

### Deployment Execution Flow
1. **Trigger**: Git event or manual trigger → `DeploymentTrigger` creates Deployment model
2. **Planning**: `Planner` generates pipeline stages based on app config
3. **Execution**: `DeploymentController` schedules stages → `Executor` runs operations
4. **Analysis**: Optional analysis stage evaluates metrics/logs
5. **Completion**: Status updates, notifications sent, metrics recorded

### Deployment Strategies
- **Kubernetes**: Canary, Blue-Green, Rolling Update
- **ECS**: Canary, Blue-Green
- **Lambda**: Canary, Linear
- **CloudRun**: Canary
- **Terraform**: Plan + Apply stages

### Datastore Abstraction
- Supports Firestore and MySQL backends
- Store interfaces: `ApplicationStore`, `DeploymentStore`, `CommandStore`, `PipedStore`
- All stores use transaction support for consistency

### Authentication
- Piped agents authenticate via `PIPED_TOKEN`
- Web users authenticate via OIDC or static admin
- API uses JWT tokens
- gRPC interceptors handle auth validation

### Stage Executors
Located in `pkg/app/piped/executor/` (legacy) or plugin modules:
- Each stage type has dedicated executor (e.g., K8S_SYNC, TERRAFORM_PLAN, ANALYSIS)
- Executors report logs and status to control plane
- Support cancellation and rollback

### Live State Reporting
- Piped continuously reports application state to control plane
- Drift detection compares desired vs actual state
- Located: `pkg/app/piped/livestatereporter/`, `pkg/app/piped/driftdetector/`

### Notifications
- Support Slack, webhooks, custom integrations
- Configured in piped config
- Located: `pkg/app/piped/notifier/`

## Tech Stack

**Backend**:
- Go 1.26+
- gRPC (service communication)
- Protobuf (data models)
- Firestore / MySQL (datastore)
- Zap (structured logging)
- Prometheus (metrics)

**Frontend**:
- TypeScript
- React
- Material-UI
- React Query (server state)
- Vite (build tool)

**Infrastructure**:
- Kubernetes (runtime platform)
- Helm (packaging)
- Docker (containerization)

**Cloud SDKs**:
- AWS SDK v2 (ECS, Lambda, S3, Secrets Manager)
- Google Cloud (Firestore, Storage, Secret Manager, Cloud Run)
- Kubernetes client-go

## Common Pitfall Prevention

### When working with modules
- If editing plugin SDK, remember it syncs to separate repo
- Don't import main module (`github.com/pipe-cd/pipecd`) from plugin code
- Use `go -C <module-path>` for plugin module operations

### When editing protos
- Always run `make gen/code` after editing `.proto` files
- Commit both proto and generated files together
- Never manually edit `.pb.go` or `.pb.validate.go` files

### When writing tests
- Use `make setup-envtest` for controller-runtime tests
- Piped tests require `KUBEBUILDER_ASSETS` env var
- Use mock interfaces from `golang/mock`

### When working with datastores
- Always use store interfaces, not concrete implementations
- Use transactions for multi-operation consistency
- Test with `datastoretest` utilities
