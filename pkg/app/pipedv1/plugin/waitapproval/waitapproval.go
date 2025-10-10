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
	"iter"
	"strings"
	"time"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"go.uber.org/zap"
)

// executeWaitApproval waits for approvals.
func (p *plugin) executeWaitApproval(ctx context.Context, in *sdk.ExecuteStageInput[struct{}]) sdk.StageStatus {
	lp, err := in.Client.StageLogPersister()
	if err != nil {
		in.Logger.Error("No stage log persister available", zap.Error(err))
		return sdk.StageStatusFailure
	}
	opts, err := decode(in.Request.StageConfig)
	if err != nil {
		lp.Errorf("Failed to decode stage config: %v", err)
		return sdk.StageStatusFailure
	}
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	lp.Infof("Waiting for approval from at least %d user(s)...", opts.MinApproverNum)
	for {
		select {
		case <-ticker.C: // on ticker interval
			if approved := p.checkApproval(ctx, opts.MinApproverNum, lp, in.Client); approved {
				return sdk.StageStatusSuccess
			}

		case <-ctx.Done(): // on cancelled
			lp.Info("Wait approval cancelled")
			return sdk.StageStatusFailure
		}
	}
}

// checkApproval checks if there are enough approval commands.
func (p *plugin) checkApproval(ctx context.Context, minApproverNum int, lp sdk.StageLogPersister, client StageClient) bool {
	existingApprovedUsers, err := p.getApprovedUsers(ctx, client)
	if err != nil {
		lp.Errorf("Failed to get approved users: %v", err)
		return false
	}
	// approvedUsersMap contains previously approved users and newly approved users.
	approvedUsersMap := make(map[string]bool)
	for _, user := range existingApprovedUsers {
		approvedUsersMap[user] = true
	}
	for cmd, err := range client.ListStageCommands(ctx, sdk.CommandTypeApproveStage) {
		if err != nil {
			lp.Errorf("Failed to list stage commands: %v", err)
			return false
		}
		if approvedUsersMap[cmd.Commander] {
			continue
		}
		lp.Infof("Got approval from %q", cmd.Commander)
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
		if err := client.PutStageMetadata(ctx, sdk.MetadataKeyStageApprovedUsers, aus); err != nil {
			lp.Errorf("failed to save approved users: %v", err)
		}
		if err := client.PutStageMetadata(ctx, sdk.MetadataKeyStageDisplay, displayMsg); err != nil {
			lp.Errorf("failed to save display message: %v", err)
		}
		if remain := minApproverNum - len(approvedUsers); remain > 0 {
			lp.Infof("Waiting for %d more approver(s)...", remain)
			return false
		}
		lp.Infof("Received all needed approvals")
		lp.Infof("This stage has been approved by %d user(s) (%s)", minApproverNum, aus)
		return true
	}
	return false
}

// getApprovedUsers gets the list of approved users.
func (p *plugin) getApprovedUsers(ctx context.Context, client StageClient) ([]string, error) {
	val, exists, err := client.GetStageMetadata(ctx, sdk.MetadataKeyStageApprovedUsers)
	if err != nil {
		return nil, err
	}
	if val == "" || !exists {
		return []string{}, nil
	}
	return strings.Split(val, ", "), nil
}

type StageClient interface {
	GetStageMetadata(ctx context.Context, key string) (string, bool, error)
	PutStageMetadata(ctx context.Context, key string, value string) error
	ListStageCommands(ctx context.Context, commandTypes ...sdk.CommandType) iter.Seq2[*sdk.StageCommand, error]
}
