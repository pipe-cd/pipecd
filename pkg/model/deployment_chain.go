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

import "fmt"

func GetNotCompletedDeploymentChainStatuses() []ChainStatus {
	return []ChainStatus{
		ChainStatus_DEPLOYMENT_CHAIN_PENDING,
		ChainStatus_DEPLOYMENT_CHAIN_RUNNING,
	}
}

func (dc *DeploymentChain) IsCompleted() bool {
	switch dc.Status {
	case ChainStatus_DEPLOYMENT_CHAIN_SUCCESS,
		ChainStatus_DEPLOYMENT_CHAIN_FAILURE,
		ChainStatus_DEPLOYMENT_CHAIN_CANCELLED:
		return true
	default:
		return false
	}
}

func (dc *DeploymentChain) DesiredStatus() ChainStatus {
	if dc.IsCompleted() {
		return dc.Status
	}

	var (
		successBlockCtn   int
		failedBlockCtn    int
		cancelledBlockCtn int
		runningBlockCtn   int
	)
	for _, block := range dc.Blocks {
		switch block.Status {
		case ChainBlockStatus_DEPLOYMENT_BLOCK_SUCCESS:
			successBlockCtn++
		case ChainBlockStatus_DEPLOYMENT_BLOCK_FAILURE:
			failedBlockCtn++
		case ChainBlockStatus_DEPLOYMENT_BLOCK_CANCELLED:
			cancelledBlockCtn++
		case ChainBlockStatus_DEPLOYMENT_BLOCK_RUNNING:
			runningBlockCtn++
		}
	}

	// Determine desired deployment chain status based on counted blocks' status.
	// If all blocks have finished with SUCCESS status, the chain counted as SUCCESS.
	if successBlockCtn == len(dc.Blocks) {
		return ChainStatus_DEPLOYMENT_CHAIN_SUCCESS
	}
	// If one of the chain's blocks has finished with FAILURE status, the chain counted as FAILURE.
	if failedBlockCtn > 0 {
		return ChainStatus_DEPLOYMENT_CHAIN_FAILURE
	}
	// If one of the chain's blocks has finished with CANCELLED status, the chain counted as CANCELLED.
	if cancelledBlockCtn > 0 {
		return ChainStatus_DEPLOYMENT_CHAIN_CANCELLED
	}
	// If one of the chain's blocks is at RUNNING status, the chain counted as RUNNING.
	if runningBlockCtn > 0 {
		return ChainStatus_DEPLOYMENT_CHAIN_RUNNING
	}
	// Otherwise, the status of deployment chain is remained.
	return dc.Status
}

func (dc *DeploymentChain) ListAllInChainApplicationDeploymentsMap() map[string]*ChainDeploymentRef {
	deployments := make(map[string]*ChainDeploymentRef, 0)
	for _, block := range dc.Blocks {
		for _, node := range block.Nodes {
			if node.DeploymentRef != nil {
				deployments[node.ApplicationRef.ApplicationId] = node.DeploymentRef
			}
		}
	}
	return deployments
}

func (dc *DeploymentChain) ListAllInChainApplications() []*ChainApplicationRef {
	applications := make([]*ChainApplicationRef, 0)
	for _, block := range dc.Blocks {
		for _, node := range block.Nodes {
			applications = append(applications, node.ApplicationRef)
		}
	}
	return applications
}

func (dc *DeploymentChain) SetUpdatedAt(t int64) {
	dc.UpdatedAt = t
}

func (b *ChainBlock) IsCompleted() bool {
	switch b.Status {
	case ChainBlockStatus_DEPLOYMENT_BLOCK_SUCCESS,
		ChainBlockStatus_DEPLOYMENT_BLOCK_FAILURE,
		ChainBlockStatus_DEPLOYMENT_BLOCK_CANCELLED:
		return true
	default:
		return false
	}
}

func (b *ChainBlock) DesiredStatus() ChainBlockStatus {
	if b.IsCompleted() {
		return b.Status
	}

	var (
		successDeploymentCtn   int
		failedDeploymentCtn    int
		cancelledDeploymentCtn int
		runningDeploymentCtn   int
	)
	for _, node := range b.Nodes {
		if node.DeploymentRef == nil {
			continue
		}
		// Count values to determine block status.
		switch node.DeploymentRef.Status {
		case DeploymentStatus_DEPLOYMENT_SUCCESS:
			successDeploymentCtn++
		case DeploymentStatus_DEPLOYMENT_FAILURE:
			failedDeploymentCtn++
		case DeploymentStatus_DEPLOYMENT_CANCELLED:
			cancelledDeploymentCtn++
		case DeploymentStatus_DEPLOYMENT_RUNNING:
			runningDeploymentCtn++
		}
	}

	// Determine block status based on its deployments' state.
	// If all the nodes in block is completed successfully, the block counted as SUCCESS.
	if successDeploymentCtn == len(b.Nodes) {
		return ChainBlockStatus_DEPLOYMENT_BLOCK_SUCCESS
	}
	// If one of the node in the block is completed with FAILURE status, the block counted as FAILURE.
	if failedDeploymentCtn > 0 {
		return ChainBlockStatus_DEPLOYMENT_BLOCK_FAILURE
	}
	// If one of the node in the block is completed with CANCELLED status, the block counted as CANCELLED.
	if cancelledDeploymentCtn > 0 {
		return ChainBlockStatus_DEPLOYMENT_BLOCK_CANCELLED
	}
	// If there is at least a deployment in chain which has RUNNING status,
	// and the block passed all above filters, the block counted as RUNNING.
	if runningDeploymentCtn > 0 {
		return ChainBlockStatus_DEPLOYMENT_BLOCK_RUNNING
	}
	// Otherwise, the block status is remained.
	return b.Status
}

func (b *ChainBlock) GetNodeByDeploymentID(deploymentID string) (*ChainNode, error) {
	for _, node := range b.Nodes {
		if node.DeploymentRef == nil {
			continue
		}
		if node.DeploymentRef.DeploymentId == deploymentID {
			return node, nil
		}
	}
	return nil, fmt.Errorf("unable to find node with the given deployment id (%s)", deploymentID)
}
