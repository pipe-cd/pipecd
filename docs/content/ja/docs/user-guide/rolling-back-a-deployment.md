---
title: "Deployment のロールバック"
linkTitle: "Deployment のロールバック"
weight: 2
description: >
  Deployment の自動ロールバックのタイミングや手動でのロールバックについて
---

Deployment の自動ロールバックは設定ファイルの `autoRollback` フィールドで有効にすることができます。自動ロールバックの実行条件は以下のようになっています。
- Deployment の Pipeline の任意の Stage で失敗した場合
- Analysis stage において Deployment が悪影響を及ぼすと判断した場合
- Deployment 中に何らかのエラーが起こった場合

ロールバック実行時には `ROLLBACK` Stage が Deployment の Pipeline に追加され全ての変更が取り消され Deployment 実行前の状態に戻ります。

![](/images/rolled-back-deployment.png)
<p style="text-align: center;">
ロールバックの実行例
</p>

`Cancel with rollback` ボタンをクリックして手動でロールバックを実行することもできます。
