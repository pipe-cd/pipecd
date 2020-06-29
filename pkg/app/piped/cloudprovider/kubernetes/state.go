// Copyright 2020 The PipeCD Authors.
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

package kubernetes

import (
	"fmt"
	"sort"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/pipe-cd/pipe/pkg/model"
)

func MakeKubernetesResourceState(uid string, key ResourceKey, obj *unstructured.Unstructured, now time.Time) model.KubernetesResourceState {
	var (
		owners       = obj.GetOwnerReferences()
		ownerIDs     = make([]string, 0, len(owners))
		creationTime = obj.GetCreationTimestamp()
		status, desc = determineResourceHealth(key, obj)
	)

	for _, owner := range owners {
		ownerIDs = append(ownerIDs, string(owner.UID))
	}
	sort.Strings(ownerIDs)

	state := model.KubernetesResourceState{
		Id:         uid,
		OwnerIds:   ownerIDs,
		Name:       key.Name,
		ApiVersion: key.APIVersion,
		Kind:       key.Kind,
		Namespace:  obj.GetNamespace(),

		HealthStatus:      status,
		HealthDescription: desc,

		CreatedAt: creationTime.Unix(),
		UpdatedAt: now.Unix(),
	}

	return state
}

func determineResourceHealth(key ResourceKey, obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	if !IsKubernetesBuiltInResource(key.APIVersion) {
		desc = fmt.Sprintf("Unreadable resource kind %s/%s", key.APIVersion, key.Kind)
		return
	}

	switch key.Kind {
	case KindDeployment:
		return determineDeploymentHealth(obj)
	case KindStatefulSet:
		return determineStatefulSetHealth(obj)
	case KindDaemonSet:
		return determineDaemonSetHealth(obj)
	case KindReplicaSet:
		return determineReplicaSetHealth(obj)
	case KindPod:
		return determinePodHealth(obj)
	case KindService:
		return determineServiceHealth(obj)
	case KindIngress:
		return determineIngressHealth(obj)
	}

	return
}

func determineDeploymentHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	d := &appsv1.Deployment{}
	err := scheme.Scheme.Convert(obj, d, nil)
	if err != nil {
		desc = fmt.Sprintf("Failed while convert %T to %T: %v", obj, d, err)
		return
	}

	status = model.KubernetesResourceState_OTHER
	if d.Spec.Paused {
		desc = "Deployment is paused"
		return
	}
	if d.Status.UpdatedReplicas < d.Status.Replicas {
		desc = fmt.Sprintf("Waiting for remaining %d/%d replicas to be updated", d.Status.Replicas-d.Status.UpdatedReplicas, d.Status.Replicas)
		return
	}
	if d.Status.AvailableReplicas < d.Status.Replicas {
		desc = fmt.Sprintf("Waiting for remaining %d/%d replicas to be available", d.Status.Replicas-d.Status.AvailableReplicas, d.Status.Replicas)
		return
	}

	status = model.KubernetesResourceState_HEALTHY
	return
}

func determineStatefulSetHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	s := &appsv1.StatefulSet{}
	err := scheme.Scheme.Convert(obj, s, nil)
	if err != nil {
		desc = fmt.Sprintf("Failed while convert %T to %T: %v", obj, s, err)
		return
	}
	status = model.KubernetesResourceState_OTHER
	if s.Status.ObservedGeneration == 0 || s.Generation > s.Status.ObservedGeneration {
		desc = "Waiting for statefulset spec update to be observed"
		return
	}

	if s.Spec.Replicas == nil {
		desc = "The number of desired replicas is unspecified"
		return
	}
	if *s.Spec.Replicas != s.Status.ReadyReplicas {
		desc = fmt.Sprintf("The number of ready replicas (%d) is different from the desired number (%d)", s.Status.ReadyReplicas, *s.Spec.Replicas)
		return
	}

	// Check if the partitioned roll out is in progress.
	if s.Spec.UpdateStrategy.Type == appsv1.RollingUpdateStatefulSetStrategyType && s.Spec.UpdateStrategy.RollingUpdate != nil {
		if s.Spec.Replicas != nil && s.Spec.UpdateStrategy.RollingUpdate.Partition != nil {
			if s.Status.UpdatedReplicas < (*s.Spec.Replicas - *s.Spec.UpdateStrategy.RollingUpdate.Partition) {
				desc = fmt.Sprintf("Waiting for partitioned roll out to finish because %d out of %d new pods have been updated",
					s.Status.UpdatedReplicas, (*s.Spec.Replicas - *s.Spec.UpdateStrategy.RollingUpdate.Partition))
				return
			}
		}
		status = model.KubernetesResourceState_HEALTHY
		return
	}

	if s.Status.UpdateRevision != s.Status.CurrentRevision {
		desc = fmt.Sprintf("Waiting for statefulset rolling update to complete %d pods at revision %s", s.Status.UpdatedReplicas, s.Status.UpdateRevision)
		return
	}

	status = model.KubernetesResourceState_HEALTHY
	return
}

func determineDaemonSetHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	d := &appsv1.DaemonSet{}
	err := scheme.Scheme.Convert(obj, d, nil)
	if err != nil {
		desc = fmt.Sprintf("Failed while convert %T to %T: %v", obj, d, err)
		return
	}

	status = model.KubernetesResourceState_OTHER
	if d.Status.NumberMisscheduled > 0 {
		desc = fmt.Sprintf("%d nodes that are running the daemon pod, but are not supposed to run the daemon pod", d.Status.NumberMisscheduled)
		return
	}
	if d.Status.NumberUnavailable > 0 {
		desc = fmt.Sprintf("%d nodes that should be running the daemon pod and have none of the daemon pod running and available", d.Status.NumberUnavailable)
		return
	}

	status = model.KubernetesResourceState_HEALTHY
	return
}

func determineReplicaSetHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	r := &appsv1.ReplicaSet{}
	err := scheme.Scheme.Convert(obj, r, nil)
	if err != nil {
		desc = fmt.Sprintf("Failed while convert %T to %T: %v", obj, r, err)
		return
	}

	status = model.KubernetesResourceState_OTHER
	if r.Status.ObservedGeneration == 0 || r.Generation > r.Status.ObservedGeneration {
		desc = "Waiting for rollout to finish because observed replica set generation less then desired generation"
		return
	}

	var cond *appsv1.ReplicaSetCondition
	for i := range r.Status.Conditions {
		c := r.Status.Conditions[i]
		if c.Type == appsv1.ReplicaSetReplicaFailure {
			cond = &c
			break
		}
	}
	if cond != nil && cond.Status == corev1.ConditionTrue {
		desc = cond.Message
		return
	} else if r.Spec.Replicas == nil {
		desc = "The number of desired replicas is unspecified"
		return
	} else if r.Status.AvailableReplicas < *r.Spec.Replicas {
		desc = fmt.Sprintf("Waiting for rollout to finish because only %d/%d replicas are available", r.Status.AvailableReplicas, *r.Spec.Replicas)
		return
	} else if *r.Spec.Replicas != r.Status.ReadyReplicas {
		desc = fmt.Sprintf("The number of ready replicas (%d) is different from the desired number (%d)", r.Status.ReadyReplicas, *r.Spec.Replicas)
		return
	}

	status = model.KubernetesResourceState_HEALTHY
	return
}

func determinePodHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	p := &corev1.Pod{}
	err := scheme.Scheme.Convert(obj, p, nil)
	if err != nil {
		desc = fmt.Sprintf("Failed while convert %T to %T: %v", obj, p, err)
		return
	}

	if p.Status.Phase == corev1.PodRunning {
		status = model.KubernetesResourceState_HEALTHY
	} else {
		status = model.KubernetesResourceState_OTHER
	}
	desc = p.Status.Message
	return
}

func determineIngressHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	i := &networkingv1beta1.Ingress{}
	err := scheme.Scheme.Convert(obj, i, nil)
	if err != nil {
		desc = fmt.Sprintf("Failed while convert %T to %T: %v", obj, i, err)
		return
	}

	status = model.KubernetesResourceState_OTHER
	if len(i.Status.LoadBalancer.Ingress) <= 0 {
		desc = "Ingress points for the load-balancer are in progress"
		return
	}
	status = model.KubernetesResourceState_HEALTHY
	return
}

func determineServiceHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	s := &corev1.Service{}
	err := scheme.Scheme.Convert(obj, s, nil)
	if err != nil {
		desc = fmt.Sprintf("Failed while convert %T to %T: %v", obj, s, err)
		return
	}

	status = model.KubernetesResourceState_HEALTHY
	if s.Spec.Type != corev1.ServiceTypeLoadBalancer {
		return
	}
	if len(s.Status.LoadBalancer.Ingress) <= 0 {
		status = model.KubernetesResourceState_OTHER
		desc = "Ingress points for the load-balancer are in progress"
		return
	}
	return
}
