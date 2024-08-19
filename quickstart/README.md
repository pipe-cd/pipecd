## Quickstart

This directory contains configuration files for `control-plane` and `piped` for [Quickstart](https://pipecd.dev/docs/quickstart) Guide.


The manifests directory contains raw Kubernetes manifests files. The 2 files are built using `helm template` command.

For `control-plane.yaml`

```shell
$ helm template pipecd oci://ghcr.io/pipe-cd/chart/pipecd --version v0.48.6 -n pipecd --create-namespace -f quickstart/control-plane-values.yaml
```

For `piped.yaml`

```shell
$ helm template piped oci://ghcr.io/pipe-cd/chart/piped --version v0.48.6 -n pipecd --set quickstart.enabled=true --set quickstart.pipedId=\<YOUR_PIPED_ID\> --set secret.data.piped-key=\<YOUR_PIPED_KEY\> --set quickstart.gitRepoRemote=\<YOUR_MANIFESTS_REPO\>
```
