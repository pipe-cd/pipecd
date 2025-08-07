// Copyright 2025 The PipeCD Authors.
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

package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/model"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

// executeWaitApproval waits for approvals.
func (p *plugin) executeWaitApproval(ctx context.Context, in *sdk.ExecuteStageInput[struct{}]) sdk.StageStatus {
	opts, err := decode(in.Request.StageConfig)
	if err != nil {
		in.Client.LogPersister().Errorf("Failed to decode stage config: %v", err)
		return sdk.StageStatusFailure
	}
	in.Client.LogPersister().Infof("Waiting for approval from at least %d user(s)...", opts.MinApproverNum)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C: // on ticker interval
			if approved, approvers := p.checkApproval(ctx, in, opts.MinApproverNum); approved {
				in.Client.LogPersister().Infof("This stage has been approved by %d user(s) (%s)", opts.MinApproverNum, approvers)
				return sdk.StageStatusSuccess
			}

		case <-ctx.Done(): // on cancelled
			in.Client.LogPersister().Info("Wait approval cancelled")
			return sdk.StageStatusFailure
		}
	}
}

// checkApproval checks if there are enough approval commands.
func (p *plugin) checkApproval(ctx context.Context, in *sdk.ExecuteStageInput[struct{}], minApproverNum int) (bool, string) {
	existingApprovedUsers := p.getApprovedUsers(ctx, in)
	approvedUsersMap := make(map[string]bool)
	for _, user := range existingApprovedUsers {
		approvedUsersMap[user] = true
	}
	cmds := in.Client.ListStageCommands(ctx, model.Command_APPROVE_STAGE)
	for cmd, err := range cmds {
		if err != nil {
			in.Client.LogPersister().Errorf("Failed to list stage commands: %v", err)
			return false, ""
		}
		if approvedUsersMap[cmd.Commander] {
			in.Client.LogPersister().Infof("Approval from the same user (%s) will not be counted", cmd.Commander)
			continue
		}
		in.Client.LogPersister().Infof("Got approval from %q", cmd.Commander)
		approvedUsersMap[cmd.Commander] = true
		if len(approvedUsersMap) >= minApproverNum {
			break
		}
	}
	approvedUsers := make([]string, 0, len(approvedUsersMap))
	for user := range approvedUsersMap {
		approvedUsers = append(approvedUsers, user)
	}
	if len(approvedUsers) > 0 {
		aus := strings.Join(approvedUsers, ", ")
		displayMsg := fmt.Sprintf("Approved by: %s", aus)
		remain := minApproverNum - len(approvedUsers)
		if err := in.Client.PutStageMetadata(ctx, sdk.MetadataKeyStageApprovedUsers, aus); err != nil {
			in.Logger.Error("failed to save approver information", zap.Error(err))
		}
		if err := in.Client.PutStageMetadata(ctx, sdk.MetadataKeyStageDisplay, displayMsg); err != nil {
			in.Logger.Error("failed to save approver display information", zap.Error(err))
		}
		if remain > 0 {
			in.Client.LogPersister().Infof("Waiting for %d more approver(s)...", remain)
			return false, aus
		}
		in.Client.LogPersister().Info("Received all needed approvals")
		return true, aus
	}

	return false, ""
}

// getApprovedUsers gets the list of approved users.
func (p *plugin) getApprovedUsers(ctx context.Context, in *sdk.ExecuteStageInput[struct{}]) []string {
	val, err := in.Client.GetStageMetadata(ctx, sdk.MetadataKeyStageApprovedUsers)
	if err != nil || val == "" {
		return []string{}
	}
	return strings.Split(val, ", ")
}
