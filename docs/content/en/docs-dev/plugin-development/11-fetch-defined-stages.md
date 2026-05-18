---
title: "Implementation: FetchDefinedStages"
weight: 11
description: >
  Implementing FetchDefinedStages to define custom deployment stages.
---

`FetchDefinedStages` returns the names of the stages supported by this plugin.

Our `file` plugin will support three stages: `FILE_DIFF`, `FILE_SYNC`, and `FILE_ROLLBACK`. Defining these stage names as constants is best practice so we can reuse them across other methods.

Update the `FetchDefinedStages` implementation and define the constants as follows:

```go
const (
	stageDiff     = "FILE_DIFF"
	stageSync     = "FILE_SYNC"
	stageRollback = "FILE_ROLLBACK"
)

func (plugin) FetchDefinedStages() []string {
	return []string{
		stageDiff,
		stageSync,
		stageRollback,
	}
}
```
