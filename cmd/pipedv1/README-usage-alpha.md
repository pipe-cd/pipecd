# Usage of pipedv1 and plugins (alpha)

_This page is still in preparation. The content will be changed and might not work well yet._

This page shows how to run pipedv1 and plugins of alpha status.

See [cmd/pipedv1/README.md](https://github.com/pipe-cd/pipecd/blob/master/cmd/pipedv1/README.md) if you want to develop or debug pipedv1.

_This page might be moved to another place in the future._


## Prerequisites

- kubectl and a k8s cluster (They are not required if you won't use the kubernetes plugin)

## 1. Setup Control Plane

1. Run a Control Plane that your piped will connect to. If you want to run a Control Plane locally, see [How to run Control Plane locally](https://github.com/pipe-cd/pipecd/blob/master/cmd/pipecd/README.md#how-to-run-control-plane-locally).
    - The Control Plane version must be v0.52.0 or later.

2. Generate a new piped key/ID.
    2.1. Access the Control Plane console.
    2.2. Go to the piped list page. (https://{console-address}/settings/piped)
    2.4. Add a new piped via the `+ADD` button.
    2.5. Copy the generated piped ID and base64 encoded key.

## 2. Run pipedv1

1. Create a piped config file like the following.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  apiAddress: {CONTROL_PLANE_API_ADDRESS} # like "localhost:443"
  projectID: {PROJECT_ID}
  pipedID: {PIPED_ID}
  pipedKeyData: {BASE64_ENCODED_PIPED_KEY} # or use pipedKeyFile
  repositories:
    - repoID: repo1
      remote: https://github.com/your-account/your-repo
      branch: xxx
  # See https://pipecd.dev/docs/user-guide/managing-piped/configuration-reference/ for details of above.
  # platformProviders is not necessary.

  plugins:
    - name: kubernetes
      port: 7001 # Any unused port
      url: https://github.com/pipe-cd/pipecd/xxxxxxxxxxx # TODO: Ref to the Release 
      deployTargets:
        - name: cluster1
          config: 
            masterURL: https://127.0.0.1:61337   # shown by kubectl cluster-info
            kubeConfigPath: /path/to/kubeconfig
            kubectlVersion: 1.33.0
    - name: wait
      port: 7002 # Any unused port
      url: https://github.com/pipe-cd/pipecd/xxxxxxxxxxx # TODO: Ref to the Release 

    - name: example-stage
      port: 7003 # Any unused port
      url: https://github.com/pipe-cd/community-plugins/xxxxxxxxxxx # TODO: Ref to the Release 
      config:
        - commonMessage: "[common message]"
```

2. Run pipedv1

```sh
pipedv1 piped --config-file=/path/to/piped-config.yaml --tools-dir=/tmp/piped-bin
```

- The pipedv1 version must be v0.52.0 or later.
- If your Control Plane runs on local, add `INSECURE=true` to the command to skip TLS certificate checks.


## 3. Deploy an application

1. Create an app.pipecd.yaml like the following.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: canary
  labels:
    env: example
    team: product
  pipeline:
    stages:
      # Deploy the workloads of CANARY variant. In this case, the number of
      # workload replicas of CANARY variant is 10% of the replicas number of PRIMARY variant.
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 10%
      # Wait 10 seconds before going to the next stage.
      - name: WAIT
        with:
          duration: 10s
      # Update the workload of PRIMARY variant to the new version.
      - name: K8S_PRIMARY_ROLLOUT
      # Destroy all workloads of CANARY variant.
      - name: K8S_CANARY_CLEAN
```

2. Push the app.pipecd.yaml to your remote repository.
3. On the Control Plane console, register the application via `PIPED V1 ADD FROM SUGGESTIONS` tab.

## See also

<!-- TODO: Link to each config reference -->
- kubernetes plugin: [README.md](/pkg/app/pipedv1/plugin/kubernetes/README.md)
- wait stage plugin: [README.md](/pkg/app/pipedv1/plugin/wait/README.md)
- example-stage plugin: TBA

## Note

- Currently, on the Deployment detail UI, each stage is not visible until it starts.
  - The cause is that the `Visible` field of each stage is set to `false` by default in pipedv1. Instead, the new field `Rollback` is used to determine if the stage is a rollback stage.
  - This will be modified before releasing pipedv1 beta.
