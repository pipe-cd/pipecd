---
title: "Implementation: ExecuteStage"
weight: 16
description: >
  Implementing the main entrypoint for stage execution routing.
---

`ExecuteStage` is the core method where actual synchronization, diffing, or rolling back takes place.

We will route execution to dedicated helper methods based on the `input.Request.StageName`.

### Crucial Implementation Notes

1. **Web UI Logs**: To persist logs that are viewable by users in the PipeCD Web Console, do **not** use `input.Logger`. Instead, fetch the specialized log persister using `input.Client.LogPersister()`.
2. **Stage Failure Handling**: When a stage fails due to a soft error (such as a deployment conflict), do **not** return a Go `error`. Instead, return a successful response with `Status` set to `sdk.StageStatusFailure`. Returning a Go `error` indicates a hard plugin crash.

Update `ExecuteStage` and scaffold the helper methods as follows:

```go
func (p plugin) ExecuteStage(ctx context.Context, _ *config, _ []*sdk.DeployTarget[deployTargetConfig], input *sdk.ExecuteStageInput[applicationConfig]) (*sdk.ExecuteStageResponse, error) {
	switch input.Request.StageName {
	case stageDiff:
		return p.executeStageDiff(ctx, input)
	case stageSync:
		return p.executeStageSync(ctx, input)
	case stageRollback:
		return p.executeStageRollback(ctx, input)
	default:
		return nil, fmt.Errorf("unknown stage: %s", input.Request.StageName)
	}
}

func (plugin) executeStageDiff(ctx context.Context, input *sdk.ExecuteStageInput[applicationConfig]) (*sdk.ExecuteStageResponse, error) {
	panic("unimplemented")
}

func (plugin) executeStageSync(ctx context.Context, input *sdk.ExecuteStageInput[applicationConfig]) (*sdk.ExecuteStageResponse, error) {
	panic("unimplemented")
}

func (plugin) executeStageRollback(ctx context.Context, input *sdk.ExecuteStageInput[applicationConfig]) (*sdk.ExecuteStageResponse, error) {
	panic("unimplemented")
}
```
