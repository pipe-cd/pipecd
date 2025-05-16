// Copyright 2025 The PipeCD Authors.
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

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/config"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

func executeSyncStage(ctx context.Context, input *sdk.ExecuteStageInput[config.TerraformApplicationSpec], dts []*sdk.DeployTarget[config.TerraformDeployTargetConfig]) sdk.StageStatus {
	panic("not implemented")
}

func executePlanStage(ctx context.Context, input *sdk.ExecuteStageInput[config.TerraformApplicationSpec], dts []*sdk.DeployTarget[config.TerraformDeployTargetConfig]) sdk.StageStatus {
	panic("not implemented")
}

func executeApplyStage(ctx context.Context, input *sdk.ExecuteStageInput[config.TerraformApplicationSpec], dts []*sdk.DeployTarget[config.TerraformDeployTargetConfig]) sdk.StageStatus {
	panic("not implemented")
}

func executeRollbackStage(ctx context.Context, input *sdk.ExecuteStageInput[config.TerraformApplicationSpec], dts []*sdk.DeployTarget[config.TerraformDeployTargetConfig]) sdk.StageStatus {
	panic("not implemented")
}
