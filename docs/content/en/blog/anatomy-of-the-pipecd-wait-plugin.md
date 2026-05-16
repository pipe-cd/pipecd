---
date: 2026-05-17
title: "Reading the PipeCD WAIT plugin: an English walkthrough"
linkTitle: "Reading the WAIT plugin"
weight: 965
description: "A file-by-file tour of PipeCD's smallest official plugin, written for contributors who want to read a real pipedv1 plugin end-to-end before building their own."
author: Sridhar Panigrahi ([@sridhar-panigrahi](https://github.com/sridhar-panigrahi))
categories: ["Tutorial"]
tags: ["Plugin", "PipeCD", "pipedv1", "Tutorial"]
---

The PipeCD plugin architecture went alpha in [June 2025](/blog/2025/06/16/plugin-architecture-piped-alpha-version-has-been-released/), and the docs around it are catching up in pieces. There is a [great Japanese book](https://zenn.dev/warashi/books/try-and-learn-pipecd-plugin) by [@Warashi](https://github.com/Warashi) that walks through plugin development, and the v1 docs page on [plugin development resources](/docs-v1.0.x/contribution-guidelines/contributing-plugins/plugin-development-resources/) calls out that a full English translation is on the way. Until that lands, the fastest way to actually *get* what a plugin looks like is to read one.

I spent an evening reading the `wait` plugin — the smallest official plugin in the tree — and writing down what every file does, why it does it that way, and how `piped` runs all of it. This post is that read-along. By the end you should be able to open `pkg/app/pipedv1/plugin/wait/` and understand every line, and have a reasonable map for jumping into the bigger ones like `kubernetes` or `terraform`.

I am not going to assume you have written a plugin before. I will assume you have read or skimmed the [plugin architecture overview](/blog/2024/11/28/overview-of-the-plan-for-pluginnable-pipecd/) and know roughly what `piped` is.

## Why WAIT is the right first plugin to read

The `wait` plugin pauses a deployment pipeline for a configured duration. That is the whole feature. You drop it into a pipeline between, say, a canary stage and a primary rollout, and it sits there for thirty minutes before letting the next stage run. No platform-specific code, no Kubernetes client, no Terraform shelling out. Just a timer.

That tiny surface is exactly why it is the right plugin to start from. Every PipeCD plugin has to do the same handful of things — register with `piped`, expose a gRPC server, declare which stages it handles, execute a stage, persist progress so it can survive a restart, stream logs to the UI. The `wait` plugin does all of those things and nothing else. Read it once and you have a clean mental model. Then when you open `kubernetes/`, the plugin-specific complexity is the only new thing.

The `kubernetes_multicluster` plugin from [Mohammed's LFX work](/blog/2026/04/10/building-the-kubernetes-multi-cluster-plugin-for-pipecd-lfx-mentorship/) is a great next read once you have this one under your belt.

## The shape of a pipedv1 plugin in one paragraph

A plugin is a Go binary. `piped` knows about it because you list it in `piped-config.yaml`. On startup `piped` runs the binary as a subprocess and the plugin starts a gRPC server on a port `piped` assigned it. From then on, `piped` is the client and the plugin is the server. `piped` asks the plugin what stages it handles, asks it to plan, asks it to execute, and the plugin reports status back. When `piped` shuts down, it shuts the plugin down with it. That's the entire protocol at the level you need to think about it.

The SDK that makes this nice to write is `github.com/pipe-cd/piped-plugin-sdk-go`, [released independently from PipeCD itself](https://github.com/pipe-cd/pipecd/blob/master/RELEASES.md#piped-plugin-sdk-releases). You implement a Go interface and the SDK handles the gRPC wiring, the registration, the lifecycle.

## The nine files

Run `ls pkg/app/pipedv1/plugin/wait/` and you get this:

```
main.go         options.go       wait.go
plugin.go       options_test.go  wait_test.go
README.md       go.mod           go.sum
```

That is it. No vendored protos, no codegen, no helpers directory. A `wait` plugin authored from scratch could plausibly fit on a single screen if you removed the comments. Here is what each file is for:

| File | What it holds |
|---|---|
| `main.go` | Process entrypoint. Hands control to the SDK. |
| `plugin.go` | The struct that implements `sdk.StagePlugin`. |
| `wait.go` | The actual waiting, including the resume-after-restart logic. |
| `options.go` | The typed config schema for a WAIT stage, plus its decoder. |
| `wait_test.go` | Tests for the wait function itself. |
| `options_test.go` | Tests for the config decoder. |
| `README.md` | User-facing docs for configuring the plugin. |
| `go.mod`, `go.sum` | Plugin's own Go module — separate from the main PipeCD module. |

Three things to note before we start opening files. First, the plugin is its *own* Go module. That is the convention in the v1 plugin tree; each plugin can pull a different SDK version if it needs to. Second, there is no `proto/` directory — the gRPC types live inside the SDK, you never touch them directly. Third, every file is short. The longest one is `wait.go` at about 120 lines.

## main.go — the ten-line entrypoint

```go
package main

import (
    "log"

    sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

func main() {
    plugin, err := sdk.NewPlugin("0.0.1", sdk.WithStagePlugin(&plugin{}))
    if err != nil {
        log.Fatalln(err)
    }

    if err := plugin.Run(); err != nil {
        log.Fatalln(err)
    }
}
```

That is the whole file. There are two things happening.

`sdk.NewPlugin("0.0.1", sdk.WithStagePlugin(&plugin{}))` registers a single capability — that this binary is a *stage plugin*. The version string is the plugin's own version; it has nothing to do with the SDK version. The SDK supports a few capabilities besides stage plugins (`WithDeploymentPlugin`, `WithLiveStatePlugin`, `WithDriftPlugin`), and a real plugin can register more than one of them. The Kubernetes plugin registers all of them. `wait` registers only the stage capability because it has nothing to say about deploy targets, live state, or drift.

`plugin.Run()` blocks. Inside, the SDK opens the gRPC port `piped` told it about, serves the methods on `&plugin{}` over gRPC, and handles graceful shutdown when `piped` closes the connection. You will never write any of that code yourself.

Once you have read this file you have seen the full set of moving parts you need to start your own plugin: a `main.go` exactly like this one, with `WithStagePlugin` (or one of the other capabilities) pointing at the struct you are about to write.

## plugin.go — the contract

```go
package main

import (
    "context"

    sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

const (
    stageWait string = "WAIT"
)

type plugin struct{}

func (p *plugin) FetchDefinedStages() []string {
    return []string{stageWait}
}

func (p *plugin) BuildPipelineSyncStages(ctx context.Context, _ sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
    stages := make([]sdk.PipelineStage, 0, len(input.Request.Stages))
    for _, rs := range input.Request.Stages {
        stage := sdk.PipelineStage{
            Index:              rs.Index,
            Name:               rs.Name,
            Rollback:           false,
            Metadata:           map[string]string{},
            AvailableOperation: sdk.ManualOperationNone,
        }
        stages = append(stages, stage)
    }

    return &sdk.BuildPipelineSyncStagesResponse{
        Stages: stages,
    }, nil
}

func (p *plugin) ExecuteStage(ctx context.Context, _ sdk.ConfigNone, _ sdk.DeployTargetsNone, input *sdk.ExecuteStageInput[struct{}]) (*sdk.ExecuteStageResponse, error) {
    status := p.executeWait(ctx, input)
    return &sdk.ExecuteStageResponse{
        Status: status,
    }, nil
}
```

This is the `StagePlugin` contract. The SDK requires you to implement three methods, and `wait` implements them in the most minimal way possible.

`FetchDefinedStages` answers the question *which stage names does this plugin own?* When you write a pipeline in your application config, you put stage names in it — `WAIT`, `K8S_CANARY_ROLLOUT`, `ANALYSIS`, and so on. `piped` needs to know which plugin is responsible for which name, so it calls `FetchDefinedStages` on each plugin at startup and builds a lookup table. The `wait` plugin owns one name: `"WAIT"`. That's it.

`BuildPipelineSyncStages` is the planning step. When `piped` starts a deployment, it calls this on every plugin that appears in the pipeline, asking *given this set of pipeline stages, what do you actually want to execute, in what order, with what metadata?* For most plugins this is where you might inject auxiliary stages — a canary plugin might insert a "scale down" cleanup at the end, for example. The `wait` plugin has nothing to add, so it loops over the request's stages and returns them back as-is with empty metadata and no manual operation needed.

Pay attention to the two type parameters tucked into the method signatures: `sdk.ConfigNone` and `sdk.DeployTargetsNone`. These are sentinel types from the SDK that say *this plugin has no plugin-level config and no deploy targets.* Compare to the Kubernetes plugin, which would have a real config type there and a real deploy targets type. The `_` is what you write when you accept a sentinel — `wait` has nothing to do with either.

`ExecuteStage` is where the real work happens. It is called once per stage instance. Look at the third type parameter on `*sdk.ExecuteStageInput[struct{}]` — `struct{}`. That signals to the SDK *do not try to decode my stage config for me, I will do it myself.* This is intentional. `wait` decodes its own config inside `wait.go` so that the validation errors can be surfaced via the stage log persister and shown to the user in the UI, rather than failing at the gRPC layer. We'll see that in a moment. For now, the `ExecuteStage` method on `plugin.go` is one line of real work — delegate to `executeWait`.

## options.go — typed config, three steps

```go
type WaitStageOptions struct {
    Duration unit.Duration `json:"duration,omitempty"`
}

func (o WaitStageOptions) validate() error {
    if o.Duration <= 0 {
        return fmt.Errorf("duration must be greater than 0")
    }
    return nil
}

func decode(data json.RawMessage) (WaitStageOptions, error) {
    var opts WaitStageOptions
    if err := json.Unmarshal(data, &opts); err != nil {
        return WaitStageOptions{}, fmt.Errorf("failed to unmarshal the config: %w", err)
    }
    if err := defaults.Set(&opts); err != nil {
        return WaitStageOptions{}, fmt.Errorf("failed to set default values for stage config: %w", err)
    }
    if err := opts.validate(); err != nil {
        return WaitStageOptions{}, fmt.Errorf("failed to validate the config: %w", err)
    }
    return opts, nil
}
```

This is the smallest possible example of the pattern every plugin uses for its stage config. There is a typed struct, there is a validator on it, and there is a `decode` function that does three things in a fixed order: unmarshal JSON into the struct, fill in defaults via [`creasty/defaults`](https://github.com/creasty/defaults), then validate.

If you read the other plugins you will see this exact triple repeated. `kubernetes` does it. `terraform` does it. It is a convention worth noticing because it is what you should reach for in your own plugin instead of inventing a new pattern.

The `Duration` field uses `unit.Duration` from the SDK, not Go's standard `time.Duration`. The SDK type accepts string values like `"30s"` and `"5m"` from YAML/JSON, which `time.Duration` cannot do directly. Always use the SDK's `unit` types when accepting human-authored duration or size values.

Why decode here, inside `ExecuteStage`, rather than push it earlier? Two reasons. Decoding here means a bad config surfaces as a friendly error written to the stage log persister, which shows up in the deployment UI for the user to read — rather than as a generic decode failure further up the stack. And `decode` is reusable from `BuildPipelineSyncStages` too if a future change ever needs the config at planning time. The function is generic about *when* it runs.

## wait.go — the actual work

This is the longest file in the plugin, and it is still under 120 lines. Three things are worth pulling out: the orchestration in `executeWait`, the `wait` function itself with its cancellation handling, and the metadata persistence that makes the plugin survive a `piped` restart.

### Orchestration

```go
func (p *plugin) executeWait(ctx context.Context, in *sdk.ExecuteStageInput[struct{}]) sdk.StageStatus {
    slp, err := in.Client.StageLogPersister()
    if err != nil {
        in.Logger.Error("No stage log persister available", zap.Error(err))
        return sdk.StageStatusFailure
    }
    opts, err := decode(in.Request.StageConfig)
    if err != nil {
        slp.Errorf("failed to decode the stage config: %v", err)
        return sdk.StageStatusFailure
    }

    duration := opts.Duration.Duration()

    initialStart := p.retrieveStartTime(ctx, in.Client, in.Logger)
    if initialStart.IsZero() {
        initialStart = time.Now()
    }
    p.saveStartTime(ctx, in.Client, initialStart, in.Logger)

    return wait(ctx, duration, initialStart, slp)
}
```

Notice the two distinct logging surfaces. `in.Logger` is the plugin's own structured logger — what you would see in the plugin's stdout, useful for plugin operators. `in.Client.StageLogPersister()` is what shows up *in the UI for the user running the deployment*. They serve different audiences, and `executeWait` uses both.

The decode failure is reported via the log persister so the deployment operator can see why their stage failed. The "no log persister available" error is reported via the system logger because the user can't help with that — it's an infrastructure problem.

### The select loop

```go
func wait(ctx context.Context, duration time.Duration, initialStart time.Time, slp sdk.StageLogPersister) sdk.StageStatus {
    remaining := duration - time.Since(initialStart)
    if remaining <= 0 {
        slp.Infof("Already waited for %v since %v", duration, initialStart.Local())
        return sdk.StageStatusSuccess
    }

    timer := time.NewTimer(remaining)
    defer timer.Stop()

    ticker := time.NewTicker(logInterval)
    defer ticker.Stop()

    slp.Infof("Waiting for %v since %v...", duration, initialStart.Local())
    for {
        select {
        case <-timer.C:
            slp.Infof("Waited for %v", duration)
            return sdk.StageStatusSuccess

        case <-ticker.C:
            slp.Infof("%v elapsed...", time.Since(initialStart))

        case <-ctx.Done():
            slp.Info("Wait cancelled")
            return sdk.StageStatusFailure
        }
    }
}
```

A three-way select is a small thing, but it is doing real work. The timer is the actual wait. The ticker logs progress to the UI every ten seconds so the user does not stare at a frozen screen. The `ctx.Done()` case handles cancellation — when the user clicks "cancel deployment" in the UI, the gRPC call is cancelled, the context fires, and the plugin returns immediately.

There is one detail in the source worth quoting in full:

> `case <-ctx.Done():` *We can return any status here because the piped handles this case as cancelled by a user, ignoring the result from a plugin.*

That comment teaches you something the type signature does not. The status you return from a cancelled stage does not matter — `piped` records it as cancelled regardless. The plugin author returned `StageStatusFailure` because that is the safe default, but you could equally return `StageStatusSuccess` and `piped` would still mark the stage as cancelled. Worth remembering when you write your own plugins.

### Surviving a restart

This is the most interesting part of the whole plugin, and it would have been easy to skip writing it.

```go
const (
    startTimeKey = "startTime"
)

func (p *plugin) retrieveStartTime(ctx context.Context, client *sdk.Client, logger *zap.Logger) time.Time {
    sec, ok, err := client.GetStageMetadata(ctx, startTimeKey)
    if err != nil || !ok {
        return time.Time{}
    }
    ut, err := strconv.ParseInt(sec, 10, 64)
    if err != nil {
        return time.Time{}
    }
    return time.Unix(ut, 0)
}

func (p *plugin) saveStartTime(ctx context.Context, client *sdk.Client, t time.Time, logger *zap.Logger) {
    value := strconv.FormatInt(t.Unix(), 10)
    if err := client.PutStageMetadata(ctx, startTimeKey, value); err != nil {
        logger.Error("failed to store stage metadata", zap.Error(err))
    }
}
```

A plugin process can die. `piped` can be upgraded. The user can `kubectl rollout restart` the agent. Any of those things can happen while a `WAIT` stage is mid-flight, with twenty-eight of its thirty minutes remaining.

When `piped` restarts, it eventually calls `ExecuteStage` again on the same stage. If the plugin naively did `time.Now()` and waited another full thirty minutes, every restart would reset the clock. That is the bug `retrieveStartTime` and `saveStartTime` exist to prevent.

The first time `ExecuteStage` runs, `retrieveStartTime` returns zero, the plugin sets `initialStart = time.Now()`, and persists it. On any subsequent invocation — whether it is one second later or two hours later — `retrieveStartTime` returns the original start, and the `wait` function calculates `remaining` against *that*, not now. If two hours passed and the duration was one second, `wait` returns success immediately and the pipeline moves on. If the user cancelled deployments while `piped` was down for an hour, the wait simply expires the moment it resumes.

`client.GetStageMetadata` and `client.PutStageMetadata` are the SDK's bridge to per-stage persistence on the control plane. Anything you want to remember across a plugin restart goes here. Use it sparingly — it is durable storage, not scratch space.

## The tests that explain the design

The test file is worth reading not for the assertions but for what they *describe*. There are four test cases, and the names alone tell you what the plugin author cared about:

- `TestWait_Complete` — the happy path.
- `TestWait_Cancel` — context cancellation returns promptly.
- `TestWait_RestartAfterLongTime` — if more than `duration` has already elapsed when the stage resumes, return success immediately.
- `TestWait_RestartAndContinue` — if some time has elapsed but not enough, wait only the remaining time.

The last two only make sense if you've understood the metadata-persistence story above. They are the unit tests for the restart-survival behavior, and they exist precisely because that behavior is subtle and easy to break in a refactor.

`options_test.go` is more boring — it's a table test for the decoder that checks valid, invalid, empty, negative, and zero durations. Mention worth: this is a good template for the kind of test every plugin should have for its config decoder. Cheap to write, catches the gnarly user-facing failure modes.

## How piped runs all of this — the protocol

Here is the lifecycle in one ASCII diagram. Read it top-to-bottom.

```
                piped                                  wait plugin (binary)
                  │                                              │
  startup ───────►│  reads piped-config.yaml                     │
                  │  finds plugins[] entry for "wait"            │
                  │                                              │
                  │  fork+exec the wait binary, port=7002 ──────►│ main.go runs
                  │                                              │ sdk.NewPlugin(...)
                  │                                              │ plugin.Run() blocks
                  │  gRPC dial localhost:7002 ─────────────────►│ gRPC server up
                  │                                              │
                  │  call FetchDefinedStages() ─────────────────►│
                  │◄──────────── ["WAIT"] ───────────────────────│
                  │                                              │
                  │  build stage→plugin lookup table             │
                  │                                              │
  deployment ────►│                                              │
                  │  for each WAIT in the pipeline:              │
                  │  call BuildPipelineSyncStages() ────────────►│
                  │◄────────── stage descriptors ────────────────│
                  │                                              │
                  │  schedule stage, when its turn comes:        │
                  │  call ExecuteStage(stageConfig) ────────────►│ executeWait()
                  │                                              │   decode config
                  │                                              │   load/save startTime
                  │                                              │   wait(ctx, ...)
                  │                                              │
                  │◄──── log lines via StageLogPersister ────────│   slp.Infof(...)
                  │                                              │
                  │◄──── metadata via PutStageMetadata ──────────│   client.PutStageMetadata(...)
                  │                                              │
                  │◄────── StageStatusSuccess ───────────────────│ return
                  │                                              │
  user cancels ──►│  context cancellation propagates ───────────►│ <-ctx.Done()
                  │◄──────── (status ignored) ───────────────────│ return
                  │                                              │
  shutdown ──────►│  close gRPC, SIGTERM the subprocess ────────►│ plugin.Run() unblocks
                  │                                              │ process exits
```

A few things become obvious from the diagram that are not obvious from reading the code. The gRPC channel is bidirectional — `piped` calls plugin methods, but the plugin can also call back into `piped`'s service through `in.Client` for things like log persistence and metadata. The plugin binary's lifetime is bound to `piped`'s; you do not run plugins as separate services. And every `ExecuteStage` call is independent — the plugin is allowed to (and in fact must) be idempotent across restarts of the same stage.

## A thirty-minute exercise: add a jitter option

If you want to actually *do* something with what you just read, here is a small change that touches every file you would touch in a real plugin contribution.

The goal: extend the `WAIT` stage with an optional `jitter` field, so `duration: 5m, jitter: 30s` waits somewhere between 4m30s and 5m30s. Useful for spreading out cluster-wide cooldowns.

Roughly the steps:

1. **`options.go`**: add a `Jitter unit.Duration` field to `WaitStageOptions`, default it to zero, validate that `jitter <= duration`. Update the JSON tag.
2. **`wait.go`**: in `executeWait`, after `duration := opts.Duration.Duration()`, compute `jittered := duration + randomInRange(-opts.Jitter, +opts.Jitter)`. Pass `jittered` to `wait()`. Persist the *jittered* value, not the original — otherwise the restart logic uses the wrong target.
3. **`options_test.go`**: add cases for `jitter > duration` (should fail validation), `jitter == 0` (should behave exactly like today), and a positive case.
4. **`wait_test.go`**: nothing changes, because the jittered value is computed before `wait()` runs. That is actually a nice property of how the code is factored.
5. **`README.md`**: add the new field to the table.

Building and running the plugin locally is documented in [Plugin development resources](/docs-v1.0.x/contribution-guidelines/contributing-plugins/plugin-development-resources/) — the short version is `make build/plugin` from the repo root, then `make run/piped CONFIG_FILE=...  EXPERIMENTAL=true INSECURE=true`. Remember to clear `~/.piped/plugins` after rebuilding or you will keep loading the old binary.

If you do this exercise, you have effectively shipped your first plugin change. The pattern generalizes — every plugin you contribute to will follow the same loop of options → behavior → tests → README.

## Where I'm headed next

I am applying for the [LFX 2026 Term 2 PipeCD mentorship](https://mentorship.lfx.linuxfoundation.org/) on the *Plugin Development Book, Docs DX, and Adoption Growth* track. The Japanese book by [@Warashi](https://github.com/Warashi) is the foundation; the plan is to bring it into English inside the PipeCD docs, build out the missing chapters around plugin development, and write more posts like this one to help new contributors get from "I just heard about PipeCD" to "I just shipped my first plugin PR" with as little friction as possible.

If you read this far and any of it was unclear, that is useful feedback — please tell me in [#pipecd on CNCF Slack](https://cloud-native.slack.com/archives/C01B27F9T0X) or open an issue. The next walkthrough will be the `wait` plugin's slightly bigger sibling, `waitapproval`, and after that I want to take a swing at the kubernetes plugin's deployment interface — the part where things stop being a timer and start being a real CD platform.

Thanks for reading.
