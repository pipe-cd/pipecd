// Copyright 2026 The PipeCD Authors.
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

package deployment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/config"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

func TestPlugin_executePlanStage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   *sdk.ExecuteStageInput[config.ApplicationConfigSpec]
		dts     []*sdk.DeployTarget[config.DeployTargetConfig]
		want    sdk.StageStatus
	}{
		{
			name: "failure when StageLogPersister is unavailable",
			input: &sdk.ExecuteStageInput[config.ApplicationConfigSpec]{
				Client: sdk.NewClient(nil, "terraform", "app1", "stage1", nil, nil),
				Logger: zap.NewNop(),
			},
			dts:  []*sdk.DeployTarget[config.DeployTargetConfig]{{}},
			want: sdk.StageStatusFailure,
		},
	}

	p := &Plugin{}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := p.executePlanStage(context.Background(), tc.input, tc.dts)
			assert.Equal(t, tc.want, got)
		})
	}
}
