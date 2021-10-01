// Copyright 2020 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package waitapproval

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
	"github.com/pipe-cd/pipe/pkg/model"
)

const (
	approvedByKey = "ApprovedBy"
)

type Executor struct {
	executor.Input
}

type registerer interface {
	Register(stage model.Stage, f executor.Factory) error
}

// Register registers this executor factory into a given registerer.
func Register(r registerer) {
	f := func(in executor.Input) executor.Executor {
		return &Executor{
			Input: in,
		}
	}
	r.Register(model.StageWaitApproval, f)
}

// Execute starts waiting until an approval from one of the specified users.
func (e *Executor) Execute(sig executor.StopSignal) model.StageStatus {
	var (
		originalStatus = e.Stage.Status
		ctx            = sig.Context()
		ticker         = time.NewTicker(5 * time.Second)
	)
	defer ticker.Stop()
	timeout := e.StageConfig.WaitApprovalStageOptions.Timeout.Duration()
	timer := time.NewTimer(timeout)

	e.reportRequiringApproval(ctx)
	e.LogPersister.Info("Waiting for an approval...")
	for {
		select {
		case <-ticker.C:
			if commander, ok := e.checkApproval(ctx); ok {
				e.LogPersister.Infof("Got an approval from %s", commander)
				return model.StageStatus_STAGE_SUCCESS
			}

		case s := <-sig.Ch():
			switch s {
			case executor.StopSignalCancel:
				return model.StageStatus_STAGE_CANCELLED
			case executor.StopSignalTerminate:
				return originalStatus
			default:
				return model.StageStatus_STAGE_FAILURE
			}
		case <-timer.C:
			e.LogPersister.Errorf("Timed out %v", timeout)
			return model.StageStatus_STAGE_FAILURE
		}
	}
}

func (e *Executor) checkApproval(ctx context.Context) (string, bool) {
	var approveCmd *model.ReportableCommand
	commands := e.CommandLister.ListCommands()

	for i, cmd := range commands {
		if cmd.GetApproveStage() != nil {
			approveCmd = &commands[i]
			break
		}
	}
	if approveCmd == nil {
		return "", false
	}

	metadata := map[string]string{
		approvedByKey: approveCmd.Commander,
	}
	if ori, ok := e.MetadataStore.GetStageMetadata(e.Stage.Id); ok {
		for k, v := range ori {
			metadata[k] = v
		}
	}
	if err := e.MetadataStore.SetStageMetadata(ctx, e.Stage.Id, metadata); err != nil {
		e.LogPersister.Errorf("Unabled to save approver information to deployment, %v", err)
		return "", false
	}

	if err := approveCmd.Report(ctx, model.CommandStatus_COMMAND_SUCCEEDED, nil, nil); err != nil {
		e.Logger.Error("failed to report handled command", zap.Error(err))
	}
	return approveCmd.Commander, true
}

func (e *Executor) reportRequiringApproval(ctx context.Context) {
	ds, err := e.TargetDSP.GetReadOnly(ctx, e.LogPersister)
	if err != nil {
		e.LogPersister.Errorf("Failed to prepare running deploy source data (%v)", err)
		return
	}

	var approvers []string

	for _, v := range ds.GenericDeploymentConfig.DeploymentNotification.Mentions {
		if v.Event == "DEPLOYMENT_WAIT_APPROVAL" {
			approvers = v.Slack
		}
	}

	e.Notifier.Notify(model.NotificationEvent{
		Type: model.NotificationEventType_EVENT_DEPLOYMENT_WAIT_APPROVAL,
		Metadata: &model.NotificationEventDeploymentWaitApproval{
			Deployment:        e.Deployment,
			EnvName:           e.EnvName,
			MentionedAccounts: approvers,
		},
	})
}
