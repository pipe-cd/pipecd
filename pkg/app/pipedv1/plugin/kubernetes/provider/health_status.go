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

package provider

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"

	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

func (m Manifest) calculateHealthStatus() (sdk.ResourceHealthStatus, string) {
	switch {
	case m.IsDeployment():
		obj := &appsv1.Deployment{}
		if err := m.ConvertToStructuredObject(obj); err != nil {
			return sdk.ResourceHealthStateUnknown, ""
		}
		return deploymentHealthStatus(obj)
	default:
		// TODO: Implement health status calculation for other resource types.
		return sdk.ResourceHealthStateUnknown, fmt.Sprintf("Unimplemented or unknown resource: %s", m.body.GroupVersionKind())
	}
}

func deploymentHealthStatus(obj *appsv1.Deployment) (sdk.ResourceHealthStatus, string) {
	if obj.Spec.Paused {
		return sdk.ResourceHealthStateUnknown, "Deployment is paused"
	}
	if obj.Generation > obj.Status.ObservedGeneration {
		return sdk.ResourceHealthStateUnknown, "Waiting for rollout to finish because observed deployment generation less than desired generation"
	}
	const (
		reasonTimeout = "ProgressDeadlineExceeded"
	)
	for _, cond := range obj.Status.Conditions {
		if cond.Type == appsv1.DeploymentProgressing && cond.Reason == reasonTimeout {
			return sdk.ResourceHealthStateUnhealthy, fmt.Sprintf("Deployment %q exceeded its progress deadline", obj.GetName())
		}
	}

	if obj.Spec.Replicas == nil {
		return sdk.ResourceHealthStateUnknown, "The number of desired replicas is unspecified"
	}
	if obj.Status.Replicas < *obj.Spec.Replicas {
		return sdk.ResourceHealthStateUnhealthy, fmt.Sprintf("Waiting for remaining %d/%d replicas to be updated", obj.Status.Replicas-obj.Status.AvailableReplicas, obj.Status.Replicas)
	}
	if obj.Status.UpdatedReplicas < obj.Status.Replicas {
		return sdk.ResourceHealthStateUnhealthy, fmt.Sprintf("%d old replicas are pending termination", obj.Status.Replicas-obj.Status.UpdatedReplicas)
	}
	if obj.Status.AvailableReplicas < obj.Status.Replicas {
		return sdk.ResourceHealthStateUnhealthy, fmt.Sprintf("Waiting for remaining %d/%d replicas to be available", obj.Status.Replicas-obj.Status.AvailableReplicas, obj.Status.Replicas)
	}
	return sdk.ResourceHealthStateHealthy, ""
}
