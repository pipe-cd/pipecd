---
title: "Implementation: BuildPipelineSyncStages"
weight: 14
description: >
  Implementing BuildPipelineSyncStages to construct user-defined pipelines.
---

`BuildPipelineSyncStages` is triggered when a `PipelineSync` strategy is selected. This method takes the pipeline defined by the user in `app.pipecd.yaml` and constructs the actual execution graph.

### Key Considerations

1. **`Index` Matching**: The returned stages must maintain the same `Index` values provided in `input.Request.Stages`. Returning an invalid or mismatched `Index` will result in a deployment error.
2. **Rollback Stage Integration**: If `input.Request.Rollback` is `true`, we must append a rollback stage (`FILE_ROLLBACK`). To ensure rollback starts immediately in case of failure, set its `Index` to the minimum index among all defined stages.

Since our plugin does not require custom stage metadata, we can map stage names to our constants using a straightforward `switch` block.

Update `BuildPipelineSyncStages` as follows:

```go
import "fmt"

func (plugin) BuildPipelineSyncStages(_ context.Context, _ *config, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	if len(input.Request.Stages) == 0 {
		return nil, fmt.Errorf("no stages defined in the request")
	}

	stages := make([]sdk.PipelineStage, 0, len(input.Request.Stages)+1) // +1 for the rollback stage
	for _, s := range input.Request.Stages {
		switch s.Name {
		case stageDiff:
			stages = append(stages, sdk.PipelineStage{
				Index: s.Index,
				Name:  stageDiff,
			})
		case stageSync:
			stages = append(stages, sdk.PipelineStage{
				Index: s.Index,
				Name:  stageSync,
			})
		default:
			return nil, fmt.Errorf("unknown stage: %s", s.Name)
		}
	}

	if input.Request.Rollback {
		// Find the minimum index to assign to the rollback stage
		idx := input.Request.Stages[0].Index
		for _, s := range input.Request.Stages[1:] {
			if s.Index < idx {
				idx = s.Index
			}
		}
		stages = append(stages, sdk.PipelineStage{
			Index:    idx,
			Name:     stageRollback,
			Rollback: true,
		})
	}

	return &sdk.BuildPipelineSyncStagesResponse{
		Stages: stages,
	}, nil
}
```
