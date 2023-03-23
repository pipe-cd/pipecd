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

package applicationlivestatestore

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestMergeKubernetesResourceStatesOnAddOrUpdated(t *testing.T) {
	testcases := []struct {
		name       string
		prevStates []*model.KubernetesResourceState
		event      *model.KubernetesResourceStateEvent

		expectedStetes []*model.KubernetesResourceState
	}{
		{
			name: "event.State was not found in prevStates",
			prevStates: []*model.KubernetesResourceState{
				{
					Id:           "resource-01",
					ApiVersion:   "v1",
					Name:         "unit-test",
					Kind:         "Service",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
				{
					Id:           "resource-02",
					ApiVersion:   "apps/v1",
					Name:         "unit-test",
					Kind:         "Deployment",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
				{
					Id:           "resource-03",
					ApiVersion:   "v1",
					Name:         "unit-test",
					Kind:         "ConfigMap",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
			},
			event: &model.KubernetesResourceStateEvent{
				Id:            "event-id",
				ApplicationId: "application-id",
				State: &model.KubernetesResourceState{
					Id:           "resource-04",
					ApiVersion:   "batch/v1",
					Name:         "unit-test",
					Kind:         "Job",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
			},

			expectedStetes: []*model.KubernetesResourceState{
				{
					Id:           "resource-01",
					ApiVersion:   "v1",
					Name:         "unit-test",
					Kind:         "Service",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
				{
					Id:           "resource-02",
					ApiVersion:   "apps/v1",
					Name:         "unit-test",
					Kind:         "Deployment",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
				{
					Id:           "resource-03",
					ApiVersion:   "v1",
					Name:         "unit-test",
					Kind:         "ConfigMap",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
				{
					Id:           "resource-04",
					ApiVersion:   "batch/v1",
					Name:         "unit-test",
					Kind:         "Job",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
			},
		},

		{
			name: "event.State was found in prevStates",
			prevStates: []*model.KubernetesResourceState{
				{
					Id:           "resource-01",
					ApiVersion:   "v1",
					Name:         "unit-test",
					Kind:         "Service",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
				{
					Id:           "resource-02",
					ApiVersion:   "apps/v1",
					Name:         "unit-test",
					Kind:         "Deployment",
					HealthStatus: model.KubernetesResourceState_UNKNOWN,
				},
				{
					Id:           "resource-03",
					ApiVersion:   "v1",
					Name:         "unit-test",
					Kind:         "ConfigMap",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
			},
			event: &model.KubernetesResourceStateEvent{
				Id:            "event-id",
				ApplicationId: "application-id",
				State: &model.KubernetesResourceState{
					Id:           "resource-02",
					ApiVersion:   "apps/v1",
					Name:         "unit-test",
					Kind:         "Deployment",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
			},

			expectedStetes: []*model.KubernetesResourceState{
				{
					Id:           "resource-01",
					ApiVersion:   "v1",
					Name:         "unit-test",
					Kind:         "Service",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
				{
					Id:           "resource-02",
					ApiVersion:   "apps/v1",
					Name:         "unit-test",
					Kind:         "Deployment",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
				{
					Id:           "resource-03",
					ApiVersion:   "v1",
					Name:         "unit-test",
					Kind:         "ConfigMap",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			states := mergeKubernetesResourceStatesOnAddOrUpdated(tc.prevStates, tc.event)
			assert.Equal(t, tc.expectedStetes, states)
		})
	}
}

func TestMergeKubernetesResourceStatesOnDeleted(t *testing.T) {
	testcases := []struct {
		name       string
		prevStates []*model.KubernetesResourceState
		event      *model.KubernetesResourceStateEvent

		expectedStetes []*model.KubernetesResourceState
	}{
		{
			name: "event.State was not found in prevStates",
			prevStates: []*model.KubernetesResourceState{
				{
					Id:           "resource-01",
					ApiVersion:   "v1",
					Name:         "unit-test",
					Kind:         "Service",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
				{
					Id:           "resource-02",
					ApiVersion:   "apps/v1",
					Name:         "unit-test",
					Kind:         "Deployment",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
				{
					Id:           "resource-03",
					ApiVersion:   "v1",
					Name:         "unit-test",
					Kind:         "ConfigMap",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
			},
			event: &model.KubernetesResourceStateEvent{
				Id:            "event-id",
				ApplicationId: "application-id",
				State: &model.KubernetesResourceState{
					Id:           "resource-99",
					ApiVersion:   "batch/v1",
					Name:         "unit-test",
					Kind:         "Job",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
			},

			expectedStetes: []*model.KubernetesResourceState{
				{
					Id:           "resource-01",
					ApiVersion:   "v1",
					Name:         "unit-test",
					Kind:         "Service",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
				{
					Id:           "resource-02",
					ApiVersion:   "apps/v1",
					Name:         "unit-test",
					Kind:         "Deployment",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
				{
					Id:           "resource-03",
					ApiVersion:   "v1",
					Name:         "unit-test",
					Kind:         "ConfigMap",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
			},
		},

		{
			name: "event.State was found in prevStates",
			prevStates: []*model.KubernetesResourceState{
				{
					Id:           "resource-01",
					ApiVersion:   "v1",
					Name:         "unit-test",
					Kind:         "Service",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
				{
					Id:           "resource-02",
					ApiVersion:   "apps/v1",
					Name:         "unit-test",
					Kind:         "Deployment",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
				{
					Id:           "resource-03",
					ApiVersion:   "v1",
					Name:         "unit-test",
					Kind:         "ConfigMap",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
			},
			event: &model.KubernetesResourceStateEvent{
				Id:            "event-id",
				ApplicationId: "application-id",
				State: &model.KubernetesResourceState{
					Id:           "resource-03",
					ApiVersion:   "v1",
					Name:         "unit-test",
					Kind:         "ConfigMap",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
			},

			expectedStetes: []*model.KubernetesResourceState{
				{
					Id:           "resource-01",
					ApiVersion:   "v1",
					Name:         "unit-test",
					Kind:         "Service",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
				{
					Id:           "resource-02",
					ApiVersion:   "apps/v1",
					Name:         "unit-test",
					Kind:         "Deployment",
					HealthStatus: model.KubernetesResourceState_HEALTHY,
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			states := mergeKubernetesResourceStatesOnDeleted(tc.prevStates, tc.event)
			assert.Equal(t, tc.expectedStetes, states)
		})
	}
}
