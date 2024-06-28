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

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeploymentChainDesireStatus(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name                 string
		deploymentChain      DeploymentChain
		expectedDesireStatus ChainStatus
	}{
		{
			name: "reach SUCCESS state from RUNNING after all block finished successfully",
			deploymentChain: DeploymentChain{
				Status: ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
				Blocks: []*ChainBlock{
					{
						Status: ChainBlockStatus_DEPLOYMENT_BLOCK_SUCCESS,
					},
					{
						Status: ChainBlockStatus_DEPLOYMENT_BLOCK_SUCCESS,
					},
				},
			},
			expectedDesireStatus: ChainStatus_DEPLOYMENT_CHAIN_SUCCESS,
		},
		{
			name: "finished chain keeps its finished state",
			deploymentChain: DeploymentChain{
				Status: ChainStatus_DEPLOYMENT_CHAIN_FAILURE,
				Blocks: []*ChainBlock{
					{
						Status: ChainBlockStatus_DEPLOYMENT_BLOCK_SUCCESS,
					},
					{
						Status: ChainBlockStatus_DEPLOYMENT_BLOCK_FAILURE,
					},
				},
			},
			expectedDesireStatus: ChainStatus_DEPLOYMENT_CHAIN_FAILURE,
		},
		{
			name: "one FAILURE block makes the whole chain as FAILURE",
			deploymentChain: DeploymentChain{
				Status: ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
				Blocks: []*ChainBlock{
					{
						Status: ChainBlockStatus_DEPLOYMENT_BLOCK_FAILURE,
					},
					{
						Status: ChainBlockStatus_DEPLOYMENT_BLOCK_CANCELLED,
					},
				},
			},
			expectedDesireStatus: ChainStatus_DEPLOYMENT_CHAIN_FAILURE,
		},
		{
			name: "one CANCELLED block makes the whole chain as CANCELLED",
			deploymentChain: DeploymentChain{
				Status: ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
				Blocks: []*ChainBlock{
					{
						Status: ChainBlockStatus_DEPLOYMENT_BLOCK_SUCCESS,
					},
					{
						Status: ChainBlockStatus_DEPLOYMENT_BLOCK_CANCELLED,
					},
				},
			},
			expectedDesireStatus: ChainStatus_DEPLOYMENT_CHAIN_CANCELLED,
		},
		{
			name: "one RUNNING block make not yet finished chain as RUNNING",
			deploymentChain: DeploymentChain{
				Status: ChainStatus_DEPLOYMENT_CHAIN_PENDING,
				Blocks: []*ChainBlock{
					{
						Status: ChainBlockStatus_DEPLOYMENT_BLOCK_RUNNING,
					},
					{
						Status: ChainBlockStatus_DEPLOYMENT_BLOCK_PENDING,
					},
				},
			},
			expectedDesireStatus: ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			desireStatus := tc.deploymentChain.DesiredStatus()
			assert.Equal(t, tc.expectedDesireStatus, desireStatus)
		})
	}
}
func TestListAllInChainApplicationDeploymentsMap(t *testing.T) {
	t.Parallel()
	dc := &DeploymentChain{
		Blocks: []*ChainBlock{
			{
				Nodes: []*ChainNode{
					{
						ApplicationRef: &ChainApplicationRef{ApplicationId: "app1"},
						DeploymentRef:  &ChainDeploymentRef{DeploymentId: "dep1"},
					},
					{
						ApplicationRef: &ChainApplicationRef{ApplicationId: "app2"},
						DeploymentRef:  nil,
					},
				},
			},
			{
				Nodes: []*ChainNode{
					{
						ApplicationRef: &ChainApplicationRef{ApplicationId: "app2"},
						DeploymentRef:  &ChainDeploymentRef{DeploymentId: "dep2"},
					},
					{
						ApplicationRef: &ChainApplicationRef{ApplicationId: "app3"},
						DeploymentRef:  &ChainDeploymentRef{DeploymentId: "dep3"},
					},
				},
			},
		},
	}

	want := map[string]*ChainDeploymentRef{
		"app1": {DeploymentId: "dep1"},
		"app2": {DeploymentId: "dep2"},
		"app3": {DeploymentId: "dep3"},
	}

	got := dc.ListAllInChainApplicationDeploymentsMap()

	assert.Equal(t, want, got)
}
func TestListAllInChainApplications(t *testing.T) {
	t.Parallel()
	dc := &DeploymentChain{
		Blocks: []*ChainBlock{
			{
				Nodes: []*ChainNode{
					{
						ApplicationRef: &ChainApplicationRef{ApplicationId: "app1"},
						DeploymentRef:  &ChainDeploymentRef{DeploymentId: "dep1"},
					},
					{
						ApplicationRef: &ChainApplicationRef{ApplicationId: "app2"},
						DeploymentRef:  nil,
					},
				},
			},
			{
				Nodes: []*ChainNode{
					{
						ApplicationRef: &ChainApplicationRef{ApplicationId: "app2"},
						DeploymentRef:  &ChainDeploymentRef{DeploymentId: "dep2"},
					},
					{
						ApplicationRef: &ChainApplicationRef{ApplicationId: "app3"},
						DeploymentRef:  &ChainDeploymentRef{DeploymentId: "dep3"},
					},
				},
			},
		},
	}

	want := []*ChainApplicationRef{
		{ApplicationId: "app1"},
		{ApplicationId: "app2"},
		{ApplicationId: "app2"},
		{ApplicationId: "app3"},
	}

	got := dc.ListAllInChainApplications()

	assert.Equal(t, want, got)
}

