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

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

// executeWaitApproval waits for approvals.
func (p *plugin) executeWaitApproval(ctx context.Context, in *sdk.ExecuteStageInput[struct{}]) sdk.StageStatus {
	opts, err := decode(in.Request.StageConfig)
	if err != nil {
		in.Client.LogPersister().Errorf("Failed to decode stage config: %v", err)
		return sdk.StageStatusFailure
	}
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	in.Client.LogPersister().Infof("Waiting for approval from at least %d user(s)...", opts.MinApproverNum)
	for {
		select {
		case <-ticker.C: // on ticker interval
			if approved, approvers := p.checkApproval(ctx, opts.MinApproverNum, in.Client, in.Logger); approved {
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
func (p *plugin) checkApproval(ctx context.Context, minApproverNum int, client *sdk.Client, logger *zap.Logger) (bool, string) {
	existingApprovedUsers := p.getApprovedUsers(ctx, client)
	approvedUsersMap := make(map[string]bool)
	for _, user := range existingApprovedUsers {
		approvedUsersMap[user] = true
	}
	for cmd, err := range client.ListStageCommands(ctx, sdk.CommandTypeApproveStage) {
		if err != nil {
			logger.Error("Failed to list stage commands", zap.Error(err))
			return false, ""
		}
		if approvedUsersMap[cmd.Commander] {
			client.LogPersister().Infof("Approval from the same user (%s) will not be counted", cmd.Commander)
			continue
		}
		client.LogPersister().Infof("Got approval from %q", cmd.Commander)
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
		if err := client.PutStageMetadata(ctx, sdk.MetadataKeyStageApprovedUsers, aus); err != nil {
			logger.Error("failed to save approver information", zap.Error(err))
		}
		if err := client.PutStageMetadata(ctx, sdk.MetadataKeyStageDisplay, displayMsg); err != nil {
			logger.Error("failed to save approver display information", zap.Error(err))
		}
		if remain > 0 {
			client.LogPersister().Infof("Waiting for %d more approver(s)...", remain)
			return false, aus
		}
		client.LogPersister().Infof("Received all needed approvals")
		return true, aus
	}

	return false, ""
}

// getApprovedUsers gets the list of approved users.
func (p *plugin) getApprovedUsers(ctx context.Context, client *sdk.Client) []string {
	val, exists, err := client.GetStageMetadata(ctx, sdk.MetadataKeyStageApprovedUsers)
	if err != nil || val == "" || !exists {
		return []string{}
	}
	return strings.Split(val, ", ")
}
