---
title: "Contribute to PipeCD Blogs"
linkTitle: "Contribute to PipeCD Blogs"
description: >
  A guide for how to submit and contribute blog posts to the PipeCD website.
weight: 4
---

PipeCD accepts technical articles and community updates for publication on the official website. This guide outlines the structure and process for submitting a new blog post.

## 1. Location and Naming
Blog post files are Markdown (`.md`) and should be created in the following location:

/docs/content/en/blog/


The filename should be descriptive and use kebab-case (e.g., `my-new-blog-post.md`).

## 2. Required Frontmatter (Metadata)
Every blog post must include metadata at the top of the file. You can copy this template and fill in the details:

```yaml
---
date: YYYY-MM-DD
title: "Your Post Title Here"
linkTitle: "Your Post Title Here"
weight: 985 # Use a high number to keep the post at the top until a new one is published.
author: Your Name ([@your-github-username](https://github.com/your-github-username))
categories: ["Announcement", "Tutorial", "Community"] # Choose one or more
tags: ["New Feature", "Tutorial", "Kubernetes"] # Use specific tags
---