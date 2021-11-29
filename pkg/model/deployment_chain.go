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

import (
	"errors"
)

func (dc *DeploymentChain) IsSuccessfullyCompletedBlock(blockIndex uint32) (bool, error) {
	if blockIndex >= uint32(len(dc.Blocks)) {
		return false, errors.New("invalid block index given")
	}

	block := dc.Blocks[blockIndex]
	for _, node := range block.Nodes {
		if !IsSuccessfullyCompletedDeployment(node.DeploymentRef.Status) {
			return false, nil
		}
	}
	return true, nil
}
