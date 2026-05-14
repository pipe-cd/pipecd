# Documentation

The source files for the PipeCD documentation are placed in the [docs/content](https://github.com/pipe-cd/pipecd/tree/master/docs/content) directory.

# Website

The PipeCD documentation website is built with [Hugo](https://gohugo.io/) and published at https://pipecd.dev.

# Docs contribution workflow

The PipeCD official site contains multiple versions of documentation under `docs/content/en`, including:
- `docs-dev`: experimental documentation for unreleased features or changes.
- `docs-vX.Y.x`: contains docs for a specific released version family.

Use the following workflow for common documentation contribution scenarios:
1. Update docs that are related to a specified version (which is not the latest released version):
In this case, update the docs under `docs-vX.Y.x`.
2. Update docs for not yet released features or changes:
In this case, update the docs under `docs-dev`.
3. Update docs that are related to the latest released docs version:
In this case, update the corresponding files under `docs-dev` and `docs-vX.Y.x`.

If you find any issues related to the docs, we're happy to accept your help.

# How to run website locally

## Prerequisite
- [Hugo Extended 0.92.1+](https://gohugo.io/)

## Commands
Run `make run/site` at the repository root, and then access http://localhost:1313.
