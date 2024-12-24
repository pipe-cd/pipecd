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

package execute

import (
	"context"

	"github.com/pipe-cd/pipecd/pkg/app/piped/logpersister"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
)

type Stage string

const (
	stageWait Stage = "WAIT"
)

// Execute starts waiting for the specified duration.
func (s *deploymentServiceServer) execute(ctx context.Context, in *deployment.ExecutePluginInput, slp logpersister.StageLogPersister) model.StageStatus {
	// TOD: implement the logic of waiting
	return model.StageStatus_STAGE_FAILURE
}
