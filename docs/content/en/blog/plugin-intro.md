---
date: 2024-11-28
title: "Overview of the Plan for Pluginnable PipeCD"
linkTitle: "Overview of the Plan for Pluginnable PipeCD"
weight: 985
author: Tetsuya Kikuchi ([@t-kikuc](https://github.com/t-kikuc))
categories: ["Announcement"]
tags: ["Architecture"]
---

In this article, we introduce the overview of the "Pluginnable" feature, a significant progress currently under development.

**Note:** The content may be change in the future since it's still under development.

## Summary

- In the Pluginnable PipeCD, Plugins execute deployments instead of Piped's core.
- Anyone can develop Plugins and **you can use someone's Plugins**.
- It will enable users to deploy to **more various platform** with **more various Stages**.
  - platform example:  Cloudflare, Azure, and AWS CloudFormation/CDK
  - stage example:  [Automated Analysis](https://pipecd.dev/docs/user-guide/managing-application/customizing-deployment/automated-deployment-analysis/) by New Relic or Amazon CloudWatch


## What's the Pluginnable Architecture?

### Role of Plugins

In the Plugginable Architecture, Plugins are actors who execute deployments on behalf of a Piped.

**Current: A Piped deploys to each platform by itself**
![h:3.3em w:18em](/images/plugin-intro-mechanism-cur.drawio.png)
<!-- *Current Piped* -->

**In Plugginable version: Plugins deploy to each platform**
Piped's core will control deployment flows.

![h:6.1em w:18em](/images/plugin-intro-mechanism-new.drawio.png)
<!-- *Piped and Plugins* -->

### Where Plugins come from?

Plugins are fetched from outside a Piped. They are not PipeCD built-in code.

Plugin's binary is placed in GitHub or something, and a Piped loads it on launch.


Configuration would be like this:
```yaml
# Piped's config YAML
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
    plugin: # [New Area]
       - name: k8s_plugin
         # [HERE] (The URL format might be changed)
         sourceURL: https://github.com/xxxx/k8s-plugin
         ...
```

The key point is that you can develop your own plugins, or **use plugins developed by others**.

Addiitonally, deployments to Kubernetes, CloudRun, Terraform, ECS, and Lambda, which are currently supported by PipeCD, will be supported as official plugins.

### How does a Plugin run?

A Plugin runs as a gRPC server.
On a Piped starts, the Piped core will load the binary of each specified plugin and launch it as a gRPC server.
While a deployment, a Plugin communicate with Piped core.

![h:11em](/images/plugin-intro-running.drawio.png)
<!-- *動作時のイメージ* -->

### How to develop a Plugin?

**Anyone can develop a Plugin freely.** We will pubilsh a guide on how to develop a Plugin in the future.

Developing a Plugin will not be so tough. Here I introduce somes reasons:

- The main task is implementing what you want to execute in each stage.
  - Piped's core will control deployment flows and handle git. You don't need to deeply understand the mechanism of GitOps.
  - To be exact, if you want Plan Preview or Drift Detection for the plugin, you need to implement them.

- Plugins can be also written by languages other than Golang.

- You don't need to put your plugin in PipeCD official repository.
  - You can determine review policy and release timing as you want.

## Advantages of the Pluginnable Architecture

There are three big advantages. Thanks to them, users will not need to manage additional CD tools when adding a kind of platform or customizing behaviors.

### 1. Users can deploy to platforms other than K8s/CloudRun/Terraform/ECS/Lambda

Currently, PipeCD supports deploying to Kubernetes, CloudRun, Terraform, ECS, and Lambda.

Plugins enable you to deploy to the other platforms, for example, Cloudflare, Azure, and Amazon CloudFormation/CDK.

### 2. Users can choose behaviors as they want

- Example 1: Deploy in a way different from built-in.
- Example 2: For Kubernetes (or any platform), deploy simply for some Applications with Plugin-A, and deploy in a complicated way for the other Applications with Plugin-B.


### 3. Users can develop and use custom Stages
- Example 1: Analyze by New Relic or Amazon CloudWatch in the `ANALYSIS` stage.
- Example 2: Run a Kubernetes Job or a Lambda function to execute a job.

These effects will be more powerful as more Plugins come out since you can use plugins made by others.

## Schedule overview

![](/images/plugin-intro-schedule.drawio.png)

Pluginnable PipeCD will be released around Feburary 2025.

The Plugin development guide will be published after March 2025. After that, new Plugins will come out.

For the current PipeCD users, please migrate to the Pluginnable version.  New features will be added to the Pluginnable version.

We will share the details of migration ways later.  We've been making efforts to minimize the amount of the migration tasks. The main tasks will be modifying Piped configs and Application configs.

## Conclusion

The Pluginnable Architecture is a big step toward the PipeCD's vision - "**The One CD for All {applications, platforms, operations}**".  The community energy will be more important.

To see the development progress, visit [the RFC](https://github.com/pipe-cd/pipecd/blob/master/docs/rfcs/0015-pipecd-plugin-arch-meta.md) or [the Meta Issue](https://github.com/pipe-cd/pipecd/issues/5259).

If you have any questions or opinions, feel free to talk in [CNCF Slack](https://cloud-native.slack.com/) > #pipecd-jp or [PipeCD Community Meeting](https://docs.google.com/document/d/1AtE0CQYbUV5wLfvAcl9mo9MyTCH52BuU7AngVUvE7vg/edit).
