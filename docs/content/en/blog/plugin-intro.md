---
date: 2024-11-28
title: "Overview of the Plan for Pluginnable PipeCD"
linkTitle: "Overview of the Plan for Pluginnable PipeCD"
weight: 985
author: Tetsuya Kikuchi ([@t-kikuc](https://github.com/t-kikuc))
categories: ["Announcement"]
tags: ["Plugin", "Development"]
---


<!-- TODO -->
In this article, our big progress "Pluginnable PipeCD"

<!-- この記事では、[PipeCD](https://pipecd.dev/)で絶賛開発中かつ大きな進歩である「プラグイン化」の概要を紹介します。 -->

**Note:** The content may be change in the future since it's still under development.

## Summary

- In the Pluginnable PipeCD, Plugins execute deployments instead of Piped's core.
<!-- - プラグイン化とは、デプロイ処理をPiped本体ではなく**プラグイン**が行うようにする構想 -->
- Anyone can develop Plugins and **you can use someone's Plugins**.
- It will enable users to deploy to **more various platform** with **more various Stages**.
  - platform example:  Cloudflare, Azure, and AWS CloudFormation/CDK
  - stage example:  [Automated Analysis](https://pipecd.dev/docs/user-guide/managing-application/customizing-deployment/automated-deployment-analysis/) by New Relic or Amazon CloudWatch


## What's the Pluginnable Architecture?

### Role of Plugins

まず「プラグイン」がPipeCD全体でどこに位置付けられるかというと、「Pipedのデプロイ処理を代理する」立ち位置になります。

**現在: Pipedが各デプロイ先へのデプロイ全てを担っている**
![h:3.3em w:18em](/images/pipecd-plugin-intro/mechanism-cur.drawio.png)
*現在のPiped*

**プラグイン対応後: 各デプロイ先へのデプロイはプラグインが担う**
Piped本体は主にデプロイのフロー制御を担います。

![h:6.1em w:18em](/images/pipecd-plugin-intro/mechanism-new.drawio.png)
*Piped本体とプラグインとで分担*

### Where Plugins come from?


Plugins are fetched  not PipeCD built-in code.

各プラグインはPipeCD組込みのコードではなく、外部から取得します。
GitHubなどにプラグインのバイナリが配置され、それをPipedが起動時にロードする形式です。

設定イメージ:
```yaml
# Pipedのconfig YAMLファイル
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
    plugin: # [新設項目]
       - name: k8s_plugin
         # [ここに指定] ※ URLの指定方法は今後変更の可能性もあります
         sourceURL: https://github.com/xxxx/k8s-plugin
         ...
```

ポイントとして、プラグインは自作でもよいし、**他の誰かが開発したものも利用可能**です。

なお、現在PipeCDがサポートしているK8s/CloudRun/Terraform/ECS/Lambdaへのデプロイについては、公式プラグインとしてサポートします。

### How does a Plugin run?

A Plugin runs as a gRPC server.
On a Piped starts, the Piped core will load the binary of each specified plugin and launch it as a gRPC server.
While a deployment, a Plugin communicate with Piped core.

![h:11em](/images/pipecd-plugin-intro/running.drawio.png)
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

Plugins enable you to deploy to the other platforms. For example, Cloudflare, Azure, and Amazon CloudFormation/CDK.

### 2. Users can choose behaviors different from built-in ones
<!-- - 例）標準ではK8sはXXなデプロイだが、YYな挙動でデプロイしたい -->

<!-- For example, by PipeCD default -->


### 3. Users can develop and use custom Stages
- Example 1: Analyze by New Relic or Amazon CloudWatch in the `ANALYSIS` stage.
- Example 2: Run a Kubernetes Job or a Lambda function to execute a job.

These effects will be more powerful as more Plugins come out since you can use plugins made by others.

## Schedule overview

![](/images/pipecd-plugin-intro/schedule.drawio.png)

Pluginnable PipeCD will be released around Feburary 2025.

The Plugin development guide will be published after March 2025. After that, new Plugins will come out.

For the current PipeCD users, please migrate to the Pluginnable version.  New features will be added to the Pluginnable version.

We will share the details of migration ways later.  We've been making efforts to minimize the amount of the migration tasks. The main tasks will be modifying Piped configs and Application configs.

## Conclusion

The Pluginnable Architecture is a big step toward the PipeCD's vision - "**The One CD for All {applications, platforms, operations}**".  The community energy will be more important.

To see the development progress, visit [the RFC](https://github.com/pipe-cd/pipecd/blob/master/docs/rfcs/0015-pipecd-plugin-arch-meta.md) or [the Meta Issue](https://github.com/pipe-cd/pipecd/issues/5259).

If you have any questions or opinions, feel free to talk in [CNCF Slack](https://cloud-native.slack.com/) > #pipecd-jp or [PipeCD Community Meeting](https://docs.google.com/document/d/1AtE0CQYbUV5wLfvAcl9mo9MyTCH52BuU7AngVUvE7vg/edit).
