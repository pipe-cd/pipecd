---
title: "Command-line tool: pipectl"
linkTitle: "Command-line tool: pipectl"
weight: 9
description: >
  This page describes how to install and use pipectl to manage PipeCD's resources.
---

Besides using web UI, PipeCD also provides a command-line tool, pipectl, which allows you to run commands against your project's resources.
You can use pipectl to add and sync applications, wait for a deployment status.

## Installation

### Binary

1. Download the appropriate version for your platform from [PipeCD Releases](https://github.com/pipe-cd/pipecd/releases).

    We recommend using the latest version of pipectl to avoid unforeseen issues.
    Run the following script:

    ``` console
    # OS="darwin" or "linux"
    curl -Lo ./pipectl https://github.com/pipe-cd/pipecd/releases/download/{{< blocks/latest_version >}}/pipectl_{{< blocks/latest_version >}}_{OS}_amd64
    ```

2. Make the pipectl binary executable.

    ``` console
    chmod +x ./pipectl
    ```

3. Move the binary to your PATH.

    ``` console
    sudo mv ./pipectl /usr/local/bin/pipectl
    ```

4. Test to ensure the version you installed is up-to-date.

    ``` console
    pipectl version
    ```

### [Asdf](https://asdf-vm.com/)

1. Add pipectl plugin to asdf. (If you have not yet `asdf add plugin add pipectl`.)
    ```console
    asdf plugin add pipectl
    ```

2. Install pipectl. Available versions are [here](https://github.com/pipe-cd/pipecd/releases).
    ```console
    asdf install pipectl {VERSION}
    ```

3. Set a version.
    ```console
    asdf global pipectl {VERSION}
    ```

4. Test to ensure the version you installed is up-to-date.

    ``` console
    pipectl version
    ```

### Docker
We are storing every version of docker image for pipectl on Google Cloud Container Registry.
Available versions are [here](https://github.com/pipe-cd/pipecd/releases).

```
docker run --rm gcr.io/pipecd/pipectl:{VERSION} -h
```

## Authentication

In order for pipectl to authenticate with PipeCD's Control Plane, it needs an API key, which can be created from `Settings/API Key` tab on the web UI.
There are two kinds of key role: `READ_ONLY` and `READ_WRITE`. Depending on the command, it might require an appropriate role to execute.

![](/images/settings-api-key.png)
<p style="text-align: center;">
Adding a new API key from Settings tab
</p>

When executing a command of pipectl you have to specify either a string of API key via `--api-key` flag or a path to the API key file via `--api-key-file` flag. 

## Usage

### Help

Run `help` to know the available commands:

``` console
$ pipectl --help

The command line tool for PipeCD.

Usage:
  pipectl [command]

Available Commands:
  application  Manage application resources.
  deployment   Manage deployment resources.
  encrypt      Encrypt the plaintext entered in either stdin or the --input-file flag.
  event        Manage event resources.
  help         Help about any command
  init         Generate an application config (app.pipecd.yaml) easily and interactively.
  piped        Manage piped resources.
  plan-preview Show plan preview against the specified commit.
  quickstart   Quick prepare PipeCD control plane in quickstart mode.
  version      Print the information of current binary.

Flags:
  -h, --help                               help for pipectl
      --log-encoding string                The encoding type for logger [json|console|humanize]. (default "humanize")
      --log-level string                   The minimum enabled logging level. (default "info")
      --metrics                            Whether metrics is enabled or not. (default true)
      --profile                            If true enables uploading the profiles to Stackdriver.
      --profile-debug-logging              If true enables logging debug information of profiler.
      --profiler-credentials-file string   The path to the credentials file using while sending profiles to Stackdriver.

Use "pipectl [command] --help" for more information about a command.
```

### Adding a new application

Add a new application into the project:

``` console
pipectl application add \
    --address=CONTROL_PLANE_API_ADDRESS \
    --api-key=API_KEY \
    --app-name=simple \
    --app-kind=KUBERNETES \
    --piped-id=PIPED_ID \
    --platform-provider=kubernetes-default \
    --repo-id=examples \
    --app-dir=kubernetes/simple
```

Run `help` to know what command flags should be specified:

``` console
$ pipectl application add --help

Add a new application.

Usage:
  pipectl application add [flags]

Flags:
      --app-dir string            The relative path from the root of repository to the application directory.
      --app-kind string           The kind of application. (KUBERNETES|TERRAFORM|LAMBDA|CLOUDRUN)
      --app-name string           The application name.
      --platform-provider string  The platform provider name. One of the registered providers in the piped configuration. The previous name of this field is cloud-provider.
      --config-file-name string   The configuration file name. (default "app.pipecd.yaml")
      --description string        The description of the application.
  -h, --help                      help for add
      --piped-id string           The ID of piped that should handle this application.
      --repo-id string            The repository ID. One the registered repositories in the piped configuration.

Global Flags:
      --address string                     The address to Control Plane api.
      --api-key string                     The API key used while authenticating with Control Plane.
      --api-key-file string                Path to the file containing API key used while authenticating with Control Plane.
      --cert-file string                   The path to the TLS certificate file.
      --insecure                           Whether disabling transport security while connecting to Control Plane.
      --log-encoding string                The encoding type for logger [json|console|humanize]. (default "humanize")
      --log-level string                   The minimum enabled logging level. (default "info")
      --metrics                            Whether metrics is enabled or not. (default true)
      --profile                            If true enables uploading the profiles to Stackdriver.
      --profile-debug-logging              If true enables logging debug information of profiler.
      --profiler-credentials-file string   The path to the credentials file using while sending profiles to Stackdriver.
```

### Syncing an application

- Send a request to sync an application and exit immediately when the deployment is triggered:

  ``` console
  pipectl application sync \
      --address={CONTROL_PLANE_API_ADDRESS} \
      --api-key={API_KEY} \
      --app-id={APPLICATION_ID}
  ```

- Send a request to sync an application and wait until the triggered deployment reaches one of the specified statuses:

  ``` console
  pipectl application sync \
      --address={CONTROL_PLANE_API_ADDRESS} \
      --api-key={API_KEY} \
      --app-id={APPLICATION_ID} \
      --wait-status=DEPLOYMENT_SUCCESS,DEPLOYMENT_FAILURE
  ```

### Getting an application

Display the information of a given application in JSON format:

``` console
pipectl application get \
    --address={CONTROL_PLANE_API_ADDRESS} \
    --api-key={API_KEY} \
    --app-id={APPLICATION_ID}
```

### Listing applications

Find and display the information of matching applications in JSON format:

``` console
pipectl application list \
    --address={CONTROL_PLANE_API_ADDRESS} \
    --api-key={API_KEY} \
    --app-name={APPLICATION_NAME} \
    --app-kind=KUBERNETES \
```

### Disable an application

Disable an application with given id:

``` console
pipectl application disable \
    --address={CONTROL_PLANE_API_ADDRESS} \
    --api-key={API_KEY} \
    --app-id={APPLICATION_ID}
```

### Deleting an application

Delete an application with given id:

``` console
pipectl application delete \
    --address={CONTROL_PLANE_API_ADDRESS} \
    --api-key={API_KEY} \
    --app-id={APPLICATION_ID}
```

### List deployments

Show the list of deployments based on filters.

```console
pipectl deployment list \
    --address={CONTROL_PLANE_API_ADDRESS} \
    --api-key={API_KEY} \
    --app-id={APPLICATION_ID}
```

### Waiting a deployment status

Wait until a given deployment reaches one of the specified statuses:

``` console
pipectl deployment wait-status \
    --address={CONTROL_PLANE_API_ADDRESS} \
    --api-key={API_KEY} \
    --deployment-id={DEPLOYMENT_ID} \
    --status=DEPLOYMENT_SUCCESS
```

### Get deployment stages log

Get deployment stages log.

```console
pipectl deployment logs \
    --address={CONTROL_PLANE_API_ADDRESS} \
    --api-key={API_KEY} \
    --deployment-id={DEPLOYMENT_ID}
```

### Registering an event for EventWatcher

Register an event that can be used by EventWatcher:

``` console
pipectl event register \
    --address={CONTROL_PLANE_API_ADDRESS} \
    --api-key={API_KEY} \
    --name=example-image-pushed \
    --data=gcr.io/pipecd/example:v0.1.0
```

### Encrypting the data you want to use when deploying

Encrypt the plaintext entered either in stdin or via the `--input-file` flag.

You can encrypt it the same way you do [from the web](../managing-application/secret-management/#encrypting-secret-data).

- From stdin:

  ``` console
  pipectl encrypt \
      --address={CONTROL_PLANE_API_ADDRESS} \
      --api-key={API_KEY} \
      --piped-id={PIPED_ID} <{PATH_TO_SECRET_FILE}
  ```

- From the `--input-file` flag:

  ``` console
  pipectl encrypt \
      --address={CONTROL_PLANE_API_ADDRESS} \
      --api-key={API_KEY} \
      --piped-id={PIPED_ID} \
      --input-file={PATH_TO_SECRET_FILE}
  ```

Note: The docs for pipectl available command is maybe outdated, we suggest users use the `help` command for the updated usage while using pipectl.

### Generating an application config (app.pipecd.yaml)

Generate an app.pipecd.yaml interactively:

``` console
$ pipectl init 
Which platform? Enter the number [0]Kubernetes [1]ECS: 1
Name of the application: myApp
...
```

After the above interaction, you can get the config YAML:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: ECSApp
spec:
  name: myApp
  input:
    serviceDefinitionFile: serviceDef.yaml
    taskDefinitionFile: taskDef.yaml
    targetGroups:
      primary:
        targetGroupArn: arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/xxx/xxx
        containerName: web
        containerPort: 80
  description: Generated by `pipectl init`. See https://pipecd.dev/docs/user-guide/configuration-reference/ for more.
```


### You want more?

We always want to add more needed commands into pipectl. Please let us know what command you want to add by creating issues in the [pipe-cd/pipecd](https://github.com/pipe-cd/pipecd/issues) repository. We also welcome your pull request to add the command.
