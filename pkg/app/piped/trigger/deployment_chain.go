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

package trigger

import (
	"context"
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func (t *Trigger) triggerDeploymentChain(
	ctx context.Context,
	dc *config.DeploymentChain,
	firstDeployment *model.Deployment,
) error {
	matchers := make([]*pipedservice.CreateDeploymentChainRequest_ApplicationMatcher, 0, len(dc.ApplicationMatchers))
	for _, m := range dc.ApplicationMatchers {
		matchers = append(matchers, &pipedservice.CreateDeploymentChainRequest_ApplicationMatcher{
			Name:   m.Name,
			Kind:   m.Kind,
			Labels: m.Labels,
		})
	}

	if _, err := t.apiClient.CreateDeploymentChain(ctx, &pipedservice.CreateDeploymentChainRequest{
		Matchers:        matchers,
		FirstDeployment: firstDeployment,
	}); err != nil {
		return fmt.Errorf("could not create new deployment chain: %w", err)
	}
	return nil
}
