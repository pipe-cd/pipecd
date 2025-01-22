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

package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestDeploymentTraceID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		d1     *model.Deployment
		d2     *model.Deployment
		assert assert.ComparisonAssertionFunc
	}{
		{
			name: "same deployment id",
			d1: &model.Deployment{
				Id: "example-deployment-id",
			},
			d2: &model.Deployment{
				Id: "example-deployment-id",
			},
			assert: assert.Equal,
		},
		{
			name: "different deployment id",
			d1: &model.Deployment{
				Id: "example-deployment-id",
			},
			d2: &model.Deployment{
				Id: "example-deployment-id-other",
			},
			assert: assert.NotEqual,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id1 := deploymentTraceID(tt.d1)
			id2 := deploymentTraceID(tt.d2)

			tt.assert(t, id1, id2)
		})
	}
}
