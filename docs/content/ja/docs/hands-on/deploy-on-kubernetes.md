---
title: "Kubernetesへのデプロイ"
linkTitle: "Kubernetesへのデプロイ"
weight: 3
description: >
  Kubernetesクラスタの上に最初のアプリケーションをデプロイします
---

このガイドでは、PipeCDを用いてKubernetesクラスタの上にアプリケーションをデプロイします。

手元に以下のものがあることを確認します。

- Kubernetesクラスタ
- デプロイするアプリケーションの設定ファイルを管理するgitリポジトリ（以降、Config-repoと呼ぶ）

## Pipedのインストール

このセクションでは、お使いのKubernetesクラスタ上でPiped Podが正常に実行されることを目指します。

[Piped](/docs/concepts/#piped)とは、デプロイメントタスクを処理する単一のバイナリコンポーネントです。このコンポーネントはステートレスに設計されているので、単一のVMやローカルマシンでも実行できます。

### ログイン
Control-plane operatorにプロジェクトの作成を依頼し、以下の情報を聞いてください（あなたがControl-plane operatorである場合は、[ガイド](/docs/operator-manual/control-plane/adding-a-project)に従ってプロジェクトを作成してください）

- Control-planeのアドレス
- プロジェクトID
- Static admin用のユーザー名
- Static admin用のパスワード


Control-planeにアクセスし、Static adminとしてログインしてください。

（任意）PipeCDはRBACによる権限管理機構を提供しています。あなたのチームに適切な権限を与えたい場合は、[ガイド](/docs/operator-manual/control-plane/auth/#role-based-access-control-rbac)に従って設定を行ってください。

### Pipedの登録

Pipedをインストールする前に、Pipedがcontrol-planeと通信するための認証情報をWeb UIから取得します。

まず、Pipedが属する[Environment](/docs/concepts/#environment)を作成し、その後Pipedを登録します。
Environmentとはアプリケーションを論理的にグルーピングするためのものです。通常、アプリケーションの実行環境に応じて `dev`, `stg`, `prod` などが利用されます。

- `https://{CONTROL_PLANE_ADDRESS}/settings/environment` を開き、+ADDボタンを押してEnvironmentを作成してください。
- `https://{CONTROL_PLANE_ADDRESS}/settings/piped` を開き、+ADDボタンを押してPipedを登録してください。作成後に表示されるpiped-idとpiped-keyは後ほど利用するので控えておいてください。

### SSH keyの登録
PipeCDにおけるデプロイフローは、PipedがConfig-repoを監視することから始まります。Config-repoが更新されたら（例えばDeploymentマニフェストの更新）Pipedはデプロイメントを開始します。
そのためにもまずは、PipedがConfig-repoにアクセス出来るように設定しておく必要があります。

- SSHキーペアを作成してください
- 公開鍵をConfig-repoに登録してください（Githubをお使いの場合は、[デプロイキー](https://docs.github.com/en/developers/overview/managing-deploy-keys)として登録することが出来ます）
- 秘密鍵は後ほど利用するので控えておいてください

### クラスタへPipedをインストールする
お使いのクラスタ上にPipedをインストールします。前のステップで取得した以下のものが手元にあることを確認してください：

- piped-id
- piped-key
- ssh-key

[インストールガイド](/docs/operator-manual/piped/installation/#in-the-cluster-wide-mode)に従ってインストールを完了させてください。
こちらのインストールガイドはクラスタレベルのリソースを必要とすることに注意してください。
NamespaceレベルのリソースのみでPipedを構成したい場合は、[こちら](https://pipecd.dev/docs/operator-manual/piped/installation/#in-the-namespaced-mode)に従ってインストールを進めてください。

Gitホストはデフォルトで`github.com`です。Config-repoがそれ以外のホストで管理されている場合は、Piped設定ファイルに以下のフィールドを追加してください：

```yaml
git:
  host: ghe.foo.org
```

Piped Podの状態が`Running`になっていれば成功です。

## アプリケーションのデプロイ

このセクションでは、最初のアプリケーションが正常にデプロイされることを目指します。

### アプリケーションの登録

デプロイを行う前に、デプロイするアプリケーション情報をPipeCDに登録する必要があります。

[アプリケーション](/docs/concepts/#application)とは、一緒に管理されるk8sリソースの集合です。マイクロサービスアーキテクチャの場合、通常一つのマイクロサービスが一つのアプリケーションとなります。
PipeCDでは、アプリケーションのマニフェスト群とそのデプロイ設定（`.pipe.yaml`）を1つのディレクトリにまとめる必要があります。 そのディレクトリはアプリケーションディレクトリと呼ばれます。

- [ガイド](/docs/user-guide/adding-an-application/#registering-a-new-application-from-web-ui)に従って、デプロイするアプリケーションをPipeCDに登録してください
- [ガイド](/docs/user-guide/adding-an-application/#adding-deployment-configuration-file)に従って、デプロイ設定(`.pipe.yaml`)を追加してください

しばらくすると、最初のデプロイが実行されるはずです。`https://{CONTROL_PLANE_ADDRESS}/deployments` を開いて結果を確認します。ステータスが`SUCCESS`であれば成功です。

お疲れさまでした。チュートリアルは以上となります。