func TestIsCompleted(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status ChainBlockStatus
		want   bool
	}{
		{
			name:   "returns true for DEPLOYMENT_BLOCK_SUCCESS",
			status: ChainBlockStatus_DEPLOYMENT_BLOCK_SUCCESS,
			want:   true,
		},
		{
			name:   "returns true for DEPLOYMENT_BLOCK_FAILURE",
			status: ChainBlockStatus_DEPLOYMENT_BLOCK_FAILURE,
			want:   true,
		},
		{
			name:   "returns true for DEPLOYMENT_BLOCK_CANCELLED",
			status: ChainBlockStatus_DEPLOYMENT_BLOCK_CANCELLED,
			want:   true,
		},
		{
			name:   "returns false for DEPLOYMENT_BLOCK_PENDING",
			status: ChainBlockStatus_DEPLOYMENT_BLOCK_PENDING,
			want:   false,
		},
		{
			name:   "returns false for DEPLOYMENT_BLOCK_RUNNING",
			status: ChainBlockStatus_DEPLOYMENT_BLOCK_RUNNING,
			want:   false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			b := &ChainBlock{Status: tt.status}
			got := b.IsCompleted()
			assert.Equal(t, tt.want, got)
		})
	}
}
func TestDesiredStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		block          *ChainBlock
		wantDesired    ChainBlockStatus
		wantIsComplete bool
	}{
		{
			name: "returns DEPLOYMENT_BLOCK_SUCCESS for all successful deployments",
			block: &ChainBlock{
				Nodes: []*ChainNode{
					{DeploymentRef: &ChainDeploymentRef{Status: DeploymentStatus_DEPLOYMENT_SUCCESS}},
					{DeploymentRef: &ChainDeploymentRef{Status: DeploymentStatus_DEPLOYMENT_SUCCESS}},
				},
			},
			wantDesired: ChainBlockStatus_DEPLOYMENT_BLOCK_SUCCESS,
		},
		{
			name: "returns DEPLOYMENT_BLOCK_FAILURE for at least one failed deployment",
			block: &ChainBlock{
				Nodes: []*ChainNode{
					{DeploymentRef: &ChainDeploymentRef{Status: DeploymentStatus_DEPLOYMENT_SUCCESS}},
					{DeploymentRef: &ChainDeploymentRef{Status: DeploymentStatus_DEPLOYMENT_FAILURE}},
				},
			},
			wantDesired: ChainBlockStatus_DEPLOYMENT_BLOCK_FAILURE,
		},
		{
			name: "returns DEPLOYMENT_BLOCK_CANCELLED for at least one cancelled deployment",
			block: &ChainBlock{
				Nodes: []*ChainNode{
					{DeploymentRef: &ChainDeploymentRef{Status: DeploymentStatus_DEPLOYMENT_SUCCESS}},
					{DeploymentRef: &ChainDeploymentRef{Status: DeploymentStatus_DEPLOYMENT_CANCELLED}},
				},
			},
			wantDesired: ChainBlockStatus_DEPLOYMENT_BLOCK_CANCELLED,
		},
		{
			name: "returns DEPLOYMENT_BLOCK_RUNNING for at least one running deployment",
			block: &ChainBlock{
				Nodes: []*ChainNode{
					{DeploymentRef: &ChainDeploymentRef{Status: DeploymentStatus_DEPLOYMENT_SUCCESS}},
					{DeploymentRef: &ChainDeploymentRef{Status: DeploymentStatus_DEPLOYMENT_RUNNING}},
				},
			},
			wantDesired: ChainBlockStatus_DEPLOYMENT_BLOCK_RUNNING,
		},
		{
			name: "returns original status if no deployments",
			block: &ChainBlock{
				Status: ChainBlockStatus_DEPLOYMENT_BLOCK_PENDING,
			},
			wantDesired: ChainBlockStatus_DEPLOYMENT_BLOCK_SUCCESS,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotDesired := tt.block.DesiredStatus()
			assert.Equal(t, tt.wantDesired, gotDesired)
		})
	}
}

func TestGetNodeByDeploymentID(t *testing.T) {
	t.Parallel()
	block := &ChainBlock{
		Nodes: []*ChainNode{
			{
				DeploymentRef: &ChainDeploymentRef{DeploymentId: "dep1"},
			},
			{
				DeploymentRef: &ChainDeploymentRef{DeploymentId: "dep2"},
			},
			{
				DeploymentRef: nil,
			},
		},
	}

	tests := []struct {
		name         string
		deploymentID string
		wantNode     *ChainNode
		wantErr      bool
	}{
		{
			name:         "returns node with matching deployment ID",
			deploymentID: "dep1",
			wantNode:     block.Nodes[0],
			wantErr:      false,
		},
		{
			name:         "returns error for non-existent deployment ID",
			deploymentID: "dep3",
			wantNode:     nil,
			wantErr:      true,
		},
		{
			name:         "returns error for nil deployment ref",
			deploymentID: "",
			wantNode:     nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotNode, err := block.GetNodeByDeploymentID(tt.deploymentID)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantNode, gotNode)
		})
	}
}
