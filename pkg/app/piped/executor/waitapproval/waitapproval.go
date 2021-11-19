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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

const (
	approvedByKey    = "ApprovedBy"
	minApproverNum   = "MinApproverNum"
	approversKey     = "CurrentApprovers"
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

	e.reportRequiringApproval()
	e.LogPersister.Info("Waiting for an approval...")
	for {
		select {
		case <-ticker.C:
			if commander, ok := e.checkApproval(ctx); ok {
				e.reportApproved(commander)
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

	as, ok := e.validateApproverNum(ctx, approveCmd.Commander)
	if !ok {
		if len(as) > 0 {
			if err := e.MetadataStore.Stage(e.Stage.Id).Put(ctx, approversKey, as); err != nil {
				e.LogPersister.Errorf("Unabled to save approver information to deployment, %v", err)
			}
		}
		return "", false
	}

	metadata := map[string]string{
		approvedByKey: as,
	}
	if err := e.MetadataStore.Stage(e.Stage.Id).PutMulti(ctx, metadata); err != nil {
		e.LogPersister.Errorf("Unabled to save approver information to deployment, %v", err)
		return "", false
	}

	if err := approveCmd.Report(ctx, model.CommandStatus_COMMAND_SUCCEEDED, nil, nil); err != nil {
		e.Logger.Error("failed to report handled command", zap.Error(err))
	}
	return as, true
}

func (e *Executor) reportApproved(approver string) {
	accounts, err := e.getMentionedAccounts(model.NotificationEventType_EVENT_DEPLOYMENT_APPROVED)
	if err != nil {
		e.Logger.Error("failed to get the list of accounts", zap.Error(err))
	}

	e.Notifier.Notify(model.NotificationEvent{
		Type: model.NotificationEventType_EVENT_DEPLOYMENT_APPROVED,
		Metadata: &model.NotificationEventDeploymentApproved{
			Deployment:        e.Deployment,
			EnvName:           e.EnvName,
			Approver:          approver,
			MentionedAccounts: accounts,
		},
	})
}

func (e *Executor) reportRequiringApproval() {
	accounts, err := e.getMentionedAccounts(model.NotificationEventType_EVENT_DEPLOYMENT_WAIT_APPROVAL)
	if err != nil {
		e.Logger.Error("failed to get the list of accounts", zap.Error(err))
	}

	e.Notifier.Notify(model.NotificationEvent{
		Type: model.NotificationEventType_EVENT_DEPLOYMENT_WAIT_APPROVAL,
		Metadata: &model.NotificationEventDeploymentWaitApproval{
			Deployment:        e.Deployment,
			EnvName:           e.EnvName,
			MentionedAccounts: accounts,
		},
	})
}

func (e *Executor) getMentionedAccounts(event model.NotificationEventType) ([]string, error) {
	n, ok := e.MetadataStore.Shared().Get(model.MetadataKeyDeploymentNotification)
	if !ok {
		return []string{}, nil
	}

	var notification config.DeploymentNotification
	if err := json.Unmarshal([]byte(n), &notification); err != nil {
		return nil, fmt.Errorf("could not extract mentions config: %w", err)
	}

	return notification.FindSlackAccounts(event), nil
}

func (e *Executor) validateApproverNum(ctx context.Context, approver string) (string, bool) {
	n, ok := e.MetadataStore.Stage(e.Stage.Id).Get(minApproverNum)
	if !ok {
		return "", false
	}
	num, err := strconv.Atoi(n)
	if err != nil {
		e.LogPersister.Errorf("%s could not be converted to integer: %v", num, err)
		return "", false
	}
	if num <= 1 {
		return approver, true
	}
	as, ok := e.MetadataStore.Stage(e.Stage.Id).Get(approversKey)
	if !ok {
		e.LogPersister.Infof("%d more approvals are needed", num-1)
		return approver, false
	}
	approvers := strings.Split(as, ",")
	for _, a := range approvers {
		if a == approver {
			e.LogPersister.Errorf("%s has already approved. %d more approvals are needed", num-1)
			return "", false
		}
	}
	if d := num - len(approvers) - 1; d > 0 {
		e.LogPersister.Infof("%d more approvals are needed", d)
		return as + "," + approver, false
	}
	return as + "," + approver, true
}
