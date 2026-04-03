---
title: "Installing on Cloud Run"
linkTitle: "Installing on Cloud Run"
weight: 3
description: >
  This page describes how to install `piped` on Cloud Run.
---

## Prerequisites

### A registered `piped`

- Make sure your `piped` is registered in the Control Plane and that you have its **PIPED_ID** and **PIPED_KEY**.  
- If not, follow the guide to [register a new `piped`](../../../user-guide/managing-controlplane/registering-a-piped/).

### SSH key for Git repositories

- If your Git repositories are private, `piped` requires a private SSH key to access those repositories.
- Please check out [this documentation](https://help.github.com/en/github/authenticating-to-github/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent) for how to generate a new SSH key pair. Then add the public key to your repositories. (If you are using GitHub, you can add it to Deploy Keys at the repository’s **Settings** page.)

## Installation

### Preparing the `piped` configuration file

Plugins are external binaries that have to be referenced in the `piped` configuration file. There are no plugins set by default. In PipeCD v1, deployment targets are configured under `plugins`, not legacy `platformProviders`. See [Configuring a plugin](../../../user-guide/managing-piped/configuring-a-plugin/) for how to add Kubernetes, Terraform, or other plugins.

An example of a minimal `piped` v1 configuration file:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: {PROJECT_ID}
  pipedID: {PIPED_ID}
  pipedKeyData: {BASE64_ENCODED_PIPED_KEY}
  # Write in a format like "host:443" because the communication is done via gRPC.
  apiAddress: {CONTROL_PLANE_API_ADDRESS}

  git:
    sshKeyData: {BASE64_ENCODED_PRIVATE_SSH_KEY}

  repositories:
    - repoId: {REPO_ID_OR_NAME}
      remote: git@github.com:{GIT_ORG}/{GIT_REPO}.git
      branch: {GIT_BRANCH}
  syncInterval: 1m
  plugins: {}

  # Optional
  # Uncomment this if you want to enable Secret Management.
  # See: https://pipecd.dev/docs/user-guide/managing-application/secret-management/
  # secretManagement:
  #   type: KEY_PAIR
  #   config:
  #     privateKeyData: {BASE64_ENCODED_PRIVATE_KEY}
  #     publicKeyData: {BASE64_ENCODED_PUBLIC_KEY}
```

>**Note:**
> `piped`'s plugins are versioned independently from PipeCD. See the [latest releases](https://github.com/pipe-cd/pipecd/releases) for more information.
>
> There is also a [Community plugins repository on GitHub](https://github.com/pipe-cd/community-plugins) for plugins built by the community.

See the [configuration reference](../../../user-guide/managing-piped/configuration-reference/) for the full configuration.

### Storing the configuration in Secret Manager

Create a new secret in [Secret Manager](https://cloud.google.com/secret-manager/docs/creating-and-accessing-secrets) to store the configuration securely:

```console
gcloud secrets create cloudrun-piped-config --data-file={PATH_TO_CONFIG_FILE}
```

Then make sure that Cloud Run has permission to access that secret as described in the [Cloud Run secret access guide](https://cloud.google.com/run/docs/configuring/secrets#access-secret).

### Running `piped` on Cloud Run

Prepare a Cloud Run service manifest as below.

{{< tabpane >}}
{{< tab lang="yaml" header="Piped with Remote-upgrade" >}}
# Enable remote-upgrade feature of piped.
# https://pipecd.dev/docs/user-guide/managing-piped/remote-upgrade-remote-config/#remote-upgrade
# This allows upgrading piped to a new version from the web console.

apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: piped
  annotations:
    run.googleapis.com/ingress: internal
    run.googleapis.com/ingress-status: internal
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/maxScale: "1"          # This must be 1.
        autoscaling.knative.dev/minScale: "1"          # This must be 1.
        run.googleapis.com/cpu-throttling: "false"     # This is required.
    spec:
      containerConcurrency: 1                          # This must be 1 to ensure piped works correctly.
      containers:
        - image: gcr.io/pipecd/launcher:{{< blocks/latest_version >}}
          args:
            - launcher
            - --launcher-admin-port=9086
            - --config-file=/etc/piped-config/config.yaml
          ports:
            - containerPort: 9086
          volumeMounts:
            - mountPath: /etc/piped-config
              name: piped-config
          resources:
            limits:
              cpu: 1000m
              memory: 2Gi
      volumes:
        - name: piped-config
          secret:
            secretName: cloudrun-piped-config
            items:
              - key: latest
                path: config.yaml
{{< /tab >}}

{{< tab lang="yaml" header="Piped without Remote-upgrade" >}}
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: piped
  annotations:
    run.googleapis.com/ingress: internal
    run.googleapis.com/ingress-status: internal
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/maxScale: "1"
        autoscaling.knative.dev/minScale: "1"
        run.googleapis.com/cpu-throttling: "false"
    spec:
      containerConcurrency: 1
      containers:
        - image: gcr.io/pipecd/launcher:{{< blocks/latest_version >}}
          args:
            - launcher
            - --launcher-admin-port=9086
            - --config-file=/etc/piped-config/config.yaml
          ports:
            - containerPort: 9086
          volumeMounts:
            - mountPath: /etc/piped-config
              name: piped-config
          resources:
            limits:
              cpu: 1000m
              memory: 2Gi
      volumes:
        - name: piped-config
          secret:
            secretName: cloudrun-piped-config
            items:
              - key: latest
                path: config.yaml
{{< /tab >}}
{{< /tabpane >}}

Apply the Cloud Run service:

```console
gcloud run services replace {PATH_TO_CLOUD_RUN_SERVICE_MANIFEST}
```

Once the service is created, Cloud Run will run the `piped` agent as a stateless service that connects to your PipeCD Control Plane and deploys applications according to your configuration. Make sure that the created secret is accessible from this `piped` service as described in the [Cloud Run secret access guide](https://cloud.google.com/run/docs/configuring/secrets#access-secret).


