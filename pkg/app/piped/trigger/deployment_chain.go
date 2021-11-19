// Copyright 2021 The PipeCD Authors.
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

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/config"
)

func (t *Trigger) triggerDeploymentChain(ctx context.Context, dc *config.DeploymentChain) error {
	filters := make([]*pipedservice.CreateDeploymentChainRequest_ApplicationsFilter, 0, len(dc.Nodes))
	for _, node := range dc.Nodes {
		filters = append(filters, &pipedservice.CreateDeploymentChainRequest_ApplicationsFilter{
			AppName:   node.AppName,
			AppKind:   node.AppKind,
			AppLabels: node.AppLabels,
		})
	}

	if _, err := t.apiClient.CreateDeploymentChain(ctx, &pipedservice.CreateDeploymentChainRequest{
		Filters: filters,
	}); err != nil {
		return fmt.Errorf("could not create new deployment chain: %w", err)
	}
	return nil
}
