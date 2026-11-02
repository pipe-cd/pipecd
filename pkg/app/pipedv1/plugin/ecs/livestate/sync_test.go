// Copyright 2026 The PipeCD Authors.
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

package livestate

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/provider"
)

func ecsControllerService(commit string) *types.Service {
	svc := &types.Service{
		ServiceArn:   aws.String("arn:service"),
		ServiceName:  aws.String("svc"),
		Status:       aws.String("ACTIVE"),
		DesiredCount: 2,
		RunningCount: 2,
	}
	if commit != "" {
		svc.Tags = []types.Tag{{Key: aws.String(provider.LabelCommitHash), Value: aws.String(commit)}}
	}
	return svc
}

func TestComputeSyncState_ECSController(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		svc             *types.Service
		commitHash      string
		wantStatus      sdk.ApplicationSyncStatus
		wantShortReason string
		wantReason      string
	}{
		{
			name:       "same commit",
			svc:        ecsControllerService("abc"),
			commitHash: "abc",
			wantStatus: sdk.ApplicationSyncStateSynced,
		},
		{
			name:            "different commit",
			svc:             ecsControllerService("old"),
			commitHash:      "abc",
			wantStatus:      sdk.ApplicationSyncStateOutOfSync,
			wantShortReason: "Deployed commit does not match current commit",
			wantReason:      "deployed: old, expected: abc",
		},
		{
			name:            "never deployed successfully",
			svc:             ecsControllerService(""),
			commitHash:      "abc",
			wantStatus:      sdk.ApplicationSyncStateOutOfSync,
			wantShortReason: "No successful deployment found for the service",
		},
		{
			name:       "no commit to compare",
			svc:        ecsControllerService(""),
			commitHash: "",
			wantStatus: sdk.ApplicationSyncStateSynced,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := computeSyncState(&queryResourcesResult{Service: tc.svc}, types.Service{ServiceName: aws.String("svc")}, tc.commitHash)
			assert.Equal(t, tc.wantStatus, got.Status)
			assert.Equal(t, tc.wantShortReason, got.ShortReason)
			assert.Equal(t, tc.wantReason, got.Reason)
		})
	}
}

func TestComputeSyncState_ExternalControllerUnchanged(t *testing.T) {
	t.Parallel()

	svc := ecsControllerService("abc")
	svc.DeploymentController = &types.DeploymentController{Type: types.DeploymentControllerTypeExternal}

	got := computeSyncState(&queryResourcesResult{Service: svc}, types.Service{ServiceName: aws.String("svc")}, "abc")
	assert.Equal(t, sdk.ApplicationSyncStateOutOfSync, got.Status)
	assert.Equal(t, "No PRIMARY task set found", got.ShortReason)
}
