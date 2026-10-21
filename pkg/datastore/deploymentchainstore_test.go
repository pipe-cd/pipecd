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

func TestDeploymentChainAddDeploymentToBlock(t *testing.T) {
	testcases := []struct {
		name            string
		deploymentChain model.DeploymentChain
		deployment      *model.Deployment

		expectedErr    bool
		expectedRefID  string
		expectedStatus model.DeploymentStatus
	}{
		{
			name: "invalid block index given",
			deploymentChain: model.DeploymentChain{
				Blocks: []*model.ChainBlock{{}},
			},
			deployment: &model.Deployment{
				Id:                        "deploy-1",
				ApplicationId:             "app-1",
				DeploymentChainBlockIndex: 5,
			},
			expectedErr: true,
		},
		{
			name: "no node of the given application in block",
			deploymentChain: model.DeploymentChain{
				Blocks: []*model.ChainBlock{
					{
						Nodes: []*model.ChainNode{
							{ApplicationRef: &model.ChainApplicationRef{ApplicationId: "app-other"}},
						},
					},
				},
			},
			deployment: &model.Deployment{
				Id:                        "deploy-1",
				ApplicationId:             "app-1",
				DeploymentChainBlockIndex: 0,
			},
			expectedErr: true,
		},
		{
			name: "link deployment to its node",
			deploymentChain: model.DeploymentChain{
				Blocks: []*model.ChainBlock{
					{
						Nodes: []*model.ChainNode{
							{ApplicationRef: &model.ChainApplicationRef{ApplicationId: "app-1"}},
						},
					},
				},
			},
			deployment: &model.Deployment{
				Id:                        "deploy-1",
				ApplicationId:             "app-1",
				DeploymentChainBlockIndex: 0,
				Status:                    model.DeploymentStatus_DEPLOYMENT_PENDING,
			},
			expectedRefID:  "deploy-1",
			expectedStatus: model.DeploymentStatus_DEPLOYMENT_PENDING,
		},
		{
			// Retrying CreateDeployment must not reset a status which has
			// already been advanced by a status report.
			name: "relinking the same deployment keeps the reported status",
			deploymentChain: model.DeploymentChain{
				Blocks: []*model.ChainBlock{
					{
						Nodes: []*model.ChainNode{
							{
								ApplicationRef: &model.ChainApplicationRef{ApplicationId: "app-1"},
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-1",
									Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
								},
							},
						},
					},
				},
			},
			deployment: &model.Deployment{
				Id:                        "deploy-1",
				ApplicationId:             "app-1",
				DeploymentChainBlockIndex: 0,
				Status:                    model.DeploymentStatus_DEPLOYMENT_PENDING,
			},
			expectedRefID:  "deploy-1",
			expectedStatus: model.DeploymentStatus_DEPLOYMENT_RUNNING,
		},
		{
			name: "a newly triggered deployment of the same application overwrites the ref",
			deploymentChain: model.DeploymentChain{
				Blocks: []*model.ChainBlock{
					{
						Nodes: []*model.ChainNode{
							{
								ApplicationRef: &model.ChainApplicationRef{ApplicationId: "app-1"},
								DeploymentRef: &model.ChainDeploymentRef{
									DeploymentId: "deploy-old",
									Status:       model.DeploymentStatus_DEPLOYMENT_SUCCESS,
								},
							},
						},
					},
				},
			},
			deployment: &model.Deployment{
				Id:                        "deploy-new",
				ApplicationId:             "app-1",
				DeploymentChainBlockIndex: 0,
				Status:                    model.DeploymentStatus_DEPLOYMENT_PENDING,
			},
			expectedRefID:  "deploy-new",
			expectedStatus: model.DeploymentStatus_DEPLOYMENT_PENDING,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			updater := addDeploymentToBlockUpdateFunc(tc.deployment)
			err := updater(&tc.deploymentChain)
			if tc.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			ref := tc.deploymentChain.Blocks[tc.deployment.DeploymentChainBlockIndex].Nodes[0].DeploymentRef
			assert.Equal(t, tc.expectedRefID, ref.DeploymentId)
			assert.Equal(t, tc.expectedStatus, ref.Status)
		})
	}
}

