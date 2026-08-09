# Documentation

The source files for the documentation are located in the [content](https://github.com/pipe-cd/pipecd/tree/master/docs/content) directory.

# Website

The PipeCD documentation website is built with [hugo](https://gohugo.io/) and hosted on [Netlify](https://www.netlify.com/), published at https://pipecd.dev

# Docs workflow and versioning

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
