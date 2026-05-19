---
title: "PipeCD Plugin Development Book"
linkTitle: "Plugin Development Book"
weight: 1
description: >

 A practical, chapter-by-chapter guide to building PipeCD plugins from scratch using the `pipedv1` plugin SDK.
---

PipeCD Plugin Development Book — English translation and expansion of the original Japanese tutorial by Warashi.

This book is an English translation and expansion of the original Japanese book by Warashi (https://zenn.dev/warashi/books/try-and-learn-pipecd-plugin). It has been updated to reflect the current `pipedv1` plugin SDK and expanded with new content for English-speaking contributors.

Credit: Warashi — https://zenn.dev/warashi/books/try-and-learn-pipecd-plugin

Code examples in this book are verified against the current PipeCD codebase where possible; when uncertain, placeholders marked `[TODO: verify against current SDK]` are used.

Chapters:

1. 00. Prerequisites and Setup
2. 01. Introduction
3. 02. What We Will Build
4. 03. Technology Choices
5. 04. Project Initialization
6. 05. Adding Dependencies
7. 06. Understanding Plugin Types
8. 07. First Steps in Plugin Implementation
9. 08. Satisfying the DeploymentPlugin Interface
10. 09. Defining Configuration Types
11. 10. Stub Implementation
12. 11. Implementing FetchDefinedStages
13. 12. Implementing DetermineVersions
14. 13. Implementing DetermineStrategy
15. 14. Implementing BuildPipelineSyncStages
16. 15. Implementing BuildQuickSyncStages
17. 16. Implementing ExecuteStage
18. 17. ExecuteStage: The DIFF Stage
19. 18. ExecuteStage: The SYNC Stage
20. 19. ExecuteStage: The ROLLBACK Stage
21. 20. Updating the Main Function
22. 21. Running Your Plugin with Piped
23. 22. Conclusion and Next Steps

For more implementation resources see [Plugin development resources](./plugin-development-resources/).
