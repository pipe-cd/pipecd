---
title: "Implementation: Defining Configuration Types"
weight: 9
description: >
  Defining the configuration Go structs for our file plugin.
---

Our `file` plugin needs to allow users to specify a unique target directory path per application. While we could theoretically define this in a Deploy Target, requiring a separate Deploy Target for every single directory would be tedious for users. Therefore, we will define the target directory path within each application's configuration.

For future extensions—for instance, if we want this plugin to deploy files not just to the local machine but to remote hosts via SSH—the SSH connection details (host, keys, etc.) would be defined within the Deploy Target configuration, whereas the specific path on that host would remain in the Application configuration.

Add the following type definitions to `main.go` (or a separate file under `package main`):

```go
type (
	// config represents the global plugin configuration.
	// Since our file plugin does not require global settings, we leave it empty.
	config struct{}

	// deployTargetConfig represents the deployment target settings.
	// Since our file plugin only interacts with the local machine, we leave it empty.
	deployTargetConfig struct{}

	// applicationConfig represents the configuration defined per application in app.pipecd.yaml.
	// Here, we define the target directory path where the files should be synchronized.
	applicationConfig struct {
		// Path specifies the absolute path on the local file system where the files should be synced.
		Path string `json:"path"`
	}
)
```
