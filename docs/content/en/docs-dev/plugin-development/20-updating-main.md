---
title: "Modifying the Main Function"
weight: 20
description: >
  Updating the main function to register the completed DeploymentPlugin.
---

Now that our plugin struct satisfies all interface methods, we can complete the `main` function in `main.go`.

We will supply our `plugin{}` instance to `sdk.NewPlugin` using the registration option `sdk.WithDeploymentPlugin(plugin{})`.

Update `main.go` as follows:

```go
package main

import (
	"log"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go/pkg/plugin/sdk"
)

func main() {
	plugin, err := sdk.NewPlugin("0.0.1", sdk.WithDeploymentPlugin(plugin{}))
	if err != nil {
		log.Fatalln(err)
	}

	if err := plugin.Run(); err != nil {
		log.Fatalln(err)
	}
}
```
