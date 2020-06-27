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
	desc = "Not implemented yet"
	return
}

func determineDaemonSetHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	desc = "Not implemented yet"
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
	if r.Spec.Replicas == nil {
		desc = "The number of desired replicas is unspecified"
		return
	}
	if *r.Spec.Replicas != r.Status.ReadyReplicas {
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
	desc = "Not implemented yet"
	return
}

func determineServiceHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	desc = "Not implemented yet"
	return
}
