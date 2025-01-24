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

package provider

import (
	"context"
	"fmt"
)

// GetLiveResources returns all live resources that belong to the given application.
func GetLiveResources(ctx context.Context, kubectl *Kubectl, kubeconfig string, appID string, selector ...string) (namespaceScoped []Manifest, clusterScoped []Manifest, _ error) {
	namespacedLiveResources, err := kubectl.GetAll(ctx, kubeconfig,
		"",
		fmt.Sprintf("%s=%s", LabelManagedBy, ManagedByPiped),
		fmt.Sprintf("%s=%s", LabelApplication, appID),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed while listing all namespace-scoped resources (%v)", err)
	}

	clusterScopedLiveResources, err := kubectl.GetAllClusterScoped(ctx, kubeconfig,
		fmt.Sprintf("%s=%s", LabelManagedBy, ManagedByPiped),
		fmt.Sprintf("%s=%s", LabelApplication, appID),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed while listing all cluster-scoped resources (%v)", err)
	}

	return namespacedLiveResources, clusterScopedLiveResources, nil
}
