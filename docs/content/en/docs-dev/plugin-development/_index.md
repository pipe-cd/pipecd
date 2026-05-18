---
title: "Plugin Development Book"
linkTitle: "Plugin Development Book"
weight: 10
description: >
  A comprehensive, step-by-step guide to understanding, building, and implementing custom plugins for PipeCD.
---

Welcome to the **PipeCD Plugin Development Book**. This book is a translated and adapted version of the excellent Japanese resource _作って学ぶ PipeCD プラグイン_ (Try and Learn PipeCD Plugins) by Warashi.

In this book, you will learn the internal mechanisms of PipeCD's Pluggable Architecture and walk through building a fully functional custom plugin from scratch using the Go SDK.

### Table of Contents

1. [Introduction](01-introduction/)
2. [Plugin Features to Implement](02-plugin-features/)
3. [Technology Selection](03-tech-selection/)
4. [Initializing the Project](04-project-init/)
5. [Adding Dependency Libraries](05-dependencies/)
6. [Understanding Plugin Types](06-plugin-types/)
7. [First Step of Plugin Implementation](07-first-steps/)
8. [Satisfying the DeploymentPlugin Interface](08-deployment-plugin-interface/)
9. [Implementation: Defining Configuration Types](09-defining-config-types/)
10. [Implementation: Empty Implementation to Satisfy the Interface](10-empty-implementation/)
11. [Implementation: FetchDefinedStages](11-fetch-defined-stages/)
12. [Implementation: DetermineVersions](12-determine-versions/)
13. [Implementation: DetermineStrategy](13-determine-strategy/)
14. [Implementation: BuildPipelineSyncStages](14-build-pipeline-sync-stages/)
15. [Implementation: BuildQuickSyncStages](15-build-quick-sync-stages/)
16. [Implementation: ExecuteStage](16-execute-stage/)
17. [ExecuteStage Implementation: DIFF Stage](17-execute-stage-diff/)
18. [ExecuteStage Implementation: SYNC Stage](18-execute-stage-sync/)
19. [ExecuteStage Implementation: ROLLBACK Stage](19-execute-stage-rollback/)
20. [Modifying the Main Function](20-updating-main/)
21. [Running Locally with Piped](21-trying-with-piped/)
22. [Conclusion](22-conclusion/)
