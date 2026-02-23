---
title: "Remote upgrade and remote config"
linkTitle: "Remote upgrade and remote config"
weight: 1
description: >
  This page describes how to use remote upgrade and remote config features.
---

>**NOTE:**
>
>Remote Upgrade and Remote Config features are not available right now for PipeCD versions 1.0.0 and above.
> We are working on them, and will update this page as soon as they are ready.

This page will be updated when these features are available for PipeCD v1.

<!-- ## Remote upgrade

The remote upgrade is the ability to restart the currently running `piped` with another version from the web console.
This reduces the effort involved in updating `piped` to newer versions.
All `piped` instances that are running by the provided `piped` container image can be enabled to use this feature.
It means `piped` instances running on a Kubernetes cluster, a virtual machine, a serverless service can be upgraded remotely from the web console.

In order to use this feature you must run `piped` with `/launcher` command instead of `/piped` command as usual.
Please check the [installation](../../../installation/install-piped/) guide on each environment to see the details.

After starting `piped` with the remote-upgrade feature, you can go to the Settings page then click on `UPGRADE` button on the top-right corner.
A dialog will be shown for selecting which `piped` instances you want to upgrade and what version they should run.

![Remote Upgrade Settings](/images/settings-remote-upgrade.png)
<p style="text-align: center;">
Select a list of piped instances to upgrade from Settings page
</p>

## Remote config

Although the remote-upgrade allows you remotely restart your `piped` instances to run any new version you want, if your `piped` is loading its config locally where `piped` is running, you still need to manually restart `piped` after adding any change on that config data. Remote-config is for you to remove that kind of manual operation.

Remote-config is the ability to load `piped` config data from a remote location such as a Git repository. Not only that, but it also watches the config periodically to detect any changes on that config and restarts `piped` to reflect the new configuration automatically.

This feature requires the remote-upgrade feature to be enabled simultaneously. Please check [Runtime Options](runtime-options.md) and the [installation](../../../installation/install-piped/) guide on each environment to know how to configure `piped` to load a remote config file.

## Summary

- By `remote-upgrade` you can upgrade your `piped` to a newer version by clicking on the web console
- By `remote-config` you can enforce your `piped` to use the latest config data just by updating its config file stored in a Git repository -->
