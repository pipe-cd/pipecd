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

package datastore

import (
	"fmt"
	"testing"

	"github.com/pipe-cd/pipecd/pkg/model"

	"github.com/stretchr/testify/assert"
)

func TestDeploymentChainNodeDeploymentStatusUpdater(t *testing.T) {
	testcases := []struct {
		name             string
		deploymentChain  model.DeploymentChain
		blockIndex       uint32
		deploymentID     string
		deploymentStatus model.DeploymentStatus

		expectedBlockStatus model.ChainBlockStatus
		expectedChainStatus model.ChainStatus
		expectedErr         error
	}{
		{
			name: "invalid blockIndex given",
			deploymentChain: model.DeploymentChain{
				Status: model.ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
				Blocks: []*model.ChainBlock{
					{
						Nodes: []*model.ChainNode{
							{
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-1",
									Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
								},
							},
						},
					},
				},
			},
			blockIndex:       1,
			deploymentID:     "deploy-1",
			deploymentStatus: model.DeploymentStatus_DEPLOYMENT_SUCCESS,
			expectedErr:      fmt.Errorf("invalid block index 1 provided"),
		},
		{
			name: "deployment id not found in the block",
			deploymentChain: model.DeploymentChain{
				Status: model.ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
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
			},
			blockIndex:       0,
			deploymentID:     "deploy-1",
			deploymentStatus: model.DeploymentStatus_DEPLOYMENT_SUCCESS,
			expectedErr:      fmt.Errorf("unable to find the right node in chain to assign deployment to"),
		},
		{
			name: "try to update a finished block",
			deploymentChain: model.DeploymentChain{
				Status: model.ChainStatus_DEPLOYMENT_CHAIN_FAILURE,
				Blocks: []*model.ChainBlock{
					{
						Status: model.ChainBlockStatus_DEPLOYMENT_BLOCK_FAILURE,
						Nodes: []*model.ChainNode{
							{
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-0",
									Status:       model.DeploymentStatus_DEPLOYMENT_FAILURE,
								},
							},
							{
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-1",
									Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
								},
							},
						},
					},
				},
			},
			blockIndex:          0,
			deploymentID:        "deploy-1",
			deploymentStatus:    model.DeploymentStatus_DEPLOYMENT_CANCELLED,
			expectedBlockStatus: model.ChainBlockStatus_DEPLOYMENT_BLOCK_FAILURE,
			expectedChainStatus: model.ChainStatus_DEPLOYMENT_CHAIN_FAILURE,
		},
		{
			name: "block reaches SUCCESS status after update its last deployment with SUCCESS status",
			deploymentChain: model.DeploymentChain{
				Status: model.ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
				Blocks: []*model.ChainBlock{
					{
						Nodes: []*model.ChainNode{
							{
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-0",
									Status:       model.DeploymentStatus_DEPLOYMENT_SUCCESS,
								},
							},
							{
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-1",
									Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
								},
							},
						},
						Status: model.ChainBlockStatus_DEPLOYMENT_BLOCK_RUNNING,
					},
				},
			},
			blockIndex:          0,
			deploymentID:        "deploy-1",
			deploymentStatus:    model.DeploymentStatus_DEPLOYMENT_SUCCESS,
			expectedBlockStatus: model.ChainBlockStatus_DEPLOYMENT_BLOCK_SUCCESS,
			expectedChainStatus: model.ChainStatus_DEPLOYMENT_CHAIN_SUCCESS,
		},
		{
			name: "block is marked as FAILURE after update its deployment with FAILURE status",
			deploymentChain: model.DeploymentChain{
				Status: model.ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
				Blocks: []*model.ChainBlock{
					{
						Nodes: []*model.ChainNode{
							{
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-0",
									Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
								},
							},
							{
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-1",
									Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
								},
							},
						},
						Status: model.ChainBlockStatus_DEPLOYMENT_BLOCK_RUNNING,
					},
				},
			},
			blockIndex:          0,
			deploymentID:        "deploy-1",
			deploymentStatus:    model.DeploymentStatus_DEPLOYMENT_FAILURE,
			expectedBlockStatus: model.ChainBlockStatus_DEPLOYMENT_BLOCK_FAILURE,
			expectedChainStatus: model.ChainStatus_DEPLOYMENT_CHAIN_FAILURE,
		},
		{
			name: "block is marked CANCELLED status after update its deployment with CANCELLED status",
			deploymentChain: model.DeploymentChain{
				Status: model.ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
				Blocks: []*model.ChainBlock{
					{
						Nodes: []*model.ChainNode{
							{
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-0",
									Status:       model.DeploymentStatus_DEPLOYMENT_SUCCESS,
								},
							},
							{
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-1",
									Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
								},
							},
						},
						Status: model.ChainBlockStatus_DEPLOYMENT_BLOCK_RUNNING,
					},
				},
			},
			blockIndex:          0,
			deploymentID:        "deploy-1",
			deploymentStatus:    model.DeploymentStatus_DEPLOYMENT_CANCELLED,
			expectedBlockStatus: model.ChainBlockStatus_DEPLOYMENT_BLOCK_CANCELLED,
			expectedChainStatus: model.ChainStatus_DEPLOYMENT_CHAIN_CANCELLED,
		},
		{
			name: "block keep it status on update not the first started deployment status to RUNNING",
			deploymentChain: model.DeploymentChain{
				Status: model.ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
				Blocks: []*model.ChainBlock{
					{
						Nodes: []*model.ChainNode{
							{
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-0",
									Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
								},
							},
							{
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-1",
									Status:       model.DeploymentStatus_DEPLOYMENT_PENDING,
								},
							},
							{
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-2",
									Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
								},
							},
						},
						Status: model.ChainBlockStatus_DEPLOYMENT_BLOCK_RUNNING,
					},
				},
			},
			blockIndex:          0,
			deploymentID:        "deploy-1",
			deploymentStatus:    model.DeploymentStatus_DEPLOYMENT_RUNNING,
			expectedBlockStatus: model.ChainBlockStatus_DEPLOYMENT_BLOCK_RUNNING,
			expectedChainStatus: model.ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
		},
		{
			name: "block keep it status on update not the first started deployment status to SUCCESS",
			deploymentChain: model.DeploymentChain{
				Status: model.ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
				Blocks: []*model.ChainBlock{
					{
						Nodes: []*model.ChainNode{
							{
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-0",
									Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
								},
							},
							{
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-1",
									Status:       model.DeploymentStatus_DEPLOYMENT_PENDING,
								},
							},
							{
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-2",
									Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
								},
							},
						},
						Status: model.ChainBlockStatus_DEPLOYMENT_BLOCK_RUNNING,
					},
				},
			},
			blockIndex:          0,
			deploymentID:        "deploy-2",
			deploymentStatus:    model.DeploymentStatus_DEPLOYMENT_SUCCESS,
			expectedBlockStatus: model.ChainBlockStatus_DEPLOYMENT_BLOCK_RUNNING,
			expectedChainStatus: model.ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
		},
		{
			name: "block changes its status as RUNNING once the it has one RUNNING deployment",
			deploymentChain: model.DeploymentChain{
				Status: model.ChainStatus_DEPLOYMENT_CHAIN_PENDING,
				Blocks: []*model.ChainBlock{
					{
						Nodes: []*model.ChainNode{
							{
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-0",
									Status:       model.DeploymentStatus_DEPLOYMENT_PENDING,
								},
							},
							{
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-1",
									Status:       model.DeploymentStatus_DEPLOYMENT_PENDING,
								},
							},
							{
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-2",
									Status:       model.DeploymentStatus_DEPLOYMENT_PENDING,
								},
							},
						},
						Status: model.ChainBlockStatus_DEPLOYMENT_BLOCK_PENDING,
					},
				},
			},
			blockIndex:          0,
			deploymentID:        "deploy-1",
			deploymentStatus:    model.DeploymentStatus_DEPLOYMENT_RUNNING,
			expectedBlockStatus: model.ChainBlockStatus_DEPLOYMENT_BLOCK_RUNNING,
			expectedChainStatus: model.ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			updater := nodeDeploymentStatusUpdateFunc(tc.blockIndex, tc.deploymentID, tc.deploymentStatus, "")
			err := updater(&tc.deploymentChain)
			if err != nil {
				if tc.expectedErr == nil {
					assert.NoError(t, err)
					return
				}
				assert.Error(t, err, tc.expectedErr)
				return
			}
			assert.Equal(t, tc.expectedBlockStatus, tc.deploymentChain.Blocks[tc.blockIndex].Status)
			assert.Equal(t, tc.expectedChainStatus, tc.deploymentChain.Status)
		})
	}
}
