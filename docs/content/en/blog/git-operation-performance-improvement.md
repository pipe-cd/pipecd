---
date: 2024-09-11
title: "Performance improvement in Git operations on PipeCD v0.48.9"
linkTitle: "Performance improvement in Git operations on PipeCD v0.48.9"
weight: 986
author: Shinnosuke Sawada-Dazai ([@Warashi](https://github.com/Warashi))
categories: ["Announcement"]
tags: ["News"]
---

We have implemented performance improvements for Git operations in PipeCD v0.48.9, which was released on September 10, 2024.
In this post, we will explain how we enhanced the performance of Git operations in PipeCD and when you can expect to see these improvements in your daily use.

### Background

The PipeCD agent (piped) maintains a git repository cache in the local filesystem to accelerate Git operations.
When piped needs to clone a repository, it first checks if the repository is already cloned in the cache.
If the repository is already cloned, piped will fetch the latest changes from the remote repository.
If the repository is not cloned, piped will clone it from the remote repository.
Piped uses sync.Mutex to ensure that only one goroutine can access the cache at a time.

### Problem

When multiple deployments are triggered simultaneously within a single repository, piped waits for other goroutines to finish their Git operations
because sync.Mutex ensures that only one goroutine can access the cache at a time.
This creates a bottleneck in Git operations, causing delays in deployments.

### Solution

We have improved the performance of Git operations by implementing singleflight.Group from the golang.org/x/sync/singleflight package.
The singleflight.Group is a package for duplicate function call suppression and caching.
It provides a mechanism to suppress duplicate function calls, meaning that if multiple goroutines call the same function with the same key, only one call will be made to the function.
The singleflight.Group also provides a caching mechanism, where the result of the function call is cached and returned to other goroutines that call the same function with the same key.
We use this singleflight.Group to ensure that only one goroutine can clone the repository from the remote repository, while other goroutines wait for the result of the function call.
This makes Git operations faster and more efficient, preventing delays in deployments.

### Conclusion

We have enhanced the performance of Git operations in PipeCD v0.48.9 by implementing singleflight.Group from the golang.org/x/sync/singleflight package.
This improvement makes Git operations faster and more efficient, preventing delays in deployments even when multiple deployments are triggered simultaneously within a single repository.
You can expect to experience these improvements in your daily use of PipeCD v0.48.9.
