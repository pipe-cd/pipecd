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

package trigger

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type lastTriggeredCommitStore struct {
	apiClient apiClient
	cache     cache.Cache
}

func (s *lastTriggeredCommitStore) Get(ctx context.Context, applicationID string) (string, error) {
	// Firstly, find from memory cache.
	commit, err := s.cache.Get(applicationID)
	if err == nil {
		return commit.(string), nil
	}

	// No data in memorycache so we have to cost a RPC call to get from control-plane.
	deploy, err := s.getLastTriggeredDeployment(ctx, applicationID)
	switch {
	case err == nil:
		return deploy.Trigger.Commit.Hash, nil

	case status.Code(err) == codes.NotFound:
		// It seems this application has not been deployed anytime.
		return "", nil

	default:
		return "", err
	}
}

func (s *lastTriggeredCommitStore) Put(applicationID, commit string) error {
	return s.cache.Put(applicationID, commit)
}

func (s *lastTriggeredCommitStore) getLastTriggeredDeployment(ctx context.Context, applicationID string) (*model.ApplicationDeploymentReference, error) {
	var (
		err   error
		resp  *pipedservice.GetApplicationMostRecentDeploymentResponse
		retry = pipedservice.NewRetry(3)
		req   = &pipedservice.GetApplicationMostRecentDeploymentRequest{
			ApplicationId: applicationID,
			Status:        model.DeploymentStatus_DEPLOYMENT_PENDING,
		}
	)

	for retry.WaitNext(ctx) {
		if resp, err = s.apiClient.GetApplicationMostRecentDeployment(ctx, req); err == nil {
			return resp.Deployment, nil
		}
		if !pipedservice.Retriable(err) {
			return nil, err
		}
	}
	return nil, err
}
