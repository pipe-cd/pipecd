---
title: "Secret management"
linkTitle: "Secret management"
weight: 9
description: >
  Storing secrets safely in the Git repository.
---

GitOps workflows use Git as the single source of truth for application configurations. Storing sensitive data such as credentials, API keys, and secrets directly in Git repositories poses security risks.

PipeCD's secret management feature allows you to store encrypted secrets in your Git repository alongside application manifests. The encrypted secrets are decrypted by `piped` during deployment operations.

## Prerequisites

Before using this feature, `piped` needs to be started with a key pair for secret encryption.

You can use the following command to generate a key pair:

``` console
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out private-key
openssl pkey -in private-key -pubout -out public-key
```

Then specify them while [installing](../../../installation/install-piped/installing-on-kubernetes) `piped` with these options:

``` console
--set-file secret.data.secret-public-key=PATH_TO_PUBLIC_KEY_FILE \
--set-file secret.data.secret-private-key=PATH_TO_PRIVATE_KEY_FILE
```

Finally, enable this feature in the Piped configuration file with the `secretManagement` field as below:

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  pipedID: your-piped-id
  ...
  secretManagement:
    type: KEY_PAIR
    config:
      privateKeyFile: /etc/piped-secret/secret-private-key
      publicKeyFile: /etc/piped-secret/secret-public-key
```

## How it works

The secret management workflow is as follows:

- Encrypt secret data using PipeCD's Web UI and store the encrypted data in Git
- `piped` automatically decrypts the encrypted secrets before performing deployment tasks

## Encrypting secret data

To encrypt secret data, navigate to the Applications page and click the "Encrypt Secret" button located in the top-left corner. Then, select a piped from the dropdown list, enter your secret data, and click the "ENCRYPT" button.
Copy the encrypted data to store in Git.

![Sealed Secret Button](/images/sealed-secret-button.png)
<p style="text-align: center;">
Applications page
</p>

<br>

![Sealed Secret Encrypting Drawer Form](/images/sealed-secret-encrypting-drawer-form.png)
<p style="text-align: center;">
The form for encrypting secret data
</p>

## Storing encrypted secrets in Git

To make encrypted secrets available to an application, specify them in the application configuration file.

- `encryptedSecrets` contains a list of the encrypted secrets.
- `decryptionTargets` contains a list of files that use one of the encrypted secrets and should be decrypted by `piped`.

``` yaml
apiVersion: pipecd.dev/v1beta1
# One of Piped defined app, for example: using the Kubernetes plugin
kind: Application
spec:
  encryption:
    encryptedSecrets:
      password: encrypted-data
    decryptionTargets:
      - secret.yaml
```

## Accessing encrypted secrets

Any file in the application directory can use the `.encryptedSecrets` context to access secrets you have encrypted and stored in the application configuration.

For example:

- Accessing by a Kubernetes Secret manifest

``` yaml
apiVersion: v1
kind: Secret
metadata:
  name: simple-sealed-secret
data:
  password: "{{ .encryptedSecrets.password }}"
```

- Configuring an ENV variable of a Lambda function to use an encrypted secret

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: LambdaFunction
spec:
  name: HelloFunction
  environments:
    KEY: "{{ .encryptedSecrets.key }}"
```

In all cases, `piped` decrypts the encrypted secrets and renders the decryption target files before using them to handle any deployment tasks.

<!-- ## Examples

- [examples/kubernetes/secret-management](https://github.com/pipe-cd/examples/tree/master/kubernetes/secret-management)
- [examples/cloudrun/secret-management](https://github.com/pipe-cd/examples/tree/master/cloudrun/secret-management)
- [examples/lambda/secret-management](https://github.com/pipe-cd/examples/tree/master/lambda/secret-management)
- [examples/terraform/secret-management](https://github.com/pipe-cd/examples/tree/master/terraform/secret-management)
- [examples/ecs/secret-management](https://github.com/pipe-cd/examples/tree/master/ecs/secret-management) -->

