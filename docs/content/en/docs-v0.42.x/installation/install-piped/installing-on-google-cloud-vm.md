---
title: "Installing on Google Cloud VM"
linkTitle: "Installing on Google Cloud VM"
weight: 2
description: >
  This page describes how to install Piped on Google Cloud VM.
---

## Prerequisites

##### Having piped's ID and Key strings
- Ensure that the `piped` has been registered and you are having its `PIPED_ID` and `PIPED_KEY` strings.
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
    # Uncomment this if you want to enable this Piped to handle Cloud Run application.
    # platformProviders:
    #  - name: cloudrun-in-project
    #    type: CLOUDRUN
    #    config:
    #      project: {GCP_PROJECT_ID}
    #      region: {GCP_PROJECT_REGION}

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

- Creating a new secret in [SecretManager](https://cloud.google.com/secret-manager/docs/creating-and-accessing-secrets) to store above configuration data securely

  ``` shell
  gcloud secrets create vm-piped-config --data-file={PATH_TO_CONFIG_FILE}
  ```

- Creating a new Service Account for Piped and giving it needed roles

  ``` shell
  gcloud iam service-accounts create vm-piped \
    --description="Using by Piped running on Google Cloud VM" \
    --display-name="vm-piped"

  # Allow Piped to access the created secret.
  gcloud secrets add-iam-policy-binding vm-piped-config \
    --member="serviceAccount:vm-piped@{GCP_PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/secretmanager.secretAccessor"

  # Allow Piped to write its log messages to Google Cloud Logging service.
  gcloud projects add-iam-policy-binding {GCP_PROJECT_ID} \
    --member="serviceAccount:vm-piped@{GCP_PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/logging.logWriter"

  # Optional
  # If you want to use this Piped to handle Cloud Run application
  # run the following command to give it the needed roles.
  # https://cloud.google.com/run/docs/reference/iam/roles#additional-configuration
  #
  # gcloud projects add-iam-policy-binding {GCP_PROJECT_ID} \
  #   --member="serviceAccount:vm-piped@{GCP_PROJECT_ID}.iam.gserviceaccount.com" \
  #   --role="roles/run.developer"
  #
  # gcloud iam service-accounts add-iam-policy-binding {GCP_PROJECT_NUMBER}-compute@developer.gserviceaccount.com \
  #   --member="serviceAccount:vm-piped@{GCP_PROJECT_ID}.iam.gserviceaccount.com" \
  #   --role="roles/iam.serviceAccountUser"
  ```

- Running Piped on a Google Cloud VM

  {{< tabpane >}}
  {{< tab lang="console" header="Piped with Remote-upgrade" >}}
# Enable remote-upgrade feature of Piped.
# https://pipecd.dev/docs/user-guide/managing-piped/remote-upgrade-remote-config/#remote-upgrade
# This allows upgrading Piped to a new version from the web console.

  gcloud compute instances create-with-container vm-piped \
    --container-image="ghcr.io/pipe-cd/launcher:{{< blocks/latest_version >}}" \
    --container-arg="launcher" \
    --container-arg="--config-from-gcp-secret=true" \
    --container-arg="--gcp-secret-id=projects/{GCP_PROJECT_ID}/secrets/vm-piped-config/versions/{SECRET_VERSION}" \
    --network="{VPC_NETWORK}" \
    --subnet="{VPC_SUBNET}" \
    --scopes="cloud-platform" \
    --service-account="vm-piped@{GCP_PROJECT_ID}.iam.gserviceaccount.com"
  {{< /tab >}}
  {{< tab lang="console" header="Piped" >}}
# This just installs a Piped with the specified version.
# Whenever you want to upgrade that Piped to a new version or update its config data you have to restart it.

  gcloud compute instances create-with-container vm-piped \
    --container-image="ghcr.io/pipe-cd/piped:{{< blocks/latest_version >}}" \
    --container-arg="piped" \
    --container-arg="--config-gcp-secret=projects/{GCP_PROJECT_ID}/secrets/vm-piped-config/versions/{SECRET_VERSION}" \
    --network="{VPC_NETWORK}" \
    --subnet="{VPC_SUBNET}" \
    --scopes="cloud-platform" \
    --service-account="vm-piped@{GCP_PROJECT_ID}.iam.gserviceaccount.com"
  {{< /tab >}}
  {{< /tabpane >}}

After that, you can see on PipeCD web at `Settings` page that Piped is connecting to the Control Plane.
You can also view Piped log as described [here](https://cloud.google.com/compute/docs/containers/deploying-containers#viewing_logs).
