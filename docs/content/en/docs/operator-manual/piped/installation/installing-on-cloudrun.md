---
title: "Installing on CloudRun"
linkTitle: "Installing on CloudRun"
weight: 3
description: >
  This page describes how to install Piped on CloudRun.
---

## Prerequisites

##### Having piped's ID and Key strings
- Ensure that the `piped` has been registered and you are having its PIPED_ID and PIPED_KEY strings.
- If you are not having them, this [page](/docs/operator-manual/control-plane/registering-a-piped/) guides you how to register a new one.

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
    webAddress: {CONTROL_PLANE_WEB_ADDRESS}
    # Write in a format like "host:443" because the communication is done via gRPC.
    apiAddress: {CONTROL_PLANE_API_ADDRESS}

    git:
      sshKeyData: {BASE64_ENCODED_PRIVATE_SSH_KEY}

    repositories:
      - repoId: {REPO_ID_OR_NAME}
        remote: git@github.com:{GIT_ORG}/{GIT_REPO}.git
        branch: {GIT_BRANCH}

    # Optional
    # Enable this Piped to handle CLOUD_RUN application.
    cloudProviders:
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
    # https://pipecd.dev/docs/user-guide/secret-management
    # secretManagement:
    #   type: KEY_PAIR
    #   config:
    #     privateKeyData: {BASE64_ENCODED_PRIVATE_KEY}
    #     publicKeyData: {BASE64_ENCODED_PUBLIC_KEY}
  ```

- Creating a new secret in [SecretManager](https://cloud.google.com/secret-manager/docs/creating-and-accessing-secrets) to store above configuration data securely

  ``` console
  gcloud secrets create cloudrun-piped-config --data-file={PATH_TO_CONFIG_FILE}
  ```

- Running Piped in CloudRun

  Prepare a CloudRun service manifest file as below.

  **Note**: Fields which set to `1` are strict to be set with that value to ensure piped work correctly.

  ``` yaml
  apiVersion: serving.knative.dev/v1
  kind: Service
  metadata:
    name: piped
  spec:
    template:
      metadata:
        annotations:
          autoscaling.knative.dev/maxScale: '1' # This must be 1.
          autoscaling.knative.dev/minScale: '1' # This must be 1.
          run.googleapis.com/ingress: internal
          run.googleapis.com/ingress-status: internal
      spec:
        containerConcurrency: 1 # This must be 1.
        containers:
          - image: gcr.io/pipecd/piped:{{< blocks/latest_version >}}
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
                memory: 512Mi
        volumes:
          - name: piped-config
            secret:
              secretName: cloudrun-piped-config
              items:
                - path: config.yaml
                  key: latest
  ```

  Create Piped service on CloudRun with the following command. Please note to use `no-cpu-throttling` flag to disable CPU throttling on its container.

  ``` console
  gcloud beta run services replace cloudrun-piped-service.yaml --no-cpu-throttling
  ```

  Note: Make sure that the created secret is accessible from this Piped service. See more [here](https://cloud.google.com/run/docs/configuring/secrets#access-secret).
