---
title: "コンセプト"
linkTitle: "コンセプト"
weight: 2
description: >
  PipeCDにおける重要な幾つかのコンセプト
---

### Piped

`piped` (’d’は’daemon’の略称です。)はDeploymentタスクを実行することのできるシングルバイナリコンポーネントです。
PodまたはDeploymentを開始するだけでKubernetesクラスタ内で起動させることができます。
ステートレスに設計されているため、一つの仮想マシンやお使いの物理マシンでも実行することができます。

### Control Plane

Deploymentのデータを管理し、`piped`に接続する為のgRPC APIを提供する集中型コンポーネントです。認証、DeploymentやApplicationのリスト・詳細の表示と言った、PipeCDのWeb機能と同様の機能を搭載しています。

### Project

複数のApplicationやEnvironmentを持つことができる論理的なグループです。
また、プロジェクトごとに複数の`piped`インスタンスを異なるクラウドやEnvironmentから持つことができます。
プロジェクトの権限には3つの種類があります。
- **Viewer** はDeploymentとApplicationを見ることしかできません。
- **Editor** はViewerの権限に加えて、Deploymentにおける手動のトリガーやキャンセルといった変更を行うことができます。
- **Admin** はEditorの権限に加えてプロジェクトのデータや`piped`を管理することができます。

### Application

リソースの集合や設定があわせて管理されているもので、`KUBERNETES`や`TERRAFORM`が例として挙げられます。

### Environment

一つのプロジェクトに属する複数のApplicationを持つことができる論理的なグループのことです。
Applicationは一つのEnvironmentのみに属することしかできませんが、`piped`は複数のEnvironmentに属することができます。

### Deployment

ある特定のApplicationにおいて現在の状態（デプロイ環境）と望ましい状態（Gitで管理されている環境）での差分を埋める処理を指します。
Deploymentが成功すると、対象となるコミットがデプロイ環境へと反映されます。

### Deployment Configuration

`.pipe.yaml`にはApplicationのデプロイ方法が書かれています。
各Applicationには設定ファイルがApplicationディレクトリにそれぞれ作られている必要があります。

### Application Directory

Deploymentの設定ファイル(`.pipe.yaml`)とApplicationマニフェストが格納されているGitリポジトリです。
Applicationごとに一つのApplication Directoryを持っています。

### Quick Sync

Quick Syncを使えば、最先端のStrategyや手動によるアプルーブなしでApplicationのDeploymentを素早く行うことができます。
いくつかの例を以下にあげます。
- Kubernetesの場合全てのManifestsを適用するだけです。
- Terraformの場合検出された全ての変更を自動で適用します。
- CloudRun・Lambdaの場合新しいバージョンのロールアウトと全てのトラフィックのルーティングをします。

### Pipeline

Applicationのデプロイ方法を`piped`に指示する為の設定ファイルの中に書かれるStageのリストです。
Pipelineが特に指定されない場合、ApplicationはQuick Syncの方法を使ってデプロイされます。

### Stage

Deploymentプロセスにおいて現在の状態と望ましい状態の間に位置する一時的な中間状態を表します。

### Cloud Provider

Cloud Providerはどのクラウドを用いてどこにApplicationをデプロイすべきかを定義します。

### Analysis Provider

`Prometheus`, `Datadog`, `Stackdriver`, `CloudWatch`などのメトリクスやログを提供する外部のDeployment監視ツールのことです。
[Automated deployment analysis](/docs/user-guide/automated-deployment-analysis/)に関する話題で主に使用されます。
