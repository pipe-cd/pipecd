---
title: "Implementation: DetermineStrategy"
weight: 13
description: >
  Implementing DetermineStrategy to configure sync execution strategies.
---

`DetermineStrategy` decides whether to run a `PipelineSync` or a `QuickSync` strategy.

- **`PipelineSync`**: Executes the deployment using a pipeline explicitly defined by the user in `app.pipecd.yaml`.
- **`QuickSync`**: Bypasses any user-defined pipelines and executes a default set of sync stages defined by the plugin. For example, in Kubernetes deployments, a minor change (like updating replicas) might bypass canary release pipelines using a Quick Sync.

If your plugin does not contain complex logic to dynamically choose between `QuickSync` and `PipelineSync`, you can simply return `nil, nil`. In this case, Piped will fall back to the strategy configured in the application configuration.

Update `DetermineStrategy` as follows:

```go
func (plugin) DetermineStrategy(context.Context, *config, *sdk.DetermineStrategyInput[applicationConfig]) (*sdk.DetermineStrategyResponse, error) {
	return nil, nil
}
```
