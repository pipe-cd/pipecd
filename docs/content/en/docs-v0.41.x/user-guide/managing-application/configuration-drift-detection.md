---
title: "Configuration drift detection"
linkTitle: "Configuration drift detection"
weight: 8
description: >
  Automatically detecting the configuration drift.
---

Configuration Drift is a phenomenon where running resources of service become more and more different from the definitions in Git as time goes on, due to manual ad-hoc changes and updates.
As PipeCD is using Git as a single source of truth, all application resources and infrastructure changes should be done by making a pull request to Git. Whenever a configuration drift occurs it should be notified to the developers and be fixed.

PipeCD includes `Configuration Drift Detection` feature, which periodically compares running resources/configurations with the definitions in Git to detect the configuration drift and shows the comparing result in the application details web page as well as sends the notifications to the developers.

### Detection Result
There are three statuses for the drift detection result: `SYNCED`, `OUT_OF_SYNC`, `DEPLOYING`.

###### SYNCED

This status means no configuration drift was detected. All resources/configurations are synced from the definitions in Git. From the application details page, this status is shown by a green "Synced" mark.

![](/images/application-synced.png)
<p style="text-align: center;">
Application is in SYNCED state
</p>

###### OUT_OF_SYNC

This status means a configuration drift was detected. An application is in this status when at least one of the following conditions is satisfied:
- at least one resource is defined in Git but NOT running in the cluster
- at least one resource is NOT defined in Git but running in the cluster
- at least one resource that is both defined in Git and running in the cluster but NOT in the same configuration

This status is shown by a red "Out of Sync" mark on the application details page.

![](/images/application-out-of-sync.png)
<p style="text-align: center;">
Application is in OUT_OF_SYNC state
</p>

Click on the "SHOW DETAILS" button to see more details about why the application is in the `OUT_OF_SYNC` status. In the below example, the replicas number of a Deployment was not matching, it was `300` in Git but `3` in the cluster.

![](/images/application-out-of-sync-details.png)
<p style="text-align: center;">
The details shows why the application is in OUT_OF_SYNC state
</p>

###### DEPLOYING

This status means the application is deploying and the configuration drift detection is not running a white. Whenever a new deployment of the application was started, the detection process will temporarily be stopped until that deployment finishes and will be continued after that.

### How to enable

This feature is automatically enabled for all applications.

You can change the checking interval as well as [configure the notification](../../managing-piped/configuring-notifications/) for these events in `piped` configuration.

Note: If you want to trigger deployment automatically when `OUT_OF_SYNC` occurs, see [Trigger configuration](./triggering-a-deployment/#trigger-configuration).