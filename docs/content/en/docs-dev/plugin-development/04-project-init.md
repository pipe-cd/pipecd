---
title: "Initializing the Project"
weight: 4
description: >
  Initializing the Git repository and Go module.
---

Let's begin by creating our Git repository and initializing a Go module.

Run the following commands in your terminal in a suitable working directory. Replace `<YOUR_USERNAME>` with your GitHub username or any appropriate identifier:

```console
$ git init pipecd-plugin-file
$ cd pipecd-plugin-file
$ go mod init github.com/<YOUR_USERNAME>/pipecd-plugin-file
```

Once you have initialized the project, make your initial git commit to save your progress. We highly recommend committing your work at small, logical milestones throughout this guide.
