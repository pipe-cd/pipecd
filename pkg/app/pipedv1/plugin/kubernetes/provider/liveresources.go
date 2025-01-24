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
	"time"

	"github.com/pipe-cd/pipecd/pkg/model"
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

// BuildApplicationLiveState builds the live state of the application from the given manifests.
func BuildApplicationLiveState(deploytarget string, manifests []Manifest, now time.Time) *model.ApplicationLiveState {
	states := make([]*model.ResourceState, 0, len(manifests))
	for _, m := range manifests {
		states = append(states, buildResourceState(m, now))
	}

	return &model.ApplicationLiveState{
		Resources:    states,
		HealthStatus: model.ApplicationLiveState_UNKNOWN, // TODO: Implement health status calculation
	}
}

// buildResourceState builds the resource state from the given manifest.
func buildResourceState(m Manifest, now time.Time) *model.ResourceState {
	parents := make([]string, 0, len(m.body.GetOwnerReferences()))
	for _, o := range m.body.GetOwnerReferences() {
		parents = append(parents, string(o.UID))
	}

	return &model.ResourceState{
		Id:                string(m.body.GetUID()),
		Name:              m.body.GetName(),
		ParentIds:         parents,
		HealthStatus:      model.ResourceState_UNKNOWN, // TODO: Implement health status calculation
		HealthDescription: "",                          // TODO: Implement health status calculation
		ResourceType:      m.body.GetKind(),
		ResourceMetadata: map[string]string{
			"Namespace":   m.body.GetNamespace(),
			"API Version": m.body.GetAPIVersion(),
			"Kind":        m.body.GetKind(),
		},
		CreatedAt: m.body.GetCreationTimestamp().Unix(),
		UpdatedAt: now.Unix(),
	}
}
