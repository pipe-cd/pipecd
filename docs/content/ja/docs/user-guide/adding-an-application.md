---
title: "Applicationの追加"
linkTitle: "Applicationの追加"
weight: 1
description: >
  新しいApplicationの追加方法
---

Applicationはリソースの集合や設定が一緒に管理されているもので、デプロイ予定のサービスを表しています。
PipeCDの場合、全てのApplicationマニフェストとDeploymentの設定ファイル（`.pipe.yaml`）はGitリポジトリのディレクトリ（Applicationディレクトリ）にコミットされる必要があります。

Applicationはデプロイされる前に、PipeCDのWebコンソールでの登録に加えて、Deploymentの設定ファイル（`.pipe.yaml`）がApplicationディレクトリにコミットされている必要があります。
`piped`はApplicationを管理することができ、それぞれのApplicationは一つのEnvironmentに所属しなければなりません。現在、PipeCDは以下のApplicationをサポートしています。

- Kubernetes application
- Terraform application
- CloudRun application
- Lambda application
- ECS application

## Webコンソールから新しいApplicationを登録

Applicationを登録することでPipeCDはApplicationに関する基本情報を知ることができます。
具体的には、Application設定がされている場所、`piped`がどのように処理するべきか、どのクラウドにデプロイされるべきなのかなどといった情報です。

Applicationリストページにある`+ADD`ボタンをクリックすることで、以下のようにポップアップウインドウが表示されます。

![](/images/registering-an-application.png)
<p style="text-align: center;">
ポップアップウインドウの表示例
</p>

全ての入力必須項目を入力後、`Save`ボタンをクリックすればApplication登録は完了です。

以下のリストが登録フォームの入力項目です。

| フィールド | 説明 | 項目 |
|-|-|-|-|
| Name | Applicationの名前 | 必須 |
| Kind | `Kubernetes`、`Terraform`、`CloudRun`、`Lambda`、`ECS`の中から選ぶ。 | 必須 |
| Env | Applicationが所属するEnvironment。`Settings/Environment`ページで登録したものの中から選ぶ。 | 必須 |
| Piped | `Settings/Piped`ページで登録したものの中から選ぶ。 | 必須 |
| Repository | Application設定とDeployment設定が格納されたGitリポジトリ。`piped`の設定で登録したものの中から選ぶ。 | 必須 |
| Path | Repositoryを基準とした相対パス。 `./`はリポジトリルートを意味する。 | 必須 |
| Config Filename | Deploymentの設定ファイルの名前。デフォルトは`.pipe.yaml`。 | 任意 |
| Cloud Provider | Applicationがデプロイされる場所。`piped`の設定で登録したものの中から選ぶ。 | 必須 |

## Deploymentの設定ファイルの追加

Applicationの登録後はDeploymentの設定ファイル（`.pipe.yaml`）をGitリポジトリのApplicationディレクトリに追加していきます。
`piped`はこのファイルによってcanary/blue-greenストラテジーや手動による承認などといったアプリケーションのデプロイ方法について知ることができます。
Deploymentの設定ファイルは以下のようにYAMLフォーマットで書くことができます。

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: ApplicationKind
spec:
  ...
```

- `kind`にはApplication Kindを設定します。先述したように、`Kubernetes`、`Terraform`、`CloudRun`、`Lambda`、`ECS`がサポート対象のKindとなっています。
- `spec`ではそれぞれのApplication Kind特有の設定をします。

サポートされているDeploymentについては[pipecd/examples](https://pipecd.dev/docs/user-guide/examples/)をご覧ください。
