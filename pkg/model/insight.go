// Copyright 2022 The PipeCD Authors.
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
	"google.golang.org/protobuf/proto"
)

func (c *InsightDeploymentChunk) Size() (int64, error) {
	raw, err := proto.Marshal(c)
	if err != nil {
		return 0, err
	}

	return int64(len(raw)), nil
}

func (m *InsightDeploymentChunkMetadata) ListChunkNames(from, to int64) []string {
	var paths []string
	for _, m := range m.Chunks {
		if overlap(from, to, m.From, m.To) {
			paths = append(paths, m.Name)
		}
	}
	return paths
}

func overlap(lhsFrom, lhsTo, rhsFrom, rhsTo int64) bool {
	return (rhsFrom < lhsTo && lhsFrom < rhsTo)
}

func (c *InsightDeploymentChunk) ExtractDeployments(from, to int64) []*InsightDeployment {
	var result []*InsightDeployment
	for _, d := range c.Deployments {
		if from <= d.CompletedAt && d.CompletedAt <= to {
			result = append(result, d)
		}
	}
	return result
}
