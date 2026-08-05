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

package applicationlivestatestore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/cache/memorycache"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// TestPatchKubernetesApplicationLiveState_OutOfOrderEvents proves that a single
// call carrying two events for the same application, where an older event
// happens to come after a newer one (e.g. because it was redelivered by a
// retry), does not let the older event overwrite the state the newer one
// already applied.
//
// This can genuinely happen: the piped-side reporter batches events and
// resends unacknowledged ones, so a duplicate/stale event landing after a
// fresher one for the same resource is not a made-up scenario.
func TestPatchKubernetesApplicationLiveState_OutOfOrderEvents(t *testing.T) {
	s := &store{
		cache: &applicationLiveStateCache{backend: memorycache.NewCache()},
		logger: zap.NewNop(),
	}

	const appID = "app-1"
	baseVersion := model.ApplicationLiveStateVersion{Timestamp: 100, Index: 0}
	newerVersion := model.ApplicationLiveStateVersion{Timestamp: 100, Index: 1}

	initial := &model.ApplicationLiveStateSnapshot{
		ApplicationId: appID,
		Kind:          model.ApplicationKind_KUBERNETES,
		Version:       &baseVersion,
		Kubernetes: &model.KubernetesApplicationLiveState{
			Resources: []*model.KubernetesResourceState{
				{Id: "pod-1", HealthStatus: model.KubernetesResourceState_UNKNOWN},
			},
		},
	}
	require.NoError(t, s.cache.Put(appID, initial))

	// The true, later event: the pod recovered. Its SnapshotVersion is the
	// version the app was at right before this change (newerVersion),
	// i.e. it is chronologically the second change.
	recoveredEvent := &model.KubernetesResourceStateEvent{
		Id:              "event-recovered",
		ApplicationId:   appID,
		Type:            model.KubernetesResourceStateEvent_ADD_OR_UPDATED,
		SnapshotVersion: &newerVersion,
		State:           &model.KubernetesResourceState{Id: "pod-1", HealthStatus: model.KubernetesResourceState_HEALTHY},
	}
	// A stale redelivery of the first change (the crash), arriving in the
	// same batch AFTER the recovery event above.
	staleCrashEvent := &model.KubernetesResourceStateEvent{
		Id:              "event-crash-retry",
		ApplicationId:   appID,
		Type:            model.KubernetesResourceStateEvent_ADD_OR_UPDATED,
		SnapshotVersion: &baseVersion,
		State:           &model.KubernetesResourceState{Id: "pod-1", HealthStatus: model.KubernetesResourceState_OTHER},
	}

	s.PatchKubernetesApplicationLiveState(context.Background(), []*model.KubernetesResourceStateEvent{
		recoveredEvent,
		staleCrashEvent,
	})

	got, err := s.GetStateSnapshot(context.Background(), appID)
	require.NoError(t, err)
	require.Len(t, got.Kubernetes.Resources, 1)
	assert.Equal(t, model.KubernetesResourceState_HEALTHY, got.Kubernetes.Resources[0].HealthStatus,
		"the stale retried event must not overwrite the newer state that was already applied in this same batch")
}

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
