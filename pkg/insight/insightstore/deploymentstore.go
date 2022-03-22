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

package insightstore

import (
	"context"
	"errors"

	"github.com/pipe-cd/pipecd/pkg/model"
)

var (
	errUnimplemented = errors.New("unimplemented")
)

type DeploymentStore interface {
	// GetDailyDeployments returns slice of DailyDeployment sorted by DailyDeployment.Date ASC.
	GetDailyDeployments(ctx context.Context, projectID string, dateRage *model.ChunkDateRange) ([]*model.DailyDeployment, error)

	PutDailyDeployment(ctx context.Context, projectID string, deployments *model.DailyDeployment) error
}

// GetDailyDeployments returns slice of DailyDeployment sorted by DailyDeployment.Date ASC.
func (s *store) GetDailyDeployments(_ context.Context, _ string, _ *model.ChunkDateRange) ([]*model.DailyDeployment, error) {
	return nil, errUnimplemented
}

func (s *store) PutDailyDeployment(_ context.Context, _ string, _ *model.DailyDeployment) error {
	return errUnimplemented
}
