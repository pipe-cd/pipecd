---
title: "Contribute to PipeCD Blogs"
linkTitle: "Contribute to PipeCD Blogs"
description: >
  This page describes how to submit and contribute blog posts to the PipeCD website.
---


Welcome! We are happy you want to share your knowledge with the PipeCD community.


This guide explains how you can submit a blog post to the PipeCD website.


## Where to find the blogs


Our blog posts are located in the `/docs/content/en/blog/` folder within the [pipe-cd/pipecd repository](https://github.com/pipe-cd/pipecd).


## How to add a new blog post


1. **Create a file:** Create a new Markdown file in the `docs/content/en/blog/` directory. The filename should use kebab-case (e.g., `my-new-feature.md`).


2. **Add Frontmatter:** Every blog post requires a header (frontmatter) at the top of the file. Copy and customize this template:
```yaml
---
date: 2023-10-27
title: "Your Blog Post Title"
linkTitle: "Your Blog Post Title"
author: Your Name ([@your-github-username](https://github.com/your-github-username))
categories: ["Tutorial"]
tags: ["Kubernetes", "PipeCD"]
---
```
3. **Write your content:** Write your article in standard Markdown below the frontmatter.


## How to preview your blog


1. **Run the server:** From the root of the repository, run:
```bash
make run/site
```
2. **Preview:** Open `http://localhost:1313/blog/` in your browser.


## How to submit (The PR Process)


1. **Create a branch:** from the master branch, create a new branch for your post:
```bash
git checkout -b blog/my-new-post
```
2. **Commit and push:** Commit your changes and push your new file to your fork.


3. **Open a Pull Request:** Submit your PR to the PipeCD repository.


## What happens next?


After you submit your pull request:


- **Review:** Maintainers will review your post. They might suggest changes or ask clarifying questions.
- **Update:** Address any feedback by pushing new commits to your branch.
- **Merge:** Once approved, a maintainer will merge your PR. Your post will appear on the website shortly after.


## Need help?


If you have questions or need ideas for blog posts, join our community:


- **Slack:** Join the `#blog` channel on our [PipeCD Slack workspace](https://pipecd.dev/slack).
- **Community Meeting:** Join our community meeting (held bi-weekly on Wednesdays) to discuss your ideas directly with the maintainers.

