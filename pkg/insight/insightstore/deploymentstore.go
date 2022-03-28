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
	// List returns slice of Deployment sorted by startedAt ASC.
	List(ctx context.Context, projectID string, from, to int64, minimumVersion *model.InsightDeploymentVersion) ([]*model.InsightDeployment, error)

	Put(ctx context.Context, projectID string, deployments []*model.InsightDeployment, version *model.InsightDeploymentVersion) error
}

// List returns slice of Deployment sorted by startedAt ASC.
func (s *store) List(_ context.Context, _ string, _, _ int64, minimumVersion *model.InsightDeploymentVersion) ([]*model.InsightDeployment, error) {
	return nil, errUnimplemented
}

func (s *store) Put(_ context.Context, _ string, _ []*model.InsightDeployment, version *model.InsightDeploymentVersion) error {
	return errUnimplemented
}
