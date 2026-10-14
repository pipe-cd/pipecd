---
title: "Building sync stages: pipeline sync and quick sync"
linkTitle: "Building sync stages: pipeline sync and quick sync"
weight: 5
description: >
  Implement BuildPipelineSyncStages and BuildQuickSyncStages.
---

In the previous chapter the plugin told `piped` which stages it provides and which strategy a deployment uses. This chapter builds the stage lists that each strategy runs. `piped` calls `BuildPipelineSyncStages` when a deployment uses pipeline sync, and `BuildQuickSyncStages` when it uses quick sync. Neither method runs a deployment; each one assembles the ordered plan that `piped` executes later.

## Build the pipeline sync stages

`piped` calls `BuildPipelineSyncStages` for a pipeline sync deployment. It passes the stages the user defined in the application's pipeline, and the plugin returns the stage list `piped` runs.

For example, suppose the user defines this pipeline:

```yaml
pipeline:
  stages:
    - name: FILE_DIFF
    - name: FILE_SYNC
```

`piped` passes those stages to the plugin with their positions as indexes:

```go
// input.Request.Stages
[]sdk.StageConfig{
    {Name: "FILE_DIFF", Index: 0},
    {Name: "FILE_SYNC", Index: 1},
}
```

`BuildPipelineSyncStages` returns them unchanged. When `Rollback` is true, it also appends a rollback stage:

```go
// returned Stages, Rollback = true
[]sdk.PipelineStage{
    {Name: "FILE_DIFF", Index: 0},
    {Name: "FILE_SYNC", Index: 1},
    {Name: "FILE_ROLLBACK", Index: 0, Rollback: true}, // Index matches a requested stage; piped still runs it last
}
```

Three rules shape the result:

- Each returned stage carries an `Index`, which sets the order the stage runs in. Return the same `Index` that came in on the matching request stage. Returning an `Index` that was not in the request is an error, so copy it straight through.
- The file plugin only knows two pipeline stages, `FILE_DIFF` and `FILE_SYNC`. Treat any other name as an error rather than passing it on.
- When `input.Request.Rollback` is true, append a rollback stage in addition to the user's stages. `piped` always runs rollback stages after every normal stage, so this `Index` does not position the rollback within the pipeline. It must match one of the `Index` values from the request, and it only orders rollback stages relative to one another when several plugins contribute them. Set it to the smallest requested `Index` as a safe, valid choice.

Replace the empty method with the following:

```go
func (p *plugin) BuildPipelineSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	if len(input.Request.Stages) == 0 {
		return nil, fmt.Errorf("no stages defined in the request")
	}

	stages := make([]sdk.PipelineStage, 0, len(input.Request.Stages)+1) // +1 leaves room for the rollback stage without a second allocation
	for _, s := range input.Request.Stages {
		switch s.Name {
		case stageDiff, stageSync:
			stages = append(stages, sdk.PipelineStage{
				Index: s.Index,
				Name:  s.Name,
			})
		default:
			return nil, fmt.Errorf("unknown stage: %s", s.Name)
		}
	}

	if input.Request.Rollback {
		// The rollback stage's Index must match one from the request; it only orders
		// rollback stages across plugins, since piped runs them after the normal stages.
		minIndex := input.Request.Stages[0].Index
		for _, s := range input.Request.Stages[1:] {
			if s.Index < minIndex {
				minIndex = s.Index
			}
		}
		stages = append(stages, sdk.PipelineStage{
			Index:    minIndex,
			Name:     stageRollback,
			Rollback: true,
		})
	}

	return &sdk.BuildPipelineSyncStagesResponse{Stages: stages}, nil
}
```

## Build the quick sync stages

`piped` calls `BuildQuickSyncStages` for a quick sync deployment. Quick sync ignores the user's pipeline and runs a fixed, minimal set of stages the plugin defines. For the file plugin that is a single `FILE_SYNC` stage, plus a `FILE_ROLLBACK` stage when a rollback is requested.

The result differs from pipeline sync in two ways:

- There is no `Index`. Quick sync stages are not ordered, so the plugin does not set one.
- The plugin sets `Description` itself. In pipeline sync the description comes from the user's configuration, but quick sync reads no user configuration, so the plugin supplies the text.

Replace the empty method with the following:

```go
func (p *plugin) BuildQuickSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	stages := make([]sdk.QuickSyncStage, 0, 2)
	stages = append(stages, sdk.QuickSyncStage{
		Name:        stageSync,
		Description: "Sync by applying the files in the deployment source",
	})

	if input.Request.Rollback {
		stages = append(stages, sdk.QuickSyncStage{
			Name:        stageRollback,
			Description: "Rollback to the previously applied files",
			Rollback:    true,
		})
	}

	return &sdk.BuildQuickSyncStagesResponse{Stages: stages}, nil
}
```

Both methods use `fmt`, so add it to the import block:

```go
import (
	"context"
	"fmt"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)
```

Build the project again to confirm the two methods compile:

```bash
go build ./...
```

The plugin can now plan a deployment either way: it returns the user's pipeline for a pipeline sync, and a fixed sync stage for a quick sync. In the next chapter you start running these stages with `ExecuteStage`, beginning with the `FILE_DIFF` stage.
