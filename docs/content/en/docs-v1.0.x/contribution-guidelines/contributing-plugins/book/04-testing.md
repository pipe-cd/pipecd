---
title: "Testing and Debugging"
linkTitle: "Testing and Debugging"
weight: 40
description: >
  Running and verifying your plugin locally.
---

Testing a plugin involves building it and then telling a `piped` agent to use it.

## 1. Build your Plugin

Since plugins are standalone binaries, you can build them using standard Go commands:

```bash
go build -o my-plugin .
```

If you are contributing to the official PipeCD repository, you can use the provided Makefile:

```bash
make build/plugin
```

## 2. Configure Piped to use your Plugin

To test your plugin locally, you need a `piped-config.yaml` that registers your plugin's address.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  plugins:
    - name: my-custom-plugin
      url: localhost:7001 # Address where your plugin will listen
```

## 3. Run Piped and your Plugin

First, start your plugin. By default, the SDK listens on a port (you can configure this via flags or environment variables, usually `--port`).

```bash
./my-plugin --port=7001
```

In another terminal, run `piped` with your local configuration:

```bash
# Using the make command from the pipecd repo
make run/piped CONFIG_FILE=piped-config.yaml EXPERIMENTAL=true INSECURE=true
```

## 4. Verify in the Web UI

1. Create an application in PipeCD that uses your custom stage.
2. Trigger a deployment.
3. Open the deployment details in the Web UI.
4. You should see your custom stage executing. Click on it to see the logs sent via `LogPersister`.

## 5. Debugging Tips

- **Check Piped Logs**: `piped` will log any gRPC connection errors or handshake failures.
- **Verbose Logging**: Use `lp.Debugf` in your plugin for detailed logs that only appear if the deployment is in debug mode.
- **Restarting**: If you change your plugin code, you must rebuild and restart the plugin binary. `piped` will automatically try to reconnect.

---

Congratulations! You've built and tested your first PipeCD plugin. For more complex examples, explore the [official plugins](https://github.com/pipe-cd/pipecd/tree/master/pkg/app/pipedv1/plugin) in the repository.
