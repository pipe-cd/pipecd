# Documentation

The source files for the documentation is placing in [content](https://github.com/pipe-cd/pipecd/tree/master/docs/content) directory.

# Website

The PipeCD documentation website is built with [hugo](https://gohugo.io/) and published at https://pipecd.dev

# Docs and workaround with docs

PipeCD official site contains multiple versions of documentation, all placed under the `/docs/content/en` directory, which are:
- `/docs-dev`: experimental version docs, contains docs for not yet released features or changes.
- `/docs-v0.x.x`: contains docs for specified version family (a version family is all versions which in the same major release).

Here are the flow of docs contribution regard some known scenarios:
1. Update docs that are related to a specified version (which is not the latest released version):
In such case, update the docs under `/docs-v0.x.x` is enough.
2. Update docs for not yet released features or changes:
In such case, update the docs under `/docs-dev` is enough.
3. Update docs that are related to the latest released docs version:
Change the docs' content that fixes the issue under `/docs-dev` and `/docs-v0.x.x`, they share the same file structure so you should find the right files in both directories.

If you find any issues related to the docs, we're happy to accept your help.

# How to run website locally

## Prerequisite
- [Hugo 0.92.1+extended](https://gohugo.io/)

## Commands
Run `make run/site` at the root directory of the repository and then access http://localhost:1313
