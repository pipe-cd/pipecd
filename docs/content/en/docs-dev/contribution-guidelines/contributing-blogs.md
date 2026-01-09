---
title: "Contribute to PipeCD Blogs"
linkTitle: "Contribute to PipeCD Blogs"
weight: 3
description: >
  This page describes how to contribute blog posts to PipeCD.
---

We welcome blog contributions from the community! Blog posts are a great way to share your experiences, tutorials, and insights about PipeCD with other users.

## What makes a good blog post?

Blog posts can cover a variety of topics, including:

- **Announcements**: New features, releases, or project updates
- **Tutorials**: Step-by-step guides on using PipeCD features
- **Use cases**: How you or your organization uses PipeCD
- **Best practices**: Tips and tricks for getting the most out of PipeCD
- **Integrations**: How to integrate PipeCD with other tools

## Where blog posts live

Blog posts are located in the `/docs/content/en/blog/` directory within the [pipe-cd/pipecd repository](https://github.com/pipe-cd/pipecd).

## Blog post format

Each blog post is a Markdown file with YAML front matter. Here's the structure:

```yaml
---
date: 2025-01-03
title: "Your Blog Post Title"
linkTitle: "Short Title for Navigation"
weight: 980
description: "A brief description of your post"
author: Your Name ([@your-github-handle](https://github.com/your-github-handle))
categories: ["Announcement"]
tags: ["Tag1", "Tag2"]
---

Your content here...
```

### Front matter fields

| Field | Required | Description |
|-------|----------|-------------|
| `date` | Yes | Publication date in `YYYY-MM-DD` format |
| `title` | Yes | Full title of your blog post |
| `linkTitle` | Yes | Shorter title used in navigation menus |
| `weight` | Yes | Controls ordering (lower = newer, typically use ~980-990) |
| `description` | No | Brief summary for SEO and previews |
| `author` | Yes | Your name with GitHub profile link |
| `categories` | Yes | One of: `Announcement`, `Tutorial`, `Release` |
| `tags` | No | Relevant keywords for your post |

### Content guidelines

- Use clear, concise language
- Include code examples where appropriate (use fenced code blocks with language identifiers)
- Add images to `/docs/static/images/` and reference them as `![alt text](/images/your-image.png)`
- Structure your post with headings (`##`, `###`) for readability
- Include a conclusion or summary section

## How to submit your blog post

1. **Fork and clone** the [pipecd repository](https://github.com/pipe-cd/pipecd)

2. **Create a branch** for your blog post:
   ```bash
   git checkout -b blog/your-post-title
   ```

3. **Create your blog post** file in `/docs/content/en/blog/`:
   ```bash
   touch docs/content/en/blog/your-post-title.md
   ```

4. **Add any images** to `/docs/static/images/`

5. **Preview locally** by running:
   ```bash
   make run/site
   ```
   Then visit `http://localhost:1313/blog/` to see your post.

6. **Commit and push** your changes:
   ```bash
   git add .
   git commit -s -m "blog: add post about your-topic"
   git push origin blog/your-post-title
   ```

7. **Open a Pull Request** against the `master` branch

## Review process

A maintainer will review your blog post for:

- Technical accuracy
- Clarity and readability
- Adherence to the format guidelines
- Appropriate use of images and code examples

Feel free to reach out on the [#pipecd Slack channel](https://cloud-native.slack.com/) if you have questions or want feedback on your draft before submitting.

Thank you for contributing to the PipeCD blog!
