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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestListStageCommands(t *testing.T) {
	t.Parallel()

	store := store{
		stageCommands: stageCommandsMap{
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
					{
						Id:           "command-3",
						DeploymentId: "deployment-1",
						StageId:      "stage-1",
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
		want         []*model.Command
		wantErr      error
	}{
		{
			name:         "valid arguments",
			deploymentID: "deployment-1",
			stageID:      "stage-1",
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
				{
					Id:           "command-3",
					DeploymentId: "deployment-1",
					StageId:      "stage-1",
					Type:         model.Command_SKIP_STAGE,
				},
			},
			wantErr: nil,
		},
		{
			name:         "deploymentID not exist",
			deploymentID: "xxx",
			stageID:      "stage-1",
			want:         nil,
			wantErr:      nil,
		},
		{
			name:         "stageID not exist",
			deploymentID: "deployment-1",
			stageID:      "stage-999",
			want:         nil,
			wantErr:      nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := store.ListStageCommands(tc.deploymentID, tc.stageID)
			assert.Equal(t, tc.want, got)
		})
	}
}
