# Documentation

The source files for the documentation are located in the [content](https://github.com/pipe-cd/pipecd/tree/master/docs/content) directory.

# Website

The PipeCD documentation website is built with [Hugo](https://gohugo.io/) and published at https://pipecd.dev

# Docs and working with docs

PipeCD’s official site contains multiple versions of documentation, all placed under the `/docs/content/en` directory:

- `/docs-dev`: experimental docs for not-yet-released features or changes.
- `/docs-vX.Y.x`: docs for a specific released version family (e.g., `docs-v0.56.x`, `docs-v1.0.x`).

Here are the recommended flows for common documentation updates:

1. **Update docs related to an older released version (not the latest released version):**  
   Update the docs under the corresponding `/docs-vX.Y.x` directory.

2. **Update docs for not-yet-released features or changes:**  
   Update the docs under `/docs-dev`.

3. **Update docs related to the latest released docs version:**  
   Apply the change in both `/docs-dev` and the latest `/docs-vX.Y.x` directory (they share the same structure, so you can find the same page in both).

If you find any issues related to the docs, we're happy to accept your help.

# How to run the website locally

## Prerequisite

- [Hugo 0.92.1+extended](https://gohugo.io/)

## Commands

Run `make run/site` at the root directory of the repository and then access http://localhost:1313
