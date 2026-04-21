# GitHub Copilot Instructions for PipeCD

PipeCD is a GitOps-style continuous delivery platform. It has two main runtime components: **Control Plane** (`cmd/pipecd`) and **Piped agent** (`cmd/piped`/`cmd/pipedv1`), plus a CLI (`cmd/pipectl`) and a React web frontend (`web/`).

## Build, Test, and Lint

### Pre-commit check (run this before submitting a PR)
```bash
make check   # runs build + lint + test + generated code check + DCO check
```

### Go
```bash
make build/go                          # build all Go binaries into .artifacts/
make build/go MOD=pipecd               # build a single binary (pipecd, piped, pipectl, launcher)
make test/go                           # run all Go tests
go test -run TestFooBar ./pkg/foo/...  # run a single test
make lint/go                           # lint via Docker (golangci-lint)
```

To test or build a specific plugin module (each plugin is its own Go module):
```bash
go -C pkg/app/pipedv1/plugin/kubernetes test -race ./...
go -C pkg/app/pipedv1/plugin/ecs test -race ./...
```

### Web (React/TypeScript with Yarn)
```bash
make run/web            # start dev server at localhost:9090 with MSW mocks
make test/web           # run tests with coverage
yarn --cwd web lint     # lint frontend
yarn --cwd web typecheck
```

### Code generation (Protobuf / API)
```bash
make gen/code           # regenerate .pb.go and .pb.validate.go (runs via Docker)
```

### Local development environment
```bash
make up/local-cluster   # start local kind cluster + registry
make run/pipecd         # build and deploy control plane to local cluster
# then: kubectl port-forward -n pipecd svc/pipecd 8080
make run/piped CONFIG_FILE=path/to/piped-config.yaml INSECURE=true
make down/local-cluster # teardown
```

## Architecture

```
cmd/
  pipecd/     Control Plane: gRPC server for piped connections, web auth, deployment management
  piped/      Legacy piped agent (platform-specific deployment logic built-in)
  pipedv1/    Next-gen piped agent (plugin-based architecture)
  pipectl/    CLI tool for interacting with control plane
  launcher/   Enables remote upgrade of the piped agent

pkg/
  model/      Protobuf-defined domain models (.proto → .pb.go + .pb.validate.go)
  config/     Application and piped configuration (legacy, for piped v0)
  configv1/   Application and piped configuration (for pipedv1/plugin arch)
  app/
    server/   Control plane application logic; gRPC services in service/
    piped/    Legacy piped agent logic
    pipedv1/  Next-gen piped agent logic
      plugin/ Platform plugins (kubernetes, ecs, cloudrun, lambda, etc.)
  rpc/        gRPC server/client utilities, interceptors
  plugin/
    api/      Plugin gRPC API definitions (protobuf)
    sdk/      Plugin SDK (Go)

web/          React + TypeScript frontend
manifests/    Helm charts for all components
```

### Control Plane ↔ Piped communication
The control plane exposes a gRPC API defined in `pkg/app/server/service/pipedservice/service.proto`. Piped agents connect to this API. Web clients use `pkg/app/server/service/webservice/service.proto`. External API consumers use `pkg/app/server/service/apiservice/service.proto`.

### Plugin architecture (pipedv1)
Each plugin (`kubernetes`, `ecs`, `cloudrun`, `lambda`, `scriptrun`, etc.) lives in `pkg/app/pipedv1/plugin/<name>/` as a **separate Go module** with its own `go.mod`. Plugins implement the SDK interfaces for `deployment`, `livestate`, and `planpreview`. At runtime they are separate binaries communicating via gRPC.

## Key Conventions

### License header
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

### Banned imports (enforced by depguard lint)
- `sync/atomic` → use `go.uber.org/atomic` instead
- `io/ioutil` → use `os` or `io` functions instead
- `pipedv1` code must import `github.com/pipe-cd/pipecd/pkg/configv1`, not `pkg/config`
- Plugin code under `pkg/app/pipedv1/plugin/` must NOT import from `github.com/pipe-cd/pipecd` (the main module). Only the `github.com/pipe-cd/piped-plugin-sdk-go` SDK is permitted.

### Protobuf / generated files
Models in `pkg/model/` are defined in `.proto` files and compiled to `.pb.go` and `.pb.validate.go`. Do not manually edit generated files — run `make gen/code` instead.

### Commits
- Sign off every commit: `git commit -s` (DCO required)
- Commit message: single sentence, present tense, capital first letter (e.g., `Add imports to Terraform plan result`)
- PRs target the `master` branch
