---
title: "Running Locally with Piped"
weight: 21
description: >
  Step-by-step walkthrough to run, test, and deploy using the custom plugin.
---

Now that the custom plugin is fully built, let's run it locally and test it against a real PipeCD deployment agent!

---

### 1. Setting Up the Control Plane

To deploy applications, we need a running PipeCD Control Plane and a Piped agent. While production instances are typically deployed on Kubernetes, we will use Docker Compose for our local environment.

Clone the official `pipe-cd/tutorial` repository:

```console
$ git clone https://github.com/pipe-cd/tutorial.git
$ cd tutorial/src/install/control-plane
```

Open `docker-compose.yaml` and update the `ghcr.io/pipe-cd/pipecd` image tag to `v0.52.0-54-g8a12400` (the version compiled with Pluggable Architecture support).

Start the control plane:

```console
$ docker compose up -d
```

Open [http://localhost:8080/login](http://localhost:8080/login) in your browser. Log in with the default credentials:

- **Project Name**: `tutorial`
- **Username**: `hello-pipecd`
- **Password**: `hello-pipecd`

---

### 2. Registering Piped

1. Open [http://localhost:8080/settings/piped?project=tutorial](http://localhost:8080/settings/piped?project=tutorial).
2. Click **+ ADD** in the top left corner.
3. Fill in a descriptive Name and Description, and click **SAVE**.
4. Securely copy the generated **Piped Id** and **Base64 Encoded Piped Key**.

---

### 3. Setting Up the Git Manifest Repository

Piped continuously syncs state with a Git repository. We will set up a local bare Git repository so we do not have to push to GitHub during local testing.

In a new terminal window:

```console
$ git init --bare pipecd-manifest.git
$ git clone ./pipecd-manifest.git
$ cd pipecd-manifest
$ mkdir demo-file-app
```

Under `demo-file-app`, create an `app.pipecd.yaml` configuration file specifying our `file` plugin:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: demo-file-app
  description: |
    Demo application synchronized using our custom file plugin.
  pipeline:
    stages:
      - name: FILE_DIFF
      - name: FILE_SYNC
  plugins:
    file:
      path: /tmp/try-pipecd-file-plugin
```

Commit and push these changes:

```console
$ git add demo-file-app/app.pipecd.yaml
$ git commit -m "Initialize demo file application"
$ git branch -M main
$ git push origin main
```

---

### 4. Compiling the Plugin Binary

Compile the plugin in your plugin development directory:

```console
$ go build -o pipecd-plugin-file
```

---

### 5. Configuring Piped

Create Piped's local configuration file (`piped.config.yaml`). Ensure the directory paths are absolute:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: tutorial
  pipedID: <<YOUR_PIPED_ID>>
  pipedKeyData: <<YOUR_PIPED_KEY_DATA>>
  apiAddress: localhost:8080
  repositories:
    - repoId: local-manifest
      remote: file://<<ABSOLUTE_PATH_TO_pipecd-manifest.git>>
      branch: main
  plugins:
    - name: file
      port: 7001
      url: file://<<ABSOLUTE_PATH_TO_pipecd-plugin-file_binary>>
      deployTargets:
        - name: local
          config: {}
```

---

### 6. Starting Piped

Download the Pluggable Architecture Piped binary released under [kubecon-jp-2025](https://github.com/pipe-cd/pipecd/releases/tag/kubecon-jp-2025) matching your OS and CPU architecture.

Make it executable and launch Piped:

```console
$ chmod +x piped_kubecon_jp_2025_${os}_${arch}
$ ./piped_kubecon_jp_2025_${os}_${arch} piped --config-file=./piped.config.yaml --insecure
```

Check the logs to verify that Piped successfully connects to the Control Plane and registers the `file` plugin via localhost port `7001`.

---

### 7. Registering the Application

1. Open [http://localhost:8080/applications?project=tutorial](http://localhost:8080/applications?project=tutorial).
2. Click **+ ADD** and navigate to the **ADD FROM SUGGESTIONS** tab.
3. Select your Piped and the suggested application, then click **SAVE**.

---

### 8. Testing Synchronization

Since no sync files exist in our Git repository yet, the first deployment will create an empty `/tmp/try-pipecd-file-plugin` folder.

Add a file to the Git repository to test syncing:

```console
$ cd pipecd-manifest/demo-file-app
$ echo "Hello from PipeCD Custom Plugin!" > hello.txt
$ git add hello.txt
$ git commit -m "Add hello.txt"
$ git push origin main
```

Piped will automatically detect the commit, trigger a deployment, and run the `FILE_DIFF` and `FILE_SYNC` stages.

Verify that the file is synced successfully:

```console
$ cat /tmp/try-pipecd-file-plugin/hello.txt
Hello from PipeCD Custom Plugin!
```
