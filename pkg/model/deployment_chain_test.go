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

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeploymentChainDesireStatus(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			desireStatus := tc.deploymentChain.DesiredStatus()
			assert.Equal(t, tc.expectedDesireStatus, desireStatus)
		})
	}
}
