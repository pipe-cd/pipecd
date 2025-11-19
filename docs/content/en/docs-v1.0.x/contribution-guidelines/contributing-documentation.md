---
title: "Contribute to PipeCD Documentation"
linkTitle: "Contribute to PipeCD Documentation"
description: >
  This page describes how to contribute to the PipeCD Documentation.
---

Welcome! We are so happy you're interested in helping improve our documentation. Your contributions make the project better for everyone.

This guide explains how you can contribute to the PipeCD Documentation, which resides on our official website, https://pipecd.dev.

## Where to find the docs

Our documentation is located in the `/docs` folder within the [pipe-cd/pipecd repository](https://github.com/pipe-cd/pipecd).

The content files are written in Markdown and live inside `/docs/content/en/`. You'll notice two types of documentation folders:

* `/docs-dev/`: This is for documentation related to unreleased, in-development features.
* `/docs-v0.x.x/` (and `/docs-v1.0.x/`): These folders contain the documentation for specific released versions of PipeCD.

## How to build the docs locally

To preview your changes as you work, you must run the documentation website on your local machine.

1.  **Install Prerequisite:** You must have the **extended** version of [Hugo (v0.92.1 or higher)](https://gohugo.io/getting-started/installing/) installed.
2.  **Run the Server:** From the root of the `pipecd` repository, run the following command:
    ```bash
    make run/site
    ```
3.  **Preview:** Open your browser and go to `http://localhost:1313` to see the live-reloading site.

## How to submit your changes (The PR Process)

1.  **Create a Branch:** Create a new branch for your changes (e.g., `git checkout -b my-docs-fix`).
2.  **Make Your Changes:** Edit the necessary documentation files. If you are fixing an issue in the current documentation, remember to edit the file in both the `/docs-dev/` and the latest `/docs-vx.y.z/` folders.
3.  **Commit and Push:** Commit your changes with a clear message and push your branch to your fork.
4.  **Open a Pull Request:** Go to the PipeCD repository and open a Pull Request. In the description, please link to the issue you are fixing (e.g., `Addresses #6124`).
5.  **Review:** A maintainer will review your PR, provide feedback, and merge it.

Thank you for contributing!
