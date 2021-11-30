// Copyright 2021 The PipeCD Authors.
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

func (dc *DeploymentChain) IsCompletedBlock(blockIndex uint32) (bool, error) {
	if blockIndex >= uint32(len(dc.Blocks)) {
		return false, fmt.Errorf("invalid block index %d given", blockIndex)
	}

	block := dc.Blocks[blockIndex]
	switch block.Status {
	case ChainBlockStatus_DEPLOYMENT_BLOCK_SUCCESS,
		ChainBlockStatus_DEPLOYMENT_BLOCK_FAILURE,
		ChainBlockStatus_DEPLOYMENT_BLOCK_CANCELLED:
		return true, nil
	default:
		return false, nil
	}
}
