---
title: "Development"
linkTitle: "Development"
weight: 2
description: >
  This page describes how to build, test PipeCD source code at your local environment.
---

## Prerequisites

- [Go 1.17](https://go.dev/)
- [Docker Decktop](https://www.docker.com/products/docker-desktop/)

## Repositories
- [pipecd](https://github.com/pipe-cd/pipecd): contains all source code and documentation of PipeCD project.
- [examples](https://github.com/pipe-cd/examples): contains various generated examples to demonstrate how to use PipeCD.

## Commands

- `make build/backend`: builds all binaries of all backend modules.
- `make build/frontend`: builds the static files for frontend.

- `make test/backend`: runs all unit tests of backend.
- `make test/frontend`: runs all unit tests of backend.
- `make test/integration`: runs integration tests.

- `make run/pipecd`: runs Control Plane locally.
- `make run/piped`: runs Piped Agent locally.
- `make run/site`: runs PipeCD site locally (requires [hugo](https://github.com/gohugoio/hugo) with `_extended` version `0.92.1` or later to be installed).

- `make gen/code`: generate Go and Typescript code from protos and mock configs. You need to run it if you modified any proto or mock definition files.

For the full list of available commands, please see the Makefile at the root of repository.

## Docs and workaround with docs

PipeCD official site contains multiple versions of documentation, all placed under the `/docs/content/en` directory, which are:
- `/docs`: stable version docs, usually synced with the latest released version docs.
- `/docs-dev`: experimental version docs, contains docs for not yet released features or changes.
- `/docs-v0.x.x`: contains docs for specified version family (a version family is all versions which in the same major release).

Basically, we have two simple rules:
- Do not touch to the `/docs` content directly.
- Keep stable docs version synced with the latest released docs version.

Here are the flow of docs contribution regard some known scenarios:
1. Update docs that are related to a specified version (which is not the latest released version):
In such case, update the docs under `/docs-v0.x.x` is enough.
2. Update docs for not yet released features or changes:
In such case, update the docs under `/docs-dev` is enough.
3. Update docs that are related to the latest released docs version:
- Change the docs' content that fixes the issue under `/docs-dev` and `/docs-v0.x.x`, they share the same file structure so you should find the right files in both directories.
- Use `make gen/stable-docs` command to sync the latest released version docs under `/docs-v0.x.x` to `/docs`

If you find any issues related to the docs, we're happy to accept your help.