// TestDeploymentChainSettlesOnLastDeploymentCompletion is the regression test for the
// deployment chain which never reached a terminal status. Previously nothing wrote the
// chain status once every deployment had settled; now the report which completes the
// last deployment recalculates and settles both its block and the chain.
func TestDeploymentChainSettlesOnLastDeploymentCompletion(t *testing.T) {
	// A two block chain, the first block already succeeded and the second one has a
	// single remaining running deployment.
	dc := &model.DeploymentChain{
		Status: model.ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
		Blocks: []*model.ChainBlock{
			{
				Status: model.ChainBlockStatus_DEPLOYMENT_BLOCK_SUCCESS,
				Nodes: []*model.ChainNode{
					{
						ApplicationRef: &model.ChainApplicationRef{ApplicationId: "app-1"},
						DeploymentRef: &model.ChainDeploymentRef{
							DeploymentId: "deploy-1",
							Status:       model.DeploymentStatus_DEPLOYMENT_SUCCESS,
						},
					},
				},
			},
			{
				Status: model.ChainBlockStatus_DEPLOYMENT_BLOCK_RUNNING,
				Nodes: []*model.ChainNode{
					{
						ApplicationRef: &model.ChainApplicationRef{ApplicationId: "app-2"},
						DeploymentRef: &model.ChainDeploymentRef{
							DeploymentId: "deploy-2",
							Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
						},
					},
				},
			},
		},
	}

	err := nodeDeploymentStatusUpdateFunc(1, "deploy-2", model.DeploymentStatus_DEPLOYMENT_SUCCESS, "")(dc)
	assert.NoError(t, err)

	assert.Equal(t, model.ChainBlockStatus_DEPLOYMENT_BLOCK_SUCCESS, dc.Blocks[1].Status)
	assert.Equal(t, model.ChainStatus_DEPLOYMENT_CHAIN_SUCCESS, dc.Status)
	assert.NotZero(t, dc.CompletedAt, "a settled chain must have CompletedAt set")

	// Applying the very same report again must not change anything.
	completedAt := dc.CompletedAt
	err = nodeDeploymentStatusUpdateFunc(1, "deploy-2", model.DeploymentStatus_DEPLOYMENT_SUCCESS, "")(dc)
	assert.NoError(t, err)
	assert.Equal(t, model.ChainStatus_DEPLOYMENT_CHAIN_SUCCESS, dc.Status)
	assert.Equal(t, completedAt, dc.CompletedAt, "reapplying a terminal report must not move CompletedAt")
}

// TestDeploymentChainBlockNeedsEveryNodeLinked locks in the reason restoring the linking
// in CreateDeployment is load-bearing: a block holding a node without a deployment ref
// can never be counted as succeeded.
func TestDeploymentChainBlockNeedsEveryNodeLinked(t *testing.T) {
	dc := &model.DeploymentChain{
		Status: model.ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
		Blocks: []*model.ChainBlock{
			{
				Status: model.ChainBlockStatus_DEPLOYMENT_BLOCK_RUNNING,
				Nodes: []*model.ChainNode{
					{
						ApplicationRef: &model.ChainApplicationRef{ApplicationId: "app-1"},
						DeploymentRef: &model.ChainDeploymentRef{
							DeploymentId: "deploy-1",
							Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
						},
					},
					// Never linked, e.g. its piped has not triggered the deployment yet.
					{ApplicationRef: &model.ChainApplicationRef{ApplicationId: "app-2"}},
				},
			},
		},
	}

	err := nodeDeploymentStatusUpdateFunc(0, "deploy-1", model.DeploymentStatus_DEPLOYMENT_SUCCESS, "")(dc)
	assert.NoError(t, err)

	assert.NotEqual(t, model.ChainBlockStatus_DEPLOYMENT_BLOCK_SUCCESS, dc.Blocks[0].Status)
	assert.NotEqual(t, model.ChainStatus_DEPLOYMENT_CHAIN_SUCCESS, dc.Status)
}

// TestDeploymentChainUpdateUnknownDeployment locks in that reporting a status for a
// deployment which is not linked to any node fails instead of silently doing nothing.
func TestDeploymentChainUpdateUnknownDeployment(t *testing.T) {
	dc := &model.DeploymentChain{
		Status: model.ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
		Blocks: []*model.ChainBlock{
			{
				Status: model.ChainBlockStatus_DEPLOYMENT_BLOCK_RUNNING,
				Nodes: []*model.ChainNode{
					{ApplicationRef: &model.ChainApplicationRef{ApplicationId: "app-1"}},
				},
			},
		},
	}

	err := nodeDeploymentStatusUpdateFunc(0, "deploy-unknown", model.DeploymentStatus_DEPLOYMENT_SUCCESS, "")(dc)
	assert.Error(t, err)
	assert.Equal(t, model.ChainStatus_DEPLOYMENT_CHAIN_RUNNING, dc.Status)
}
