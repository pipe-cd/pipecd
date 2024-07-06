// Copyright 2024 The PipeCD Authors.
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
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
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

	e.reportRequiringApproval()

	num := e.StageConfig.WaitApprovalStageOptions.MinApproverNum
	e.LogPersister.Infof("Waiting for approval from at least %d user(s)...", num)
	for {
		select {
		case <-ticker.C:
			if e.checkApproval(ctx, num) {
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

func (e *Executor) checkApproval(ctx context.Context, num int) bool {
	var approveCmd *model.ReportableCommand
	commands := e.CommandLister.ListCommands()

	for i, cmd := range commands {
		if cmd.GetApproveStage() != nil {
			approveCmd = &commands[i]
			break
		}
	}
	if approveCmd == nil {
		return false
	}

	reached := e.validateApproverNum(ctx, approveCmd.Commander, num)
	if err := approveCmd.Report(ctx, model.CommandStatus_COMMAND_SUCCEEDED, nil, nil); err != nil {
		e.Logger.Error("failed to report handled command", zap.Error(err))
	}
	return reached
}

func (e *Executor) reportApproved(approver string) {
	users, err := e.getMentionedUsers(model.NotificationEventType_EVENT_DEPLOYMENT_APPROVED)
	if err != nil {
		e.Logger.Error("failed to get the list of users", zap.Error(err))
	}

	groups, err := e.getMentionedGroups(model.NotificationEventType_EVENT_DEPLOYMENT_APPROVED)
	if err != nil {
		e.Logger.Error("failed to get the list of accounts", zap.Error(err))
	}

	e.Notifier.Notify(model.NotificationEvent{
		Type: model.NotificationEventType_EVENT_DEPLOYMENT_APPROVED,
		Metadata: &model.NotificationEventDeploymentApproved{
			Deployment:        e.Deployment,
			Approver:          approver,
			MentionedAccounts: users,
			MentionedGroups:   groups,
		},
	})
}

func (e *Executor) reportRequiringApproval() {
	users, err := e.getMentionedUsers(model.NotificationEventType_EVENT_DEPLOYMENT_WAIT_APPROVAL)
	if err != nil {
		e.Logger.Error("failed to get the list of users", zap.Error(err))
	}

	groups, err := e.getMentionedGroups(model.NotificationEventType_EVENT_DEPLOYMENT_WAIT_APPROVAL)
	if err != nil {
		e.Logger.Error("failed to get the list of groups", zap.Error(err))
	}

	e.Notifier.Notify(model.NotificationEvent{
		Type: model.NotificationEventType_EVENT_DEPLOYMENT_WAIT_APPROVAL,
		Metadata: &model.NotificationEventDeploymentWaitApproval{
			Deployment:        e.Deployment,
			MentionedAccounts: users,
			MentionedGroups:   groups,
		},
	})
}

func (e *Executor) getMentionedUsers(event model.NotificationEventType) ([]string, error) {
	n, ok := e.MetadataStore.Shared().Get(model.MetadataKeyDeploymentNotification)
	if !ok {
		return []string{}, nil
	}

	var notification config.DeploymentNotification
	if err := json.Unmarshal([]byte(n), &notification); err != nil {
		return nil, fmt.Errorf("could not extract mentions users config: %w", err)
	}

	return notification.FindSlackUsers(event), nil
}

func (e *Executor) getMentionedGroups(event model.NotificationEventType) ([]string, error) {
	n, ok := e.MetadataStore.Shared().Get(model.MetadataKeyDeploymentNotification)
	if !ok {
		return []string{}, nil
	}

	var notification config.DeploymentNotification
	if err := json.Unmarshal([]byte(n), &notification); err != nil {
		return nil, fmt.Errorf("could not extract mentions groups config: %w", err)
	}

	return notification.FindSlackGroups(event), nil
}

// validateApproverNum checks if number of approves is valid.
func (e *Executor) validateApproverNum(ctx context.Context, approver string, minApproverNum int) bool {
	if minApproverNum == 1 {
		if err := e.MetadataStore.Stage(e.Stage.Id).Put(ctx, approvedByKey, approver); err != nil {
			e.LogPersister.Errorf("Unable to save approver information to deployment, %v", err)
		}
		e.LogPersister.Infof("Got approval from %q", approver)
		e.reportApproved(approver)
		e.LogPersister.Infof("This stage has been approved by %d user (%s)", minApproverNum, approver)
		return true
	}

	const delimiter = ", "
	as, _ := e.MetadataStore.Stage(e.Stage.Id).Get(approvedByKey)
	var approvedUsers []string
	if as != "" {
		approvedUsers = strings.Split(as, delimiter)
	}

	for _, u := range approvedUsers {
		if u == approver {
			e.LogPersister.Infof("Approval from the same user (%s) will not be counted", approver)
			return false
		}
	}
	e.LogPersister.Infof("Got approval from %q", approver)
	approvedUsers = append(approvedUsers, approver)
	aus := strings.Join(approvedUsers, delimiter)

	if err := e.MetadataStore.Stage(e.Stage.Id).Put(ctx, approvedByKey, aus); err != nil {
		e.LogPersister.Errorf("Unable to save approver information to deployment, %v", err)
	}
	if remain := minApproverNum - len(approvedUsers); remain > 0 {
		e.LogPersister.Infof("Waiting for %d other approvers...", remain)
		return false
	}
	e.reportApproved(aus)
	e.LogPersister.Info("Received all needed approvals")
	e.LogPersister.Infof("This stage has been approved by %d users (%s)", minApproverNum, aus)
	return true
}
