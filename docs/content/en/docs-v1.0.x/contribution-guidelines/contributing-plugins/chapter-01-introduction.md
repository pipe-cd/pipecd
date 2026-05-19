---
title: "Introduction"
linkTitle: "Introduction"
weight: 3
description: >
  An introduction to the PipeCD Plugin Development Book: goals, audience, and what you'll build.
---

← Back to Book Index: ./_index.md

# Introduction

## 1 — What This Book Is About

This book walks you through building a real PipeCD plugin from scratch. It uses a learn-by-doing approach: every concept is introduced alongside working code and minimal runnable examples. By the end of the book you will have a working DeploymentPlugin that `piped` can load and execute.

## 2 — What Is PipeCD

PipeCD is a CNCF open source GitOps continuous delivery platform that focuses on safe, auditable, and automated delivery flows. It decouples control-plane responsibilities from platform-specific deployment logic by using plugins.

In the `pipedv1` architecture, each deployment platform (for example Kubernetes, ECS, or Terraform) can be implemented as a separate plugin process. Plugins run as independent gRPC servers and communicate with `piped` over localhost. This design makes PipeCD extensible: anyone can write a plugin to support a new platform or deployment strategy.

## 3 — Who This Book Is For

- Go developers who want to contribute to PipeCD
- Platform engineers who want to extend PipeCD for a custom deployment target
- Anyone curious about plugin architectures, gRPC-based integrations, and building small, testable Go programs

## 4 — What We Will Build

- A working `DeploymentPlugin` implementation that `piped` can load and execute
- The plugin will implement required SDK interfaces and provide simple stage implementations you can run locally
- Later chapters cover wiring the plugin into `piped`, running it against a local control plane, and testing behavior

## 5 — A Note on Versions

This English translation and expansion is based on the PipeCD codebase as of May 19, 2026. The original Japanese book was written in June 2025 by Warashi. This edition updates examples to match the current plugin SDK and adds material aimed at English-speaking contributors. Always check the official PipeCD documentation for the latest information: https://pipecd.dev/docs

Credit and thanks to Warashi for the original work: https://zenn.dev/warashi/books/try-and-learn-pipecd-plugin

## 6 — How to Use This Book

- Read chapters in order — each chapter builds on the previous one
- Every chapter includes working code you can copy and run locally
- The final plugin examples are available in the PipeCD examples repository: https://github.com/pipe-cd/pipe-cd-examples
- If you get stuck, ask for help in the #pipecd channel on the Cloud Native Slack or open an issue in the PipeCD repository

---

Next Chapter → ./chapter-02-what-we-will-build.md
