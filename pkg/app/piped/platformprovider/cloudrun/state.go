// Copyright 2023 The PipeCD Authors.
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

package cloudrun

import (
	"sort"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func MakeResourceStates(svc *Service, revs []*Revision, updatedAt time.Time) []*model.CloudRunResourceState {
	states := make([]*model.CloudRunResourceState, 0, len(revs)+1)

	// Set service state.
	sm, err := svc.ServiceManifest()
	if err == nil {
		status, desc := svc.StatusConditions().HealthStatus()
		states = append(states, makeResourceState(sm.u, status, desc, updatedAt))
	}

	// Set active revision states.
	for _, r := range revs {
		rm, err := r.RevisionManifest()
		if err != nil {
			continue
		}

		status, desc := r.StatusConditions().HealthStatus()
		states = append(states, makeResourceState(rm.u, status, desc, updatedAt))
	}
	return states
}

func makeResourceState(obj *unstructured.Unstructured, status model.CloudRunResourceState_HealthStatus, desc string, updatedAt time.Time) *model.CloudRunResourceState {
	var (
		owners       = obj.GetOwnerReferences()
		ownerIDs     = make([]string, 0, len(owners))
		creationTime = obj.GetCreationTimestamp()
	)

	for _, owner := range owners {
		ownerIDs = append(ownerIDs, string(owner.UID))
	}
	sort.Strings(ownerIDs)

	state := &model.CloudRunResourceState{
		Id:         string(obj.GetUID()),
		OwnerIds:   ownerIDs,
		ParentIds:  ownerIDs,
		Name:       obj.GetName(),
		ApiVersion: obj.GetAPIVersion(),
		Kind:       obj.GetKind(),
		Namespace:  obj.GetNamespace(),

		HealthStatus:      status,
		HealthDescription: desc,

		CreatedAt: creationTime.Unix(),
		UpdatedAt: updatedAt.Unix(),
	}

	return state
}
