---
title: "Secret management"
linkTitle: "Secret management"
weight: 9
description: >
  Storing secrets safely in the Git repository.
---

When doing GitOps, users want to use Git as a single source of truth. But storing credentials like Kubernetes Secret or Terraform's credentials directly in Git is not safe.
This feature helps you keep that sensitive information safely in Git, right next to your application manifests.

Basically, the flow will look like this:
- user encrypts their secret data via the PipeCD's Web UI and stores the encrypted data in Git
- `Piped` decrypts them before doing deployment tasks

## Prerequisites

Before using this feature, `Piped` needs to be started with a key pair for secret encryption.

You can use the following command to generate a key pair:

``` console
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out private-key
openssl pkey -in private-key -pubout -out public-key
```

Then specify them while [installing](../../../installation/install-piped/installing-on-kubernetes) the `Piped` with these options:

``` console
--set-file secret.data.secret-public-key=PATH_TO_PUBLIC_KEY_FILE \
--set-file secret.data.secret-private-key=PATH_TO_PRIVATE_KEY_FILE
```

Finally, enable this feature in Piped configuration file with `secretManagement` field as below:

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

In order to encrypt the secret data, navigate to the Applications page and click the “Encrypt Secret“ button located in the top-left corner. Then, select a piped from the dropdown list, after that enter your secret data, and click the “ENCRYPT“ button.
The encrypted data should be shown for you. Copy it to store in Git.

![](/images/sealed-secret-button.png)
<p style="text-align: center;">
Applications page
</p>

<br>

![](/images/sealed-secret-encrypting-drawer-form.png)
<p style="text-align: center;">
The form for encrypting secret data
</p>

## Storing encrypted secrets in Git

To make encrypted secrets available to an application, they must be specified in the application configuration file of that application.

- `encryptedSecrets` contains a list of the encrypted secrets.
- `decryptionTargets` contains a list of files that are using one of the encrypted secrets and should be decrypted by `Piped`.

``` yaml
apiVersion: pipecd.dev/v1beta1
# One of Piped defined app kind such as: KubernetesApp
kind: {APPLICATION_KIND}
spec:
  encryption:
    encryptedSecrets:
      password: encrypted-data
    decryptionTargets:
      - secret.yaml
```

## Accessing encrypted secrets

Any file in the application directory can use `.encryptedSecrets` context to access secrets you have encrypted and stored in the application configuration.

For example,

- Accessing by a Kubernetes Secret manifest

``` yaml
apiVersion: v1
kind: Secret
metadata:
  name: simple-sealed-secret
data:
  password: "{{ .encryptedSecrets.password }}"
```

- Configuring ENV variable of a Lambda function to use an encrypted secret

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: LambdaFunction
spec:
  name: HelloFunction
  environments:
    KEY: "{{ .encryptedSecrets.key }}"
```

In all cases, `Piped` will decrypt the encrypted secrets and render the decryption target files before using them to handle any deployment tasks.

## Examples

- [examples/kubernetes/secret-management](https://github.com/pipe-cd/examples/tree/master/kubernetes/secret-management)
- [examples/cloudrun/secret-management](https://github.com/pipe-cd/examples/tree/master/cloudrun/secret-management)
- [examples/lambda/secret-management](https://github.com/pipe-cd/examples/tree/master/lambda/secret-management)
- [examples/terraform/secret-management](https://github.com/pipe-cd/examples/tree/master/terraform/secret-management)
- [examples/ecs/secret-management](https://github.com/pipe-cd/examples/tree/master/ecs/secret-management)
