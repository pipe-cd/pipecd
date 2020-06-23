
# piped

## Development

### How to run a piped at your local environment with using a fake control-plane

1. Prepare a `.dev` directory at the root of repository that contains
- a `piped-key` file containing key for piped
- a `piped-config.yaml` file containing configuration for piped

2. Ensure that your kube-context is connecting to right cluster

2. Run one of the following commands to start running piped

``` console
bazelisk run --run_under="cd $PWD && " //cmd/piped:piped -- piped \
--log-encoding=humanize \
--use-fake-api-client=true \
--config-file=.dev/piped-config.yaml \
--bin-dir=/tmp/piped-bin
```

``` console
bazelisk run --run_under="cd $PWD && " //cmd/piped:piped -- piped \
--log-encoding=humanize \
--metrics \
--metrics-exporter=prometheus \
--use-fake-api-client=false \
--control-plane-address=localhost:8080 \
--bin-dir=/tmp/piped-bin \
--config-file=.dev/piped-config.yaml
```
