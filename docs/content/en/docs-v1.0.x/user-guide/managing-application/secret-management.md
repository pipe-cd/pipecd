---
title: "Secret management"
linkTitle: "Secret management"
weight: 9
description: >
  Storing secrets safely in the Git repository.
---

When doing GitOps, you want to use Git as a single source of truth. However, storing credentials like Kubernetes Secrets or Terraform credentials directly in Git is not safe.

This feature allows you to keep sensitive information safely in Git, right next to your application manifests.

The basic flow works as follows:

- You encrypt your secret data via PipeCD's Web UI and store the encrypted data in Git
- `piped` decrypts them before performing deployment tasks

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

## Encrypting secret data

To encrypt secret data, navigate to the Applications page and click the "Encrypt Secret" button located in the top-left corner. Then, select a piped from the dropdown list, enter your secret data, and click the "ENCRYPT" button.
The encrypted data is displayed for you. Copy it to store in Git.

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

To make encrypted secrets available to an application, specify them in the application configuration file of that application.

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

