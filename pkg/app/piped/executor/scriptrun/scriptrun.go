package scriptrun

import (
	"os"
	"os/exec"

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type registerer interface {
	Register(stage model.Stage, f executor.Factory) error
	RegisterRollback(kind model.RollbackKind, f executor.Factory) error
}

type Executor struct {
	executor.Input
}

func (e *Executor) Execute(sig executor.StopSignal) model.StageStatus {
	e.LogPersister.Infof("Start executing the script run stage")

	opts := e.Input.StageConfig.ScriptRunStageOptions
	if opts == nil {
		e.LogPersister.Infof("option for script run stage not found")
		return model.StageStatus_STAGE_FAILURE
	}

	if opts.Run == "" {
		return model.StageStatus_STAGE_SUCCESS
	}

	envs := make([]string, 0, len(opts.Env))
	for key, value := range opts.Env {
		envs = append(envs, key+"="+value)
	}

	cmd := exec.Command("/bin/sh", "-l", "-c", opts.Run)
	cmd.Env = append(os.Environ(), envs...)
	cmd.Stdout = e.LogPersister
	cmd.Stderr = e.LogPersister

	e.LogPersister.Infof("executing script:")
	e.LogPersister.Infof(opts.Run)

	if err := cmd.Run(); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}
	return model.StageStatus_STAGE_SUCCESS
}

type RollbackExecutor struct {
	executor.Input
}

func (e *RollbackExecutor) Execute(sig executor.StopSignal) model.StageStatus {
	return model.StageStatus_STAGE_NOT_STARTED_YET
}

// Register registers this executor factory into a given registerer.
func Register(r registerer) {
	r.Register(model.StageScriptRun, func(in executor.Input) executor.Executor {
		return &Executor{
			Input: in,
		}
	})
}
