---
title: "Wait Stage の追加"
linkTitle: "Wait Stage の追加"
weight: 4
description: >
  Wait Stage に関する説明
---

`WAIT` Stage はある時点で一定時間待ってから Deployment を続行するようにします。
どのくらい待つかについては `duration` フィールドで設定できます。

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  pipeline:
    stages:
      - name: K8S_CANARY_ROLLOUT
      - name: WAIT
        with:
          duration: 5m
      - name: K8S_PRIMARY_ROLLOUT
      - name: K8S_CANARY_CLEAN
```

![](/images/deployment-wait-stage.png)
<p style="text-align: center;">
WAIT stageの例
</p>
