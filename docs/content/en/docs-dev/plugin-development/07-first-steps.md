---
title: "First Step of Plugin Implementation"
weight: 7
description: >
  Creating the main function and starting the gRPC server.
---

Let's begin writing our plugin code by implementing the `main` function.

Create a file named `main.go` and save the following code:

```go
package main

import (
	"log"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

func main() {
	plugin, err := sdk.NewPlugin[any, any, any]("0.0.1")
	if err != nil {
		log.Fatalln(err)
	}

	if err := plugin.Run(); err != nil {
		log.Fatalln(err)
	}
}
```

The type parameters passed to `sdk.NewPlugin` are temporary placeholders. By the time we finish implementing the plugin, Go's type inference will determine these automatically, but for now we must explicitly write them to compile.

### Note on Plugin Naming

You might have noticed that the string `file` (the name of our plugin) does not appear anywhere in this code.

Indeed, a plugin binary has no concept of its own name. The naming is entirely determined in the Piped configuration written by the user. A user could configure Piped to recognize this plugin binary as `file`, or `filesystem`, or any other name. The name itself is decoupled from the binary.

### Running this Code

At this point, if you attempt to run `go run main.go`, it will fail because `sdk.NewPlugin` returns an error if no plugin implementation is registered. In the next section, we will define a struct and satisfy the `DeploymentPlugin` interface.
