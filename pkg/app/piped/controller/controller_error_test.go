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

package controller

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type reportDeploymentCompletedError struct {
	message string
}

func (e *reportDeploymentCompletedError) Error() string {
	return e.message
}

type reportDeploymentCompletedAPIClient struct {
	apiClient

	err    error
	cancel context.CancelFunc
	calls  int
}

func (c *reportDeploymentCompletedAPIClient) ReportDeploymentCompleted(context.Context, *pipedservice.ReportDeploymentCompletedRequest, ...grpc.CallOption) (*pipedservice.ReportDeploymentCompletedResponse, error) {
	c.calls++
	c.cancel()
	return nil, c.err
}

func TestCancelDeploymentWrapsReportDeploymentCompletedError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	underlying := &reportDeploymentCompletedError{
		message: "report deployment completed failed",
	}
	apiClient := &reportDeploymentCompletedAPIClient{
		err:    underlying,
		cancel: cancel,
	}
	c := &controller{
		apiClient: apiClient,
	}

	err := c.cancelDeployment(ctx, &model.Deployment{Id: "deployment-id"}, "cancel reason")

	assert.Equal(t, 1, apiClient.calls)
	assert.ErrorContains(t, err, "failed to report deployment status to control-plane")
	assert.True(t, errors.Is(err, underlying))

	var target *reportDeploymentCompletedError
	assert.True(t, errors.As(err, &target))
	assert.Same(t, underlying, target)
}
