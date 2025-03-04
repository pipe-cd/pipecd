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

package commandstore

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestListStageCommands(t *testing.T) {
	t.Parallel()

	store := store{
		stageApproveCommands: stageCommandMap{
			"deployment-1": {
				"stage-1": []*model.Command{
					{
						Id:           "command-1",
						DeploymentId: "deployment-1",
						StageId:      "stage-1",
						Type:         model.Command_APPROVE_STAGE,
						Commander:    "commander-1",
					},
					{
						Id:           "command-2",
						DeploymentId: "deployment-1",
						StageId:      "stage-1",
						Type:         model.Command_APPROVE_STAGE,
						Commander:    "commander-2",
					},
				},
			},
		},
		stageSkipCommands: stageCommandMap{
			"deployment-11": {
				"stage-11": []*model.Command{
					{
						Id:           "command-11",
						DeploymentId: "deployment-11",
						StageId:      "stage-11",
						Type:         model.Command_SKIP_STAGE,
					},
				},
			},
		},
		logger: zap.NewNop(),
	}

	testcases := []struct {
		name         string
		deploymentID string
		stageID      string
		commandType  model.Command_Type
		want         []*model.Command
		wantErr      error
	}{
		{
			name:         "valid arguments of Approve",
			deploymentID: "deployment-1",
			stageID:      "stage-1",
			commandType:  model.Command_APPROVE_STAGE,
			want: []*model.Command{
				{
					Id:           "command-1",
					DeploymentId: "deployment-1",
					StageId:      "stage-1",
					Type:         model.Command_APPROVE_STAGE,
					Commander:    "commander-1",
				},
				{
					Id:           "command-2",
					DeploymentId: "deployment-1",
					StageId:      "stage-1",
					Type:         model.Command_APPROVE_STAGE,
					Commander:    "commander-2",
				},
			},
			wantErr: nil,
		},
		{
			name:         "valid arguments of Skip",
			deploymentID: "deployment-11",
			stageID:      "stage-11",
			commandType:  model.Command_SKIP_STAGE,
			want: []*model.Command{
				{
					Id:           "command-11",
					DeploymentId: "deployment-11",
					StageId:      "stage-11",
					Type:         model.Command_SKIP_STAGE,
				},
			},
			wantErr: nil,
		},
		{
			name:         "stageID not exists",
			deploymentID: "deployment-1",
			stageID:      "stage-999",
			commandType:  model.Command_APPROVE_STAGE,
			want:         nil,
			wantErr:      nil,
		},
		{
			name:         "invalid commandType",
			deploymentID: "deployment-1",
			stageID:      "stage-1",
			commandType:  model.Command_CANCEL_DEPLOYMENT,
			want:         nil,
			wantErr:      fmt.Errorf("invalid command type: CANCEL_DEPLOYMENT"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := store.ListStageCommands(tc.deploymentID, tc.stageID, tc.commandType)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
