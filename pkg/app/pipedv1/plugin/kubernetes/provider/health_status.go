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
	corev1 "k8s.io/api/core/v1"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

func (m Manifest) calculateHealthStatus() (sdk.ResourceHealthStatus, string) {
	switch {
	case m.IsDeployment():
		obj := &appsv1.Deployment{}
		if err := m.ConvertToStructuredObject(obj); err != nil {
			return sdk.ResourceHealthStateUnknown, ""
		}
		return deploymentHealthStatus(obj)
	case m.IsStatefulSet():
		obj := &appsv1.StatefulSet{}
		if err := m.ConvertToStructuredObject(obj); err != nil {
			return sdk.ResourceHealthStateUnknown, ""
		}
		return statefulSetHealthStatus(obj)
	case m.IsReplicaSet():
		obj := &appsv1.ReplicaSet{}
		if err := m.ConvertToStructuredObject(obj); err != nil {
			return sdk.ResourceHealthStateUnknown, ""
		}
		return replicaSetHealthStatus(obj)
	case m.IsDaemonSet():
		obj := &appsv1.DaemonSet{}
		if err := m.ConvertToStructuredObject(obj); err != nil {
			return sdk.ResourceHealthStateUnknown, ""
		}
		return daemonSetHealthStatus(obj)
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
	if obj.Status.UpdatedReplicas < *obj.Spec.Replicas {
		return sdk.ResourceHealthStateUnhealthy, fmt.Sprintf("Waiting for remaining %d/%d replicas to be updated", obj.Status.UpdatedReplicas, *obj.Spec.Replicas)
	}
	if obj.Status.UpdatedReplicas < obj.Status.Replicas {
		return sdk.ResourceHealthStateUnhealthy, fmt.Sprintf("%d old replicas are pending termination", obj.Status.Replicas-obj.Status.UpdatedReplicas)
	}
	if obj.Status.AvailableReplicas < obj.Status.Replicas {
		return sdk.ResourceHealthStateUnhealthy, fmt.Sprintf("Waiting for remaining %d/%d replicas to be available", obj.Status.Replicas-obj.Status.AvailableReplicas, obj.Status.Replicas)
	}
	return sdk.ResourceHealthStateHealthy, ""
}

func statefulSetHealthStatus(obj *appsv1.StatefulSet) (sdk.ResourceHealthStatus, string) {
	// Referred to:
	//   https://github.com/kubernetes/kubernetes/blob/7942dca975b7be9386540df3c17e309c3cb2de60/staging/src/k8s.io/kubectl/pkg/polymorphichelpers/rollout_status.go#L130-L149
	if obj.Status.ObservedGeneration == 0 || obj.Generation > obj.Status.ObservedGeneration {
		return sdk.ResourceHealthStateUnhealthy, "Waiting for statefulset spec update to be observed"
	}

	if obj.Spec.Replicas == nil {
		return sdk.ResourceHealthStateUnhealthy, "The number of desired replicas is unspecified"
	}
	if *obj.Spec.Replicas != obj.Status.ReadyReplicas {
		return sdk.ResourceHealthStateUnhealthy, fmt.Sprintf("The number of ready replicas (%d) is different from the desired number (%d)", obj.Status.ReadyReplicas, *obj.Spec.Replicas)
	}

	// Check if the partitioned roll out is in progress.
	if obj.Spec.UpdateStrategy.Type == appsv1.RollingUpdateStatefulSetStrategyType && obj.Spec.UpdateStrategy.RollingUpdate != nil {
		if obj.Spec.Replicas != nil && obj.Spec.UpdateStrategy.RollingUpdate.Partition != nil {
			if obj.Status.UpdatedReplicas < (*obj.Spec.Replicas - *obj.Spec.UpdateStrategy.RollingUpdate.Partition) {
				return sdk.ResourceHealthStateUnhealthy, fmt.Sprintf("Waiting for partitioned roll out to finish because %d out of %d new pods have been updated",
					obj.Status.UpdatedReplicas, (*obj.Spec.Replicas - *obj.Spec.UpdateStrategy.RollingUpdate.Partition))
			}
		}
		return sdk.ResourceHealthStateHealthy, ""
	}

	if obj.Status.UpdateRevision != obj.Status.CurrentRevision {
		return sdk.ResourceHealthStateUnhealthy, fmt.Sprintf("Waiting for statefulset rolling update to complete %d pods at revision %s", obj.Status.UpdatedReplicas, obj.Status.UpdateRevision)
	}

	return sdk.ResourceHealthStateHealthy, ""
}

func replicaSetHealthStatus(obj *appsv1.ReplicaSet) (sdk.ResourceHealthStatus, string) {
	if obj.Status.ObservedGeneration == 0 || obj.Generation > obj.Status.ObservedGeneration {
		return sdk.ResourceHealthStateUnhealthy, "Waiting for rollout to finish because observed replica set generation less than desired generation"
	}

	var cond *appsv1.ReplicaSetCondition
	for i := range obj.Status.Conditions {
		c := obj.Status.Conditions[i]
		if c.Type == appsv1.ReplicaSetReplicaFailure {
			cond = &c
			break
		}
	}
	if cond != nil && cond.Status == corev1.ConditionTrue {
		return sdk.ResourceHealthStateUnhealthy, cond.Message
	}

	if obj.Spec.Replicas == nil {
		return sdk.ResourceHealthStateUnhealthy, "The number of desired replicas is unspecified"
	}

	if obj.Status.AvailableReplicas < *obj.Spec.Replicas {
		return sdk.ResourceHealthStateUnhealthy, fmt.Sprintf("Waiting for remaining %d/%d replicas to be available", obj.Status.Replicas-obj.Status.AvailableReplicas, obj.Status.Replicas)
	}

	if *obj.Spec.Replicas != obj.Status.ReadyReplicas {
		return sdk.ResourceHealthStateUnhealthy, fmt.Sprintf("The number of ready replicas (%d) is different from the desired number (%d)", obj.Status.ReadyReplicas, *obj.Spec.Replicas)
	}

	return sdk.ResourceHealthStateHealthy, ""
}

func daemonSetHealthStatus(obj *appsv1.DaemonSet) (sdk.ResourceHealthStatus, string) {
	if obj.Status.ObservedGeneration == 0 || obj.Generation > obj.Status.ObservedGeneration {
		return sdk.ResourceHealthStateUnhealthy, "Waiting for rollout to finish because observed daemon set generation less than desired generation"
	}

	if obj.Status.UpdatedNumberScheduled < obj.Status.DesiredNumberScheduled {
		return sdk.ResourceHealthStateUnhealthy, fmt.Sprintf("Waiting for daemon set %q rollout to finish because %d out of %d new pods have been updated", obj.GetName(), obj.Status.UpdatedNumberScheduled, obj.Status.DesiredNumberScheduled)
	}
	if obj.Status.NumberAvailable < obj.Status.DesiredNumberScheduled {
		return sdk.ResourceHealthStateUnhealthy, fmt.Sprintf("Waiting for daemon set %q rollout to finish because %d of %d updated pods are available", obj.GetName(), obj.Status.NumberAvailable, obj.Status.DesiredNumberScheduled)
	}

	if obj.Status.NumberMisscheduled > 0 {
		return sdk.ResourceHealthStateUnhealthy, fmt.Sprintf("%d nodes that are running the daemon pod, but are not supposed to run the daemon pod", obj.Status.NumberMisscheduled)
	}
	if obj.Status.NumberUnavailable > 0 {
		return sdk.ResourceHealthStateUnhealthy, fmt.Sprintf("%d nodes that should be running the daemon pod and have none of the daemon pod running and available", obj.Status.NumberUnavailable)
	}

	return sdk.ResourceHealthStateHealthy, ""
}
