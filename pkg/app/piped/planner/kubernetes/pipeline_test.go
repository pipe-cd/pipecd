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

package kubernetes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestBuildQuickSyncPipeline(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		wantAutoRollback bool
	}{
		{
			name:             "want auto rollback stage",
			wantAutoRollback: true,
		},
		{
			name:             "don't want auto rollback stage",
			wantAutoRollback: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotStages := buildQuickSyncPipeline(tc.wantAutoRollback, time.Now())
			var gotAutoRollback bool
			for _, stage := range gotStages {
				if stage.Name == string(model.StageRollback) {
					gotAutoRollback = true
				}
			}
			assert.Equal(t, tc.wantAutoRollback, gotAutoRollback)
		})
	}
}

func TestBuildProgressivePipeline(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		wantAutoRollback bool
	}{
		{
			name:             "want auto rollback stage",
			wantAutoRollback: true,
		},
		{
			name:             "don't want auto rollback stage",
			wantAutoRollback: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotStages := buildProgressivePipeline(&config.DeploymentPipeline{}, tc.wantAutoRollback, time.Now())
			var gotAutoRollback bool
			for _, stage := range gotStages {
				if stage.Name == string(model.StageRollback) {
					gotAutoRollback = true
				}
			}
			assert.Equal(t, tc.wantAutoRollback, gotAutoRollback)
		})
	}
}
