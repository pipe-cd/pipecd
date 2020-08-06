---
title: "Configuration drift detection"
linkTitle: "Configuration drift detection"
weight: 13
description: >
  Automatically detecting the configuration drift.
---

> TBA

Configuration Drift is the phenomenon where running resources become more and more different with the defined in Git as time goes on, due to manual ad-hoc changes and updates.
As PipeCD is using Git as single source of truth, all application resources and infrastructure changes must be done by making a pull request to Git.

![](/images/application-synced.png)
<p style="text-align: center;">
Application is in SYNCED state
</p>

![](/images/application-out-of-sync.png)
<p style="text-align: center;">
Application is in OUT_OF_SYNC state
</p>

![](/images/application-out-of-sync-details.png)
<p style="text-align: center;">
The details shows why the application is in OUT_OF_SYNC state
</p>
