---
title: "Manual Approval Stage の追加"
linkTitle: "Manual Approval Stage の追加"
weight: 4
description: >
  Manual Approval Stage に関する説明
---

Application をプロダクション環境にデプロイする前に、誰かがその変更内容を確認したいケースがあるかもしれません。
`WAIT_APPROVAL` Stage は Deployment のプロセスにおいてある時点で特定の人物またはチームによる承認を待ち、承認後 Deployment を再開させることができます。

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
      - name: WAIT_APPROVAL
        with:
          timeout: 6h
          approvers:
            - user-abc
      - name: K8S_PRIMARY_ROLLOUT
```

上記の例では、 `K8S_PRIMARY_ROLLOUT` Stage が実行される前にユーザーID `user-abc` の承認が必要となります。

`approvers` フィールドに記入する Value （ユーザーID）は [SSOの設定](/docs/operator-manual/control-plane/auth/) によって異なります。
例えば、Github を用いた SSO の場合 user ID は Github のユーザーIDになり、Google の場合 Gmail のアカウントになります。

`approvers` フィールドが設定されていない場合、 Project の中で `Editor` または `Admin` の権限を持った人が承認できます。

`timeout` フィールドでタイムアウトの時間（デフォルトは６時間）を設定することもできます。

![](/images/deployment-wait-approval-stage.png)
<p style="text-align: center;">
WAIT_APPROVAL Stage の例
</p>
