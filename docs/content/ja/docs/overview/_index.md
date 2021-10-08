---
title: "概要"
linkTitle: "概要"
weight: 1
description: >
  PipeCDの概要
---

## PipeCDについて

{{% pageinfo %}}
PipeCDはGitOpsを用いてデプロイメントを安全かつ高速に行う継続的デリバリーを提供します。全てのデプロイメント操作をPull Requestで行うことができ、マルチクラウドにも対応しています。
{{% /pageinfo %}}

![](/images/architecture-overview.png)
<p style="text-align: center;">
コンポーネントアーキテクチャ
</p>

## PipeCDが選ばれる理由

**Visibility**
- シンプルで明確なデザイン
- デプロイメントごとに分かれたlog viewer
- アプリケーションの状態をリアルタイムで可視化
- Slackやwebhookによる通知機能
- リードタイムやデプロイメント頻度などのデリバリーパフォーマンスを計測するためのメトリクスを表示するInsights

**Automation**
- メトリクス、ログ、リクエストに基づいたデプロイメントのインパクトを自動計測・分析
- デプロイ中の問題発生時には素早い自動ロールバック
- Configuration driftを自動検出・通知
- コンテナイメージをpushやhelm chart publishedなどの定義されたイベントが起きた際には、新しいデプロイが自動で開始

**Safety and Security**
- シングルサインオンやロールベースアクセス制御に対応
- 認証情報の一切は外部に流出せず、Control-planeにも保存されない設計
- Piped側のみがアウトバウンドリクエストを行うことによって制限されたネットワーク環境下でも運用可能
- Secrets管理を搭載

**Multi-provider & Multi-Tenancy**
- Kubernetes、Terraform、Cloud Run、AWS Lambda、Amazon ECSなどの様々なマルチクラウドに対応
- Prometheus、Datadog、Stackdriverなどの様々な分析サービスに対応
- それぞれの環境で簡単な運用が可能
**Open Source**

- オープンソースプロジェクトとして開発
- ライセンスはAPACHE 2.0 licenseを使用。詳しくは[こちらへ](https://github.com/pipe-cd/pipe/blob/master/LICENSE)
