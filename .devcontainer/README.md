# PipeCD Dev Container

Development container configuration for [GitHub Codespaces](https://github.com/features/codespaces) and VS Code [Dev Containers](https://containers.dev/).

Tool versions match CI where applicable: Go 1.26.2, Node 20.19.0, Helm 3.8.2.

## Prerequisites

- Docker available on the host (Docker Desktop, OrbStack, or Codespaces built-in Docker).
- For VS Code: install the [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension.

## Open the environment

### GitHub Codespaces

Create a codespace from this repository. A machine with at least 4 CPU and 16 GB RAM is recommended if you plan to run a local kind cluster.

### VS Code

1. Open this repository folder in VS Code.
2. Run **Dev Containers: Reopen in Container** from the command palette.
3. Wait for the image build and `post-create.sh` to finish.

Setup runs `make update/go-deps` and `make update/web-deps` automatically.

## Verify the environment

```bash
docker info
go version
make build/go
```

## Optional: local Kubernetes cluster and control plane

These steps use the same Makefile targets as a manual local setup. They are optional and not required to build or test Go code.

```bash
make up/local-cluster
kind export kubeconfig --name pipecd
kubectl get nodes
make run/pipecd
kubectl port-forward -n pipecd svc/pipecd 8080
```

Open `http://localhost:8080?project=quickstart` and sign in with username `hello-pipecd` and password `hello-pipecd`.

The `pipecd` namespace is created when you run `make run/pipecd`, not by `make up/local-cluster`.

Cleanup:

```bash
make down/local-cluster
```

## Known limitations

- Docker runs on the host via socket mount (Docker-outside-of-Docker). Rootless-only Docker setups are not supported.
- `make lint/go` and `make gen/code` use Docker on the host, same as a normal local setup.
- `make run/pipecd` builds a linux/amd64 image. This matches the existing Makefile behavior.
- A local kind cluster may be slow or fail to start on some host setups (for example OrbStack on Apple Silicon inside a dev container). If `make up/local-cluster` fails, try GitHub Codespaces or run cluster commands on the host. Pushing to `localhost:5001` has been verified from inside the dev container.
- Gitpod is not configured yet. See CONTRIBUTING.md.

## Troubleshooting

### `make up/local-cluster` fails with connection refused on port 6443

kind could not start the Kubernetes API server in time. Try:

1. `make down/local-cluster` and run `make up/local-cluster` again.
2. Allocate more memory to Docker (8 GB minimum, 16 GB recommended).
3. Use GitHub Codespaces instead of a local dev container.

### Stale kind cluster or registry

```bash
make down/local-cluster
kind delete cluster --name pipecd 2>/dev/null || true
docker rm -f kind-registry 2>/dev/null || true
```

Then run `make up/local-cluster` again.
