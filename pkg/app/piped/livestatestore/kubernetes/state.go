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
	"sort"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	provider "github.com/kapetaniosci/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/kapetaniosci/pipe/pkg/model"
)

func makeKubernetesResourceState(uid string, key provider.ResourceKey, obj *unstructured.Unstructured, now time.Time) model.KubernetesResourceState {
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

func determineResourceHealth(key provider.ResourceKey, obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	// TODO: Decide health status for each resource.
	return
}
