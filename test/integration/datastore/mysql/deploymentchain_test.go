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

func TestGetDeploymentChain(t *testing.T) {
	col := &collection{kind: "DeploymentChain"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeDeploymentChain := &model.DeploymentChain{
		Id:        "get-id",
		ProjectId: "project-id",
		Status:    model.ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
		Blocks: []*model.ChainBlock{
			{
				Nodes: []*model.ChainNode{
					{
						DeploymentRef: &model.ChainDeploymentRef{
							DeploymentId: "deploy-0",
							Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
						},
					},
				},
			},
		},
		CompletedAt: 1,
		CreatedAt:   1,
		UpdatedAt:   1,
	}
	err := client.Create(ctx, col, "get-id", fakeDeploymentChain)
	require.NoError(t, err)

	testcases := []struct {
		name    string
		id      string
		want    *model.DeploymentChain
		wantErr error
	}{
		{
			name:    "entity found",
			id:      "get-id",
			want:    fakeDeploymentChain,
			wantErr: nil,
		},
		{
			name:    "not found",
			id:      "id-wrong",
			want:    &model.DeploymentChain{},
			wantErr: datastore.ErrNotFound,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := &model.DeploymentChain{}
			err := client.Get(ctx, col, tc.id, got)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFindDeploymentChain(t *testing.T) {
	col := &collection{kind: "DeploymentChain"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeDeploymentChain := &model.DeploymentChain{
		Id:        "find-id-1",
		ProjectId: "find-project-id",
		Status:    model.ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
		Blocks: []*model.ChainBlock{
			{
				Nodes: []*model.ChainNode{
					{
						DeploymentRef: &model.ChainDeploymentRef{
							DeploymentId: "deploy-0",
							Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
						},
					},
				},
			},
		},
		CompletedAt: 1,
		CreatedAt:   1,
		UpdatedAt:   1,
	}
	fakeDeploymentChain2 := &model.DeploymentChain{
		Id:        "find-id-2",
		ProjectId: "find-project-id",
		Status:    model.ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
		Blocks: []*model.ChainBlock{
			{
				Nodes: []*model.ChainNode{
					{
						DeploymentRef: &model.ChainDeploymentRef{
							DeploymentId: "deploy-0",
							Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
						},
					},
				},
			},
		},
		CompletedAt: 1,
		CreatedAt:   1,
		UpdatedAt:   1,
	}
	err := client.Create(ctx, col, "find-id-1", fakeDeploymentChain)
	require.NoError(t, err)
	err = client.Create(ctx, col, "find-id-2", fakeDeploymentChain2)
	require.NoError(t, err)

	testcases := []struct {
		name    string
		opts    datastore.ListOptions
		want    []*model.DeploymentChain
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
			want: []*model.DeploymentChain{
				fakeDeploymentChain,
				fakeDeploymentChain2,
			},
			wantErr: false,
		},
		{
			name: "only cursor given",
			opts: datastore.ListOptions{
				Cursor: "cursor",
			},
			want:    []*model.DeploymentChain{},
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			it, err := client.Find(ctx, col, tc.opts)
			assert.Equal(t, tc.wantErr, err != nil)
			got, err := listDeploymentChains(it)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func listDeploymentChains(it datastore.Iterator) ([]*model.DeploymentChain, error) {
	ret := make([]*model.DeploymentChain, 0)
	if it == nil {
		return ret, nil
	}
	for {
		var v model.DeploymentChain
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

func TestCreateDeploymentChain(t *testing.T) {
	col := &collection{kind: "DeploymentChain"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeDeploymentChain := &model.DeploymentChain{
		Id:        "create-id",
		ProjectId: "project-id",
		Status:    model.ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
		Blocks: []*model.ChainBlock{
			{
				Nodes: []*model.ChainNode{
					{
						DeploymentRef: &model.ChainDeploymentRef{
							DeploymentId: "deploy-0",
							Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
						},
					},
				},
			},
		},
		CompletedAt: 1,
		CreatedAt:   1,
		UpdatedAt:   1,
	}

	err := client.Create(ctx, col, "create-id", fakeDeploymentChain)
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
			err := client.Create(ctx, col, tc.id, fakeDeploymentChain)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestUpdateDeploymentChain(t *testing.T) {
	col := &collection{
		kind: "DeploymentChain",
		factory: func() interface{} {
			return &model.DeploymentChain{}
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeDeploymentChain := &model.DeploymentChain{
		Id:        "update-id",
		ProjectId: "project-id",
		Status:    model.ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
		Blocks: []*model.ChainBlock{
			{
				Nodes: []*model.ChainNode{
					{
						DeploymentRef: &model.ChainDeploymentRef{
							DeploymentId: "deploy-0",
							Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
						},
					},
				},
			},
		},
		CompletedAt: 1,
		CreatedAt:   1,
		UpdatedAt:   1,
	}
	err := client.Create(ctx, col, "update-id", fakeDeploymentChain)
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
				v := e.(*model.DeploymentChain)
				v.Status = model.ChainStatus_DEPLOYMENT_CHAIN_SUCCESS
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
