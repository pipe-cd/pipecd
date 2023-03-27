---
title: "Installing on Cloud Run"
linkTitle: "Installing on Cloud Run"
weight: 3
description: >
  This page describes how to install Piped on Cloud Run.
---

## Prerequisites

##### Having piped's ID and Key strings
- Ensure that the `piped` has been registered and you are having its PIPED_ID and PIPED_KEY strings.
- If you are not having them, this [page](../../../user-guide/managing-controlplane/registering-a-piped/) guides you how to register a new one.

##### Preparing SSH key
- If your Git repositories are private, `piped` requires a private SSH key to access those repositories.
- Please checkout [this documentation](https://help.github.com/en/github/authenticating-to-github/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent) for how to generate a new SSH key pair. Then add the public key to your repositories. (If you are using GitHub, you can add it to Deploy Keys at the repository's Settings page.)

## Installation

- Preparing a piped configuration file as the following:

  ``` yaml
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

    # Optional
    # Enable this Piped to handle Cloud Run application.
    platformProviders:
      - name: cloudrun-in-project
        type: CLOUDRUN
        config:
          project: {GCP_PROJECT_ID}
          region: {GCP_PROJECT_REGION}

    # Optional
    # Uncomment this if you want to enable this Piped to handle Terraform application.
    #  - name: terraform-gcp
    #    type: TERRAFORM

    # Optional
    # Uncomment this if you want to enable SecretManagement feature.
    # https://pipecd.dev//docs/user-guide/managing-application/secret-management/
    # secretManagement:
    #   type: KEY_PAIR
    #   config:
    #     privateKeyData: {BASE64_ENCODED_PRIVATE_KEY}
    #     publicKeyData: {BASE64_ENCODED_PUBLIC_KEY}
  ```

See [ConfigurationReference](../../../user-guide/managing-piped/configuration-reference/) for the full configuration.

- Creating a new secret in [SecretManager](https://cloud.google.com/secret-manager/docs/creating-and-accessing-secrets) to store above configuration data securely

  ``` console
  gcloud secrets create cloudrun-piped-config --data-file={PATH_TO_CONFIG_FILE}
  ```

  then make sure that Cloud Run has the ability to access that secret as [this guide](https://cloud.google.com/run/docs/configuring/secrets#access-secret).

- Running Piped in Cloud Run

  Prepare a Cloud Run service manifest file as below.

  {{< tabpane >}}
  {{< tab lang="yaml" header="Piped with Remote-upgrade" >}}
# Enable remote-upgrade feature of Piped.
# https://pipecd.dev/docs/user-guide/managing-piped/remote-upgrade-remote-config/#remote-upgrade
# This allows upgrading Piped to a new version from the web console.

apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: piped
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/maxScale: '1'           # This must be 1.
        autoscaling.knative.dev/minScale: '1'           # This must be 1.
        run.googleapis.com/ingress: internal
        run.googleapis.com/ingress-status: internal
        run.googleapis.com/cpu-throttling: "false"      # This is required.
    spec:
      containerConcurrency: 1                           # This must be 1 to ensure Piped work correctly.
      containers:
        - image: ghcr.io/pipe-cd/launcher:{{< blocks/latest_version >}}
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
              - path: config.yaml
                key: latest
  {{< /tab >}}
  {{< tab lang="yaml" header="Piped" >}}
# This just installs a Piped with the specified version.
# Whenever you want to upgrade that Piped to a new version or update its config data you have to restart it.

apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: piped
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/maxScale: '1'           # This must be 1.
        autoscaling.knative.dev/minScale: '1'           # This must be 1.
        run.googleapis.com/ingress: internal
        run.googleapis.com/ingress-status: internal
        run.googleapis.com/cpu-throttling: "false"      # This is required.
    spec:
      containerConcurrency: 1                           # This must be 1.
      containers:
        - image: ghcr.io/pipe-cd/piped:{{< blocks/latest_version >}}
          args:
            - piped
            - --config-file=/etc/piped-config/config.yaml
          ports:
            - containerPort: 9085
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
              - path: config.yaml
                key: latest
  {{< /tab >}}
  {{< /tabpane >}}

  Run Piped service on Cloud Run with the following command:

  ``` console
  gcloud beta run services replace cloudrun-piped-service.yaml
  ```

  Note: Make sure that the created secret is accessible from this Piped service. See more [here](https://cloud.google.com/run/docs/configuring/secrets#access-secret).
