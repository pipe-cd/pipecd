---
date: 2025-05-26
title: "MCP Server for PipeCD Docs Has Been Released"
linkTitle: "MCP Server for PipeCD Docs"
weight: 981
description: ""
author: Tetsuya Kikuchi ([@t-kikuc](https://github.com/t-kikuc))
categories: ["Announcement"]
tags: ["AI"]
---

## Introduction

Finding the right information in the official PipeCD docs can sometimes be challenging due to its volume.
To address this, we have developed the "**PipeCD Docs Local MCP Server**," which allows you to perform full-text searches on the official documentation locally.  
This server automatically fetches the PipeCD docs from GitHub and provides simple APIs for full-text search via the [MCP protocol](https://modelcontextprotocol.io/introduction).

The repository is here: https://github.com/pipe-cd/docs-mcp-server

## Example

For example, when I ask Cursor how to configure PipeCD for ECS, Cursor will search docs via the MCP server like this:

![](/images/mcp-server-intro-example1.png)

## Usage

See [README](https://github.com/pipe-cd/docs-mcp-server/blob/main/README.md) for details.

After preparing npm, you can simply run the server by `npx @pipe-cd/docs-mcp-server@latest`.

For example, for Cursor editor, add the following to your `mcp.json`:

```json
{
  "mcpServers": {
    "pipe-cd.docs-mcp-server": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "@pipe-cd/docs-mcp-server@latest"
      ]
    }
  }
}
```

## MCP Tools

The MCP Server provides two tools for now:

- `search_docs`: Performs full-text search with keywords.
- `read_docs`: Retrieves the content of a specified page. 

## Development Highlights

- **Efficient cloning with sparse checkout**  
  Only the necessary directory (`docs/content/en`) is cloned to speed up processing.

- **Simple Search Logic**  
  The current search logic is simple: it prioritizes title matches, followed by content matches.

## Conclusion

We hope this tool helps PipeCD users to configure or find information more efficiently.
Contributions, including bug reports and feature requests are welcome via [Issues](https://github.com/pipe-cd/docs-mcp-server/issues).
