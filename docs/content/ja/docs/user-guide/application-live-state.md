---
title: "Applicationの状態"
linkTitle: "Applicationの状態"
weight: 7
description: >
  Application のコンポーネントのヘルスステータスおよび状態の確認方法に関する説明
---

デフォルトでは、デプロイされた Application の全てのリソース（コンポーネント）が監視されています。
Application の状態はリアルタイムで可視化され、詳細ページにて表示されます。
ここでは、ヘルスステータスはもちろん、クラスター内で実行中のものについても確認することができます。
具体的には以下の内容を指します。
- 各リソースのメタデータやヘルスステータスが含まれたグラフ。
- Application 全体のヘルスステータス。全てのリソースのヘルスステータスが `HEALTHY` の場合にのみ、Application のヘルスステータスは `HEALTHY` と表示されます。

![](/images/application-details.png)
<p style="text-align: center;">
Application の詳細ページ例
</p>

リソースのノードをクリックすることで、詳細情報を確認することができます。
