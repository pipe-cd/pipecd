---
title: "Deployment のトリガー"
linkTitle: "Deployment のトリガー"
weight: 2
description: >
  Deployment の自動トリガーや手動でのトリガーについて
---

PipeCD は Git を信頼できる唯一の情報源（全ての Application リソースが Git で宣言的かつ不変で定義されている）として扱っています。
開発者は Application やインフラに変更を加える際は Git のリポジトリにプルリクエストを作成する必要があります。
クラスター内で実行している Application やインフラにとって Git 内で定義されている状態は望ましい状態になります。

PipeCD は Deployment を実行することで、クラスターで実行しているリソースに Git で加えられた変更を適用します。 
Deployment のミッションはクラスター内の Application の全てのリソースと Git の最新のコミットの状態の差分をなくすことです。

デフォルトでは新しくマージされたプルリクエストが Application に変更を加える際に、クラスター内の Application と同期を取る為に新しい Deployment がトリガーされます。
新しい Deployment のトリガーのタイミングについて設定することも可能です。
例えば [`onOutOfSync`](/docs/user-guide/configuration-reference/#deploymenttrigger) を使うことで、ドリフト検出がされると `OUT_OF_SYNC` State を解決するために Deployment がトリガーされるようになります。

新しい Deployment はトリガーされるとエンキューされます（この時 Pipeline はまだ決まっていません）。
`Piped` エージェントはそれぞれの Application につき1つの Deployment を実行するように管理しています。
Application の Deployment が実行されていない時、 `Piped` エージェントは新たな Deployment を実行する為にデキューします。
`Piped` エージェントは Deployment の設定をもとに Pipeline を実行する条件は以下の通りです。

- マージされたプルリクエストによって Deployment 内のコンテナーイメージまたは ConfigMap や Secret に何らかの変更が加わる場合、 `Piped` エージェントはプログレッシブ Deployment を実行する為に特定の Pipeline を使います。
- マージされたプルリクエストによって `replicas` の数に変更が加わる場合、 `Piped` エージェントはリソースをスケールする為に Quick Syncを使います。

他にも、 [QuickSync](/docs/concepts/#quick-sync) やコミットメッセージをもとにして特定の Pipeline （[CommitMatcher](/docs/user-guide/configuration-reference/#commitmatcher) を設定する必要があります） を使うことができます。
Deployment は定められた Pipeline として実行され、それぞれの Stage のログなどを含んだ実行については Deployment 詳細ページにてリアルタイムで確認することができます。

![](/images/deployment-details.png)
<p style="text-align: center;">
A Running Deployment at the Deployment Details Page
Deployment 詳細ページ例
</p>

先述したようにデフォルトでは全ての Deployment はマージされたコミットを確認することで自動でトリガーされますが、Web コンソールから手動で新しい Deployment をトリガーすることもできます。
Application 詳細ページの `SYNC` ボタンをクリックすると、 master （デフォルト）ブランチの最新のコミットの状態と差分をなくす為に新しい Deployment をトリガーします。

![](/images/application-details.png)
<p style="text-align: center;">
Application 詳細ページ例
</p>
