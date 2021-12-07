---
title: "機密データの管理"
linkTitle: "機密データの管理"
weight: 6
description: >
  Git リポジトリ内の機密データの管理方法について
---

全てを Git で管理するのは便利ですが、Kubernetes の Secret や Terraform の Credential などといった機密データを Git に直接格納するのは安全ではありません。
ここでは、そういった機密情報をどのように Git 上で安全に管理できるかについて説明します。

基本的には以下のフローになっています。
- PipeCD の Web コンソールで暗号化した機密情報を Git に格納
- `Piped` は Deployment のタスクを実行する際に復号化

## 前提条件

この機能を使う為には公開鍵と秘密鍵のペアが必要になります。

``` console
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out private-key
openssl pkey -in private-key -pubout -out public-key
```

作成後、`Piped` の[インストール](http://localhost:1313/docs/operator-manual/piped/installation/#installing-on-a-kubernetes-cluster)時に以下のオプションを追加します。

``` console
--set-file secret.secretManagementKeyPair.publicKey.data=PATH_TO_PUBLIC_KEY_FILE \
--set-file secret.secretManagementKeyPair.privateKey.data=PATH_TO_PRIVATE_KEY_FILE
```

Piped の設定ファイルに `secretManagement` フィールドを追加すれば準備は完了です。

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  pipedID: your-piped-id
  ...
  secretManagement:
    type: KEY_PAIR
    config:
      privateKeyFile: /etc/piped-secret/secret-management-private-key
      publicKeyFile: /etc/piped-secret/secret-management-public-key
```

## 機密データの暗号化

機密データを暗号化するために、 該当する Application の右側のオプションアイコンをクリックして、 "Encrypt Secret" を選択します。
機密情報を入力後、 "ENCRYPT" ボタンをクリックします。
暗号化されたデータが表示されるのでコピーして Git に保存してください。

![](/images/sealed-secret-application-list.png)
<p style="text-align: center;">
Application リストのページ例
</p>

<br>

![](/images/sealed-secret-encrypting-form.png)
<p style="text-align: center;">
暗号化された機密データの例
</p>

## 暗号化された機密データを Git へ格納する

暗号化された機密データを利用できるようにする為には Application の `.pipe.yaml` 内で指定する必要があります。

- `encryptedSecrets` には暗号化された機密情報の文字列を記入します
- `decryptionTargets` には 暗号化された機密情報が書かれているファイル名を記入します

``` yaml
apiVersion: pipecd.dev/v1beta1
# KubernetesApp などの Piped で定義されている Application の Kind
kind: {APPLICATION_KIND}
spec:
  encryption:
    encryptedSecrets:
      password: encrypted-data
    decryptionTargets:
      - secret.yaml
```

## 暗号化された機密データへのアクセス

Application ディレクトリにある任意のファイルは `.pipe.yaml` に書かれた機密情報にアクセスする為に `.encryptedSecrets` Context を使用することができます。
いくつかの例を以下にあげます。

- Kubernets の Secret マニフェストへのアクセス

``` yaml
apiVersion: v1
kind: Secret
metadata:
  name: simple-sealed-secret
data:
  password: "{{ .encryptedSecrets.password }}"
```

- Lambda 関数の環境変数を設定

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: LambdaFunction
spec:
  name: HelloFunction
  environments:
    KEY: "{{ .encryptedSecrets.key }}"
```

全てのケースにおいて `Piped` は暗号化された機密情報を復号化してから、Deployment タスクを実行します。

## 使用例

- [examples/kubernetes/secret-management](https://github.com/pipe-cd/examples/tree/master/kubernetes/secret-management)
- [examples/cloudrun/secret-management](https://github.com/pipe-cd/examples/tree/master/cloudrun/secret-management)
- [examples/lambda/secret-management](https://github.com/pipe-cd/examples/tree/master/lambda/secret-management)
- [examples/terraform/secret-management](https://github.com/pipe-cd/examples/tree/master/terraform/secret-management)
- [examples/ecs/secret-management](https://github.com/pipe-cd/examples/tree/master/ecs/secret-management)
