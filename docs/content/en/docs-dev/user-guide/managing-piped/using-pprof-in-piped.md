---
title: "Using Pprof in Piped"
linkTitle: "Using Pprof in Piped"
weight: 10
description: >
  This guide is for developers who want to use pprof for performance profiling in Piped.
---

Piped provides built-in support for pprof, a tool for visualization and analysis of profiling data. It's a part of the standard Go library.

In Piped, several routes are registered to serve the profiling data in a format understood by the pprof tool. Here are the routes:

- `/debug/pprof/`: This route serves an index page that lists the available profiling data.
- `/debug/pprof/profile`: This route serves CPU profiling data.
- `/debug/pprof/trace`: This route serves execution trace data.

You can access these routes to get the profiling data. For example, to get the CPU profiling data, you can access the `/debug/pprof/profile` route.  

Note that using these features in a production environment may impact performance.  

This document explains the basic usage of [pprof](https://pkg.go.dev/net/http/pprof) in Piped. For more detailed information or specific use cases, please refer to the official Go documentation.

## How to use pprof

1. Access the pprof index page
    ```bash
    curl http://localhost:9085/debug/pprof/
    ```
    This will return an HTML page that lists the available profiling data.

2. Get the Cpi Profile
    ```bash
    curl http://localhost:9085/debug/pprof/profile > cpu.pprof
    ```
    This will save the CPU profiling data to a file named cpu.pprof. You can then analyze this data using the pprof tool:
    ```bash
    go tool pprof cpu.pprof
    ```

3. Get the Execution Trace
    ```bash
    curl http://localhost:9085/debug/pprof/trace > trace.out
    ```
    This will save the execution trace data to a file named trace.out. You can then view this trace using the go tool trace command:
    ```bash
    go tool trace trace.out
    ```
    Please replace localhost:9085 with the actual address and port of your Piped's admin server.

