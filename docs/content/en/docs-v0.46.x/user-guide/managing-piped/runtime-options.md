---
title: "Runtime Options"
linkTitle: "Runtime Options"
weight: 11
description: >
  This page describes configurable options for executing Piped and launcher.
---

You can configure some options when running Piped and launcher.

## Options for Piped

```
Usage:
  piped piped [flags]

Flags:
      --add-login-user-to-passwd                   Whether to add login user to $HOME/passwd. This is typically for applications running as a random user ID.
      --admin-port int                             The port number used to run a HTTP server for admin tasks such as metrics, healthz. (default 9085)
      --cert-file string                           The path to the TLS certificate file.
      --config-aws-secret string                   The ARN of secret that contains Piped config and be stored in AWS Secrets Manager.
      --config-data string                         The base64 encoded string of the configuration data.
      --config-file string                         The path to the configuration file.
      --config-gcp-secret string                   The resource ID of secret that contains Piped config and be stored in GCP SecretManager.
      --enable-default-kubernetes-cloud-provider   Whether the default kubernetes provider is enabled or not. This feature is deprecated.
      --grace-period duration                      How long to wait for graceful shutdown. (default 30s)
  -h, --help                                       help for piped
      --insecure                                   Whether disabling transport security while connecting to control-plane.
      --launcher-version string                    The version of launcher which initialized this Piped.
      --tools-dir string                           The path to directory where to install needed tools such as kubectl, helm, kustomize. (default "/Users/s24798/.piped/tools")

Global Flags:
      --log-encoding string                The encoding type for logger [json|console|humanize]. (default "humanize")
      --log-level string                   The minimum enabled logging level. (default "info")
      --metrics                            Whether metrics is enabled or not. (default true)
      --profile                            If true enables uploading the profiles to Stackdriver.
      --profile-debug-logging              If true enables logging debug information of profiler.
      --profiler-credentials-file string   The path to the credentials file using while sending profiles to Stackdriver.
```

## Options for launcher

```
Usage:
  launcher launcher [flags]

Flags:
      --aws-secret-id string           The ARN of secret that contains Piped config in AWS Secrets Manager service.
      --cert-file string               The path to the TLS certificate file.
      --check-interval duration        Interval to periodically check desired config/version to restart Piped. Default is 1m. (default 1m0s)
      --config-data string             The base64 encoded string of the configuration data.
      --config-file string             The path to the configuration file.
      --config-from-aws-secret         Whether to load Piped config that is being stored in AWS Secrets Manager service.
      --config-from-gcp-secret         Whether to load Piped config that is being stored in GCP SecretManager service.
      --config-from-git-repo           Whether to load Piped config that is being stored in a git repository.
      --default-version string         The version should be run when no desired version was specified. Empty means using the same version with Launcher.
      --gcp-secret-id string           The resource ID of secret that contains Piped config in GCP SecretManager service.
      --git-branch string              Branch of git repository to for Piped config.
      --git-piped-config-file string   Relative path within git repository to locate Piped config file.
      --git-repo-url string            The remote URL of git repository to fetch Piped config.
      --git-ssh-key-file string        The path to SSH private key to fetch private git repository.
      --grace-period duration          How long to wait for graceful shutdown. (default 30s)
  -h, --help                           help for launcher
      --home-dir string                The working directory of Launcher.
      --insecure                       Whether disabling transport security while connecting to control-plane.
      --launcher-admin-port int        The port number used to run a HTTP server for admin tasks such as metrics, healthz.

Global Flags:
      --log-encoding string                The encoding type for logger [json|console|humanize]. (default "humanize")
      --log-level string                   The minimum enabled logging level. (default "info")
      --metrics                            Whether metrics is enabled or not. (default true)
      --profile                            If true enables uploading the profiles to Stackdriver.
      --profile-debug-logging              If true enables logging debug information of profiler.
      --profiler-credentials-file string   The path to the credentials file using while sending profiles to Stackdriver.
```
