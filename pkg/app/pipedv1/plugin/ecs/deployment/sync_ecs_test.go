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

package deployment

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/provider"
)

func ecsServiceDef(controller *types.DeploymentController) types.Service {
	return types.Service{
		ClusterArn:           aws.String("cluster"),
		ServiceName:          aws.String("svc"),
		ServiceArn:           aws.String("svc-arn"),
		DeploymentController: controller,
		Tags: []types.Tag{
			{Key: aws.String("a"), Value: aws.String("b")},
			{Key: aws.String(provider.LabelCommitHash), Value: aws.String("new")},
		},
	}
}

func TestSyncECSController(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		waitErr    error
		wantErr    bool
		wantTag    bool
		wantPrev   *string
		prevDeploy *types.ServiceDeploymentBrief
	}{
		{name: "success tags commit after completion", wantTag: true},
		{
			name:       "waits relative to previous deployment",
			prevDeploy: &types.ServiceDeploymentBrief{ServiceDeploymentArn: aws.String("old")},
			wantPrev:   aws.String("old"),
			wantTag:    true,
		},
		{name: "failed deployment does not tag", waitErr: errors.New("rolled back"), wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var tagged []types.Tag
			var gotPrev *string
			var created *types.Service
			c := happyPathClient(&types.TaskDefinition{TaskDefinitionArn: aws.String("td-new")}, &types.Service{ServiceName: aws.String("svc")}, nil, nil)
			c.ServiceExistsFunc = func(context.Context, string, string) (bool, error) { return false, nil }
			c.CreateServiceFunc = func(_ context.Context, s types.Service) (*types.Service, error) {
				created = &s
				return &s, nil
			}
			c.GetLatestServiceDeploymentFunc = func(context.Context, string, string) (*types.ServiceDeploymentBrief, error) {
				return tc.prevDeploy, nil
			}
			c.WaitServiceDeploymentCompleteFunc = func(_ context.Context, _, _ string, prev *string) error {
				gotPrev = prev
				return tc.waitErr
			}
			c.TagResourceFunc = func(_ context.Context, _ string, tags []types.Tag) error {
				tagged = append(tagged, tags...)
				return nil
			}

			err := syncECSController(context.Background(), &fakeLogPersister{}, c, types.TaskDefinition{}, ecsServiceDef(nil))
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.NotNil(t, created)
			assert.Equal(t, "td-new", *created.TaskDefinition)
			for _, tag := range created.Tags {
				assert.NotEqual(t, provider.LabelCommitHash, *tag.Key, "commit tag must not be sent before success")
			}
			assert.Equal(t, tc.wantPrev, gotPrev)

			var hasCommit bool
			for _, tag := range tagged {
				if *tag.Key == provider.LabelCommitHash {
					hasCommit = true
				}
			}
			assert.Equal(t, tc.wantTag, hasCommit)
		})
	}
}

func TestResolveController(t *testing.T) {
	t.Parallel()

	live := func(c types.DeploymentControllerType) *types.Service {
		return &types.Service{DeploymentController: &types.DeploymentController{Type: c}}
	}
	tests := []struct {
		name    string
		def     *types.DeploymentController
		exists  bool
		live    *types.Service
		want    types.DeploymentControllerType
		wantErr bool
	}{
		{name: "nil means ECS for new service", want: types.DeploymentControllerTypeEcs},
		{name: "explicit EXTERNAL", def: &types.DeploymentController{Type: types.DeploymentControllerTypeExternal}, want: types.DeploymentControllerTypeExternal},
		{name: "CODE_DEPLOY rejected", def: &types.DeploymentController{Type: types.DeploymentControllerTypeCodeDeploy}, wantErr: true},
		{name: "unknown rejected", def: &types.DeploymentController{Type: "FOO"}, wantErr: true},
		{name: "live matches", exists: true, live: live(types.DeploymentControllerTypeEcs), want: types.DeploymentControllerTypeEcs},
		{name: "live mismatch", exists: true, live: live(types.DeploymentControllerTypeExternal), wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := &mockECSClient{
				ServiceExistsFunc: func(context.Context, string, string) (bool, error) { return tc.exists, nil },
				DescribeServiceFunc: func(context.Context, types.Service) (*types.Service, error) {
					return tc.live, nil
				},
			}
			got, err := resolveController(context.Background(), c, ecsServiceDef(tc.def))
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestRejectExternalOnlyStage(t *testing.T) {
	t.Parallel()
	assert.Error(t, rejectExternalOnlyStage(StageECSPrimaryRollout, ecsServiceDef(nil)))
	assert.NoError(t, rejectExternalOnlyStage(StageECSPrimaryRollout, ecsServiceDef(&types.DeploymentController{Type: types.DeploymentControllerTypeExternal})))
}
