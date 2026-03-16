---
title: "Installing on Google Cloud VM"
linkTitle: "Installing on Google Cloud VM"
weight: 2
description: >
  This page describes how to install `piped` on a Google Cloud VM.
---

## Prerequisites

### A registered `piped`

- Make sure your `piped` is registered in the Control Plane and that you have its **PIPED_ID** and **PIPED_KEY**.  
- If not, follow the guide to [register a new `piped`](../../../user-guide/managing-controlplane/registering-a-piped/).

### SSH key for Git repositories

- If your Git repositories are private, `piped` requires a private SSH key to access those repositories.
- Please check out [this documentation](https://help.github.com/en/github/authenticating-to-github/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent) for how to generate a new SSH key pair. Then add the public key to your repositories. (If you are using GitHub, you can add it to Deploy Keys at the repository's **Settings** page.)

## Installation

### Preparing the `piped` configuration file

Prepare a `piped` configuration file as the following:

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

  # Optional
  # Uncomment this if you want to enable this piped to handle Cloud Run applications.
  # platformProviders:
  #  - name: cloudrun-in-project
  #    type: CLOUDRUN
  #    config:
  #      project: {GCP_PROJECT_ID}
  #      region: {GCP_PROJECT_REGION}

  # Optional
  # Uncomment this if you want to enable this piped to handle Terraform applications.
  #  - name: terraform-gcp
  #    type: TERRAFORM

  # Optional
  # Uncomment this if you want to enable Secret Management.
  # See: https://pipecd.dev/docs/user-guide/managing-application/secret-management/
  # secretManagement:
  #   type: KEY_PAIR
  #   config:
  #     privateKeyData: {BASE64_ENCODED_PRIVATE_KEY}
  #     publicKeyData: {BASE64_ENCODED_PUBLIC_KEY}
```

See the [configuration reference](../../../user-guide/managing-piped/configuration-reference/) for the full configuration.

### Creating a Secret Manager secret

Create a new secret in [Secret Manager](https://cloud.google.com/secret-manager/docs/creating-and-accessing-secrets) to store the configuration securely:

```shell
gcloud secrets create vm-piped-config --data-file={PATH_TO_CONFIG_FILE}
```

### Creating a service account and granting roles

Create a service account for `piped` and grant it the required roles:

```shell
gcloud iam service-accounts create vm-piped \
  --description="Used by piped running on Google Cloud VM" \
  --display-name="vm-piped"

# Allow piped to access the created secret.
gcloud secrets add-iam-policy-binding vm-piped-config \
  --member="serviceAccount:vm-piped@{GCP_PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"

# Allow piped to write its log messages to Google Cloud Logging.
gcloud projects add-iam-policy-binding {GCP_PROJECT_ID} \
  --member="serviceAccount:vm-piped@{GCP_PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/logging.logWriter"

# Optional
# If you want to use this piped to handle Cloud Run applications,
# grant additional roles as described in the Cloud Run IAM documentation:
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

### Running `piped` on a Google Cloud VM

Use `gcloud compute instances create-with-container` to run `piped` as a container on a VM. The following examples show how to configure the instance to read the configuration from Secret Manager.

{{< tabpane >}}
{{< tab lang="console" header="Piped with Remote-upgrade" >}}
# Enable remote-upgrade feature of piped.
# https://pipecd.dev/docs/user-guide/managing-piped/remote-upgrade-remote-config/#remote-upgrade
# This allows upgrading piped to a new version from the web console.

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
# This just installs a piped with the specified version.
# Whenever you want to upgrade that piped to a new version or update its config data you have to restart it.

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

After that, you should see on the PipeCD web `Settings` page that `piped` is connected to the Control Plane. You can also view `piped` logs as described in the [Compute Engine container logs documentation](https://cloud.google.com/compute/docs/containers/deploying-containers#viewing_logs).

