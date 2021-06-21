---
title: "Sealed Secrets"
linkTitle: "Sealed Secrets"
weight: 1
description: >
  Storing secrets safely in the Git repository.
---

> NOTE: This feature is deprecated. It can still be used, but expected to be removed entirely sometime in the future.
> Instead, please use [Secret Management](/docs/user-guide/secret-management).

When doing GitOps, users want to use Git as a single source of truth. But storing credentials like Kubernetes Secret or Terraform's credentials in Git is not safe.
This feature helps you store those secret data safely in Git, right next to your application manifests.

Users encrypt their secret data from the Web UI and store the encrypted data in Git, `piped` will decrypt them before doing deployment tasks.

## Prerequisites

Before using this feature, `piped` needs to be started with an RSA key pair for secret encryption.

You can use the following command to generate a key pair:

``` console
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out private-key
openssl pkey -in private-key -pubout -out public-key
```

Then specify them while [installing](http://localhost:1313/docs/operator-manual/piped/installation/#installing-on-a-kubernetes-cluster) the `piped` with these options:

``` console
--set-file secret.sealedSecretSealingKey.publicKey.data=PATH_TO_PUBLIC_KEY_FILE \
--set-file secret.sealedSecretSealingKey.privateKey.data=PATH_TO_PRIVATE_KEY_FILE
```

And enable this feature in piped configuration file as below:

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  pipedID: your-piped-id
  ...
  sealedSecretManagement:
    type: SEALING_KEY
    config:
      privateKeyFile: /etc/piped-secret/sealed-secret-sealingkey-private-key
      publicKeyFile: /etc/piped-secret/sealed-secret-sealingkey-public-key
```

## Encrypting secret data

In order to encrypt the secret data, go to the application list page and click on the options icon at the right side of the application row, and choose "Encrypt Secret" option.
After that, input your secret data to the showed form and click on "ENCRYPT" button.
The encrypted data should be shown for you. Copy it to store in Git.

![](/images/sealed-secret-application-list.png)
<p style="text-align: center;">
Application list page
</p>

<br>

![](/images/sealed-secret-encrypting-form.png)
<p style="text-align: center;">
The form for encrypting secret data
</p>

## Storing the encrypted secret in Git

### Kubernetes example

Instead of Kubernetes Secret in Git, you store the `SealedSecret` file. This file contains the encrypted secret data and a template to render the needed Kubernetes Secret while decrypting.

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: SealedSecret
spec:
  template: |
    apiVersion: v1
    kind: Secret
    metadata:
        name: simple-sealed-secret
    data:
      password: {{ .encryptedItems.password }}
  encryptedItems:
    password: encrypted-data
```

In the application configuration file `.pipe.yaml`, you specify which SealedSecret files should be decrypted.

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  sealedSecrets:
    - path: sealed-secret.yaml
```

### Terraform example

You store the `SealedSecret` file containing the encrypted credentials.

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: SealedSecret
spec:
    encryptedData: encrypted-data
```

In the application configuration file `.pipe.yaml`, you specify which SealedSecret files should be decrypted.

``` yaml
apiVersion: pipecd.dev/v1beta1
kind: TerraformApp
spec:
  input:
    workspace: dev
  sealedSecrets:
    - path: service-account.yaml
      outFilename: service-account.json
      outDir: .terraform-credentials
```

And in your `tf` files, the credentials files can be accessed as below:

``` tf
provider "google" {
  project     = var.project
  credentials = ".terraform-credentials/service-account.json"
}
```
