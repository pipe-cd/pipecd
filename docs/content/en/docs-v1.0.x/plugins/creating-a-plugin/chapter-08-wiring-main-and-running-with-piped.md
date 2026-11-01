---
title: "Wiring main.go and running with Piped"
linkTitle: "Wiring main.go and running with Piped"
weight: 8
description: >
  Register the finished plugin in main.go and run it against a local piped.
---

The plugin is code-complete. This chapter registers it in the `main` function, builds it into a binary, and runs it against a local `piped` so you can watch a real deployment go through the stages you wrote.

## Register the plugin in main.go

You wrote `main.go` back when the plugin first became buildable, and it already does everything it needs to:

```go
package main

import (
	"log"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

func main() {
	p, err := sdk.NewPlugin(
		"0.0.1",
		sdk.WithDeploymentPlugin(&plugin{}),
	)
	if err != nil {
		log.Fatalln(err)
	}

	if err := p.Run(); err != nil {
		log.Fatalln(err)
	}
}
```

`sdk.NewPlugin` takes the plugin's version and registers the deployment plugin with `sdk.WithDeploymentPlugin`. Because you passed `&plugin{}` in from the start, and it now implements all six methods, there is nothing left to change here. `p.Run` starts the gRPC server that `piped` connects to and calls those methods on.

## Build the plugin binary

`piped` runs each plugin as a separate binary. Build one from the project directory:

```bash
go build -o pipecd-plugin-file
```

This produces a `pipecd-plugin-file` binary in the same directory. Note its absolute path; the `piped` configuration points at it.

## Prepare a control plane and a piped

Running the plugin needs two PipeCD components: a Control Plane, which holds projects and shows deployment status in its web UI, and a `piped` agent, which loads your plugin and runs deployments.

Set up a local Control Plane by following the [Quickstart](/docs-v1.0.x/quickstart/). Once its web UI is up, register a new `piped` from the **Settings** page and note the **Piped ID** and the **Base64 encoded Piped key**; the `piped` configuration needs both.

You will run the `piped` agent yourself, using a v1 build. From a clone of the [`pipecd`](https://github.com/pipe-cd/pipecd) repository you can run `make run/piped EXPERIMENTAL=true` (the `EXPERIMENTAL=true` flag selects the v1 agent; plain `make run/piped` runs the v0 agent, which loads no plugins), or install a v1 release binary as described in [Installing on a single machine](/docs-v1.0.x/installation/install-piped/installing-on-single-machine/). Either way, the command form is:

```bash
piped run --config-file=PATH_TO_PIPED_CONFIG --insecure=true
```

`--insecure=true` is fine here because the local Control Plane has no TLS. Older blog posts and the original Japanese book use an older `piped` with different syntax; if a command does not match your build, run `piped --help` to confirm the current flags.

## Write the application configuration

`piped` reads application configuration from a Git repository. It can read a local repository directly, which is convenient for trying things out. Create a bare repository and a working clone:

```bash
git init --bare pipecd-manifest.git
git clone ./pipecd-manifest.git
cd pipecd-manifest
mkdir -p demo-file-app
```

In `demo-file-app`, create `app.pipecd.yaml`. This tells PipeCD the application uses the `file` plugin and syncs its files to `/tmp/try-pipecd-file-plugin`:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: demo-file-app
  description: Application for trying out the file plugin.
  pipeline:
    stages:
      - name: FILE_DIFF
      - name: FILE_SYNC
  plugins:
    file:
      path: /tmp/try-pipecd-file-plugin
```

The `plugins.file.path` value is the `path` field of the `applicationConfig` type you defined in an earlier chapter. The plugin lists this directory during its diff and sync stages, so create it before the first deployment:

```bash
mkdir -p /tmp/try-pipecd-file-plugin
```

Commit the file and push it:

```bash
git add demo-file-app/app.pipecd.yaml
git commit -m "Add demo-file-app"
git push origin main
```

## Write the piped configuration

Now create the `piped` configuration file. Fill in the Piped ID and key you noted earlier. Two paths must be absolute: the `remote` of the repository, which points at the bare repository you created, and the plugin `url`, which points at the binary you built.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: quickstart
  pipedID: {PIPED_ID}
  pipedKeyData: {BASE64_ENCODED_PIPED_KEY}
  apiAddress: localhost:8080
  repositories:
    - repoId: local-manifest
      remote: file:///absolute/path/to/pipecd-manifest.git
      branch: main
  plugins:
    - name: file
      port: 7001
      url: file:///absolute/path/to/pipecd-plugin-file
      deployTargets:
        - name: local
          config: {}
```

The `name` under `plugins` must match the plugin key you used in `app.pipecd.yaml`, which is `file`. The `port` is any free port the plugin listens on; `piped` talks to the plugin over it.

## Run piped and confirm the deployment

Start `piped` with the configuration file:

```bash
piped run --config-file=PATH_TO_PIPED_CONFIG --insecure=true
```

`piped` caches plugin binaries, so while you are developing and rebuilding the plugin, add `--force-plugin-redownload` to make `piped` pick up the new binary instead of a cached copy.

If the command keeps running and streams logs instead of exiting, `piped` has connected to the Control Plane. Now register the application in the web UI: open the **Applications** page, choose **ADD**, and use **ADD FROM SUGGESTIONS**, which lists the application `piped` found in your repository. Select it and save.

Registering an application triggers a first deployment. Open the **Deployments** page to watch it run through the `FILE_DIFF` and `FILE_SYNC` stages, using the log messages you wrote with the log persister. The repository has no files to sync yet, so the sync stage copies nothing.

Add a file and push it:

```bash
cd demo-file-app
echo a > a.txt
git add a.txt
git commit -m "Add a.txt"
git push origin main
```

`piped` polls the repository, notices the change, and starts another deployment. After a short wait, refresh the **Deployments** page. The new deployment runs, and `a.txt` appears in `/tmp/try-pipecd-file-plugin`. Your plugin is deploying real files.

The plugin now works end-to-end. The final chapter covers what comes next: testing your plugin more thoroughly, publishing it to the community plugins repository, and contributing it back.
