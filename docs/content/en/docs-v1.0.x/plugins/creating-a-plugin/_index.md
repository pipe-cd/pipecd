---
title: "Creating a Plugin"
linkTitle: "Creating a Plugin"
weight: 20
description: >
  A hands-on tutorial that builds a working PipeCD plugin from scratch.
---

In PipeCD v1, every deployment runs through a plugin, and you can write your own. This tutorial builds a complete plugin step by step. You begin with an empty directory and end with a file plugin that `piped` loads and runs.

The tutorial is published in chapters. Start with the introduction, which describes what you build and how the later chapters fit together.

- [Chapter 1: Introduction](./chapter-01-introduction/)

If you want to configure an existing plugin instead of writing one, see the [official plugin pages](../official/).

## Credits

This tutorial is adapted from warashi's Zenn book, [作って学ぶ PipeCD プラグイン (Try and Learn PipeCD Plugin)](https://zenn.dev/warashi/books/try-and-learn-pipecd-plugin). It builds the same file plugin; this English version targets the v1 plugin SDK.
