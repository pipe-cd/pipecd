# Documentation

The source files for the documentation is placing in [content](https://github.com/pipe-cd/pipecd/tree/master/docs/content) directory.

# Website

The PipeCD documentation website is built with [hugo](https://gohugo.io/) and hosted on [Netlify](https://www.netlify.com/), published at https://pipecd.dev

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

# Hosting

The site is hosted on Netlify with the following setup:
- **Build tool**: Hugo (extended) via `netlify.toml` configuration
- **Deploy trigger**: Automatic on push to `master` branch and version tags
- **Deploy previews**: Automatically generated for pull requests that modify `docs/`
- **Redirects**: `/docs/` redirects to the latest released version (configured in `netlify.toml`)

# How to run website locally

## Prerequisite
- [Hugo 0.148.2+extended](https://gohugo.io/)
- [Node.js 24+](https://nodejs.org/)

## Commands
1. Install Hugo theme dependencies:
```
cd docs && npm install
```
2. Run the development server:
```
hugo server
```
3. Access http://localhost:1313

# Netlify Configuration

The Netlify build configuration is defined in [`netlify.toml`](./netlify.toml). Key settings:
- **Build command**: `npm ci && hugo --gc --minify`
- **Publish directory**: `public/`
- **Redirects**: `/docs/` → latest version docs
- **Security headers**: X-Frame-Options, X-XSS-Protection, etc.

When a new release version is created, the `hack/gen-release-docs.sh` script automatically updates the redirect target in `netlify.toml`.
