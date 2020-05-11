
# piped

## Development


### How to run a piped at your local environment with using a fake control-plane

1. Prepare a `.dev` directory at the root of repository that contains a `piped.key` file

2. Ensure that your kube-context is connecting to right cluster

2. Run the following command to start running piped

``` console
bazelisk run --run_under="cd $PWD && " //cmd/piped:piped -- piped \
--log-encoding=humanize \
--use-fake-api-client=true \
--project-id=local-dev-project \
--piped-id=local-dev-piped \
--piped-key-file=.dev/piped.key \
--kube-config=$HOME/.kube/config \
--config-file=pkg/config/testdata/piped/dev-config.yaml
```
