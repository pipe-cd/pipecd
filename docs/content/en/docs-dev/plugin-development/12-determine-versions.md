---
title: "Implementation: DetermineVersions"
weight: 12
description: >
  Implementing DetermineVersions to manage version tracking.
---

`DetermineVersions` returns the version of the resource being deployed.

Since our plugin synchronizes files from a Git repository, we will use the Git commit hash of the deployment source as our resource version.

The deployment source details are passed via `input.Request.DeploymentSource`. We will fetch `CommitHash` from it. We do not need to specify the artifact name or URL for this local plugin.

Update `DetermineVersions` as follows:

```go
func (plugin) DetermineVersions(_ context.Context, _ *config, input *sdk.DetermineVersionsInput[applicationConfig]) (*sdk.DetermineVersionsResponse, error) {
	return &sdk.DetermineVersionsResponse{
		Versions: []sdk.ArtifactVersion{{Version: input.Request.DeploymentSource.CommitHash}},
	}, nil
}
```
