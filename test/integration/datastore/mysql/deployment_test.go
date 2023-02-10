// Copyright 2023 The PipeCD Authors.
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

package mysql

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestGetDeployment(t *testing.T) {
	col := &collection{kind: "Deployment"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeDeployment := &model.Deployment{
		Id:              "get-id",
		ApplicationId:   "application-id",
		ApplicationName: "application-name",
		PipedId:         "piped-id",
		ProjectId:       "project-id",
		Kind:            model.ApplicationKind_KUBERNETES,
		StatusReason:    "status-desc",
		Status:          model.DeploymentStatus_DEPLOYMENT_PENDING,
		Stages: []*model.PipelineStage{
			{
				Id:       "stage-id1",
				Name:     "stage1",
				Desc:     "desc1",
				Index:    1,
				Status:   model.StageStatus_STAGE_SUCCESS,
				Metadata: map[string]string{"meta": "value"},
			},
			{
				Id:       "stage-id2",
				Name:     "stage2",
				Desc:     "desc2",
				Index:    2,
				Status:   model.StageStatus_STAGE_RUNNING,
				Metadata: map[string]string{"meta": "value"},
			},
		},
		CompletedAt: 1,
		CreatedAt:   1,
		UpdatedAt:   1,
	}
	err := client.Create(ctx, col, "get-id", fakeDeployment)
	require.NoError(t, err)

	testcases := []struct {
		name    string
		id      string
		want    *model.Deployment
		wantErr error
	}{
		{
			name:    "entity found",
			id:      "get-id",
			want:    fakeDeployment,
			wantErr: nil,
		},
		{
			name:    "not found",
			id:      "id-wrong",
			want:    &model.Deployment{},
			wantErr: datastore.ErrNotFound,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := &model.Deployment{}
			err := client.Get(ctx, col, tc.id, got)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFindDeployment(t *testing.T) {
	col := &collection{kind: "Deployment"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeDeployment := &model.Deployment{
		Id:              "find-id-1",
		ApplicationId:   "application-id",
		ApplicationName: "application-name",
		PipedId:         "piped-id",
		ProjectId:       "find-project-id",
		Kind:            model.ApplicationKind_KUBERNETES,
		StatusReason:    "status-desc",
		Status:          model.DeploymentStatus_DEPLOYMENT_PENDING,
		Stages: []*model.PipelineStage{
			{
				Id:       "stage-id1",
				Name:     "stage1",
				Desc:     "desc1",
				Index:    1,
				Status:   model.StageStatus_STAGE_SUCCESS,
				Metadata: map[string]string{"meta": "value"},
			},
			{
				Id:       "stage-id2",
				Name:     "stage2",
				Desc:     "desc2",
				Index:    2,
				Status:   model.StageStatus_STAGE_RUNNING,
				Metadata: map[string]string{"meta": "value"},
			},
		},
		CompletedAt: 1,
		CreatedAt:   1,
		UpdatedAt:   1,
	}
	fakeDeployment2 := &model.Deployment{
		Id:              "find-id-2",
		ApplicationId:   "application-id",
		ApplicationName: "application-name",
		PipedId:         "piped-id",
		ProjectId:       "find-project-id",
		Kind:            model.ApplicationKind_KUBERNETES,
		StatusReason:    "status-desc",
		Status:          model.DeploymentStatus_DEPLOYMENT_PENDING,
		Stages: []*model.PipelineStage{
			{
				Id:       "stage-id1",
				Name:     "stage1",
				Desc:     "desc1",
				Index:    1,
				Status:   model.StageStatus_STAGE_SUCCESS,
				Metadata: map[string]string{"meta": "value"},
			},
			{
				Id:       "stage-id2",
				Name:     "stage2",
				Desc:     "desc2",
				Index:    2,
				Status:   model.StageStatus_STAGE_RUNNING,
				Metadata: map[string]string{"meta": "value"},
			},
		},
		CompletedAt: 1,
		CreatedAt:   1,
		UpdatedAt:   1,
	}
	err := client.Create(ctx, col, "find-id-1", fakeDeployment)
	require.NoError(t, err)
	err = client.Create(ctx, col, "find-id-2", fakeDeployment2)
	require.NoError(t, err)

	testcases := []struct {
		name    string
		opts    datastore.ListOptions
		want    []*model.Deployment
		wantErr bool
	}{
		{
			name: "fetch by project_id",
			opts: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "ProjectId",
						Operator: datastore.OperatorEqual,
						Value:    "find-project-id",
					},
				},
			},
			want: []*model.Deployment{
				fakeDeployment,
				fakeDeployment2,
			},
			wantErr: false,
		},
		{
			name: "only cursor given",
			opts: datastore.ListOptions{
				Cursor: "cursor",
			},
			want:    []*model.Deployment{},
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			it, err := client.Find(ctx, col, tc.opts)
			assert.Equal(t, tc.wantErr, err != nil)
			got, err := listDeployments(it)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func listDeployments(it datastore.Iterator) ([]*model.Deployment, error) {
	ret := make([]*model.Deployment, 0)
	if it == nil {
		return ret, nil
	}
	for {
		var v model.Deployment
		err := it.Next(&v)
		if errors.Is(err, datastore.ErrIteratorDone) {
			break
		}
		if err != nil {
			return nil, err
		}
		ret = append(ret, &v)
	}
	return ret, nil
}

func TestCreateDeployment(t *testing.T) {
	col := &collection{kind: "Deployment"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeDeployment := &model.Deployment{
		Id:              "create-id",
		ApplicationId:   "application-id",
		ApplicationName: "application-name",
		PipedId:         "piped-id",
		ProjectId:       "project-id",
		Kind:            model.ApplicationKind_KUBERNETES,
		StatusReason:    "status-desc",
		Status:          model.DeploymentStatus_DEPLOYMENT_PENDING,
		Stages: []*model.PipelineStage{
			{
				Id:       "stage-id1",
				Name:     "stage1",
				Desc:     "desc1",
				Index:    1,
				Status:   model.StageStatus_STAGE_SUCCESS,
				Metadata: map[string]string{"meta": "value"},
			},
			{
				Id:       "stage-id2",
				Name:     "stage2",
				Desc:     "desc2",
				Index:    2,
				Status:   model.StageStatus_STAGE_RUNNING,
				Metadata: map[string]string{"meta": "value"},
			},
		},
		CompletedAt: 1,
		CreatedAt:   1,
		UpdatedAt:   1,
	}
	err := client.Create(ctx, col, "create-id", fakeDeployment)
	require.NoError(t, err)

	testcases := []struct {
		name    string
		id      string
		wantErr error
	}{
		{
			name:    "already exists",
			id:      "create-id",
			wantErr: datastore.ErrAlreadyExists,
		},
		{
			name:    "successful create",
			id:      "id-new",
			wantErr: nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := client.Create(ctx, col, tc.id, fakeDeployment)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestUpdateDeployment(t *testing.T) {
	col := &collection{
		kind: "Deployment",
		factory: func() interface{} {
			return &model.Deployment{}
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeDeployment := &model.Deployment{
		Id:              "update-id",
		ApplicationId:   "application-id",
		ApplicationName: "application-name",
		PipedId:         "piped-id",
		ProjectId:       "project-id",
		Kind:            model.ApplicationKind_KUBERNETES,
		StatusReason:    "status-desc",
		Status:          model.DeploymentStatus_DEPLOYMENT_PENDING,
		Stages: []*model.PipelineStage{
			{
				Id:       "stage-id1",
				Name:     "stage1",
				Desc:     "desc1",
				Index:    1,
				Status:   model.StageStatus_STAGE_SUCCESS,
				Metadata: map[string]string{"meta": "value"},
			},
			{
				Id:       "stage-id2",
				Name:     "stage2",
				Desc:     "desc2",
				Index:    2,
				Status:   model.StageStatus_STAGE_RUNNING,
				Metadata: map[string]string{"meta": "value"},
			},
		},
		CompletedAt: 1,
		CreatedAt:   1,
		UpdatedAt:   1,
	}
	err := client.Create(ctx, col, "update-id", fakeDeployment)
	require.NoError(t, err)

	testcases := []struct {
		name    string
		id      string
		updater func(interface{}) error
		wantErr error
	}{
		{
			name:    "not found",
			id:      "id-wrong",
			wantErr: datastore.ErrNotFound,
		},
		{
			name: "unable to update",
			id:   "update-id",
			updater: func(interface{}) error {
				return fmt.Errorf("error")
			},
			wantErr: fmt.Errorf("error"),
		},
		{
			name: "successful update",
			id:   "update-id",
			updater: func(e interface{}) error {
				v := e.(*model.Deployment)
				v.Status = model.DeploymentStatus_DEPLOYMENT_SUCCESS
				return nil
			},
			wantErr: nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := client.Update(ctx, col, tc.id, tc.updater)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
