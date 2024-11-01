## Quickstart

This directory contains configuration files for `control-plane` and `piped` for [Quickstart](https://pipecd.dev/docs/quickstart) Guide.



### Quickstart with raw manifests

1. create namespace `pipecd`

```
kubectl create namespace pipecd
```

2. deploy Control Plane to the namespace: pipecd

```
kubectl apply -n pipecd -f ./manifests/control-plane.yaml
```

3. create piped on the Web UI

4. overwrite some values in quickstart/manifests/piped.yaml

| placeholder | description |
| ---- | ---- |
| <YOUR_PIPED_ID> | piped id |
| <YOUR_PIPED_KEY> | base64-encoded piped key |

5. deploy piped to the namespace: pipecd

```
kubectl apply -n pipecd -f ./manifests/piped.yaml
```

**About the Manifests**

The manifests directory contains raw Kubernetes manifests files. The 2 files are built using `helm template` command.

For `control-plane.yaml`

```shell
$ helm template pipecd oci://ghcr.io/pipe-cd/chart/pipecd --version v0.48.6 -n pipecd -f quickstart/control-plane-values.yaml
```

For `piped.yaml`

```shell
$ helm template piped oci://ghcr.io/pipe-cd/chart/piped --version v0.48.6 -n pipecd --set quickstart.enabled=true --set quickstart.pipedId=\<YOUR_PIPED_ID\> --set quickstart.pipedKeyData=\<YOUR_PIPED_KEY_DATA\>
```
